package digitalocean

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDigitalOceanApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanAppCreate,
		Read:   resourceDigitalOceanAppRead,
		Update: resourceDigitalOceanAppUpdate,
		Delete: resourceDigitalOceanAppDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"spec": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: "A DigitalOcean App Platform Spec",
				Elem: &schema.Resource{
					Schema: appSpecSchema(),
				},
			},

			// Computed attributes
			"default_ingress": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The default URL to access the App",
			},

			"live_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			// TODO: The full Deployment should be a data source, not a resource
			// specify the app id for the active deployment, include a deployment
			// id for a specific one
			"active_deployment_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID the App's currently active deployment",
			},

			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of when the App was last updated",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of when the App was created",
			},
		},
	}
}

func resourceDigitalOceanAppCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
	appCreateRequest := &godo.AppCreateRequest{}
	appCreateRequest.Spec = expandAppSpec(d.Get("spec").([]interface{}))

	log.Printf("[DEBUG] App create request: %#v", appCreateRequest)
	app, _, err := client.Apps.Create(context.Background(), appCreateRequest)
	if err != nil {
		return fmt.Errorf("Error creating App: %s", err)
	}

	d.SetId(app.ID)
	log.Printf("[DEBUG] Waiting for app (%s) deployment to become active", app.ID)

	err = waitForAppDeployment(client, app.ID)
	if err != nil {
		return err
	}

	log.Printf("[INFO] App created, ID: %s", d.Id())

	return resourceDigitalOceanAppRead(d, meta)
}

func resourceDigitalOceanAppRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	app, resp, err := client.Apps.Get(context.Background(), d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[DEBUG] App (%s) was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading App: %s", err)
	}

	d.SetId(app.ID)
	d.Set("default_ingress", app.DefaultIngress)
	d.Set("live_url", app.LiveURL)
	d.Set("active_deployment_id", app.ActiveDeployment.ID)
	d.Set("updated_at", app.UpdatedAt.UTC().String())
	d.Set("created_at", app.CreatedAt.UTC().String())

	if err := d.Set("spec", flattenAppSpec(app.Spec)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting app spec: %#v", err)
	}

	return err
}

func resourceDigitalOceanAppUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	if d.HasChange("spec") {
		appUpdateRequest := &godo.AppUpdateRequest{}
		appUpdateRequest.Spec = expandAppSpec(d.Get("spec").([]interface{}))

		app, _, err := client.Apps.Update(context.Background(), d.Id(), appUpdateRequest)
		if err != nil {
			return fmt.Errorf("Error updating app (%s): %s", d.Id(), err)
		}

		log.Printf("[DEBUG] Waiting for app (%s) deployment to become active", app.ID)
		err = waitForAppDeployment(client, app.ID)
		if err != nil {
			return err
		}

		log.Printf("[INFO] Updated app (%s)", app.ID)
	}

	return resourceDigitalOceanAppRead(d, meta)
}

func resourceDigitalOceanAppDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	log.Printf("[INFO] Deleting App: %s", d.Id())
	_, err := client.Apps.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error deletingApp: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForAppDeployment(client *godo.Client, id string) error {
	tickerInterval := 10 //10s
	timeout := 1800      //1800s, 30min
	n := 0

	// No deployment seems to be available right way, sleep 10 seconds
	time.Sleep(10 * time.Second)
	app, _, err := client.Apps.Get(context.Background(), id)
	if err != nil {
		return fmt.Errorf("Error trying to read app deployment state: %s", err)
	}

	ticker := time.NewTicker(time.Duration(tickerInterval) * time.Second)
	for range ticker.C {
		deployment, _, err := client.Apps.GetDeployment(context.Background(), id, app.InProgressDeployment.ID)
		if err != nil {
			ticker.Stop()
			return fmt.Errorf("Error trying to read app deployment state: %s", err)
		}

		allSuccessful := deployment.Progress.SuccessSteps == deployment.Progress.TotalSteps
		if allSuccessful {
			ticker.Stop()
			return nil
		}

		if deployment.Progress.ErrorSteps > 0 {
			ticker.Stop()
			return fmt.Errorf("error deploying app (%s) (deployment ID: %s):\n%s", id, app.InProgressDeployment.ID, godo.Stringify(deployment.Progress))
		}

		if n*tickerInterval > timeout {
			ticker.Stop()
			break
		}

		log.Printf("[DEBUG] Waiting for app (%s) deployment (%s) to become active (%d/%d)", app.ID, app.InProgressDeployment.ID, deployment.Progress.SuccessSteps, deployment.Progress.TotalSteps)

		n++
	}

	return fmt.Errorf("timeout waiting to app (%s) deployment", app.ID)
}
