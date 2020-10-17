package digitalocean

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mitchellh/go-homedir"
)

func resourceDigitalOceanSpacesBucketObject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanSpacesBucketObjectCreate,
		ReadContext:   resourceDigitalOceanSpacesBucketObjectRead,
		UpdateContext: resourceDigitalOceanSpacesBucketObjectUpdate,
		DeleteContext: resourceDigitalOceanSpacesBucketObjectDelete,

		CustomizeDiff: resourceDigitalOceanSpacesBucketObjectCustomizeDiff,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"bucket": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"key": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"acl": {
				Type:     schema.TypeString,
				Default:  s3.ObjectCannedACLPrivate,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					s3.ObjectCannedACLPrivate,
					s3.ObjectCannedACLPublicRead,
				}, false),
			},

			"cache_control": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content_disposition": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content_encoding": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"content_language": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"metadata": {
				Type:         schema.TypeMap,
				ValidateFunc: validateMetadataIsLowerCase,
				Optional:     true,
				Elem:         &schema.Schema{Type: schema.TypeString},
			},

			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"source": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"content", "content_base64"},
			},

			"content": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source", "content_base64"},
			},

			"content_base64": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source", "content"},
			},

			"etag": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"website_redirect": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"force_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func s3connFromResourceData(d *schema.ResourceData, meta interface{}) (*s3.S3, error) {
	region := d.Get("region").(string)

	client, err := meta.(*CombinedConfig).spacesClient(region)
	if err != nil {
		return nil, err
	}

	svc := s3.New(client)
	return svc, nil
}

func resourceDigitalOceanSpacesBucketObjectPut(d *schema.ResourceData, meta interface{}) error {
	s3conn, err := s3connFromResourceData(d, meta)
	if err != nil {
		return err
	}

	var body io.ReadSeeker

	if v, ok := d.GetOk("source"); ok {
		source := v.(string)
		path, err := homedir.Expand(source)
		if err != nil {
			return fmt.Errorf("Error expanding homedir in source (%s): %s", source, err)
		}
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("Error opening Spaces bucket object source (%s): %s", path, err)
		}

		body = file
		defer func() {
			err := file.Close()
			if err != nil {
				log.Printf("[WARN] Error closing Spaces bucket object source (%s): %s", path, err)
			}
		}()
	} else if v, ok := d.GetOk("content"); ok {
		content := v.(string)
		body = bytes.NewReader([]byte(content))
	} else if v, ok := d.GetOk("content_base64"); ok {
		content := v.(string)
		// We can't do streaming decoding here (with base64.NewDecoder) because
		// the AWS SDK requires an io.ReadSeeker but a base64 decoder can't seek.
		contentRaw, err := base64.StdEncoding.DecodeString(content)
		if err != nil {
			return fmt.Errorf("error decoding content_base64: %s", err)
		}
		body = bytes.NewReader(contentRaw)
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	putInput := &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		ACL:    aws.String(d.Get("acl").(string)),
		Body:   body,
	}

	if v, ok := d.GetOk("cache_control"); ok {
		putInput.CacheControl = aws.String(v.(string))
	}

	if v, ok := d.GetOk("content_type"); ok {
		putInput.ContentType = aws.String(v.(string))
	}

	if v, ok := d.GetOk("metadata"); ok {
		putInput.Metadata = stringMapToPointers(v.(map[string]interface{}))
	}

	if v, ok := d.GetOk("content_encoding"); ok {
		putInput.ContentEncoding = aws.String(v.(string))
	}

	if v, ok := d.GetOk("content_language"); ok {
		putInput.ContentLanguage = aws.String(v.(string))
	}

	if v, ok := d.GetOk("content_disposition"); ok {
		putInput.ContentDisposition = aws.String(v.(string))
	}

	if v, ok := d.GetOk("website_redirect"); ok {
		putInput.WebsiteRedirectLocation = aws.String(v.(string))
	}

	if _, err := s3conn.PutObject(putInput); err != nil {
		return fmt.Errorf("Error putting object in Spaces bucket (%s): %s", bucket, err)
	}

	d.SetId(key)
	return resourceDigitalOceanSpacesBucketObjectRead(d, meta)
}

func resourceDigitalOceanSpacesBucketObjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	return resourceDigitalOceanSpacesBucketObjectPut(d, meta)
}

func resourceDigitalOceanSpacesBucketObjectRead(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	s3conn, err := s3connFromResourceData(d, meta)
	if err != nil {
		return err
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	resp, err := s3conn.HeadObject(
		&s3.HeadObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		})

	if err != nil {
		// If S3 returns a 404 Request Failure, mark the object as destroyed
		if awsErr, ok := err.(awserr.RequestFailure); ok && awsErr.StatusCode() == 404 {
			d.SetId("")
			log.Printf("[WARN] Error Reading Object (%s), object not found (HTTP status 404)", key)
			return nil
		}
		return err
	}
	log.Printf("[DEBUG] Reading Spaces Bucket Object meta: %s", resp)

	d.Set("cache_control", resp.CacheControl)
	d.Set("content_disposition", resp.ContentDisposition)
	d.Set("content_encoding", resp.ContentEncoding)
	d.Set("content_language", resp.ContentLanguage)
	d.Set("content_type", resp.ContentType)
	metadata := pointersMapToStringList(resp.Metadata)

	// AWS Go SDK capitalizes metadata, this is a workaround. https://github.com/aws/aws-sdk-go/issues/445
	for k, v := range metadata {
		delete(metadata, k)
		metadata[strings.ToLower(k)] = v
	}

	if err := d.Set("metadata", metadata); err != nil {
		return fmt.Errorf("error setting metadata: %s", err)
	}
	d.Set("version_id", resp.VersionId)
	d.Set("website_redirect", resp.WebsiteRedirectLocation)

	// See https://forums.aws.amazon.com/thread.jspa?threadID=44003
	d.Set("etag", strings.Trim(aws.StringValue(resp.ETag), `"`))

	return nil
}

func resourceDigitalOceanSpacesBucketObjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	// Changes to any of these attributes requires creation of a new object version (if bucket is versioned):
	for _, key := range []string{
		"cache_control",
		"content_base64",
		"content_disposition",
		"content_encoding",
		"content_language",
		"content_type",
		"content",
		"etag",
		"metadata",
		"source",
		"website_redirect",
	} {
		if d.HasChange(key) {
			return resourceDigitalOceanSpacesBucketObjectPut(d, meta)
		}
	}

	conn, err := s3connFromResourceData(d, meta)
	if err != nil {
		return err
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)

	if d.HasChange("acl") {
		_, err := conn.PutObjectAcl(&s3.PutObjectAclInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
			ACL:    aws.String(d.Get("acl").(string)),
		})
		if err != nil {
			return fmt.Errorf("error putting Spaces object ACL: %s", err)
		}
	}

	return resourceDigitalOceanSpacesBucketObjectRead(d, meta)
}

func resourceDigitalOceanSpacesBucketObjectDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	s3conn, err := s3connFromResourceData(d, meta)
	if err != nil {
		return err
	}

	bucket := d.Get("bucket").(string)
	key := d.Get("key").(string)
	// We are effectively ignoring any leading '/' in the key name as aws.Config.DisableRestProtocolURICleaning is false
	key = strings.TrimPrefix(key, "/")

	if _, ok := d.GetOk("version_id"); ok {
		err = deleteAllS3ObjectVersions(s3conn, bucket, key, d.Get("force_destroy").(bool), false)
	} else {
		err = deleteS3ObjectVersion(s3conn, bucket, key, "", false)
	}

	if err != nil {
		return fmt.Errorf("error deleting Spaces Bucket (%s) Object (%s): %s", bucket, key, err)
	}

	return nil
}

func validateMetadataIsLowerCase(v interface{}, k string) (ws []string, errors []error) {
	value := v.(map[string]interface{})

	for k := range value {
		if k != strings.ToLower(k) {
			errors = append(errors, fmt.Errorf(
				"Metadata must be lowercase only. Offending key: %q", k))
		}
	}
	return
}

func resourceDigitalOceanSpacesBucketObjectCustomizeDiff(
	ctx context.Context,
	d *schema.ResourceDiff,
	meta interface{},
) error {
	if d.HasChange("etag") {
		d.SetNewComputed("version_id")
	}

	return nil
}

// deleteAllS3ObjectVersions deletes all versions of a specified key from an S3 bucket.
// If key is empty then all versions of all objects are deleted.
// Set force to true to override any S3 object lock protections on object lock enabled buckets.
func deleteAllS3ObjectVersions(conn *s3.S3, bucketName, key string, force, ignoreObjectErrors bool) error {
	input := &s3.ListObjectVersionsInput{
		Bucket: aws.String(bucketName),
	}
	if key != "" {
		input.Prefix = aws.String(key)
	}

	var lastErr error
	err := conn.ListObjectVersionsPages(input, func(page *s3.ListObjectVersionsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, objectVersion := range page.Versions {
			objectKey := aws.StringValue(objectVersion.Key)
			objectVersionID := aws.StringValue(objectVersion.VersionId)

			if key != "" && key != objectKey {
				continue
			}

			err := deleteS3ObjectVersion(conn, bucketName, objectKey, objectVersionID, force)
			if isAWSErr(err, "AccessDenied", "") && force {
				// Remove any legal hold.
				resp, err := conn.HeadObject(&s3.HeadObjectInput{
					Bucket:    aws.String(bucketName),
					Key:       objectVersion.Key,
					VersionId: objectVersion.VersionId,
				})

				if err != nil {
					log.Printf("[ERROR] Error getting Spaces Bucket (%s) Object (%s) Version (%s) metadata: %s", bucketName, objectKey, objectVersionID, err)
					lastErr = err
					continue
				}

				if aws.StringValue(resp.ObjectLockLegalHoldStatus) == s3.ObjectLockLegalHoldStatusOn {
					_, err := conn.PutObjectLegalHold(&s3.PutObjectLegalHoldInput{
						Bucket:    aws.String(bucketName),
						Key:       objectVersion.Key,
						VersionId: objectVersion.VersionId,
						LegalHold: &s3.ObjectLockLegalHold{
							Status: aws.String(s3.ObjectLockLegalHoldStatusOff),
						},
					})

					if err != nil {
						log.Printf("[ERROR] Error putting Spaces Bucket (%s) Object (%s) Version(%s) legal hold: %s", bucketName, objectKey, objectVersionID, err)
						lastErr = err
						continue
					}

					// Attempt to delete again.
					err = deleteS3ObjectVersion(conn, bucketName, objectKey, objectVersionID, force)

					if err != nil {
						lastErr = err
					}

					continue
				}

				// AccessDenied for another reason.
				lastErr = fmt.Errorf("AccessDenied deleting Spaces Bucket (%s) Object (%s) Version: %s", bucketName, objectKey, objectVersionID)
				continue
			}

			if err != nil {
				lastErr = err
			}
		}

		return !lastPage
	})

	if isAWSErr(err, s3.ErrCodeNoSuchBucket, "") {
		err = nil
	}

	if err != nil {
		return err
	}

	if lastErr != nil {
		if !ignoreObjectErrors {
			return fmt.Errorf("error deleting at least one object version, last error: %s", lastErr)
		}

		lastErr = nil
	}

	err = conn.ListObjectVersionsPages(input, func(page *s3.ListObjectVersionsOutput, lastPage bool) bool {
		if page == nil {
			return !lastPage
		}

		for _, deleteMarker := range page.DeleteMarkers {
			deleteMarkerKey := aws.StringValue(deleteMarker.Key)
			deleteMarkerVersionID := aws.StringValue(deleteMarker.VersionId)

			if key != "" && key != deleteMarkerKey {
				continue
			}

			// Delete markers have no object lock protections.
			err := deleteS3ObjectVersion(conn, bucketName, deleteMarkerKey, deleteMarkerVersionID, false)

			if err != nil {
				lastErr = err
			}
		}

		return !lastPage
	})

	if isAWSErr(err, s3.ErrCodeNoSuchBucket, "") {
		err = nil
	}

	if err != nil {
		return err
	}

	if lastErr != nil {
		if !ignoreObjectErrors {
			return fmt.Errorf("error deleting at least one object delete marker, last error: %s", lastErr)
		}

		lastErr = nil
	}

	return nil
}

// deleteS3ObjectVersion deletes a specific bucket object version.
// Set force to true to override any S3 object lock protections.
func deleteS3ObjectVersion(conn *s3.S3, b, k, v string, force bool) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(b),
		Key:    aws.String(k),
	}

	if v != "" {
		input.VersionId = aws.String(v)
	}

	if force {
		input.BypassGovernanceRetention = aws.Bool(true)
	}

	log.Printf("[INFO] Deleting Spaces Bucket (%s) Object (%s) Version: %s", b, k, v)
	_, err := conn.DeleteObject(input)

	if err != nil {
		log.Printf("[WARN] Error deleting Spaces Bucket (%s) Object (%s) Version (%s): %s", b, k, v, err)
	}

	if isAWSErr(err, s3.ErrCodeNoSuchBucket, "") || isAWSErr(err, s3.ErrCodeNoSuchKey, "") {
		return nil
	}

	return err
}

func stringMapToPointers(m map[string]interface{}) map[string]*string {
	list := make(map[string]*string, len(m))
	for i, v := range m {
		list[i] = aws.String(v.(string))
	}
	return list
}

func pointersMapToStringList(pointers map[string]*string) map[string]interface{} {
	list := make(map[string]interface{}, len(pointers))
	for i, v := range pointers {
		list[i] = *v
	}
	return list
}
