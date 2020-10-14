---
page_title: "DigitalOcean: digitalocean_container_registry"
---

# digitalocean\_container_registry

Provides a DigitalOcean Container Registry resource. A Container Registry is
a secure, private location to store your containers for rapid deployment. 

## Example Usage

```hcl
# Create a new container registry
resource "digitalocean_container_registry" "foobar" {
  name = "foobar"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the container_registry

## Attributes Reference

The following attributes are exported:

* `id` - The id of the container registry
* `name` - The name of the container registry
* `endpoint`: The URL endpoint of the container registry. Ex: `registry.digitalocean.com/my_registry`
* `server_url`: The domain of the container registry. Ex: `registry.digitalocean.com`


## Import

Container Registries can be imported using the `name`, e.g.

```
terraform import digitalocean_container_registry.myregistry registryname
```
