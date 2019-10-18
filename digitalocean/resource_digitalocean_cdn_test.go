package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const originSuffix = ".ams3.digitaloceanspaces.com"

func TestAccDigitalOceanCDN_Create(t *testing.T) {

	bucketName := generateBucketName()
	cdnCreateConfig := fmt.Sprintf(testAccCheckDigitalOceanCDNConfig_Create, bucketName)

	expectedOrigin := bucketName + originSuffix
	expectedTTL := "3600"

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
						"digitalocean_cdn.foobar", "origin", expectedOrigin),
					resource.TestCheckResourceAttr("digitalocean_cdn.foobar", "ttl", expectedTTL),
				),
			},
		},
	})
}

func TestAccDigitalOceanCDN_Create_with_TTL(t *testing.T) {

	bucketName := generateBucketName()
	ttl := 1800
	cdnCreateConfig := fmt.Sprintf(testAccCheckDigitalOceanCDNConfig_Create_with_TTL, bucketName, ttl)

	expectedOrigin := bucketName + originSuffix
	expectedTTL := "1800"

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
						"digitalocean_cdn.foobar", "origin", expectedOrigin),
					resource.TestCheckResourceAttr("digitalocean_cdn.foobar", "ttl", expectedTTL),
				),
			},
		},
	})
}

func TestAccDigitalOceanCDN_Create_and_Update(t *testing.T) {

	bucketName := generateBucketName()
	ttl := 1800

	cdnCreateConfig := fmt.Sprintf(testAccCheckDigitalOceanCDNConfig_Create, bucketName)
	cdnUpdateConfig := fmt.Sprintf(testAccCheckDigitalOceanCDNConfig_Create_with_TTL, bucketName, ttl)

	expectedOrigin := bucketName + originSuffix
	expectedTTL := "1800"

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
						"digitalocean_cdn.foobar", "origin", expectedOrigin),
					resource.TestCheckResourceAttr("digitalocean_cdn.foobar", "ttl", "3600"),
				),
			},
			{
				Config: cdnUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanCDNExists("digitalocean_cdn.foobar"),
					resource.TestCheckResourceAttr(
						"digitalocean_cdn.foobar", "origin", expectedOrigin),
					resource.TestCheckResourceAttr("digitalocean_cdn.foobar", "ttl", expectedTTL),
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

func generateBucketName() string {
	return fmt.Sprintf("tf-cdn-test-bucket-%d", acctest.RandInt())
}

const testAccCheckDigitalOceanCDNConfig_Create = `
resource "digitalocean_spaces_bucket" "bucket" {
	name = "%s"
	region = "ams3"
	acl = "public-read"
}

resource "digitalocean_cdn" "foobar" {
	origin = "${digitalocean_spaces_bucket.bucket.bucket_domain_name}"
}`

const testAccCheckDigitalOceanCDNConfig_Create_with_TTL = `
resource "digitalocean_spaces_bucket" "bucket" {
	name = "%s"
	region = "ams3"
	acl = "public-read"
}

resource "digitalocean_cdn" "foobar" {
	origin = "${digitalocean_spaces_bucket.bucket.bucket_domain_name}"
	ttl = %d
}`
