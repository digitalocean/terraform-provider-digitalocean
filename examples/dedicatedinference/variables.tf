variable "name" {
  description = "A human-readable name for the dedicated inference endpoint."
  type        = string
  default     = "my-inference-endpoint"
}

variable "region" {
  description = "The region slug where the dedicated inference endpoint will be deployed."
  type        = string
  default     = "tor1"
}

variable "model_slug" {
  description = "The slug identifier for the model to deploy."
  type        = string
  default     = "deepseek-r1-distill-qwen-14b"
}

variable "model_provider" {
  description = "The provider of the model."
  type        = string
  default     = "digitalocean"
}

variable "accelerator_slug" {
  description = "The slug identifier for the GPU accelerator type."
  type        = string
  default     = "gpu-h100x1-80gb"
}

variable "accelerator_scale" {
  description = "The number of accelerator units to allocate."
  type        = number
  default     = 1
}

variable "accelerator_type" {
  description = "The accelerator type."
  type        = string
  default     = "nvidia_h100"
}
