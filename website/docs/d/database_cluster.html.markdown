---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_database_cluster"
sidebar_current: "docs-do-datasource-database-cluster"
description: |-
  Get information on a database cluster resource.
---

# digitalocean\_database\_cluster

Provides information on a DigitalOcean database cluster resource.

## Example Usage

```hcl
# Create a new database cluster
data "digitalocean_database_cluster" "example" {
  name = "example-cluster"
  engine = "pg"
  region = "fra1"
  node_count = "2"
  size = "db-s-1vcpu-1gb"

  maintenance_window {
    day  = "friday"
	hour = "13:00:00"
  }
}

output "database_output" {
  value = data.digitalocean_database_cluster.example.uri
}
```

## Argument Reference

The following arguments are supported (always as String):

* `name` - (Required) The name of the database cluster.
* `engine` - (Required) Database engine used by the cluster (ex. `pg` for PostreSQL).
* `region` - (Required) DigitalOcean region where the cluster will reside. ([Regional Availability Matrix](https://www.digitalocean.com/docs/platform/availability-matrix/))
* `size` - (Required) Database droplet size associated with the cluster (ex. `db-s-1vcpu-1gb`).
* `node_count` - (Required) Number of nodes that will be included in the cluster.
* `version` - Engine version used by the cluster (ex. `11` for PostgreSQL 11).
* `maintenance_window` - Defines when the automatic maintenance should be performed for the database cluster.
* `eviction_policy` - (Redis only) Specify the eviction policy for a Redis cluster.
* `sql_mode` - (MySql only) A comma separated string specifying the  SQL modes for a MySQL cluster.


`maintenance_window` supports the following:

* `day` - The day of the week on which to apply maintenance updates.
* `hour` - The hour in UTC at which maintenance updates will be applied in 24 hour format.

`eviction_policy` supports the following:
* `noeviction`
* `allkeys_lru`
* `allkeys_random`
* `volatile_lru`
* `volatile_random`
* `volatile_ttl`

## Attributes Reference

The following attributes are computed by DigitalOcean:

* `urn` - The uniform resource name of the database cluster.
* `host` - Database cluster's hostname.
* `private_host` - Same as `host`, but only accessible from resources within the account and in the same region.
* `port` - Network port that the database cluster is listening on.
* `uri` - The full URI for connecting to the database cluster.
* `private_uri` - Same as `uri`, but only accessible from resources within the account and in the same region.
* `database` - Name of the cluster's default database.
* `user` - Username for the cluster's default user.
* `password` - Password for the cluster's default user.
