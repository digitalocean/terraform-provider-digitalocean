package digitalocean

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDigitalOceanDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanDomainRead,
		Schema: map[string]*schema.Schema{

			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of the domain",
			},
			// computed attributes
			"ttl": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ttl of the domain",
			},
			"zone_file": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "zone file of the domain",
			},
		},
	}
}

func dataSourceDigitalOceanDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 1,
	}

	domainList := []godo.Domain{}

	for {
		domains, resp, err := client.Domains.List(context.Background(), opts)
		if err != nil {
			d.SetId("")
			return err
		}

		for _, domain := range domains {
			domainList = append(domainList, domain)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			d.SetId("")
			return err
		}

		opts.Page = page + 1
	}

	domain, err := findDomainByName(domainList, d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(domain.Name)
	d.Set("name", domain.Name)
	d.Set("ttl", domain.TTL)
	d.Set("zone_file", domain.ZoneFile)

	return nil
}

func findDomainByName(domains []godo.Domain, name string) (*godo.Domain, error) {
	results := make([]godo.Domain, 0)
	for _, v := range domains {
		if v.Name == name {
			results = append(results, v)
		}
	}
	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no domains found with name %s", name)
	}
	return nil, fmt.Errorf("too many domains found (found %d, expected 1)", len(results))
}
