---
page_title: "DigitalOcean: digitalocean_database_cluster"
subcategory: "Databases"
---

# digitalocean\_database\_cluster

Reads a DigitalOcean database cluster's connection credentials as an [ephemeral resource](https://developer.hashicorp.com/terraform/language/resources/ephemeral). The credentials are retrieved on each Terraform run and are never stored in the plan or state.

This resource works with any managed database engine; which attributes are populated depends on the engine.

~> **Note:** Ephemeral resources require Terraform 1.10 or later.

## Example Usage

```hcl
resource "digitalocean_database_cluster" "example" {
  name       = "example-cluster"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}

ephemeral "digitalocean_database_cluster" "example" {
  id = digitalocean_database_cluster.example.id
}

provider "postgresql" {
  host     = ephemeral.digitalocean_database_cluster.example.host
  port     = ephemeral.digitalocean_database_cluster.example.port
  database = ephemeral.digitalocean_database_cluster.example.database
  username = ephemeral.digitalocean_database_cluster.example.user
  password = ephemeral.digitalocean_database_cluster.example.password
  sslmode  = "require"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the database cluster.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the database cluster.
* `engine` - Database engine used by the cluster (ex. `pg` for PostgreSQL).
* `version` - Engine version used by the cluster.
* `region` - DigitalOcean region where the cluster is running.
* `host` - Database cluster's hostname.
* `private_host` - Same as `host`, but only accessible from resources within the account and in the same region.
* `port` - Network port that the database cluster is listening on.
* `database` - Name of the cluster's default database.
* `user` - Username for the cluster's default user.
* `password` - Password for the cluster's default user.
* `uri` - The full URI for connecting to the database cluster.
* `private_uri` - Same as `uri`, but only accessible from resources within the account and in the same region.

OpenSearch clusters will have the following additional attributes with connection details for their dashboard:

* `ui_host` - Hostname for the OpenSearch dashboard.
* `ui_port` - Network port that the OpenSearch dashboard is listening on.
* `ui_uri` - The full URI for connecting to the OpenSearch dashboard.
* `ui_database` - Name of the OpenSearch dashboard db.
* `ui_user` - Username for OpenSearch dashboard's default user.
* `ui_password` - Password for the OpenSearch dashboard's default user.

~> **Note:** MongoDB cluster passwords are only returned when the cluster is created. Reading an existing MongoDB cluster returns no password, so this ephemeral resource cannot retrieve it. The other connection fields are still returned.

~> **Note:** When the provider is configured with a read-only DigitalOcean API token, the API redacts database passwords. In that case `password` will be null and `uri` or `private_uri` may not contain a usable password.

~> **Note:** This ephemeral resource does not store the values it returns, but the `digitalocean_database_cluster` resource and data source still store their own copies of the cluster credentials in state.

Use the [`digitalocean_database_cluster` data source](/providers/digitalocean/digitalocean/latest/docs/data-sources/database_cluster) for non-sensitive metadata.
