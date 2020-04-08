package digitalocean

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"regexp"
	"sort"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

const (
	testAccDigitalOceanSpacesBucketObject_TestRegion = "nyc3"
)

func TestAccDigitalOceanSpacesBucketObject_noNameNoKey(t *testing.T) {
	bucketError := regexp.MustCompile(`bucket must not be empty`)
	keyError := regexp.MustCompile(`key must not be empty`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanSpacesBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig:   func() {},
				Config:      testAccDigitalOceanSpacesBucketObjectConfigBasic("", "a key"),
				ExpectError: bucketError,
			},
			{
				PreConfig:   func() {},
				Config:      testAccDigitalOceanSpacesBucketObjectConfigBasic("a name", ""),
				ExpectError: keyError,
			},
		},
	})
}
func TestAccDigitalOceanSpacesBucketObject_empty(t *testing.T) {
	var obj s3.GetObjectOutput
	resourceName := "digitalocean_spaces_bucket_object.object"
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanSpacesBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {},
				Config:    testAccDigitalOceanSpacesBucketObjectConfigEmpty(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketObjectExists(resourceName, &obj),
					testAccCheckDigitalOceanSpacesBucketObjectBody(&obj, ""),
				),
			},
		},
	})
}

func TestAccDigitalOceanSpacesBucketObject_source(t *testing.T) {
	var obj s3.GetObjectOutput
	resourceName := "digitalocean_spaces_bucket_object.object"
	rInt := acctest.RandInt()

	source := testAccDigitalOceanSpacesBucketObjectCreateTempFile(t, "{anything will do }")
	defer os.Remove(source)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanSpacesBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanSpacesBucketObjectConfigSource(rInt, source),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketObjectExists(resourceName, &obj),
					testAccCheckDigitalOceanSpacesBucketObjectBody(&obj, "{anything will do }"),
				),
			},
		},
	})
}

func TestAccDigitalOceanSpacesBucketObject_content(t *testing.T) {
	var obj s3.GetObjectOutput
	resourceName := "digitalocean_spaces_bucket_object.object"
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanSpacesBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {},
				Config:    testAccDigitalOceanSpacesBucketObjectConfigContent(rInt, "some_bucket_content"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketObjectExists(resourceName, &obj),
					testAccCheckDigitalOceanSpacesBucketObjectBody(&obj, "some_bucket_content"),
				),
			},
		},
	})
}

func TestAccDigitalOceanSpacesBucketObject_contentBase64(t *testing.T) {
	var obj s3.GetObjectOutput
	resourceName := "digitalocean_spaces_bucket_object.object"
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanSpacesBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {},
				Config:    testAccDigitalOceanSpacesBucketObjectConfigContentBase64(rInt, base64.StdEncoding.EncodeToString([]byte("some_bucket_content"))),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketObjectExists(resourceName, &obj),
					testAccCheckDigitalOceanSpacesBucketObjectBody(&obj, "some_bucket_content"),
				),
			},
		},
	})
}

func TestAccDigitalOceanSpacesBucketObject_withContentCharacteristics(t *testing.T) {
	var obj s3.GetObjectOutput
	resourceName := "digitalocean_spaces_bucket_object.object"
	rInt := acctest.RandInt()

	source := testAccDigitalOceanSpacesBucketObjectCreateTempFile(t, "{anything will do }")
	defer os.Remove(source)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanSpacesBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanSpacesBucketObjectConfig_withContentCharacteristics(rInt, source),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketObjectExists(resourceName, &obj),
					testAccCheckDigitalOceanSpacesBucketObjectBody(&obj, "{anything will do }"),
					resource.TestCheckResourceAttr(resourceName, "content_type", "binary/octet-stream"),
					resource.TestCheckResourceAttr(resourceName, "website_redirect", "http://google.com"),
				),
			},
		},
	})
}

func TestAccDigitalOceanSpacesBucketObject_NonVersioned(t *testing.T) {
	sourceInitial := testAccDigitalOceanSpacesBucketObjectCreateTempFile(t, "initial object state")
	defer os.Remove(sourceInitial)

	var originalObj s3.GetObjectOutput
	resourceName := "digitalocean_spaces_bucket_object.object"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanSpacesBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanSpacesBucketObjectConfig_NonVersioned(acctest.RandInt(), sourceInitial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketObjectExists(resourceName, &originalObj),
					testAccCheckDigitalOceanSpacesBucketObjectBody(&originalObj, "initial object state"),
					resource.TestCheckResourceAttr(resourceName, "version_id", ""),
				),
			},
		},
	})
}

func TestAccDigitalOceanSpacesBucketObject_updates(t *testing.T) {
	var originalObj, modifiedObj s3.GetObjectOutput
	resourceName := "digitalocean_spaces_bucket_object.object"
	rInt := acctest.RandInt()

	sourceInitial := testAccDigitalOceanSpacesBucketObjectCreateTempFile(t, "initial object state")
	defer os.Remove(sourceInitial)
	sourceModified := testAccDigitalOceanSpacesBucketObjectCreateTempFile(t, "modified object")
	defer os.Remove(sourceInitial)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanSpacesBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanSpacesBucketObjectConfig_updateable(rInt, false, sourceInitial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketObjectExists(resourceName, &originalObj),
					testAccCheckDigitalOceanSpacesBucketObjectBody(&originalObj, "initial object state"),
					resource.TestCheckResourceAttr(resourceName, "etag", "647d1d58e1011c743ec67d5e8af87b53"),
				),
			},
			{
				Config: testAccDigitalOceanSpacesBucketObjectConfig_updateable(rInt, false, sourceModified),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketObjectExists(resourceName, &modifiedObj),
					testAccCheckDigitalOceanSpacesBucketObjectBody(&modifiedObj, "modified object"),
					resource.TestCheckResourceAttr(resourceName, "etag", "1c7fd13df1515c2a13ad9eb068931f09"),
				),
			},
		},
	})
}

func TestAccDigitalOceanSpacesBucketObject_updateSameFile(t *testing.T) {
	var originalObj, modifiedObj s3.GetObjectOutput
	resourceName := "digitalocean_spaces_bucket_object.object"
	rInt := acctest.RandInt()

	startingData := "lane 8"
	changingData := "chicane"

	filename := testAccDigitalOceanSpacesBucketObjectCreateTempFile(t, startingData)
	defer os.Remove(filename)

	rewriteFile := func(*terraform.State) error {
		if err := ioutil.WriteFile(filename, []byte(changingData), 0644); err != nil {
			os.Remove(filename)
			t.Fatal(err)
		}
		return nil
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanSpacesBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanSpacesBucketObjectConfig_updateable(rInt, false, filename),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketObjectExists(resourceName, &originalObj),
					testAccCheckDigitalOceanSpacesBucketObjectBody(&originalObj, startingData),
					resource.TestCheckResourceAttr(resourceName, "etag", "aa48b42f36a2652cbee40c30a5df7d25"),
					rewriteFile,
				),
				ExpectNonEmptyPlan: true,
			},
			{
				Config: testAccDigitalOceanSpacesBucketObjectConfig_updateable(rInt, false, filename),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketObjectExists(resourceName, &modifiedObj),
					testAccCheckDigitalOceanSpacesBucketObjectBody(&modifiedObj, changingData),
					resource.TestCheckResourceAttr(resourceName, "etag", "fafc05f8c4da0266a99154681ab86e8c"),
				),
			},
		},
	})
}

func TestAccDigitalOceanSpacesBucketObject_updatesWithVersioning(t *testing.T) {
	var originalObj, modifiedObj s3.GetObjectOutput
	resourceName := "digitalocean_spaces_bucket_object.object"
	rInt := acctest.RandInt()

	sourceInitial := testAccDigitalOceanSpacesBucketObjectCreateTempFile(t, "initial versioned object state")
	defer os.Remove(sourceInitial)
	sourceModified := testAccDigitalOceanSpacesBucketObjectCreateTempFile(t, "modified versioned object")
	defer os.Remove(sourceInitial)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanSpacesBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanSpacesBucketObjectConfig_updateable(rInt, true, sourceInitial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketObjectExists(resourceName, &originalObj),
					testAccCheckDigitalOceanSpacesBucketObjectBody(&originalObj, "initial versioned object state"),
					resource.TestCheckResourceAttr(resourceName, "etag", "cee4407fa91906284e2a5e5e03e86b1b"),
				),
			},
			{
				Config: testAccDigitalOceanSpacesBucketObjectConfig_updateable(rInt, true, sourceModified),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketObjectExists(resourceName, &modifiedObj),
					testAccCheckDigitalOceanSpacesBucketObjectBody(&modifiedObj, "modified versioned object"),
					resource.TestCheckResourceAttr(resourceName, "etag", "00b8c73b1b50e7cc932362c7225b8e29"),
					testAccCheckDigitalOceanSpacesBucketObjectVersionIdDiffers(&modifiedObj, &originalObj),
				),
			},
		},
	})
}

func TestAccDigitalOceanSpacesBucketObject_acl(t *testing.T) {
	var obj1, obj2, obj3 s3.GetObjectOutput
	resourceName := "digitalocean_spaces_bucket_object.object"
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanSpacesBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanSpacesBucketObjectConfig_acl(rInt, "some_bucket_content", "private"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketObjectExists(resourceName, &obj1),
					testAccCheckDigitalOceanSpacesBucketObjectBody(&obj1, "some_bucket_content"),
					resource.TestCheckResourceAttr(resourceName, "acl", "private"),
					testAccCheckDigitalOceanSpacesBucketObjectAcl(resourceName, []string{"FULL_CONTROL"}),
				),
			},
			{
				Config: testAccDigitalOceanSpacesBucketObjectConfig_acl(rInt, "some_bucket_content", "public-read"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketObjectExists(resourceName, &obj2),
					testAccCheckDigitalOceanSpacesBucketObjectVersionIdEquals(&obj2, &obj1),
					testAccCheckDigitalOceanSpacesBucketObjectBody(&obj2, "some_bucket_content"),
					resource.TestCheckResourceAttr(resourceName, "acl", "public-read"),
					testAccCheckDigitalOceanSpacesBucketObjectAcl(resourceName, []string{"FULL_CONTROL", "READ"}),
				),
			},
			{
				Config: testAccDigitalOceanSpacesBucketObjectConfig_acl(rInt, "changed_some_bucket_content", "private"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketObjectExists(resourceName, &obj3),
					testAccCheckDigitalOceanSpacesBucketObjectVersionIdDiffers(&obj3, &obj2),
					testAccCheckDigitalOceanSpacesBucketObjectBody(&obj3, "changed_some_bucket_content"),
					resource.TestCheckResourceAttr(resourceName, "acl", "private"),
					testAccCheckDigitalOceanSpacesBucketObjectAcl(resourceName, []string{"FULL_CONTROL"}),
				),
			},
		},
	})
}

func TestAccDigitalOceanSpacesBucketObject_metadata(t *testing.T) {
	rInt := acctest.RandInt()
	var obj s3.GetObjectOutput
	resourceName := "digitalocean_spaces_bucket_object.object"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanSpacesBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanSpacesBucketObjectConfig_withMetadata(rInt, "key1", "value1", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketObjectExists(resourceName, &obj),
					resource.TestCheckResourceAttr(resourceName, "metadata.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key2", "value2"),
				),
			},
			{
				Config: testAccDigitalOceanSpacesBucketObjectConfig_withMetadata(rInt, "key1", "value1updated", "key3", "value3"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketObjectExists(resourceName, &obj),
					resource.TestCheckResourceAttr(resourceName, "metadata.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "metadata.key3", "value3"),
				),
			},
			{
				Config: testAccDigitalOceanSpacesBucketObjectConfigEmpty(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketObjectExists(resourceName, &obj),
					resource.TestCheckResourceAttr(resourceName, "metadata.%", "0"),
				),
			},
		},
	})
}

func TestAccDigitalOceanSpacesBucketObject_storageClass(t *testing.T) {
	var obj s3.GetObjectOutput
	resourceName := "digitalocean_spaces_bucket_object.object"
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanSpacesBucketObjectDestroy,
		Steps: []resource.TestStep{
			{
				PreConfig: func() {},
				Config:    testAccDigitalOceanSpacesBucketObjectConfigContent(rInt, "some_bucket_content"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesBucketObjectExists(resourceName, &obj),
					resource.TestCheckResourceAttr(resourceName, "storage_class", "STANDARD"),
					testAccCheckDigitalOceanSpacesBucketObjectStorageClass(resourceName, "STANDARD"),
				),
			},
		},
	})
}

func testAccGetS3Conn() (*s3.S3, error) {
	client, err := testAccProvider.Meta().(*CombinedConfig).spacesClient(testAccDigitalOceanSpacesBucketObject_TestRegion)
	if err != nil {
		return nil, err
	}

	s3conn := s3.New(client)

	return s3conn, nil
}

func testAccCheckDigitalOceanSpacesBucketObjectVersionIdDiffers(first, second *s3.GetObjectOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if first.VersionId == nil {
			return fmt.Errorf("Expected first object to have VersionId: %s", first)
		}
		if second.VersionId == nil {
			return fmt.Errorf("Expected second object to have VersionId: %s", second)
		}

		if *first.VersionId == *second.VersionId {
			return fmt.Errorf("Expected Version IDs to differ, but they are equal (%s)", *first.VersionId)
		}

		return nil
	}
}

func testAccCheckDigitalOceanSpacesBucketObjectVersionIdEquals(first, second *s3.GetObjectOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if first.VersionId == nil {
			return fmt.Errorf("Expected first object to have VersionId: %s", first)
		}
		if second.VersionId == nil {
			return fmt.Errorf("Expected second object to have VersionId: %s", second)
		}

		if *first.VersionId != *second.VersionId {
			return fmt.Errorf("Expected Version IDs to be equal, but they differ (%s, %s)", *first.VersionId, *second.VersionId)
		}

		return nil
	}
}

func testAccCheckDigitalOceanSpacesBucketObjectDestroy(s *terraform.State) error {
	s3conn, err := testAccGetS3Conn()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		switch rs.Type {
		case "digitalocean_spaces_bucket_object":
			_, err := s3conn.HeadObject(
				&s3.HeadObjectInput{
					Bucket:  aws.String(rs.Primary.Attributes["bucket"]),
					Key:     aws.String(rs.Primary.Attributes["key"]),
					IfMatch: aws.String(rs.Primary.Attributes["etag"]),
				})
			if err == nil {
				return fmt.Errorf("Spaces Bucket Object still exists: %s", rs.Primary.ID)
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

func testAccCheckDigitalOceanSpacesBucketObjectExists(n string, obj *s3.GetObjectOutput) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No S3 Bucket Object ID is set")
		}

		s3conn, err := testAccGetS3Conn()
		if err != nil {
			return err
		}

		out, err := s3conn.GetObject(
			&s3.GetObjectInput{
				Bucket:  aws.String(rs.Primary.Attributes["bucket"]),
				Key:     aws.String(rs.Primary.Attributes["key"]),
				IfMatch: aws.String(rs.Primary.Attributes["etag"]),
			})
		if err != nil {
			return fmt.Errorf("S3Bucket Object error: %s", err)
		}

		*obj = *out

		return nil
	}
}

func testAccCheckDigitalOceanSpacesBucketObjectBody(obj *s3.GetObjectOutput, want string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		body, err := ioutil.ReadAll(obj.Body)
		if err != nil {
			return fmt.Errorf("failed to read body: %s", err)
		}
		obj.Body.Close()

		if got := string(body); got != want {
			return fmt.Errorf("wrong result body %q; want %q", got, want)
		}

		return nil
	}
}

func testAccCheckDigitalOceanSpacesBucketObjectAcl(n string, expectedPerms []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]

		s3conn, err := testAccGetS3Conn()
		if err != nil {
			return err
		}

		out, err := s3conn.GetObjectAcl(&s3.GetObjectAclInput{
			Bucket: aws.String(rs.Primary.Attributes["bucket"]),
			Key:    aws.String(rs.Primary.Attributes["key"]),
		})

		if err != nil {
			return fmt.Errorf("GetObjectAcl error: %v", err)
		}

		var perms []string
		for _, v := range out.Grants {
			perms = append(perms, *v.Permission)
		}
		sort.Strings(perms)

		if !reflect.DeepEqual(perms, expectedPerms) {
			return fmt.Errorf("Expected ACL permissions to be %v, got %v", expectedPerms, perms)
		}

		return nil
	}
}

func testAccCheckDigitalOceanSpacesBucketObjectStorageClass(n, expectedClass string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs := s.RootModule().Resources[n]

		s3conn, err := testAccGetS3Conn()
		if err != nil {
			return err
		}

		out, err := s3conn.HeadObject(&s3.HeadObjectInput{
			Bucket: aws.String(rs.Primary.Attributes["bucket"]),
			Key:    aws.String(rs.Primary.Attributes["key"]),
		})

		if err != nil {
			return fmt.Errorf("HeadObject error: %v", err)
		}

		// The "STANDARD" (which is also the default) storage
		// class when set would not be included in the results.
		storageClass := s3.StorageClassStandard
		if out.StorageClass != nil {
			storageClass = *out.StorageClass
		}

		if storageClass != expectedClass {
			return fmt.Errorf("Expected Storage Class to be %v, got %v",
				expectedClass, storageClass)
		}

		return nil
	}
}

func testAccDigitalOceanSpacesBucketObjectCreateTempFile(t *testing.T, data string) string {
	tmpFile, err := ioutil.TempFile("", "tf-acc-s3-obj")
	if err != nil {
		t.Fatal(err)
	}
	filename := tmpFile.Name()

	err = ioutil.WriteFile(filename, []byte(data), 0644)
	if err != nil {
		os.Remove(filename)
		t.Fatal(err)
	}

	return filename
}

func testAccDigitalOceanSpacesBucketObjectConfigBasic(bucket, key string) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket_object" "object" {
  region = "%s"
  bucket = "%s"
  key = "%s"
}
`, testAccDigitalOceanSpacesBucketObject_TestRegion, bucket, key)
}

func testAccDigitalOceanSpacesBucketObjectConfigEmpty(randInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "object_bucket" {
  region = "%s"
  name   = "tf-object-test-bucket-%d"
  force_destroy = true
}

resource "digitalocean_spaces_bucket_object" "object" {
  region = digitalocean_spaces_bucket.object_bucket.region
  bucket = digitalocean_spaces_bucket.object_bucket.name
  key = "test-key"
}
`, testAccDigitalOceanSpacesBucketObject_TestRegion, randInt)
}

func testAccDigitalOceanSpacesBucketObjectConfigSource(randInt int, source string) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "object_bucket" {
  region = "%s"
  name   = "tf-object-test-bucket-%d"
  force_destroy = true
}

resource "digitalocean_spaces_bucket_object" "object" {
  region       = digitalocean_spaces_bucket.object_bucket.region
  bucket       = digitalocean_spaces_bucket.object_bucket.name
  key          = "test-key"
  source       = "%s"
  content_type = "binary/octet-stream"
}
`, testAccDigitalOceanSpacesBucketObject_TestRegion, randInt, source)
}

func testAccDigitalOceanSpacesBucketObjectConfig_withContentCharacteristics(randInt int, source string) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "object_bucket" {
  region = "%s"
  name   = "tf-object-test-bucket-%d"
  force_destroy = true
}

resource "digitalocean_spaces_bucket_object" "object" {
  region           = digitalocean_spaces_bucket.object_bucket.region
  bucket           = digitalocean_spaces_bucket.object_bucket.name
  key              = "test-key"
  source           = "%s"
  content_language = "en"
  content_type     = "binary/octet-stream"
  website_redirect = "http://google.com"
}
`, testAccDigitalOceanSpacesBucketObject_TestRegion, randInt, source)
}

func testAccDigitalOceanSpacesBucketObjectConfigContent(randInt int, content string) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "object_bucket" {
  region = "%s"
  name   = "tf-object-test-bucket-%d"
  force_destroy = true
}

resource "digitalocean_spaces_bucket_object" "object" {
  region  = digitalocean_spaces_bucket.object_bucket.region
  bucket  = digitalocean_spaces_bucket.object_bucket.name
  key     = "test-key"
  content = "%s"
}
`, testAccDigitalOceanSpacesBucketObject_TestRegion, randInt, content)
}

func testAccDigitalOceanSpacesBucketObjectConfigContentBase64(randInt int, contentBase64 string) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "object_bucket" {
  region = "%s"
  name   = "tf-object-test-bucket-%d"
  force_destroy = true
}

resource "digitalocean_spaces_bucket_object" "object" {
  region         = digitalocean_spaces_bucket.object_bucket.region
  bucket         = digitalocean_spaces_bucket.object_bucket.name
  key            = "test-key"
  content_base64 = "%s"
}
`, testAccDigitalOceanSpacesBucketObject_TestRegion, randInt, contentBase64)
}

func testAccDigitalOceanSpacesBucketObjectConfig_updateable(randInt int, bucketVersioning bool, source string) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "object_bucket_3" {
  region = "%s"
  name   = "tf-object-test-bucket-%d"
  force_destroy = true

  versioning {
    enabled = %t
  }
}

resource "digitalocean_spaces_bucket_object" "object" {
  region = digitalocean_spaces_bucket.object_bucket_3.region
  bucket = digitalocean_spaces_bucket.object_bucket_3.name
  key    = "updateable-key"
  source = "%s"
  etag   = "${filemd5("%s")}"
}
`, testAccDigitalOceanSpacesBucketObject_TestRegion, randInt, bucketVersioning, source, source)
}

func testAccDigitalOceanSpacesBucketObjectConfig_acl(randInt int, content, acl string) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "object_bucket" {
  region = "%s"
  name   = "tf-object-test-bucket-%d"
  force_destroy = true

  versioning {
    enabled = true
  }
}

resource "digitalocean_spaces_bucket_object" "object" {
  region  = digitalocean_spaces_bucket.object_bucket.region
  bucket  = digitalocean_spaces_bucket.object_bucket.name
  key     = "test-key"
  content = "%s"
  acl     = "%s"
}
`, testAccDigitalOceanSpacesBucketObject_TestRegion, randInt, content, acl)
}

func testAccDigitalOceanSpacesBucketObjectConfig_storageClass(randInt int, storage_class string) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "object_bucket" {
  region = "%s"
  name   = "tf-object-test-bucket-%d"
  force_destroy = true
}

resource "digitalocean_spaces_bucket_object" "object" {
  region        = digitalocean_spaces_bucket.object_bucket.region
  bucket        = digitalocean_spaces_bucket.object_bucket.name
  key           = "test-key"
  content       = "some_bucket_content"
  storage_class = "%s"
}
`, testAccDigitalOceanSpacesBucketObject_TestRegion, randInt, storage_class)
}

func testAccDigitalOceanSpacesBucketObjectConfig_withMetadata(randInt int, metadataKey1, metadataValue1, metadataKey2, metadataValue2 string) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "object_bucket" {
  region = "%s"
  name   = "tf-object-test-bucket-%d"
  force_destroy = true
}

resource "digitalocean_spaces_bucket_object" "object" {
  region = digitalocean_spaces_bucket.object_bucket.region
  bucket  = digitalocean_spaces_bucket.object_bucket.name
  key     = "test-key"

  metadata = {
    %[3]s = %[4]q
    %[5]s = %[6]q
  }
}
`, testAccDigitalOceanSpacesBucketObject_TestRegion, randInt, metadataKey1, metadataValue1, metadataKey2, metadataValue2)
}

func testAccDigitalOceanSpacesBucketObjectConfig_NonVersioned(randInt int, source string) string {
	return fmt.Sprintf(`
resource "digitalocean_spaces_bucket" "object_bucket_3" {
  region = "%s"
  name   = "tf-object-test-bucket-%d"
  force_destroy = true
}

resource "digitalocean_spaces_bucket_object" "object" {
  region = digitalocean_spaces_bucket.object_bucket_3.region
  bucket = digitalocean_spaces_bucket.object_bucket_3.name
  key    = "updateable-key"
  source = "%s"
  etag   = "${filemd5("%s")}"
}
`, testAccDigitalOceanSpacesBucketObject_TestRegion, randInt, source, source)
}

func testSweepSkipSweepError(err error) bool {
	// Ignore missing API endpoints
	if isAWSErr(err, "RequestError", "send request failed") {
		return true
	}
	// Ignore unsupported API calls
	if isAWSErr(err, "UnsupportedOperation", "") {
		return true
	}
	// Ignore more unsupported API calls
	// InvalidParameterValue: Use of cache security groups is not permitted in this API version for your account.
	if isAWSErr(err, "InvalidParameterValue", "not permitted in this API version for your account") {
		return true
	}
	// InvalidParameterValue: Access Denied to API Version: APIGlobalDatabases
	if isAWSErr(err, "InvalidParameterValue", "Access Denied to API Version") {
		return true
	}
	// GovCloud has endpoints that respond with (no message provided):
	// AccessDeniedException:
	// Since acceptance test sweepers are best effort and this response is very common,
	// we allow bypassing this error globally instead of individual test sweeper fixes.
	if isAWSErr(err, "AccessDeniedException", "") {
		return true
	}
	// Example: BadRequestException: vpc link not supported for region us-gov-west-1
	if isAWSErr(err, "BadRequestException", "not supported") {
		return true
	}
	// Example: InvalidAction: The action DescribeTransitGatewayAttachments is not valid for this web service
	if isAWSErr(err, "InvalidAction", "is not valid") {
		return true
	}
	return false
}

func testS3BucketRegion(conn *s3.S3, bucket string) (string, error) {
	input := &s3.GetBucketLocationInput{
		Bucket: aws.String(bucket),
	}

	output, err := conn.GetBucketLocation(input)

	if err != nil {
		return "", err
	}

	if output == nil || output.LocationConstraint == nil {
		return "nyc3", nil
	}

	return aws.StringValue(output.LocationConstraint), nil
}
