terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = ">= 2.8.0"
    }
  }
}

provider "digitalocean" {
  # You need to set this in your .bashrc
  # export DIGITALOCEAN_TOKEN="Your API TOKEN"
  #
}

resource "digitalocean_database_logsink" "logsink-01" {
  cluster_id = digitalocean_database_cluster.doby.id
  sink_name  = "fox2"
  sink_type  = "opensearch"


  config {
    url            = "https://user:passwd@192.168.0.1:25060"
    index_prefix   = "opensearch-logs"
    index_days_max = 5
  }
}

resource "digitalocean_database_cluster" "doby" {
  name       = "dobydb"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-2gb"
  region     = "nyc1"
  node_count = 1
  tags       = ["production"]
}
