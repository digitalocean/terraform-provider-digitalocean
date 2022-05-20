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
	client := meta.(*CombinedConfig).godoClient()

	ipAddress := d.Get("ip_address").(string)

	reservedIP, resp, err := client.ReservedIPs.Get(context.Background(), ipAddress)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return diag.Errorf("reserved ip not found: %s", err)
		}
		return diag.Errorf("Error retrieving reserved ip: %s", err)
	}

	d.SetId(reservedIP.IP)
	d.Set("ip_address", reservedIP.IP)
	d.Set("urn", reservedIP.URN())
	d.Set("region", reservedIP.Region.Slug)

	if reservedIP.Droplet != nil {
		d.Set("droplet_id", reservedIP.Droplet.ID)
	}

	return nil
}
