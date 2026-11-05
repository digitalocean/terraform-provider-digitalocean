---
page_title: "DigitalOcean: digitalocean_vector_database"
subcategory: "Databases"
---

# digitalocean\_vector\_database

Provides information on a DigitalOcean vector database resource.

## Example Usage

```hcl
data "digitalocean_vector_database" "example" {
  name = "example-vector-db"
}

output "vector_db_http_endpoint" {
  value = data.digitalocean_vector_database.example.endpoints[0].http
}
```

A vector database may also be looked up by its `id`:

```hcl
data "digitalocean_vector_database" "example" {
  id = "245bcfd0-7f31-4ce6-a2bc-475a116cca97"
}
```

## Argument Reference

One of the following arguments must be provided:

* `id` - The ID of the vector database.
* `name` - The name of the vector database.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the vector database.
* `name` - The name of the vector database.
* `region` - The slug identifier for the region where the vector database is located.
* `size` - The slug identifier representing the size of the vector database.
* `status` - The current status of the vector database (ex. `active`).
* `owner_uuid` - The UUID of the account that owns the vector database.
* `tags` - A list of tag names applied to the vector database.
* `config` - Advanced configuration for the vector database. The structure is documented below.
* `endpoints` - The connection endpoints for the vector database. The structure is documented below.
* `created_at` - The date and time when the vector database was created.
* `updated_at` - The date and time when the vector database was last updated.

`config` exports the following:

* `default_quantization` - The default vector quantization method applied to new collections.
* `enable_auto_schema` - Whether Weaviate's auto-schema feature is enabled.
* `weaviate_version` - The Weaviate engine version used by the vector database.

`endpoints` exports the following:

* `http` - The HTTP endpoint used to connect to the vector database.
* `grpc` - The gRPC endpoint used to connect to the vector database.
