---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_container_registry"
sidebar_current: "docs-do-datasource-container-registry"
description: |-
  Get information on a container registry.
---

# digitalocean_container_registry

Get information on a container registry. This data source provides the name as 
configured on your DigitalOcean account. This is useful if the container 
registry name in question is not managed by Terraform or you need validate if 
the container registry exists in the account.

An error is triggered if the provided container registry name does not exist.

## Example Usage

### Basic Example

Get the container registry:

```hcl
data "digitalocean_container_registry" "example" {
  name = "example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the container registry.

## Attributes Reference

The following attributes are exported:

* `id`: The ID of the tag. This is the same as the name.
* `name` - The name of the container registry
* `endpoint`: The URL endpoint of the container registry. Ex: `registry.digitalocean.com/my_registry`
* `server_url`: The domain of the container registry. Ex: `registry.digitalocean.com`
