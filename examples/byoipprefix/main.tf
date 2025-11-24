# Query the existing BYOIP prefix
data "digitalocean_byoip_prefix" "example" {
  uuid = var.byoip_prefix_uuid
}

# List all IP addresses already assigned from the BYOIP prefix
data "digitalocean_byoip_prefix_resources" "example" {
  byoip_prefix_uuid = data.digitalocean_byoip_prefix.example.uuid
}

# Create a Droplet in the same region as the BYOIP prefix
resource "digitalocean_droplet" "web" {
  name   = var.droplet_name
  size   = var.droplet_size
  image  = "ubuntu-22-04-x64"
  region = data.digitalocean_byoip_prefix.example.region
}

# Assign a BYOIP IP address to the Droplet
resource "digitalocean_reserved_ip_assignment" "byoip_ip" {
  ip_address = "192.0.2.2"
  droplet_id = digitalocean_droplet.web.id
}
