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


resource "digitalocean_droplet" "web" {
  image  = "ubuntu-20-04-x64"
  name   = "web2"
  region = "nyc2"
  size   = "s-1vcpu-1gb"
}
