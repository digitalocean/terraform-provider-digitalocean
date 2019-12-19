---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_database_cluster"
sidebar_current: "docs-do-resource-database-cluster"
description: |-
  Provides a DigitalOcean database cluster resource.
---

# digitalocean\_database\_cluster

Provides a DigitalOcean database cluster resource.

## Example Usage

### Create a new PostgreSQL database cluster
```hcl
resource "digitalocean_database_cluster" "postgres-example" {
  name       = "example-postgres-cluster"
  engine     = "pg"
  version    = "11"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}
```

### Create a new MySQL database cluster
```hcl
resource "digitalocean_database_cluster" "mysql-example" {
  name       = "example-mysql-cluster"
  engine     = "mysql"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}
```

### Create a new Redis database cluster
```hcl
resource "digitalocean_database_cluster" "redis-example" {
  name       = "example-redis-cluster"
  engine     = "redis"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the database cluster.
* `engine` - (Required) Database engine used by the cluster (ex. `pg` for PostreSQL, `mysql` for MySQL, or `redis` for Redis).
* `size` - (Required) Database Droplet size associated with the cluster (ex. `db-s-1vcpu-1gb`).
* `region` - (Required) DigitalOcean region where the cluster will reside.
* `node_count` - (Required) Number of nodes that will be included in the cluster.
* `version` - (Optional) Engine version used by the cluster (ex. `11` for PostgreSQL 11).
* `tags` - (Optional) A list of tag names to be applied to the database cluster.
* `eviction_policy` - (Optional) A string specifying the eviction policy for a Redis cluster. Valid values are: `noeviction`, `allkeys_lru`, `allkeys_random`, `volatile_lru`, `volatile_random`, or `volatile_ttl`.
* `sql_mode` - (Optional) A comma separated string specifying the  SQL modes for a MySQL cluster.
* `maintenance_window` - (Optional) Defines when the automatic maintenance should be performed for the database cluster.

`maintenance_window` supports the following:

* `day` - (Required) The day of the week on which to apply maintenance updates.
* `hour` - (Required) The hour in UTC at which maintenance updates will be applied in 24 hour format.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The ID of the database cluster.
* `urn` - The uniform resource name of the database cluster.
* `host` - Database cluster's hostname.
* `private_host` - Same as `host`, but only accessible from resources within the account and in the same region.
* `port` - Network port that the database cluster is listening on.
* `uri` - The full URI for connecting to the database cluster.
* `private_uri` - Same as `uri`, but only accessible from resources within the account and in the same region.
* `database` - Name of the cluster's default database.
* `user` - Username for the cluster's default user.
* `password` - Password for the cluster's default user.

## Import

Database clusters can be imported using the `id` returned from DigitalOcean, e.g.

```
terraform import digitalocean_database_cluster.mycluster 245bcfd0-7f31-4ce6-a2bc-475a116cca97
```
