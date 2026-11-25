---
page_title: "DigitalOcean: digitalocean_database_user"
subcategory: "Databases"
---

# digitalocean\_database\_user

Reads a DigitalOcean database user's credentials as an [ephemeral resource](https://developer.hashicorp.com/terraform/language/resources/ephemeral). The credentials are retrieved on each Terraform run and are never stored in the plan or state, unlike the `password` attribute of the `digitalocean_database_user` resource and data source.

This resource works with any managed database engine; which attributes are populated depends on the engine.

~> **Note:** Ephemeral resources require Terraform 1.10 or later.

## Example Usage

```hcl
ephemeral "digitalocean_database_user" "example" {
  cluster_id = digitalocean_database_cluster.example.id
  name       = "example-user"
}

provider "postgresql" {
  host     = digitalocean_database_cluster.example.host
  port     = digitalocean_database_cluster.example.port
  username = ephemeral.digitalocean_database_user.example.name
  password = ephemeral.digitalocean_database_user.example.password
  sslmode  = "require"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The ID of the database cluster.
* `name` - (Required) The name of the database user.

## Attributes Reference

The following attributes are exported:

* `role` - The role of the database user. The value will be either `primary` or `normal`.
* `password` - The password of the database user. This is not set for MongoDB users.
* `access_cert` - Access certificate for TLS client authentication. (Kafka only)
* `access_key` - Access key for TLS client authentication. (Kafka only)
* `mysql_auth_plugin` - The authentication method of the MySQL user. The value will be `mysql_native_password` or `caching_sha2_password`.

~> **Note:** MongoDB user passwords are only available when the user is created. Reading an existing MongoDB user returns no password, so this ephemeral resource cannot retrieve it. Other fields the API returns are still populated.

~> **Note:** When the provider is configured with a read-only DigitalOcean API token, the API may also omit the password. In that case `password` will be null.

Use the [`digitalocean_database_user` data source](/providers/digitalocean/digitalocean/latest/docs/data-sources/database_user) for non-sensitive metadata such as Kafka or OpenSearch ACL settings.
