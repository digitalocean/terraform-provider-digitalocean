package digitalocean

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/minio/minio-go"
)

func TestAccDigitalOceanBucket_Basic(t *testing.T) {
	var bucket minio.BucketInfo

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanBucketConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists("digitalocean_bucket.foobar"),
					testAccCheckDigitalOceanBucketAttributes(&bucket),
					resource.TestCheckResourceAttr(
						"digitalocean_bucket.foobar", "name", "foobar"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanBucketDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*minio.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_bucket" {
			continue
		}

		// Try to find the bucket
		_, err := client.BucketExists(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Bucket still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanBucketAttributes(bucket *minio.BucketInfo) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if bucket.Name != "foobar" {
			return fmt.Errorf("Bad name: %s", bucket.Name)
		}

		return nil
	}
}

func testAccCheckDigitalOceanBucketExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Bucket name is set")
		}

		client := testAccProvider.Meta().(*minio.Client)

		// Try to find the bucket
		_, err := client.BucketExists(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Bucket exists")
		}

		return nil
	}
}

var testAccCheckDigitalOceanBucketConfig_basic = fmt.Sprintf(`
resource "digitalocean_bucket" "foobar" {
  name 			= "foobar"
  endpoint 	= "nyc3"
}`)
