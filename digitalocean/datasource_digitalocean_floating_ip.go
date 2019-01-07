package digitalocean

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceDigitalOceanFloatingIp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanFloatingIpRead,
		Schema: map[string]*schema.Schema{

			"ip_address": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "floating ip address",
				ValidateFunc: validation.NoZeroValues,
			},
			// computed attributes
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the region that the floating ip is reserved to",
			},
			"droplet_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the droplet id that the floating ip has been assigned to.",
			},
		},
	}
}

func dataSourceDigitalOceanFloatingIpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	ipAddress := d.Get("ip_address").(string)

	floatingIp, resp, err := client.FloatingIPs.Get(context.Background(), ipAddress)
	if err != nil {
		if resp.StatusCode == 404 {
			return fmt.Errorf("floating ip not found: %s", err)
		}
		return fmt.Errorf("Error retrieving floating ip: %s", err)
	}

	d.SetId(floatingIp.IP)
	d.Set("ip_address", floatingIp.IP)
	d.Set("region", floatingIp.Region.Slug)

	if floatingIp.Droplet != nil {
		d.Set("droplet_id", floatingIp.Droplet.ID)
	}

	return nil
}
