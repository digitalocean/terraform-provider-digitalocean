package genai

import (
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanAgents() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        AgentSchemaRead(),
		ResultAttributeName: "agents",
		FlattenRecord:       flattenDigitalOceanAgent,
		GetRecords:          getDigitalOceanAgents,
	}

	return datalist.NewResource(dataListConfig)
}
