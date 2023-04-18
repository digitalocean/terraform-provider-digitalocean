---
page_title: "DigitalOcean: digitalocean_kubernetes_node_pool"
---

# digitalocean\_kubernetes\_node\_pool

Provides a DigitalOcean Kubernetes node pool resource. While the default node pool must be defined in the `digitalocean_kubernetes_cluster` resource, this resource can be used to add additional ones to a cluster.

## Example Usage

### Basic Example

```hcl
resource "digitalocean_kubernetes_cluster" "foo" {
  name    = "foo"
  region  = "nyc1"
  version = "1.22.8-do.1"

  node_pool {
    name       = "front-end-pool"
    size       = "s-2vcpu-2gb"
    node_count = 3
  }
}

resource "digitalocean_kubernetes_node_pool" "bar" {
  cluster_id = digitalocean_kubernetes_cluster.foo.id

  name       = "backend-pool"
  size       = "c-2"
  node_count = 2
  tags       = ["backend"]

  labels = {
    service  = "backend"
    priority = "high"
  }

  taint {
    key    = "workloadKind"
    value  = "database"
    effect = "NoSchedule"
  }
}
```

### Autoscaling Example

Node pools may also be configured to [autoscale](https://www.digitalocean.com/docs/kubernetes/how-to/autoscale/).
For example:

```hcl
resource "digitalocean_kubernetes_node_pool" "autoscale-pool-01" {
  cluster_id = digitalocean_kubernetes_cluster.foo.id
  name       = "autoscale-pool-01"
  size       = "s-1vcpu-2gb"
  auto_scale = true
  min_nodes  = 1
  max_nodes  = 5
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The ID of the Kubernetes cluster to which the node pool is associated.
* `name` - (Required) A name for the node pool.
* `size` - (Required) The slug identifier for the type of Droplet to be used as workers in the node pool.
* `node_count` - (Optional) The number of Droplet instances in the node pool. If auto-scaling is enabled, this should only be set if the desired result is to explicitly reset the number of nodes to this value. If auto-scaling is enabled, and the node count is outside of the given min/max range, it will use the min nodes value.
* `auto_scale` - (Optional) Enable auto-scaling of the number of nodes in the node pool within the given min/max range.
* `min_nodes` - (Optional) If auto-scaling is enabled, this represents the minimum number of nodes that the node pool can be scaled down to.
* `max_nodes` - (Optional) If auto-scaling is enabled, this represents the maximum number of nodes that the node pool can be scaled up to.
* `tags` - (Optional) A list of tag names to be applied to the Kubernetes cluster.
* `labels` - (Optional) A map of key/value pairs to apply to nodes in the pool. The labels are exposed in the Kubernetes API as labels in the metadata of the corresponding [Node resources](https://kubernetes.io/docs/concepts/architecture/nodes/).
* `taint` - (Optional) A list of taints applied to all nodes in the pool.

This resource supports [customized create timeouts](https://www.terraform.io/docs/language/resources/syntax.html#operation-timeouts). The default timeout is 30 minutes.

## Attributes Reference

In addition to the arguments listed above, the following additional attributes are exported:

* `id` -  A unique ID that can be used to identify and reference the node pool.
* `actual_node_count` - A computed field representing the actual number of nodes in the node pool, which is especially useful when auto-scaling is enabled.
* `nodes` - A list of nodes in the pool. Each node exports the following attributes:
  - `id` -  A unique ID that can be used to identify and reference the node.
  - `name` - The auto-generated name for the node.
  - `status` -  A string indicating the current status of the individual node.
  - `droplet_id` - The id of the node's droplet
  - `created_at` - The date and time when the node was created.
  - `updated_at` - The date and time when the node was last updated.
* `taint` - A list of taints applied to all nodes in the pool. Each taint exports the following attributes:
  - `key` - An arbitrary string. The "key" and "value" fields of the "taint" object form a key-value pair.
  - `value` - An arbitrary string. The "key" and "value" fields of the "taint" object form a key-value pair.
  - `effect` - How the node reacts to pods that it won't tolerate. Available effect values are: "NoSchedule", "PreferNoSchedule", "NoExecute".

## Import

If you are importing an existing Kubernetes cluster with a single node pool, just
import the cluster. Additional node pools can be imported by using their `id`, e.g.

```
terraform import digitalocean_kubernetes_node_pool.mynodepool 9d76f410-9284-4436-9633-4066852442c8
```

Note: If the node pool has the `terraform:default-node-pool` tag, then it is a default node pool for an
existing cluster. The provider will refuse to import the node pool in that case because the node pool
is managed by the `digitalocean_kubernetes_cluster` resource and not by this
`digitalocean_kubernetes_node_pool` resource.
