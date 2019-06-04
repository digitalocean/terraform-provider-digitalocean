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
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the database cluster.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the database cluster.
* `engine` - Database engine used by the cluster (ex. `pg` for PostreSQL).
* `version` - Engine version used by the cluster (ex. `11` for PostgreSQL 11).
* `size` - Database droplet size associated with the cluster (ex. `db-s-1vcpu-1gb`).
* `region` - DigitalOcean region where the cluster will reside.
* `node_count` - Number of nodes that will be included in the cluster.
* `maintenance_window` - Defines when the automatic maintenance should be performed for the database cluster.

* `host` - Database cluster's hostname.
* `port` - Network port that the database cluster is listening on.
* `uri` - The full URI for connecting to the database cluster.
* `database` - Name of the cluster's default database.
* `user` - Username for the cluster's default user.

`maintenance_window` supports the following:

* `day` - The day of the week on which to apply maintenance updates.
* `hour` - The hour in UTC at which maintenance updates will be applied in 24 hour format.

