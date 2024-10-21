---
page_title: "DigitalOcean: digitalocean_database_user"
subcategory: "Databases"
---

# digitalocean\_database\_user

Provides information on a DigitalOcean database user resource.

## Example Usage

```hcl
data "digitalocean_database_cluster" "main" {
  name = "main-cluster"
}

data "digitalocean_database_user" "example" {
  cluster_id = data.digitalocean_database_cluster.main.id
  name       = "example-user"
}

output "database_user_password" {
  value     = data.digitalocean_database_user.example.password
  sensitive = true
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The ID of the database cluster.
* `name` - (Required) The name of the database user.

## Attributes Reference

The following attributes are exported:

* `role` - The role of the database user. The value will be either `primary` or `normal`.
* `password` - The password of the database user. This will not be set for MongoDB users.
* `access_cert` - Access certificate for TLS client authentication. (Kafka only)
* `access_key` - Access key for TLS client authentication. (Kafka only)
* `mysql_auth_plugin` - The authentication method of the MySQL user. The value will be `mysql_native_password` or `caching_sha2_password`.
