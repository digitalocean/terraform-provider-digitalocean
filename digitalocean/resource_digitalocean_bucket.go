package digitalocean

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDigitalOceanBucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanBucketCreate,
		Read:   resourceDigitalOceanBucketRead,
		Update: resourceDigitalOceanBucketUpdate,
		Delete: resourceDigitalOceanBucketDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Bucket name",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Bucket region",
				Default:     "nyc3",
			},
			"profile": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    false,
				Description: "Spaces Access Profile",
			},
			"acl": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Canned ACL applied on bucket creation (Change requires replacement)",
				Default:     "private",
			},
		},
	}
}

func resourceDigitalOceanBucketCreate(d *schema.ResourceData, meta interface{}) error {
	endpoint := fmt.Sprintf("https://%s.digitaloceanspaces.com", d.Get("region").(string))
	sesh, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewSharedCredentials("", d.Get("profile").(string)),
		Endpoint:    aws.String(endpoint)},
	)
	svc := s3.New(sesh)

	if err != nil {
		return err
	}

	input := &s3.CreateBucketInput{
		Bucket: aws.String(d.Get("bucket").(string)),
		ACL:    aws.String(d.Get("acl").(string)),
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		log.Printf("[DEBUG] Trying to create new Spaces bucket: %q", d.Get("bucket").(string))
		_, err := svc.CreateBucket(input)
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "OperationAborted" {
				log.Printf("[WARN] Got an error while trying to create Spaces bucket %s: %s", d.Get("bucket").(string), err)
				return resource.RetryableError(
					fmt.Errorf("[WARN] Error creating Spaces bucket %s, retrying: %s",
						d.Get("bucket").(string), err))
			}
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("Error creating Spaces bucket: %s", err)
	}
	log.Println("Bucket created")

	d.SetId(d.Get("bucket").(string))
	return resourceDigitalOceanBucketRead(d, meta)
}

func resourceDigitalOceanBucketUpdate(d *schema.ResourceData, meta interface{}) error {
	endpoint := fmt.Sprintf("https://%s.digitaloceanspaces.com", d.Get("region").(string))
	sesh, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewSharedCredentials("", d.Get("profile").(string)),
		Endpoint:    aws.String(endpoint)},
	)
	svc := s3.New(sesh)

	if err != nil {
		return err
	}

	if d.HasChange("acl") {
		if err := resourceDigitalOceanBucketAclUpdate(svc, d); err != nil {
			return err
		}
	}

	return resourceDigitalOceanBucketRead(d, meta)
}

func resourceDigitalOceanBucketRead(d *schema.ResourceData, meta interface{}) error {
	endpoint := fmt.Sprintf("https://%s.digitaloceanspaces.com", d.Get("region").(string))
	sesh, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewSharedCredentials("", d.Get("profile").(string)),
		Endpoint:    aws.String(endpoint)},
	)
	svc := s3.New(sesh)

	if err != nil {
		return err
	}

	_, err = retryOnAwsCode("NoSuchBucket", func() (interface{}, error) {
		return svc.HeadBucket(&s3.HeadBucketInput{
			Bucket: aws.String(d.Id()),
		})
	})
	if err != nil {
		if awsError, ok := err.(awserr.RequestFailure); ok && awsError.StatusCode() == 404 {
			log.Printf("[WARN] Spaces Bucket (%s) not found, error code (404)", d.Id())
			d.SetId("")
			return nil
		} else {
			// some of the AWS SDK's errors can be empty strings, so let's add
			// some additional context.
			return fmt.Errorf("error reading Spaces bucket \"%s\": %s", d.Id(), err)
		}
	}

	// In the import case, we won't have this
	if _, ok := d.GetOk("bucket"); !ok {
		d.Set("bucket", d.Id())
	}

	d.Set("bucket_domain_name", bucketDomainName(d.Get("bucket").(string), d.Get("region").(string)))

	// Add the region as an attribute

	locationResponse, err := retryOnAwsCode("NoSuchBucket", func() (interface{}, error) {
		return svc.GetBucketLocation(
			&s3.GetBucketLocationInput{
				Bucket: aws.String(d.Id()),
			},
		)
	})
	if err != nil {
		return err
	}
	location := locationResponse.(*s3.GetBucketLocationOutput)
	var region string
	if location.LocationConstraint != nil {
		region = *location.LocationConstraint
	}
	region = normalizeRegion(region)
	if err := d.Set("region", region); err != nil {
		return err
	}

	d.Set("bucket", d.Get("bucket").(string))

	return nil
}

func resourceDigitalOceanBucketDelete(d *schema.ResourceData, meta interface{}) error {
	endpoint := fmt.Sprintf("https://%s.digitaloceanspaces.com", d.Get("region").(string))
	sesh, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewSharedCredentials("", d.Get("profile").(string)),
		Endpoint:    aws.String(endpoint)},
	)
	svc := s3.New(sesh)

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Spaces Delete Bucket: %s", d.Id())
	_, err = svc.DeleteBucket(&s3.DeleteBucketInput{
		Bucket: aws.String(d.Id()),
	})
	if err != nil {
		ec2err, ok := err.(awserr.Error)
		if ok && ec2err.Code() == "BucketNotEmpty" {
			if d.Get("force_destroy").(bool) {
				// bucket may have things delete them
				log.Printf("[DEBUG] Spaces Bucket attempting to forceDestroy %+v", err)

				bucket := d.Get("bucket").(string)
				resp, err := svc.ListObjectVersions(
					&s3.ListObjectVersionsInput{
						Bucket: aws.String(bucket),
					},
				)

				if err != nil {
					return fmt.Errorf("Error Spaces Bucket list Object Versions err: %s", err)
				}

				objectsToDelete := make([]*s3.ObjectIdentifier, 0)

				if len(resp.DeleteMarkers) != 0 {

					for _, v := range resp.DeleteMarkers {
						objectsToDelete = append(objectsToDelete, &s3.ObjectIdentifier{
							Key:       v.Key,
							VersionId: v.VersionId,
						})
					}
				}

				if len(resp.Versions) != 0 {
					for _, v := range resp.Versions {
						objectsToDelete = append(objectsToDelete, &s3.ObjectIdentifier{
							Key:       v.Key,
							VersionId: v.VersionId,
						})
					}
				}

				params := &s3.DeleteObjectsInput{
					Bucket: aws.String(bucket),
					Delete: &s3.Delete{
						Objects: objectsToDelete,
					},
				}

				_, err = svc.DeleteObjects(params)

				if err != nil {
					return fmt.Errorf("Error Spaces Bucket force_destroy error deleting: %s", err)
				}

				// this line recurses until all objects are deleted or an error is returned
				return resourceDigitalOceanBucketDelete(d, meta)
			}
		}
		return fmt.Errorf("Error deleting Spaces Bucket: %s %q", err, d.Get("bucket").(string))
	}
	log.Println("Bucket destroyed")

	d.SetId("")
	return nil
}

func resourceDigitalOceanBucketAclUpdate(svc *s3.S3, d *schema.ResourceData) error {
	acl := d.Get("acl").(string)
	bucket := d.Get("bucket").(string)

	i := &s3.PutBucketAclInput{
		Bucket: aws.String(bucket),
		ACL:    aws.String(acl),
	}
	log.Printf("[DEBUG] Spaces put bucket ACL: %#v", i)

	_, err := retryOnAwsCode("NoSuchBucket", func() (interface{}, error) {
		return svc.PutBucketAcl(i)
	})
	if err != nil {
		return fmt.Errorf("Error putting Spaces ACL: %s", err)
	}

	return nil
}

func bucketDomainName(bucket string, region string) string {
	return fmt.Sprintf("%q.%q.digitaloceanspaces.com", bucket, region)
}

func retryOnAwsCode(code string, f func() (interface{}, error)) (interface{}, error) {
	var resp interface{}
	err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		var err error
		resp, err = f()
		if err != nil {
			awsErr, ok := err.(awserr.Error)
			if ok && awsErr.Code() == code {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	return resp, err
}

func normalizeRegion(region string) string {
	// Default to nyc3 if the bucket doesn't have a region:
	if region == "" {
		region = "nyc3"
	}

	return region
}
