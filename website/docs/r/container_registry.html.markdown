---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_container_registry"
sidebar_current: "docs-do-resource-container-registry"
description: |-
  Provides a DigitalOcean Tag resource.
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
* `write` - (Optional) Boolean  to retrieve read/write credentials, suitable for use with the Docker client or in a CI system. Defaults to false.
* `expiry_seconds` - (Optional) The number of seconds to pass before expiring the access token. Defaults to 2147483647, which is roughly 68 years.

## Attributes Reference

The following attributes are exported:

* `id` - The id of the container registry
* `name` - The name of the container registry
* `endpoint`: The URL endpoint of the container registry. Ex: `registry.digitalocean.com/my_registry`
* `server_url`: The domain of the container registry. Ex: `registry.digitalocean.com`
* `docker_credentials`: Credentials for the container registry.
* `expiry_seconds`: Number of seconds after creation for token to expire.
* `credential_expiration_time`: The date and time the registry access token will expire.


## Import

Container Registries can be imported using the `name`, e.g.

```
terraform import digitalocean_container_registry.myregistry registryname
```