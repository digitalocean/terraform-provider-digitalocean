package spaces

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/sweep"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func init() {
	resource.AddTestSweepers("digitalocean_spaces_bucket", &resource.Sweeper{
		Name: "digitalocean_spaces_bucket",
		F:    sweepSpaces,
	})

}

func sweepSpaces(region string) error {
	meta, err := sweep.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	for _, r := range SpacesRegions {
		client, err := meta.(*config.CombinedConfig).SpacesClient(r)
		if err != nil {
			return fmt.Errorf("Error building Spaces client: %s", err)
		}

		svc := s3.New(client)

		buckets, err := getSpacesBucketsInRegion(meta, r)
		if err != nil {
			return err
		}

		for _, b := range buckets {
			if strings.HasPrefix(*b.Name, sweep.TestNamePrefix) {
				log.Printf("[DEBUG] Destroying Spaces bucket %s in %s", *b.Name, r)

				_, err = svc.DeleteBucket(&s3.DeleteBucketInput{
					Bucket: b.Name,
				})
				if err != nil {
					if IsAWSErr(err, "BucketNotEmpty", "") {
						log.Printf("[DEBUG] Deleting objects in Spaces bucket %s in %s", *b.Name, r)
						err := spacesBucketForceDelete(svc, *b.Name)
						if err != nil {
							return err
						}
					} else {
						return err
					}
				}
			}
		}
	}

	return nil
}
