package microvm

import (
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceDigitalOceanMicroVMs returns a plural data source over the
// MicroVMs endpoint with optional `region` and `name` server-side
// filters. `filter` and `sort` (from the datalist framework) apply on top.
func DataSourceDigitalOceanMicroVMs() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        microVMDataSourceSchema(),
		ResultAttributeName: "micro_vms",
		GetRecords:          getDigitalOceanMicroVMs,
		FlattenRecord:       flattenMicroVM,
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
