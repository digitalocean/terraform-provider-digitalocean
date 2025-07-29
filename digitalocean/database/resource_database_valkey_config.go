package database

import (
	"context"
	"fmt"
	"log"
	"regexp"

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
				Description:  "A unique identifier for the database cluster.",
			},

			"pubsub_client_output_buffer_limit": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "Set output buffer limit for pub / sub clients in MB. The value is the hard limit, the soft limit is 1/4 of the hard limit. When setting the limit, be mindful of the available memory in the selected service plan.",
				ValidateFunc: validation.IntBetween(32, 512),
			},

			"number_of_databases": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "The number of logical databases in the Valkey cluster. Must be between 1 and 128.",
				ValidateFunc: validation.IntBetween(1, 128),
			},

			"io_threads": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "The number of IO threads used by Valkey. Must be between 1 and 32.",
				ValidateFunc: validation.IntBetween(1, 32),
			},

			"lfu_log_factor": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "The log factor for Valkey's LFU (Least Frequently Used) cache eviction. Must be between 1 and 100.",
				ValidateFunc: validation.IntBetween(1, 100),
			},

			"lfu_decay_time": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "The decay time for Valkey's LFU cache eviction. Must be between 1 and 120.",
				ValidateFunc: validation.IntBetween(1, 120),
			},

			"ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Whether to enable SSL/TLS for connections to the Valkey cluster.",
			},

			"timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: "The timeout (in seconds) for Valkey client connections.",
			},

			"notify_keyspace_events": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Set notify-keyspace-events option. Requires at least K or E and accepts any combination of the following options. Setting the parameter to \"\" disables notifications.\n\nK — Keyspace events\nE — Keyevent events\ng — Generic commands (e.g. DEL, EXPIRE, RENAME, ...)\n$ — String commands\nl — List commands\ns — Set commands\nh — Hash commands\nz — Sorted set commands\nt — Stream commands\nd — Module key type events\nx — Expired events\ne — Evicted events\nm — Key miss events\nn — New key events\nA — Alias for \"g$lshztxed\"",
				ValidateFunc: validation.All(
					validation.StringLenBetween(0, 32),
					validation.StringMatch(regexp.MustCompile(`^[KEg$lshzxeA]*$`), "must only contain: K, E, g, $, l, s, h, z, x, e, A"),
				),
			},

			"persistence": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "When persistence is 'rdb', Valkey does RDB dumps each 10 minutes if any key is changed. Also RDB dumps are done according to backup schedule for backup purposes. When persistence is 'off', no RDB dumps and backups are done, so data can be lost at any moment if service is restarted for any reason, or if service is powered off. Also service can't be forked.",
				ValidateFunc: validation.StringInSlice(
					[]string{
						"off",
						"rdb",
					},
					true,
				),
			},

			"acl_channels_default": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Determines default pub/sub channels' ACL for new users if ACL is not supplied. When this option is not defined, all_channels is assumed to keep backward compatibility. This option doesn't affect Valkey configuration acl-pubsub-default.",
				ValidateFunc: validation.StringInSlice(
					[]string{
						"allchannels",
						"resetchannels",
					},
					true,
				),
			},

			"frequent_snapshots": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: "Frequent RDB snapshots. When enabled, Valkey will create frequent local RDB snapshots. When disabled, Valkey will only take RDB snapshots when a backup is created, based on the backup schedule. This setting is ignored when valkey_persistence is set to off.",
			},

			"valkey_active_expire_effort": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				Description:  "Active expire effort. Valkey reclaims expired keys both when accessed and in the background. The background process scans for expired keys to free memory. Increasing the active-expire-effort setting (default 1, max 10) uses more CPU to reclaim expired keys faster, reducing memory usage but potentially increasing latency.",
				ValidateFunc: validation.IntBetween(1, 10),
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

	// Check if the cluster is active before updating config
	cluster, _, err := client.Databases.Get(ctx, clusterID)
	if err != nil {
		return fmt.Errorf("failed to fetch cluster status: %w", err)
	}
	if cluster.Status != "online" {
		return fmt.Errorf("cannot update config: cluster status is '%s' (must be 'online')", cluster.Status)
	}

	opts := &godo.ValkeyConfig{}

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
	_, err = client.Databases.UpdateValkeyConfig(ctx, clusterID, opts)
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
