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
		GetRecords: func(meta interface{}) ([]interface{}, error) {
			client := meta.(*CombinedConfig).godoClient()
			projects, err := getDigitalOceanProjects(client)
			if err == nil {
				var rawProjects []interface{}
				for _, project := range projects {
					rawProjects = append(rawProjects, project)
				}
				return rawProjects, nil
			} else {
				return nil, err
			}
		},
	}

	return datalist.NewResource(dataListConfig)
}
