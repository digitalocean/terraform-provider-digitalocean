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

func ResourceDigitalOceanDatabaseValkeyConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDatabaseValkeyConfigCreate,
		ReadContext:   resourceDigitalOceanDatabaseValkeyConfigRead,
		UpdateContext: resourceDigitalOceanDatabaseValkeyConfigUpdate,
		DeleteContext: resourceDigitalOceanDatabaseValkeyConfigDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDigitalOceanDatabaseValkeyConfigImport,
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

func resourceDigitalOceanDatabaseValkeyConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	err := updateValkeyConfig(ctx, d, client)
	if err != nil {
		return diag.Errorf("Error updating Valkey configuration: %s", err)
	}

	d.SetId(makeDatabaseValkeyConfigID(clusterID))

	return resourceDigitalOceanDatabaseValkeyConfigRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseValkeyConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	err := updateValkeyConfig(ctx, d, client)
	if err != nil {
		return diag.Errorf("Error updating Valkey configuration: %s", err)
	}

	return resourceDigitalOceanDatabaseValkeyConfigRead(ctx, d, meta)
}

func updateValkeyConfig(ctx context.Context, d *schema.ResourceData, client *godo.Client) error {
	clusterID := d.Get("cluster_id").(string)

	opts := &godo.ValkeyConfig{}

	if v, ok := d.GetOk("maxmemory_policy"); ok {
		opts.ValkeyMaxmemoryPolicy = godo.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("pubsub_client_output_buffer_limit"); ok {
		opts.ValkeyPubsubClientOutputBufferLimit = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("number_of_databases"); ok {
		opts.ValkeyNumberOfDatabases = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("io_threads"); ok {
		opts.ValkeyIOThreads = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("lfu_log_factor"); ok {
		opts.ValkeyLFULogFactor = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("lfu_decay_time"); ok {
		opts.ValkeyLFUDecayTime = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOkExists("ssl"); ok {
		opts.ValkeySSL = godo.PtrTo(v.(bool))
	}

	if v, ok := d.GetOkExists("timeout"); ok {
		opts.ValkeyTimeout = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("notify_keyspace_events"); ok {
		opts.ValkeyNotifyKeyspaceEvents = godo.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("persistence"); ok {
		opts.ValkeyPersistence = godo.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("acl_channels_default"); ok {
		opts.ValkeyACLChannelsDefault = godo.PtrTo(v.(string))
	}

	log.Printf("[DEBUG] Valkey configuration: %s", godo.Stringify(opts))
	_, err := client.Databases.UpdateValkeyConfig(ctx, clusterID, opts)
	if err != nil {
		return err
	}

	return nil
}

func resourceDigitalOceanDatabaseValkeyConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	config, resp, err := client.Databases.GetValkeyConfig(ctx, clusterID)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving Valkey configuration: %s", err)
	}

	d.Set("maxmemory_policy", config.ValkeyMaxmemoryPolicy)
	d.Set("pubsub_client_output_buffer_limit", config.ValkeyPubsubClientOutputBufferLimit)
	d.Set("number_of_databases", config.ValkeyNumberOfDatabases)
	d.Set("io_threads", config.ValkeyIOThreads)
	d.Set("lfu_log_factor", config.ValkeyLFULogFactor)
	d.Set("lfu_decay_time", config.ValkeyLFUDecayTime)
	d.Set("ssl", config.ValkeySSL)
	d.Set("timeout", config.ValkeyTimeout)
	d.Set("notify_keyspace_events", config.ValkeyNotifyKeyspaceEvents)
	d.Set("persistence", config.ValkeyPersistence)
	d.Set("acl_channels_default", config.ValkeyACLChannelsDefault)

	return nil
}

func resourceDigitalOceanDatabaseValkeyConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")
	warn := []diag.Diagnostic{
		{
			Severity: diag.Warning,
			Summary:  "digitalocean_database_valkey_config removed from state",
			Detail:   "Database configurations are only removed from state when destroyed. The remote configuration is not unset.",
		},
	}
	return warn
}

func resourceDigitalOceanDatabaseValkeyConfigImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	clusterID := d.Id()
	d.SetId(makeDatabaseValkeyConfigID(clusterID))
	d.Set("cluster_id", clusterID)

	return []*schema.ResourceData{d}, nil
}

func makeDatabaseValkeyConfigID(clusterID string) string {
	return fmt.Sprintf("%s/valkey-config", clusterID)
}
