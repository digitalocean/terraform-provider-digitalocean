data "digitalocean_kubernetes_versions" "example" {
  version_prefix = "1.19."
}

resource "digitalocean_kubernetes_cluster" "foobar" {
  name    = "testing"
  region  = "lon1"
  version = data.digitalocean_kubernetes_versions.example.latest_version
  tags    = ["foo", "bar"]

  node_pool {
    name       = "foobar"
    size       = "s-1vcpu-2gb"
    node_count = 1
    tags       = ["one", "two"] // Tags from cluster are automatically added to node pools
  }
}

// Add a second node pool
resource "digitalocean_kubernetes_node_pool" "barfoo" {
  cluster_id = digitalocean_kubernetes_cluster.foobar.id

  name       = "barfoo"
  size       = "s-1vcpu-2gb"
  node_count = 1
  tags       = ["three", "four"] // Tags from cluster are automatically added to node pools
}
