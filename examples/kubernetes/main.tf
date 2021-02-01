terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = ">= 2.4.0"
    }
    kubernetes = {
      source = "hashicorp/kubernetes"
      version = ">= 2.0.0"
    }
    helm = {
      source  = "hashicorp/helm"
      version = ">= 2.0.1"
    }
  }
}

resource "random_id" "cluster_name" {
  byte_length = 5
}

locals {
  cluster_name = "tf-k8s-${random_id.cluster_name.hex}"
}

module "doks-cluster" {
  source             = "./doks-cluster"
  cluster_name       = local.cluster_name
  cluster_region     = "nyc3"
  cluster_version    = var.cluster_version

  worker_size        = var.worker_size
  worker_count       = var.worker_count
}

module "kubernetes-config" {
  source           = "./kubernetes-config"
  cluster_name     = module.doks-cluster.cluster_name
  cluster_id       = module.doks-cluster.cluster_id

  write_kubeconfig = var.write_kubeconfig
}
