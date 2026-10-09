package microvm

import (
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// DataSourceDigitalOceanMicroVMs returns a plural data source over the
// MicroVMs endpoint with optional `region`, `name`, and `tag_name`
// server-side filters that combine with logical AND. `filter` and `sort`
// (from the datalist framework) apply on top.
func DataSourceDigitalOceanMicroVMs() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        microVMDataSourceSchema(),
		ResultAttributeName: "micro_vms",
		GetRecords:          getDigitalOceanMicroVMs,
		FlattenRecord:       flattenMicroVM,
		ExtraQuerySchema: map[string]*schema.Schema{
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Only return MicroVMs in this region.",
				ValidateFunc: validation.NoZeroValues,
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Only return MicroVMs with this exact name.",
				ValidateFunc: validation.NoZeroValues,
			},
			"tag_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Only return MicroVMs that carry this tag.",
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
	return datalist.NewResource(dataListConfig)
}
