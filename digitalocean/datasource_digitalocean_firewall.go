package digitalocean

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDigitalOceanFirewall() *schema.Resource {
	fwSchema := firewallSchema()

	for _, f := range fwSchema {
		if !f.Required {
			f.Computed = true
		}
	}

	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanFirewallRead,
		Schema:      fwSchema,
	}
}

func dataSourceDigitalOceanFirewallRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("firewall_id").(string))
	return resourceDigitalOceanFirewallRead(ctx, d, meta)
}
