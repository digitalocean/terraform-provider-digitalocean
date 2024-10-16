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

func ResourceDigitalOceanDatabaseMySQLConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDatabaseMySQLConfigCreate,
		ReadContext:   resourceDigitalOceanDatabaseMySQLConfigRead,
		UpdateContext: resourceDigitalOceanDatabaseMySQLConfigUpdate,
		DeleteContext: resourceDigitalOceanDatabaseMySQLConfigDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDigitalOceanDatabaseMySQLConfigImport,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"connect_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"default_time_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"innodb_log_buffer_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"innodb_online_alter_log_max_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"innodb_lock_wait_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"interactive_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"max_allowed_packet": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"net_read_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"sort_buffer_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"sql_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sql_require_primary_key": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"wait_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"net_write_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"group_concat_max_len": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"information_schema_stats_expiry": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"innodb_ft_min_token_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"innodb_ft_server_stopword_table": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"innodb_print_all_deadlocks": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"innodb_rollback_on_timeout": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"internal_tmp_mem_storage_engine": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"TempTable",
						"MEMORY",
					},
					false,
				),
			},
			"max_heap_table_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"tmp_table_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"slow_query_log": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"long_query_time": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
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
			"binlog_retention_period": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceDigitalOceanDatabaseMySQLConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	if err := updateMySQLConfig(ctx, d, client); err != nil {
		return diag.Errorf("Error updating MySQL configuration: %s", err)
	}

	d.SetId(makeDatabaseMySQLConfigID(clusterID))

	return resourceDigitalOceanDatabaseMySQLConfigRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseMySQLConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	if err := updateMySQLConfig(ctx, d, client); err != nil {
		return diag.Errorf("Error updating MySQL configuration: %s", err)
	}

	return resourceDigitalOceanDatabaseMySQLConfigRead(ctx, d, meta)
}

func updateMySQLConfig(ctx context.Context, d *schema.ResourceData, client *godo.Client) error {
	clusterID := d.Get("cluster_id").(string)

	opts := &godo.MySQLConfig{}

	if v, ok := d.GetOk("connect_timeout"); ok {
		opts.ConnectTimeout = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("default_time_zone"); ok {
		opts.DefaultTimeZone = godo.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("innodb_log_buffer_size"); ok {
		opts.InnodbLogBufferSize = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("innodb_online_alter_log_max_size"); ok {
		opts.InnodbOnlineAlterLogMaxSize = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("innodb_lock_wait_timeout"); ok {
		opts.InnodbLockWaitTimeout = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("interactive_timeout"); ok {
		opts.InteractiveTimeout = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("max_allowed_packet"); ok {
		opts.MaxAllowedPacket = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("net_read_timeout"); ok {
		opts.NetReadTimeout = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("sort_buffer_size"); ok {
		opts.SortBufferSize = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("sql_mode"); ok {
		opts.SQLMode = godo.PtrTo(v.(string))
	}

	if v, ok := d.GetOkExists("sql_require_primary_key"); ok {
		opts.SQLRequirePrimaryKey = godo.PtrTo(v.(bool))
	}

	if v, ok := d.GetOk("wait_timeout"); ok {
		opts.WaitTimeout = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("net_write_timeout"); ok {
		opts.NetWriteTimeout = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("group_concat_max_len"); ok {
		opts.GroupConcatMaxLen = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("information_schema_stats_expiry"); ok {
		opts.InformationSchemaStatsExpiry = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("innodb_ft_min_token_size"); ok {
		opts.InnodbFtMinTokenSize = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("innodb_ft_server_stopword_table"); ok {
		opts.InnodbFtServerStopwordTable = godo.PtrTo(v.(string))
	}

	if v, ok := d.GetOkExists("innodb_print_all_deadlocks"); ok {
		opts.InnodbPrintAllDeadlocks = godo.PtrTo(v.(bool))
	}

	if v, ok := d.GetOkExists("innodb_rollback_on_timeout"); ok {
		opts.InnodbRollbackOnTimeout = godo.PtrTo(v.(bool))
	}

	if v, ok := d.GetOk("internal_tmp_mem_storage_engine"); ok {
		opts.InternalTmpMemStorageEngine = godo.PtrTo(v.(string))
	}

	if v, ok := d.GetOk("max_heap_table_size"); ok {
		opts.MaxHeapTableSize = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("tmp_table_size"); ok {
		opts.TmpTableSize = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("slow_query_log"); ok {
		opts.SlowQueryLog = godo.PtrTo(v.(bool))
	}

	if v, ok := d.GetOk("long_query_time"); ok {
		opts.LongQueryTime = godo.PtrTo(float32(v.(float64)))
	}

	if v, ok := d.GetOk("backup_hour"); ok {
		opts.BackupHour = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("backup_minute"); ok {
		opts.BackupMinute = godo.PtrTo(v.(int))
	}

	if v, ok := d.GetOk("binlog_retention_period"); ok {
		opts.BinlogRetentionPeriod = godo.PtrTo(v.(int))
	}

	log.Printf("[DEBUG] MySQL configuration: %s", godo.Stringify(opts))

	if _, err := client.Databases.UpdateMySQLConfig(ctx, clusterID, opts); err != nil {
		return err
	}

	return nil
}

func resourceDigitalOceanDatabaseMySQLConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	config, resp, err := client.Databases.GetMySQLConfig(ctx, clusterID)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving MySQL configuration: %s", err)
	}

	d.Set("connect_timeout", config.ConnectTimeout)
	d.Set("default_time_zone", config.DefaultTimeZone)
	d.Set("innodb_log_buffer_size", config.InnodbLogBufferSize)
	d.Set("innodb_online_alter_log_max_size", config.InnodbOnlineAlterLogMaxSize)
	d.Set("innodb_lock_wait_timeout", config.InnodbLockWaitTimeout)
	d.Set("interactive_timeout", config.InteractiveTimeout)
	d.Set("max_allowed_packet", config.MaxAllowedPacket)
	d.Set("net_read_timeout", config.NetReadTimeout)
	d.Set("sort_buffer_size", config.SortBufferSize)
	d.Set("sql_mode", config.SQLMode)
	d.Set("sql_require_primary_key", config.SQLRequirePrimaryKey)
	d.Set("wait_timeout", config.WaitTimeout)
	d.Set("net_write_timeout", config.NetWriteTimeout)
	d.Set("group_concat_max_len", config.GroupConcatMaxLen)
	d.Set("information_schema_stats_expiry", config.InformationSchemaStatsExpiry)
	d.Set("innodb_ft_min_token_size", config.InnodbFtMinTokenSize)
	d.Set("innodb_ft_server_stopword_table", config.InnodbFtServerStopwordTable)
	d.Set("innodb_print_all_deadlocks", config.InnodbPrintAllDeadlocks)
	d.Set("innodb_rollback_on_timeout", config.InnodbRollbackOnTimeout)
	d.Set("internal_tmp_mem_storage_engine", config.InternalTmpMemStorageEngine)
	d.Set("max_heap_table_size", config.MaxHeapTableSize)
	d.Set("tmp_table_size", config.TmpTableSize)
	d.Set("slow_query_log", config.SlowQueryLog)
	d.Set("long_query_time", config.LongQueryTime)
	d.Set("backup_hour", config.BackupHour)
	d.Set("backup_minute", config.BackupMinute)
	d.Set("binlog_retention_period", config.BinlogRetentionPeriod)

	return nil
}

func resourceDigitalOceanDatabaseMySQLConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	warn := []diag.Diagnostic{
		{
			Severity: diag.Warning,
			Summary:  "digitalocean_database_mysql_config removed from state",
			Detail:   "Database configurations are only removed from state when destroyed. The remote configuration is not unset.",
		},
	}

	return warn
}

func resourceDigitalOceanDatabaseMySQLConfigImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	clusterID := d.Id()

	d.SetId(makeDatabaseMySQLConfigID(clusterID))
	d.Set("cluster_id", clusterID)

	return []*schema.ResourceData{d}, nil
}

func makeDatabaseMySQLConfigID(clusterID string) string {
	return fmt.Sprintf("%s/mysql-config", clusterID)
}
