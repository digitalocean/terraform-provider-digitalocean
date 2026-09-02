package microdroplet

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceDigitalOceanMicroDropletImage returns the
// digitalocean_microdroplet_image resource schema. Images are immutable —
// there is no UpdateContext. All fields are ForceNew.
func ResourceDigitalOceanMicroDropletImage() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanMicroDropletImageCreate,
		ReadContext:   resourceDigitalOceanMicroDropletImageRead,
		DeleteContext: resourceDigitalOceanMicroDropletImageDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: microDropletImageResourceSchema(),
	}
}

func resourceDigitalOceanMicroDropletImageCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	req := &godo.MicroDropletImageCreateRequest{
		Name:   d.Get("name").(string),
		Source: d.Get("source").(string),
	}

	log.Printf("[DEBUG] MicroDroplet image create request: %+v", req)
	image, _, err := client.MicroDropletImages.Create(ctx, req)
	if err != nil {
		return diag.Errorf("Error creating MicroDroplet image: %s", err)
	}

	d.SetId(image.ID)
	log.Printf("[INFO] MicroDroplet image created, ID: %s", d.Id())

	if _, err := waitForMicroDropletImage(ctx, client, d.Id(), godo.MicroDropletImageStatusAvailable, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("Error waiting for MicroDroplet image (%s) to become available: %s", d.Id(), err)
	}

	return resourceDigitalOceanMicroDropletImageRead(ctx, d, meta)
}

func resourceDigitalOceanMicroDropletImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	image, resp, err := client.MicroDropletImages.Get(ctx, d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			log.Printf("[WARN] MicroDroplet image (%s) not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error retrieving MicroDroplet image: %s", err)
	}

	if image.Status == godo.MicroDropletImageStatusDeleted {
		d.SetId("")
		return nil
	}

	d.SetId(image.ID)
	d.Set("name", image.Name)
	d.Set("source", image.Source)
	d.Set("status", string(image.Status))
	d.Set("created_at", image.Created)
	d.Set("urn", image.URN())
	return nil
}

func resourceDigitalOceanMicroDropletImageDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	resp, err := client.MicroDropletImages.Delete(ctx, d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return diag.Errorf("Error deleting MicroDroplet image: %s", err)
	}

	log.Printf("[INFO] MicroDroplet image deleted, ID: %s", d.Id())
	d.SetId("")
	return nil
}

// waitForMicroDropletImage polls MicroDropletImages.Get until the observed
// status matches target, or the timeout elapses.
func waitForMicroDropletImage(ctx context.Context, client *godo.Client, id string, target godo.MicroDropletImageStatus, timeout time.Duration) (interface{}, error) {
	log.Printf("[INFO] Waiting for MicroDroplet image (%s) to have status %s", id, target)

	stateConf := &retry.StateChangeConf{
		Pending: []string{
			string(godo.MicroDropletImageStatusUnknown),
			string(godo.MicroDropletImageStatusImporting),
		},
		Target:     []string{string(target)},
		Refresh:    microDropletImageStateRefreshFunc(ctx, client, id),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 10 * time.Second,
	}
	return stateConf.WaitForStateContext(ctx)
}

func microDropletImageStateRefreshFunc(ctx context.Context, client *godo.Client, id string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		image, _, err := client.MicroDropletImages.Get(ctx, id)
		if err != nil {
			return nil, "", err
		}
		if image.Status == godo.MicroDropletImageStatusFailed {
			return image, string(image.Status), fmt.Errorf("MicroDroplet image %s entered failed state", id)
		}
		return image, string(image.Status), nil
	}
}
