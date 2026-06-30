package database

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const vectorDBActiveStatus = "active"

func ResourceDigitalOceanVectorDatabase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanVectorDatabaseCreate,
		ReadContext:   resourceDigitalOceanVectorDatabaseRead,
		UpdateContext: resourceDigitalOceanVectorDatabaseUpdate,
		DeleteContext: resourceDigitalOceanVectorDatabaseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					// DO API V2 region slug is always lowercase
					return strings.ToLower(val.(string))
				},
				ValidateFunc: validation.NoZeroValues,
			},

			"size": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"tags": tag.TagsSchema(),

			"config": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"default_quantization": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"enable_auto_schema": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"weaviate_version": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			// Computed attributes
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"owner_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"http": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"grpc": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},
	}
}

func resourceDigitalOceanVectorDatabaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	opts := &godo.VectorDBCreateRequest{
		Name:      d.Get("name").(string),
		Region:    strings.ToLower(d.Get("region").(string)),
		Size:      d.Get("size").(string),
		Tags:      tag.ExpandTags(d.Get("tags").(*schema.Set).List()),
		ProjectID: d.Get("project_id").(string),
	}

	vectorDB, _, err := client.VectorDBs.Create(ctx, opts)
	if err != nil {
		return diag.Errorf("Error creating vector database: %s", err)
	}

	d.SetId(vectorDB.ID)

	if _, err := waitForVectorDB(ctx, client, d, vectorDBActiveStatus); err != nil {
		return diag.Errorf("Error waiting for vector database (%s) to become active: %s", d.Id(), err)
	}

	// Configuration can only be set after the cluster is provisioned, so apply
	// any user-supplied config as a follow-up update.
	if _, ok := d.GetOk("config"); ok {
		updateOpts := &godo.VectorDBUpdateRequest{
			ID:     d.Id(),
			Config: expandVectorDBConfig(d.Get("config").([]interface{})),
		}
		if _, _, err := client.VectorDBs.Update(ctx, d.Id(), updateOpts); err != nil {
			return diag.Errorf("Error setting vector database config: %s", err)
		}
	}

	return resourceDigitalOceanVectorDatabaseRead(ctx, d, meta)
}

func resourceDigitalOceanVectorDatabaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	vectorDB, resp, err := client.VectorDBs.Get(ctx, d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error retrieving vector database: %s", err)
	}

	d.Set("name", vectorDB.Name)
	d.Set("region", vectorDB.Region)
	d.Set("size", vectorDB.Size)
	d.Set("status", vectorDB.Status)
	d.Set("owner_uuid", vectorDB.OwnerUUID)
	d.Set("created_at", vectorDB.CreatedAt.UTC().String())
	d.Set("updated_at", vectorDB.UpdatedAt.UTC().String())
	d.Set("tags", tag.FlattenTags(vectorDB.Tags))

	if err := d.Set("config", flattenVectorDBConfig(vectorDB.Config)); err != nil {
		return diag.Errorf("Error setting vector database config: %s", err)
	}

	if err := d.Set("endpoints", flattenVectorDBEndpoints(vectorDB.Endpoints)); err != nil {
		return diag.Errorf("Error setting vector database endpoints: %s", err)
	}

	return nil
}

func resourceDigitalOceanVectorDatabaseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	if d.HasChange("size") {
		opts := &godo.VectorDBResizeRequest{
			ID:   d.Id(),
			Size: d.Get("size").(string),
		}
		if _, _, err := client.VectorDBs.Resize(ctx, d.Id(), opts); err != nil {
			return diag.Errorf("Error resizing vector database: %s", err)
		}

		if _, err := waitForVectorDB(ctx, client, d, vectorDBActiveStatus); err != nil {
			return diag.Errorf("Error waiting for vector database (%s) to become active after resize: %s", d.Id(), err)
		}
	}

	if d.HasChange("tags") {
		opts := &godo.VectorDBUpdateTagsRequest{
			ID:   d.Id(),
			Tags: tag.ExpandTags(d.Get("tags").(*schema.Set).List()),
		}
		if _, _, err := client.VectorDBs.UpdateTags(ctx, d.Id(), opts); err != nil {
			return diag.Errorf("Error updating vector database tags: %s", err)
		}
	}

	if d.HasChange("config") {
		opts := &godo.VectorDBUpdateRequest{
			ID:     d.Id(),
			Config: expandVectorDBConfig(d.Get("config").([]interface{})),
		}
		if _, _, err := client.VectorDBs.Update(ctx, d.Id(), opts); err != nil {
			return diag.Errorf("Error updating vector database config: %s", err)
		}
	}

	return resourceDigitalOceanVectorDatabaseRead(ctx, d, meta)
}

func resourceDigitalOceanVectorDatabaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	if _, err := client.VectorDBs.Delete(ctx, d.Id()); err != nil {
		return diag.Errorf("Error deleting vector database: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForVectorDB(ctx context.Context, client *godo.Client, d *schema.ResourceData, status string) (*godo.VectorDB, error) {
	var (
		tickerInterval = 15 * time.Second
		timeoutSeconds = d.Timeout(schema.TimeoutCreate).Seconds()
		timeout        = int(timeoutSeconds / tickerInterval.Seconds())
		n              = 0
		ticker         = time.NewTicker(tickerInterval)
	)
	defer ticker.Stop()

	for range ticker.C {
		vectorDB, resp, err := client.VectorDBs.Get(ctx, d.Id())
		if resp != nil && resp.StatusCode == 404 {
			continue
		}

		if err != nil {
			return nil, fmt.Errorf("error trying to read vector database state: %s", err)
		}

		if vectorDB.Status == status {
			return vectorDB, nil
		}

		if n >= timeout {
			break
		}

		n++
	}

	return nil, fmt.Errorf("timeout waiting for vector database to become %s", status)
}

func expandVectorDBConfig(config []interface{}) *godo.VectorDBConfig {
	if len(config) == 0 || config[0] == nil {
		return nil
	}

	configMap := config[0].(map[string]interface{})

	return &godo.VectorDBConfig{
		DefaultQuantization: configMap["default_quantization"].(string),
		EnableAutoSchema:    configMap["enable_auto_schema"].(bool),
		WeaviateVersion:     configMap["weaviate_version"].(string),
	}
}

func flattenVectorDBConfig(config *godo.VectorDBConfig) []interface{} {
	if config == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"default_quantization": config.DefaultQuantization,
			"enable_auto_schema":   config.EnableAutoSchema,
			"weaviate_version":     config.WeaviateVersion,
		},
	}
}

func flattenVectorDBEndpoints(endpoints *godo.VectorDBEndpoints) []interface{} {
	if endpoints == nil {
		return []interface{}{}
	}

	return []interface{}{
		map[string]interface{}{
			"http": endpoints.HTTP,
			"grpc": endpoints.GRPC,
		},
	}
}
