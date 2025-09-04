package genai

import (
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanModels() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        ModelSchemaRead().Schema,
		ResultAttributeName: "models",
		FlattenRecord:       flattenDigitalOceanModel,
		GetRecords:          getDigitalOceanModels,
	}

	return datalist.NewResource(dataListConfig)
}
