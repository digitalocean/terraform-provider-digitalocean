package digitalocean

import (
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDigitalOceanDroplets() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        dropletSchema(),
		ResultAttributeName: "droplets",
		GetRecords:          getDigitalOceanDroplets,
		FlattenRecord:       flattenDigitalOceanDroplet,
	}

	return datalist.NewResource(dataListConfig)
}
