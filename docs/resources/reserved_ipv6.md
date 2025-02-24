---
page_title: "DigitalOcean: digitalocean_reserved_ipv6"
subcategory: "Networking"
---

# digitalocean\_reserved_ipv6

Provides a DigitalOcean reserved IPv6 to represent a publicly-accessible static IPv6 addresses that can be mapped to one of your Droplets.

~> **NOTE:** Reserved IPv6s can be assigned to a Droplet using 
`digitalocean_reserved_ipv6_assignment` resource only.

## Example Usage

```hcl
resource "digitalocean_reserved_ipv6" "foobar" {
  region_slug = "nyc3"
}`
```

## Argument Reference

The following arguments are supported:

* `region_slug` - (Required) The region that the reserved IPv6 needs to be reserved to.


## Import

Reserved IPv6s can be imported using the `ip`, e.g.

```
terraform import digitalocean_reserved_ipv6.myip 
2409:40d0:fa:27dd:9b24:7074:7b85:eee6
```
