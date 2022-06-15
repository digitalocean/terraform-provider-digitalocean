---
page_title: "DigitalOcean: digitalocean_reserved_ip_assignment"
---

# digitalocean\_reserved_ip_assignment

Provides a resource for assigning an existing DigitalOcean reserved IP to a Droplet. This
makes it easy to provision reserved IP addresses that are not tied to the lifecycle of your
Droplet.

## Example Usage

```hcl
resource "digitalocean_reserved_ip" "example" {
  region = "nyc3"
}

resource "digitalocean_droplet" "example" {
  name               = "baz"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}

resource "digitalocean_reserved_ip_assignment" "example" {
  ip_address = digitalocean_reserved_ip.example.ip_address
  droplet_id = digitalocean_droplet.example.id
}
```

## Argument Reference

The following arguments are supported:

* `ip_address` - (Required) The reserved IP to assign to the Droplet.
* `droplet_id` - (Optional) The ID of Droplet that the reserved IP will be assigned to.

## Import

Reserved IP assignments can be imported using the reserved IP itself and the `id` of
the Droplet joined with a comma. For example:

```
terraform import digitalocean_reserved_ip_assignment.foobar 192.0.2.1,123456
```
