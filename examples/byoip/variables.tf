variable "byoip_prefix_uuid" {
  description = "The UUID of your BYOIP prefix (created outside of Terraform)"
  type        = string
}

variable "droplet_name" {
  description = "Name for the Droplet"
  type        = string
  default     = "byoip-example"
}

variable "droplet_size" {
  description = "Size of the Droplet"
  type        = string
  default     = "s-1vcpu-1gb"
}
