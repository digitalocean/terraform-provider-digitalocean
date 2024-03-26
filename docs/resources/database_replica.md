---
page_title: "DigitalOcean: digitalocean_database_replica"
---

# digitalocean\_database\_replica

Provides a DigitalOcean database replica resource.

## Example Usage

### Create a new PostgreSQL database replica
```hcl
output "UUID" {
  value = digitalocean_database_replica.replica-example.uuid
}

resource "digitalocean_database_cluster" "postgres-example" {
  name       = "example-postgres-cluster"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

resource "digitalocean_database_replica" "replica-example" {
  cluster_id = digitalocean_database_cluster.postgres-example.id
  name       = "replica-example"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
}

# Create firewall rule for database replica
resource "digitalocean_database_firewall" "example-fw" {
  cluster_id = digitalocean_database_replica.replica-example.uuid

  rule {
    type  = "ip_addr"
    value = "192.168.1.1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The ID of the original source database cluster.
* `name` - (Required) The name for the database replica.
* `size` - (Required) Database Droplet size associated with the replica (ex. `db-s-1vcpu-1gb`). Note that when resizing an existing replica, its size can only be increased. Decreasing its size is not supported.
* `region` - (Required) DigitalOcean region where the replica will reside.
* `tags` - (Optional) A list of tag names to be applied to the database replica.
* `private_network_uuid` - (Optional) The ID of the VPC where the database replica will be located.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The ID of the database replica created by Terraform.
* `uuid` - The UUID of the database replica. The uuid can be used to reference the database replica as the target database cluster in other resources. See example  "Create firewall rule for database replica" above.
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
