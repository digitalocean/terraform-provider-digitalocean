package image

import (
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanImages() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        imageSchema(),
		ResultAttributeName: "images",
		FlattenRecord:       flattenDigitalOceanImage,
		GetRecords:          getDigitalOceanImages,
	}

	return datalist.NewResource(dataListConfig)
}
