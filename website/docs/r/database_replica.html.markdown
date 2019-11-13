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
  cluster_id = digitalocean_database_cluster.postgres-example.id
  name       = "read-replica"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
}

resource "digitalocean_database_cluster" "postgres-example" {
  name       = "example-postgres-cluster"
  engine     = "pg"
  version    = "11"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The ID of the original source database cluster.
* `name` - (Required) The name for the database replica.
* `size` - (Required) Database Droplet size associated with the replica (ex. `db-s-1vcpu-1gb`).
* `region` - (Required) DigitalOcean region where the replica will reside.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The ID of the database replica.
* `host` - Database replica's hostname.
* `private_host` - Same as `host`, but only accessible from resources within the account and in the same region.
* `port` - Network port that the database replica is listening on.
* `uri` - The full URI for connecting to the database replica.
* `private_uri` - Same as `uri`, but only accessible from resources within the account and in the same region.
* `database` - Name of the replica's default database.
* `user` - Username for the replica's default user.
* `password` - Password for the replica's default user.

## Import

Database replicas can be imported using the `id` of the source database cluster
and the `name` of the replica joined with a comma. For example:

```
terraform import digitalocean_database_replica.read-replica 245bcfd0-7f31-4ce6-a2bc-475a116cca97,read-replica
```
