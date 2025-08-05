terraform {
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = ">= 2.44.1"
    }
  }
}

provider "digitalocean" {
  token = "DO_TOKEN"
}


resource "digitalocean_vpc" "example_vpc" {
  name     = "vpc-1"
  region   = "nyc3"
}

resource "digitalocean_partner_attachment" "parent_pia" {
  name                         = "parent-partner-attachment"
  connection_bandwidth_in_mbps = 1000
  region                       = "nyc"
  naas_provider                = "MEGAPORT"
  vpc_ids = [
    digitalocean_vpc.example_vpc.id
  ]
}

resource "digitalocean_partner_attachment" "child_pia" {
  name                         = "child-partner-attachment"
  connection_bandwidth_in_mbps = 1000
  region                       = "nyc"
  naas_provider                = "MEGAPORT"
  vpc_ids = [
    digitalocean_vpc.example_vpc.id
  ]
  parent_uuid = digitalocean_partner_attachment.parent_pia.id
}