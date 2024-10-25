package droplet

import (
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanDroplets() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        dropletSchema(),
		ResultAttributeName: "droplets",
		GetRecords:          getDigitalOceanDroplets,
		FlattenRecord:       flattenDigitalOceanDroplet,
		ExtraQuerySchema: map[string]*schema.Schema{
			"gpus": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}

	return datalist.NewResource(dataListConfig)
}
