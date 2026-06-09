---
page_title: "DigitalOcean: digitalocean_database_advanced_postgresql_config"
subcategory: "Databases"
---

# digitalocean\_database\_advanced\_postgresql\_config

Provides a virtual resource that can be used to change advanced configuration
options for a DigitalOcean managed PostgreSQL Advanced Edition (`advanced_pg`)
database cluster.

-> **Note** Advanced PostgreSQL configurations are only removed from state when destroyed. The remote configuration is not unset.

## Example Usage

```hcl
resource "digitalocean_database_advanced_postgresql_config" "example" {
  cluster_id = digitalocean_database_cluster.example.id

  pg_parameters = {
    timezone = "UTC"
    work_mem = "4096"
  }
}

resource "digitalocean_database_cluster" "example" {
  name       = "example-advanced-postgresql-cluster"
  engine     = "advanced_pg"
  version    = "16"
  size       = "gd-2vcpu-8gb-intel"
  region     = "nyc1"
  node_count = 1
}
```

## Argument Reference

The following arguments are supported. See the [DigitalOcean API documentation](https://docs.digitalocean.com/reference/api/digitalocean/#tag/Databases/operation/databases_patch_config)
for additional details on each option.

* `cluster_id` - (Required) The ID of the target PostgreSQL Advanced Edition cluster.
* `pg_parameters` - (Optional) A map of PostgreSQL GUC parameter names to their string values. Only parameters included in this map are managed by Terraform. Values use PostgreSQL GUC syntax (for example, `work_mem` is specified in kilobytes unless a unit suffix is provided).

## Supported `pg_parameters`

The following PostgreSQL GUC parameters can be set on `advanced_pg` clusters. Default values and restart requirements are returned by the DigitalOcean API and may vary by cluster size.

* `application_name` - (Optional) Sets the application name to be reported in statistics and logs.
* `array_nulls` - (Optional) Enables input of NULL elements in arrays. Default: `on`.
* `authentication_timeout` - (Optional) Sets the maximum allowed time to complete client authentication. Default: `60`.
* `autovacuum` - (Optional) Starts the autovacuum subprocess. Default: `on`.
* `autovacuum_analyze_scale_factor` - (Optional) Number of tuple inserts, updates, or deletes prior to analyze as a fraction of reltuples. Default: `0.1`.
* `autovacuum_analyze_threshold` - (Optional) Minimum number of tuple inserts, updates, or deletes prior to analyze. Default: `50`.
* `autovacuum_freeze_max_age` - (Optional) Age at which to autovacuum a table to prevent transaction ID wraparound. Default: `200000000`. **Requires restart.**
* `autovacuum_multixact_freeze_max_age` - (Optional) Multixact age at which to autovacuum a table to prevent multixact wraparound. Default: `400000000`. **Requires restart.**
* `autovacuum_naptime` - (Optional) Time to sleep between autovacuum runs. Default: `60`.
* `autovacuum_vacuum_cost_delay` - (Optional) Vacuum cost delay in milliseconds, for autovacuum. Default: `2`.
* `autovacuum_vacuum_cost_limit` - (Optional) Vacuum cost amount available before napping, for autovacuum. Default: `-1`.
* `autovacuum_vacuum_insert_scale_factor` - (Optional) Number of tuple inserts prior to vacuum as a fraction of reltuples. Default: `0.2`.
* `autovacuum_vacuum_insert_threshold` - (Optional) Minimum number of tuple inserts prior to vacuum. Default: `1000`.
* `autovacuum_vacuum_scale_factor` - (Optional) Number of tuple updates or deletes prior to vacuum as a fraction of reltuples. Default: `0.2`.
* `autovacuum_vacuum_threshold` - (Optional) Minimum number of tuple updates or deletes prior to vacuum. Default: `50`.
* `backend_flush_after` - (Optional) Number of pages after which previously performed writes are flushed to disk. Default: `0`.
* `backslash_quote` - (Optional) Sets whether "\'" is allowed in string literals. Default: `safe_encoding`.
* `bgwriter_delay` - (Optional) Background writer sleep time between rounds. Default: `200`.
* `bgwriter_flush_after` - (Optional) Number of pages after which previously performed writes are flushed to disk. Default: `64`.
* `bgwriter_lru_maxpages` - (Optional) Background writer maximum number of LRU pages to flush per round. Default: `100`.
* `bgwriter_lru_multiplier` - (Optional) Multiple of the average buffer usage to free per round. Default: `2`.
* `bytea_output` - (Optional) Sets the output format for bytea. Default: `hex`.
* `check_function_bodies` - (Optional) Check routine bodies during CREATE FUNCTION and CREATE PROCEDURE. Default: `on`.
* `checkpoint_completion_target` - (Optional) Time spent flushing dirty buffers during checkpoint, as fraction of checkpoint interval. Default: `0.9`.
* `checkpoint_flush_after` - (Optional) Number of pages after which previously performed writes are flushed to disk. Default: `32`.
* `checkpoint_timeout` - (Optional) Sets the maximum time between automatic WAL checkpoints. Default: `300`.
* `checkpoint_warning` - (Optional) Sets the maximum time before warning if checkpoints triggered by WAL volume happen too frequently. Default: `30`.
* `client_connection_check_interval` - (Optional) Sets the time interval between checks for disconnection while running queries. Default: `0`.
* `client_encoding` - (Optional) Sets the client's character set encoding. Default: `SQL_ASCII`.
* `client_min_messages` - (Optional) Sets the message levels that are sent to the client. Default: `notice`.
* `commit_delay` - (Optional) Sets the delay in microseconds between transaction commit and flushing WAL to disk. Default: `0`.
* `commit_siblings` - (Optional) Sets the minimum number of concurrent open transactions required before performing "commit_delay". Default: `5`.
* `compute_query_id` - (Optional) Enables in-core computation of query identifiers. Default: `auto`.
* `constraint_exclusion` - (Optional) Enables the planner to use constraints to optimize queries. Default: `partition`.
* `cpu_index_tuple_cost` - (Optional) Sets the planner's estimate of the cost of processing each index entry during an index scan. Default: `0.005`.
* `cpu_operator_cost` - (Optional) Sets the planner's estimate of the cost of processing each operator or function call. Default: `0.0025`.
* `cpu_tuple_cost` - (Optional) Sets the planner's estimate of the cost of processing each tuple (row). Default: `0.01`.
* `createrole_self_grant` - (Optional) Sets whether a CREATEROLE user automatically grants the role to themselves, and with which options.
* `cursor_tuple_fraction` - (Optional) Sets the planner's estimate of the fraction of a cursor's rows that will be retrieved. Default: `0.1`.
* `datestyle` - (Optional) Sets the display format for date and time values. Default: `ISO, MDY`.
* `deadlock_timeout` - (Optional) Sets the time to wait on a lock before checking for deadlock. Default: `1000`.
* `debug_pretty_print` - (Optional) Indents parse and plan tree displays. Default: `on`.
* `debug_print_parse` - (Optional) Logs each query's parse tree. Default: `off`.
* `debug_print_plan` - (Optional) Logs each query's execution plan. Default: `off`.
* `debug_print_rewritten` - (Optional) Logs each query's rewritten parse tree. Default: `off`.
* `default_statistics_target` - (Optional) Sets the default statistics target. Default: `100`.
* `default_tablespace` - (Optional) Sets the default tablespace to create tables and indexes in.
* `default_toast_compression` - (Optional) Sets the default compression method for compressible values. Default: `pglz`.
* `default_transaction_deferrable` - (Optional) Sets the default deferrable status of new transactions. Default: `off`.
* `default_transaction_isolation` - (Optional) Sets the transaction isolation level of each new transaction. Default: `read committed`.
* `default_transaction_read_only` - (Optional) Sets the default read-only status of new transactions. Default: `off`.
* `effective_io_concurrency` - (Optional) Number of simultaneous requests that can be handled efficiently by the disk subsystem. Default: `16`.
* `enable_async_append` - (Optional) Enables the planner's use of async append plans. Default: `on`.
* `enable_bitmapscan` - (Optional) Enables the planner's use of bitmap-scan plans. Default: `on`.
* `enable_gathermerge` - (Optional) Enables the planner's use of gather merge plans. Default: `on`.
* `enable_group_by_reordering` - (Optional) Enables reordering of GROUP BY keys. Default: `on`.
* `enable_hashagg` - (Optional) Enables the planner's use of hashed aggregation plans. Default: `on`.
* `enable_hashjoin` - (Optional) Enables the planner's use of hash join plans. Default: `on`.
* `enable_incremental_sort` - (Optional) Enables the planner's use of incremental sort steps. Default: `on`.
* `enable_indexonlyscan` - (Optional) Enables the planner's use of index-only-scan plans. Default: `on`.
* `enable_indexscan` - (Optional) Enables the planner's use of index-scan plans. Default: `on`.
* `enable_material` - (Optional) Enables the planner's use of materialization. Default: `on`.
* `enable_memoize` - (Optional) Enables the planner's use of memoization. Default: `on`.
* `enable_mergejoin` - (Optional) Enables the planner's use of merge join plans. Default: `on`.
* `enable_nestloop` - (Optional) Enables the planner's use of nested-loop join plans. Default: `on`.
* `enable_parallel_append` - (Optional) Enables the planner's use of parallel append plans. Default: `on`.
* `enable_parallel_hash` - (Optional) Enables the planner's use of parallel hash plans. Default: `on`.
* `enable_partition_pruning` - (Optional) Enables plan-time and execution-time partition pruning. Default: `on`.
* `enable_partitionwise_aggregate` - (Optional) Enables partitionwise aggregation and grouping. Default: `off`.
* `enable_partitionwise_join` - (Optional) Enables partitionwise join. Default: `off`.
* `enable_presorted_aggregate` - (Optional) Enables the planner's ability to produce plans that provide presorted input for ORDER BY / DISTINCT aggregate functions. Default: `on`.
* `enable_seqscan` - (Optional) Enables the planner's use of sequential-scan plans. Default: `on`.
* `enable_sort` - (Optional) Enables the planner's use of explicit sort steps. Default: `on`.
* `enable_tidscan` - (Optional) Enables the planner's use of TID scan plans. Default: `on`.
* `escape_string_warning` - (Optional) Warn about backslash escapes in ordinary string literals. Default: `on`.
* `event_triggers` - (Optional) Enables event triggers. Default: `on`.
* `extra_float_digits` - (Optional) Sets the number of digits displayed for floating-point values. Default: `1`.
* `from_collapse_limit` - (Optional) Sets the FROM-list size beyond which subqueries are not collapsed. Default: `8`.
* `geqo` - (Optional) Enables genetic query optimization. Default: `on`.
* `geqo_effort` - (Optional) GEQO: effort is used to set the default for other GEQO parameters. Default: `5`.
* `geqo_generations` - (Optional) GEQO: number of iterations of the algorithm. Default: `0`.
* `geqo_pool_size` - (Optional) GEQO: number of individuals in the population. Default: `0`.
* `geqo_seed` - (Optional) GEQO: seed for random path selection. Default: `0`.
* `geqo_selection_bias` - (Optional) GEQO: selective pressure within the population. Default: `2`.
* `geqo_threshold` - (Optional) Sets the threshold of FROM items beyond which GEQO is used. Default: `12`.
* `gin_fuzzy_search_limit` - (Optional) Sets the maximum allowed result for exact search by GIN. Default: `0`.
* `gin_pending_list_limit` - (Optional) Sets the maximum size of the pending list for GIN index. Default: `4096`.
* `gss_accept_delegation` - (Optional) Sets whether GSSAPI delegation should be accepted from the client. Default: `off`.
* `hash_mem_multiplier` - (Optional) Multiple of "work_mem" to use for hash tables. Default: `2`.
* `icu_validation_level` - (Optional) Log level for reporting invalid ICU locale strings. Default: `warning`.
* `idle_in_transaction_session_timeout` - (Optional) Sets the maximum allowed idle time between queries, when in a transaction. Default: `0`.
* `idle_session_timeout` - (Optional) Sets the maximum allowed idle time between queries, when not in a transaction. Default: `0`.
* `intervalstyle` - (Optional) Sets the display format for interval values. Default: `postgres`.
* `io_combine_limit` - (Optional) Limit on the size of data reads and writes. Default: `16`.
* `jit` - (Optional) Allow JIT compilation. Default: `on`.
* `jit_above_cost` - (Optional) Perform JIT compilation if query is more expensive. Default: `100000`.
* `jit_inline_above_cost` - (Optional) Perform JIT inlining if query is more expensive. Default: `500000`.
* `jit_optimize_above_cost` - (Optional) Optimize JIT-compiled functions if query is more expensive. Default: `500000`.
* `join_collapse_limit` - (Optional) Sets the FROM-list size beyond which JOIN constructs are not flattened. Default: `8`.
* `lc_messages` - (Optional) Sets the language in which messages are displayed.
* `lc_monetary` - (Optional) Sets the locale for formatting monetary amounts. Default: `C`.
* `lc_numeric` - (Optional) Sets the locale for formatting numbers. Default: `C`.
* `lc_time` - (Optional) Sets the locale for formatting date and time values. Default: `C`.
* `log_autovacuum_min_duration` - (Optional) Sets the minimum execution time above which autovacuum actions will be logged. Default: `600000`.
* `log_checkpoints` - (Optional) Logs each checkpoint. Default: `on`.
* `log_connections` - (Optional) Logs specified aspects of connection establishment and setup.
* `log_disconnections` - (Optional) Logs end of a session, including duration. Default: `off`.
* `log_duration` - (Optional) Logs the duration of each completed SQL statement. Default: `off`.
* `log_error_verbosity` - (Optional) Sets the verbosity of logged messages. Default: `default`.
* `log_executor_stats` - (Optional) Writes executor performance statistics to the server log. Default: `off`.
* `log_hostname` - (Optional) Logs the host name in the connection logs. Default: `off`.
* `log_lock_waits` - (Optional) Logs long lock waits. Default: `off`.
* `log_min_duration_sample` - (Optional) Sets the minimum execution time above which a sample of statements will be logged. Sampling is determined by "log_statement_sample_rate". Default: `-1`.
* `log_min_duration_statement` - (Optional) Sets the minimum execution time above which all statements will be logged. Default: `-1`.
* `log_min_error_statement` - (Optional) Causes all statements generating error at or above this level to be logged. Default: `error`.
* `log_min_messages` - (Optional) Sets the message levels that are logged. Default: `warning`.
* `log_parameter_max_length` - (Optional) Sets the maximum length in bytes of data logged for bind parameter values when logging statements. Default: `-1`.
* `log_parameter_max_length_on_error` - (Optional) Sets the maximum length in bytes of data logged for bind parameter values when logging statements, on error. Default: `0`.
* `log_parser_stats` - (Optional) Writes parser performance statistics to the server log. Default: `off`.
* `log_planner_stats` - (Optional) Writes planner performance statistics to the server log. Default: `off`.
* `log_recovery_conflict_waits` - (Optional) Logs standby recovery conflict waits. Default: `off`.
* `log_replication_commands` - (Optional) Logs each replication command. Default: `off`.
* `log_rotation_size` - (Optional) Sets the maximum size a log file can reach before being rotated. Default: `10240`.
* `log_startup_progress_interval` - (Optional) Time between progress updates for long-running startup operations. Default: `10000`.
* `log_statement` - (Optional) Sets the type of statements logged. Default: `none`.
* `log_statement_sample_rate` - (Optional) Fraction of statements exceeding "log_min_duration_sample" to be logged. Default: `1`.
* `log_statement_stats` - (Optional) Writes cumulative performance statistics to the server log. Default: `off`.
* `log_temp_files` - (Optional) Log the use of temporary files larger than this number of kilobytes. Default: `-1`.
* `log_transaction_sample_rate` - (Optional) Sets the fraction of transactions from which to log all statements. Default: `0`.
* `logical_decoding_work_mem` - (Optional) Sets the maximum memory to be used for logical decoding. Default: `65536`.
* `maintenance_io_concurrency` - (Optional) A variant of "effective_io_concurrency" that is used for maintenance work. Default: `16`.
* `max_files_per_process` - (Optional) Sets the maximum number of files each server process is allowed to open simultaneously. Default: `1000`. **Requires restart.**
* `max_locks_per_transaction` - (Optional) Sets the maximum number of locks per transaction. Default: `64`. **Requires restart.**
* `max_logical_replication_workers` - (Optional) Maximum number of logical replication worker processes. Default: `4`. **Requires restart.**
* `max_notify_queue_pages` - (Optional) Sets the maximum number of allocated pages for NOTIFY / LISTEN queue. Default: `1048576`. **Requires restart.**
* `max_parallel_apply_workers_per_subscription` - (Optional) Maximum number of parallel apply workers per subscription. Default: `2`.
* `max_parallel_maintenance_workers` - (Optional) Sets the maximum number of parallel processes per maintenance operation. Default: `2`.
* `max_parallel_workers` - (Optional) Sets the maximum number of parallel workers that can be active at one time. Default: `8`.
* `max_parallel_workers_per_gather` - (Optional) Sets the maximum number of parallel processes per executor node. Default: `2`.
* `max_pred_locks_per_page` - (Optional) Sets the maximum number of predicate-locked tuples per page. Default: `2`.
* `max_pred_locks_per_relation` - (Optional) Sets the maximum number of predicate-locked pages and tuples per relation. Default: `-2`.
* `max_pred_locks_per_transaction` - (Optional) Sets the maximum number of predicate locks per transaction. Default: `64`. **Requires restart.**
* `max_prepared_transactions` - (Optional) Sets the maximum number of simultaneously prepared transactions. Default: `0`. **Requires restart.**
* `max_slot_wal_keep_size` - (Optional) Sets the maximum WAL size that can be reserved by replication slots. Default: `-1`.
* `max_standby_archive_delay` - (Optional) Sets the maximum delay before canceling queries when a hot standby server is processing archived WAL data. Default: `30000`.
* `max_standby_streaming_delay` - (Optional) Sets the maximum delay before canceling queries when a hot standby server is processing streamed WAL data. Default: `30000`.
* `max_sync_workers_per_subscription` - (Optional) Maximum number of table synchronization workers per subscription. Default: `2`.
* `min_dynamic_shared_memory` - (Optional) Amount of dynamic shared memory reserved at startup. Default: `0`. **Requires restart.**
* `min_parallel_index_scan_size` - (Optional) Sets the minimum amount of index data for a parallel scan. Default: `64`.
* `min_parallel_table_scan_size` - (Optional) Sets the minimum amount of table data for a parallel scan. Default: `1024`.
* `parallel_leader_participation` - (Optional) Controls whether Gather and Gather Merge also run subplans. Default: `on`.
* `parallel_setup_cost` - (Optional) Sets the planner's estimate of the cost of starting up worker processes for parallel query. Default: `1000`.
* `parallel_tuple_cost` - (Optional) Sets the planner's estimate of the cost of passing each tuple (row) from worker to leader backend. Default: `0.1`.
* `plan_cache_mode` - (Optional) Controls the planner's selection of custom or generic plan. Default: `auto`.
* `quote_all_identifiers` - (Optional) When generating SQL fragments, quote all identifiers. Default: `off`.
* `random_page_cost` - (Optional) Sets the planner's estimate of the cost of a nonsequentially fetched disk page. Default: `4`.
* `recursive_worktable_factor` - (Optional) Sets the planner's estimate of the average size of a recursive query's working table. Default: `10`.
* `remove_temp_files_after_crash` - (Optional) Remove temporary files after backend crash. Default: `on`.
* `row_security` - (Optional) Enables row security. Default: `on`.
* `scram_iterations` - (Optional) Sets the iteration count for SCRAM secret generation. Default: `4096`.
* `search_path` - (Optional) Sets the schema search order for names that are not schema-qualified. Default: `"$user", public`.
* `seq_page_cost` - (Optional) Sets the planner's estimate of the cost of a sequentially fetched disk page. Default: `1`.
* `session_replication_role` - (Optional) Sets the session's behavior for triggers and rewrite rules. Default: `origin`.
* `shared_buffers` - (Optional) Sets the number of shared memory buffers used by the server. Default: `16384`. **Requires restart.**
* `standard_conforming_strings` - (Optional) Causes '...' strings to treat backslashes literally. Default: `on`.
* `statement_timeout` - (Optional) Sets the maximum allowed duration of any statement. Default: `0`.
* `stats_fetch_consistency` - (Optional) Sets the consistency of accesses to statistics data. Default: `cache`.
* `synchronize_seqscans` - (Optional) Enables synchronized sequential scans. Default: `on`.
* `synchronous_commit` - (Optional) Sets the current transaction's synchronization level. Default: `on`.
* `tcp_keepalives_count` - (Optional) Maximum number of TCP keepalive retransmits. Default: `0`.
* `tcp_keepalives_idle` - (Optional) Time between issuing TCP keepalives. Default: `0`.
* `tcp_keepalives_interval` - (Optional) Time between TCP keepalive retransmits. Default: `0`.
* `temp_buffers` - (Optional) Sets the maximum number of temporary buffers used by each session. Default: `1024`.
* `temp_file_limit` - (Optional) Limits the total size of all temporary files used by each process. Default: `-1`.
* `temp_tablespaces` - (Optional) Sets the tablespace(s) to use for temporary tables and sort files.
* `timezone` - (Optional) Sets the time zone for displaying and interpreting time stamps. Default: `GMT`.
* `trace_connection_negotiation` - (Optional) Logs details of pre-authentication connection handshake. Default: `off`. **Requires restart.**
* `track_activities` - (Optional) Collects information about executing commands. Default: `on`.
* `track_activity_query_size` - (Optional) Sets the size reserved for pg_stat_activity.query, in bytes. Default: `1024`. **Requires restart.**
* `track_functions` - (Optional) Collects function-level statistics on database activity. Default: `none`.
* `track_io_timing` - (Optional) Collects timing statistics for database I/O activity. Default: `off`.
* `track_wal_io_timing` - (Optional) Collects timing statistics for WAL I/O activity. Default: `off`.
* `transform_null_equals` - (Optional) Treats "expr=NULL" as "expr IS NULL". Default: `off`.
* `vacuum_buffer_usage_limit` - (Optional) Sets the buffer pool size for VACUUM, ANALYZE, and autovacuum. Default: `2048`.
* `vacuum_cost_delay` - (Optional) Vacuum cost delay in milliseconds. Default: `0`.
* `vacuum_cost_limit` - (Optional) Vacuum cost amount available before napping. Default: `200`.
* `vacuum_cost_page_dirty` - (Optional) Vacuum cost for a page dirtied by vacuum. Default: `20`.
* `vacuum_cost_page_hit` - (Optional) Vacuum cost for a page found in the buffer cache. Default: `1`.
* `vacuum_cost_page_miss` - (Optional) Vacuum cost for a page not found in the buffer cache. Default: `2`.
* `vacuum_failsafe_age` - (Optional) Age at which VACUUM should trigger failsafe to avoid a wraparound outage. Default: `1600000000`.
* `vacuum_freeze_min_age` - (Optional) Minimum age at which VACUUM should freeze a table row. Default: `50000000`.
* `vacuum_freeze_table_age` - (Optional) Age at which VACUUM should scan whole table to freeze tuples. Default: `150000000`.
* `vacuum_multixact_failsafe_age` - (Optional) Multixact age at which VACUUM should trigger failsafe to avoid a wraparound outage. Default: `1600000000`.
* `vacuum_multixact_freeze_min_age` - (Optional) Minimum age at which VACUUM should freeze a MultiXactId in a table row. Default: `5000000`.
* `vacuum_multixact_freeze_table_age` - (Optional) Multixact age at which VACUUM should scan whole table to freeze tuples. Default: `150000000`.
* `wal_compression` - (Optional) Compresses full-page writes written in WAL file with specified method. Default: `off`.
* `wal_receiver_status_interval` - (Optional) Sets the maximum interval between WAL receiver status reports to the sending server. Default: `10`.
* `wal_receiver_timeout` - (Optional) Sets the maximum wait time to receive data from the sending server. Default: `60000`.
* `wal_sender_timeout` - (Optional) Sets the maximum time to wait for WAL replication. Default: `60000`.
* `wal_skip_threshold` - (Optional) Minimum size of new file to fsync instead of writing WAL. Default: `2048`.
* `wal_writer_delay` - (Optional) Time between WAL flushes performed in the WAL writer. Default: `200`.
* `wal_writer_flush_after` - (Optional) Amount of WAL written out by WAL writer that triggers a flush. Default: `128`.
* `work_mem` - (Optional) Sets the maximum memory to be used for query workspaces. Default: `4096`.
* `xmlbinary` - (Optional) Sets how binary values are to be encoded in XML. Default: `base64`.
* `xmloption` - (Optional) Sets whether XML data in implicit parsing and serialization operations is to be considered as documents or content fragments. Default: `content`.

## Attributes Reference

All above attributes are exported.

## Import

An advanced PostgreSQL database cluster's configuration can be imported using the `id` of the parent cluster, e.g.

```bash
terraform import digitalocean_database_advanced_postgresql_config.example 52556c07-788e-4d41-b8a7-c796432197d1
```
