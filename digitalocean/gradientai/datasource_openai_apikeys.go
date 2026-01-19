package gradientai

import (
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanOpenAIApiKeys() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        OpenAIApiKeySchemaRead(),
		ResultAttributeName: "openai_api_keys",
		FlattenRecord:       flattenOpenAIApiKeyInfo,
		GetRecords:          getDigitalOceanOpenAIApiKeys,
	}

	return datalist.NewResource(dataListConfig)
}

func DataSourceDigitalOceanAgentsByOpenAIApiKey() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanAgentsByOpenAIApiKeyRead,
		Schema: map[string]*schema.Schema{
			"uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The UUID of the OpenAI API key.",
			},
			"agents": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of agents associated with the OpenAI API key.",
				Elem:        &schema.Resource{Schema: AgentSchemaRead()},
			},
		},
	}
}
