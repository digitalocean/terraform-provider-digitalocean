variable "primary_cluster" {
  type = object({
    id                     = string
    name                   = string
    endpoint               = string
    token                  = string
    cluster_ca_certificate = string
    raw_config             = string
  })
}

variable "cluster_id" {
  type = string
}

variable "write_kubeconfig" {
  type        = bool
  default     = false
}