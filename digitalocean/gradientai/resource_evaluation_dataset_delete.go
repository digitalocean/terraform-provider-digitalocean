package gradientai

import (
	"context"
	"fmt"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceDigitalOceanEvaluationDatasetDelete defines the DigitalOcean GradientAI
// Evaluation Dataset Delete resource.
func ResourceDigitalOceanEvaluationDatasetDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanEvaluationDatasetDeleteCreate,
		ReadContext:   resourceDigitalOceanEvaluationDatasetDeleteRead,
		DeleteContext: resourceDigitalOceanEvaluationDatasetDeleteDelete,
		// No update context since this is a one-time action.

		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The UUID of the evaluation dataset to delete.",
			},
		},
	}
}

func resourceDigitalOceanEvaluationDatasetDeleteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	uuid := d.Get("uuid").(string)

	if _, _, err := client.GradientAI.DeleteEvaluationDataset(ctx, uuid); err != nil {
		return diag.FromErr(fmt.Errorf("error deleting evaluation dataset (%s): %s", uuid, err))
	}

	d.SetId(uuid)

	return nil
}

func resourceDigitalOceanEvaluationDatasetDeleteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// Deleting an evaluation dataset is a one-time action; the dataset no longer
	// exists after creation of this resource, so there is nothing to refresh.
	return nil
}

func resourceDigitalOceanEvaluationDatasetDeleteDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// The deletion already happened on create and is permanent, so removing this
	// resource from Terraform only drops it from state.
	d.SetId("")
	return nil
}
