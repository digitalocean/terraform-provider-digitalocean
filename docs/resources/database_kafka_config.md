---
page_title: "DigitalOcean: digitalocean_database_kafka_config"
---

# digitalocean\_database\_kafka\_config

Provides a virtual resource that can be used to change advanced configuration
options for a DigitalOcean managed Kafka database cluster.

-> **Note** Kafka configurations are only removed from state when destroyed. The remote configuration is not unset.

## Example Usage

```hcl
resource "digitalocean_database_kafka_config" "example" {
  cluster_id                         = digitalocean_database_cluster.example.id
  group_initial_rebalance_delay_ms = 3000
  group_min_session_timeout_ms = 6000
  group_max_session_timeout_ms = 1800000
  message_max_bytes = 1048588
  log_cleaner_delete_retention_ms = 86400000
  log_cleaner_min_compaction_lag_ms = 0
  log_flush_interval_ms = 9223372036854775807
  log_index_interval_bytes = 4096
  log_message_downconversion_enable = true
  log_message_timestamp_difference_max_ms = 9223372036854775807
  log_preallocate = false
  log_retention_bytes = -1
  log_retention_hours = 168
  log_retention_ms = 604800000
  log_roll_jitter_ms = 0
  log_segment_delete_delay_ms = 60000
  auto_create_topics_enable = true
}

resource "digitalocean_database_cluster" "example" {
  name       = "example-kafka-cluster"
  engine     = "kafka"
  version    = "3.7"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc3"
  node_count = 3
}
```


## Argument Reference

The following arguments are supported. See the [DigitalOcean API documentation](https://docs.digitalocean.com/reference/api/api-reference/#operation/databases_patch_config)
for additional details on each option.

* `cluster_id` - (Required)  The ID of the target Kafka cluster.
* `group_initial_rebalance_delay_ms` - (Optional) The amount of time, in milliseconds, the group coordinator will wait for more consumers to join a new group before performing the first rebalance. A longer delay means potentially fewer rebalances, but increases the time until processing begins. The default value for this is 3 seconds. During development and testing it might be desirable to set this to 0 in order to not delay test execution time.
* `group_min_session_timeout_ms` - (Optional) The minimum allowed session timeout for registered consumers. Longer timeouts give consumers more time to process messages in between heartbeats at the cost of a longer time to detect failures.
* `group_max_session_timeout_ms` - (Optional) The maximum allowed session timeout for registered consumers. Longer timeouts give consumers more time to process messages in between heartbeats at the cost of a longer time to detect failures.
* `message_max_bytes` - (Optional) The maximum size of message that the server can receive.
* `log_cleaner_delete_retention_ms` - (Optional) How long are delete records retained?
* `log_cleaner_min_compaction_lag_ms` - (Optional) The minimum time a message will remain uncompacted in the log. Only applicable for logs that are being compacted.
* `log_flush_interval_ms` - (Optional) The maximum time in ms that a message in any topic is kept in memory before flushed to disk. If not set, the value in log.flush.scheduler.interval.ms is used.
* `log_index_interval_bytes` - (Optional) The interval with which Kafka adds an entry to the offset index.
* `log_message_downconversion_enable` - (Optional) This configuration controls whether down-conversion of message formats is enabled to satisfy consume requests.
* `log_message_timestamp_difference_max_ms` - (Optional) The maximum difference allowed between the timestamp when a broker receives a message and the timestamp specified in the message.
* `log_preallocate` - (Optional) Controls whether to preallocate a file when creating a new segment.
* `log_retention_bytes` - (Optional) The maximum size of the log before deleting messages.
* `log_retention_hours` - (Optional) The number of hours to keep a log file before deleting it.
* `log_retention_ms` - (Optional) The number of milliseconds to keep a log file before deleting it (in milliseconds), If not set, the value in log.retention.minutes is used. If set to -1, no time limit is applied.
* `log_roll_jitter_ms` - (Optional) The maximum jitter to subtract from logRollTimeMillis (in milliseconds). If not set, the value in log.roll.jitter.hours is used.
* `log_segment_delete_delay_ms` - (Optional) The amount of time to wait before deleting a file from the filesystem.
* `auto_create_topics_enable` - (Optional) Enable auto creation of topics.

## Attributes Reference

All above attributes are exported. If an attribute was set outside of Terraform, it will be computed.

## Import

A Kafka database cluster's configuration can be imported using the `id` the parent cluster, e.g.

```
terraform import digitalocean_database_kafka_config.example 4b62829a-9c42-465b-aaa3-84051048e712
```
