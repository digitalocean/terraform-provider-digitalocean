package digitalocean

import (
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceDigitalOceanProjects() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        projectSchema(),
		ResultAttributeName: "projects",
		FlattenRecord:       flattenDigitalOceanProject,
		GetRecords:          getDigitalOceanProjects,
	}

	return datalist.NewResource(dataListConfig)
}
