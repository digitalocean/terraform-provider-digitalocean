package digitalocean

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceDigitalOceanBucket() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanBucketCreate,
		Read:   resourceDigitalOceanBucketRead,
		Update: resourceDigitalOceanBucketUpdate,
		Delete: resourceDigitalOceanBucketDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDigitalOceanBucketImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Bucket name",
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the uniform resource name for the bucket",
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Bucket region",
				Default:     "nyc3",
				StateFunc: func(val interface{}) string {
					// DO API V2 region slug is always lowercase
					return strings.ToLower(val.(string))
				},
			},
			"acl": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Canned ACL applied on bucket creation",
				Default:     "private",
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
			// This is structured as a subobject in case Spaces supports more of s3.VersioningConfiguration
			// than just enabling bucket versioning.
			"versioning": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old == "1" && new == "0"
				},
			},

			"bucket_domain_name": {
				Type:        schema.TypeString,
				Description: "The FQDN of the bucket",
				Computed:    true,
			},

			"lifecycle_rule": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.StringLenBetween(0, 255),
						},
						"prefix": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"abort_incomplete_multipart_upload_days": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"expiration": {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      expirationHash,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validateS3BucketLifecycleTimestamp,
									},
									"days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(0),
									},
									"expired_object_delete_marker": {
										Type:     schema.TypeBool,
										Optional: true,
									},
								},
							},
						},
						"noncurrent_version_expiration": {
							Type:     schema.TypeSet,
							MaxItems: 1,
							Optional: true,
							Set:      expirationHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(1),
									},
								},
							},
						},
						"transition": {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      transitionHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"date": {
										Type:         schema.TypeString,
										Optional:     true,
										ValidateFunc: validateS3BucketLifecycleTimestamp,
									},
									"days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(0),
									},
									"storage_class": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateS3BucketLifecycleTransitionStorageClass(),
									},
								},
							},
						},
						"noncurrent_version_transition": {
							Type:     schema.TypeSet,
							Optional: true,
							Set:      transitionHash,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"days": {
										Type:         schema.TypeInt,
										Optional:     true,
										ValidateFunc: validation.IntAtLeast(0),
									},
									"storage_class": {
										Type:         schema.TypeString,
										Required:     true,
										ValidateFunc: validateS3BucketLifecycleTransitionStorageClass(),
									},
								},
							},
						},
					},
				},
			},

			"force_destroy": {
				Type:        schema.TypeBool,
				Description: "Unless true, the bucket will only be destroyed if empty",
				Optional:    true,
				Default:     false,
			},
		},
	}
}

func resourceDigitalOceanBucketCreate(d *schema.ResourceData, meta interface{}) error {
	region := d.Get("region").(string)
	client, err := meta.(*CombinedConfig).spacesClient(region)

	if err != nil {
		return fmt.Errorf("Error creating bucket: %s", err)
	}

	svc := s3.New(client)

	input := &s3.CreateBucketInput{
		Bucket: aws.String(d.Get("name").(string)),
		ACL:    aws.String(d.Get("acl").(string)),
	}

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		log.Printf("[DEBUG] Trying to create new Spaces bucket: %q", d.Get("name").(string))
		_, err := svc.CreateBucket(input)
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == "OperationAborted" {
				log.Printf("[WARN] Got an error while trying to create Spaces bucket %s: %s", d.Get("name").(string), err)
				return resource.RetryableError(
					fmt.Errorf("[WARN] Error creating Spaces bucket %s, retrying: %s",
						d.Get("name").(string), err))
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

	d.SetId(d.Get("name").(string))
	return resourceDigitalOceanBucketUpdate(d, meta)
}

func resourceDigitalOceanBucketUpdate(d *schema.ResourceData, meta interface{}) error {
	region := d.Get("region").(string)
	client, err := meta.(*CombinedConfig).spacesClient(region)

	if err != nil {
		return fmt.Errorf("Error updating bucket: %s", err)
	}

	svc := s3.New(client)

	if d.HasChange("acl") {
		if err := resourceDigitalOceanBucketACLUpdate(svc, d); err != nil {
			return err
		}
	}

	if d.HasChange("cors_rule") {
		if err := resourceDigitalOceanBucketCorsUpdate(svc, d); err != nil {
			return err
		}
	}

	if d.HasChange("versioning") {
		if err := resourceDigitalOceanSpacesBucketVersioningUpdate(svc, d); err != nil {
			return err
		}
	}

	if d.HasChange("lifecycle_rule") {
		if err := resourceDigitalOceanBucketLifecycleUpdate(svc, d); err != nil {
			return err
		}
	}

	return resourceDigitalOceanBucketRead(d, meta)
}

func resourceDigitalOceanBucketRead(d *schema.ResourceData, meta interface{}) error {
	region := d.Get("region").(string)
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

	// In the import case, we won't have this
	if _, ok := d.GetOk("name"); !ok {
		d.Set("name", d.Id())
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

	// Read the versioning configuration
	versioningResponse, err := retryOnAwsCode(s3.ErrCodeNoSuchBucket, func() (interface{}, error) {
		return svc.GetBucketVersioning(&s3.GetBucketVersioningInput{
			Bucket: aws.String(d.Id()),
		})
	})
	if err != nil {
		return err
	}
	vcl := make([]map[string]interface{}, 0, 1)
	if versioning, ok := versioningResponse.(*s3.GetBucketVersioningOutput); ok {
		vc := make(map[string]interface{})
		if versioning.Status != nil && *versioning.Status == s3.BucketVersioningStatusEnabled {
			vc["enabled"] = true
		} else {
			vc["enabled"] = false
		}
		vcl = append(vcl, vc)
	}
	if err := d.Set("versioning", vcl); err != nil {
		return fmt.Errorf("error setting versioning: %s", err)
	}

	// Read the lifecycle configuration
	lifecycleResponse, err := retryOnAwsCode(s3.ErrCodeNoSuchBucket, func() (interface{}, error) {
		return svc.GetBucketLifecycleConfiguration(&s3.GetBucketLifecycleConfigurationInput{
			Bucket: aws.String(d.Id()),
		})
	})
	if err != nil && !isAWSErr(err, "NoSuchLifecycleConfiguration", "") {
		return err
	}

	lifecycleRules := make([]map[string]interface{}, 0)
	if lifecycle, ok := lifecycleResponse.(*s3.GetBucketLifecycleConfigurationOutput); ok && len(lifecycle.Rules) > 0 {
		lifecycleRules = make([]map[string]interface{}, 0, len(lifecycle.Rules))

		for _, lifecycleRule := range lifecycle.Rules {
			log.Printf("[DEBUG] Spaces bucket: %s, read lifecycle rule: %v", d.Id(), lifecycleRule)
			rule := make(map[string]interface{})

			// ID
			if lifecycleRule.ID != nil && *lifecycleRule.ID != "" {
				rule["id"] = *lifecycleRule.ID
			}
			filter := lifecycleRule.Filter
			if filter != nil {
				if filter.And != nil {
					// Prefix
					if filter.And.Prefix != nil && *filter.And.Prefix != "" {
						rule["prefix"] = *filter.And.Prefix
					}
					// Tag
					if len(filter.And.Tags) > 0 {
						return fmt.Errorf("unsupported tag found on lifecycle_rule")
					}
				} else {
					// Prefix
					if filter.Prefix != nil && *filter.Prefix != "" {
						rule["prefix"] = *filter.Prefix
					}
					// Tag
					if filter.Tag != nil {
						return fmt.Errorf("unsupported tag found on lifecycle_rule")
					}
				}
			} else {
				if lifecycleRule.Prefix != nil {
					rule["prefix"] = *lifecycleRule.Prefix
				}
			}

			// Enabled
			if lifecycleRule.Status != nil {
				if *lifecycleRule.Status == s3.ExpirationStatusEnabled {
					rule["enabled"] = true
				} else {
					rule["enabled"] = false
				}
			}

			// AbortIncompleteMultipartUploadDays
			if lifecycleRule.AbortIncompleteMultipartUpload != nil {
				if lifecycleRule.AbortIncompleteMultipartUpload.DaysAfterInitiation != nil {
					rule["abort_incomplete_multipart_upload_days"] = int(*lifecycleRule.AbortIncompleteMultipartUpload.DaysAfterInitiation)
				}
			}

			// expiration
			if lifecycleRule.Expiration != nil {
				e := make(map[string]interface{})
				if lifecycleRule.Expiration.Date != nil {
					e["date"] = (*lifecycleRule.Expiration.Date).Format("2006-01-02")
				}
				if lifecycleRule.Expiration.Days != nil {
					e["days"] = int(*lifecycleRule.Expiration.Days)
				}
				if lifecycleRule.Expiration.ExpiredObjectDeleteMarker != nil {
					e["expired_object_delete_marker"] = *lifecycleRule.Expiration.ExpiredObjectDeleteMarker
				}
				rule["expiration"] = schema.NewSet(expirationHash, []interface{}{e})
			}
			// noncurrent_version_expiration
			if lifecycleRule.NoncurrentVersionExpiration != nil {
				e := make(map[string]interface{})
				if lifecycleRule.NoncurrentVersionExpiration.NoncurrentDays != nil {
					e["days"] = int(*lifecycleRule.NoncurrentVersionExpiration.NoncurrentDays)
				}
				rule["noncurrent_version_expiration"] = schema.NewSet(expirationHash, []interface{}{e})
			}
			//// transition
			if len(lifecycleRule.Transitions) > 0 {
				transitions := make([]interface{}, 0, len(lifecycleRule.Transitions))
				for _, v := range lifecycleRule.Transitions {
					t := make(map[string]interface{})
					if v.Date != nil {
						t["date"] = (*v.Date).Format("2006-01-02")
					}
					if v.Days != nil {
						t["days"] = int(*v.Days)
					}
					if v.StorageClass != nil {
						t["storage_class"] = *v.StorageClass
					}
					transitions = append(transitions, t)
				}
				rule["transition"] = schema.NewSet(transitionHash, transitions)
			}
			// noncurrent_version_transition
			if len(lifecycleRule.NoncurrentVersionTransitions) > 0 {
				transitions := make([]interface{}, 0, len(lifecycleRule.NoncurrentVersionTransitions))
				for _, v := range lifecycleRule.NoncurrentVersionTransitions {
					t := make(map[string]interface{})
					if v.NoncurrentDays != nil {
						t["days"] = int(*v.NoncurrentDays)
					}
					if v.StorageClass != nil {
						t["storage_class"] = *v.StorageClass
					}
					transitions = append(transitions, t)
				}
				rule["noncurrent_version_transition"] = schema.NewSet(transitionHash, transitions)
			}

			lifecycleRules = append(lifecycleRules, rule)
		}
	}
	if err := d.Set("lifecycle_rule", lifecycleRules); err != nil {
		return fmt.Errorf("error setting lifecycle_rule: %s", err)
	}

	// Set the bucket's name.
	d.Set("name", d.Get("name").(string))

	// Set the URN attribute.
	urn := fmt.Sprintf("do:space:%s", d.Get("name"))
	d.Set("urn", urn)

	return nil
}

func resourceDigitalOceanBucketDelete(d *schema.ResourceData, meta interface{}) error {
	region := d.Get("region").(string)
	client, err := meta.(*CombinedConfig).spacesClient(region)

	if err != nil {
		return fmt.Errorf("Error deleting bucket: %s", err)
	}

	svc := s3.New(client)

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

				bucket := d.Get("name").(string)
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
		return fmt.Errorf("Error deleting Spaces Bucket: %s %q", err, d.Get("name").(string))
	}
	log.Println("Bucket destroyed")

	d.SetId("")
	return nil
}

func resourceDigitalOceanBucketACLUpdate(svc *s3.S3, d *schema.ResourceData) error {
	acl := d.Get("acl").(string)
	bucket := d.Get("name").(string)

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

func resourceDigitalOceanBucketCorsUpdate(svc *s3.S3, d *schema.ResourceData) error {
	rawCors := d.Get("cors_rule").([]interface{})
	bucket := d.Get("name").(string)

	if len(rawCors) == 0 {
		// Delete CORS
		log.Printf("[DEBUG] Spaces bucket: %s, delete CORS", bucket)
		_, err := svc.DeleteBucketCors(&s3.DeleteBucketCorsInput{
			Bucket: aws.String(bucket),
		})
		if err != nil {
			return fmt.Errorf("Error deleting Spaces CORS: %s", err)
		}
	} else {
		// Put CORS
		rules := make([]*s3.CORSRule, 0, len(rawCors))
		for _, cors := range rawCors {
			corsMap := cors.(map[string]interface{})
			r := &s3.CORSRule{}
			for k, v := range corsMap {
				log.Printf("[DEBUG] Spaces bucket: %s, put CORS: %#v, %#v", bucket, k, v)
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
					}
				}
			}
			rules = append(rules, r)
		}
		corsInput := &s3.PutBucketCorsInput{
			Bucket: aws.String(bucket),
			CORSConfiguration: &s3.CORSConfiguration{
				CORSRules: rules,
			},
		}
		log.Printf("[DEBUG] Spaces bucket: %s, put CORS: %#v", bucket, corsInput)
		_, err := svc.PutBucketCors(corsInput)
		if err != nil {
			return fmt.Errorf("Error putting Spaces CORS: %s", err)
		}
	}

	return nil
}

func resourceDigitalOceanSpacesBucketVersioningUpdate(s3conn *s3.S3, d *schema.ResourceData) error {
	v := d.Get("versioning").([]interface{})
	bucket := d.Get("name").(string)
	vc := &s3.VersioningConfiguration{}

	if len(v) > 0 {
		c := v[0].(map[string]interface{})

		if c["enabled"].(bool) {
			vc.Status = aws.String(s3.BucketVersioningStatusEnabled)
		} else {
			vc.Status = aws.String(s3.BucketVersioningStatusSuspended)
		}
	} else {
		vc.Status = aws.String(s3.BucketVersioningStatusSuspended)
	}

	i := &s3.PutBucketVersioningInput{
		Bucket:                  aws.String(bucket),
		VersioningConfiguration: vc,
	}
	log.Printf("[DEBUG] Spaces PUT bucket versioning: %#v", i)

	_, err := retryOnAwsCode(s3.ErrCodeNoSuchBucket, func() (interface{}, error) {
		return s3conn.PutBucketVersioning(i)
	})
	if err != nil {
		return fmt.Errorf("Error putting Spaces versioning: %s", err)
	}

	return nil
}

func resourceDigitalOceanBucketLifecycleUpdate(s3conn *s3.S3, d *schema.ResourceData) error {
	bucket := d.Get("name").(string)

	lifecycleRules := d.Get("lifecycle_rule").([]interface{})

	if len(lifecycleRules) == 0 {
		i := &s3.DeleteBucketLifecycleInput{
			Bucket: aws.String(bucket),
		}

		_, err := s3conn.DeleteBucketLifecycle(i)
		if err != nil {
			return fmt.Errorf("Error removing S3 lifecycle: %s", err)
		}
		return nil
	}

	rules := make([]*s3.LifecycleRule, 0, len(lifecycleRules))

	for i, lifecycleRule := range lifecycleRules {
		r := lifecycleRule.(map[string]interface{})

		rule := &s3.LifecycleRule{}

		// Filter
		//tags := keyvaluetags.New(r["tags"]).IgnoreAws().S3Tags()
		filter := &s3.LifecycleRuleFilter{}
		//if len(tags) > 0 {
		//	lifecycleRuleAndOp := &s3.LifecycleRuleAndOperator{}
		//	lifecycleRuleAndOp.SetPrefix(r["prefix"].(string))
		//	//lifecycleRuleAndOp.SetTags(tags)
		//	filter.SetAnd(lifecycleRuleAndOp)
		//} else {
		//	filter.SetPrefix(r["prefix"].(string))
		//}
		filter.SetPrefix(r["prefix"].(string))
		rule.SetFilter(filter)

		// ID
		if val, ok := r["id"].(string); ok && val != "" {
			rule.ID = aws.String(val)
		} else {
			rule.ID = aws.String(resource.PrefixedUniqueId("tf-s3-lifecycle-"))
		}

		// Enabled
		if val, ok := r["enabled"].(bool); ok && val {
			rule.Status = aws.String(s3.ExpirationStatusEnabled)
		} else {
			rule.Status = aws.String(s3.ExpirationStatusDisabled)
		}

		// AbortIncompleteMultipartUpload
		if val, ok := r["abort_incomplete_multipart_upload_days"].(int); ok && val > 0 {
			rule.AbortIncompleteMultipartUpload = &s3.AbortIncompleteMultipartUpload{
				DaysAfterInitiation: aws.Int64(int64(val)),
			}
		}

		// Expiration
		expiration := d.Get(fmt.Sprintf("lifecycle_rule.%d.expiration", i)).(*schema.Set).List()
		if len(expiration) > 0 {
			e := expiration[0].(map[string]interface{})
			i := &s3.LifecycleExpiration{}

			if val, ok := e["date"].(string); ok && val != "" {
				t, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", val))
				if err != nil {
					return fmt.Errorf("Error Parsing AWS S3 Bucket Lifecycle Expiration Date: %s", err.Error())
				}
				i.Date = aws.Time(t)
			} else if val, ok := e["days"].(int); ok && val > 0 {
				i.Days = aws.Int64(int64(val))
			} else if val, ok := e["expired_object_delete_marker"].(bool); ok {
				i.ExpiredObjectDeleteMarker = aws.Bool(val)
			}
			rule.Expiration = i
		}

		// NoncurrentVersionExpiration
		nc_expiration := d.Get(fmt.Sprintf("lifecycle_rule.%d.noncurrent_version_expiration", i)).(*schema.Set).List()
		if len(nc_expiration) > 0 {
			e := nc_expiration[0].(map[string]interface{})

			if val, ok := e["days"].(int); ok && val > 0 {
				rule.NoncurrentVersionExpiration = &s3.NoncurrentVersionExpiration{
					NoncurrentDays: aws.Int64(int64(val)),
				}
			}
		}

		// Transitions
		transitions := d.Get(fmt.Sprintf("lifecycle_rule.%d.transition", i)).(*schema.Set).List()
		if len(transitions) > 0 {
			rule.Transitions = make([]*s3.Transition, 0, len(transitions))
			for _, transition := range transitions {
				transition := transition.(map[string]interface{})
				i := &s3.Transition{}
				if val, ok := transition["date"].(string); ok && val != "" {
					t, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", val))
					if err != nil {
						return fmt.Errorf("Error Parsing AWS S3 Bucket Lifecycle Expiration Date: %s", err.Error())
					}
					i.Date = aws.Time(t)
				} else if val, ok := transition["days"].(int); ok && val >= 0 {
					i.Days = aws.Int64(int64(val))
				}
				if val, ok := transition["storage_class"].(string); ok && val != "" {
					i.StorageClass = aws.String(val)
				}

				rule.Transitions = append(rule.Transitions, i)
			}
		}
		// NoncurrentVersionTransitions
		nc_transitions := d.Get(fmt.Sprintf("lifecycle_rule.%d.noncurrent_version_transition", i)).(*schema.Set).List()
		if len(nc_transitions) > 0 {
			rule.NoncurrentVersionTransitions = make([]*s3.NoncurrentVersionTransition, 0, len(nc_transitions))
			for _, transition := range nc_transitions {
				transition := transition.(map[string]interface{})
				i := &s3.NoncurrentVersionTransition{}
				if val, ok := transition["days"].(int); ok && val >= 0 {
					i.NoncurrentDays = aws.Int64(int64(val))
				}
				if val, ok := transition["storage_class"].(string); ok && val != "" {
					i.StorageClass = aws.String(val)
				}

				rule.NoncurrentVersionTransitions = append(rule.NoncurrentVersionTransitions, i)
			}
		}

		rules = append(rules, rule)
	}

	i := &s3.PutBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucket),
		LifecycleConfiguration: &s3.BucketLifecycleConfiguration{
			Rules: rules,
		},
	}

	_, err := retryOnAwsCode(s3.ErrCodeNoSuchBucket, func() (interface{}, error) {
		return s3conn.PutBucketLifecycleConfiguration(i)
	})
	if err != nil {
		return fmt.Errorf("Error putting S3 lifecycle: %s", err)
	}

	return nil
}

func resourceDigitalOceanBucketImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")

		d.SetId(s[1])
		d.Set("region", s[0])
	}

	return []*schema.ResourceData{d}, nil
}

func bucketDomainName(bucket string, region string) string {
	return fmt.Sprintf("%s.%s.digitaloceanspaces.com", bucket, region)
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

func expirationHash(v interface{}) int {
	var buf bytes.Buffer
	m, ok := v.(map[string]interface{})

	if !ok {
		return 0
	}

	if v, ok := m["date"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["days"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}
	if v, ok := m["expired_object_delete_marker"]; ok {
		buf.WriteString(fmt.Sprintf("%t-", v.(bool)))
	}
	return hashcode.String(buf.String())
}

func transitionHash(v interface{}) int {
	var buf bytes.Buffer
	m, ok := v.(map[string]interface{})

	if !ok {
		return 0
	}

	if v, ok := m["date"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	if v, ok := m["days"]; ok {
		buf.WriteString(fmt.Sprintf("%d-", v.(int)))
	}
	if v, ok := m["storage_class"]; ok {
		buf.WriteString(fmt.Sprintf("%s-", v.(string)))
	}
	return hashcode.String(buf.String())
}

func validateS3BucketLifecycleTimestamp(v interface{}, k string) (ws []string, errors []error) {
	value := v.(string)
	_, err := time.Parse(time.RFC3339, fmt.Sprintf("%sT00:00:00Z", value))
	if err != nil {
		errors = append(errors, fmt.Errorf(
			"%q cannot be parsed as RFC3339 Timestamp Format", value))
	}

	return
}

func validateS3BucketLifecycleTransitionStorageClass() schema.SchemaValidateFunc {
	return validation.StringInSlice([]string{
		s3.TransitionStorageClassStandardIa,
	}, false)
}

func isAWSErr(err error, code string, message string) bool {
	if err, ok := err.(awserr.Error); ok {
		return err.Code() == code && strings.Contains(err.Message(), message)
	}
	return false
}
