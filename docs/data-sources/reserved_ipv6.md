---
page_title: "DigitalOcean: digitalocean_reserved_ipv6"
subcategory: "Networking"
---

# digitalocean_reserved_ipv6

Get information on a reserved IPv6. This data source provides the region_slug and droplet id as configured on your DigitalOcean account. This is useful if the reserved IPv6 in question is not managed by Terraform or you need to find the Droplet the IP is
attached to.

An error is triggered if the provided reserved IPv6 does not exist.

## Example Usage

Get the reserved IPv6:

```hcl

resource "digitalocean_reserved_ipv6" "foo" {
  region_slug = "nyc3"
}

data "digitalocean_reserved_ipv6" "foobar" {
  ip = digitalocean_reserved_ipv6.foo.ip
}

```

## Argument Reference

The following arguments are supported:

* `ip` - (Required) The allocated IPv6 address of the specific reserved IPv6 to retrieve.

## Attributes Reference

The following attributes are exported:

* `region_slug`: The region that the reserved IPv6 is reserved to.
* `urn`: The uniform resource name of the reserved IPv6.
* `droplet_id`: The Droplet id that the reserved IP has been assigned to.
