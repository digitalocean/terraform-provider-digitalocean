terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = ">= 2.61.0"
    }
  }
}

resource "digitalocean_kubernetes_cluster" "main" {
  name                             = "tf-k8s-random122"
  tags                             = ["k8sss", "foosss1ss2"]
  region                           = "nyc3"
  version                          = "1.33.1-do.2"
  destroy_all_associated_resources = false

  cluster_autoscaler_configuration {
    scale_down_utilization_threshold = 0.5
    scale_down_unneeded_time         = "1m30s"
  }

  node_pool {
    name       = "foo"
    size       = "s-2vcpu-2gb"
    node_count = 1
    auto_scale = false
  }
}