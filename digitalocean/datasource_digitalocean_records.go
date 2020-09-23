package digitalocean

import (
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDigitalOceanRecords() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        recordsSchema(),
		ResultAttributeName: "records",
		ExtraQuerySchema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		FlattenRecordWithExtraQuery: flattenDigitalOceanRecord,
		GetRecordsWithExtraQuery:    getDigitalOceanRecords,
	}

	return datalist.NewResource(dataListConfig)
}
