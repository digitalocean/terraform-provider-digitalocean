terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = ">= 2.8.0"
    }
  }
}

# Set the variable value in a .tfvars file or using -var="do_token=..."
# The token can also be set via the DIGITALOCEAN_TOKEN environment variable.
variable "do_token" {
  type = string
}

provider "digitalocean" {
  token = var.do_token
}

# Look up available GPU sizes and model configs
data "digitalocean_dedicated_inference_sizes" "available" {}

data "digitalocean_dedicated_inference_gpu_model_config" "available" {}

# Create a dedicated inference endpoint
resource "digitalocean_dedicated_inference" "example" {
  name   = var.name
  region = var.region

  model_deployments {
    model_slug     = var.model_slug
    model_provider = var.model_provider

    accelerators {
      accelerator_slug = var.accelerator_slug
      scale            = var.accelerator_scale
      type             = var.accelerator_type
    }
  }
}

# Create an API token for the endpoint
resource "digitalocean_dedicated_inference_token" "example" {
  dedicated_inference_id = digitalocean_dedicated_inference.example.id
  name                   = "example-token"
}
