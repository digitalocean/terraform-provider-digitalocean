---
page_title: "DigitalOcean: digitalocean_database_replica"
---

# digitalocean\_database\_replica

Provides information on a DigitalOcean database replica.

## Example Usage

```hcl
data "digitalocean_database_cluster" "example" {
  name = "example-cluster"
}
data "digitalocean_database_replica" "read-only" {
  cluster_id = data.digitalocean_database_cluster.example.id
  name       = "terra-test-ro"
}
output "replica_output" {
  value = data.digitalocean_database_replica.read-only.uri
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The ID of the original source database cluster.
* `name` - (Required) The name for the database replica.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the database replica.
* `host` - Database replica's hostname.
* `private_host` - Same as `host`, but only accessible from resources within the account and in the same region.
* `port` - Network port that the database replica is listening on.
* `uri` - The full URI for connecting to the database replica.
* `private_uri` - Same as `uri`, but only accessible from resources within the account and in the same region.
* `tags` - A list of tag names to be applied to the database replica.
* `database` - Name of the replica's default database.
* `user` - Username for the replica's default user.
* `password` - Password for the replica's default user.
