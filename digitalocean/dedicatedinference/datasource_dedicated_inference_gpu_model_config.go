package dedicatedinference

import (
	"context"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanDedicatedInferenceGPUModelConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanDedicatedInferenceGPUModelConfigRead,
		Schema: map[string]*schema.Schema{
			"gpu_model_configs": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of supported GPU and model combinations for dedicated inference endpoints.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gpu_slugs": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The GPU slugs that support this model.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"model_slug": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The slug identifier for the model.",
						},
						"model_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The human-readable name of the model.",
						},
						"is_model_gated": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the model requires gated access (e.g. a HuggingFace token).",
						},
					},
				},
			},
		},
	}
}

func dataSourceDigitalOceanDedicatedInferenceGPUModelConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	resp, _, err := client.DedicatedInference.GetGPUModelConfig(ctx)
	if err != nil {
		return diag.Errorf("error retrieving dedicated inference GPU model configs: %s", err)
	}

	d.SetId("dedicated_inference_gpu_model_config")

	flatConfigs := make([]map[string]interface{}, 0, len(resp.GPUModelConfigs))
	for _, c := range resp.GPUModelConfigs {
		flatConfigs = append(flatConfigs, map[string]interface{}{
			"gpu_slugs":      c.GPUSlugs,
			"model_slug":     c.ModelSlug,
			"model_name":     c.ModelName,
			"is_model_gated": c.IsModelGated,
		})
	}

	if err := d.Set("gpu_model_configs", flatConfigs); err != nil {
		return diag.Errorf("error setting gpu_model_configs: %s", err)
	}

	return nil
}
