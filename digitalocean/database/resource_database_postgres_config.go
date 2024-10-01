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

func ResourceDigitalOceanDatabasePostgreSQLConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDatabasePostgreSQLConfigCreate,
		ReadContext:   resourceDigitalOceanDatabasePostgreSQLConfigRead,
		UpdateContext: resourceDigitalOceanDatabasePostgreSQLConfigUpdate,
		DeleteContext: resourceDigitalOceanDatabasePostgreSQLConfigDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDigitalOceanDatabasePostgreSQLConfigImport,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"autovacuum_freeze_max_age": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"autovacuum_max_workers": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"autovacuum_naptime": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"autovacuum_vacuum_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"autovacuum_analyze_threshold": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"autovacuum_vacuum_scale_factor": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"autovacuum_analyze_scale_factor": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"autovacuum_vacuum_cost_delay": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"autovacuum_vacuum_cost_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"bgwriter_delay": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"bgwriter_flush_after": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"bgwriter_lru_maxpages": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"bgwriter_lru_multiplier": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"deadlock_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"default_toast_compression": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"lz4",
						"pglz",
					},
					false,
				),
			},
			"idle_in_transaction_session_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"jit": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"log_autovacuum_min_duration": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"log_error_verbosity": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"TERSE",
						"DEFAULT",
						"VERBOSE",
					},
					false,
				),
			},
			"log_line_prefix": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"pid=%p,user=%u,db=%d,app=%a,client=%h",
						"%m [%p] %q[user=%u,db=%d,app=%a]",
						"%t [%p]: [%l-1] user=%u,db=%d,app=%a,client=%h",
					},
					false,
				),
			},
			"log_min_duration_statement": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"max_files_per_process": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"max_prepared_transactions": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"max_pred_locks_per_transaction": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"max_locks_per_transaction": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"max_stack_depth": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"max_standby_archive_delay": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"max_standby_streaming_delay": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"max_replication_slots": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"max_logical_replication_workers": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"max_parallel_workers": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"max_parallel_workers_per_gather": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"max_worker_processes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"pg_partman_bgw_role": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"pg_partman_bgw_interval": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"pg_stat_statements_track": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"temp_file_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"track_activity_query_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"track_commit_timestamp": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"off",
						"on",
					},
					false,
				),
			},
			"track_functions": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"all",
						"pl",
						"none",
					},
					false,
				),
			},
			"track_io_timing": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"off",
						"on",
					},
					false,
				),
			},
			"max_wal_senders": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"wal_sender_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"wal_writer_delay": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"shared_buffers_percentage": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.FloatBetween(
					20.0,
					60.0,
				),
			},
			"pgbouncer": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_reset_query_always": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"ignore_startup_parameters": {
							Type:     schema.TypeSet,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Optional: true,
							Computed: true,
						},
						"min_pool_size": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"server_lifetime": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"server_idle_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"autodb_pool_size": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"autodb_pool_mode": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"autodb_max_db_connections": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"autodb_idle_timeout": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"backup_hour": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"backup_minute": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"work_mem": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"timescaledb": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MinItems: 1,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"max_background_workers": {
							Type:     schema.TypeInt,
							Optional: true,
						}},
				},
			},
		},
	}
}

func resourceDigitalOceanDatabasePostgreSQLConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	if err := updatePostgreSQLConfig(ctx, d, client); err != nil {
		return diag.Errorf("Error updating PostgreSQL configuration: %s", err)
	}

	d.SetId(makeDatabasePostgreSQLConfigID(clusterID))

	return resourceDigitalOceanDatabasePostgreSQLConfigRead(ctx, d, meta)
}

func resourceDigitalOceanDatabasePostgreSQLConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	if err := updatePostgreSQLConfig(ctx, d, client); err != nil {
		return diag.Errorf("Error updating PostgreSQL configuration: %s", err)
	}

	return resourceDigitalOceanDatabasePostgreSQLConfigRead(ctx, d, meta)
}

func updatePostgreSQLConfig(ctx context.Context, d *schema.ResourceData, client *godo.Client) error {
	clusterID := d.Get("cluster_id").(string)

	opts := &godo.PostgreSQLConfig{}

	if v, ok := d.GetOk("autovacuum_freeze_max_age"); ok {
		opts.AutovacuumFreezeMaxAge = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("autovacuum_max_workers"); ok {
		opts.AutovacuumMaxWorkers = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("autovacuum_naptime"); ok {
		opts.AutovacuumNaptime = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("autovacuum_vacuum_threshold"); ok {
		opts.AutovacuumVacuumThreshold = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("autovacuum_analyze_threshold"); ok {
		opts.AutovacuumAnalyzeThreshold = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("autovacuum_vacuum_scale_factor"); ok {
		opts.AutovacuumVacuumScaleFactor = godo.PtrTo(float32(v.(float64)))
	}

	if v, ok := d.GetOk("autovacuum_analyze_scale_factor"); ok {
		opts.AutovacuumAnalyzeScaleFactor = godo.PtrTo(float32(v.(float64)))
	}

	if v, ok := d.GetOk("autovacuum_vacuum_cost_delay"); ok {
		opts.AutovacuumVacuumCostDelay = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("autovacuum_vacuum_cost_limit"); ok {
		opts.AutovacuumVacuumCostLimit = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("bgwriter_delay"); ok {
		opts.BGWriterDelay = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("bgwriter_flush_after"); ok {
		opts.BGWriterFlushAfter = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("bgwriter_lru_maxpages"); ok {
		opts.BGWriterLRUMaxpages = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("bgwriter_lru_multiplier"); ok {
		opts.BGWriterLRUMultiplier = godo.PtrTo(float32(v.(float64)))
	}

	if v, ok := d.GetOk("deadlock_timeout"); ok {
		opts.DeadlockTimeoutMillis = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("default_toast_compression"); ok {
		opts.DefaultToastCompression = godo.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("idle_in_transaction_session_timeout"); ok {
		opts.IdleInTransactionSessionTimeout = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("jit"); ok {
		opts.JIT = godo.PtrTo(v.(bool))
	}

	if v, ok := d.GetOk("log_autovacuum_min_duration"); ok {
		opts.LogAutovacuumMinDuration = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("log_error_verbosity"); ok {
		opts.LogErrorVerbosity = godo.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("log_line_prefix"); ok {
		opts.LogLinePrefix = godo.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("log_min_duration_statement"); ok {
		opts.LogMinDurationStatement = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("max_files_per_process"); ok {
		opts.MaxFilesPerProcess = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("max_prepared_transactions"); ok {
		opts.MaxPreparedTransactions = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("max_pred_locks_per_transaction"); ok {
		opts.MaxPredLocksPerTransaction = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("max_locks_per_transaction"); ok {
		opts.MaxLocksPerTransaction = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("max_stack_depth"); ok {
		opts.MaxStackDepth = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("max_standby_archive_delay"); ok {
		opts.MaxStandbyArchiveDelay = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("max_standby_streaming_delay"); ok {
		opts.MaxStandbyStreamingDelay = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("max_replication_slots"); ok {
		opts.MaxReplicationSlots = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("max_logical_replication_workers"); ok {
		opts.MaxLogicalReplicationWorkers = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("max_parallel_workers"); ok {
		opts.MaxParallelWorkers = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("max_parallel_workers_per_gather"); ok {
		opts.MaxParallelWorkersPerGather = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("max_worker_processes"); ok {
		opts.MaxWorkerProcesses = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("pg_partman_bgw_role"); ok {
		opts.PGPartmanBGWRole = godo.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("pg_partman_bgw_interval"); ok {
		opts.PGPartmanBGWInterval = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("pg_stat_statements_track"); ok {
		opts.PGStatStatementsTrack = godo.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("temp_file_limit"); ok {
		opts.TempFileLimit = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("timezone"); ok {
		opts.Timezone = godo.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("track_activity_query_size"); ok {
		opts.TrackActivityQuerySize = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("track_commit_timestamp"); ok {
		opts.TrackCommitTimestamp = godo.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("track_functions"); ok {
		opts.TrackFunctions = godo.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("track_io_timing"); ok {
		opts.TrackIOTiming = godo.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("max_wal_senders"); ok {
		opts.MaxWalSenders = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("wal_sender_timeout"); ok {
		opts.WalSenderTimeout = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("wal_writer_delay"); ok {
		opts.WalWriterDelay = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("shared_buffers_percentage"); ok {
		opts.SharedBuffersPercentage = godo.PtrTo(float32(v.(float64)))
	}

	if v, ok := d.GetOk("pgbouncer"); ok {
		opts.PgBouncer = expandPgBouncer(v.([]interface{}))
	}

	if v, ok := d.GetOk("backup_hour"); ok {
		opts.BackupHour = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("backup_minute"); ok {
		opts.BackupMinute = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("work_mem"); ok {
		opts.WorkMem = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("timescaledb"); ok {
		opts.TimeScaleDB = expandTimeScaleDB(v.([]interface{}))
	}

	log.Printf("[DEBUG] PostgreSQL configuration: %s", godo.Stringify(opts))

	if _, err := client.Databases.UpdatePostgreSQLConfig(ctx, clusterID, opts); err != nil {
		return err
	}

	return nil
}

func resourceDigitalOceanDatabasePostgreSQLConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	config, resp, err := client.Databases.GetPostgreSQLConfig(ctx, clusterID)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving PostgreSQL configuration: %s", err)
	}

	d.Set("autovacuum_freeze_max_age", config.AutovacuumFreezeMaxAge)
	d.Set("autovacuum_max_workers", config.AutovacuumMaxWorkers)
	d.Set("autovacuum_naptime", config.AutovacuumNaptime)
	d.Set("autovacuum_vacuum_threshold", config.AutovacuumVacuumThreshold)
	d.Set("autovacuum_analyze_threshold", config.AutovacuumAnalyzeThreshold)
	d.Set("autovacuum_vacuum_scale_factor", config.AutovacuumVacuumScaleFactor)
	d.Set("autovacuum_analyze_scale_factor", config.AutovacuumAnalyzeScaleFactor)
	d.Set("autovacuum_vacuum_cost_delay", config.AutovacuumVacuumCostDelay)
	d.Set("autovacuum_vacuum_cost_limit", config.AutovacuumVacuumCostLimit)
	d.Set("bgwriter_delay", config.BGWriterDelay)
	d.Set("bgwriter_flush_after", config.BGWriterFlushAfter)
	d.Set("bgwriter_lru_maxpages", config.BGWriterLRUMaxpages)
	d.Set("bgwriter_lru_multiplier", config.BGWriterLRUMultiplier)
	d.Set("deadlock_timeout", config.DeadlockTimeoutMillis)
	d.Set("default_toast_compression", config.DefaultToastCompression)
	d.Set("idle_in_transaction_session_timeout", config.IdleInTransactionSessionTimeout)
	d.Set("jit", config.JIT)
	d.Set("log_autovacuum_min_duration", config.LogAutovacuumMinDuration)
	d.Set("log_error_verbosity", config.LogErrorVerbosity)
	d.Set("log_line_prefix", config.LogLinePrefix)
	d.Set("log_min_duration_statement", config.LogMinDurationStatement)
	d.Set("max_files_per_process", config.MaxFilesPerProcess)
	d.Set("max_prepared_transactions", config.MaxPreparedTransactions)
	d.Set("max_pred_locks_per_transaction", config.MaxPredLocksPerTransaction)
	d.Set("max_locks_per_transaction", config.MaxLocksPerTransaction)
	d.Set("max_stack_depth", config.MaxStackDepth)
	d.Set("max_standby_archive_delay", config.MaxStandbyArchiveDelay)
	d.Set("max_standby_streaming_delay", config.MaxStandbyStreamingDelay)
	d.Set("max_replication_slots", config.MaxReplicationSlots)
	d.Set("max_logical_replication_workers", config.MaxLogicalReplicationWorkers)
	d.Set("max_parallel_workers", config.MaxParallelWorkers)
	d.Set("max_parallel_workers_per_gather", config.MaxParallelWorkersPerGather)
	d.Set("max_worker_processes", config.MaxWorkerProcesses)
	d.Set("pg_partman_bgw_role", config.PGPartmanBGWRole)
	d.Set("pg_partman_bgw_interval", config.PGPartmanBGWInterval)
	d.Set("pg_stat_statements_track", config.PGStatStatementsTrack)
	d.Set("temp_file_limit", config.TempFileLimit)
	d.Set("timezone", config.Timezone)
	d.Set("track_activity_query_size", config.TrackActivityQuerySize)
	d.Set("track_commit_timestamp", config.TrackCommitTimestamp)
	d.Set("track_functions", config.TrackFunctions)
	d.Set("track_io_timing", config.TrackIOTiming)
	d.Set("max_wal_senders", config.MaxWalSenders)
	d.Set("wal_sender_timeout", config.WalSenderTimeout)
	d.Set("wal_writer_delay", config.WalWriterDelay)
	d.Set("shared_buffers_percentage", config.SharedBuffersPercentage)
	d.Set("backup_hour", config.BackupHour)
	d.Set("backup_minute", config.BackupMinute)
	d.Set("work_mem", config.WorkMem)

	if _, ok := d.GetOk("pgbouncer"); ok {
		if err := d.Set("pgbouncer", flattenPGBouncerOpts(*config.PgBouncer)); err != nil {
			return diag.Errorf("[DEBUG] Error setting pgbouncer - error: %#v", err)
		}
	}

	if _, ok := d.GetOk("timescaledb"); ok {
		if err := d.Set("timescaledb", flattenTimeScaleDBOpts(*config.TimeScaleDB)); err != nil {
			return diag.Errorf("[DEBUG] Error setting timescaledb - error: %#v", err)
		}
	}

	return nil
}

func resourceDigitalOceanDatabasePostgreSQLConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	warn := []diag.Diagnostic{
		{
			Severity: diag.Warning,
			Summary:  "digitalocean_database_postgresql_config removed from state",
			Detail:   "Database configurations are only removed from state when destroyed. The remote configuration is not unset.",
		},
	}

	return warn
}

func resourceDigitalOceanDatabasePostgreSQLConfigImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	clusterID := d.Id()

	d.SetId(makeDatabasePostgreSQLConfigID(clusterID))
	d.Set("cluster_id", clusterID)

	return []*schema.ResourceData{d}, nil
}

func makeDatabasePostgreSQLConfigID(clusterID string) string {
	return fmt.Sprintf("%s/postgresql-config", clusterID)
}

func expandPgBouncer(config []interface{}) *godo.PostgreSQLBouncerConfig {
	configMap := config[0].(map[string]interface{})

	pgBouncerConfig := &godo.PostgreSQLBouncerConfig{
		ServerResetQueryAlways:  godo.PtrTo(configMap["server_reset_query_always"].(bool)),
		IgnoreStartupParameters: godo.PtrTo(configMap["ignore_startup_parameters"].([]string)),
		MinPoolSize:             godo.PtrTo(configMap["min_pool_size"].(int)),
		ServerLifetime:          godo.PtrTo(configMap["server_lifetime"].(int)),
		ServerIdleTimeout:       godo.PtrTo(configMap["server_idle_timeout"].(int)),
		AutodbPoolSize:          godo.PtrTo(configMap["autodb_pool_size"].(int)),
		AutodbPoolMode:          godo.PtrTo(configMap["autodb_pool_mode"].(string)),
		AutodbMaxDbConnections:  godo.PtrTo(configMap["autodb_max_db_connections"].(int)),
		AutodbIdleTimeout:       godo.PtrTo(configMap["autodb_idle_timeout"].(int)),
	}

	return pgBouncerConfig
}

func expandTimeScaleDB(config []interface{}) *godo.PostgreSQLTimeScaleDBConfig {
	configMap := config[0].(map[string]interface{})

	timeScaleDBConfig := &godo.PostgreSQLTimeScaleDBConfig{
		MaxBackgroundWorkers: godo.PtrTo(configMap["max_background_workers"].(int)),
	}

	return timeScaleDBConfig
}

func flattenPGBouncerOpts(opts godo.PostgreSQLBouncerConfig) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	item := make(map[string]interface{})

	item["server_reset_query_always"] = opts.ServerResetQueryAlways
	item["ignore_startup_parameters"] = opts.IgnoreStartupParameters
	item["min_pool_size"] = opts.MinPoolSize
	item["server_lifetime"] = opts.ServerLifetime
	item["server_idle_timeout"] = opts.ServerIdleTimeout
	item["autodb_pool_size"] = opts.AutodbPoolSize
	item["autodb_pool_mode"] = opts.AutodbPoolMode
	item["autodb_max_db_connections"] = opts.AutodbMaxDbConnections
	item["autodb_idle_timeout"] = opts.AutodbIdleTimeout

	result = append(result, item)

	return result
}

func flattenTimeScaleDBOpts(opts godo.PostgreSQLTimeScaleDBConfig) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	item := make(map[string]interface{})

	item["max_background_workers"] = opts.MaxBackgroundWorkers

	result = append(result, item)

	return result
}
