package digitalocean

import (
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceDigitalOceanImages() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        imageSchema(),
		ResultAttributeName: "images",
		FlattenRecord:       flattenDigitalOceanImage,
		GetRecords:          getDigitalOceanImages,
	}

	return datalist.NewResource(dataListConfig)
}
