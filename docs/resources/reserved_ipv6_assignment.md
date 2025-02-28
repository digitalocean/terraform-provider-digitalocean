---
page_title: "DigitalOcean: digitalocean_reserved_ipv6_assignment"
subcategory: "Networking"
---

# digitalocean\_reserved_ipv6_assignment

Provides a resource for assigning an existing DigitalOcean reserved IPv6 to a Droplet. This
makes it easy to provision reserved IPv6 addresses that are not tied to the lifecycle of your Droplet.

## Example Usage

```hcl
resource "digitalocean_reserved_ipv6" "foobar" {
  region_slug = "nyc3"
}

resource "digitalocean_droplet" "foobar" {
  image  = "ubuntu-22-04-x64"
  name   = "tf-acc-test-01"
  region = "nyc3"
  size   = "s-1vcpu-1gb"
  ipv6   = true
}

resource "digitalocean_reserved_ipv6_assignment" "foobar" {
  ip         = digitalocean_reserved_ipv6.foobar.ip
  droplet_id = digitalocean_droplet.foobar.id

  lifecycle {
    create_before_destroy = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `ip` - (Required) The reserved IPv6 to assign to the Droplet.
* `droplet_id` - (Required) The ID of Droplet that the reserved IPv6 will be assigned to.

## Import

Reserved IPv6 assignments can be imported using the reserved IPv6 itself and the `id` of
the Droplet joined with a comma. For example:

```
terraform import digitalocean_reserved_ipv6_assignment.foobar 2409:40d0:fa:27dd:9b24:7074:7b85:eee6,123456
```
