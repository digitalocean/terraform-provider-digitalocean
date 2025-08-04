---
page_title: "DigitalOcean: digitalocean_database_redis_config"
subcategory: "Databases"
---

# digitalocean\_database\_redis\_config

Provides a virtual resource that can be used to change advanced configuration
options for a DigitalOcean managed Redis database cluster.

-> **Note** DigitalOcean managed Redis cluster product is discontinued as of 30 June 2025 and is replaced by the Managed Valkey product. Use the `digitalocean_database_valkey_config` resource instead of `digitalocean_database_redis_config`

-> **Note** Redis configurations are only removed from state when destroyed. The remote configuration is not unset.

## Example Usage

```hcl
resource "digitalocean_database_redis_config" "example" {
  cluster_id             = digitalocean_database_cluster.example.id
  maxmemory_policy       = "allkeys-lru"
  notify_keyspace_events = "KEA"
  timeout                = 90
}

resource "digitalocean_database_cluster" "example" {
  name       = "example-redis-cluster"
  engine     = "redis"
  version    = "7"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}
```


## Argument Reference

The following arguments are supported. See the [DigitalOcean API documentation](https://docs.digitalocean.com/reference/api/digitalocean/#tag/Databases/operation/databases_patch_config)
for additional details on each option. 


* `cluster_id` - (Required)  The ID of the target Redis cluster.
* `maxmemory_policy` - (Optional) A string specifying the desired eviction policy for the Redis cluster.Supported values are: `noeviction`, `allkeys-lru`, `allkeys-random`, `volatile-lru`, `volatile-random`, `volatile-ttl`
* `pubsub_client_output_buffer_limit` - (Optional) The output buffer limit for pub/sub clients in MB. The value is the hard limit, the soft limit is 1/4 of the hard limit. When setting the limit, be mindful of the available memory in the selected service plan.
* `number_of_databases` - (Optional) The number of Redis databases. Changing this will cause a restart of Redis service.
* `io_threads` - (Optional) The Redis IO thread count.
* `lfu_log_factor` - (Optional) The counter logarithm factor for volatile-lfu and allkeys-lfu maxmemory policies.
* `lfu_decay_time` - (Optional) The LFU maxmemory policy counter decay time in minutes.
* `ssl` - (Optional) A boolean indicating whether to require SSL to access Redis.
 - When enabled, Redis accepts only SSL connections on port `25061`.
 - When disabled, port `25060` is opened for non-SSL connections, while port `25061` remains available for SSL connections.
* `timeout` - (Optional) The Redis idle connection timeout in seconds.
* `notify_keyspace_events` - (Optional) The `notify-keyspace-events` option. Requires at least `K` or `E`.
* `persistence` - (Optional) When persistence is `rdb`, Redis does RDB dumps each 10 minutes if any key is changed. Also RDB dumps are done according to backup schedule for backup purposes. When persistence is `off`, no RDB dumps and backups are done, so data can be lost at any moment if service is restarted for any reason, or if service is powered off. Also service can't be forked.
* `acl_channels_default` - (Optional) Determines default pub/sub channels' ACL for new users if an ACL is not supplied. When this option is not defined, `allchannels` is assumed to keep backward compatibility. This option doesn't affect Redis' `acl-pubsub-default` configuration. Supported values are: `allchannels` and `resetchannels`

## Attributes Reference

All above attributes are exported. If an attribute was set outside of Terraform, it will be computed.

## Import

A Redis database cluster's configuration can be imported using the `id` the parent cluster, e.g.

```
terraform import digitalocean_database_redis_config.example 245bcfd0-7f31-4ce6-a2bc-475a116cca97
```
