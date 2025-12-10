---
page_title: "DigitalOcean: digitalocean_kubernetes_cluster"
subcategory: "Kubernetes"
---

# digitalocean\_kubernetes\_cluster

Provides a DigitalOcean Kubernetes cluster resource. This can be used to create, delete, and modify clusters. For more information see the [official documentation](https://www.digitalocean.com/docs/kubernetes/).

## Example Usage

### Basic Example

```hcl
resource "digitalocean_kubernetes_cluster" "foo" {
  name   = "foo"
  region = "nyc1"
  # Grab the latest version slug from `doctl kubernetes options versions` (e.g. "1.14.6-do.1"
  # If set to "latest", latest published version will be used.
  version = "latest"

  node_pool {
    name       = "worker-pool"
    size       = "s-2vcpu-2gb"
    node_count = 3

    taint {
      key    = "workloadKind"
      value  = "database"
      effect = "NoSchedule"
    }
  }
}
```

### Autoscaling Example

Node pools may also be configured to [autoscale](https://www.digitalocean.com/docs/kubernetes/how-to/autoscale/).
For example:

```hcl
resource "digitalocean_kubernetes_cluster" "foo" {
  name    = "foo"
  region  = "nyc1"
  version = "1.22.8-do.1"

  node_pool {
    name       = "autoscale-worker-pool"
    size       = "s-2vcpu-2gb"
    auto_scale = true
    min_nodes  = 1
    max_nodes  = 5
  }
}
```

Note that, currently, each node pool must always have at least one node and when using autoscaling the min_nodes must be greater than or equal to 1.
> Autoscaling to zero (`min_nodes=0`) is in [private preview](https://docs.digitalocean.com/release-notes/kubernetes/#2025-01-07) and not available for public use.

### Auto Upgrade Example

DigitalOcean Kubernetes clusters may also be configured to [auto upgrade](https://www.digitalocean.com/docs/kubernetes/how-to/upgrade-cluster/#automatically) patch versions. You may explicitly specify the maintenance window policy.
For example:

```hcl
data "digitalocean_kubernetes_versions" "example" {
  version_prefix = "1.22."
}

resource "digitalocean_kubernetes_cluster" "foo" {
  name         = "foo"
  region       = "nyc1"
  auto_upgrade = true
  version      = data.digitalocean_kubernetes_versions.example.latest_version

  maintenance_policy {
    start_time = "04:00"
    day        = "sunday"
  }

  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 3
  }
}
```

Note that a data source is used to supply the version. This is needed to prevent configuration diff whenever a cluster is upgraded.

### Kubernetes Terraform Provider Example

The cluster's kubeconfig is exported as an attribute allowing you to use it with
the [Kubernetes Terraform provider](https://registry.terraform.io/providers/hashicorp/kubernetes/latest/docs).

~> When using interpolation to pass credentials from a `digitalocean_kubernetes_cluster`
resource to the Kubernetes provider, the cluster resource generally should not
be created in the same Terraform module where Kubernetes provider resources are
also used. This can lead to unpredictable errors which are hard to debug and
diagnose. The root issue lies with the order in which Terraform itself evaluates
the provider blocks vs. actual resources.

When using the Kubernetes provider with a cluster created in a separate Terraform
module or configuration, use the [`digitalocean_kubernetes_cluster` data-source](https://registry.terraform.io/providers/digitalocean/digitalocean/latest/docs/data-sources/kubernetes_cluster)
to access the cluster's credentials. [See here for a full example](https://github.com/digitalocean/terraform-provider-digitalocean/tree/main/examples/kubernetes).

```hcl
data "digitalocean_kubernetes_cluster" "example" {
  name = "prod-cluster-01"
}

provider "kubernetes" {
  host  = data.digitalocean_kubernetes_cluster.example.endpoint
  token = data.digitalocean_kubernetes_cluster.example.kube_config[0].token
  cluster_ca_certificate = base64decode(
    data.digitalocean_kubernetes_cluster.example.kube_config[0].cluster_ca_certificate
  )
}
```

#### Exec credential plugin

Another method to ensure that the Kubernetes provider is receiving valid credentials
is to use an exec plugin. In order to use use this approach, the DigitalOcean
CLI (`doctl`) must be present. `doctl` will renew the token if needed before
initializing the provider.

```hcl
provider "kubernetes" {
  host = data.digitalocean_kubernetes_cluster.foo.endpoint
  cluster_ca_certificate = base64decode(
    data.digitalocean_kubernetes_cluster.foo.kube_config[0].cluster_ca_certificate
  )

  exec {
    api_version = "client.authentication.k8s.io/v1beta1"
    command     = "doctl"
    args = ["kubernetes", "cluster", "kubeconfig", "exec-credential",
    "--version=v1beta1", data.digitalocean_kubernetes_cluster.foo.id]
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A name for the Kubernetes cluster.
* `region` - (Required) The slug identifier for the region where the Kubernetes cluster will be created.
* `version` - (Required) The slug identifier for the version of Kubernetes used for the cluster. Use [doctl](https://github.com/digitalocean/doctl) to find the available versions `doctl kubernetes options versions`. (**Note:** A cluster may only be upgraded to newer versions in-place. If the version is decreased, a new resource will be created.)
* `cluster_subnet` - (Optional) The range of IP addresses in the overlay network of the Kubernetes cluster. For more information, see [here](https://docs.digitalocean.com/products/kubernetes/how-to/create-clusters/#create-with-vpc-native).
* `service_subnet` - (Optional) The range of assignable IP addresses for services running in the Kubernetes cluster. For more information, see [here](https://docs.digitalocean.com/products/kubernetes/how-to/create-clusters/#create-with-vpc-native).
* `control_plane_firewall` - (Optional) A block representing the cluster's control plane firewall
  - `enabled` - (Required) Boolean flag whether the firewall should be enabled or not.
  - `allowed_addresses` - (Required) A list of addresses allowed (CIDR notation).
* `vpc_uuid` - (Optional) The ID of the VPC where the Kubernetes cluster will be located.
* `auto_upgrade` - (Optional) A boolean value indicating whether the cluster will be automatically upgraded to new patch releases during its maintenance window.
* `surge_upgrade` - (Optional) Enable/disable surge upgrades for a cluster. Default: true
* `ha` - (Optional) Enable/disable the high availability control plane for a cluster. Once enabled for a cluster, high availability cannot be disabled. Default: false
* `registry_integration` - (optional) Enables or disables the DigitalOcean container registry integration for the cluster. This requires that a container registry has first been created for the account. Default: false
* `node_pool` - (Required) A block representing the cluster's default node pool. Additional node pools may be added to the cluster using the `digitalocean_kubernetes_node_pool` resource. The following arguments may be specified:
  - `name` - (Required) A name for the node pool.
  - `size` - (Required) The slug identifier for the type of Droplet to be used as workers in the node pool.
  - `node_count` - (Optional) The number of Droplet instances in the node pool. If auto-scaling is enabled, this should only be set if the desired result is to explicitly reset the number of nodes to this value. If auto-scaling is enabled, and the node count is outside of the given min/max range, it will use the min nodes value.
  - `auto_scale` - (Optional) Enable auto-scaling of the number of nodes in the node pool within the given min/max range.
  - `min_nodes` - (Optional) If auto-scaling is enabled, this represents the minimum number of nodes that the node pool can be scaled down to.
  - `max_nodes` - (Optional) If auto-scaling is enabled, this represents the maximum number of nodes that the node pool can be scaled up to.
  - `tags` - (Optional) A list of tag names applied to the node pool.
  - `labels` - (Optional) A map of key/value pairs to apply to nodes in the pool. The labels are exposed in the Kubernetes API as labels in the metadata of the corresponding [Node resources](https://kubernetes.io/docs/concepts/architecture/nodes/).
* `tags` - (Optional) A list of tag names to be applied to the Kubernetes cluster.
* `maintenance_policy` - (Optional) A block representing the cluster's maintenance window. Updates will be applied within this window. If not specified, a default maintenance window will be chosen. `auto_upgrade` must be set to `true` for this to have an effect.
  - `day` - (Required) The day of the maintenance window policy. May be one of "monday" through "sunday", or "any" to indicate an arbitrary week day.
  - `start_time` (Required) The start time in UTC of the maintenance window policy in 24-hour clock format / HH:MM notation (e.g., 15:00).
* `destroy_all_associated_resources` - (Optional) **Use with caution.** When set to true, all associated DigitalOcean resources created via the Kubernetes API (load balancers, volumes, and volume snapshots) will be destroyed along with the cluster when it is destroyed.
* `kubeconfig_expire_seconds` - (Optional) The duration in seconds that the returned Kubernetes credentials will be valid. If not set or 0, the credentials will have a 7 day expiry.
* `routing_agent` - (Optional) Block containing options for the routing-agent component. If not specified, the routing-agent component will not be installed in the cluster.
  - `enabled` - (Required) Boolean flag whether the routing-agent should be enabled or not.
* `amd_gpu_device_plugin` - (Optional) Block containing options for the AMD GPU device plugin component. If not specified, the component will be enabled by default for clusters with AMD GPU nodes.
  - `enabled` - (Required) Boolean flag whether the component should be enabled or not.
`amd_gpu_device_metrics_exporter_plugin` - (Optional) Block containing options for the AMD GPU device metrics exporter component. If not specified, the component will not be installed in the cluster.
    - `enabled` - (Required) Boolean flag whether the component should be enabled or not.
* `nvidia_gpu_device_plugin` - (Optional) Block containing options for the NVIDIA GPU device plugin component. If not specified, the component will be enabled by default for clusters with NVIDIA GPU nodes.
  - `enabled` - (Required) Boolean flag whether the component should be enabled or not.
`rdma_shared_device_plugin` - (Optional) Block containing options for the RDMA Shared Device Plugin (k8s-rdma-shared-dev-plugin) component. If not specified, the component will be enabled by default for clusters with GPU nodes connected to a dedicated high-speed networking fabric.
    - `enabled` - (Required) Boolean flag whether the component should be enabled or not.
* `cluster_autoscaler_configuration` - (Optional) Block containing options for cluster auto-scaling.
  - `scale_down_utilization_threshold` - (Optional) Float setting the Node utilization level, defined as sum of requested resources divided by capacity, in which a node can be considered for scale down.
  - `scale_down_unneeded_time` - (Optional) String setting how long a node should be unneeded before it's eligible for scale down.

This resource supports [customized create timeouts](https://www.terraform.io/docs/language/resources/syntax.html#operation-timeouts). The default timeout is 30 minutes.

## Attributes Reference

In addition to the arguments listed above, the following additional attributes are exported:

* `id` - A unique ID that can be used to identify and reference a Kubernetes cluster.
* `ipv4_address` - The public IPv4 address of the Kubernetes master node. This will not be set if high availability is configured on the cluster (v1.21+)
* `endpoint` - The base URL of the API server on the Kubernetes master node.
* `status` -  A string indicating the current status of the cluster. Potential values include running, provisioning, and errored.
* `created_at` - The date and time when the Kubernetes cluster was created.
* `updated_at` - The date and time when the Kubernetes cluster was last updated.
* `auto_upgrade` - A boolean value indicating whether the cluster will be automatically upgraded to new patch releases during its maintenance window.
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
  - `taint` - A block representing a taint applied to all nodes in the pool. Each taint exports the following attributes (taints must be unique by key and effect pair):
    + `key` - An arbitrary string. The "key" and "value" fields of the "taint" object form a key-value pair.
    + `value` - An arbitrary string. The "key" and "value" fields of the "taint" object form a key-value pair.
    + `effect` - How the node reacts to pods that it won't tolerate. Available effect values are: "NoSchedule", "PreferNoSchedule", "NoExecute".
* `urn` - The uniform resource name (URN) for the Kubernetes cluster.
* `maintenance_policy` - A block representing the cluster's maintenance window. Updates will be applied within this window. If not specified, a default maintenance window will be chosen.
  - `day` - The day of the maintenance window policy. May be one of "monday" through "sunday", or "any" to indicate an arbitrary week day.
  - `duration` A string denoting the duration of the service window, e.g., "04:00".
  - `start_time` The hour in UTC when maintenance updates will be applied, in 24 hour format (e.g. “16:00”).
* `routing_agent` - Block containing options for the routing-agent component.
  - `enabled` - Boolean flag whether the routing-agent is enabled or not.
* `amd_gpu_device_plugin` - Block containing options for the AMD GPU device plugin component. If not specified, the component will be enabled by default for clusters with AMD GPU nodes.
  - `enabled` - Boolean flag whether the component is enabled or not.
* `amd_gpu_device_metrics_exporter_plugin` - Block containing options for the AMD GPU device metrics exporter component.
  - `enabled` - Boolean flag whether the component is enabled or not.

## Import

Before importing a Kubernetes cluster, the cluster's default node pool must be tagged with
the `terraform:default-node-pool` tag. The provider will automatically add this tag if
the cluster only has a single node pool. Clusters with more than one node pool, however, will require
that you manually add the `terraform:default-node-pool` tag to the node pool that you intend to be
the default node pool.

Then the Kubernetes cluster and its default node pool can be imported using the cluster's `id`, e.g.

```
terraform import digitalocean_kubernetes_cluster.mycluster 1b8b2100-0e9f-4e8f-ad78-9eb578c2a0af
```

Additional node pools must be imported separately as `digitalocean_kubernetes_cluster`
resources, e.g.

```
terraform import digitalocean_kubernetes_node_pool.mynodepool 9d76f410-9284-4436-9633-4066852442c8
```
