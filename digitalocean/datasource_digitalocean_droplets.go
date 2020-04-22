package digitalocean

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-digitalocean/internal/datalist"
)

func dataSourceDigitalOceanDroplets() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema: dropletSchema(),
		FilterKeys: []string{
			"id",
			"name",
			"created_at",
			"urn",
			"region",
			"image",
			"size",
			"disk",
			"vcpus",
			"memory",
			"price_hourly",
			"price_monthly",
			"status",
			"locked",
			"ipv4_address",
			"ipv4_address_private",
			"ipv6_address",
			"ipv6_address_private",
			"backups",
			"ipv6",
			"private_networking",
			"monitoring",
			"volume_ids",
			"tags",
			"vpc_uuid",
		},
		SortKeys: []string{
			"id",
			"name",
			"created_at",
			"urn",
			"region",
			"image",
			"size",
			"disk",
			"vcpus",
			"memory",
			"price_hourly",
			"price_monthly",
			"status",
			"locked",
			"ipv4_address",
			"ipv4_address_private",
			"ipv6_address",
			"ipv6_address_private",
			"backups",
			"ipv6",
			"private_networking",
			"monitoring",
			"vpc_uuid",
		},
		ResultAttributeName: "droplets",
		GetRecords:          getDigitalOceanDroplets,
		FlattenRecord:       flattenDigitalOceanDroplet,
	}

	return datalist.NewResource(dataListConfig)
}
