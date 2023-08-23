package spaces_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	testAccDigitalOceanSpacesBucketCorsConfiguration_TestRegion = "nyc3"
)

func TestAccDigitalOceanBucketCorsConfiguration_basic(t *testing.T) {
	name := acceptance.RandomTestName()
	ctx := context.Background()
	region := testAccDigitalOceanSpacesBucketCorsConfiguration_TestRegion
	resourceName := "digitalocean_spaces_bucket_cors_configuration.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanSpacesBucketCorsConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpacesBucketCORSConfigurationConfig_basic(name, region, "https://www.example.com"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketCorsConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "cors_rule.*", map[string]string{
						"allowed_methods.#": "1",
						"allowed_origins.#": "1",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_rule.*.allowed_methods.*", "PUT"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_rule.*.allowed_origins.*", "https://www.example.com"),
				),
			},
		},
	})
}

func TestAccS3BucketCorsConfiguration_SingleRule(t *testing.T) {
	resourceName := "digitalocean_spaces_bucket_cors_configuration.test"
	rName := acceptance.RandomTestName()
	ctx := context.Background()
	region := testAccDigitalOceanSpacesBucketCorsConfiguration_TestRegion

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanSpacesBucketCorsConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpacesBucketCORSConfigurationConfig_completeSingleRule(rName, region, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketCorsConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "bucket", "digitalocean_spaces_bucket.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "cors_rule.*", map[string]string{
						"allowed_headers.#": "1",
						"allowed_methods.#": "3",
						"allowed_origins.#": "1",
						"expose_headers.#":  "1",
						"id":                rName,
						"max_age_seconds":   "3000",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_rule.*.allowed_headers.*", "*"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_rule.*.allowed_methods.*", "DELETE"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_rule.*.allowed_methods.*", "POST"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_rule.*.allowed_methods.*", "PUT"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_rule.*.allowed_origins.*", "https://www.example.com"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_rule.*.expose_headers.*", "ETag"),
				),
			},
		},
	})
}

func TestAccS3BucketCorsConfiguration_MultipleRules(t *testing.T) {
	resourceName := "digitalocean_spaces_bucket_cors_configuration.test"
	rName := acceptance.RandomTestName()
	ctx := context.Background()
	region := testAccDigitalOceanSpacesBucketCorsConfiguration_TestRegion

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanSpacesBucketCorsConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpacesBucketCORSConfigurationConfig_multipleRules(rName, region, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketCorsConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "bucket", "digitalocean_spaces_bucket.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "cors_rule.*", map[string]string{
						"allowed_headers.#": "1",
						"allowed_methods.#": "3",
						"allowed_origins.#": "1",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_rule.*.allowed_headers.*", "*"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_rule.*.allowed_methods.*", "DELETE"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_rule.*.allowed_methods.*", "POST"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_rule.*.allowed_methods.*", "PUT"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_rule.*.allowed_origins.*", "https://www.example.com"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "cors_rule.*", map[string]string{
						"allowed_methods.#": "1",
						"allowed_origins.#": "1",
					}),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_rule.*.allowed_methods.*", "GET"),
					resource.TestCheckTypeSetElemAttr(resourceName, "cors_rule.*.allowed_origins.*", "*"),
				),
			},
		},
	})
}

func TestAccS3BucketCorsConfiguration_update(t *testing.T) {
	resourceName := "digitalocean_spaces_bucket_cors_configuration.test"
	rName := acceptance.RandomTestName()
	ctx := context.Background()
	region := testAccDigitalOceanSpacesBucketCorsConfiguration_TestRegion

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanSpacesBucketCorsConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSpacesBucketCORSConfigurationConfig_completeSingleRule(rName, region, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketCorsConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "bucket", "digitalocean_spaces_bucket.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "cors_rule.*", map[string]string{
						"allowed_headers.#": "1",
						"allowed_methods.#": "3",
						"allowed_origins.#": "1",
						"expose_headers.#":  "1",
						"id":                rName,
						"max_age_seconds":   "3000",
					}),
				),
			},
			{
				Config: testAccSpacesBucketCORSConfigurationConfig_multipleRules(rName, region, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketCorsConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttrPair(resourceName, "bucket", "digitalocean_spaces_bucket.test", "id"),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.#", "2"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "cors_rule.*", map[string]string{
						"allowed_headers.#": "1",
						"allowed_methods.#": "3",
						"allowed_origins.#": "1",
					}),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "cors_rule.*", map[string]string{
						"allowed_methods.#": "1",
						"allowed_origins.#": "1",
					}),
				),
			},
			{
				Config: testAccSpacesBucketCORSConfigurationConfig_basic(rName, region, "https://www.example.com"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketCorsConfigurationExists(ctx, resourceName),
					resource.TestCheckResourceAttr(resourceName, "cors_rule.#", "1"),
					resource.TestCheckTypeSetElemNestedAttrs(resourceName, "cors_rule.*", map[string]string{
						"allowed_methods.#": "1",
						"allowed_origins.#": "1",
					}),
				),
			},
		},
	})
}

func testAccGetS3CorsConfigurationConn() (*s3.S3, error) {
	client, err := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).SpacesClient(testAccDigitalOceanSpacesBucketCorsConfiguration_TestRegion)
	if err != nil {
		return nil, err
	}

	s3conn := s3.New(client)

	return s3conn, nil
}

func testAccCheckDigitalOceanSpacesBucketCorsConfigurationExists(ctx context.Context, resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not Found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource (%s) ID not set", resourceName)
		}

		s3conn, err := testAccGetS3CorsConfigurationConn()
		if err != nil {
			return err
		}

		response, err := s3conn.GetBucketCorsWithContext(context.Background(),
			&s3.GetBucketCorsInput{
				Bucket: aws.String(rs.Primary.Attributes["bucket"]),
			})
		if err != nil {
			return fmt.Errorf("S3Bucket CORs error: %s", err)
		}

		if len(response.CORSRules) == 0 {
			return fmt.Errorf("S3 Bucket CORS configuration (%s) not found", rs.Primary.ID)
		}
		return nil
	}
}

func testAccCheckDigitalOceanSpacesBucketCorsConfigurationDestroy(s *terraform.State) error {
	s3conn, err := testAccGetS3CorsConfigurationConn()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "digitalocean_spaces_bucket_cors_configuration":
			_, err := s3conn.GetBucketCorsWithContext(context.Background(), &s3.GetBucketCorsInput{
				Bucket: aws.String(rs.Primary.Attributes["bucket"]),
			})
			if err == nil {
				return fmt.Errorf("Spaces Bucket Cors Configuration still exists: %s", rs.Primary.ID)
			}

		case "digitalocean_spaces_bucket":
			_, err = s3conn.HeadBucket(&s3.HeadBucketInput{
				Bucket: aws.String(rs.Primary.ID),
			})
			if err == nil {
				return fmt.Errorf("Spaces Bucket still exists: %s", rs.Primary.ID)
			}

		default:
			continue
		}
	}

	return nil
}

func testAccSpacesBucketCORSConfigurationConfig_basic(rName string, region string, origin string) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "foobar" {
  name   = "%s"
  region = "%s"
}

resource "digitalocean_spaces_bucket_cors_configuration" "test" {
  bucket = digitalocean_spaces_bucket.foobar.id
  region = "nyc3"

  cors_rule {
    allowed_methods = ["PUT"]
    allowed_origins = ["%s"]
  }
}
`, rName, region, origin)
}

func testAccSpacesBucketCORSConfigurationConfig_completeSingleRule(rName string, region string, Name string) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "foobar" {
  name          = "%s"
  region        = "%s"
  force_destroy = true
}

resource "digitalocean_spaces_bucket_cors_configuration" "test" {
  bucket = digitalocean_spaces_bucket.foobar.id
  region = "nyc3"

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST", "DELETE"]
    allowed_origins = ["https://www.example.com"]
    expose_headers  = ["ETag"]
    id              = "%s"
    max_age_seconds = 3000
  }
}
`, rName, region, Name)
}

func testAccSpacesBucketCORSConfigurationConfig_multipleRules(rName string, region string, Name string) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "foobar" {
  name          = "%s"
  region        = "%s"
  force_destroy = true
}

resource "digitalocean_spaces_bucket_cors_configuration" "test" {
  bucket = digitalocean_spaces_bucket.foobar.id
  region = "nyc3"

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST", "DELETE"]
    allowed_origins = ["https://www.example.com"]
    expose_headers  = ["ETag"]
    id              = "%s"
    max_age_seconds = 3000
  }

  cors_rule {
    allowed_methods = ["GET"]
    allowed_origins = ["*"]
  }
}
`, rName, region, Name)
}
