---
page_title: "DigitalOcean: digitalocean_reserved_ip"
---

# digitalocean\_reserved_ip

Provides a DigitalOcean reserved IP to represent a publicly-accessible static IP addresses that can be mapped to one of your Droplets.

~> **NOTE:** Reserved IPs can be assigned to a Droplet either directly on the `digitalocean_reserved_ip` resource by setting a `droplet_id` or using the `digitalocean_reserved_ip_assignment` resource, but the two cannot be used together.

## Example Usage

```hcl
resource "digitalocean_droplet" "example" {
  name               = "example"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}

resource "digitalocean_reserved_ip" "example" {
  droplet_id = digitalocean_droplet.example.id
  region     = digitalocean_droplet.example.region
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required) The region that the reserved IP is reserved to.
* `droplet_id` - (Optional) The ID of Droplet that the reserved IP will be assigned to.

## Attributes Reference

The following attributes are exported:

* `ip_address` - The IP Address of the resource
* `urn` - The uniform resource name of the reserved ip

## Import

Reserved IPs can be imported using the `ip`, e.g.

```
terraform import digitalocean_reserved_ip.myip 192.168.0.1
```
