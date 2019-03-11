---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_kubernetes_node_pool"
sidebar_current: "docs-do-resource-kubernetes-node-pool"
description: |-
  Provides a DigitalOcean Kubernetes node pool.
---

# digitalocean\_kubernetes\_node\_pool

~> **NOTE:** DigitalOcean Kubernetes is currently in [Limited Availability](https://www.digitalocean.com/docs/platform/product-lifecycle/). In order to access its API, you must first enable Kubernetes on your account by opting-in via the [cloud control panel](https://cloud.digitalocean.com/kubernetes/clusters). While the Kubernetes Cluster functionality is currently in limited availability the structure of this resource may change over time. Please share any feedback you may have by [opening an issue on GitHub](https://github.com/terraform-providers/terraform-provider-digitalocean/issues).

Provides a DigitalOcean Kubernetes node pool resource. While the default node pool must be defined in the `digitalocean_kubernetes_cluster` resource, this resource can be used to add additional ones to a cluster.

## Example Usage

```hcl
resource "digitalocean_kubernetes_cluster" "foo" {
  name    = "foo"
  region  = "nyc1"
  version = "1.12.1-do.2"

  node_pool {
    name       = "front-end-pool"
    size       = "s-2vcpu-2gb"
    node_count = 3
  }
}

resource "digitalocean_kubernetes_node_pool" "bar" {
  cluster_id = "${digitalocean_kubernetes_cluster.foo.id}"

  name       = "backend-pool"
  size       = "c-2"
  node_count = 2
  tags       = ["backend"]
}
```


## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The ID of the Kubernetes cluster to which the node pool is associated.
* `name` - (Required) A name for the node pool.
* `size` - (Required) The slug identifier for the type of Droplet to be used as workers in the node pool.
* `node_count` - (Required) The number of Droplet instances in the node pool.
* `tags` - (Optional) A list of tag names to be applied to the Kubernetes cluster.

## Attributes Reference

In addition to the arguments listed above, the following additional attributes are exported:

* `id` -  A unique ID that can be used to identify and reference the node pool.
* `nodes` - A list of nodes in the pool. Each node exports the following attributes:
  - `id` -  A unique ID that can be used to identify and reference the node.
  - `name` - The auto-generated name for the node.
  - `status` -  A string indicating the current status of the individual node.
  - `created_at` - The date and time when the node was created.
  - `updated_at` - The date and time when the node was last updated.

## Import

Kubernetes node pools can not be imported at this time.
