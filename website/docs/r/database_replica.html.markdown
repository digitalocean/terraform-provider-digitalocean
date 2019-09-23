---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_database_replica"
sidebar_current: "docs-do-resource-database-replica"
description: |-
  Provides a DigitalOcean database replica resource.
---

# digitalocean\_database\_replica

Provides a DigitalOcean database replica resource.

## Example Usage

### Create a new PostgreSQL database replica
```hcl
resource "digitalocean_database_replica" "read-replica" {
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  cluster_id = "${digitalocean_database_cluster.some_db.id}"
}
```

## Argument Reference

The following arguments are supported:

* `size` - (Required) Database droplet size associated with the replica (ex. `db-s-1vcpu-1gb`).
* `region` - (Required) DigitalOcean region where the replica will reside.
* `cluster_id` - (Required) The ID of the original source database cluster.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The ID of the database replica.
* `host` - Database replica's hostname.
* `port` - Network port that the database replica is listening on.
* `uri` - The full URI for connecting to the database replica.
* `database` - Name of the replica's default database.
* `user` - Username for the replica's default user.
* `password` - Password for the replica's default user.

## Import

Database replicas can be imported using the `id` returned from DigitalOcean, e.g.

```
terraform import digitalocean_database_replica.myreplica 245bcfd0-7f31-4ce6-a2bc-475a116cca97
```
