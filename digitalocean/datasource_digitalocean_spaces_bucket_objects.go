package digitalocean

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const keyRequestPageSize = 1000

func dataSourceDigitalOceanSpacesBucketObjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanSpacesBucketObjectsRead,

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"prefix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"delimiter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encoding_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_keys": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1000,
			},

			// computed attributes

			"keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"common_prefixes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"owners": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceDigitalOceanSpacesBucketObjectsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	region := d.Get("region").(string)
	client, err := meta.(*CombinedConfig).spacesClient(region)
	if err != nil {
		return diag.FromErr(err)
	}

	conn := s3.New(client)

	bucket := d.Get("bucket").(string)
	prefix := d.Get("prefix").(string)

	d.SetId(resource.UniqueId())

	listInput := s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	}

	if prefix != "" {
		listInput.Prefix = aws.String(prefix)
	}

	if s, ok := d.GetOk("delimiter"); ok {
		listInput.Delimiter = aws.String(s.(string))
	}

	if s, ok := d.GetOk("encoding_type"); ok {
		listInput.EncodingType = aws.String(s.(string))
	}

	// "listInput.MaxKeys" refers to max keys returned in a single request
	// (i.e., page size), not the total number of keys returned if you page
	// through the results. "maxKeys" does refer to total keys returned.
	maxKeys := int64(d.Get("max_keys").(int))
	if maxKeys <= keyRequestPageSize {
		listInput.MaxKeys = aws.Int64(maxKeys)
	}

	var commonPrefixes []string
	var keys []string
	var owners []string

	err = conn.ListObjectsPages(&listInput, func(page *s3.ListObjectsOutput, lastPage bool) bool {
		for _, commonPrefix := range page.CommonPrefixes {
			commonPrefixes = append(commonPrefixes, aws.StringValue(commonPrefix.Prefix))
		}

		for _, object := range page.Contents {
			keys = append(keys, aws.StringValue(object.Key))

			if object.Owner != nil {
				owners = append(owners, aws.StringValue(object.Owner.ID))
			}
		}

		maxKeys = maxKeys - int64(len(page.Contents))

		if maxKeys <= keyRequestPageSize {
			listInput.MaxKeys = aws.Int64(maxKeys)
		}

		return !lastPage
	})

	if err != nil {
		return diag.Errorf("error listing Spaces Bucket (%s) Objects: %s", bucket, err)
	}

	if err := d.Set("common_prefixes", commonPrefixes); err != nil {
		return diag.Errorf("error setting common_prefixes: %s", err)
	}

	if err := d.Set("keys", keys); err != nil {
		return diag.Errorf("error setting keys: %s", err)
	}

	if err := d.Set("owners", owners); err != nil {
		return diag.Errorf("error setting owners: %s", err)
	}

	return nil
}
