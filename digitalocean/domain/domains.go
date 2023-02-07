package domain

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func domainSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "name of the domain",
		},
		"urn": {
			Type:        schema.TypeString,
			Description: "the uniform resource name for the domain",
		},
		"ttl": {
			Type:        schema.TypeInt,
			Description: "ttl of the domain",
		},
	}
}

func getDigitalOceanDomains(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.CombinedConfig).GodoClient()

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	var allDomains []interface{}

	for {
		domains, resp, err := client.Domains.List(context.Background(), opts)

		if err != nil {
			return nil, fmt.Errorf("Error retrieving domains: %s", err)
		}

		for _, domain := range domains {
			allDomains = append(allDomains, domain)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("Error retrieving domains: %s", err)
		}

		opts.Page = page + 1
	}

	return allDomains, nil
}

func flattenDigitalOceanDomain(rawDomain, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	domain := rawDomain.(godo.Domain)

	flattenedDomain := map[string]interface{}{
		"name": domain.Name,
		"urn":  domain.URN(),
		"ttl":  domain.TTL,
	}

	return flattenedDomain, nil
}
