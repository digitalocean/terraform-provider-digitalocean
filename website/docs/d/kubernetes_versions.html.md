---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_kubernetes_versions"
sidebar_current: "docs-do-datasource-kubernetes-versions"
description: |-
  Get available DigitalOcean Kubernetes versions.
---

# digitalocean\_kubernetes\_versions

Provides access to the available DigitalOcean Kubernetes Service versions.

## Example Usage

### Output a list of all available versions

```hcl
data "digitalocean_kubernetes_versions" "example" {}

output "k8s-versions" {
  value = data.digitalocean_kubernetes_versions.example.valid_versions
}
```

### Create a Kubernetes cluster using the most recent version available

```hcl
data "digitalocean_kubernetes_versions" "example" {}

resource "digitalocean_kubernetes_cluster" "example-cluster" {
  name    = "example-cluster"
  region  = "lon1"
  version = data.digitalocean_kubernetes_versions.example.latest_version

  node_pool {
    name = "default"
    size  = "s-1vcpu-2gb"
    node_count = 3
  }
}
```

### Pin a Kubernetes cluster to a specific minor version

```hcl
data "digitalocean_kubernetes_versions" "example" {
  version_prefix = "1.16."
}

resource "digitalocean_kubernetes_cluster" "example-cluster" {
  name    = "example-cluster"
  region  = "lon1"
  version = data.digitalocean_kubernetes_versions.example.latest_version

  node_pool {
    name = "default"
    size  = "s-1vcpu-2gb"
    node_count = 3
  }
}
```

## Argument Reference

The following arguments are supported:

* `version_prefix` - (Optional) If provided, Terraform will only return versions that match the string prefix. For example, `1.15.` will match all 1.15.x series releases.

## Attributes Reference

The following attributes are exported:

* `valid_versions` - A list of available versions.
* `latest_version` - The most recent version available.
