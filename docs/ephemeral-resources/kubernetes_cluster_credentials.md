---
page_title: "DigitalOcean: digitalocean_kubernetes_cluster_credentials"
subcategory: "Kubernetes"
---

# digitalocean\_kubernetes\_cluster\_credentials

Fetches credentials for a DigitalOcean Kubernetes cluster as an [ephemeral resource](https://developer.hashicorp.com/terraform/language/resources/ephemeral). The credentials are retrieved on each Terraform run and are never stored in the plan or state, unlike the `kube_config` attribute of the `digitalocean_kubernetes_cluster` resource and data source.

~> **Note:** Ephemeral resources require Terraform 1.10 or later.

## Example Usage

```hcl
ephemeral "digitalocean_kubernetes_cluster_credentials" "example" {
  cluster_id = digitalocean_kubernetes_cluster.example.id
}

provider "kubernetes" {
  host  = ephemeral.digitalocean_kubernetes_cluster_credentials.example.host
  token = ephemeral.digitalocean_kubernetes_cluster_credentials.example.token
  cluster_ca_certificate = base64decode(
    ephemeral.digitalocean_kubernetes_cluster_credentials.example.cluster_ca_certificate
  )
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The ID of the Kubernetes cluster.
* `expiry_seconds` - (Optional) The duration in seconds that the returned credentials will be valid. Defaults to the API's default expiry.

## Attributes Reference

The following attributes are exported:

* `raw_config` - The full contents of the Kubernetes cluster's kubeconfig file.
* `host` - The URL of the API server on the Kubernetes master node.
* `cluster_ca_certificate` - The base64 encoded public certificate for the cluster's certificate authority.
* `token` - The DigitalOcean API access token used by clients to access the cluster.
* `client_key` - The PEM-encoded private key used by clients to access the cluster. Only available if token authentication is not supported on your cluster.
* `client_certificate` - The PEM-encoded public certificate used by clients to access the cluster. Only available if token authentication is not supported on your cluster.
* `expires_at` - The date and time when the credentials will expire and need to be regenerated.

Credentials are not available for clusters that require SSO, because the API does not issue token credentials for them.
