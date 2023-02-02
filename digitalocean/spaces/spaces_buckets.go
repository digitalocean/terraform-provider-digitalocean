package spaces

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	awspolicy "github.com/hashicorp/awspolicyequivalence"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var (
	// SpacesRegions is a list of DigitalOcean regions that support Spaces.
	SpacesRegions = []string{"ams3", "fra1", "nyc3", "sfo2", "sfo3", "sgp1", "syd1"}
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
		"endpoint": {
			Type:        schema.TypeString,
			Description: "The FQDN of the bucket without the bucket name",
		},
	}
}

func getSpacesBucketsInRegion(meta interface{}, region string) ([]*s3.Bucket, error) {
	client, err := meta.(*config.CombinedConfig).SpacesClient(region)
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
	var buckets []interface{}

	for _, region := range SpacesRegions {
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
	flattenedBucket["bucket_domain_name"] = BucketDomainName(name, region)
	flattenedBucket["urn"] = fmt.Sprintf("do:space:%s", name)
	flattenedBucket["endpoint"] = BucketEndpoint(region)

	return flattenedBucket, nil
}

func CompareSpacesBucketPolicy(policy1, policy2 string) bool {
	equivalent, err := awspolicy.PoliciesAreEquivalent(policy1, policy2)
	if err != nil {
		return false
	}
	return equivalent
}
