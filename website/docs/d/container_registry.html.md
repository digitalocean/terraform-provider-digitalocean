---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_container_registry"
sidebar_current: "docs-do-datasource-container-registry"
description: |-
  Get information on a container registry.
---

# digitalocean_container_registry

-> **Note**: DigitalOcean Container Registry is currently in private beta.
This feature will become available to the general public soon.

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

### Docker Provider Example

Use the `endpoint` and `docker_credentials` with the Docker provider:

```
data "digitalocean_container_registry" "example" {
  name = "example"
}

provider "docker" {
  host = "tcp:localhost:2376"

  registry_auth {
    address = data.digitalocean_container_registry.example.server_url
    config_file_content = data.digitalocean_container_registry.example.docker_credentials
  }
}
```

### Kubernetes Example

Combined with the Kubernetes Provider's `kubernetes_secret` resource, you can 
access the registry from inside your cluster:

```
data "digitalocean_container_registry" "example" {
  name = "example"
}

resource "digitalocean_kubernetes_cluster" "example" {
  name    = "example"
  region  = "nyc1"
  # Grab the latest version slug from `doctl kubernetes options versions`
  version = "1.16.6-do.0"

  node_pool {
    name       = "worker-pool"
    size       = "s-2vcpu-2gb"
    node_count = 3
  }
}

provider "kubernetes" {
  load_config_file = false
  host  = digitalocean_kubernetes_cluster.example.endpoint
  token = digitalocean_kubernetes_cluster.example.kube_config[0].token
  cluster_ca_certificate = base64decode(
    digitalocean_kubernetes_cluster.example.kube_config[0].cluster_ca_certificate
  )
}

resource "kubernetes_secret" "example" {
  metadata {
    name = "docker-cfg"
  }

  data = {
    ".dockerconfigjson" = data.digitalocean_container_registry.example.docker_credentials
  }

  type = "kubernetes.io/dockerconfigjson"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the container registry.
* `write` - (Optional) Boolean  to retrieve read/write credentials, suitable for use with the Docker client or in a CI system. Defaults to false.

## Attributes Reference

The following attributes are exported:

* `id`: The ID of the tag. This is the same as the name.
* `name` - The name of the container registry
* `endpoint`: The URL endpoint of the container registry. Ex: `registry.digitalocean.com/my_registry`
* `server_url`: The domain of the container registry. Ex: `registry.digitalocean.com`
* `docker_credentials`: Credentials for the container registry.
