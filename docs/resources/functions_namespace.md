---
page_title: "DigitalOcean: digitalocean_functions_namespace"
subcategory: "Functions"
---

# digitalocean_functions_namespace

Provides a DigitalOcean Functions namespace resource. A namespace is the
container that functions are deployed into, and each namespace lives in a
single region.

Functions are deployed into a namespace with `doctl serverless deploy` (see the
[Functions documentation](https://docs.digitalocean.com/products/functions/));
this resource manages the namespace itself so it can be created and wired up
alongside the rest of your infrastructure.

## Example Usage

```hcl
resource "digitalocean_functions_namespace" "example" {
  label  = "my-functions"
  region = "nyc1"
}

output "functions_api_host" {
  value = digitalocean_functions_namespace.example.api_host
}
```

## Argument Reference

The following arguments are supported:

* `label` - (Required) The namespace's name. Changing this creates a new namespace.
* `region` - (Required) The region the namespace is deployed in, e.g. `nyc1`. Changing this creates a new namespace.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The ID of the namespace, e.g. `fn-b90faf52-2b42-49c2-9792-75edfbb6f397`.
* `uuid` - The namespace's UUID.
* `api_host` - The API host used to deploy and invoke functions in the namespace.
* `key` - (Sensitive) The namespace's access key. Together with `api_host` it is used by the serverless tooling.
* `created_at` - The time the namespace was created.
* `updated_at` - The time the namespace was last updated.

Deleting a namespace also deletes all functions and triggers in it.

## Import

Functions namespaces can be imported using their ID, e.g.

```
terraform import digitalocean_functions_namespace.example fn-b90faf52-2b42-49c2-9792-75edfbb6f397
```
