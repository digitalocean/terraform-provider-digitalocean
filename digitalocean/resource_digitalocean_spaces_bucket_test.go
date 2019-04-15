package digitalocean

import (
	"fmt"
	"log"
	"os"
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

	expectedRegion := "ams3"
	expectedBucketName := testAccBucketName(rInt)
	expectBucketURN := fmt.Sprintf("do:space:%s", expectedBucketName)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		/*
			IDRefreshName:   "digitalocean_spaces_bucket.bucket",
			IDRefreshIgnore: []string{"force_destroy"},
		*/
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanBucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists("digitalocean_spaces_bucket.bucket"),
					resource.TestCheckResourceAttr(
						"digitalocean_spaces_bucket.bucket", "region", expectedRegion),
					resource.TestCheckResourceAttr(
						"digitalocean_spaces_bucket.bucket", "name", expectedBucketName),
					resource.TestCheckResourceAttr("digitalocean_spaces_bucket.bucket", "urn", expectBucketURN),
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
					testAccCheckDigitalOceanBucketExists("digitalocean_spaces_bucket.bucket"),
					resource.TestCheckResourceAttr("digitalocean_spaces_bucket.bucket", "region", "ams3"),
				),
			},
		},
	})
}

func TestAccDigitalOceanBucket_UpdateAcl(t *testing.T) {
	ri := acctest.RandInt()
	preConfig := fmt.Sprintf(testAccDigitalOceanBucketConfigWithACL, ri)
	postConfig := fmt.Sprintf(testAccDigitalOceanBucketConfigWithACLUpdate, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists("digitalocean_spaces_bucket.bucket"),
					resource.TestCheckResourceAttr(
						"digitalocean_spaces_bucket.bucket", "acl", "public-read"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists("digitalocean_spaces_bucket.bucket"),
					resource.TestCheckResourceAttr(
						"digitalocean_spaces_bucket.bucket", "acl", "private"),
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
					testAccCheckDigitalOceanBucketExists("digitalocean_spaces_bucket.bucket"),
					testAccCheckDigitalOceanDestroyBucket("digitalocean_spaces_bucket.bucket"),
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

	for _, rs := range s.RootModule().Resources {
		sesh, err := session.NewSession(&aws.Config{
			Region:      aws.String(rs.Primary.Attributes["region"]),
			Credentials: credentials.NewStaticCredentials(os.Getenv("SPACES_ACCESS_KEY_ID"), os.Getenv("SPACES_SECRET_ACCESS_KEY"), "")},
		)

		svc := s3.New(sesh, &aws.Config{
			Endpoint: aws.String(fmt.Sprintf("https://%s.digitaloceanspaces.com", rs.Primary.Attributes["region"]))},
		)

		if err != nil {
			log.Fatal(err)
		}

		if rs.Type != "digitalocean_spaces_bucket" {
			continue
		}
		_, err = svc.DeleteBucket(&s3.DeleteBucketInput{
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
			Region:      aws.String(rs.Primary.Attributes["region"]),
			Credentials: credentials.NewStaticCredentials(os.Getenv("SPACES_ACCESS_KEY_ID"), os.Getenv("SPACES_SECRET_ACCESS_KEY"), "")},
		)
		svc := s3.New(sesh, &aws.Config{
			Endpoint: aws.String(fmt.Sprintf("https://%s.digitaloceanspaces.com", rs.Primary.Attributes["region"]))},
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
			Region:      aws.String(rs.Primary.Attributes["region"]),
			Credentials: credentials.NewStaticCredentials(os.Getenv("SPACES_ACCESS_KEY_ID"), os.Getenv("SPACES_SECRET_ACCESS_KEY"), "")},
		)
		svc := s3.New(sesh, &aws.Config{
			Endpoint: aws.String(fmt.Sprintf("https://%s.digitaloceanspaces.com", rs.Primary.Attributes["region"]))},
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
func testAccBucketName(randInt int) string {
	return fmt.Sprintf("tf-test-bucket-%d", randInt)
}

func testAccDigitalOceanBucketConfig(randInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "bucket" {
	name = "tf-test-bucket-%d"
	acl = "public-read"
	region = "ams3"
}
`, randInt)
}

func testAccDigitalOceanBucketDestroyedConfig(randInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "bucket" {
	name = "tf-test-bucket-%d"
	acl = "public-read"
}
`, randInt)
}

func testAccDigitalOceanBucketConfigWithRegion(randInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "bucket" {
	name = "tf-test-bucket-%d"
	region = "ams3"
}
`, randInt)
}

func testAccDigitalOceanBucketConfigImport(randInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "bucket" {
	name = "tf-test-bucket-%d"
}
`, randInt)
}

var testAccDigitalOceanBucketConfigWithACL = `
resource "digitalocean_spaces_bucket" "bucket" {
	name = "tf-test-bucket-%d"
	acl = "public-read"
}
`

var testAccDigitalOceanBucketConfigWithACLUpdate = `
resource "digitalocean_spaces_bucket" "bucket" {
	name = "tf-test-bucket-%d"
	acl = "private"
}
`
