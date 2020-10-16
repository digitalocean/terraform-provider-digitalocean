package digitalocean

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDigitalOceanProjectResources() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanProjectResourcesUpdate,
		Update: resourceDigitalOceanProjectResourcesUpdate,
		Read:   resourceDigitalOceanProjectResourcesRead,
		Delete: resourceDigitalOceanProjectResourcesDelete,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "project ID",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"resources": {
				Type:        schema.TypeSet,
				Required:    true,
				Description: "the resources associated with the project",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceDigitalOceanProjectResourcesUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	projectId := d.Get("project").(string)

	_, resp, err := client.Projects.Get(context.Background(), projectId)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			// Project does not exist. Mark this resource as not existing.
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error while retrieving project %s: %v", projectId, err)
	}

	if d.HasChange("resources") {
		oldURNs, newURNs := d.GetChange("resources")

		if oldURNs.(*schema.Set).Len() > 0 {
			_, err = assignResourcesToDefaultProject(client, oldURNs.(*schema.Set))
			if err != nil {
				return fmt.Errorf("Error assigning resources to default project: %s", err)
			}
		}

		var urns *[]interface{}

		if newURNs.(*schema.Set).Len() > 0 {
			urns, err = assignResourcesToProject(client, projectId, newURNs.(*schema.Set))
			if err != nil {
				return fmt.Errorf("Error assigning resources to project %s: %s", projectId, err)
			}
		}

		if err = d.Set("resources", urns); err != nil {
			return err
		}
	}

	d.SetId(projectId)

	return resourceDigitalOceanProjectResourcesRead(d, meta)
}

func resourceDigitalOceanProjectResourcesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	projectId := d.Id()

	_, resp, err := client.Projects.Get(context.Background(), projectId)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			// Project does not exist. Mark this resource as not existing.
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error while retrieving project: %v", err)
	}

	if err = d.Set("project", projectId); err != nil {
		return err
	}

	apiURNs, err := loadResourceURNs(client, projectId)
	if err != nil {
		return fmt.Errorf("Error while retrieving project resources: %s", err)
	}

	var newURNs []string

	configuredURNs := d.Get("resources").(*schema.Set).List()
	for _, rawConfiguredURN := range configuredURNs {
		configuredURN := rawConfiguredURN.(string)

		for _, apiURN := range *apiURNs {
			if configuredURN == apiURN {
				newURNs = append(newURNs, configuredURN)
			}
		}
	}

	if err = d.Set("resources", newURNs); err != nil {
		return err
	}

	return nil
}

func resourceDigitalOceanProjectResourcesDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	projectId := d.Get("project").(string)
	urns := d.Get("resources").(*schema.Set)

	_, resp, err := client.Projects.Get(context.Background(), projectId)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			// Project does not exist. Mark this resource as not existing.
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error while retrieving project: %s", err)
	}

	if urns.Len() > 0 {
		if _, err = assignResourcesToDefaultProject(client, urns); err != nil {
			return fmt.Errorf("Error assigning resources to default project: %s", err)
		}
	}

	d.SetId("")
	return nil
}
