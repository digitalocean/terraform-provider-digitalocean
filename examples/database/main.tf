terraform {
  required_version = "~> 1"
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2"
    }
  }
}

provider "digitalocean" {
  # You need to set this in your .bashrc
  # export DIGITALOCEAN_TOKEN="Your API TOKEN"
}

data "digitalocean_vpc" "default" {
  region = var.region
}

data "digitalocean_database_metrics_credentials" "test" {}

resource "digitalocean_database_cluster" "test" {
  name                 = var.db_name
  engine               = var.db_engine
  version              = var.db_version
  size                 = var.db_size
  region               = var.region
  node_count           = var.db_node_count
  private_network_uuid = data.digitalocean_vpc.default.id
}
