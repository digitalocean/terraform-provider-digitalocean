package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

// func init() {
// 	resource.AddTestSweepers("digitalocean_cdn", &resource.Sweeper{
// 		Name: "digitalocean_cdn",
// 		F:    testSweepCertificate,
// 	})

// }

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

// func TestAccDigitalOceanCDN_withCustomDomain(t *testing.T) {

// 	rInt := acctest.RandInt()
// 	privateKeyMaterial, leafCertMaterial, certChainMaterial := generateTestCertMaterial(t)
// 	//domainName := fmt.Sprintf("trenttest%d.com", rInt)
// 	domainName := "trenttest1.com"
// 	bucketName := fmt.Sprintf("tf-cdn-test-bucket-%d", rInt)

// 	cdnConfig := testAccCheckDigitalOceanCDNConfig_withDomain(rInt, privateKeyMaterial, leafCertMaterial, certChainMaterial, domainName, bucketName)

// 	resource.Test(t, resource.TestCase{
// 		PreCheck:     func() { testAccPreCheck(t) },
// 		Providers:    testAccProviders,
// 		CheckDestroy: testAccCheckDigitalOceanCDNDestroy,
// 		Steps: []resource.TestStep{
// 			{
// 				Config: cdnConfig,
// 				Check: resource.ComposeTestCheckFunc(
// 					testAccCheckDigitalOceanCDNExists("digitalocean_cdn.mycdn"),
// 					resource.TestCheckResourceAttr(
// 						"digitalocean_cdn.mycdn", "origin", bucketName+".ams3.digitaloceanspaces.com"),
// 					resource.TestCheckResourceAttr("digitalocean_cdn.mycdn", "ttl", "3600"),
// 				),
// 			},
// 		},
// 	})
// }

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

// func testAccCheckDigitalOceanCDNConfig_withDomain(rInt int, privateKeyMaterial, leafCert, certChain, domainName, bucketName string) string {

// 	return fmt.Sprintf(`
// resource "digitalocean_certificate" "mycert" {
//   	name = "certificate-%d"
//   	private_key = <<EOF
// %s
// EOF
//   	leaf_certificate = <<EOF
// %s
// EOF
//   	certificate_chain = <<EOF
// %s
// EOF
// }

// resource "digitalocean_domain" "mydomain" {
// 	name       = "%s"
// }

// resource "digitalocean_spaces_bucket" "mybucket" {
// 	name = "%s"
// 	region = "ams3"
// 	acl = "public-read"
// }

// resource "digitalocean_cdn" "mycdn" {
// 	origin = "${digitalocean_spaces_bucket.mybucket.bucket_domain_name}"
// 	custom_domain = "${digitalocean_domain.mydomain.name}"
// 	certificate_id = "${digitalocean_certificate.mycert.id}"
// }`, rInt, privateKeyMaterial, leafCert, certChain, domainName, bucketName)
//}
