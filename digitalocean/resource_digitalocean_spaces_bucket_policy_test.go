package digitalocean

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	testAccDigitalOceanSpacesBucketPolicy_TestRegion = "nyc3"
)

func TestAccDigitalOceanBucketPolicy_basic(t *testing.T) {
	randInt := acctest.RandInt()

	bucketPolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"s3:*","Resource":"*"}]}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanSpacesBucketPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanSpacesBucketPolicy(randInt, bucketPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketPolicy("digitalocean_spaces_bucket_policy.policy", bucketPolicy),
					resource.TestCheckResourceAttr("digitalocean_spaces_bucket_policy.policy", "region", testAccDigitalOceanSpacesBucketPolicy_TestRegion),
				),
			},
		},
	})
}

func TestAccDigitalOceanBucketPolicy_update(t *testing.T) {
	randInt := acctest.RandInt()

	initialBucketPolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"s3:*","Resource":"*"}]}`
	updatedBucketPolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"s3:*","Resource":"*"},{"Effect":"Allow","Principal":"*","Action":"s3:GetObject","Resource":"*"}]}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanSpacesBucketPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanSpacesBucketPolicy(randInt, initialBucketPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketPolicy("digitalocean_spaces_bucket_policy.policy", initialBucketPolicy),
					resource.TestCheckResourceAttr("digitalocean_spaces_bucket_policy.policy", "region", testAccDigitalOceanSpacesBucketPolicy_TestRegion),
				),
			},
			{
				Config: testAccDigitalOceanSpacesBucketPolicy(randInt, updatedBucketPolicy),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketPolicy("digitalocean_spaces_bucket_policy.policy", updatedBucketPolicy),
					resource.TestCheckResourceAttr("digitalocean_spaces_bucket_policy.policy", "region", testAccDigitalOceanSpacesBucketPolicy_TestRegion),
				),
			},
		},
	})
}

func TestAccDigitalOceanBucketPolicy_invalidJson(t *testing.T) {
	randInt := acctest.RandInt()

	bucketPolicy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":"*","Action":"s3:*","Resource":"*"}}`
	expectError := regexp.MustCompile(`"policy" contains an invalid JSON`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanSpacesBucketPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDigitalOceanSpacesBucketPolicy(randInt, bucketPolicy),
				ExpectError: expectError,
			},
		},
	})
}

func TestAccDigitalOceanBucketPolicy_emptyPolicy(t *testing.T) {
	randInt := acctest.RandInt()

	expectError := regexp.MustCompile(`policy must not be empty`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanSpacesBucketPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDigitalOceanSpacesBucketEmptyPolicy(randInt),
				ExpectError: expectError,
			},
		},
	})
}

func TestAccDigitalOceanBucketPolicy_unknownBucket(t *testing.T) {
	expectError := regexp.MustCompile(`bucket 'unknown' does not exist`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanSpacesBucketPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccDigitalOceanSpacesBucketUnknownBucket(),
				ExpectError: expectError,
			},
		},
	})
}

func testAccGetS3PolicyConn() (*s3.S3, error) {
	client, err := testAccProvider.Meta().(*CombinedConfig).spacesClient(testAccDigitalOceanSpacesBucketPolicy_TestRegion)
	if err != nil {
		return nil, err
	}

	s3conn := s3.New(client)

	return s3conn, nil
}

func testAccCheckDigitalOceanSpacesBucketPolicy(n string, expectedPolicy string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No S3 Bucket Policy ID is set")
		}

		s3conn, err := testAccGetS3PolicyConn()
		if err != nil {
			return err
		}

		response, err := s3conn.GetBucketPolicy(
			&s3.GetBucketPolicyInput{
				Bucket: aws.String(rs.Primary.Attributes["bucket"]),
			})
		if err != nil {
			return fmt.Errorf("S3Bucket policy error: %s", err)
		}

		actualPolicy := aws.StringValue(response.Policy)
		equivalent := compareSpacesBucketPolicy(expectedPolicy, actualPolicy)
		if !equivalent {
			return fmt.Errorf("Expected policy to be '%v', got '%v'", expectedPolicy, actualPolicy)
		}
		return nil
	}
}

func testAccCheckDigitalOceanSpacesBucketPolicyDestroy(s *terraform.State) error {
	s3conn, err := testAccGetS3PolicyConn()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "digitalocean_spaces_bucket_policy":
			_, err := s3conn.GetBucketPolicy(&s3.GetBucketPolicyInput{
				Bucket: aws.String(rs.Primary.Attributes["bucket"]),
			})
			if err == nil {
				return fmt.Errorf("Spaces Bucket policy still exists: %s", rs.Primary.ID)
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

func testAccDigitalOceanSpacesBucketPolicy(randInt int, policy string) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "policy_bucket" {
  region = "%s"
  name   = "tf-policy-test-bucket-%d"
  force_destroy = true
}

resource "digitalocean_spaces_bucket_policy" "policy" {
  region = digitalocean_spaces_bucket.policy_bucket.region
  bucket = digitalocean_spaces_bucket.policy_bucket.name
  policy = <<EOF
%s
EOF
}

`, testAccDigitalOceanSpacesBucketPolicy_TestRegion, randInt, policy)
}

func testAccDigitalOceanSpacesBucketEmptyPolicy(randInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "policy_bucket" {
  region = "%s"
  name   = "tf-policy-test-bucket-%d"
  force_destroy = true
}

resource "digitalocean_spaces_bucket_policy" "policy" {
  region = digitalocean_spaces_bucket.policy_bucket.region
  bucket = digitalocean_spaces_bucket.policy_bucket.name
  policy = ""
}

`, testAccDigitalOceanSpacesBucketPolicy_TestRegion, randInt)
}

func testAccDigitalOceanSpacesBucketUnknownBucket() string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket_policy" "policy" {
  region = "%s"
  bucket = "unknown"
  policy = "{}"
}

`, testAccDigitalOceanSpacesBucketPolicy_TestRegion)
}
