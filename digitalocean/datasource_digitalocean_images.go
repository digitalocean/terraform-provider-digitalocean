package digitalocean

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-digitalocean/internal/datalist"
)

func dataSourceDigitalOceanImages() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema: imageSchema(),
		FilterKeys: []string{
			"id",
			"name",
			"type",
			"distribution",
			"slug",
			"image",
			"private",
			"min_disk_size",
			"size_gigabytes",
			"regions",
			"tags",
			"status",
			"error_message",
		},
		SortKeys: []string{
			"id",
			"name",
			"type",
			"distribution",
			"slug",
			"image",
			"private",
			"min_disk_size",
			"size_gigabytes",
			"status",
			"error_message",
		},
		ResultAttributeName: "images",
		FlattenRecord:       flattenDigitalOceanImage,
		GetRecords:          getDigitalOceanImages,
	}

	return datalist.NewResource(dataListConfig)
}
