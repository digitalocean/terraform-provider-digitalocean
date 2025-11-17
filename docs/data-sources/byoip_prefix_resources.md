---
page_title: "DigitalOcean: digitalocean_byoip_prefix_resources"
subcategory: "Networking"
---

# digitalocean_byoip_prefix_resources

Get information about IP addresses that have been **already assigned** from a 
BYOIP (Bring Your Own IP) prefix. This data source provides a list of all IP addresses 
that are currently assigned to resources from a specific BYOIP prefix.

**Note:** This data source only lists IPs that are already assigned to resources (like Droplets
or Load Balancers). To allocate new IPs from the BYOIP prefix, you need to use `digitalocean_reserved_ip` resource.

## Example Usage

List all assigned IP addresses from a BYOIP prefix:

```hcl
data "digitalocean_byoip_prefix" "example" {
  uuid = "506f78a4-e098-11e5-ad9f-000f53306ae1"
}

data "digitalocean_byoip_prefix_resources" "example" {
  byoip_prefix_uuid = data.digitalocean_byoip_prefix.example.uuid
}

# Output the assigned IPs
output "assigned_byoip_ips" {
  value = [
    for addr in data.digitalocean_byoip_prefix_resources.example.addresses : {
      ip       = addr.ip_address
      region   = addr.region
      assigned = addr.assigned_at
    }
  ]
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
