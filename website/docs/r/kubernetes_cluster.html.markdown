---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_kubernetes_cluster"
sidebar_current: "docs-do-resource-kubernetes-cluster"
description: |-
  Provides a DigitalOcean Kubernetes cluster resource.
---

# digitalocean\_kubernetes\_cluster

Provides a DigitalOcean Kubernetes cluster resource. This can be used to create, delete, and modify clusters. For more information see the [official documentation](https://www.digitalocean.com/docs/kubernetes/).

## Example Usage

### Basic Example

```hcl
resource "digitalocean_kubernetes_cluster" "foo" {
  name    = "foo"
  region  = "nyc1"
  # Grab the latest version slug from `doctl kubernetes options versions`
  version = "1.15.5-do.1"

  node_pool {
    name       = "worker-pool"
    size       = "s-2vcpu-2gb"
    node_count = 3
  }
}
```

### Kubernetes Terraform Provider Example

The cluster's kubeconfig is exported as an attribute allowing you to use it with the [Kubernetes Terraform provider](https://www.terraform.io/docs/providers/kubernetes/index.html). For example:

```hcl
resource "digitalocean_kubernetes_cluster" "foo" {
  name    = "foo"
  region  = "nyc1"
  # Grab the latest version slug from `doctl kubernetes options versions`
  version = "1.15.5-do.1"
  tags    = ["staging"]

  node_pool {
    name       = "worker-pool"
    size       = "s-2vcpu-2gb"
    node_count = 3
  }
}

provider "kubernetes" {
  load_config_file = false
  host  = digitalocean_kubernetes_cluster.foo.endpoint
  token = digitalocean_kubernetes_cluster.foo.kube_config[0].token
  cluster_ca_certificate = base64decode(
    digitalocean_kubernetes_cluster.foo.kube_config[0].cluster_ca_certificate
  )
}
```

### Autoscaling Example

Node pools may also be configured to [autoscale](https://www.digitalocean.com/docs/kubernetes/how-to/autoscale/).
For example:

```
resource "digitalocean_kubernetes_cluster" "foo" {
  name    = "foo"
  region  = "nyc1"
  version = "1.15.5-do.1"

  node_pool {
    name       = "autoscale-worker-pool"
    size       = "s-2vcpu-2gb"
    auto_scale = true
    min_nodes  = 1
    max_nodes  = 5
  }
}
```

Note that, while individual node pools may scale to 0, a cluster must always include at least one node.

## Argument Reference

The following arguments are supported:

* `name` - (Required) A name for the Kubernetes cluster.
* `region` - (Required) The slug identifier for the region where the Kubernetes cluster will be created.
* `version` - (Required) The slug identifier for the version of Kubernetes used for the cluster. Use [doctl](https://github.com/digitalocean/doctl) to find the available versions `doctl kubernetes options versions`. (**Note:** A cluster may only be upgraded to newer versions in-place. If the version is decreased, a new resource will be created.)
* `node_pool` - (Required) A block representing the cluster's default node pool. Additional node pools may be added to the cluster using the `digitalocean_kubernetes_node_pool` resource. The following arguments may be specified:
  - `name` - (Required) A name for the node pool.
  - `size` - (Required) The slug identifier for the type of Droplet to be used as workers in the node pool.
  - `node_count` - (Optional) The number of Droplet instances in the node pool. If auto-scaling is enabled, this should only be set if the desired result is to explicitly reset the number of nodes to this value. If auto-scaling is enabled, and the node count is outside of the given min/max range, it will use the min nodes value.
  - `auto_scale` - (Optional) Enable auto-scaling of the number of nodes in the node pool within the given min/max range.
  - `min_nodes` - (Optional) If auto-scaling is enabled, this represents the minimum number of nodes that the node pool can be scaled down to.
  - `max_nodes` - (Optional) If auto-scaling is enabled, this represents the maximum number of nodes that the node pool can be scaled up to.
  - `tags` - (Optional) A list of tag names to be applied to the Kubernetes cluster.
* `tags` - (Optional) A list of tag names to be applied to the Kubernetes cluster.

## Attributes Reference

In addition to the arguments listed above, the following additional attributes are exported:

* `id` - A unique ID that can be used to identify and reference a Kubernetes cluster.
* `cluster_subnet` - The range of IP addresses in the overlay network of the Kubernetes cluster.
* `service_subnet` - The range of assignable IP addresses for services running in the Kubernetes cluster.
* `ipv4_address` - The public IPv4 address of the Kubernetes master node.
* `endpoint` - The base URL of the API server on the Kubernetes master node.
* `status` -  A string indicating the current status of the cluster. Potential values include running, provisioning, and errored.
* `created_at` - The date and time when the Kubernetes cluster was created.
* `updated_at` - The date and time when the Kubernetes cluster was last updated.
* `kube_config.0` - A representation of the Kubernetes cluster's kubeconfig with the following attributes:
  - `raw_config` - The full contents of the Kubernetes cluster's kubeconfig file.
  - `host` - The URL of the API server on the Kubernetes master node.
  - `cluster_ca_certificate` - The base64 encoded public certificate for the cluster's certificate authority.
  - `token` - The DigitalOcean API access token used by clients to access the cluster.
  - `client_key` - The base64 encoded private key used by clients to access the cluster. Only available if token authentication is not supported on your cluster.
  - `client_certificate` - The base64 encoded public certificate used by clients to access the cluster. Only available if token authentication is not supported on your cluster.
  - `expires_at` - The date and time when the credentials will expire and need to be regenerated.
* `node_pool` - In addition to the arguments provided, these additional attributes about the cluster's default node pool are exported:
  - `id` -  A unique ID that can be used to identify and reference the node pool.
  - `actual_node_count` - A computed field representing the actual number of nodes in the node pool, which is especially useful when auto-scaling is enabled.
  - `nodes` - A list of nodes in the pool. Each node exports the following attributes:
     + `id` -  A unique ID that can be used to identify and reference the node.
     + `name` - The auto-generated name for the node.
     + `status` -  A string indicating the current status of the individual node.
     + `droplet_id` - The id of the node's droplet
     + `created_at` - The date and time when the node was created.
     + `updated_at` - The date and time when the node was last updated.

## Import

Kubernetes clusters can not be imported at this time.
