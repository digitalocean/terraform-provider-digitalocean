package digitalocean

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-digitalocean/internal/datalist"
)

func dataSourceDigitalOceanProjects() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema: projectSchema(),
		FilterKeys: []string{
			"name",
			"purpose",
			"description",
			"environment",
			"is_default",
		},
		SortKeys: []string{
			"name",
			"purpose",
			"description",
			"environment",
		},
		ResultAttributeName: "projects",
		FlattenRecord:       flattenDigitalOceanProject,
		GetRecords:          getDigitalOceanProjects,
	}

	return datalist.NewResource(dataListConfig)
}
