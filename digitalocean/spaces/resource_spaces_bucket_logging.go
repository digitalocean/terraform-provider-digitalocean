package spaces

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanSpacesBucketLogging() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSpacesBucketLoggingCreate,
		ReadContext:   resourceSpacesBucketLoggingRead,
		UpdateContext: resourceSpacesBucketLoggingUpdate,
		DeleteContext: resourceSpacesBucketLoggingDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDigitalOceanAccessLoggingImport,
		},
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(SpacesRegions, true),
			},
			"target_bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"target_prefix": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSpacesBucketLoggingCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	region := d.Get("region").(string)
	client, err := meta.(*config.CombinedConfig).SpacesClient(region)
	if err != nil {
		return diag.Errorf("Error creating Spaces client: %s", err)
	}

	svc := s3.New(client)

	bucket := d.Get("bucket").(string)
	input := &s3.PutBucketLoggingInput{
		Bucket: aws.String(bucket),
		BucketLoggingStatus: &s3.BucketLoggingStatus{
			LoggingEnabled: &s3.LoggingEnabled{
				TargetBucket: aws.String(d.Get("target_bucket").(string)),
				TargetPrefix: aws.String(d.Get("target_prefix").(string)),
			},
		},
	}

	err = retry.RetryContext(ctx, 5*time.Minute, func() *retry.RetryError {
		log.Printf("[DEBUG] Trying to enable access logging for Spaces Bucket (%s)", bucket)
		_, err := svc.PutBucketLogging(input)
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "NotFound" {
				log.Printf("[DEBUG] Waiting for Spaces Bucket (%s) to be available, retrying: %v", bucket, awsErr.Message())
				return retry.RetryableError(err)
			}
		}
		if err != nil {
			return retry.NonRetryableError(err)
		}

		return nil
	})

	if err != nil {
		return diag.Errorf("Error enabling Spaces access logging: %s", err)
	}

	d.SetId(bucket)

	return resourceSpacesBucketLoggingRead(ctx, d, meta)
}

func resourceSpacesBucketLoggingRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	region := d.Get("region").(string)
	client, err := meta.(*config.CombinedConfig).SpacesClient(region)
	if err != nil {
		return diag.Errorf("Error creating Spaces client: %s", err)
	}

	svc := s3.New(client)
	bucket := d.Id()
	input := &s3.GetBucketLoggingInput{
		Bucket: aws.String(bucket),
	}

	response, err := svc.GetBucketLogging(input)
	if err != nil {
		if awsError, ok := err.(awserr.RequestFailure); ok && awsError.StatusCode() == 404 {
			log.Printf("[WARN] Spaces Bucket (%s) not found; removing access logging from state", d.Id())
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("Error reading access logging for Spaces Bucket (%s): %s", d.Id(), err)
		}
	}

	d.Set("bucket", bucket)
	if response.LoggingEnabled != nil {
		d.Set("target_bucket", response.LoggingEnabled.TargetBucket)
		d.Set("target_prefix", response.LoggingEnabled.TargetPrefix)
	} else {
		d.Set("target_bucket", "")
		d.Set("target_prefix", "")
	}

	return nil
}

func resourceSpacesBucketLoggingUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	region := d.Get("region").(string)
	client, err := meta.(*config.CombinedConfig).SpacesClient(region)
	if err != nil {
		return diag.Errorf("Error creating Spaces client: %s", err)
	}

	svc := s3.New(client)
	bucket := d.Id()

	input := &s3.PutBucketLoggingInput{
		Bucket: aws.String(bucket),
		BucketLoggingStatus: &s3.BucketLoggingStatus{
			LoggingEnabled: &s3.LoggingEnabled{
				TargetBucket: aws.String(d.Get("target_bucket").(string)),
				TargetPrefix: aws.String(d.Get("target_prefix").(string)),
			},
		},
	}

	_, err = svc.PutBucketLogging(input)
	if err != nil {
		return diag.Errorf("Error updating Spaces Bucket (%s) logging: %s", bucket, err)
	}

	return nil
}

func resourceSpacesBucketLoggingDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	region := d.Get("region").(string)
	client, err := meta.(*config.CombinedConfig).SpacesClient(region)
	if err != nil {
		return diag.Errorf("Error creating Spaces client: %s", err)
	}

	svc := s3.New(client)
	bucket := d.Id()

	input := &s3.PutBucketLoggingInput{
		Bucket:              aws.String(bucket),
		BucketLoggingStatus: &s3.BucketLoggingStatus{},
	}

	_, err = svc.PutBucketLogging(input)
	if err != nil {
		if awsError, ok := err.(awserr.RequestFailure); ok && awsError.StatusCode() == 404 {
			log.Printf("[WARN] Spaces Bucket (%s) not found; removing access logging from state", d.Id())
			d.SetId("")
			return nil
		} else {
			return diag.Errorf("Error disabling access logging for Spaces Bucket (%s): %s", d.Id(), err)
		}
	}

	return nil
}

func resourceDigitalOceanAccessLoggingImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")

		d.SetId(s[1])
		d.Set("region", s[0])
	}

	if d.Id() == "" || d.Get("region") == "" {
		return nil, fmt.Errorf("importing a Spaces Bucket access logging configuration requires the format: <region>,<bucket>")
	}

	return []*schema.ResourceData{d}, nil
}
