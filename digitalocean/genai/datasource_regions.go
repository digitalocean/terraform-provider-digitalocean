package genai

import (
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanRegions() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        RegionSchemaRead().Schema,
		ResultAttributeName: "regions",
		FlattenRecord:       flattenDigitalOceanRegion,
		GetRecords:          getDigitalOceanRegions,
	}

	return datalist.NewResource(dataListConfig)
}
