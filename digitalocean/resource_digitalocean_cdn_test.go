package digitalocean

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDigitalOceanCDN_Basic(t *testing.T) {
	digitalOceanBucketName := fmt.Sprintf("tf-test-bucket-%d", acctest.RandInt())
	cdnConfig := fmt.Sprintf(testAccCheckDigitalOceanCDNConfig_basic, digitalOceanBucketName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanCDNDestroy,
		Steps: []resource.TestStep{
			{
				Config: cdnConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanCDNExists("digitalocean_cdn.foobar"),
					resource.TestCheckResourceAttr(
						"digitalocean_cdn.foobar", "origin", digitalOceanBucketName+".ams3.digitaloceanspaces.com"),
					resource.TestCheckResourceAttr("digitalocean_cdn.foobar", "ttl", "3600"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanCDNDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).godoClient()

	return nil
}

func testAccCheckDigitalOceanCDNExists(resource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		return nil
	}
}

const testAccCheckDigitalOceanCDNConfig_basic = `
resource "digitalocean_spaces_bucket" "bucket" {
	name = "%s"
	region = "ams3"
	acl = "public-read"
}

resource "digitalocean_cdn" "foobar" {
	origin = "${digitalocean_spaces_bucket.bucket.bucket_domain_name}"
}`
