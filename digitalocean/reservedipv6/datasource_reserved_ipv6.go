package reservedipv6

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDigitalOceanReservedIPV6() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanReservedIPV6Read,
		Schema: map[string]*schema.Schema{
			"ip": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "reserved ipv6 address",
				ValidateFunc: validation.NoZeroValues,
			},
			// computed attributes
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the uniform resource name for the reserved ipv6",
			},
			"region_slug": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the regionSlug that the reserved ipv6 is reserved to",
			},
			"droplet_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the droplet id that the reserved ipv6 has been assigned to.",
			},
		},
	}
}

func dataSourceDigitalOceanReservedIPV6Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	ipAddress := d.Get("ip").(string)
	d.SetId(ipAddress)

	return resourceDigitalOceanReservedIPV6Read(ctx, d, meta)
}
