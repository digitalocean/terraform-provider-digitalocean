package digitalocean

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceDigitalOceanSpacesBuckets_Basic(t *testing.T) {
	bucketName1 := testAccBucketName(acctest.RandInt())
	bucketRegion1 := "nyc3"

	bucketName2 := testAccBucketName(acctest.RandInt())
	bucketRegion2 := "ams3"

	bucketsConfig := fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "bucket1" {
  name = "%s"
  region = "%s"
}

resource "digitalocean_spaces_bucket" "bucket2" {
  name = "%s"
  region = "%s"
}
`, bucketName1, bucketRegion1, bucketName2, bucketRegion2)

	datasourceConfig := fmt.Sprintf(`
data "digitalocean_spaces_buckets" "result" {
  filter {
    key = "name"
    values = ["%s"]
  }
}
`, bucketName1)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: bucketsConfig,
			},
			{
				Config: bucketsConfig + datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_spaces_buckets.result", "buckets.#", "1"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_buckets.result", "buckets.0.name", bucketName1),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_buckets.result", "buckets.0.region", bucketRegion1),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_buckets.result", "buckets.0.bucket_domain_name", bucketDomainName(bucketName1, bucketRegion1)),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_buckets.result", "buckets.0.urn", fmt.Sprintf("do:space:%s", bucketName1)),
				),
			},
			{
				// Remove the datasource from the config so Terraform trying to refresh it does not race with
				// deleting the bucket resources. By removing the datasource from the config here, this ensures
				// that the buckets are deleted after the datasource has been removed from the state.
				Config: bucketsConfig,
			},
		},
	})
}
