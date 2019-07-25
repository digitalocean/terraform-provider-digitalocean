---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_floating_ip"
sidebar_current: "docs-do-resource-floating-ip"
description: |-
  Provides a DigitalOcean Floating IP resource.
---

# digitalocean\_floating_ip

Provides a DigitalOcean Floating IP to represent a publicly-accessible static IP addresses that can be mapped to one of your Droplets.

~> **NOTE:** Floating IPs can be assigned to a Droplet either directly on the `digitalocean_floating_ip` resource by setting a `droplet_id` or using the `digitalocean_floating_ip_assignment` resource, but the two cannot be used together.

## Example Usage

```hcl
resource "digitalocean_droplet" "foobar" {
  name               = "baz"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-18-04-x64"
  region             = "sgp1"
  ipv6               = true
  private_networking = true
}

resource "digitalocean_floating_ip" "foobar" {
  droplet_id = digitalocean_droplet.foobar.id
  region     = digitalocean_droplet.foobar.region
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required) The region that the Floating IP is reserved to.
* `droplet_id` - (Optional) The ID of Droplet that the Floating IP will be assigned to.

## Attributes Reference

The following attributes are exported:

* `ip_address` - The IP Address of the resource
* `urn` - The uniform resource name of the floating ip

## Import

Floating IPs can be imported using the `ip`, e.g.

```
terraform import digitalocean_floating_ip.myip 192.168.0.1
```
