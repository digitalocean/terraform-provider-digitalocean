---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_database_cluster"
sidebar_current: "docs-do-resource-database-cluster"
description: |-
  Provides a DigitalOcean database cluster resource.
---

# digitalocean\_database\_cluster

Provides a DigitalOcean database cluster resource.

## Example Usage

```hcl
# Create a new database cluster
resource "digitalocean_database_cluster" "example" {
  name       = "example-cluster"
  engine     = "pg"
  version    = "11"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the database cluster.
* `engine` - (Required) Database engine used by the cluster (ex. `pg` for PostreSQL).
* `version` - (Required) Engine version used by the cluster (ex. `11` for PostgreSQL 11).
* `size` - (Required) Database droplet size associated with the cluster (ex. `db-s-1vcpu-1gb`).
* `region` - (Required) DigitalOcean region where the cluster will reside.
* `node_count` - (Required) Number of nodes that will be included in the cluster.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The ID of the database cluster.

## Import

Database clusters can be imported using the `id` returned from DigitalOcean, e.g.

```
terraform import digitalocean_database_cluster.mycluster 245bcfd0-7f31-4ce6-a2bc-475a116cca97
```
