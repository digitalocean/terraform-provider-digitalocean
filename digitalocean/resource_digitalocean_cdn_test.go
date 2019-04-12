package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDigitalOceanCDN_Basic(t *testing.T) {
	digitalOceanBucketName := fmt.Sprintf("tf-cdn-test-bucket-%d", acctest.RandInt())
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

func TestAccDigitalOceanCDN_withTTL(t *testing.T) {
	digitalOceanBucketName := fmt.Sprintf("tf-cdn-test-bucket-%d", acctest.RandInt())
	cdnConfig := fmt.Sprintf(testAccCheckDigitalOceanCDNConfig_withTTL, digitalOceanBucketName)

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
					resource.TestCheckResourceAttr("digitalocean_cdn.foobar", "ttl", "1800"),
				),
			},
		},
	})
}

func TestAccDigitalOceanCDN_Create_And_Update(t *testing.T) {
	digitalOceanBucketName := fmt.Sprintf("tf-cdn-test-bucket-%d", acctest.RandInt())
	cdnCreateConfig := fmt.Sprintf(testAccCheckDigitalOceanCDNConfig_basic, digitalOceanBucketName)
	cdnUpdateConfig := fmt.Sprintf(testAccCheckDigitalOceanCDNConfig_withTTL, digitalOceanBucketName)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanCDNDestroy,
		Steps: []resource.TestStep{
			{
				Config: cdnCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanCDNExists("digitalocean_cdn.foobar"),
					resource.TestCheckResourceAttr(
						"digitalocean_cdn.foobar", "origin", digitalOceanBucketName+".ams3.digitaloceanspaces.com"),
					resource.TestCheckResourceAttr("digitalocean_cdn.foobar", "ttl", "3600"),
				),
			},
			{
				Config: cdnUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanCDNExists("digitalocean_cdn.foobar"),
					resource.TestCheckResourceAttr(
						"digitalocean_cdn.foobar", "origin", digitalOceanBucketName+".ams3.digitaloceanspaces.com"),
					resource.TestCheckResourceAttr("digitalocean_cdn.foobar", "ttl", "1800"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanCDNDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).godoClient()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "digitalocean_cdn" {
			continue
		}

		_, _, err := client.CDNs.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("CDN resource still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanCDNExists(resource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		rs, ok := s.RootModule().Resources[resource]

		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID set for resource: %s", resource)
		}

		foundCDN, _, err := client.CDNs.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundCDN.ID != rs.Primary.ID {
			return fmt.Errorf("Resource not found: %s : %s", resource, rs.Primary.ID)
		}

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

const testAccCheckDigitalOceanCDNConfig_withTTL = `
resource "digitalocean_spaces_bucket" "bucket" {
	name = "%s"
	region = "ams3"
	acl = "public-read"
}

resource "digitalocean_cdn" "foobar" {
	origin = "${digitalocean_spaces_bucket.bucket.bucket_domain_name}"
	ttl = 1800
}`
