---
page_title: "DigitalOcean: digitalocean_database_cluster"
---

# digitalocean\_database\_cluster

Provides information on a DigitalOcean database cluster resource.

## Example Usage

```hcl
# Create a new database cluster
data "digitalocean_database_cluster" "example" {
  name = "example-cluster"
}

output "database_output" {
  value = data.digitalocean_database_cluster.example.uri
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the database cluster.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the database cluster.
* `urn` - The uniform resource name of the database cluster.
* `engine` - Database engine used by the cluster (ex. `pg` for PostreSQL).
* `version` - Engine version used by the cluster (ex. `11` for PostgreSQL 11).
* `size` - Database droplet size associated with the cluster (ex. `db-s-1vcpu-1gb`).
* `region` - DigitalOcean region where the cluster will reside.
* `node_count` - Number of nodes that will be included in the cluster.
* `maintenance_window` - Defines when the automatic maintenance should be performed for the database cluster.
* `private_network_uuid` - The ID of the VPC where the database cluster is located.
* `host` - Database cluster's hostname.
* `private_host` - Same as `host`, but only accessible from resources within the account and in the same region.
* `port` - Network port that the database cluster is listening on.
* `uri` - The full URI for connecting to the database cluster.
* `private_uri` - Same as `uri`, but only accessible from resources within the account and in the same region.
* `database` - Name of the cluster's default database.
* `user` - Username for the cluster's default user.
* `password` - Password for the cluster's default user.

`maintenance_window` supports the following:

* `day` - The day of the week on which to apply maintenance updates.
* `hour` - The hour in UTC at which maintenance updates will be applied in 24 hour format.
