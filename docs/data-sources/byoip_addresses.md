---
page_title: "DigitalOcean: digitalocean_byoip_addresses"
subcategory: "Networking"
---

# digitalocean_byoip_addresses

Get information about IP addresses allocated from a BYOIP (Bring Your Own IP) prefix.
This data source provides a list of all IP addresses that have been allocated from a
specific BYOIP prefix.

This is useful when you need to reference existing BYOIP-allocated IP addresses or
manage resources that depend on BYOIP addresses.

## Example Usage

Get all addresses from a BYOIP prefix:

```hcl
data "digitalocean_byoip_prefix" "example" {
  uuid = "506f78a4-e098-11e5-ad9f-000f53306ae1"
}

data "digitalocean_byoip_addresses" "example" {
  byoip_prefix_uuid = data.digitalocean_byoip_prefix.example.uuid
}

# Use a BYOIP address with a reserved IP
resource "digitalocean_reserved_ip" "example" {
  ip_address = data.digitalocean_byoip_addresses.example.addresses[0].ip_address
  region     = data.digitalocean_byoip_prefix.example.region
}
```

## Argument Reference

The following arguments are supported:

* `byoip_prefix_uuid` - (Required) The UUID of the BYOIP prefix to list addresses from.

## Attributes Reference

The following attributes are exported:

* `id`: The UUID of the BYOIP prefix.
* `addresses`: A list of IP addresses allocated from the BYOIP prefix. Each address has the following attributes:
  * `id`: The unique identifier of the IP address allocation.
  * `ip_address`: The IP address.
  * `region`: The region where the IP is allocated.
  * `assigned_at`: The timestamp when the IP was assigned.
