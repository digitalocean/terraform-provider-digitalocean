package digitalocean

import (
	"context"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceDigitalOceanFloatingIP() *schema.Resource {
	return &schema.Resource{
		// TODO: Uncomment when dates for deprecation timeline are set.
		// DeprecationMessage: "This data source is deprecated and will be removed in a future release. Please use digitalocean_reserved_ip instead.",
		ReadContext: dataSourceDigitalOceanFloatingIPRead,
		Schema: map[string]*schema.Schema{

			"ip_address": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "floating ip address",
				ValidateFunc: validation.NoZeroValues,
			},
			// computed attributes
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the uniform resource name for the floating ip",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the region that the floating ip is reserved to",
			},
			"droplet_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the droplet id that the floating ip has been assigned to.",
			},
		},
	}
}

func dataSourceDigitalOceanFloatingIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	dataSourceDigitalOceanReservedIPRead(ctx, d, meta)
	// Re-format reserved IP URN as floating IP URN
	// TODO: Remove when the projects' API changes return values.
	ip := d.Get("ip_address")
	d.Set("urn", godo.FloatingIP{IP: ip.(string)}.URN())

	return nil
}
