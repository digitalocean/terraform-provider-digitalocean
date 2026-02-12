terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = ">= 2.4.0"
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
  name             = "nfs-test"
  region           = "atl1"
  size             = 50
  vpc_id           = digitalocean_vpc.test.id
  performance_tier = "high" # Options: "standard" or "high". Can be changed to switch tiers after creation.

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

output "nfs_host" {
  value       = digitalocean_nfs.test.host
  description = "NFS server IP address"
}

output "nfs_mount_path" {
  value       = digitalocean_nfs.test.mount_path
  description = "NFS export path"
}

output "nfs_mount_command" {
  value       = "mount -t nfs ${digitalocean_nfs.test.host}:${digitalocean_nfs.test.mount_path} /mnt/nfs"
  description = "Example mount command"
}

output "nfs_performance_tier" {
  value       = digitalocean_nfs.test.performance_tier
  description = "Performance tier of the NFS share"
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
# Kindly note that Snapshots can be created after the NFS share is created.


resource "digitalocean_nfs_snapshot" "test" {
  share_id = digitalocean_nfs.test.id
  name     = "nfs-test-snapshot"
  region   = "atl1"
}

output "snapshot_id" {
  value = digitalocean_nfs_snapshot.test.id
}

output "snapshot_name" {
  value = digitalocean_nfs_snapshot.test.name
}

output "snapshot_status" {
  value = digitalocean_nfs_snapshot.test.status
}
