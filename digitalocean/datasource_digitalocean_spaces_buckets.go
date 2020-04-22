package digitalocean

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-digitalocean/internal/datalist"
)

func dataSourceDigitalOceanSpacesBuckets() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema: spacesBucketSchema(),
		FilterKeys: []string{
			"bucket_domain_name",
			"name",
			"region",
			"urn",
		},
		SortKeys: []string{
			"bucket_domain_name",
			"name",
			"region",
			"urn",
		},
		ResultAttributeName: "buckets",
		FlattenRecord:       flattenSpacesBucket,
		GetRecords:          getDigitalOceanBuckets,
	}

	return datalist.NewResource(dataListConfig)
}
