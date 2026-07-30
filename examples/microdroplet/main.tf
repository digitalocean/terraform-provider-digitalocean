terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = ">= 2.49.0"
    }
  }
}

provider "digitalocean" {
  # export DIGITALOCEAN_TOKEN="Your API TOKEN"
}

# Import a MicroDroplet image from a public OCI reference. Image import is
# asynchronous; Terraform waits until the import completes before proceeding.
resource "digitalocean_microdroplet_image" "web" {
  name   = "example-nginx"
  source = "docker.io/library/nginx:latest"
}

# Provision a MicroDroplet backed by the imported image. auto_pause causes the
# MicroDroplet to pause automatically after 5 minutes of idleness, and
# auto_resume brings it back up on the next request.
resource "digitalocean_microdroplet" "web" {
  name   = "example-microdroplet"
  region = "nyc3"
  size   = "microdroplet-1"
  image  = digitalocean_microdroplet_image.web.id

  http_port     = 80
  http_protocol = "http"

  auto_pause {
    enabled      = true
    idle_timeout = "5m"
  }

  auto_resume = true

  environment = {
    LOG_LEVEL = "info"
  }

  tags = ["example", "microdroplet"]
}

# Read all checkpoints belonging to the MicroDroplet. Checkpoints are captured
# automatically by DigitalOcean each time the MicroDroplet pauses.
data "digitalocean_microdroplet_checkpoints" "web" {
  microdroplet_id = digitalocean_microdroplet.web.id
}

output "endpoint" {
  description = "Public endpoint URL for the MicroDroplet"
  value       = digitalocean_microdroplet.web.endpoint
}

output "checkpoint_ids" {
  description = "IDs of all checkpoints captured for this MicroDroplet"
  value       = [for c in data.digitalocean_microdroplet_checkpoints.web.checkpoints : c.id]
}
