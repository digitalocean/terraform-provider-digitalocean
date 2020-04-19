package digitalocean

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type bucketMetadataStruct struct {
	bucket *s3.Bucket
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

func getSpacesRegions(meta interface{}) ([]godo.Region, error) {
	client := meta.(*CombinedConfig).godoClient()

	var spacesRegions []godo.Region

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	for {
		regions, resp, err := client.Regions.List(context.Background(), opts)

		if err != nil {
			return nil, fmt.Errorf("Error retrieving regions: %s", err)
		}

		for _, region := range regions {
			supportsSpaces := false
			for _, feature := range region.Features {
				if feature == "spaces" {
					supportsSpaces = true
				}
			}

			if supportsSpaces {
				spacesRegions = append(spacesRegions, region)
			}
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("Error retrieving regions: %s", err)
		}

		opts.Page = page + 1
	}

	return spacesRegions, nil
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

func getDigitalOceanBuckets(meta interface{}) ([]interface{}, error) {
	// Retrieve the regions with Spaces enabled.
	spacesRegions, err := getSpacesRegions(meta)
	if err != nil {
		return nil, err
	}

	var buckets []interface{}

	for _, region := range spacesRegions {
		bucketsInRegion, err := getSpacesBucketsInRegion(meta, region.Slug)
		if err != nil {
			return nil, err
		}

		for _, bucketInRegion := range bucketsInRegion {
			metadata := &bucketMetadataStruct{
				bucket: bucketInRegion,
				region: region.Slug,
			}
			buckets = append(buckets, metadata)
		}
	}

	return buckets, nil
}

func flattenSpacesBucket(rawBucketMetadata, meta interface{}) (map[string]interface{}, error) {
	bucketMetadata := rawBucketMetadata.(*bucketMetadataStruct)

	name := *bucketMetadata.bucket.Name
	region := bucketMetadata.region

	flattenedBucket := map[string]interface{}{}
	flattenedBucket["name"] = name
	flattenedBucket["region"] = region
	flattenedBucket["bucket_domain_name"] = bucketDomainName(name, region)
	flattenedBucket["urn"] = fmt.Sprintf("do:space:%s", name)

	return flattenedBucket, nil
}
