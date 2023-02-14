package spaces_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/spaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanSpacesBucket_Basic(t *testing.T) {
	bucketName := acceptance.RandomTestName()
	bucketRegion := "nyc3"

	resourceConfig := fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "bucket" {
  name   = "%s"
  region = "%s"
}
`, bucketName, bucketRegion)

	datasourceConfig := fmt.Sprintf(`
data "digitalocean_spaces_bucket" "bucket" {
  name   = "%s"
  region = "%s"
}
`, bucketName, bucketRegion)

	config1 := resourceConfig
	config2 := config1 + datasourceConfig

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: config1,
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket.bucket", "name", bucketName),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket.bucket", "region", bucketRegion),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket.bucket", "bucket_domain_name", spaces.BucketDomainName(bucketName, bucketRegion)),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket.bucket", "endpoint", spaces.BucketEndpoint(bucketRegion)),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket.bucket", "urn", fmt.Sprintf("do:space:%s", bucketName)),
				),
			},
			{
				// Remove the datasource from the config so Terraform trying to refresh it does not race with
				// deleting the bucket resource. By removing the datasource from the config here, this ensures
				// that the bucket will be deleted after the datasource has been removed from the state.
				Config: config1,
			},
		},
	})
}

func TestAccDataSourceDigitalOceanSpacesBucket_NotFound(t *testing.T) {
	datasourceConfig := `
data "digitalocean_spaces_bucket" "bucket" {
  name   = "no-such-bucket"
  region = "nyc3"
}
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config:      datasourceConfig,
				ExpectError: regexp.MustCompile("Spaces Bucket.*not found"),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanSpacesBucket_RegionError(t *testing.T) {
	badRegion := "ny2"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
data "digitalocean_spaces_bucket" "bucket" {
  name   = "tf-test-bucket"
  region = "%s"
}`, badRegion),
				ExpectError: regexp.MustCompile(`expected region to be one of`),
			},
		},
	})
}
