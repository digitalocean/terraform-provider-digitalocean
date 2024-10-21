---
page_title: "DigitalOcean: digitalocean_database_mysql_config"
subcategory: "Databases"
---

# digitalocean\_database\_mysql\_config

Provides a virtual resource that can be used to change advanced configuration
options for a DigitalOcean managed MySQL database cluster.

-> **Note** MySQL configurations are only removed from state when destroyed. The remote configuration is not unset.

## Example Usage

```hcl
resource "digitalocean_database_mysql_config" "example" {
  cluster_id        = digitalocean_database_cluster.example.id
  connect_timeout   = 10
  default_time_zone = "UTC"
}

resource "digitalocean_database_cluster" "example" {
  name       = "example-mysql-cluster"
  engine     = "mysql"
  version    = "8"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}
```


## Argument Reference

The following arguments are supported. See the [DigitalOcean API documentation](https://docs.digitalocean.com/reference/api/api-reference/#operation/databases_patch_config)
for additional details on each option.

* `cluster_id` - (Required)  The ID of the target MySQL cluster.
* `connect_timeout` - (Optional) The number of seconds that the mysqld server waits for a connect packet before responding with bad handshake.
* `default_time_zone` - (Optional) Default server time zone, in the form of an offset from UTC (from -12:00 to +12:00), a time zone name (EST), or `SYSTEM` to use the MySQL server default.
* `innodb_log_buffer_size` - (Optional) The size of the buffer, in bytes, that InnoDB uses to write to the log files. on disk.
* `innodb_online_alter_log_max_size` - (Optional) The upper limit, in bytes, of the size of the temporary log files used during online DDL operations for InnoDB tables.
* `innodb_lock_wait_timeout` - (Optional) The time, in seconds, that an InnoDB transaction waits for a row lock. before giving up.
* `interactive_timeout` - (Optional) The time, in seconds, the server waits for activity on an interactive. connection before closing it.
* `max_allowed_packet` - (Optional) The size of the largest message, in bytes, that can be received by the server. Default is `67108864` (64M).
* `net_read_timeout` - (Optional) The time, in seconds, to wait for more data from an existing connection. aborting the read.
* `sort_buffer_size` - (Optional) The sort buffer size, in bytes, for `ORDER BY` optimization. Default is `262144`. (256K).
* `sql_mode` - (Optional) Global SQL mode. If empty, uses MySQL server defaults. Must only include uppercase alphabetic characters, underscores, and commas.
* `sql_require_primary_key` - (Optional) Require primary key to be defined for new tables or old tables modified with ALTER TABLE and fail if missing. It is recommended to always have primary keys because various functionality may break if any large table is missing them.
* `wait_timeout` - (Optional) The number of seconds the server waits for activity on a noninteractive connection before closing it.
* `net_write_timeout` - (Optional) The number of seconds to wait for a block to be written to a connection before aborting the write.
* `group_concat_max_len` - (Optional) The maximum permitted result length, in bytes, for the `GROUP_CONCAT()` function.
* `information_schema_stats_expiry` - (Optional) The time, in seconds, before cached statistics expire.
* `innodb_ft_min_token_size` - (Optional) The minimum length of words that an InnoDB FULLTEXT index stores.
* `innodb_ft_server_stopword_table` - (Optional) The InnoDB FULLTEXT index stopword list for all InnoDB tables.
* `innodb_print_all_deadlocks` - (Optional) When enabled, records information about all deadlocks in InnoDB user transactions in the error log. Disabled by default.
* `innodb_rollback_on_timeout` - (Optional) When enabled, transaction timeouts cause InnoDB to abort and roll back the entire transaction.
* `internal_tmp_mem_storage_engine` - (Optional) The storage engine for in-memory internal temporary tables. Supported values are: `TempTable`, `MEMORY`.
* `max_heap_table_size` - (Optional) The maximum size, in bytes, of internal in-memory tables. Also set `tmp_table_size`. Default is `16777216` (16M)
* `tmp_table_size` - (Optional) The maximum size, in bytes, of internal in-memory tables. Also set `max_heap_table_size`. Default is `16777216` (16M).
* `slow_query_log` - (Optional) When enabled, captures slow queries. When disabled, also truncates the mysql.slow_log table. Default is false.
* `long_query_time` - (Optional) The time, in seconds, for a query to take to execute before being captured by `slow_query_logs`. Default is `10` seconds.
* `backup_hour` - (Optional) The hour of day (in UTC) when backup for the service starts. New backup only starts if previous backup has already completed.
* `backup_minute` - (Optional) The minute of the backup hour when backup for the service starts. New backup only starts if previous backup has already completed.
* `binlog_retention_period` - (Optional) The minimum amount of time, in seconds, to keep binlog entries before deletion. This may be extended for services that require binlog entries for longer than the default, for example if using the MySQL Debezium Kafka connector.

## Attributes Reference

All above attributes are exported. If an attribute was set outside of Terraform, it will be computed.

## Import

A MySQL database cluster's configuration can be imported using the `id` the parent cluster, e.g.

```
terraform import digitalocean_database_mysql_config.example 4b62829a-9c42-465b-aaa3-84051048e712
```
