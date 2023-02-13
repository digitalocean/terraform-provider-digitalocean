package cdn_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const originSuffix = ".ams3.digitaloceanspaces.com"

func TestAccDigitalOceanCDN_Create(t *testing.T) {

	bucketName := generateBucketName()
	cdnCreateConfig := fmt.Sprintf(testAccCheckDigitalOceanCDNConfig_Create, bucketName)

	expectedOrigin := bucketName + originSuffix
	expectedTTL := "3600"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanCDNDestroy,
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
	ttl := 600
	cdnCreateConfig := fmt.Sprintf(testAccCheckDigitalOceanCDNConfig_Create_with_TTL, bucketName, ttl)

	expectedOrigin := bucketName + originSuffix
	expectedTTL := "600"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanCDNDestroy,
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
	ttl := 600

	cdnCreateConfig := fmt.Sprintf(testAccCheckDigitalOceanCDNConfig_Create, bucketName)
	cdnUpdateConfig := fmt.Sprintf(testAccCheckDigitalOceanCDNConfig_Create_with_TTL, bucketName, ttl)

	expectedOrigin := bucketName + originSuffix
	expectedTTL := "600"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanCDNDestroy,
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

func TestAccDigitalOceanCDN_CustomDomain(t *testing.T) {
	spaceName := generateBucketName()
	certName := acceptance.RandomTestName()
	updatedCertName := generateBucketName()
	domain := acceptance.RandomTestName() + ".com"
	config := testAccCheckDigitalOceanCDNConfig_CustomDomain(domain, spaceName, certName)
	updatedConfig := testAccCheckDigitalOceanCDNConfig_CustomDomain(domain, spaceName, updatedCertName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanCDNDestroy,
		ExternalProviders: map[string]resource.ExternalProvider{
			"tls": {
				Source:            "hashicorp/tls",
				VersionConstraint: "3.0.0",
			},
		},
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanCDNExists("digitalocean_cdn.space_cdn"),
					resource.TestCheckResourceAttr(
						"digitalocean_cdn.space_cdn", "certificate_name", certName),
					resource.TestCheckResourceAttr(
						"digitalocean_cdn.space_cdn", "custom_domain", "foo."+domain),
				),
			},
			{
				Config: updatedConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanCDNExists("digitalocean_cdn.space_cdn"),
					resource.TestCheckResourceAttr(
						"digitalocean_cdn.space_cdn", "certificate_name", updatedCertName),
					resource.TestCheckResourceAttr(
						"digitalocean_cdn.space_cdn", "custom_domain", "foo."+domain),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanCDNDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

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
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

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
	return acceptance.RandomTestName("cdn")
}

const testAccCheckDigitalOceanCDNConfig_Create = `
resource "digitalocean_spaces_bucket" "bucket" {
  name   = "%s"
  region = "ams3"
  acl    = "public-read"
}

resource "digitalocean_cdn" "foobar" {
  origin = digitalocean_spaces_bucket.bucket.bucket_domain_name
}`

const testAccCheckDigitalOceanCDNConfig_Create_with_TTL = `
resource "digitalocean_spaces_bucket" "bucket" {
  name   = "%s"
  region = "ams3"
  acl    = "public-read"
}

resource "digitalocean_cdn" "foobar" {
  origin = digitalocean_spaces_bucket.bucket.bucket_domain_name
  ttl    = %d
}`

func testAccCheckDigitalOceanCDNConfig_CustomDomain(domain string, spaceName string, certName string) string {
	return fmt.Sprintf(`
resource "tls_private_key" "example" {
  algorithm = "RSA"
}

resource "tls_self_signed_cert" "example" {
  key_algorithm   = "RSA"
  private_key_pem = tls_private_key.example.private_key_pem
  dns_names       = ["foo.%s"]
  subject {
    common_name  = "foo.%s"
    organization = "%s"
  }

  validity_period_hours = 24

  allowed_uses = [
    "key_encipherment",
    "digital_signature",
    "server_auth",
  ]
}

resource "digitalocean_spaces_bucket" "space" {
  name   = "%s"
  region = "sfo3"
}

resource "digitalocean_certificate" "spaces_cert" {
  name             = "%s"
  type             = "custom"
  private_key      = tls_private_key.example.private_key_pem
  leaf_certificate = tls_self_signed_cert.example.cert_pem

  lifecycle {
    create_before_destroy = true
  }
}

resource digitalocean_domain "domain" {
  name = "%s"
}

resource digitalocean_record "record" {
  domain = digitalocean_domain.domain.name
  type   = "CNAME"
  name   = "foo"
  value  = "${digitalocean_spaces_bucket.space.bucket_domain_name}."
}

resource "digitalocean_cdn" "space_cdn" {
  depends_on = [
    digitalocean_spaces_bucket.space,
    digitalocean_certificate.spaces_cert,
    digitalocean_record.record
  ]

  origin           = digitalocean_spaces_bucket.space.bucket_domain_name
  ttl              = 600
  certificate_name = digitalocean_certificate.spaces_cert.name
  custom_domain    = "foo.%s"
}`, domain, domain, certName, spaceName, certName, domain, domain)
}
