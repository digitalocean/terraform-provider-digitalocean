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

Get the container registry:

```hcl
data "digitalocean_container_registry" "example" {
  name = "example"
}
```

Use the endpoint and docker_credentials in the docker provider

```
provider "docker" {
  host = "tcp:localhost:2376"

  registry_auth {
    address = "registry.digitalocean.com"
    config_file_content = data.digitalocean_container_registry.registry.docker_credentials
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the container registry.

## Attributes Reference

The following attributes are exported:

* `id`: The ID of the tag. This is the same as the name.
* `endpoint`: The URL endpoint of the container registry.
* `docker_credentials`: Credentials for the container registry.
