---
page_title: "DigitalOcean: digitalocean_container_registry_docker_credentials"
---

# digitalocean_container_registry_docker_credentials

Get Docker credentials for your DigitalOcean container registry.

An error is triggered if the provided container registry name does not exist.

## Example Usage

### Basic Example

Get the container registry:

```hcl
resource "digitalocean_container_registry_docker_credentials" "example" {
  registry_name = "example"
}
```

### Docker Provider Example

Use the `endpoint` and `docker_credentials` with the Docker provider:

```hcl
data "digitalocean_container_registry" "example" {
  name = "example"
}

resource "digitalocean_container_registry_docker_credentials" "example" {
  registry_name = "example"
}

provider "docker" {
  host = "unix://var/run/docker.sock"

  registry_auth {
    address             = data.digitalocean_container_registry.example.server_url
    config_file_content = digitalocean_container_registry_docker_credentials.example.docker_credentials
  }
}
```

### Kubernetes Example

Combined with the Kubernetes Provider's `kubernetes_secret` resource, you can
access the registry from inside your cluster:

```hcl
resource "digitalocean_container_registry_docker_credentials" "example" {
  registry_name = "example"
}

data "digitalocean_kubernetes_cluster" "example" {
  name = "prod-cluster-01"
}

provider "kubernetes" {
  host             = data.digitalocean_kubernetes_cluster.example.endpoint
  token            = data.digitalocean_kubernetes_cluster.example.kube_config[0].token
  cluster_ca_certificate = base64decode(
    data.digitalocean_kubernetes_cluster.example.kube_config[0].cluster_ca_certificate
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

* `registry_name` - (Required) The name of the container registry.
* `write` - (Optional) Allow for write access to the container registry. Defaults to false.
* `expiry_seconds` - (Optional) The amount of time to pass before the Docker credentials expire in seconds. Defaults to 1576800000, or roughly 50 years. Must be greater than 0 and less than 1576800000.

## Attributes Reference

The following attributes are exported:

* `id`: The ID of the tag. This is the same as the name.
* `registry_name` - The name of the container registry
* `docker_credentials`: Credentials for the container registry.
* `expiry_seconds`: Number of seconds after creation for token to expire.
* `credential_expiration_time`: The date and time the registry access token will expire.
