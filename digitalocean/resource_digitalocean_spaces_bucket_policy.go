package digitalocean

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDigitalOceanSpacesBucketPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanBucketPolicyCreate,
		ReadContext:   resourceDigitalOceanBucketPolicyRead,
		UpdateContext: resourceDigitalOceanBucketPolicyUpdate,
		DeleteContext: resourceDigitalOceanBucketPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDigitalOceanBucketPolicyImport,
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
			"policy": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsJSON,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return compareSpacesBucketPolicy(old, new)
				},
				StateFunc: func(v interface{}) string {
					json, _ := structure.NormalizeJsonString(v)
					return json
				},
			},
		},
	}
}

func s3connFromSpacesBucketPolicyResourceData(d *schema.ResourceData, meta interface{}) (*s3.S3, error) {
	region := d.Get("region").(string)

	client, err := meta.(*CombinedConfig).spacesClient(region)
	if err != nil {
		return nil, err
	}

	svc := s3.New(client)
	return svc, nil
}

func resourceDigitalOceanBucketPolicyImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")

		d.SetId(s[1])
		d.Set("region", s[0])
	}

	if d.Id() == "" || d.Get("region") == "" {
		return nil, fmt.Errorf("importing a Spaces bucket policy requires the format: <region>,<bucket>")
	}

	return []*schema.ResourceData{d}, nil
}

func resourceDigitalOceanBucketPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := s3connFromSpacesBucketPolicyResourceData(d, meta)
	if err != nil {
		return diag.Errorf("Error occurred while creating new Spaces bucket policy: %s", err)
	}

	bucket := d.Get("bucket").(string)
	policy := d.Get("policy").(string)

	if policy == "" {
		return diag.Errorf("Spaces bucket policy must not be empty")
	}

	log.Printf("[DEBUG] Trying to create new Spaces bucket policy for bucket: %s, policy: %s", bucket, policy)
	_, err = conn.PutBucketPolicy(&s3.PutBucketPolicyInput{
		Bucket: aws.String(bucket),
		Policy: aws.String(policy),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "NoSuchKey" {
			return diag.Errorf("Unable to create new Spaces bucket policy because bucket '%s' does not exist", bucket)
		}
		return diag.Errorf("Error occurred while creating new Spaces bucket policy: %s", err)
	}

	d.SetId(bucket)
	return resourceDigitalOceanBucketPolicyRead(ctx, d, meta)
}

func resourceDigitalOceanBucketPolicyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := s3connFromSpacesBucketPolicyResourceData(d, meta)
	if err != nil {
		return diag.Errorf("Error occurred while fetching Spaces bucket policy: %s", err)
	}

	log.Printf("[DEBUG] Trying to fetch Spaces bucket policy for bucket: %s", d.Id())
	response, err := conn.GetBucketPolicy(&s3.GetBucketPolicyInput{
		Bucket: aws.String(d.Id()),
	})

	policy := ""
	if err == nil && response.Policy != nil {
		policy = aws.StringValue(response.Policy)
	}

	d.Set("bucket", d.Id())
	d.Set("policy", policy)

	return nil
}

func resourceDigitalOceanBucketPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceDigitalOceanBucketPolicyCreate(ctx, d, meta)
}

func resourceDigitalOceanBucketPolicyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := s3connFromSpacesBucketPolicyResourceData(d, meta)
	if err != nil {
		return diag.Errorf("Error occurred while deleting Spaces bucket policy: %s", err)
	}

	bucket := d.Id()

	log.Printf("[DEBUG] Trying to delete Spaces bucket policy for bucket: %s", d.Id())
	_, err = conn.DeleteBucketPolicy(&s3.DeleteBucketPolicyInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "BucketDeleted" {
			return diag.Errorf("Unable to remove Spaces bucket policy because bucket '%s' is already deleted", bucket)
		}
		return diag.Errorf("Error occurred while deleting Spaces Bucket policy: %s", err)
	}
	return nil
}
