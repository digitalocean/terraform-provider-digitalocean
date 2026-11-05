---
page_title: "DigitalOcean: digitalocean_database_advanced_mysql_config"
subcategory: "Databases"
---

# digitalocean\_database\_advanced\_mysql\_config

Provides a virtual resource that can be used to change advanced configuration
options for a DigitalOcean managed MySQL Advanced Edition (`advanced_mysql`)
database cluster.

-> **Note** Advanced MySQL configurations are only removed from state when destroyed. The remote configuration is not unset.

## Example Usage

```hcl
resource "digitalocean_database_advanced_mysql_config" "example" {
  cluster_id = digitalocean_database_cluster.example.id

  mysql_parameters = {
    time_zone       = "SYSTEM"
    connect_timeout = "10"
  }
}

resource "digitalocean_database_cluster" "example" {
  name       = "example-advanced-mysql-cluster"
  engine     = "advanced_mysql"
  version    = "8"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}
```

## Argument Reference

The following arguments are supported. See the [DigitalOcean API documentation](https://docs.digitalocean.com/reference/api/digitalocean/#tag/Databases/operation/databases_patch_config)
for additional details on each option.

* `cluster_id` - (Required) The ID of the target MySQL Advanced Edition cluster.
* `mysql_parameters` - (Optional) A map of MySQL system variable names to their string values. Only parameters included in this map are managed by Terraform.

## Supported `mysql_parameters`

The following MySQL system variables can be set on `advanced_mysql` clusters. Default values and restart requirements are returned by the DigitalOcean API and may vary by cluster size.

* `back_log` - (Optional) Pending connection queue size. Default: `512`. **Requires restart.**
* `binlog_cache_size` - (Optional) Cache for transaction binlog events. Default: `32768`.
* `binlog_checksum` - (Optional) Checksum: CRC32 or NONE. Default: `CRC32`.
* `binlog_expire_logs_seconds` - (Optional) Seconds before binlogs purged. Default: `2592000`.
* `binlog_format` - (Optional) Format: ROW, STATEMENT, or MIXED. Default: `ROW`.
* `binlog_group_commit_sync_delay` - (Optional) Microseconds delay before sync. Default: `0`.
* `binlog_group_commit_sync_no_delay_count` - (Optional) Max txns to wait before sync. Default: `0`.
* `binlog_order_commits` - (Optional) Commit in binlog write order. Default: `ON`.
* `binlog_row_image` - (Optional) Row images: full, minimal, noblob. Default: `full`.
* `binlog_row_metadata` - (Optional) Metadata: MINIMAL or FULL. Default: `FULL`.
* `binlog_stmt_cache_size` - (Optional) Cache for non-transactional binlog events. Default: `32768`.
* `binlog_transaction_compression` - (Optional) Enable binlog txn compression. Default: `OFF`.
* `binlog_transaction_compression_level_zstd` - (Optional) zstd level (1-22). Default: `3`.
* `character_set_server` - (Optional) Server character set for new databases. Default: `utf8mb4`. **Requires restart.**
* `collation_server` - (Optional) Server default collation for new databases. Default: `utf8mb4_0900_ai_ci`. **Requires restart.**
* `connect_timeout` - (Optional) Seconds for connection handshake. Default: `10`.
* `connection_memory_chunk_size` - (Optional) Connection memory accounting chunk. Default: `8192`.
* `connection_memory_limit` - (Optional) Per-connection memory limit. Default: `9223372036854775807`.
* `cte_max_recursion_depth` - (Optional) Max CTE recursion depth. Default: `1000`.
* `eq_range_index_dive_limit` - (Optional) Ranges before switching to stats. Default: `200`.
* `explicit_defaults_for_timestamp` - (Optional) Explicit TIMESTAMP defaults. Default: `ON`.
* `general_log` - (Optional) Enable general query log. Default: `OFF`.
* `global_connection_memory_limit` - (Optional) Global memory limit all connections. Default: `9223372036854775807`.
* `global_connection_memory_tracking` - (Optional) Per-connection memory tracking. Default: `OFF`.
* `group_concat_max_len` - (Optional) Max GROUP_CONCAT() result length. Default: `1024`.
* `group_replication_communication_max_message_size` - (Optional) Max message size for Group Replication. Default: `10485760`. **Requires restart.**
* `group_replication_consistency` - (Optional) Transaction consistency level for Group Replication. Default: `BEFORE_ON_PRIMARY_FAILOVER`.
* `group_replication_flow_control_mode` - (Optional) Group Replication flow control mode. Default: `QUOTA`.
* `group_replication_flow_control_period` - (Optional) Seconds between flow control quota checks. Default: `1`. **Requires restart.**
* `group_replication_message_cache_size` - (Optional) Maximum memory used by Group Replication to cache messages. Default: `1073741824`. **Requires restart.**
* `group_replication_paxos_single_leader` - (Optional) Single-leader Paxos mode for Group Replication. Default: `OFF`.
* `group_replication_poll_spin_loops` - (Optional) Spin loops before Group Replication poll. Default: `0`.
* `group_replication_unreachable_majority_timeout` - (Optional) Seconds before partitioned member action. Default: `5`.
* `information_schema_stats_expiry` - (Optional) Seconds before cached schema stats expire. Default: `86400`.
* `innodb_adaptive_flushing` - (Optional) Adaptive flushing of dirty pages. Default: `ON`.
* `innodb_adaptive_flushing_lwm` - (Optional) Low water mark % for adaptive flushing. Default: `10`.
* `innodb_adaptive_hash_index` - (Optional) Enable/disable adaptive hash index. Default: `ON`.
* `innodb_autoextend_increment` - (Optional) Tablespace auto-extend increment in MB. Default: `64`.
* `innodb_buffer_pool_size` - (Optional) Size in bytes of the InnoDB buffer pool. Default: `134217728`.
* `innodb_change_buffer_max_size` - (Optional) Max change buffer as % of pool. Default: `25`.
* `innodb_change_buffering` - (Optional) Types of operations buffered in change buffer. Default: `all`.
* `innodb_compression_failure_threshold_pct` - (Optional) Compression failure % before padding. Default: `5`.
* `innodb_compression_level` - (Optional) zlib compression level (0-9). Default: `6`.
* `innodb_compression_pad_pct_max` - (Optional) Max % page padding for compressed tables. Default: `50`.
* `innodb_concurrency_tickets` - (Optional) Tickets for thread re-entry without concurrency check. Default: `5000`.
* `innodb_corrupt_table_action` - (Optional) On corrupt table: assert, warn, salvage. Default: `warn`.
* `innodb_ddl_threads` - (Optional) Threads for DDL sort/build operations. Default: `2`. **Requires restart.**
* `innodb_deadlock_detect` - (Optional) Enable/disable deadlock detection. Default: `ON`.
* `innodb_fill_factor` - (Optional) Fill factor for B-tree bulk load. Default: `100`.
* `innodb_flush_log_at_trx_commit` - (Optional) Controls log flushing on transaction commit. Default: `1`.
* `innodb_flush_method` - (Optional) Method for flushing data to disk. Default: `O_DIRECT`. **Requires restart.**
* `innodb_flush_sync` - (Optional) Ignore io_capacity during checkpoints. Default: `ON`.
* `innodb_fsync_threshold` - (Optional) Bytes threshold for fsync on file create. Default: `0`.
* `innodb_ft_max_token_size` - (Optional) Max word length for full-text index. Default: `84`. **Requires restart.**
* `innodb_ft_min_token_size` - (Optional) Min word length for full-text index. Default: `3`. **Requires restart.**
* `innodb_ft_server_stopword_table` - (Optional) Table for full-text stopwords.
* `innodb_io_capacity` - (Optional) Background I/O operations per second. Default: `200`.
* `innodb_io_capacity_max` - (Optional) Upper limit for background I/O operations. Default: `2000`.
* `innodb_lock_wait_timeout` - (Optional) Seconds to wait for a row lock. Default: `50`.
* `innodb_log_buffer_size` - (Optional) Size of the redo log buffer in memory. Default: `33554432`. **Requires restart.**
* `innodb_log_compressed_pages` - (Optional) Log re-compressed pages to redo log. Default: `ON`.
* `innodb_lru_scan_depth` - (Optional) How deep page cleaner scans LRU list. Default: `1024`.
* `innodb_max_dirty_pages_pct` - (Optional) Max % of dirty pages before flushing. Default: `90.000000`.
* `innodb_max_dirty_pages_pct_lwm` - (Optional) Low water mark for dirty page % preflushing. Default: `10.000000`.
* `innodb_monitor_enable` - (Optional) Enable InnoDB performance schema monitors.
* `innodb_numa_interleave` - (Optional) NUMA memory interleaving for buffer pool. Default: `OFF`. **Requires restart.**
* `innodb_old_blocks_pct` - (Optional) % of buffer pool for old block sublist. Default: `37`.
* `innodb_old_blocks_time` - (Optional) ms block stays in old sublist before promotion. Default: `1000`.
* `innodb_online_alter_log_max_size` - (Optional) Max log size for online DDL operations. Default: `134217728`.
* `innodb_open_files` - (Optional) InnoDB open files limit. Default: `4000`. **Requires restart.**
* `innodb_page_cleaners` - (Optional) Number of page cleaner threads. Default: `1`. **Requires restart.**
* `innodb_parallel_read_threads` - (Optional) Threads for parallel clustered index reads. Default: `2`. **Requires restart.**
* `innodb_print_all_deadlocks` - (Optional) Print all deadlocks to error log. Default: `OFF`.
* `innodb_print_ddl_logs` - (Optional) Print DDL logs to error log. Default: `OFF`.
* `innodb_print_lock_wait_timeout_info` - (Optional) Extra lock wait timeout info. Default: `OFF`.
* `innodb_purge_threads` - (Optional) Number of background purge threads. Default: `1`. **Requires restart.**
* `innodb_random_read_ahead` - (Optional) Enable random read-ahead. Default: `OFF`.
* `innodb_read_ahead_threshold` - (Optional) Pages to trigger linear read-ahead. Default: `56`.
* `innodb_read_io_threads` - (Optional) Number of read I/O threads. Default: `4`. **Requires restart.**
* `innodb_redo_log_capacity` - (Optional) Total redo log capacity for InnoDB. Default: `104857600`.
* `innodb_rollback_on_timeout` - (Optional) Rollback entire transaction on lock wait timeout. Default: `OFF`. **Requires restart.**
* `innodb_show_locks_held` - (Optional) Max locks shown per txn. Default: `10`.
* `innodb_sort_buffer_size` - (Optional) Sort buffer size for InnoDB online DDL. Default: `1048576`. **Requires restart.**
* `innodb_spin_wait_delay` - (Optional) Delay between spin lock polls. Default: `6`.
* `innodb_stats_auto_recalc` - (Optional) Auto-recalculate persistent stats. Default: `ON`.
* `innodb_stats_persistent` - (Optional) Whether index statistics are persistent. Default: `ON`.
* `innodb_stats_persistent_sample_pages` - (Optional) Sample pages for persistent statistics. Default: `20`.
* `innodb_stats_transient_sample_pages` - (Optional) Sample pages for transient statistics. Default: `8`.
* `innodb_strict_mode` - (Optional) Strict mode (errors vs warnings). Default: `ON`.
* `innodb_sync_spin_loops` - (Optional) Spin loops before mutex wait. Default: `30`.
* `innodb_thread_concurrency` - (Optional) Max threads inside InnoDB (0=unlimited). Default: `0`.
* `innodb_use_fdatasync` - (Optional) Use fdatasync() instead of fsync(). Default: `OFF`.
* `innodb_write_io_threads` - (Optional) Number of write I/O threads. Default: `4`. **Requires restart.**
* `interactive_timeout` - (Optional) Seconds the server waits for activity on an interactive connection before closing it. Default: `28800`.
* `internal_tmp_mem_storage_engine` - (Optional) Engine for temp tables: TempTable or MEMORY. Default: `TempTable`.
* `join_buffer_size` - (Optional) Per-session buffer for joins. Default: `262144`.
* `kill_idle_transaction` - (Optional) Seconds before killing idle transactions (Percona; was catalog misspell innodb_kill_idle_transaction). Default: `0`.
* `local_infile` - (Optional) Enable LOAD DATA LOCAL INFILE. Default: `OFF`.
* `lock_wait_timeout` - (Optional) Seconds for metadata lock wait. Default: `31536000`.
* `log_error_verbosity` - (Optional) Error log verbosity (1-3). Default: `2`.
* `log_output` - (Optional) Log destination: TABLE, FILE, NONE. Default: `FILE`.
* `log_queries_not_using_indexes` - (Optional) Log no-index queries to slow log. Default: `OFF`.
* `log_slow_extra` - (Optional) Extra info in slow query log. Default: `OFF`.
* `log_slow_filter` - (Optional) PS 8.4 filter (comma-separated): full_scan, full_join, tmp_table, tmp_table_on_disk, filesort, filesort_on_disk; empty disables.
* `log_slow_rate_limit` - (Optional) Rate-limit slow log (1=every query). Default: `1`.
* `log_slow_sp_statements` - (Optional) Log stored proc statements. Default: `OFF`.
* `log_slow_verbosity` - (Optional) Detail: microtime, query_plan, innodb.
* `log_throttle_queries_not_using_indexes` - (Optional) Throttle rate of no-index log entries. Default: `0`.
* `log_timestamps` - (Optional) Log timestamps: UTC or SYSTEM. Default: `UTC`.
* `long_query_time` - (Optional) Queries that take longer than this many seconds are logged to the slow query log. Default: `2`.
* `max_allowed_packet` - (Optional) Max client/server packet size. Default: `67108864`.
* `max_connect_errors` - (Optional) Consecutive errors before blocking host. Default: `100`.
* `max_connections` - (Optional) Maximum number of simultaneous client connections. Default: `512`.
* `max_heap_table_size` - (Optional) Maximum size for user-created MEMORY tables and internal in-memory tables. Default: `16777216`.
* `max_join_size` - (Optional) Max rows for large join protection. Default: `9223372036854775807`.
* `max_prepared_stmt_count` - (Optional) Max prepared statements server-wide. Default: `16382`.
* `max_sort_length` - (Optional) Bytes sorting BLOB/TEXT values. Default: `1024`.
* `max_user_connections` - (Optional) Max simultaneous connections per user. Default: `0`.
* `net_buffer_length` - (Optional) Initial connection/result buffer size. Default: `16384`.
* `net_read_timeout` - (Optional) Seconds to wait for read data. Default: `30`.
* `net_write_timeout` - (Optional) Seconds to wait for write. Default: `60`.
* `optimizer_switch` - (Optional) Optimizer feature flags (boot_val from PS 8.4.8 live; hypergraph_optimizer=on fails on non-debug builds). Default: `index_merge=on,index_merge_union=on,index_merge_sort_union=on,index_merge_intersection=on,engine_condition_pushdown=on,index_condition_pushdown=on,mrr=on,mrr_cost_based=on,block_nested_loop=on,batched_key_access=off,materialization=on,semijoin=on,loosescan=on,firstmatch=on,duplicateweedout=on,subquery_materialization_cost_based=on,use_index_extensions=on,condition_fanout_filter=on,derived_merge=on,use_invisible_indexes=off,skip_scan=on,hash_join=on,subquery_to_derived=off,prefer_ordering_index=on,hypergraph_optimizer=off,derived_condition_pushdown=on,hash_set_operations=on,favor_range_scan=off`.
* `password_history` - (Optional) Passwords tracked for reuse prevention. Default: `0`.
* `password_reuse_interval` - (Optional) Days before password reuse. Default: `0`.
* `performance_schema_max_digest_length` - (Optional) Max Performance Schema digest length. Default: `1024`. **Requires restart.**
* `performance_schema_max_sql_text_length` - (Optional) Max Performance Schema SQL text length. Default: `1024`. **Requires restart.**
* `range_optimizer_max_mem_size` - (Optional) Max memory for range optimizer. Default: `8388608`.
* `read_buffer_size` - (Optional) Buffer for sequential scans. Default: `131072`.
* `read_rnd_buffer_size` - (Optional) Buffer for random reads after sort. Default: `262144`.
* `replica_compressed_protocol` - (Optional) Compression for replica protocol. Default: `OFF`.
* `replica_exec_mode` - (Optional) STRICT or IDEMPOTENT. Default: `STRICT`.
* `replica_parallel_type` - (Optional) LOGICAL_CLOCK or DATABASE. Default: `LOGICAL_CLOCK`.
* `replica_parallel_workers` - (Optional) Parallel applier workers. Default: `4`.
* `replica_preserve_commit_order` - (Optional) Preserve commit order on replicas. Default: `ON`.
* `require_secure_transport` - (Optional) Require SSL/TLS. Default: `OFF`.
* `slow_query_log` - (Optional) Whether the slow query log is enabled. Default: `ON`.
* `sort_buffer_size` - (Optional) Per-session buffer for sorts. Default: `262144`.
* `sql_buffer_result` - (Optional) Force results to temp tables. Default: `OFF`.
* `sql_mode` - (Optional) SQL modes that control SQL syntax and data validation. Default: `ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION`.
* `sync_binlog` - (Optional) Sync binlog every N commits. Default: `1`.
* `table_definition_cache` - (Optional) Table definitions to cache. Default: `2000`.
* `table_open_cache` - (Optional) Open tables to cache. Default: `4000`.
* `table_open_cache_instances` - (Optional) Table cache instances. Default: `16`. **Requires restart.**
* `tablespace_definition_cache` - (Optional) Tablespace definitions to cache. Default: `256`.
* `temptable_max_mmap` - (Optional) Max TempTable mmap file size. Default: `1073741824`.
* `temptable_max_ram` - (Optional) Max TempTable RAM before spill. Default: `1073741824`.
* `temptable_use_mmap` - (Optional) Allow TempTable mmap files. Default: `ON`.
* `thread_cache_size` - (Optional) Threads cached for reuse. Default: `9`.
* `thread_pool_idle_timeout` - (Optional) Seconds before idle pool thread exits. Default: `60`.
* `thread_pool_max_threads` - (Optional) Max threads in thread pool. Default: `100000`.
* `thread_pool_size` - (Optional) Thread groups in Percona thread pool. Default: `4`. **Requires restart.**
* `thread_stack` - (Optional) Stack size per thread. Default: `1048576`. **Requires restart.**
* `thread_statistics` - (Optional) Per-thread statistics. Default: `OFF`.
* `time_zone` - (Optional) Server default time zone (replaces invalid catalog name default_time_zone). Default: `SYSTEM`.
* `tmp_table_size` - (Optional) Maximum size of internal in-memory temporary tables. Default: `16777216`.
* `userstat` - (Optional) USER_STATISTICS tables. Default: `OFF`.
* `wait_timeout` - (Optional) Seconds the server waits for activity on a noninteractive connection before closing it. Default: `28800`.
* `windowing_use_high_precision` - (Optional) High precision window functions. Default: `ON`.

## Attributes Reference

All above attributes are exported.

## Import

An advanced MySQL database cluster's configuration can be imported using the `id` of the parent cluster, e.g.

```bash
terraform import digitalocean_database_advanced_mysql_config.example 52556c07-788e-4d41-b8a7-c796432197d1
```
