---
page_title: "DigitalOcean: digitalocean_database_mongodb_config"
---

# digitalocean\_database\_mongodb\_config

Provides a virtual resource that can be used to change advanced configuration
options for a DigitalOcean managed MongoDB database cluster.

-> **Note** MongoDB configurations are only removed from state when destroyed. The remote configuration is not unset.

## Example Usage

```hcl
resource "digitalocean_database_mongodb_config" "example" {
  cluster_id        = digitalocean_database_cluster.example.id
  default_read_concern               = "majority"
  default_write_concern              = "majority"
  transaction_lifetime_limit_seconds = 100
  slow_op_threshold_ms               = 100
  verbosity                          = 3
}

resource "digitalocean_database_cluster" "example" {
  name       = "example-mongodb-cluster"
  engine     = "mongodb"
  version    = "7"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc3"
  node_count = 1
}
```


## Argument Reference

The following arguments are supported. See the [DigitalOcean API documentation](https://docs.digitalocean.com/reference/api/api-reference/#operation/databases_patch_config)
for additional details on each option.

* `cluster_id` - (Required)  The ID of the target MongoDB cluster.
* `default_read_concern` - (Optional) Specifies the default consistency behavior of reads from the database. Data that is returned from the query with may or may not have been acknowledged by all nodes in the replicaset depending on this value. Learn more [here](https://www.mongodb.com/docs/manual/reference/read-concern/).
* `default_write_concern` - (Optional) Describes the level of acknowledgment requested from MongoDB for write operations clusters. This field can set to either `majority` or a number`0...n` which will describe the number of nodes that must acknowledge the write operation before it is fully accepted. Setting to `0` will request no acknowledgement of the write operation. Learn more [here](https://www.mongodb.com/docs/manual/reference/write-concern/).
* `transaction_lifetime_limit_seconds` - (Optional) Specifies the lifetime of multi-document transactions. Transactions that exceed this limit are considered expired and will be aborted by a periodic cleanup process. The cleanup process runs every `transactionLifetimeLimitSeconds/2 seconds` or at least once every 60 seconds. <em>Changing this parameter will lead to a restart of the MongoDB service.</em> Learn more [here](https://www.mongodb.com/docs/manual/reference/parameters/#mongodb-parameter-param.transactionLifetimeLimitSeconds).
* `slow_op_threshold_ms` - (Optional) Operations that run for longer than this threshold are considered slow which are then recorded to the diagnostic logs. Higher log levels (verbosity) will record all operations regardless of this threshold on the primary node. <em>Changing this parameter will lead to a restart of the MongoDB service.</em> Learn more [here](https://www.mongodb.com/docs/manual/reference/configuration-options/#mongodb-setting-operationProfiling.slowOpThresholdMs).
* `verbosity` - (Optional) The log message verbosity level. The verbosity level determines the amount of Informational and Debug messages MongoDB outputs. 0 includes informational messages while 1...5 increases the level to include debug messages. <em>Changing this parameter will lead to a restart of the MongoDB service.</em> Learn more [here](https://www.mongodb.com/docs/manual/reference/configuration-options/#mongodb-setting-systemLog.verbosity).

## Attributes Reference

All above attributes are exported. If an attribute was set outside of Terraform, it will be computed.

## Import

A MongoDB database cluster's configuration can be imported using the `id` the parent cluster, e.g.

```
terraform import digitalocean_database_mongodb_config.example 4b62829a-9c42-465b-aaa3-84051048e712
```
