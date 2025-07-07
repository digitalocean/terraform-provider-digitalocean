---
page_title: "DigitalOcean: digitalocean_database_valkey_config"
subcategory: "Databases"
---

# digitalocean\_database\_valkey\_config

Provides a virtual resource that can be used to change advanced configuration
options for a DigitalOcean managed Valkey database cluster.

-> **Note** Valkey configurations are only removed from state when destroyed. The remote configuration is not unset.

## Example Usage

```hcl
resource "digitalocean_database_valkey_config" "example" {
  cluster_id             = digitalocean_database_cluster.example.id
  maxmemory_policy       = "allkeys-lru"
  notify_keyspace_events = "KEA"
  timeout                = 90
}

resource "digitalocean_database_cluster" "example" {
  name       = "example-valkey-cluster"
  engine     = "valkey"
  version    = "8"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}
```


## Argument Reference

The following arguments are supported. See the [DigitalOcean API documentation](https://docs.digitalocean.com/reference/api/digitalocean/#tag/Databases/operation/databases_patch_config)
for additional details on each option. 


* `cluster_id` - (Required)  The ID of the target Valkey cluster.
* `maxmemory_policy` - (Optional) A string specifying the desired eviction policy for the Valkey cluster.Supported values are: `noeviction`, `allkeys-lru`, `allkeys-random`, `volatile-lru`, `volatile-random`, `volatile-ttl`
* `pubsub_client_output_buffer_limit` - (Optional) The output buffer limit for pub/sub clients in MB. The value is the hard limit, the soft limit is 1/4 of the hard limit. When setting the limit, be mindful of the available memory in the selected service plan.
* `number_of_databases` - (Optional) The number of Valkey databases. Changing this will cause a restart of Valkey service.
* `io_threads` - (Optional) The Valkey IO thread count.
* `lfu_log_factor` - (Optional) The counter logarithm factor for volatile-lfu and allkeys-lfu maxmemory policies.
* `lfu_decay_time` - (Optional) The LFU maxmemory policy counter decay time in minutes.
* `ssl` - (Optional) A boolean indicating whether to require SSL to access Valkey.
* `timeout` - (Optional) The Valkey idle connection timeout in seconds.
* `notify_keyspace_events` - (Optional) The `notify-keyspace-events` option. Requires at least `K` or `E`.
* `persistence` - (Optional) When persistence is 'rdb', Valkey does RDB dumps each 10 minutes if any key is changed. Also RDB dumps are done according to backup schedule for backup purposes. When persistence is 'off', no RDB dumps and backups are done, so data can be lost at any moment if service is restarted for any reason, or if service is powered off. Also service can't be forked.
* `acl_channels_default` - (Optional) Determines default pub/sub channels' ACL for new users if an ACL is not supplied. When this option is not defined, `allchannels` is assumed to keep backward compatibility. This option doesn't affect Valkey' `acl-pubsub-default` configuration. Supported values are: `allchannels` and `resetchannels`

## Attributes Reference

All above attributes are exported. If an attribute was set outside of Terraform, it will be computed.

## Import

A Valkey database cluster's configuration can be imported using the `id` the parent cluster, e.g.

```
terraform import digitalocean_database_valkey_config.example 245bcfd0-7f31-4ce6-a2bc-475a116cca97
```
