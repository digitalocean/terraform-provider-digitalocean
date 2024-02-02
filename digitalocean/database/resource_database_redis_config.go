package database

import (
	"context"
	"fmt"
	"log"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanDatabaseRedisConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDatabaseRedisConfigCreate,
		ReadContext:   resourceDigitalOceanDatabaseRedisConfigRead,
		UpdateContext: resourceDigitalOceanDatabaseRedisConfigUpdate,
		DeleteContext: resourceDigitalOceanDatabaseRedisConfigDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDigitalOceanDatabaseRedisConfigImport,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"maxmemory_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"pubsub_client_output_buffer_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"number_of_databases": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"io_threads": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"lfu_log_factor": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"lfu_decay_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"ssl": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"notify_keyspace_events": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"persistence": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"off",
						"rdb",
					},
					true,
				),
			},

			"acl_channels_default": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"allchannels",
						"resetchannels",
					},
					true,
				),
			},
		},
	}
}

func resourceDigitalOceanDatabaseRedisConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	err := updateRedisConfig(ctx, d, client)
	if err != nil {
		return diag.Errorf("Error updating Redis configuration: %s", err)
	}

	d.SetId(makeDatabaseRedisConfigID(clusterID))

	return resourceDigitalOceanDatabaseRedisConfigRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseRedisConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	err := updateRedisConfig(ctx, d, client)
	if err != nil {
		return diag.Errorf("Error updating Redis configuration: %s", err)
	}

	return resourceDigitalOceanDatabaseRedisConfigRead(ctx, d, meta)
}

func updateRedisConfig(ctx context.Context, d *schema.ResourceData, client *godo.Client) error {
	clusterID := d.Get("cluster_id").(string)

	opts := &godo.RedisConfig{}

	if v, ok := d.GetOk("maxmemory_policy"); ok {
		opts.RedisMaxmemoryPolicy = godo.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("pubsub_client_output_buffer_limit"); ok {
		opts.RedisPubsubClientOutputBufferLimit = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("number_of_databases"); ok {
		opts.RedisNumberOfDatabases = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("io_threads"); ok {
		opts.RedisIOThreads = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("lfu_log_factor"); ok {
		opts.RedisLFULogFactor = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("lfu_decay_time"); ok {
		opts.RedisLFUDecayTime = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOkExists("ssl"); ok {
		opts.RedisSSL = godo.PtrTo(v.(bool))
	}

	if v, ok := d.GetOkExists("timeout"); ok {
		opts.RedisTimeout = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("notify_keyspace_events"); ok {
		opts.RedisNotifyKeyspaceEvents = godo.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("persistence"); ok {
		opts.RedisPersistence = godo.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("acl_channels_default"); ok {
		opts.RedisACLChannelsDefault = godo.PtrTo(v.(string))
	}

	log.Printf("[DEBUG] Redis configuration: %s", godo.Stringify(opts))
	_, err := client.Databases.UpdateRedisConfig(ctx, clusterID, opts)
	if err != nil {
		return err
	}

	return nil
}

func resourceDigitalOceanDatabaseRedisConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	config, resp, err := client.Databases.GetRedisConfig(ctx, clusterID)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving Redis configuration: %s", err)
	}

	d.Set("maxmemory_policy", config.RedisMaxmemoryPolicy)
	d.Set("pubsub_client_output_buffer_limit", config.RedisPubsubClientOutputBufferLimit)
	d.Set("number_of_databases", config.RedisNumberOfDatabases)
	d.Set("io_threads", config.RedisIOThreads)
	d.Set("lfu_log_factor", config.RedisLFULogFactor)
	d.Set("lfu_decay_time", config.RedisLFUDecayTime)
	d.Set("ssl", config.RedisSSL)
	d.Set("timeout", config.RedisTimeout)
	d.Set("notify_keyspace_events", config.RedisNotifyKeyspaceEvents)
	d.Set("persistence", config.RedisPersistence)
	d.Set("acl_channels_default", config.RedisACLChannelsDefault)

	return nil
}

func resourceDigitalOceanDatabaseRedisConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	warn := []diag.Diagnostic{
		{
			Severity: diag.Warning,
			Summary:  "digitalocean_database_redis_config removed from state",
			Detail:   "Database configurations are only removed from state when destroyed. The remote configuration is not unset.",
		},
	}
	return warn
}

func resourceDigitalOceanDatabaseRedisConfigImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	clusterID := d.Id()
	d.SetId(makeDatabaseRedisConfigID(clusterID))
	d.Set("cluster_id", clusterID)

	return []*schema.ResourceData{d}, nil
}

func makeDatabaseRedisConfigID(clusterID string) string {
	return fmt.Sprintf("%s/redis-config", clusterID)
}
