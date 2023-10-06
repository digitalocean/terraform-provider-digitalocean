---
page_title: "DigitalOcean: digitalocean_database_user"
---

# digitalocean\_database\_user

Provides a DigitalOcean database user resource. When creating a new database cluster, a default admin user with name `doadmin` will be created. Then, this resource can be used to provide additional normal users inside the cluster.

~> **NOTE:** Any new users created will always have `normal` role, only the default user that comes with database cluster creation has `primary` role. Additional permissions must be managed manually.

## Example Usage

### Create a new PostgreSQL database user
```hcl
resource "digitalocean_database_user" "user-example" {
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

### Create a new user for a PostgreSQL database replica 
```hcl
resource "digitalocean_database_cluster" "postgres-example" {
  name       = "example-postgres-cluster"
  engine     = "pg"
  version    = "11"
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

resource "digitalocean_database_user" "user-example" {
  cluster_id = digitalocean_database_replica.replica-example.uuid
  name       = "foobar"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The ID of the original source database cluster.
* `name` - (Required) The name for the database user.
* `mysql_auth_plugin` - (Optional) The authentication method to use for connections to the MySQL user account. The valid values are `mysql_native_password` or `caching_sha2_password` (this is the default).

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `role` - Role for the database user. The value will be either "primary" or "normal".
* `password` - Password for the database user.

## Import

Database user can be imported using the `id` of the source database cluster
and the `name` of the user joined with a comma. For example:

```
terraform import digitalocean_database_user.user-example 245bcfd0-7f31-4ce6-a2bc-475a116cca97,foobar
```

~> **Note:** MongoDB user passwords are only available when the user is created. An existing MongoDB user that is imported will not have its `password` attribute exported. Recreate the user if it is necessary to access the password with Terraform.
