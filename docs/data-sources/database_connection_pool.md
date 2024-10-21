---
page_title: "DigitalOcean: digitalocean_database_connection_pool"
subcategory: "Databases"
---

# digitalocean\_database\_connection\_pool

Provides information on a DigitalOcean PostgreSQL database connection pool.

## Example Usage

```hcl
data "digitalocean_database_cluster" "example" {
  name = "example-cluster"
}
data "digitalocean_database_connection_pool" "read-only" {
  cluster_id = data.digitalocean_database_cluster.example.id
  name       = "pool-01"
}
output "connection_pool_uri_output" {
  value = data.digitalocean_database_connection_pool.read-only.uri
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The ID of the original source database cluster.
* `name` - (Required) The name for the database connection pool.


## Attributes Reference

The following attributes are exported:

* `id` - The ID of the database connection pool.
* `host` - Connection pool hostname.
* `private_host` - Same as `host`, but only accessible from resources within the account and in the same region.
* `port` - Network port that the connection pool is listening on.
* `uri` - The full URI for connecting to the database connection pool.
* `private_uri` - Same as `uri`, but only accessible from resources within the account and in the same region.
* `db_name` - Name of the connection pool's default database.
* `size` - Size of the connection pool.
* `mode` - The transaction mode for the connection pool. 
* `user` - Username for the connection pool's default user.
* `password` - Password for the connection pool's default user.
