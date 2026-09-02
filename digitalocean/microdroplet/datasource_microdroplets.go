package microdroplet

import (
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceDigitalOceanMicroDroplets returns a plural data source over the
// MicroDroplets endpoint with optional `region` and `name` server-side
// filters. `filter` and `sort` (from the datalist framework) apply on top.
func DataSourceDigitalOceanMicroDroplets() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        microDropletDataSourceSchema(),
		ResultAttributeName: "micro_droplets",
		GetRecords:          getDigitalOceanMicroDroplets,
		FlattenRecord:       flattenMicroDroplet,
		ExtraQuerySchema: map[string]*schema.Schema{
			"region": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"region"},
			},
		},
	}
	return datalist.NewResource(dataListConfig)
}
