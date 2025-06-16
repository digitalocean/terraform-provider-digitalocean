package genai

import (
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanAgents() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        AgentSchemaRead(),
		ResultAttributeName: "agents",
		ExtraQuerySchema: map[string]*schema.Schema{
			"only_deployed": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
		FlattenRecord: flattenDigitalOceanAgent,
		GetRecords:    getDigitalOceanAgents,
	}

	return datalist.NewResource(dataListConfig)
}
