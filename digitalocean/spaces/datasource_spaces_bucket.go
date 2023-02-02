package spaces

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDigitalOceanSpacesBucket() *schema.Resource {
	recordSchema := spacesBucketSchema()

	for _, f := range recordSchema {
		f.Computed = true
	}

	recordSchema["region"].Required = true
	recordSchema["region"].Computed = false
	recordSchema["region"].ValidateFunc = validation.StringInSlice(SpacesRegions, true)
	recordSchema["name"].Required = true
	recordSchema["name"].Computed = false

	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanSpacesBucketRead,
		Schema:      recordSchema,
	}
}

func dataSourceDigitalOceanSpacesBucketRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	region := d.Get("region").(string)
	name := d.Get("name").(string)

	client, err := meta.(*config.CombinedConfig).SpacesClient(region)
	if err != nil {
		return diag.Errorf("Error reading bucket: %s", err)
	}

	svc := s3.New(client)

	_, err = retryOnAwsCode("NoSuchBucket", func() (interface{}, error) {
		return svc.HeadBucket(&s3.HeadBucketInput{
			Bucket: aws.String(name),
		})
	})
	if err != nil {
		if awsError, ok := err.(awserr.RequestFailure); ok && awsError.StatusCode() == 404 {
			d.SetId("")
			return diag.Errorf("Spaces Bucket (%s) not found", name)
		} else {
			// some of the AWS SDK's errors can be empty strings, so let's add
			// some additional context.
			return diag.Errorf("error reading Spaces bucket \"%s\": %s", d.Id(), err)
		}
	}

	metadata := bucketMetadataStruct{
		name:   name,
		region: region,
	}

	flattenedBucket, err := flattenSpacesBucket(&metadata, meta, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := util.SetResourceDataFromMap(d, flattenedBucket); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(name)
	return nil
}
