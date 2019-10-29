---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_database_db"
sidebar_current: "docs-do-resource-database-db"
description: |-
  Provides a DigitalOcean database resource.
---

# digitalocean\_database\_db

Provides a DigitalOcean database resource. When creating a new database cluster, a default database with name `defaultdb` will be created. Then, this resource can be used to provide additional database inside the cluster.

## Example Usage

### Create a new PostgreSQL database
```hcl
resource "digitalocean_database_db" "database-example" {
  cluster_id = digitalocean_database_cluster.postgres-example.id
  name       = "foobar"
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
* `name` - (Required) The name for the database.

## Attributes Reference

Only the above arguments are exported.

## Import

Database can be imported using the `id` of the source database cluster
and the `name` of the database joined with a comma. For example:

```
terraform import digitalocean_database_db.database-example 245bcfd0-7f31-4ce6-a2bc-475a116cca97,foobar
```
