package digitalocean

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type bucketMetadataStruct struct {
	name   string
	region string
}

func spacesBucketSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Type:        schema.TypeString,
			Description: "Bucket name",
		},
		"urn": {
			Type:        schema.TypeString,
			Description: "the uniform resource name for the bucket",
		},
		"region": {
			Type:        schema.TypeString,
			Description: "Bucket region",
		},
		"bucket_domain_name": {
			Type:        schema.TypeString,
			Description: "The FQDN of the bucket",
		},
	}
}

func getSpacesBucketsInRegion(meta interface{}, region string) ([]*s3.Bucket, error) {
	client, err := meta.(*CombinedConfig).spacesClient(region)
	if err != nil {
		return nil, err
	}

	svc := s3.New(client)

	input := s3.ListBucketsInput{}
	output, err := svc.ListBuckets(&input)
	if err != nil {
		return nil, err
	}

	return output.Buckets, nil
}

func getDigitalOceanBuckets(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	// The DigitalOcean API does not currently return what regions have Spaces available. Thus, this
	// function hard-codes the regions in which Spaces operates.
	//
	// This list is current as of April 20, 2020 and is from:
	// https://www.digitalocean.com/docs/platform/availability-matrix/#other-product-availability
	spacesRegions := []string{"ams3", "fra1", "nyc3", "sfo2", "sgp1"}

	var buckets []interface{}

	for _, region := range spacesRegions {
		bucketsInRegion, err := getSpacesBucketsInRegion(meta, region)
		if err != nil {
			return nil, err
		}

		for _, bucketInRegion := range bucketsInRegion {
			metadata := &bucketMetadataStruct{
				name:   *bucketInRegion.Name,
				region: region,
			}
			buckets = append(buckets, metadata)
		}
	}

	return buckets, nil
}

func flattenSpacesBucket(rawBucketMetadata, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	bucketMetadata := rawBucketMetadata.(*bucketMetadataStruct)

	name := bucketMetadata.name
	region := bucketMetadata.region

	flattenedBucket := map[string]interface{}{}
	flattenedBucket["name"] = name
	flattenedBucket["region"] = region
	flattenedBucket["bucket_domain_name"] = bucketDomainName(name, region)
	flattenedBucket["urn"] = fmt.Sprintf("do:space:%s", name)

	return flattenedBucket, nil
}
