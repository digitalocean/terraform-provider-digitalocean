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

```hcl
resource "digitalocean_kubernetes_cluster" "foo" {
  name    = "foo"
  region  = "nyc1"
  // Grab the latest version slug from `doctl kubernetes options versions`
  version = "1.15.3-do.2"

  node_pool {
    name       = "worker-pool"
    size       = "s-2vcpu-2gb"
    node_count = 3
  }
}
```

The cluster's kubeconfig is exported as an attribute allowing you to use it with the [Kubernetes Terraform provider](https://www.terraform.io/docs/providers/kubernetes/index.html). For example:

```hcl
resource "digitalocean_kubernetes_cluster" "foo" {
  name    = "foo"
  region  = "nyc1"
  // Grab the latest version slug from `doctl kubernetes options versions`
  version = "1.15.3-do.2"
  tags    = ["staging"]

  node_pool {
    name       = "worker-pool"
    size       = "s-2vcpu-2gb"
    node_count = 3
  }
}

provider "kubernetes" {
  host = "${digitalocean_kubernetes_cluster.foo.endpoint}"

  client_certificate     = "${base64decode(digitalocean_kubernetes_cluster.foo.kube_config.0.client_certificate)}"
  client_key             = "${base64decode(digitalocean_kubernetes_cluster.foo.kube_config.0.client_key)}"
  cluster_ca_certificate = "${base64decode(digitalocean_kubernetes_cluster.foo.kube_config.0.cluster_ca_certificate)}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A name for the Kubernetes cluster.
* `region` - (Required) The slug identifier for the region where the Kubernetes cluster will be created.
* `version` - (Required) The slug identifier for the version of Kubernetes used for the cluster. Use [doctl](https://github.com/digitalocean/doctl) to find the available versions `doctl kubernetes options versions`.
* `node_pool` - (Required) A block representing the cluster's default node pool. Additional node pools may be added to the cluster using the `digitalocean_kubernetes_node_pool` resource. The following arguments may be specified:
  - `name` - (Required) A name for the node pool.
  - `size` - (Required) The slug identifier for the type of Droplet to be used as workers in the node pool.
  - `node_count` - (Required) The number of Droplet instances in the node pool.
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
  - `client_key` - The base64 encoded private key used by clients to access the cluster.
  - `client_certificate` - The base64 encoded public certificate used by clients to access the cluster.
  - `cluster_ca_certificate` - The base64 encoded public certificate for the cluster's certificate authority.
* `node_pool` - In addition to the arguments provided, these additional attributes about the cluster's default node pool are exported:
  - `id` -  A unique ID that can be used to identify and reference the node pool.
  - `nodes` - A list of nodes in the pool. Each node exports the following attributes:
     + `id` -  A unique ID that can be used to identify and reference the node.
     + `name` - The auto-generated name for the node.
     + `status` -  A string indicating the current status of the individual node.
     + `created_at` - The date and time when the node was created.
     + `updated_at` - The date and time when the node was last updated.

## Import

Kubernetes clusters can not be imported at this time.
