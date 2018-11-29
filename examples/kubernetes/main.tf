resource "digitalocean_kubernetes_cluster" "foobar" {
  name    = "testing"
  region  = "lon1"
  version = "1.12.1-do.2"
  tags    = ["foo", "bar"]

  node_pool {
    size  = "s-1vcpu-2gb3"
    count = 1
    tags  = ["foo", "bar"]
  }
}

resource "digitalocean_kubernetes_node_pool" "barfoo" {
  cluster_id = "${digitalocean_kubernetes_cluster.foobar.id}"

  size  = "s-1vcpu-2gb3"
  count = 1
  tags  = ["foo", "bar"]
}
