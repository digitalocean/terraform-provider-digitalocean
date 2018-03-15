package digitalocean

import (
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform/helper/schema"
)

func TestAccDigitalOceanBucket_basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		/*
			IDRefreshName:   "digitalocean_bucket.bucket",
			IDRefreshIgnore: []string{"force_destroy"},
		*/
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanBucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists("digitalocean_bucket.bucket"),
					resource.TestCheckResourceAttr(
						"digitalocean_bucket.bucket", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"digitalocean_bucket.bucket", "bucket", testAccBucketName(rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_bucket.bucket", "bucket_domain_name", testAccBucketDomainName(rInt)),
				),
			},
		},
	})
}

func TestAccDigitalOceanBucket_region(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanBucketConfigWithRegion(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists("digitalocean_bucket.bucket"),
					resource.TestCheckResourceAttr("digitalocean_bucket.bucket", "region", "nyc3"),
				),
			},
		},
	})
}

func TestAccDigitalOceanBucket_UpdateAcl(t *testing.T) {
	ri := acctest.RandInt()
	preConfig := fmt.Sprintf(testAccDigitalOceanBucketConfigWithAcl, ri)
	postConfig := fmt.Sprintf(testAccDigitalOceanBucketConfigWithAclUpdate, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists("digitalocean_bucket.bucket"),
					resource.TestCheckResourceAttr(
						"digitalocean_bucket.bucket", "acl", "public-read"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists("digitalocean_bucket.bucket"),
					resource.TestCheckResourceAttr(
						"digitalocean_bucket.bucket", "acl", "private"),
				),
			},
		},
	})
}

// Test TestAccDigitalOceanBucket_shouldFailNotFound is designed to fail with a "plan
// not empty" error in Terraform, to check against regresssions.
// See https://github.com/hashicorp/terraform/pull/2925
func TestAccDigitalOceanBucket_shouldFailNotFound(t *testing.T) {
	rInt := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanBucketDestroyedConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists("digitalocean_bucket.bucket"),
					testAccCheckDigitalOceanDestroyBucket("digitalocean_bucket.bucket"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckDigitalOceanBucketDestroy(s *terraform.State) error {
	return testAccCheckDigitalOceanBucketDestroyWithProvider(s, testAccProvider)
}

func testAccCheckDigitalOceanBucketDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	sesh, err := session.NewSession(&aws.Config{
		Region:      aws.String("nyc3"),
		Credentials: credentials.NewSharedCredentials("", "digitalocean-spaces")},
	)
	svc := s3.New(sesh, &aws.Config{
		Endpoint: aws.String(fmt.Sprintf("https://%r.digitaloceanspaces.com", "nyc3"))},
	)

	if err != nil {
		log.Fatal(err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_bucket" {
			continue
		}
		_, err := svc.DeleteBucket(&s3.DeleteBucketInput{
			Bucket: aws.String(rs.Primary.ID),
		})
		if err != nil {
			if isAWSErr(err, s3.ErrCodeNoSuchBucket, "") {
				return nil
			}
			return err
		}
	}
	return nil
}

func testAccCheckDigitalOceanBucketExists(n string) resource.TestCheckFunc {
	return testAccCheckDigitalOceanBucketExistsWithProvider(n, func() *schema.Provider { return testAccProvider })
}

func testAccCheckDigitalOceanBucketExistsWithProvider(n string, providerF func() *schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		sesh, err := session.NewSession(&aws.Config{
			Region:      aws.String("nyc3"),
			Credentials: credentials.NewSharedCredentials("", "digitalocean-spaces")},
		)
		svc := s3.New(sesh, &aws.Config{
			Endpoint: aws.String(fmt.Sprintf("https://%r.digitaloceanspaces.com", "nyc3"))},
		)

		if err != nil {
			log.Fatal(err)
		}

		_, err = svc.HeadBucket(&s3.HeadBucketInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			if isAWSErr(err, s3.ErrCodeNoSuchBucket, "") {
				return fmt.Errorf("Spaces bucket not found")
			}
			return err
		}
		return nil

	}
}

func testAccCheckDigitalOceanDestroyBucket(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Spaces Bucket ID is set")
		}

		sesh, err := session.NewSession(&aws.Config{
			Region:      aws.String("nyc3"),
			Credentials: credentials.NewSharedCredentials("", "digitalocean-spaces")},
		)
		svc := s3.New(sesh, &aws.Config{
			Endpoint: aws.String(fmt.Sprintf("https://%r.digitaloceanspaces.com", "nyc3"))},
		)

		if err != nil {
			log.Fatal(err)
		}

		_, err = svc.DeleteBucket(&s3.DeleteBucketInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return fmt.Errorf("Error destroying Bucket (%s) in testAccCheckDigitalOceanDestroyBucket: %s", rs.Primary.ID, err)
		}
		return nil
	}
}

func isAWSErr(err error, code string, message string) bool {
	if err, ok := err.(awserr.Error); ok {
		return err.Code() == code && strings.Contains(err.Message(), message)
	}
	return false
}

// These need a bit of randomness as the name can only be used once globally
// within AWS
func testAccBucketName(randInt int) string {
	return fmt.Sprintf("tf-test-bucket-%d", randInt)
}

func testAccBucketDomainName(randInt int) string {
	return fmt.Sprintf("tf-test-bucket-%d.nyc3.digitaloceanspaces.com", randInt)
}

func testAccDigitalOceanBucketConfig(randInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"
}
`, randInt)
}

func testAccDigitalOceanBucketDestroyedConfig(randInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"
}
`, randInt)
}

func testAccDigitalOceanBucketConfigWithRegion(randInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	region = "nyc3"
}
`, randInt)
}

var testAccDigitalOceanBucketConfigWithAcl = `
resource "digitalocean_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "public-read"
}
`

var testAccDigitalOceanBucketConfigWithAclUpdate = `
resource "digitalocean_bucket" "bucket" {
	bucket = "tf-test-bucket-%d"
	acl = "private"
}
`
