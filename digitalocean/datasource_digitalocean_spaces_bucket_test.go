package digitalocean

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceDigitalOceanSpacesBucket_Basic(t *testing.T) {
	rInt := acctest.RandInt()
	bucketName := testAccBucketName(rInt)
	bucketRegion := "nyc3"

	resourceConfig := fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "bucket" {
	name = "%s"
	region = "%s"
}
`, bucketName, bucketRegion)

	datasourceConfig := fmt.Sprintf(`
data "digitalocean_spaces_bucket" "bucket" {
    name = "%s"
    region = "%s"
}
`, bucketName, bucketRegion)

	config1 := resourceConfig
	config2 := config1 + datasourceConfig

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: config1,
			},
			{
				Config: config2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket.bucket", "name", bucketName),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket.bucket", "region", bucketRegion),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket.bucket", "bucket_domain_name", bucketDomainName(bucketName, bucketRegion)),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket.bucket", "urn", fmt.Sprintf("do:space:%s", bucketName)),
				),
			},
			{
				Config: config1,
			},
		},
	})
}

func TestAccDataSourceDigitalOceanSpacesBucket_NotFound(t *testing.T) {
	datasourceConfig := `
data "digitalocean_spaces_bucket" "bucket" {
    name = "no-such-bucket"
    region = "nyc3"
}
`

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config:      datasourceConfig,
				ExpectError: regexp.MustCompile("Spaces Bucket.*not found"),
			},
		},
	})
}
