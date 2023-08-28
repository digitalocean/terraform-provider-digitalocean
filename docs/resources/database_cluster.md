---
page_title: "DigitalOcean: digitalocean_database_cluster"
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
  version    = "8"
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
  version    = "6"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}
```

### Create a new MongoDB database cluster
```hcl
resource "digitalocean_database_cluster" "mongodb-example" {
  name       = "example-mongo-cluster"
  engine     = "mongodb"
  version    = "4"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc3"
  node_count = 1
}
```

## Create a new database cluster based on a backup of an existing cluster.
```hcl
resource "digitalocean_database_cluster" "doby" {
  name       = "dobydb"
  engine     = "pg"
  version    = "11"
  size       = "db-s-1vcpu-2gb"
  region     = "nyc1"
  node_count = 1
  tags       = ["production"]
}

resource "digitalocean_database_cluster" "doby_backup" {
  name       = "dobydupe"
  engine     = "pg"
  version    = "11"
  size       = "db-s-1vcpu-2gb"
  region     = "nyc1"
  node_count = 1
  tags       = ["production"]
  
  backup_restore {
    database_name  = "dobydb"
  }

  depends_on = [
    digitalocean_database_cluster.doby
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the database cluster.
* `engine` - (Required) Database engine used by the cluster (ex. `pg` for PostreSQL, `mysql` for MySQL, `redis` for Redis, or `mongodb` for MongoDB).
* `size` - (Required) Database Droplet size associated with the cluster (ex. `db-s-1vcpu-1gb`). See here for a [list of valid size slugs](https://docs.digitalocean.com/reference/api/api-reference/#tag/Databases).
* `region` - (Required) DigitalOcean region where the cluster will reside.
* `node_count` - (Required) Number of nodes that will be included in the cluster.
* `version` - (Required) Engine version used by the cluster (ex. `14` for PostgreSQL 14).
  When this value is changed, a call to the [Upgrade major Version for a Database](https://docs.digitalocean.com/reference/api/api-reference/#operation/databases_update_major_version) API operation is made with the new version.
* `tags` - (Optional) A list of tag names to be applied to the database cluster.
* `private_network_uuid` - (Optional) The ID of the VPC where the database cluster will be located.
* `project_id` - (Optional) The ID of the project that the database cluster is assigned to. If excluded when creating a new database cluster, it will be assigned to your default project.
* `eviction_policy` - (Optional) A string specifying the eviction policy for a Redis cluster. Valid values are: `noeviction`, `allkeys_lru`, `allkeys_random`, `volatile_lru`, `volatile_random`, or `volatile_ttl`.
* `sql_mode` - (Optional) A comma separated string specifying the  SQL modes for a MySQL cluster.
* `maintenance_window` - (Optional) Defines when the automatic maintenance should be performed for the database cluster.
* `number_of_databases` - (Optional) Set number of redis databases. Valid values are between 1 and 128. Default value is 16.

`maintenance_window` supports the following:

* `day` - (Required) The day of the week on which to apply maintenance updates.
* `hour` - (Required) The hour in UTC at which maintenance updates will be applied in 24 hour format.

* `backup_restore` - (Optional) Create a new database cluster based on a backup of an existing cluster.

`backup_restore` supports the following:

* `database_name` - (Required) The name of an existing database cluster from which the backup will be restored.
* `backup_created_at` - (Optional) The timestamp of an existing database cluster backup in ISO8601 combined date and time format. The most recent backup will be used if excluded.

This resource supports [customized create timeouts](https://www.terraform.io/docs/language/resources/syntax.html#operation-timeouts). The default timeout is 30 minutes.

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
