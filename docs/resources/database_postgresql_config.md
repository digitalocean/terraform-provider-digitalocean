---
page_title: "DigitalOcean: digitalocean_database_postgresql_config"
---

# digitalocean\_database\_postgresql\_config

Provides a virtual resource that can be used to change advanced configuration
options for a DigitalOcean managed PostgreSQL database cluster.

-> **Note** PostgreSQL configurations are only removed from state when destroyed. The remote configuration is not unset.

## Example Usage

```hcl
resource "digitalocean_database_postgresql_config" "example" {
  cluster_id      = digitalocean_database_cluster.example.id
  connect_timeout = 10
  time_zone       = "UTC"
  work_mem        = 16
}

resource "digitalocean_database_cluster" "example" {
  name       = "example-postgresql-cluster"
  engine     = "pg"
  version    = "8"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}
```

## Argument Reference

The following arguments are supported. See the [DigitalOcean API documentation](https://docs.digitalocean.com/reference/api/api-reference/#operation/databases_patch_config)
for additional details on each option.

* `cluster_id` - (Required)  The ID of the target PostgreSQL cluster.
* `autovacuum_freeze_max_age` - (Optional)  Specifies the maximum age (in transactions) that a table's pg_class.relfrozenxid field can attain before a VACUUM operation is forced to prevent transaction ID wraparound within the table. Note that the system will launch autovacuum processes to prevent wraparound even when autovacuum is otherwise disabled. This parameter will cause the server to be restarted.
* `autovacuum_max_workers` - (Optional)  Specifies the maximum number of autovacuum processes (other than the autovacuum launcher) that may be running at any one time. The default is three. This parameter can only be set at server start.
* `autovacuum_naptime` - (Optional)  Specifies the minimum delay, in seconds, between autovacuum runs on any given database. The default is one minute.
* `autovacuum_vacuum_threshold` - (Optional)  Specifies the minimum number of updated or deleted tuples needed to trigger a VACUUM in any one table. The default is 50 tuples.
* `autovacuum_analyze_threshold` - (Optional)  Specifies the minimum number of inserted, updated, or deleted tuples needed to trigger an ANALYZE in any one table. The default is 50 tuples.
* `autovacuum_vacuum_scale_factor` - (Optional)  Specifies a fraction, in a decimal value, of the table size to add to autovacuum_vacuum_threshold when deciding whether to trigger a VACUUM. The default is 0.2 (20% of table size).
* `autovacuum_analyze_scale_factor` - (Optional)  Specifies a fraction, in a decimal value, of the table size to add to autovacuum_analyze_threshold when deciding whether to trigger an ANALYZE. The default is 0.2 (20% of table size).
* `autovacuum_vacuum_cost_delay` - (Optional)  Specifies the cost delay value, in milliseconds, that will be used in automatic VACUUM operations. If -1, uses the regular vacuum_cost_delay value, which is 20 milliseconds.
* `autovacuum_vacuum_cost_limit` - (Optional)  Specifies the cost limit value that will be used in automatic VACUUM operations. If -1 is specified (which is the default), the regular vacuum_cost_limit value will be used.
* `bgwriter_delay` - (Optional)  Specifies the delay, in milliseconds, between activity rounds for the background writer. Default is 200 ms.
* `bgwriter_flush_after` - (Optional)  The amount of kilobytes that need to be written by the background writer before attempting to force the OS to issue these writes to underlying storage. Specified in kilobytes, default is 512. Setting of 0 disables forced writeback.
* `bgwriter_lru_maxpages` - (Optional)  The maximum number of buffers that the background writer can write. Setting this to zero disables background writing. Default is 100.
* `bgwriter_lru_multiplier` - (Optional)  The average recent need for new buffers is multiplied by bgwriter_lru_multiplier to arrive at an estimate of the number that will be needed during the next round, (up to bgwriter_lru_maxpages). 1.0 represents a “just in time” policy of writing exactly the number of buffers predicted to be needed. Larger values provide some cushion against spikes in demand, while smaller values intentionally leave writes to be done by server processes. The default is 2.0.
* `deadlock_timeout` - (Optional)  The amount of time, in milliseconds, to wait on a lock before checking to see if there is a deadlock condition.
* `default_toast_compression` - (Optional)  Specifies the default TOAST compression method for values of compressible columns (the default is lz4). Supported values are: `lz4`, `pglz`.
* `idle_in_transaction_session_timeout` - (Optional)  Time out sessions with open transactions after this number of milliseconds
* `jit` - (Optional)  Activates, in a boolean, the system-wide use of Just-in-Time Compilation (JIT).
* `log_autovacuum_min_duration` - (Optional)  Causes each action executed by autovacuum to be logged if it ran for at least the specified number of milliseconds. Setting this to zero logs all autovacuum actions. Minus-one (the default) disables logging autovacuum actions.
* `log_error_verbosity` - (Optional)  Controls the amount of detail written in the server log for each message that is logged. Supported values are: `TERSE`, `DEFAULT`, `VERBOSE`.
* `log_line_prefix` - (Optional)  Selects one of the available log-formats. These can support popular log analyzers like pgbadger, pganalyze, etc. Supported values are: `pid=%p,user=%u,db=%d,app=%a,client=%h`, `%m [%p] %q[user=%u,db=%d,app=%a]`, `%t [%p]: [%l-1] user=%u,db=%d,app=%a,client=%h`.
* `log_min_duration_statement` - (Optional)  Log statements that take more than this number of milliseconds to run. If -1, disables.
* `max_files_per_process` - (Optional)  PostgreSQL maximum number of files that can be open per process.
* `max_prepared_transactions` - (Optional)  PostgreSQL maximum prepared transactions. Once increased, this parameter cannot be lowered from its set value.
* `max_pred_locks_per_transaction` - (Optional)  PostgreSQL maximum predicate locks per transaction.
* `max_locks_per_transaction` - (Optional)  PostgreSQL maximum locks per transaction. Once increased, this parameter cannot be lowered from its set value.
* `max_stack_depth` - (Optional)  Maximum depth of the stack in bytes.
* `max_standby_archive_delay` - (Optional)  Max standby archive delay in milliseconds.
* `max_standby_streaming_delay` - (Optional)  Max standby streaming delay in milliseconds.
* `max_replication_slots` - (Optional)  PostgreSQL maximum replication slots.
* `max_logical_replication_workers` - (Optional)  PostgreSQL maximum logical replication workers (taken from the pool of max_parallel_workers).
* `max_parallel_workers` - (Optional)  Sets the maximum number of workers that the system can support for parallel queries.
* `max_parallel_workers_per_gather` - (Optional)  Sets the maximum number of workers that can be started by a single Gather or Gather Merge node.
* `max_worker_processes` - (Optional)  Sets the maximum number of background processes that the system can support. Once increased, this parameter cannot be lowered from its set value.
* `pg_partman_bgw_role` - (Optional)  Controls which role to use for pg_partman's scheduled background tasks. Must consist of alpha-numeric characters, dots, underscores, or dashes. May not start with dash or dot. Maximum of 64 characters.
* `pg_partman_bgw_interval` - (Optional)  Sets the time interval to run pg_partman's scheduled tasks.
* `pg_stat_statements_track` - (Optional)  Controls which statements are counted. Specify 'top' to track top-level statements (those issued directly by clients), 'all' to also track nested statements (such as statements invoked within functions), or 'none' to disable statement statistics collection. The default value is top. Supported values are: `all`, `top`, `none`.
* `temp_file_limit` - (Optional)  PostgreSQL temporary file limit in KiB. If -1, sets to unlimited.
* `timezone` - (Optional)  PostgreSQL service timezone
* `track_activity_query_size` - (Optional)  Specifies the number of bytes reserved to track the currently executing command for each active session.
* `track_commit_timestamp` - (Optional)  Record commit time of transactions. The default value is top. Supported values are: `off`, `on`.
* `track_functions` - (Optional)  Enables tracking of function call counts and time used. The default value is top. Supported values are: `all`, `pl`, `none`.
* `track_io_timing` - (Optional)  Enables timing of database I/O calls. This parameter is off by default, because it will repeatedly query the operating system for the current time, which may cause significant overhead on some platforms. The default value is top. Supported values are: `off`, `on`.
* `max_wal_senders` - (Optional)  PostgreSQL maximum WAL senders. Once increased, this parameter cannot be lowered from its set value.
* `wal_sender_timeout` - (Optional)  Terminate replication connections that are inactive for longer than this amount of time, in milliseconds. Setting this value to zero disables the timeout. Must be either 0 or between 5000 and 10800000.
* `wal_writer_delay` - (Optional)  WAL flush interval in milliseconds. Note that setting this value to lower than the default 200ms may negatively impact performance
* `shared_buffers_percentage` - (Optional)  Percentage of total RAM that the database server uses for shared memory buffers. Valid range is 20-60 (float), which corresponds to 20% - 60%. This setting adjusts the shared_buffers configuration value.
* `pgbouncer` - (Optional)  PGBouncer connection pooling settings
* `backup_hour` - (Optional)  The hour of day (in UTC) when backup for the service starts. New backup only starts if previous backup has already completed.
* `backup_minute` - (Optional)  The minute of the backup hour when backup for the service starts. New backup is only started if previous backup has already completed.
* `work_mem` - (Optional)  The maximum amount of memory, in MB, used by a query operation (such as a sort or hash table) before writing to temporary disk files. Default is 1MB + 0.075% of total RAM (up to 32MB).
* `timescaledb` - (Optional)  TimescaleDB extension configuration values

## Attributes Reference

All above attributes are exported. If an attribute was set outside of Terraform, it will be computed.

## Import

A PostgreSQL database cluster's configuration can be imported using the `id` the parent cluster, e.g.

```bash
terraform import digitalocean_database_postgresql_config.example 52556c07-788e-4d41-b8a7-c796432197d1
```
