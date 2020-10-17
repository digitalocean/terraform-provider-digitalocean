package digitalocean

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceDigitalOceanFloatingIp() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanFloatingIpRead,
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

func dataSourceDigitalOceanFloatingIpRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	ipAddress := d.Get("ip_address").(string)

	floatingIp, resp, err := client.FloatingIPs.Get(context.Background(), ipAddress)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return diag.Errorf("floating ip not found: %s", err)
		}
		return diag.Errorf("Error retrieving floating ip: %s", err)
	}

	d.SetId(floatingIp.IP)
	d.Set("ip_address", floatingIp.IP)
	d.Set("urn", floatingIp.URN())
	d.Set("region", floatingIp.Region.Slug)

	if floatingIp.Droplet != nil {
		d.Set("droplet_id", floatingIp.Droplet.ID)
	}

	return nil
}
