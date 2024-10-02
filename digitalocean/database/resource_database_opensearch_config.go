package database

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanDatabaseOpensearchConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDatabaseOpensearchConfigCreate,
		ReadContext:   resourceDigitalOceanDatabaseOpensearchConfigRead,
		UpdateContext: resourceDigitalOceanDatabaseOpensearchConfigUpdate,
		DeleteContext: resourceDigitalOceanDatabaseOpensearchConfigDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDigitalOceanDatabaseOpensearchConfigImport,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"ism_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ism_history_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"ism_history_max_age_hours": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"ism_history_max_docs": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ism_history_rollover_check_period_hours": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"ism_history_rollover_retention_period_days": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"http_max_content_length_bytes": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"http_max_header_size_bytes": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"http_max_initial_line_length_bytes": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1024),
			},
			"indices_query_bool_max_clause_count": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(64),
			},
			"search_max_buckets": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"indices_fielddata_cache_size_percentage": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(3),
			},
			"indices_memory_index_buffer_size_percentage": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(3),
			},
			"indices_memory_min_index_buffer_size_mb": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(3),
			},
			"indices_memory_max_index_buffer_size_mb": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(3),
			},
			"indices_queries_cache_size_percentage": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(3),
			},
			"indices_recovery_max_mb_per_sec": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(40),
			},
			"indices_recovery_max_concurrent_file_chunks": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(2),
			},
			"action_auto_create_index_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"action_destructive_requires_name": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"enable_security_audit": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"thread_pool_search_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"thread_pool_search_throttled_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"thread_pool_search_throttled_queue_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(10),
			},
			"thread_pool_search_queue_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(10),
			},
			"thread_pool_get_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"thread_pool_get_queue_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(10),
			},
			"thread_pool_analyze_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"thread_pool_analyze_queue_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(10),
			},
			"thread_pool_write_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"thread_pool_write_queue_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(10),
			},
			"thread_pool_force_merge_size": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},
			"override_main_response_version": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"script_max_compilations_rate": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_max_shards_per_node": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(100),
			},
			"cluster_routing_allocation_node_concurrent_recoveries": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntAtLeast(2),
			},
			"plugins_alerting_filter_by_backend_roles_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"reindex_remote_whitelist": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceDigitalOceanDatabaseOpensearchConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	if err := updateOpensearchConfig(ctx, d, client); err != nil {
		return diag.Errorf("Error updating Opensearch configuration: %s", err)
	}

	d.SetId(makeDatabaseOpensearchConfigID(clusterID))

	return resourceDigitalOceanDatabaseOpensearchConfigRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseOpensearchConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	if err := updateOpensearchConfig(ctx, d, client); err != nil {
		return diag.Errorf("Error updating Opensearch configuration: %s", err)
	}

	return resourceDigitalOceanDatabaseOpensearchConfigRead(ctx, d, meta)
}

func updateOpensearchConfig(ctx context.Context, d *schema.ResourceData, client *godo.Client) error {
	clusterID := d.Get("cluster_id").(string)

	opts := &godo.OpensearchConfig{}

	if v, ok := d.GetOk("ism_enabled"); ok {
		opts.IsmEnabled = godo.PtrTo(v.(bool))
	}

	if v, ok := d.GetOk("ism_history_enabled"); ok {
		opts.IsmHistoryEnabled = godo.PtrTo(v.(bool))
	}

	if v, ok := d.GetOk("ism_history_max_age_hours"); ok {
		opts.IsmHistoryMaxAgeHours = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("ism_history_max_docs"); ok {
		v, err := strconv.ParseUint(v.(string), 10, 64)
		if err == nil {
			opts.IsmHistoryMaxDocs = godo.PtrTo(v)
		}
	}

	if v, ok := d.GetOk("ism_history_rollover_check_period_hours"); ok {
		opts.IsmHistoryRolloverCheckPeriodHours = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("ism_history_rollover_retention_period_days"); ok {
		opts.IsmHistoryRolloverRetentionPeriodDays = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("http_max_content_length_bytes"); ok {
		opts.HttpMaxContentLengthBytes = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("http_max_header_size_bytes"); ok {
		opts.HttpMaxHeaderSizeBytes = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("http_max_initial_line_length_bytes"); ok {
		opts.HttpMaxInitialLineLengthBytes = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("indices_query_bool_max_clause_count"); ok {
		opts.IndicesQueryBoolMaxClauseCount = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("search_max_buckets"); ok {
		opts.SearchMaxBuckets = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("indices_fielddata_cache_size_percentage"); ok {
		opts.IndicesFielddataCacheSizePercentage = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("indices_memory_index_buffer_size_percentage"); ok {
		opts.IndicesMemoryIndexBufferSizePercentage = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("indices_memory_min_index_buffer_size_mb"); ok {
		opts.IndicesMemoryMinIndexBufferSizeMb = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("indices_memory_max_index_buffer_size_mb"); ok {
		opts.IndicesMemoryMaxIndexBufferSizeMb = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("indices_queries_cache_size_percentage"); ok {
		opts.IndicesQueriesCacheSizePercentage = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("indices_recovery_max_mb_per_sec"); ok {
		opts.IndicesRecoveryMaxMbPerSec = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("indices_recovery_max_concurrent_file_chunks"); ok {
		opts.IndicesRecoveryMaxConcurrentFileChunks = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("action_auto_create_index_enabled"); ok {
		opts.ActionAutoCreateIndexEnabled = godo.PtrTo(v.(bool))
	}

	if v, ok := d.GetOk("action_destructive_requires_name"); ok {
		opts.ActionDestructiveRequiresName = godo.PtrTo(v.(bool))
	}

	if v, ok := d.GetOk("enable_security_audit"); ok {
		opts.EnableSecurityAudit = godo.PtrTo(v.(bool))
	}

	if v, ok := d.GetOk("thread_pool_search_size"); ok {
		opts.ThreadPoolSearchSize = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("thread_pool_search_throttled_size"); ok {
		opts.ThreadPoolSearchThrottledSize = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("thread_pool_search_throttled_queue_size"); ok {
		opts.ThreadPoolSearchThrottledQueueSize = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("thread_pool_search_queue_size"); ok {
		opts.ThreadPoolSearchQueueSize = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("thread_pool_get_size"); ok {
		opts.ThreadPoolGetSize = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("thread_pool_get_queue_size"); ok {
		opts.ThreadPoolGetQueueSize = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("thread_pool_analyze_size"); ok {
		opts.ThreadPoolAnalyzeSize = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("thread_pool_analyze_queue_size"); ok {
		opts.ThreadPoolAnalyzeQueueSize = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("thread_pool_write_size"); ok {
		opts.ThreadPoolWriteSize = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("thread_pool_write_queue_size"); ok {
		opts.ThreadPoolWriteQueueSize = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("thread_pool_force_merge_size"); ok {
		opts.ThreadPoolForceMergeSize = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("override_main_response_version"); ok {
		opts.OverrideMainResponseVersion = godo.PtrTo(v.(bool))
	}

	if v, ok := d.GetOk("script_max_compilations_rate"); ok {
		opts.ScriptMaxCompilationsRate = godo.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("cluster_max_shards_per_node"); ok {
		opts.ClusterMaxShardsPerNode = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("cluster_routing_allocation_node_concurrent_recoveries"); ok {
		opts.ClusterRoutingAllocationNodeConcurrentRecoveries = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("plugins_alerting_filter_by_backend_roles_enabled"); ok {
		opts.PluginsAlertingFilterByBackendRolesEnabled = godo.PtrTo(v.(bool))
	}

	if v, ok := d.GetOk("reindex_remote_whitelist"); ok {
		opts.ReindexRemoteWhitelist = make([]string, 0, len(v.([]interface{})))
	}

	log.Printf("[DEBUG] Opensearch configuration: %s", godo.Stringify(opts))

	if _, err := client.Databases.UpdateOpensearchConfig(ctx, clusterID, opts); err != nil {
		return err
	}

	return nil
}

func resourceDigitalOceanDatabaseOpensearchConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	config, resp, err := client.Databases.GetOpensearchConfig(ctx, clusterID)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving Opensearch configuration: %s", err)
	}

	d.Set("ism_enabled", config.IsmEnabled)
	d.Set("ism_history_enabled", config.IsmHistoryEnabled)
	d.Set("ism_history_max_age_hours", config.IsmHistoryMaxAgeHours)
	d.Set("ism_history_max_docs", config.IsmHistoryMaxDocs)
	d.Set("ism_history_rollover_check_period_hours", config.IsmHistoryRolloverCheckPeriodHours)
	d.Set("ism_history_rollover_retention_period_days", config.IsmHistoryRolloverRetentionPeriodDays)
	d.Set("http_max_content_length_bytes", config.HttpMaxContentLengthBytes)
	d.Set("http_max_header_size_bytes", config.HttpMaxHeaderSizeBytes)
	d.Set("http_max_initial_line_length_bytes", config.HttpMaxInitialLineLengthBytes)
	d.Set("indices_query_bool_max_clause_count", config.IndicesQueryBoolMaxClauseCount)
	d.Set("search_max_buckets", config.SearchMaxBuckets)
	d.Set("indices_fielddata_cache_size_percentage", config.IndicesFielddataCacheSizePercentage)
	d.Set("indices_memory_index_buffer_size_percentage", config.IndicesMemoryIndexBufferSizePercentage)
	d.Set("indices_memory_min_index_buffer_size_mb", config.IndicesMemoryMinIndexBufferSizeMb)
	d.Set("indices_memory_max_index_buffer_size_mb", config.IndicesMemoryMaxIndexBufferSizeMb)
	d.Set("indices_queries_cache_size_percentage", config.IndicesQueriesCacheSizePercentage)
	d.Set("indices_recovery_max_mb_per_sec", config.IndicesRecoveryMaxMbPerSec)
	d.Set("indices_recovery_max_concurrent_file_chunks", config.IndicesRecoveryMaxConcurrentFileChunks)
	d.Set("action_auto_create_index_enabled", config.ActionAutoCreateIndexEnabled)
	d.Set("action_destructive_requires_name", config.ActionDestructiveRequiresName)
	d.Set("enable_security_audit", config.EnableSecurityAudit)
	d.Set("thread_pool_search_size", config.ThreadPoolSearchSize)
	d.Set("thread_pool_search_throttled_size", config.ThreadPoolSearchThrottledSize)
	d.Set("thread_pool_search_throttled_queue_size", config.ThreadPoolSearchThrottledQueueSize)
	d.Set("thread_pool_search_queue_size", config.ThreadPoolSearchQueueSize)
	d.Set("thread_pool_get_size", config.ThreadPoolGetSize)
	d.Set("thread_pool_get_queue_size", config.ThreadPoolGetQueueSize)
	d.Set("thread_pool_analyze_size", config.ThreadPoolAnalyzeSize)
	d.Set("thread_pool_analyze_queue_size", config.ThreadPoolAnalyzeQueueSize)
	d.Set("thread_pool_write_size", config.ThreadPoolWriteSize)
	d.Set("thread_pool_write_queue_size", config.ThreadPoolWriteQueueSize)
	d.Set("thread_pool_force_merge_size", config.ThreadPoolForceMergeSize)
	d.Set("override_main_response_version", config.OverrideMainResponseVersion)
	d.Set("script_max_compilations_rate", config.ScriptMaxCompilationsRate)
	d.Set("cluster_max_shards_per_node", config.ClusterMaxShardsPerNode)
	d.Set("cluster_routing_allocation_node_concurrent_recoveries", config.ClusterRoutingAllocationNodeConcurrentRecoveries)
	d.Set("plugins_alerting_filter_by_backend_roles_enabled", config.PluginsAlertingFilterByBackendRolesEnabled)
	d.Set("reindex_remote_whitelist", config.ReindexRemoteWhitelist)

	return nil
}

func resourceDigitalOceanDatabaseOpensearchConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	warn := []diag.Diagnostic{
		{
			Severity: diag.Warning,
			Summary:  "digitalocean_database_opensearch_config removed from state",
			Detail:   "Database configurations are only removed from state when destroyed. The remote configuration is not unset.",
		},
	}

	return warn
}

func resourceDigitalOceanDatabaseOpensearchConfigImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	clusterID := d.Id()

	d.SetId(makeDatabaseOpensearchConfigID(clusterID))
	d.Set("cluster_id", clusterID)

	return []*schema.ResourceData{d}, nil
}

func makeDatabaseOpensearchConfigID(clusterID string) string {
	return fmt.Sprintf("%s/opensearch-config", clusterID)
}
