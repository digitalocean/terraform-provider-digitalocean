terraform {
  required_providers {
    digitalocean = {
      source = "digitalocean/digitalocean"
    }
  }
}

provider "digitalocean" {
  # Set DIGITALOCEAN_TOKEN in your environment.
}

variable "region" {
  type    = string
  default = "nyc3"
}

# A MicroVM running an OCI image. It pauses after five idle minutes and
# resumes automatically on the next HTTP request.
resource "digitalocean_microvm" "web" {
  name      = "example-microvm"
  region    = var.region
  http_port = 80

  size {
    cpu    = 2
    memory = 4096
  }

  source {
    oci_ref = "docker.io/library/nginx:latest"
  }

  auto_pause {
    enabled      = true
    idle_timeout = "5m"
  }

  tags = ["example"]
}

# Capture the memory and disk state of the running MicroVM.
resource "digitalocean_microvm_checkpoint" "web" {
  microvm_id = digitalocean_microvm.web.id
  name       = "example-checkpoint"
}

# A second MicroVM restored from the checkpoint. Region and size are
# inherited from the checkpoint.
resource "digitalocean_microvm" "restored" {
  name = "example-restored"

  source {
    checkpoint_id = digitalocean_microvm_checkpoint.web.id
  }
}

output "web_url" {
  value = "http://${digitalocean_microvm.web.urls[0].hostname}"
}

output "restored_url" {
  value = "http://${digitalocean_microvm.restored.urls[0].hostname}"
}
