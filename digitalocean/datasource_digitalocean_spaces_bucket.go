package digitalocean

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceDigitalOceanSpacesBucket() *schema.Resource {
	recordSchema := spacesBucketSchema()

	for _, f := range recordSchema {
		f.Computed = true
	}

	recordSchema["region"].Required = true
	recordSchema["region"].Computed = false
	recordSchema["name"].Required = true
	recordSchema["name"].Computed = false

	return &schema.Resource{
		Read:   dataSourceDigitalOceanSpacesBucketRead,
		Schema: recordSchema,
	}
}

func dataSourceDigitalOceanSpacesBucketRead(d *schema.ResourceData, meta interface{}) error {
	region := d.Get("region").(string)
	name := d.Get("name").(string)

	client, err := meta.(*CombinedConfig).spacesClient(region)
	if err != nil {
		return fmt.Errorf("Error reading bucket: %s", err)
	}

	svc := s3.New(client)

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

	d.Set("bucket_domain_name", bucketDomainName(d.Get("name").(string), d.Get("region").(string)))

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
	if location.LocationConstraint != nil {
		region = *location.LocationConstraint
	}
	region = normalizeRegion(region)
	if err := d.Set("region", region); err != nil {
		return err
	}

	// Get the bucket's ACLs.
	aclResponse, err := svc.GetBucketAcl(
			&s3.GetBucketAclInput{
				Bucket: aws.String(d.Id()),
			})
	if err != nil {
		return err
	}
	if err = d.Set("acl", aclResponse.); err != nil {
		return err
	}



	urn := fmt.Sprintf("do:space:%s", d.Get("name"))
	d.Set("urn", urn)

	return nil
}
