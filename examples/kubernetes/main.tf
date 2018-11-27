resource "digitalocean_kubernetes_cluster" "foobar" {
  name    = "testing"
  region  = "lon1"
  version = "1.12.1-do.2"
  tags    = ["foo", "bar"]

  node_pool {
    name  = "pool1"
    size  = "s-1vcpu-2gb3"
    count = 5
    tags  = ["foo", "bar"]
  }
}
