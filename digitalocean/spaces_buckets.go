package digitalocean

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/service/s3"
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

// TODO: Hard-coding the Spaces regions for now given no way to filter out regions like nyc1
// which do not have spaces.
//
//func getSpacesRegions(meta interface{}) ([]string, error) {
//	client := meta.(*CombinedConfig).godoClient()
//
//	var spacesRegions []string
//
//	opts := &godo.ListOptions{
//		Page:    1,
//		PerPage: 200,
//	}
//
//	for {
//		regions, resp, err := client.Regions.List(context.Background(), opts)
//
//		if err != nil {
//			return nil, fmt.Errorf("Error retrieving regions: %s", err)
//		}
//
//		// TODO: Filter out regions without Spaces. It is unclear what feature is set
//		// to indicate Spaces is available in a region because, for example, both
//		// nyc1 and nyc3 have "storage" as a feature even though nyc3 is the Spaces region in NY.
//		for _, region := range regions {
//			spacesRegions = append(spacesRegions, region.Slug)
//		}
//
//		if resp.Links == nil || resp.Links.IsLastPage() {
//			break
//		}
//
//		page, err := resp.Links.CurrentPage()
//		if err != nil {
//			return nil, fmt.Errorf("Error retrieving regions: %s", err)
//		}
//
//		opts.Page = page + 1
//	}
//
//	return spacesRegions, nil
//}

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
	//spacesRegions, err := getSpacesRegions(meta)
	//if err != nil {
	//	return nil, err
	//}
	spacesRegions := []string{"ams3", "fra1", "nyc3", "sfo2", "sgp1"}
	log.Printf("[DEBUG] spacesRegions = %v", spacesRegions)

	var buckets []interface{}

	for _, region := range spacesRegions {
		bucketsInRegion, err := getSpacesBucketsInRegion(meta, region)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] bucketsInRegion(%s) = %v", region, bucketsInRegion)

		for _, bucketInRegion := range bucketsInRegion {
			metadata := &bucketMetadataStruct{
				bucket: bucketInRegion,
				region: region,
			}
			buckets = append(buckets, metadata)
		}
	}

	log.Printf("buckets = %v", buckets)
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
