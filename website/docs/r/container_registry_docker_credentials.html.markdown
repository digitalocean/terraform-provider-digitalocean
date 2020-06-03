---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_container_registry_docker_credentials"
sidebar_current: "docs-do-resource-container-registry-docker-credentials"
description: |-
  Get DigitalOcean Container Registry docker credentials
---

# digitalocean_container_registry_docker_credentials

Get docker credentials for your DigitalOcean container registry.

An error is triggered if the provided container registry name does not exist.

## Example Usage

### Basic Example

Get the container registry:

```hcl
resource "digitalocean_container_registry_docker_credentials" "example" {
  name = "example"
}
```

### Docker Provider Example

Use the `endpoint` and `docker_credentials` with the Docker provider:

```hcl
data "digitalocean_container_registry" "example" {
  name = "example"
}

resource "digitalocean_container_registry_docker_credentials" "example" {
  name = "example"
}

provider "docker" {
  host = "tcp:localhost:2376"

  registry_auth {
    address = data.digitalocean_container_registry.example.server_url
    config_file_content = digitalocean_container_registry_docker_credentials.example.docker_credentials
  }
}
```

### Kubernetes Example

Combined with the Kubernetes Provider's `kubernetes_secret` resource, you can 
access the registry from inside your cluster:

```
resource "digitalocean_container_registry_docker_credentials" "example" {
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
    ".dockerconfigjson" = digitalocean_container_registry_docker_credentials.example.docker_credentials
  }

  type = "kubernetes.io/dockerconfigjson"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the container registry.
* `write` - (Optional) Allow for write access to the container registry. Defaults to false.
* `expiry_seconds` - (Optional) The amount of time to pass before the Docker credentials expire in seconds. Defaults to 2147483647, or roughly 68 years. Must be greater than 0 and less than 2147483647.

## Attributes Reference

The following attributes are exported:

* `id`: The ID of the tag. This is the same as the name.
* `name` - The name of the container registry
* `docker_credentials`: Credentials for the container registry.
* `expiry_seconds`: Number of seconds after creation for token to expire.
* `credential_expiration_time`: The date and time the registry access token will expire.
