package database

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strconv"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanDatabaseKafkaConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDatabaseKafkaConfigCreate,
		ReadContext:   resourceDigitalOceanDatabaseKafkaConfigRead,
		UpdateContext: resourceDigitalOceanDatabaseKafkaConfigUpdate,
		DeleteContext: resourceDigitalOceanDatabaseKafkaConfigDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDigitalOceanDatabaseKafkaConfigImport,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"group_initial_rebalance_delay_ms": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"group_min_session_timeout_ms": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"group_max_session_timeout_ms": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"message_max_bytes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"log_cleaner_delete_retention_ms": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"log_cleaner_min_compaction_lag_ms": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"log_flush_interval_ms": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"log_index_interval_bytes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"log_message_downconversion_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"log_message_timestamp_difference_max_ms": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"log_preallocate": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"log_retention_bytes": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"log_retention_hours": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"log_retention_ms": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"log_roll_jitter_ms": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"log_segment_delete_delay_ms": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"auto_create_topics_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceDigitalOceanDatabaseKafkaConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	if err := updateKafkaConfig(ctx, d, client); err != nil {
		return diag.Errorf("Error updating Kafka configuration: %s", err)
	}

	d.SetId(makeDatabaseKafkaConfigID(clusterID))

	return resourceDigitalOceanDatabaseKafkaConfigRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseKafkaConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	if err := updateKafkaConfig(ctx, d, client); err != nil {
		return diag.Errorf("Error updating Kafka configuration: %s", err)
	}

	return resourceDigitalOceanDatabaseKafkaConfigRead(ctx, d, meta)
}

func updateKafkaConfig(ctx context.Context, d *schema.ResourceData, client *godo.Client) error {
	clusterID := d.Get("cluster_id").(string)

	opts := &godo.KafkaConfig{}

	if v, ok := d.GetOk("group_initial_rebalance_delay_ms"); ok {
		opts.GroupInitialRebalanceDelayMs = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("group_min_session_timeout_ms"); ok {
		opts.GroupMinSessionTimeoutMs = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("group_max_session_timeout_ms"); ok {
		opts.GroupMaxSessionTimeoutMs = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("message_max_bytes"); ok {
		opts.MessageMaxBytes = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("log_cleaner_delete_retention_ms"); ok {
		opts.LogCleanerDeleteRetentionMs = godo.PtrTo(int64(v.(int)))
	}

	if v, ok := d.GetOk("log_cleaner_min_compaction_lag_ms"); ok {
		v, err := strconv.ParseUint(v.(string), 10, 64)
		if err == nil {
			opts.LogCleanerMinCompactionLagMs = godo.PtrTo(v)
		}
	}

	if v, ok := d.GetOk("log_flush_interval_ms"); ok {
		v, err := strconv.ParseUint(v.(string), 10, 64)
		if err == nil {
			opts.LogFlushIntervalMs = godo.PtrTo(v)
		}
	}

	if v, ok := d.GetOk("log_index_interval_bytes"); ok {
		opts.LogIndexIntervalBytes = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("log_message_downconversion_enable"); ok {
		opts.LogMessageDownconversionEnable = godo.PtrTo(v.(bool))
	}

	if v, ok := d.GetOk("log_message_timestamp_difference_max_ms"); ok {
		v, err := strconv.ParseUint(v.(string), 10, 64)
		if err == nil {
			opts.LogMessageTimestampDifferenceMaxMs = godo.PtrTo(v)
		}
	}

	if v, ok := d.GetOk("log_preallocate"); ok {
		opts.LogPreallocate = godo.PtrTo(v.(bool))
	}

	if v, ok := d.GetOk("log_retention_bytes"); ok {
		if v, ok := new(big.Int).SetString(v.(string), 10); ok {
			opts.LogRetentionBytes = v
		}
	}

	if v, ok := d.GetOk("log_retention_hours"); ok {
		opts.LogRetentionHours = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("log_retention_ms"); ok {
		if v, ok := new(big.Int).SetString(v.(string), 10); ok {
			opts.LogRetentionMs = v
		}
	}

	if v, ok := d.GetOk("log_roll_jitter_ms"); ok {
		v, err := strconv.ParseUint(v.(string), 10, 64)
		if err == nil {
			opts.LogRollJitterMs = godo.PtrTo(v)
		}
	}

	if v, ok := d.GetOk("log_segment_delete_delay_ms"); ok {
		opts.LogSegmentDeleteDelayMs = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("auto_create_topics_enable"); ok {
		opts.AutoCreateTopicsEnable = godo.PtrTo(v.(bool))
	}

	log.Printf("[DEBUG] Kafka configuration: %s", godo.Stringify(opts))

	if _, err := client.Databases.UpdateKafkaConfig(ctx, clusterID, opts); err != nil {
		return err
	}

	return nil
}

func resourceDigitalOceanDatabaseKafkaConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	config, resp, err := client.Databases.GetKafkaConfig(ctx, clusterID)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving Kafka configuration: %s", err)
	}

	d.Set("group_initial_rebalance_delay_ms", config.GroupInitialRebalanceDelayMs)
	d.Set("group_min_session_timeout_ms", config.GroupMinSessionTimeoutMs)
	d.Set("group_max_session_timeout_ms", config.GroupMaxSessionTimeoutMs)
	d.Set("message_max_bytes", config.MessageMaxBytes)
	d.Set("log_cleaner_delete_retention_ms", config.LogCleanerDeleteRetentionMs)
	d.Set("log_cleaner_min_compaction_lag_ms", config.LogCleanerMinCompactionLagMs)
	d.Set("log_flush_interval_ms", config.LogFlushIntervalMs)
	d.Set("log_index_interval_bytes", config.LogIndexIntervalBytes)
	d.Set("log_message_downconversion_enable", config.LogMessageDownconversionEnable)
	d.Set("log_message_timestamp_difference_max_ms", config.LogMessageTimestampDifferenceMaxMs)
	d.Set("log_preallocate", config.LogPreallocate)
	d.Set("log_retention_bytes", config.LogRetentionBytes)
	d.Set("log_retention_hours", config.LogRetentionHours)
	d.Set("log_retention_ms", config.LogRetentionMs)
	d.Set("log_roll_jitter_ms", config.LogRollJitterMs)
	d.Set("log_segment_delete_delay_ms", config.LogSegmentDeleteDelayMs)
	d.Set("auto_create_topics_enable", config.AutoCreateTopicsEnable)

	return nil
}

func resourceDigitalOceanDatabaseKafkaConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	warn := []diag.Diagnostic{
		{
			Severity: diag.Warning,
			Summary:  "digitalocean_database_kafka_config removed from state",
			Detail:   "Database configurations are only removed from state when destroyed. The remote configuration is not unset.",
		},
	}

	return warn
}

func resourceDigitalOceanDatabaseKafkaConfigImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	clusterID := d.Id()

	d.SetId(makeDatabaseKafkaConfigID(clusterID))
	d.Set("cluster_id", clusterID)

	return []*schema.ResourceData{d}, nil
}

func makeDatabaseKafkaConfigID(clusterID string) string {
	return fmt.Sprintf("%s/kafka-config", clusterID)
}
