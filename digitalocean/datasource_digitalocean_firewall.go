package digitalocean

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDigitalOceanFirewall() *schema.Resource {
	fwSchema := firewallSchema()

	for _, f := range fwSchema {
		f.Computed = true
		f.Required = false
	}

	fwSchema["name"].ValidateFunc = nil

	fwSchema["firewall_id"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
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
