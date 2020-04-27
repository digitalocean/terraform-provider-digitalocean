package digitalocean

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceDigitalOceanSpacesBucketObject_basic(t *testing.T) {
	rInt := acctest.RandInt()
	resourceOnlyConf, conf := testAccDataSourceDigitalOceanSpacesObjectConfig_basic(rInt)

	var rObj s3.GetObjectOutput
	var dsObj s3.GetObjectOutput

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: resourceOnlyConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketObjectExists("aws_s3_bucket_object.object", &rObj),
				),
			},
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesObjectDataSourceExists("data.aws_s3_bucket_object.obj", &dsObj),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "content_length", "11"),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "content_type", "binary/octet-stream"),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "etag", "b10a8db164e0754105b7a99be72e3fe5"),
					resource.TestMatchResourceAttr("data.aws_s3_bucket_object.obj", "last_modified",
						regexp.MustCompile("^[a-zA-Z]{3}, [0-9]+ [a-zA-Z]+ [0-9]{4} [0-9:]+ [A-Z]+$")),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "object_lock_legal_hold_status", ""),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "object_lock_mode", ""),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "object_lock_retain_until_date", ""),
					resource.TestCheckNoResourceAttr("data.aws_s3_bucket_object.obj", "body"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanSpacesBucketObject_readableBody(t *testing.T) {
	rInt := acctest.RandInt()
	resourceOnlyConf, conf := testAccDataSourceDigitalOceanSpacesObjectConfig_readableBody(rInt)

	var rObj s3.GetObjectOutput
	var dsObj s3.GetObjectOutput

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: resourceOnlyConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketObjectExists("aws_s3_bucket_object.object", &rObj),
				),
			},
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesObjectDataSourceExists("data.aws_s3_bucket_object.obj", &dsObj),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "content_length", "3"),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "content_type", "text/plain"),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "etag", "a6105c0a611b41b08f1209506350279e"),
					resource.TestMatchResourceAttr("data.aws_s3_bucket_object.obj", "last_modified",
						regexp.MustCompile("^[a-zA-Z]{3}, [0-9]+ [a-zA-Z]+ [0-9]{4} [0-9:]+ [A-Z]+$")),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "object_lock_legal_hold_status", ""),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "object_lock_mode", ""),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "object_lock_retain_until_date", ""),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "body", "yes"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanSpacesBucketObject_allParams(t *testing.T) {
	rInt := acctest.RandInt()
	resourceOnlyConf, conf := testAccDataSourceDigitalOceanSpacesObjectConfig_allParams(rInt)

	var rObj s3.GetObjectOutput
	var dsObj s3.GetObjectOutput

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: resourceOnlyConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketObjectExists("aws_s3_bucket_object.object", &rObj),
				),
			},
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesObjectDataSourceExists("data.aws_s3_bucket_object.obj", &dsObj),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "content_length", "21"),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "content_type", "application/unknown"),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "etag", "723f7a6ac0c57b445790914668f98640"),
					resource.TestMatchResourceAttr("data.aws_s3_bucket_object.obj", "last_modified",
						regexp.MustCompile("^[a-zA-Z]{3}, [0-9]+ [a-zA-Z]+ [0-9]{4} [0-9:]+ [A-Z]+$")),
					resource.TestMatchResourceAttr("data.aws_s3_bucket_object.obj", "version_id", regexp.MustCompile("^.{32}$")),
					resource.TestCheckNoResourceAttr("data.aws_s3_bucket_object.obj", "body"),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "cache_control", "no-cache"),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "content_disposition", "attachment"),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "content_encoding", "identity"),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "content_language", "en-GB"),
					// Encryption is off
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "server_side_encryption", ""),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "sse_kms_key_id", ""),
					// Supported, but difficult to reproduce in short testing time
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "storage_class", "STANDARD"),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "expiration", ""),
					// Currently unsupported in aws_s3_bucket_object resource
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "expires", ""),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "website_redirect_location", ""),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "metadata.%", "0"),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "tags.%", "1"),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "object_lock_legal_hold_status", ""),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "object_lock_mode", ""),
					resource.TestCheckResourceAttr("data.aws_s3_bucket_object.obj", "object_lock_retain_until_date", ""),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanSpacesBucketObject_LeadingSlash(t *testing.T) {
	var rObj s3.GetObjectOutput
	var dsObj1, dsObj2, dsObj3 s3.GetObjectOutput
	resourceName := "aws_s3_bucket_object.object"
	dataSourceName1 := "data.aws_s3_bucket_object.obj1"
	dataSourceName2 := "data.aws_s3_bucket_object.obj2"
	dataSourceName3 := "data.aws_s3_bucket_object.obj3"
	rInt := acctest.RandInt()
	resourceOnlyConf, conf := testAccDataSourceDigitalOceanSpacesObjectConfig_leadingSlash(rInt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: resourceOnlyConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketObjectExists(resourceName, &rObj),
				),
			},
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesObjectDataSourceExists(dataSourceName1, &dsObj1),
					resource.TestCheckResourceAttr(dataSourceName1, "content_length", "3"),
					resource.TestCheckResourceAttr(dataSourceName1, "content_type", "text/plain"),
					resource.TestCheckResourceAttr(dataSourceName1, "etag", "a6105c0a611b41b08f1209506350279e"),
					resource.TestMatchResourceAttr(dataSourceName1, "last_modified",
						regexp.MustCompile("^[a-zA-Z]{3}, [0-9]+ [a-zA-Z]+ [0-9]{4} [0-9:]+ [A-Z]+$")),
					resource.TestCheckResourceAttr(dataSourceName1, "body", "yes"),
					testAccCheckDigitalOceanSpacesObjectDataSourceExists(dataSourceName2, &dsObj2),
					resource.TestCheckResourceAttr(dataSourceName2, "content_length", "3"),
					resource.TestCheckResourceAttr(dataSourceName2, "content_type", "text/plain"),
					resource.TestCheckResourceAttr(dataSourceName2, "etag", "a6105c0a611b41b08f1209506350279e"),
					resource.TestMatchResourceAttr(dataSourceName2, "last_modified",
						regexp.MustCompile("^[a-zA-Z]{3}, [0-9]+ [a-zA-Z]+ [0-9]{4} [0-9:]+ [A-Z]+$")),
					resource.TestCheckResourceAttr(dataSourceName2, "body", "yes"),
					testAccCheckDigitalOceanSpacesObjectDataSourceExists(dataSourceName3, &dsObj3),
					resource.TestCheckResourceAttr(dataSourceName3, "content_length", "3"),
					resource.TestCheckResourceAttr(dataSourceName3, "content_type", "text/plain"),
					resource.TestCheckResourceAttr(dataSourceName3, "etag", "a6105c0a611b41b08f1209506350279e"),
					resource.TestMatchResourceAttr(dataSourceName3, "last_modified",
						regexp.MustCompile("^[a-zA-Z]{3}, [0-9]+ [a-zA-Z]+ [0-9]{4} [0-9:]+ [A-Z]+$")),
					resource.TestCheckResourceAttr(dataSourceName3, "body", "yes"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanSpacesBucketObject_MultipleSlashes(t *testing.T) {
	var rObj1, rObj2 s3.GetObjectOutput
	var dsObj1, dsObj2, dsObj3 s3.GetObjectOutput
	resourceName1 := "aws_s3_bucket_object.object1"
	resourceName2 := "aws_s3_bucket_object.object2"
	dataSourceName1 := "data.aws_s3_bucket_object.obj1"
	dataSourceName2 := "data.aws_s3_bucket_object.obj2"
	dataSourceName3 := "data.aws_s3_bucket_object.obj3"
	rInt := acctest.RandInt()
	resourceOnlyConf, conf := testAccDataSourceDigitalOceanSpacesObjectConfig_multipleSlashes(rInt)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		Providers:                 testAccProviders,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: resourceOnlyConf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAWSS3BucketObjectExists(resourceName1, &rObj1),
					testAccCheckAWSS3BucketObjectExists(resourceName2, &rObj2),
				),
			},
			{
				Config: conf,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesObjectDataSourceExists(dataSourceName1, &dsObj1),
					resource.TestCheckResourceAttr(dataSourceName1, "content_length", "3"),
					resource.TestCheckResourceAttr(dataSourceName1, "content_type", "text/plain"),
					resource.TestCheckResourceAttr(dataSourceName1, "body", "yes"),
					testAccCheckDigitalOceanSpacesObjectDataSourceExists(dataSourceName2, &dsObj2),
					resource.TestCheckResourceAttr(dataSourceName2, "content_length", "3"),
					resource.TestCheckResourceAttr(dataSourceName2, "content_type", "text/plain"),
					resource.TestCheckResourceAttr(dataSourceName2, "body", "yes"),
					testAccCheckDigitalOceanSpacesObjectDataSourceExists(dataSourceName3, &dsObj3),
					resource.TestCheckResourceAttr(dataSourceName3, "content_length", "2"),
					resource.TestCheckResourceAttr(dataSourceName3, "content_type", "text/plain"),
					resource.TestCheckResourceAttr(dataSourceName3, "body", "no"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanSpacesObjectDataSourceExists(n string, obj *s3.GetObjectOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find S3 object data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("S3 object data source ID not set")
		}

		s3conn := testAccProvider.Meta().(*AWSClient).s3conn
		out, err := s3conn.GetObject(
			&s3.GetObjectInput{
				Bucket: aws.String(rs.Primary.Attributes["bucket"]),
				Key:    aws.String(rs.Primary.Attributes["key"]),
			})
		if err != nil {
			return fmt.Errorf("Failed getting S3 Object from %s: %s",
				rs.Primary.Attributes["bucket"]+"/"+rs.Primary.Attributes["key"], err)
		}

		*obj = *out

		return nil
	}
}

func testAccDataSourceDigitalOceanSpacesObjectConfig_basic(randInt int) (string, string) {
	resources := fmt.Sprintf(`
resource "aws_s3_bucket" "object_bucket" {
	bucket = "tf-object-test-bucket-%d"
}
resource "aws_s3_bucket_object" "object" {
	bucket = "${aws_s3_bucket.object_bucket.bucket}"
	key = "tf-testing-obj-%d"
	content = "Hello World"
}
`, randInt, randInt)

	both := fmt.Sprintf(`%s
data "aws_s3_bucket_object" "obj" {
	bucket = "tf-object-test-bucket-%d"
	key = "tf-testing-obj-%d"
}
`, resources, randInt, randInt)

	return resources, both
}

func testAccDataSourceDigitalOceanSpacesObjectConfig_readableBody(randInt int) (string, string) {
	resources := fmt.Sprintf(`
resource "aws_s3_bucket" "object_bucket" {
	bucket = "tf-object-test-bucket-%d"
}
resource "aws_s3_bucket_object" "object" {
	bucket = "${aws_s3_bucket.object_bucket.bucket}"
	key = "tf-testing-obj-%d-readable"
	content = "yes"
	content_type = "text/plain"
}
`, randInt, randInt)

	both := fmt.Sprintf(`%s
data "aws_s3_bucket_object" "obj" {
	bucket = "tf-object-test-bucket-%d"
	key = "tf-testing-obj-%d-readable"
}
`, resources, randInt, randInt)

	return resources, both
}

func testAccDataSourceDigitalOceanSpacesObjectConfig_allParams(randInt int) (string, string) {
	resources := fmt.Sprintf(`
resource "aws_s3_bucket" "object_bucket" {
	bucket = "tf-object-test-bucket-%d"
	versioning {
		enabled = true
	}
}

resource "aws_s3_bucket_object" "object" {
	bucket = "${aws_s3_bucket.object_bucket.bucket}"
	key = "tf-testing-obj-%d-all-params"
	content = <<CONTENT
{"msg": "Hi there!"}
CONTENT
	content_type = "application/unknown"
	cache_control = "no-cache"
	content_disposition = "attachment"
	content_encoding = "identity"
	content_language = "en-GB"
	tags = {
		Key1 = "Value 1"
	}
}
`, randInt, randInt)

	both := fmt.Sprintf(`%s
data "aws_s3_bucket_object" "obj" {
	bucket = "tf-object-test-bucket-%d"
	key = "tf-testing-obj-%d-all-params"
}
`, resources, randInt, randInt)

	return resources, both
}

func testAccDataSourceDigitalOceanSpacesObjectConfig_leadingSlash(randInt int) (string, string) {
	resources := fmt.Sprintf(`
resource "aws_s3_bucket" "object_bucket" {
  bucket = "tf-object-test-bucket-%d"
}
resource "aws_s3_bucket_object" "object" {
  bucket = "${aws_s3_bucket.object_bucket.bucket}"
  key = "//tf-testing-obj-%d-readable"
  content = "yes"
  content_type = "text/plain"
}
`, randInt, randInt)

	both := fmt.Sprintf(`%s
data "aws_s3_bucket_object" "obj1" {
  bucket = "tf-object-test-bucket-%d"
  key = "tf-testing-obj-%d-readable"
}
data "aws_s3_bucket_object" "obj2" {
  bucket = "tf-object-test-bucket-%d"
  key = "/tf-testing-obj-%d-readable"
}
data "aws_s3_bucket_object" "obj3" {
  bucket = "tf-object-test-bucket-%d"
  key = "//tf-testing-obj-%d-readable"
}
`, resources, randInt, randInt, randInt, randInt, randInt, randInt)

	return resources, both
}

func testAccDataSourceDigitalOceanSpacesObjectConfig_multipleSlashes(randInt int) (string, string) {
	resources := fmt.Sprintf(`
resource "aws_s3_bucket" "object_bucket" {
  bucket = "tf-object-test-bucket-%d"
}
resource "aws_s3_bucket_object" "object1" {
  bucket = "${aws_s3_bucket.object_bucket.bucket}"
  key = "first//second///third//"
  content = "yes"
  content_type = "text/plain"
}
# Without a trailing slash.
resource "aws_s3_bucket_object" "object2" {
  bucket = "${aws_s3_bucket.object_bucket.bucket}"
  key = "/first////second/third"
  content = "no"
  content_type = "text/plain"
}
`, randInt)

	both := fmt.Sprintf(`%s
data "aws_s3_bucket_object" "obj1" {
  bucket = "tf-object-test-bucket-%d"
  key = "first/second/third/"
}
data "aws_s3_bucket_object" "obj2" {
  bucket = "tf-object-test-bucket-%d"
  key = "first//second///third//"
}
data "aws_s3_bucket_object" "obj3" {
  bucket = "tf-object-test-bucket-%d"
  key = "first/second/third"
}
`, resources, randInt, randInt, randInt)

	return resources, both
}
