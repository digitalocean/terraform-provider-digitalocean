# BYOIP (Bring Your Own IP) Example

This example demonstrates how to use BYOIP prefixes with DigitalOcean resources.

## Prerequisites

- A BYOIP prefix must be created and verified outside of Terraform (via the DigitalOcean control panel or API)
- You need the UUID of your BYOIP prefix

## Usage

1. Query your existing BYOIP prefix:

```hcl
data "digitalocean_byoip_prefix" "example" {
  uuid = var.byoip_prefix_uuid
}
```

2. List IP addresses from the BYOIP prefix:

```hcl
data "digitalocean_byoip_addresses" "example" {
  byoip_prefix_uuid = data.digitalocean_byoip_prefix.example.uuid
}
```

3. Use a BYOIP IP address with a Droplet:

```hcl
resource "digitalocean_droplet" "web" {
  name   = "web-1"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = data.digitalocean_byoip_prefix.example.region
}

resource "digitalocean_reserved_ip" "byoip_ip" {
  ip_address = data.digitalocean_byoip_addresses.example.addresses[0].ip_address
  region     = data.digitalocean_byoip_prefix.example.region
  droplet_id = digitalocean_droplet.web.id
}
```

## Notes

- BYOIP prefixes are created outside of Terraform and are read-only in Terraform
- The prefix must be in a "verified" status before you can use addresses from it
- IP addresses from the BYOIP prefix can be used with `digitalocean_reserved_ip` resource
