---
page_title: "DigitalOcean: digitalocean_functions_namespace"
subcategory: "Functions"
---

# digitalocean_functions_namespace

Get information about a DigitalOcean Functions namespace, looked up either by
its ID or by its label. This is useful when the namespace is not managed by
Terraform, for example to reference it from a `digitalocean_functions_trigger`.

## Example Usage

```hcl
data "digitalocean_functions_namespace" "example" {
  label = "my-functions"
}

resource "digitalocean_functions_trigger" "hourly" {
  namespace_id = data.digitalocean_functions_namespace.example.id
  name         = "hourly"
  function     = "sample/hello"

  scheduled_details {
    cron = "0 * * * *"
  }
}
```

## Argument Reference

Exactly one of the following arguments must be provided:

* `namespace_id` - The ID of the namespace, e.g. `fn-b90faf52-2b42-49c2-9792-75edfbb6f397`.
* `label` - The namespace's label. An error is returned if no namespace, or more than one, has this label.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the namespace.
* `namespace_id` - The ID of the namespace.
* `label` - The namespace's label.
* `region` - The region the namespace is deployed in.
* `uuid` - The namespace's UUID.
* `api_host` - The API host used to deploy and invoke functions in the namespace.
* `key` - (Sensitive) The namespace's access key.
* `created_at` - The time the namespace was created.
* `updated_at` - The time the namespace was last updated.
