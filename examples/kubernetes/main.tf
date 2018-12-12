resource "digitalocean_kubernetes_cluster" "foobar" {
  name    = "testing"
  region  = "lon1"
  version = "1.12.1-do.2"
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
  cluster_id = "${digitalocean_kubernetes_cluster.foobar.id}"

  size       = "s-1vcpu-2gb3"
  node_count = 1
  tags       = ["three", "four"] // Tags from cluster are automatically added to node pools
}
