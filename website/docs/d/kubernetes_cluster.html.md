---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_kubernetes_cluster"
sidebar_current: "docs-do-datasource-kubernetes-cluster"
description: |-
  Get information on a DigitalOcean Kubernetes cluster.
---

# digitalocean\_kubernetes\_cluster

Retrieves information about a DigitalOcean Kubernetes cluster for use in other resources. This data source provides all of the cluster's properties as configured on your DigitalOcean account. This is useful if the cluster in question is not managed by Terraform.

## Example Usage

```hcl
data "digitalocean_kubernetes_cluster" "example" {
  name = "prod-cluster-01"
}

provider "kubernetes" {
  host = "${data.digitalocean_kubernetes_cluster.example.endpoint}"

  client_certificate     = "${base64decode(data.digitalocean_kubernetes_cluster.example.kube_config.0.client_certificate)}"
  client_key             = "${base64decode(data.digitalocean_kubernetes_cluster.example.kube_config.0.client_key)}"
  cluster_ca_certificate = "${base64decode(data.digitalocean_kubernetes_cluster.example.kube_config.0.cluster_ca_certificate)}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of Kubernetes cluster.

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID that can be used to identify and reference a Kubernetes cluster.
* `region` - The slug identifier for the region where the Kubernetes cluster is located.
* `version` - The slug identifier for the version of Kubernetes used for the cluster.
* `tags` - A list of tag names to be applied to the Kubernetes cluster.
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
* `node_pool` - A list of node pools associated with the cluster. Each node pool exports the following attributes:
  - `id` -  The unique ID that can be used to identify and reference the node pool.
  - `name` - The name of the node pool.
  - `size` - The slug identifier for the type of Droplet used as workers in the node pool.
  - `node_count` - The number of Droplet instances in the node pool.
  - `tags` - A list of tag names applied to the node pool.
  - `nodes` - A list of nodes in the pool. Each node exports the following attributes:
     + `id` -  A unique ID that can be used to identify and reference the node.
     + `name` - The auto-generated name for the node.
     + `status` -  A string indicating the current status of the individual node.
     + `created_at` - The date and time when the node was created.
     + `updated_at` - The date and time when the node was last updated.
