package microdroplet

import (
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceDigitalOceanMicroDropletImages returns a plural data source over
// the MicroDroplet images endpoint. No extra server-side filters — godo has
// no filtered list endpoint for images. Use `filter` and `sort` from the
// datalist framework for client-side narrowing.
func DataSourceDigitalOceanMicroDropletImages() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema:        microDropletImageDataSourceSchema(),
		ResultAttributeName: "micro_droplet_images",
		GetRecords:          getDigitalOceanMicroDropletImages,
		FlattenRecord:       flattenMicroDropletImage,
	}
	return datalist.NewResource(dataListConfig)
}
