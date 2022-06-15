package digitalocean

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceDigitalOceanReservedIP() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanReservedIPRead,
		Schema: map[string]*schema.Schema{
			"ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "reserved ip address",
				ValidateFunc: validation.NoZeroValues,
			},
			// computed attributes
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the uniform resource name for the reserved ip",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the region that the reserved ip is reserved to",
			},
			"droplet_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the droplet id that the reserved ip has been assigned to.",
			},
		},
	}
}

func dataSourceDigitalOceanReservedIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ipAddress := d.Get("ip_address").(string)
	d.SetId(ipAddress)

	return resourceDigitalOceanReservedIPRead(ctx, d, meta)
}
