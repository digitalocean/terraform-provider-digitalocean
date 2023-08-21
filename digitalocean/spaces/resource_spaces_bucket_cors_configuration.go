package spaces

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanBucketCorsConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanBucketCorsConfigurationCreate,
		ReadContext:   resourceDigitalOceanBucketCorsConfigurationRead,
		UpdateContext: resourceDigitalOceanBucketCorsConfigurationUpdate,
		DeleteContext: resourceBucketCorsConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDigitalOceanBucketImport,
		},

		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Bucket ID",
			},
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(SpacesRegions, true),
			},
			"cors_rule": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A container holding a list of elements describing allowed methods for a specific origin.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allowed_methods": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "A list of HTTP methods (e.g. GET) which are allowed from the specified origin.",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"allowed_origins": {
							Type:        schema.TypeList,
							Required:    true,
							Description: "A list of hosts from which requests using the specified methods are allowed. A host may contain one wildcard (e.g. http://*.example.com).",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"allowed_headers": {
							Type:        schema.TypeList,
							Optional:    true,
							Description: "A list of headers that will be included in the CORS preflight request's Access-Control-Request-Headers. A header may contain one wildcard (e.g. x-amz-*).",
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"max_age_seconds": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceDigitalOceanBucketCorsConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := s3connFromSpacesBucketCorsResourceData(d, meta)
	if err != nil {
		return diag.Errorf("Error occurred while creating new Spaces bucket policy: %s", err)
	}

	bucket := d.Get("bucket").(string)
	rawCors := d.Get("cors_rule").([]interface{})

	if len(rawCors) == 0 {
		return diag.Errorf("Spaces bucket CORS must not be empty")
	}

	input := &s3.PutBucketCorsInput{
		Bucket: aws.String(bucket),
		CORSConfiguration: &s3.CORSConfiguration{
			CORSRules: expandBucketCorsConfigurationCorsRules(rawCors, bucket),
		},
	}

	log.Printf("[DEBUG] Trying to configure CORS for Spaces bucket: %s", bucket)
	_, err = conn.PutBucketCorsWithContext(ctx, input)

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "NoSuchKey" {
			return diag.Errorf("Unable to configure CORS for Spaces bucket because the bucket does not exist: '%s'", bucket)
		}
		return diag.Errorf("Error occurred while configuring CORS for Spaces bucket: %s", err)
	}

	d.SetId(bucket)
	return resourceDigitalOceanBucketPolicyRead(ctx, d, meta)
}

func resourceDigitalOceanBucketCorsConfigurationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := s3connFromSpacesBucketCorsResourceData(d, meta)
	if err != nil {
		return diag.Errorf("Error occurred while fetching Spaces bucket CORS configuration: %s", err)
	}

	log.Printf("[DEBUG] Trying to fetch Spaces bucket CORS configuration for bucket: %s", d.Id())
	response, err := conn.GetBucketCorsWithContext(ctx, &s3.GetBucketCorsInput{
		Bucket: aws.String(d.Id()),
	})

	if err != nil {
		return diag.Errorf("Error occurred while fetching Spaces bucket CORS configuration: %s", err)
	}

	d.Set("bucket", d.Id())

	if err := d.Set("cors_rule", flattenBucketCorsConfigurationCorsRules(response.CORSRules)); err != nil {
		return diag.Errorf("setting cors_rule: %s", err)
	}

	return nil
}

func resourceDigitalOceanBucketCorsConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceDigitalOceanBucketCorsConfigurationCreate(ctx, d, meta)
}

func resourceBucketCorsConfigurationDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conn, err := s3connFromSpacesBucketCorsResourceData(d, meta)
	if err != nil {
		return diag.Errorf("Error occurred while deleting Spaces bucket CORS configuration: %s", err)
	}

	bucket := d.Id()

	log.Printf("[DEBUG] Trying to delete Spaces bucket CORS Configuration for bucket: %s", d.Id())
	_, err = conn.DeleteBucketCorsWithContext(ctx, &s3.DeleteBucketCorsInput{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "BucketDeleted" {
			return diag.Errorf("Unable to remove Spaces bucket CORS configuration because bucket '%s' is already deleted", bucket)
		}
		return diag.Errorf("Error occurred while deleting Spaces Bucket CORS configuration: %s", err)
	}
	return nil

}

func expandBucketCorsConfigurationCorsRules(rawCors []interface{}, bucketName string) []*s3.CORSRule {
	if len(rawCors) == 0 {
		return nil
	}

	rules := make([]*s3.CORSRule, 0, len(rawCors))
	for _, cors := range rawCors {
		corsMap := cors.(map[string]interface{})
		r := &s3.CORSRule{}
		for k, v := range corsMap {
			log.Printf("[DEBUG] Spaces bucket: %s, put CORS: %#v, %#v", bucketName, k, v)
			if k == "max_age_seconds" {
				r.MaxAgeSeconds = aws.Int64(int64(v.(int)))
			} else {
				vMap := make([]*string, len(v.([]interface{})))
				for i, vv := range v.([]interface{}) {
					str := vv.(string)
					vMap[i] = aws.String(str)
				}
				switch k {
				case "allowed_headers":
					r.AllowedHeaders = vMap
				case "allowed_methods":
					r.AllowedMethods = vMap
				case "allowed_origins":
					r.AllowedOrigins = vMap
				case "expose_headers":
					r.ExposeHeaders = vMap
				}
			}
		}
		rules = append(rules, r)
	}
	return rules
}

func s3connFromSpacesBucketCorsResourceData(d *schema.ResourceData, meta interface{}) (*s3.S3, error) {
	region := d.Get("region").(string)

	client, err := meta.(*config.CombinedConfig).SpacesClient(region)
	if err != nil {
		return nil, err
	}

	svc := s3.New(client)
	return svc, nil
}

func flattenBucketCorsConfigurationCorsRules(rules []*s3.CORSRule) []interface{} {
	var results []interface{}

	for _, rule := range rules {
		if rule == nil {
			continue
		}

		m := make(map[string]interface{})

		if len(rule.AllowedHeaders) > 0 {
			m["allowed_headers"] = flattenStringSet(rule.AllowedHeaders)
		}

		if len(rule.AllowedMethods) > 0 {
			m["allowed_methods"] = flattenStringSet(rule.AllowedMethods)
		}

		if len(rule.AllowedOrigins) > 0 {
			m["allowed_origins"] = flattenStringSet(rule.AllowedOrigins)
		}

		if len(rule.ExposeHeaders) > 0 {
			m["expose_headers"] = flattenStringSet(rule.ExposeHeaders)
		}

		if rule.ID != nil {
			m["id"] = aws.StringValue(rule.ID)
		}

		if rule.MaxAgeSeconds != nil {
			m["max_age_seconds"] = aws.Int64Value(rule.MaxAgeSeconds)
		}

		results = append(results, m)
	}

	return results
}

func flattenStringSet(list []*string) *schema.Set {
	return schema.NewSet(schema.HashString, flattenStringList(list)) // nosemgrep:ci.helper-schema-Set-extraneous-NewSet-with-FlattenStringList
}

// flattenStringList takes list of pointers to strings. Expand to an array
// of raw strings and returns a []interface{}
// to keep compatibility w/ schema.NewSetschema.NewSet
func flattenStringList(list []*string) []interface{} {
	vs := make([]interface{}, 0, len(list))
	for _, v := range list {
		vs = append(vs, *v)
	}
	return vs
}
