package gradientai

import (
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanCustomModels() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        CustomModelSchemaRead().Schema,
		ResultAttributeName: "custom_models",
		FlattenRecord:       flattenDigitalOceanCustomModel,
		GetRecords:          getDigitalOceanCustomModels,
		ExtraQuerySchema: map[string]*schema.Schema{
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional status filter forwarded to the list API (e.g. STATUS_READY).",
			},
		},
	}

	return datalist.NewResource(dataListConfig)
}
