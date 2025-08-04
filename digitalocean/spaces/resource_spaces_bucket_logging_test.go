package spaces_test

import (
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
	accessLoggingTestRegion = "nyc3"
)

func TestAccDigitalOceanSpacesBucketLogging_basic(t *testing.T) {
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanSpacesBucketLoggingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanSpacesBucketLogging(name, "logs/"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_spaces_bucket_logging.example", "region", accessLoggingTestRegion),
					resource.TestCheckResourceAttr("digitalocean_spaces_bucket_logging.example", "bucket", name+"-source"),
					resource.TestCheckResourceAttr("digitalocean_spaces_bucket_logging.example", "target_bucket", name+"-target"),
					resource.TestCheckResourceAttr("digitalocean_spaces_bucket_logging.example", "target_prefix", "logs/"),
				),
			},
		},
	})
}

func TestAccDigitalOceanSpacesBucketLogging_update(t *testing.T) {
	name := acceptance.RandomTestName()
	initialPrefix := "logs/"
	updatedPrefix := "updated/"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanSpacesBucketLoggingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanSpacesBucketLogging(name, initialPrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_spaces_bucket_logging.example", "region", accessLoggingTestRegion),
					resource.TestCheckResourceAttr("digitalocean_spaces_bucket_logging.example", "bucket", name+"-source"),
					resource.TestCheckResourceAttr("digitalocean_spaces_bucket_logging.example", "target_bucket", name+"-target"),
					resource.TestCheckResourceAttr("digitalocean_spaces_bucket_logging.example", "target_prefix", initialPrefix),
				),
			},
			{
				Config: testAccDigitalOceanSpacesBucketLogging(name, updatedPrefix),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_spaces_bucket_logging.example", "region", accessLoggingTestRegion),
					resource.TestCheckResourceAttr("digitalocean_spaces_bucket_logging.example", "bucket", name+"-source"),
					resource.TestCheckResourceAttr("digitalocean_spaces_bucket_logging.example", "target_bucket", name+"-target"),
					resource.TestCheckResourceAttr("digitalocean_spaces_bucket_logging.example", "target_prefix", updatedPrefix),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanSpacesBucketLoggingDestroy(s *terraform.State) error {
	s3conn, err := testAccGetS3LoggingConn()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "digitalocean_spaces_bucket_logging":
			_, err := s3conn.GetBucketLogging(&s3.GetBucketLoggingInput{
				Bucket: aws.String(rs.Primary.Attributes["bucket"]),
			})
			if err == nil {
				return fmt.Errorf("Spaces Bucket logging still exists: %s", rs.Primary.ID)
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

func testAccGetS3LoggingConn() (*s3.S3, error) {
	client, err := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).SpacesClient(accessLoggingTestRegion)
	if err != nil {
		return nil, err
	}

	s3conn := s3.New(client)

	return s3conn, nil
}

func testAccDigitalOceanSpacesBucketLogging(name, prefix string) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "source" {
  region        = "%s"
  name          = "%s-source"
  force_destroy = true
}

resource "digitalocean_spaces_bucket" "target" {
  region        = "%s"
  name          = "%s-target"
  force_destroy = true
}

resource "digitalocean_spaces_bucket_logging" "example" {
  region = "%s"
  bucket = digitalocean_spaces_bucket.source.id

  target_bucket = digitalocean_spaces_bucket.target.id
  target_prefix = "%s"
}
`, accessLoggingTestRegion, name, accessLoggingTestRegion, name, accessLoggingTestRegion, prefix)
}
