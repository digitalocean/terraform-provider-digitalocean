package digitalocean

import (
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceDigitalOceanDomains() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema: domainSchema(),
		FilterKeys: []string{
			"name",
			"urn",
			"ttl",
		},
		SortKeys: []string{
			"name",
			"urn",
			"ttl",
		},
		ResultAttributeName: "domains",
		GetRecords:          getDigitalOceanDomains,
		FlattenRecord:       flattenDigitalOceanDomain,
	}

	return datalist.NewResource(dataListConfig)
}
