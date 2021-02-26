package digitalocean

import (
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDigitalOceanFirewalls() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        firewallSchema(),
		ResultAttributeName: "firewalls",
		GetRecords:          getDigitalOceanFirewalls,
	}

	return datalist.NewResource(dataListConfig)
}