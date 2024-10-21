---
page_title: "DigitalOcean: digitalocean_container_registry"
subcategory: "Container Registry"
---

# digitalocean\_container_registry

Provides a DigitalOcean Container Registry resource. A Container Registry is
a secure, private location to store your containers for rapid deployment.

## Example Usage

```hcl
# Create a new container registry
resource "digitalocean_container_registry" "foobar" {
  name                   = "foobar"
  subscription_tier_slug = "starter"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the container_registry
* `subscription_tier_slug` - (Required) The slug identifier for the subscription tier to use (`starter`, `basic`, or `professional`)
* `region` - (Optional) The slug identifier of for region where registry data will be stored. When not provided, a region will be selected automatically.

## Attributes Reference

The following attributes are exported:

* `id` - The id of the container registry
* `name` - The name of the container registry
* `subscription_tier_slug` - The slug identifier for the subscription tier
* `region` - The slug identifier for the  region
* `endpoint` - The URL endpoint of the container registry. Ex: `registry.digitalocean.com/my_registry`
* `server_url` - The domain of the container registry. Ex: `registry.digitalocean.com`
* `storage_usage_bytes` - The amount of storage used in the registry in bytes.
* `created_at` - The date and time when the registry was created


## Import

Container Registries can be imported using the `name`, e.g.

```
terraform import digitalocean_container_registry.myregistry registryname
```
