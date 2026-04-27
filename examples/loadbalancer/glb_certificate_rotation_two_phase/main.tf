# Two-phase GLB custom certificate rotation — certificate resources only.
#
# Phase 1: add digitalocean_certificate.new; keep digitalocean_certificate.old;
#   in your EXISTING digitalocean_loadbalancer, set:
#     domains { ... certificate_name = digitalocean_certificate.new.name }
#   then terraform apply (no destroy of old cert).
#
# Phase 2: remove the digitalocean_certificate.old block from this file (and
#   its variables); terraform apply again to delete the old cert only.

terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = ">= 2.50.0"
    }
  }
}

variable "do_token" {
  type        = string
  sensitive   = true
  description = "DigitalOcean API token (or use DIGITALOCEAN_TOKEN and omit)."
}

variable "old_cert_name" {
  type        = string
  description = "Existing certificate name (Phase 1 only)."
}

variable "new_cert_name" {
  type        = string
  description = "New unique certificate name for rotation."
}

variable "old_private_key" {
  type      = string
  sensitive = true
}

variable "old_leaf_certificate" {
  type      = string
  sensitive = true
}

variable "old_certificate_chain" {
  type      = string
  sensitive = true
}

variable "new_private_key" {
  type      = string
  sensitive = true
}

variable "new_leaf_certificate" {
  type      = string
  sensitive = true
}

variable "new_certificate_chain" {
  type      = string
  sensitive = true
}

provider "digitalocean" {
  token = var.do_token
}

# Phase 2: remove this entire resource and its variables.
resource "digitalocean_certificate" "old" {
  name              = var.old_cert_name
  type              = "custom"
  private_key       = var.old_private_key
  leaf_certificate  = var.old_leaf_certificate
  certificate_chain = var.old_certificate_chain
}

resource "digitalocean_certificate" "new" {
  name              = var.new_cert_name
  type              = "custom"
  private_key       = var.new_private_key
  leaf_certificate  = var.new_leaf_certificate
  certificate_chain = var.new_certificate_chain
}

output "attach_this_certificate_name_to_glb_domains" {
  description = "Set domains.certificate_name in your Global LB to this value for Phase 1."
  value       = digitalocean_certificate.new.name
}
