---
page_title: "DigitalOcean: digitalocean_vector_database"
subcategory: "Databases"
---

# digitalocean\_vector\_database

Provides a DigitalOcean vector database resource. Vector databases are powered by
[Weaviate](https://weaviate.io/) and are managed independently from standard
managed database clusters.

## Example Usage

### Create a new vector database
```hcl
resource "digitalocean_vector_database" "example" {
  name   = "example-vector-db"
  region = "nyc1"
  size   = "db-s-1vcpu-1gb"
  tags   = ["production"]
}
```

### Create a vector database with advanced configuration
```hcl
resource "digitalocean_vector_database" "example" {
  name   = "example-vector-db"
  region = "nyc1"
  size   = "db-s-2vcpu-2gb"

  config {
    default_quantization = "none"
    enable_auto_schema   = true
    weaviate_version     = "1.25.0"
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the vector database. Changing this forces a new resource to be created.
* `region` - (Required) The slug identifier for the region where the vector database will be created (ex. `nyc1`). Changing this forces a new resource to be created.
* `size` - (Required) The slug identifier representing the size of the vector database (ex. `db-s-1vcpu-1gb`).
* `project_id` - (Optional) The ID of the project that the vector database is assigned to. If excluded, the database will be assigned to your default project. Changing this forces a new resource to be created.
* `tags` - (Optional) A list of tag names to be applied to the vector database.
* `config` - (Optional) Advanced configuration for the vector database. The structure is documented below.

`config` supports the following:

* `default_quantization` - (Optional) The default vector quantization method applied to new collections.
* `enable_auto_schema` - (Optional) Whether Weaviate's auto-schema feature is enabled.
* `weaviate_version` - (Optional) The Weaviate engine version used by the vector database.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The ID of the vector database.
* `status` - The current status of the vector database (ex. `active`).
* `owner_uuid` - The UUID of the account that owns the vector database.
* `endpoints` - The connection endpoints for the vector database. The structure is documented below.
* `created_at` - The date and time when the vector database was created.
* `updated_at` - The date and time when the vector database was last updated.

`endpoints` exports the following:

* `http` - The HTTP endpoint used to connect to the vector database.
* `grpc` - The gRPC endpoint used to connect to the vector database.

## Import

Vector databases can be imported using their `id`, e.g.

```
terraform import digitalocean_vector_database.example 245bcfd0-7f31-4ce6-a2bc-475a116cca97
```
