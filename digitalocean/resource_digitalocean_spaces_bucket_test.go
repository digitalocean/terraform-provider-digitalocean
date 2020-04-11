package digitalocean

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

func TestAccDigitalOceanBucket_UpdateCors(t *testing.T) {
	ri := acctest.RandInt()
	preConfig := fmt.Sprintf(testAccDigitalOceanBucketConfigWithCORS, ri)
	postConfig := fmt.Sprintf(testAccDigitalOceanBucketConfigWithCORSUpdate, ri)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: preConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists("digitalocean_spaces_bucket.bucket"),
					resource.TestCheckNoResourceAttr(
						"digitalocean_spaces_bucket.bucket", "cors_rule"),
				),
			},
			{
				Config: postConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists("digitalocean_spaces_bucket.bucket"),
					testAccCheckDigitalOceanBucketCors(
						"digitalocean_spaces_bucket.bucket",
						[]*s3.CORSRule{
							{
								AllowedHeaders: []*string{aws.String("*")},
								AllowedMethods: []*string{aws.String("PUT"), aws.String("POST")},
								AllowedOrigins: []*string{aws.String("https://www.example.com")},
								MaxAgeSeconds:  aws.Int64(3000),
							},
						},
					),
				),
			},
		},
	})
}

func TestAccDigitalOceanBucket_WithCors(t *testing.T) {
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDigitalOceanBucketConfigWithCORSUpdate, ri),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists("digitalocean_spaces_bucket.bucket"),
					testAccCheckDigitalOceanBucketCors(
						"digitalocean_spaces_bucket.bucket",
						[]*s3.CORSRule{
							{
								AllowedHeaders: []*string{aws.String("*")},
								AllowedMethods: []*string{aws.String("PUT"), aws.String("POST")},
								AllowedOrigins: []*string{aws.String("https://www.example.com")},
								MaxAgeSeconds:  aws.Int64(3000),
							},
						},
					),
				),
			},
		},
	})
}

func TestAccDigitalOceanBucket_WithMultipleCorsRules(t *testing.T) {
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDigitalOceanBucketConfigWithMultiCORS, ri),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists("digitalocean_spaces_bucket.bucket"),
					testAccCheckDigitalOceanBucketCors(
						"digitalocean_spaces_bucket.bucket",
						[]*s3.CORSRule{
							{
								AllowedHeaders: []*string{aws.String("*")},
								AllowedMethods: []*string{aws.String("GET")},
								AllowedOrigins: []*string{aws.String("*")},
								MaxAgeSeconds:  aws.Int64(3000),
							},
							{
								AllowedHeaders: []*string{aws.String("*")},
								AllowedMethods: []*string{aws.String("PUT"), aws.String("DELETE"), aws.String("POST")},
								AllowedOrigins: []*string{aws.String("https://www.example.com")},
								MaxAgeSeconds:  aws.Int64(3000),
							},
						},
					),
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

func TestAccDigitalOceanBucket_Versioning(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "digitalocean_spaces_bucket.bucket"

	makeConfig := func(includeClause, versioning bool) string {
		versioningClause := ""
		if includeClause {
			versioningClause = fmt.Sprintf(`
  versioning {
    enabled = %v
  }
`, versioning)
		}
		return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "bucket" {
  name = "tf-test-bucket-%d"
  region = "ams3"
%s
}
`, rInt, versioningClause)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				// No versioning configured.
				Config: makeConfig(false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists(resourceName),
					testAccCheckDigitalOceanBucketVersioning(
						resourceName, ""),
				),
			},
			{
				// Enable versioning
				Config: makeConfig(true, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists(resourceName),
					testAccCheckDigitalOceanBucketVersioning(
						resourceName, s3.BucketVersioningStatusEnabled),
				),
			},
			{
				// Explicitly disable versioning
				Config: makeConfig(true, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists(resourceName),
					testAccCheckDigitalOceanBucketVersioning(
						resourceName, s3.BucketVersioningStatusSuspended),
				),
			},
			{
				// Re-enable versioning
				Config: makeConfig(true, true),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists(resourceName),
					testAccCheckDigitalOceanBucketVersioning(
						resourceName, s3.BucketVersioningStatusEnabled),
				),
			},
			{
				// Remove the clause completely. Should disable versioning.
				Config: makeConfig(false, false),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists(resourceName),
					testAccCheckDigitalOceanBucketVersioning(
						resourceName, s3.BucketVersioningStatusSuspended),
				),
			},
		},
	})
}

func TestAccDigitalOceanSpacesBucket_LifecycleBasic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "digitalocean_spaces_bucket.bucket"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanSpacesBucketConfigWithLifecycle(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.0.id", "id1"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.0.prefix", "path1/"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.0.expiration.2613713285.days", "365"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.0.expiration.2613713285.date", ""),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.0.expiration.2613713285.expired_object_delete_marker", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.1.id", "id2"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.1.prefix", "path2/"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.1.expiration.2855832418.date", "2016-01-12"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.1.expiration.2855832418.days", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.1.expiration.2855832418.expired_object_delete_marker", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.2.id", "id3"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.2.prefix", "path3/"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.2.abort_incomplete_multipart_upload_days", "30"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("nyc3,tf-test-bucket-%d", rInt),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"force_destroy", "acl"},
			},
			{
				Config: testAccDigitalOceanSpacesBucketConfigWithVersioningLifecycle(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.0.id", "id1"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.0.prefix", "path1/"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.0.enabled", "true"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.1.id", "id2"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.1.prefix", "path2/"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.1.enabled", "false"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.1.noncurrent_version_expiration.80908210.days", "365"),
				),
			},
			{
				Config: testAccDigitalOceanBucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists(resourceName),
				),
			},
		},
	})
}

func TestAccDigitalOceanSpacesBucket_LifecycleExpireMarkerOnly(t *testing.T) {
	rInt := acctest.RandInt()
	resourceName := "digitalocean_spaces_bucket.bucket"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanSpacesBucketConfigWithLifecycleExpireMarker(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.0.id", "id1"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.0.prefix", "path1/"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.0.expiration.3591068768.days", "0"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.0.expiration.3591068768.date", ""),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.0.expiration.3591068768.expired_object_delete_marker", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("nyc3,tf-test-bucket-%d", rInt),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"force_destroy", "acl"},
			},
			{
				Config: testAccDigitalOceanBucketConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists(resourceName),
				),
			},
		},
	})
}

// Reference: https://github.com/terraform-providers/terraform-provider-aws/issues/11420
func TestAccDigitalOceanSpacesBucket_LifecycleRule_Expiration_EmptyConfigurationBlock(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "digitalocean_spaces_bucket.bucket"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanSpacesBucketConfigLifecycleRuleExpirationEmptyConfigurationBlock(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists(resourceName),
				),
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

func testAccCheckDigitalOceanBucketCors(n string, corsRules []*s3.CORSRule) resource.TestCheckFunc {
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

		out, err := svc.GetBucketCors(&s3.GetBucketCorsInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return fmt.Errorf("GetBucketCors error: %v", err)
		}

		if !reflect.DeepEqual(out.CORSRules, corsRules) {
			return fmt.Errorf("bad error cors rule, expected: %v, got %v", corsRules, out.CORSRules)
		}

		return nil
	}
}

func testAccCheckDigitalOceanBucketVersioning(n string, versioningStatus string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]

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

		out, err := svc.GetBucketVersioning(&s3.GetBucketVersioningInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			return fmt.Errorf("GetBucketVersioning error: %v", err)
		}

		if v := out.Status; v == nil {
			if versioningStatus != "" {
				return fmt.Errorf("bad error versioning status, found nil, expected: %s", versioningStatus)
			}
		} else {
			if *v != versioningStatus {
				return fmt.Errorf("bad error versioning status, expected: %s, got %s", versioningStatus, *v)
			}
		}

		return nil
	}
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
	name   = "tf-test-bucket-%d"
	region = "sfo2"
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

var testAccDigitalOceanBucketConfigWithCORS = `
resource "digitalocean_spaces_bucket" "bucket" {
	name = "tf-test-bucket-%d"
}
`

var testAccDigitalOceanBucketConfigWithCORSUpdate = `
resource "digitalocean_spaces_bucket" "bucket" {
	name = "tf-test-bucket-%d"
	cors_rule {
			allowed_headers = ["*"]
			allowed_methods = ["PUT","POST"]
			allowed_origins = ["https://www.example.com"]
			max_age_seconds = 3000
	}
}
`

var testAccDigitalOceanBucketConfigWithMultiCORS = `
resource "digitalocean_spaces_bucket" "bucket" {
	name = "tf-test-bucket-%d"

	cors_rule {
			allowed_headers = ["*"]
			allowed_methods = ["GET"]
			allowed_origins = ["*"]
			max_age_seconds = 3000
	}

	cors_rule {
			allowed_headers = ["*"]
			allowed_methods = ["PUT", "DELETE", "POST"]
			allowed_origins = ["https://www.example.com"]
			max_age_seconds = 3000
	}
}
`

func testAccDigitalOceanSpacesBucketConfigWithLifecycle(randInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "bucket" {
  name = "tf-test-bucket-%d"
  acl  = "private"

  lifecycle_rule {
    id      = "id1"
    prefix  = "path1/"
    enabled = true

    expiration {
      days = 365
    }
  }

  lifecycle_rule {
    id      = "id2"
    prefix  = "path2/"
    enabled = true

    expiration {
      date = "2016-01-12"
    }
  }

  lifecycle_rule {
    id     = "id3"
    prefix = "path3/"
    enabled = true

    abort_incomplete_multipart_upload_days = 30
  }
}
`, randInt)
}

func testAccDigitalOceanSpacesBucketConfigWithLifecycleExpireMarker(randInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "bucket" {
  name = "tf-test-bucket-%d"
  acl  = "private"

  lifecycle_rule {
    id      = "id1"
    prefix  = "path1/"
    enabled = true

    expiration {
      expired_object_delete_marker = "true"
    }
  }
}
`, randInt)
}

func testAccDigitalOceanSpacesBucketConfigWithVersioningLifecycle(randInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "bucket" {
  name = "tf-test-bucket-%d"
  acl  = "private"

  versioning {
    enabled = false
  }

  lifecycle_rule {
    id      = "id1"
    prefix  = "path1/"
    enabled = true

    noncurrent_version_expiration {
      days = 365
    }
 }

  lifecycle_rule {
    id      = "id2"
    prefix  = "path2/"
    enabled = false

    noncurrent_version_expiration {
      days = 365
    }
  }
}
`, randInt)
}

func testAccDigitalOceanSpacesBucketConfigLifecycleRuleExpirationEmptyConfigurationBlock(rName string) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "bucket" {
  name = %[1]q

  lifecycle_rule {
    enabled = true
    id      = "id1"

    expiration {}
  }
}
`, rName)
}
