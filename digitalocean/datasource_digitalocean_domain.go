package digitalocean

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceDigitalOceanDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanDomainRead,
		Schema: map[string]*schema.Schema{

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "name of the domain",
				ValidateFunc: validation.NoZeroValues,
			},
			// computed attributes
			"ttl": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ttl of the domain",
			},
			"zone_file": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "zone file of the domain",
			},
		},
	}
}

func dataSourceDigitalOceanDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	name := d.Get("name").(string)

	domain, resp, err := client.Domains.Get(context.Background(), name)
	if err != nil {
		if resp.StatusCode == 404 {
			return fmt.Errorf("domain not found: %s", err)
		}
		return fmt.Errorf("Error retrieving domain: %s", err)
	}

	d.SetId(domain.Name)
	d.Set("name", domain.Name)
	d.Set("ttl", domain.TTL)
	d.Set("zone_file", domain.ZoneFile)

	return nil
}
