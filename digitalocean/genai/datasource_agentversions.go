package genai

import (
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanAgentVersions() *schema.Resource {

	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        AgentVersionSchemaRead(),
		ResultAttributeName: "agent_versions",
		ExtraQuerySchema: map[string]*schema.Schema{
			"agent_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the agent to fetch versions for",
			},
		},
		FlattenRecord: flattenDigitalOceanAgentVersion,
		GetRecords:    getDigitalOceanAgentVersions,
	}

	return datalist.NewResource(dataListConfig)
}
