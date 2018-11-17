resource "digitalocean_kubernetes_cluster" "foobar" {
  name    = "%s"
  region  = "lon1"
  version = "1.12.1-do.2"
  tags    = ["foo", "bar"]

  node_pool {
    name  = "default"
    size  = "s-1vcpu-2gb"
    count = 3
    tags  = ["one", "two"]
  }
}
