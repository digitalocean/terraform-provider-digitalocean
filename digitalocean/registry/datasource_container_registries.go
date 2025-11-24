package registry

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDigitalOceanContainerRegistries() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanContainerRegistriesRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "name of the container registries",
				ValidateFunc: validation.NoZeroValues,
			},
			"subscription_tier_slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_usage_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDigitalOceanContainerRegistriesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceDigitalOceanContainerRegistriesRead(ctx, d, meta)
}
