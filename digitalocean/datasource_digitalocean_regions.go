package digitalocean

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-digitalocean/internal/datalist"
)

func dataSourceDigitalOceanRegions() *schema.Resource {
	dataListConfig := &datalist.DataListResourceConfig{
		RecordSchema: map[string]*schema.Schema{
			"slug": {
				Type: schema.TypeString,
			},
			"name": {
				Type: schema.TypeString,
			},
			"sizes": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"features": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"available": {
				Type: schema.TypeBool,
			},
		},
		FilterKeys: []string{
			"slug",
			"name",
			"available",
			"features",
			"sizes",
		},
		SortKeys: []string{
			"slug",
			"name",
			"available",
		},
		ResultAttributeName: "regions",
		FlattenRecord:       flattenRegion,
		GetRecords:          getDigitalOceanRegions,
	}

	return datalist.NewDataListResource(dataListConfig)
}
