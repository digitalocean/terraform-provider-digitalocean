---
page_title: "DigitalOcean: digitalocean_database_mysql_config"
subcategory: "Databases"
---

# digitalocean\_database\_online\_migration

Provides a virtual resource that can be used to start an online migration 
for a DigitalOcean managed database cluster. Migrating a cluster establishes a 
connection with an existing cluster and replicates its contents to the target 
cluster. If the existing database is continuously being written to, the migration 
process will continue for up to two weeks unless it is manually stopped. 
Online migration is only available for MySQL, PostgreSQL, and Redis clusters.

## Example Usage

```hcl
resource "digitalocean_database_cluster" "source" {
  name       = "st01"
  engine     = "mysql"
  version    = "8"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
  tags       = ["production"]
}

resource "digitalocean_database_cluster" "destination" {
  name       = "dt01"
  engine     = "mysql"
  version    = "8"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
  tags       = ["production"]
}

resource "digitalocean_database_db" "source_db" {
  cluster_id = digitalocean_database_cluster.source.id
  name       = "terraform-db-om-source"
}

resource "digitalocean_database_online_migration" "foobar" {
  cluster_id = digitalocean_database_cluster.destination.id
  source {
    host     = digitalocean_database_cluster.source.host
    db_name  = digitalocean_database_db.source_db.name
    port     = digitalocean_database_cluster.source.port
    username = digitalocean_database_cluster.source.user
    password = digitalocean_database_cluster.source.password
  }
  depends_on = [digitalocean_database_cluster.destination, digitalocean_database_cluster.source, digitalocean_database_db.source_db]
}
```


## Argument Reference

The following arguments are supported. See the [DigitalOcean API documentation](https://docs.digitalocean.com/reference/api/digitalocean/#tag/Databases/operation/databases_update_onlineMigration)
for additional details on each option.

* `cluster_id` - (Required)  The ID of the target MySQL cluster.
* `source` - (Required) Configuration for migration
  * `host` - (Required) The FQDN pointing to the database cluster's current primary node.
  * `port` - (Required) The port on which the database cluster is listening.
  * `dbname` - (Required) The name of the default database.
  * `username` - (Required) The default user for the database.
  * `password` - (Required) A randomly generated password for the default user. 
* `disable_ssl` - (Optional) When set to true, enables SSL encryption when connecting to the source database.
* `ignore_dbs` - (Optional) A list of databases that should be ignored during migration.

## Attributes Reference

All above attributes are exported. If an attribute was set outside of Terraform, it will be computed.

## Import

A MySQL database cluster's online_migration can be imported using the `id` the parent cluster, e.g.

```
terraform import digitalocean_database_online_migration.example 4b62829a-9c42-465b-aaa3-84051048e712
```
