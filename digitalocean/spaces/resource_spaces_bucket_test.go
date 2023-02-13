package spaces_test

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/spaces"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanBucket_basic(t *testing.T) {
	expectedRegion := "ams3"
	expectedBucketName := acceptance.RandomTestName()
	expectBucketURN := fmt.Sprintf("do:space:%s", expectedBucketName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		/*
			IDRefreshName:   "digitalocean_spaces_bucket.bucket",
			IDRefreshIgnore: []string{"force_destroy"},
		*/
		CheckDestroy: testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanBucketConfig(expectedBucketName),
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
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanBucketConfigWithRegion(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists("digitalocean_spaces_bucket.bucket"),
					resource.TestCheckResourceAttr("digitalocean_spaces_bucket.bucket", "region", "ams3"),
				),
			},
		},
	})
}

func TestAccDigitalOceanBucket_UpdateAcl(t *testing.T) {
	name := acceptance.RandomTestName()
	preConfig := fmt.Sprintf(testAccDigitalOceanBucketConfigWithACL, name)
	postConfig := fmt.Sprintf(testAccDigitalOceanBucketConfigWithACLUpdate, name)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanBucketDestroy,
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
	name := acceptance.RandomTestName()
	preConfig := fmt.Sprintf(testAccDigitalOceanBucketConfigWithCORS, name)
	postConfig := fmt.Sprintf(testAccDigitalOceanBucketConfigWithCORSUpdate, name)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanBucketDestroy,
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
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDigitalOceanBucketConfigWithCORSUpdate, name),
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
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDigitalOceanBucketConfigWithMultiCORS, name),
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
	name := acceptance.RandomTestName()
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanBucketDestroyedConfig(name),
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
	name := acceptance.RandomTestName()
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
  name   = "%s"
  region = "ams3"
%s
}
`, name, versioningClause)
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanBucketDestroy,
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
	name := acceptance.RandomTestName()
	resourceName := "digitalocean_spaces_bucket.bucket"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanSpacesBucketConfigWithLifecycle(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.0.id", "id1"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.0.prefix", "path1/"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "lifecycle_rule.0.expiration.*",
						map[string]string{
							"days":                         "365",
							"date":                         "",
							"expired_object_delete_marker": "false",
						}),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.1.id", "id2"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.1.prefix", "path2/"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "lifecycle_rule.1.expiration.*",
						map[string]string{
							"days":                         "",
							"date":                         "2016-01-12",
							"expired_object_delete_marker": "",
						}),
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
				ImportStateId:     fmt.Sprintf("ams3,%s", name),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"force_destroy", "acl"},
			},
			{
				Config: testAccDigitalOceanSpacesBucketConfigWithVersioningLifecycle(name),
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
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "lifecycle_rule.1.noncurrent_version_expiration.*",
						map[string]string{"days": "365"}),
				),
			},
			{
				Config: testAccDigitalOceanBucketConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists(resourceName),
				),
			},
		},
	})
}

func TestAccDigitalOceanSpacesBucket_LifecycleExpireMarkerOnly(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceName := "digitalocean_spaces_bucket.bucket"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanSpacesBucketConfigWithLifecycleExpireMarker(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists(resourceName),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.0.id", "id1"),
					resource.TestCheckResourceAttr(
						resourceName, "lifecycle_rule.0.prefix", "path1/"),
					resource.TestCheckTypeSetElemNestedAttrs(
						resourceName, "lifecycle_rule.0.expiration.*",
						map[string]string{
							"days":                         "0",
							"date":                         "",
							"expired_object_delete_marker": "true",
						}),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("ams3,%s", name),
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"force_destroy", "acl"},
			},
			{
				Config: testAccDigitalOceanBucketConfig(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanBucketExists(resourceName),
				),
			},
		},
	})
}

func TestAccDigitalOceanSpacesBucket_RegionError(t *testing.T) {
	badRegion := "ny2"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanBucketDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "bucket" {
  name   = "%s"
  region = "%s"
}`, acceptance.RandomTestName(), badRegion),
				ExpectError: regexp.MustCompile(`expected region to be one of`),
			},
		},
	})
}

func testAccGetS3ConnForSpacesBucket(rs *terraform.ResourceState) (*s3.S3, error) {
	rawRegion := ""
	if actualRegion, ok := rs.Primary.Attributes["region"]; ok {
		rawRegion = actualRegion
	}
	region := spaces.NormalizeRegion(rawRegion)

	spacesAccessKeyId := os.Getenv("SPACES_ACCESS_KEY_ID")
	if spacesAccessKeyId == "" {
		return nil, fmt.Errorf("SPACES_ACCESS_KEY_ID must be set")
	}

	spacesSecretAccessKey := os.Getenv("SPACES_SECRET_ACCESS_KEY")
	if spacesSecretAccessKey == "" {
		return nil, fmt.Errorf("SPACES_SECRET_ACCESS_KEY must be set")
	}

	sesh, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(spacesAccessKeyId, spacesSecretAccessKey, "")},
	)
	if err != nil {
		return nil, fmt.Errorf("Unable to create S3 session (region=%s): %v", region, err)
	}

	svc := s3.New(sesh, &aws.Config{
		Endpoint: aws.String(fmt.Sprintf("https://%s.digitaloceanspaces.com", region))},
	)

	return svc, nil
}

func testAccCheckDigitalOceanBucketDestroy(s *terraform.State) error {
	return testAccCheckDigitalOceanBucketDestroyWithProvider(s, acceptance.TestAccProvider)
}

func testAccCheckDigitalOceanBucketDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_spaces_bucket" {
			continue
		}

		if rs.Primary.ID == "" {
			continue
		}

		svc, err := testAccGetS3ConnForSpacesBucket(rs)
		if err != nil {
			return fmt.Errorf("Unable to create S3 client: %v", err)
		}

		_, err = svc.DeleteBucket(&s3.DeleteBucketInput{
			Bucket: aws.String(rs.Primary.ID),
		})
		if err != nil {
			if spaces.IsAWSErr(err, s3.ErrCodeNoSuchBucket, "") {
				return nil
			}
			return err
		}
	}
	return nil
}

func testAccCheckDigitalOceanBucketExists(n string) resource.TestCheckFunc {
	return testAccCheckDigitalOceanBucketExistsWithProvider(n, func() *schema.Provider { return acceptance.TestAccProvider })
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

		svc, err := testAccGetS3ConnForSpacesBucket(rs)
		if err != nil {
			return fmt.Errorf("Unable to create S3 client: %v", err)
		}

		_, err = svc.HeadBucket(&s3.HeadBucketInput{
			Bucket: aws.String(rs.Primary.ID),
		})

		if err != nil {
			if spaces.IsAWSErr(err, s3.ErrCodeNoSuchBucket, "") {
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

		svc, err := testAccGetS3ConnForSpacesBucket(rs)
		if err != nil {
			return fmt.Errorf("Unable to create S3 client: %v", err)
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

		svc, err := testAccGetS3ConnForSpacesBucket(rs)
		if err != nil {
			return fmt.Errorf("Unable to create S3 client: %v", err)
		}

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

		svc, err := testAccGetS3ConnForSpacesBucket(rs)
		if err != nil {
			return fmt.Errorf("Unable to create S3 client: %v", err)
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

func testAccDigitalOceanBucketConfig(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "bucket" {
  name   = "%s"
  acl    = "public-read"
  region = "ams3"
}
`, name)
}

func testAccDigitalOceanBucketDestroyedConfig(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "bucket" {
  name = "%s"
  acl  = "public-read"
}
`, name)
}

func testAccDigitalOceanBucketConfigWithRegion(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "bucket" {
  name   = "%s"
  region = "ams3"
}
`, name)
}

func testAccDigitalOceanBucketConfigImport(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "bucket" {
  name   = "%s"
  region = "sfo3"
}
`, name)
}

var testAccDigitalOceanBucketConfigWithACL = `
resource "digitalocean_spaces_bucket" "bucket" {
  name = "%s"
  acl  = "public-read"
}
`

var testAccDigitalOceanBucketConfigWithACLUpdate = `
resource "digitalocean_spaces_bucket" "bucket" {
  name = "%s"
  acl  = "private"
}
`

var testAccDigitalOceanBucketConfigWithCORS = `
resource "digitalocean_spaces_bucket" "bucket" {
  name = "%s"
}
`

var testAccDigitalOceanBucketConfigWithCORSUpdate = `
resource "digitalocean_spaces_bucket" "bucket" {
  name = "%s"
  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST"]
    allowed_origins = ["https://www.example.com"]
    max_age_seconds = 3000
  }
}
`

var testAccDigitalOceanBucketConfigWithMultiCORS = `
resource "digitalocean_spaces_bucket" "bucket" {
  name = "%s"

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

func testAccDigitalOceanSpacesBucketConfigWithLifecycle(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "bucket" {
  name   = "%s"
  acl    = "private"
  region = "ams3"

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
    id      = "id3"
    prefix  = "path3/"
    enabled = true

    abort_incomplete_multipart_upload_days = 30
  }
}
`, name)
}

func testAccDigitalOceanSpacesBucketConfigWithLifecycleExpireMarker(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "bucket" {
  name   = "%s"
  acl    = "private"
  region = "ams3"

  lifecycle_rule {
    id      = "id1"
    prefix  = "path1/"
    enabled = true

    expiration {
      expired_object_delete_marker = "true"
    }
  }
}
`, name)
}

func testAccDigitalOceanSpacesBucketConfigWithVersioningLifecycle(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "bucket" {
  name   = "%s"
  acl    = "private"
  region = "ams3"

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
`, name)
}
