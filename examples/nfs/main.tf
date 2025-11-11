terraform {
  required_providers {
    digitalocean = {
      source  = "local/digitalocean/digitalocean"
      version = "99.0.0"
    }
  }
}

variable "do_token" {
  type        = string
  description = "DigitalOcean API token"
  sensitive   = true
}

provider "digitalocean" {
  token = var.do_token
}

resource "digitalocean_vpc" "test" {
  name   = "nfs-test-vpc"
  region = "atl1"
}

resource "digitalocean_nfs" "test" {
  name   = "nfs-test"
  region = "atl1"
  size   = 50
  vpc_id = digitalocean_vpc.test.id

  lifecycle {
    ignore_changes = [vpc_id]
  }
}

resource "digitalocean_nfs_attachment" "test" {
  vpc_id   = digitalocean_vpc.test.id
  share_id = digitalocean_nfs.test.id
  region   = "atl1"
}

output "nfs_id" {
  value = digitalocean_nfs.test.id
}

output "vpc1_id" {
  value = digitalocean_vpc.test.id
}

resource "digitalocean_vpc" "test2" {
  name   = "nfs-test-vpc-2"
  region = "atl1"
}

output "vpc2_id" {
  value = digitalocean_vpc.test2.id
}

output "current_attachment" {
  value = digitalocean_nfs_attachment.test.vpc_id
}


# SNAPSHOT

# resource "digitalocean_nfs_snapshot" "test" {
#   share_id = digitalocean_nfs.test.id
#   name     = "nfs-test-snapshot"
#   region   = "atl1"
# }
#
# output "snapshot_id" {
#   value = digitalocean_nfs_snapshot.test.id
# }
#
# output "snapshot_name" {
#   value = digitalocean_nfs_snapshot.test.name
# }
#
# output "snapshot_status"{
#     value = digitalocean_nfs_snapshot.test.status
# }