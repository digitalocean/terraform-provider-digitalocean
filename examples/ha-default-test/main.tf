# Test cluster creation WITHOUT ha - API applies default (true for 1.36+)
# Uses your locally built provider via dev_overrides in ~/.terraformrc

terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = ">= 2.44.1"
    }
  }
}

provider "digitalocean" {
  # Uses DIGITALOCEAN_ACCESS_TOKEN env var, or set token = "..." here
}

data "digitalocean_kubernetes_versions" "test" {}

resource "digitalocean_kubernetes_cluster" "test" {
  name    = "terra-ha-true1"
  region  = "nyc1"
  version = "1.35.1-do.6"
#  "1.36.0-do.0"
# "1.35.1-do.6"
  # ha is intentionally OMITTED - API applies version-dependent default
 # ha = true
 # ha = false
  node_pool {
    name       = "default"
    size       = "s-1vcpu-2gb"
    node_count = 1
  }
}

output "ha" {
  value       = digitalocean_kubernetes_cluster.test.ha
  description = "HA value from API (should be true for 1.36+)"
}

output "cluster_id" {
  value = digitalocean_kubernetes_cluster.test.id
}
