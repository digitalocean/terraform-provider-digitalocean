package genai

import (
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanKnowledgeBases() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        KnowledgeBaseSchemaRead(),
		ResultAttributeName: "knowledge_bases",
		FlattenRecord:       flattenDigitalOceanKnowledgeBase,
		GetRecords:          getDigitalOceanKnowledgeBases,
	}

	return datalist.NewResource(dataListConfig)
}
