package dedicatedinference

import (
	"context"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanDedicatedInferenceSizes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanDedicatedInferenceSizesRead,
		Schema: map[string]*schema.Schema{
			"enabled_regions": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of regions where dedicated inference endpoints can be deployed.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"sizes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of available GPU sizes with their configurations and pricing.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"gpu_slug": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The slug identifier for this GPU size.",
						},
						"price_per_hour": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The hourly price for this GPU size.",
						},
						"regions": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The regions where this GPU size is available.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"currency": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The currency for the price.",
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of vCPUs.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The amount of memory in MiB.",
						},
						"gpu": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "GPU hardware details for this size.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"count": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of GPUs.",
									},
									"vram_gb": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The VRAM per GPU in GiB.",
									},
									"slug": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The GPU model slug.",
									},
								},
							},
						},
						"size_category": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The category this size belongs to.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The display name of the size category.",
									},
									"fleet_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The fleet name associated with the size category.",
									},
								},
							},
						},
						"disks": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The disks attached to this size.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The disk type.",
									},
									"size_gb": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The disk size in GiB.",
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceDigitalOceanDedicatedInferenceSizesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	sizesResp, _, err := client.DedicatedInference.GetSizes(ctx)
	if err != nil {
		return diag.Errorf("error retrieving dedicated inference sizes: %s", err)
	}

	d.SetId("dedicated_inference_sizes")

	if err := d.Set("enabled_regions", sizesResp.EnabledRegions); err != nil {
		return diag.Errorf("error setting enabled_regions: %s", err)
	}

	flatSizes := make([]map[string]interface{}, 0, len(sizesResp.Sizes))
	for _, s := range sizesResp.Sizes {
		flat := map[string]interface{}{
			"gpu_slug":       s.GPUSlug,
			"price_per_hour": s.PricePerHour,
			"currency":       s.Currency,
			"cpu":            int(s.CPU),
			"memory":         int(s.Memory),
			"regions":        s.Regions,
		}

		if s.GPU != nil {
			flat["gpu"] = []map[string]interface{}{
				{
					"count":   int(s.GPU.Count),
					"vram_gb": int(s.GPU.VramGb),
					"slug":    s.GPU.Slug,
				},
			}
		} else {
			flat["gpu"] = []map[string]interface{}{}
		}

		if s.SizeCategory != nil {
			flat["size_category"] = []map[string]interface{}{
				{
					"name":       s.SizeCategory.Name,
					"fleet_name": s.SizeCategory.FleetName,
				},
			}
		} else {
			flat["size_category"] = []map[string]interface{}{}
		}

		disks := make([]map[string]interface{}, 0, len(s.Disks))
		for _, disk := range s.Disks {
			disks = append(disks, map[string]interface{}{
				"type":    disk.Type,
				"size_gb": int(disk.SizeGb),
			})
		}
		flat["disks"] = disks

		flatSizes = append(flatSizes, flat)
	}

	if err := d.Set("sizes", flatSizes); err != nil {
		return diag.Errorf("error setting sizes: %s", err)
	}

	return nil
}
