package digitalocean

import (
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDigitalOceanDomains() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        domainSchema(),
		ResultAttributeName: "domains",
		GetRecords:          getDigitalOceanDomains,
		FlattenRecord:       flattenDigitalOceanDomain,
	}

	return datalist.NewResource(dataListConfig)
}
