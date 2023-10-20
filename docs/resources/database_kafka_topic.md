---
page_title: "DigitalOcean: digitalocean_database_kafka_topic"
---

# digitalocean\_database\_kafka\_topic

Provides a DigitalOcean Kafka topic for Kafka clusters.

## Example Usage

### Create a new Kafka topic
```hcl
resource "digitalocean_database_kafka_topic" "topic-01" {
  cluster_id            = digitalocean_database_cluster.kafka-example.id
  name                  = "topic-01"
  partition_count       = 3
  replication_factor    = 2
  config {
    cleanup_policy                      = "compact"
    compression_type                    = "uncompressed"
    delete_retention_ms                 = 14000
    file_delete_delay_ms                = 170000
    flush_messages                      = 92233
    flush_ms                            = 92233720368
    index_interval_bytes                = 40962
    max_compaction_lag_ms               = 9223372036854775807
    max_message_bytes                   = 1048588
    message_down_conversion_enable      = true
    message_format_version              = "3.0-IV1"
    message_timestamp_difference_max_ms = 9223372036854775807
    message_timestamp_type              = "log_append_time"
    min_cleanable_dirty_ratio           = 0.5
    min_compaction_lag_ms               = 20000
    min_insync_replicas                 = 2
    preallocate                         = false
    retention_bytes                     = -1
    retention_ms                        = -1
    segment_bytes                       = 209715200
    segment_index_bytes                 = 10485760
    segment_jitter_ms                   = 0
    segment_ms                          = 604800000
    unclean_leader_election_enable      = true
  }
}

resource "digitalocean_database_cluster" "kafka-example" {
  name       = "example-kafka-cluster"
  engine     = "kafka"
  version    = "3.5"
  size       = "db-s-1vcpu-2gb"
  region     = "nyc1"
  node_count = 3
  tags       = ["production"]
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The ID of the source database cluster. Note: This must be a Kafka cluster.
* `name` - (Required) The name for the topic.
* `partition_count` - (Optional) The number of partitions for the topic. Default and minimum set at 3, maximum is 2048.
* `replication_factor` - (Optional) The number of nodes that topics are replicated across. Default and minimum set at 2, maximum is the number of nodes in the cluster.
* `config` - (Optional) A set of advanced configuration parameters. Defaults will be set for any of the parameters that are not included.
  The `config` block is documented below.

`config` supports the following:

* `cleanup_policy` - (Optional) The topic cleanup policy that decribes whether messages should be deleted, compacted, or both when retention policies are violated.
  This may be one of "delete", "compact", or "compact_delete".
* `compression_type` - (Optional) The topic compression codecs used for a given topic.
  This may be one of "uncompressed", "gzip", "snappy", "lz4", "producer", "zstd". "uncompressed" indicates that there is no compression and "producer" retains the original compression codec set by the producer.
* `delete_retention_ms` - (Optional) The amount of time, in ms, that deleted records are retained.
* `file_delete_delay_ms` - (Optional) The amount of time, in ms, to wait before deleting a topic log segment from the filesystem.
* `flush_messages` - (Optional) The number of messages accumulated on a topic partition before they are flushed to disk.
* `flush_ms` - (Optional) The maximum time, in ms, that a topic is kept in memory before being flushed to disk.
* `index_interval_bytes` - (Optional) The interval, in bytes, in which entries are added to the offset index.
* `max_compaction_lag_ms` - (Optional) The maximum time, in ms, that a particular message will remain uncompacted. This will not apply if the `compression_type` is set to "uncompressed" or it is set to `producer` and the producer is not using compression.
* `max_message_bytes` - (Optional) The maximum size, in bytes, of a message.
* `message_down_conversion_enable` - (Optional) Determines whether down-conversion of message formats for consumers is enabled.
* `message_format_version` - (Optional) The version of the inter-broker protocol that will be used. This may be one of "0.8.0", "0.8.1", "0.8.2", "0.9.0", "0.10.0", "0.10.0-IV0", "0.10.0-IV1", "0.10.1", "0.10.1-IV0", "0.10.1-IV1", "0.10.1-IV2", "0.10.2", "0.10.2-IV0", "0.11.0", "0.11.0-IV0", "0.11.0-IV1", "0.11.0-IV2", "1.0", "1.0-IV0", "1.1", "1.1-IV0", "2.0", "2.0-IV0", "2.0-IV1", "2.1", "2.1-IV0", "2.1-IV1", "2.1-IV2", "2.2", "2.2-IV0", "2.2-IV1", "2.3", "2.3-IV0", "2.3-IV1", "2.4", "2.4-IV0", "2.4-IV1", "2.5", "2.5-IV0", "2.6", "2.6-IV0", "2.7", "2.7-IV0", "2.7-IV1", "2.7-IV2", "2.8", "2.8-IV0", "2.8-IV1", "3.0", "3.0-IV0", "3.0-IV1", "3.1", "3.1-IV0", "3.2", "3.2-IV0", "3.3", "3.3-IV0", "3.3-IV1", "3.3-IV2", "3.3-IV3", "3.4", "3.4-IV0", "3.5", "3.5-IV0", "3.5-IV1", "3.5-IV2", "3.6", "3.6-IV0", "3.6-IV1", "3.6-IV2".
* `message_timestamp_difference_max_ms` - (Optional) The maximum difference, in ms, between the timestamp specific in a message and when the broker receives the message.
* `message_timestamp_type` - (Optional) Specifies which timestamp to use for the message. This may be one of "create_time" or "log_append_time".
* `min_cleanable_dirty_ratio` - (Optional) A scale between 0.0 and 1.0 which controls the frequency of the compactor. Larger values mean more frequent compactions. This is often paired with `max_compaction_lag_ms` to control the compactor frequency.
* `min_insync_replicas` - (Optional) The number of replicas that must acknowledge a write before it is considered successful. -1 is a special setting to indicate that all nodes must ack a message before a write is considered successful.
* `preallocate` - (Optional) Determines whether to preallocate a file on disk when creating a new log segment within a topic.
* `retention_bytes` - (Optional) The maximum size, in bytes, of a topic before messages are deleted. -1 is a special setting indicating that this setting has no limit.
* `retention_ms` - (Optional) The maximum time, in ms, that a topic log file is retained before deleting it. -1 is a special setting indicating that this setting has no limit.
* `segment_bytes` - (Optional) The maximum size, in bytes, of a single topic log file.
* `segment_index_bytes` - (Optional) The maximum size, in bytes, of the offset index.
* `segment_jitter_ms` - (Optional) The maximum time, in ms, subtracted from the scheduled segment disk flush time to avoid the thundering herd problem for segment flushing.
* `segment_ms` - (Optional) The maximum time, in ms, before the topic log will flush to disk.
* `unclean_leader_election_enable` - (Optional) Determines whether to allow nodes that are not part of the in-sync replica set (IRS) to be elected as leader. Note: setting this to "true" could result in data loss.



## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `state` - The current status of the topic. Possible values are 'active', 'configuring', and 'deleting'.

## Import

Topics can be imported using the `id` of the source cluster and the `name` of the topic joined with a comma. For example:

```
terraform import digitalocean_database_kafka_topic.topic-01 245bcfd0-7f31-4ce6-a2bc-475a116cca97,topic-01
```
