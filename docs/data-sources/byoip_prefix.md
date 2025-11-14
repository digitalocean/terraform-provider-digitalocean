---
page_title: "DigitalOcean: digitalocean_byoip_prefix"
subcategory: "Networking"
---

# digitalocean_byoip_prefix

Get information on a BYOIP (Bring Your Own IP) prefix. This data source provides the
prefix CIDR, region, advertisement status, and current state as configured on your
DigitalOcean account. This is useful if the BYOIP prefix in question is not managed
by Terraform or you need to utilize any of the prefix's data.

**Note:** BYOIP prefixes are created and managed outside of Terraform through the
DigitalOcean control panel or API. This data source is read-only.

An error is triggered if the provided BYOIP prefix UUID does not exist.

## Example Usage

Get the BYOIP prefix:

```hcl
data "digitalocean_byoip_prefix" "example" {
  uuid = "506f78a4-e098-11e5-ad9f-000f53306ae1"
}
```

List assigned IP addresses from a BYOIP prefix:

```hcl
data "digitalocean_byoip_prefix" "example" {
  uuid = "506f78a4-e098-11e5-ad9f-000f53306ae1"
}

data "digitalocean_byoip_addresses" "example" {
  byoip_prefix_uuid = data.digitalocean_byoip_prefix.example.uuid
}

# Output information about the BYOIP prefix and its assigned IPs
output "byoip_info" {
  value = {
    prefix          = data.digitalocean_byoip_prefix.example.prefix
    region          = data.digitalocean_byoip_prefix.example.region
    status          = data.digitalocean_byoip_prefix.example.status
    assigned_count  = length(data.digitalocean_byoip_addresses.example.addresses)
  }
}
```

## Argument Reference

The following arguments are supported:

* `uuid` - (Required) The UUID of the BYOIP prefix.

## Attributes Reference

The following attributes are exported:

* `id`: The UUID of the BYOIP prefix.
* `uuid`: The UUID of the BYOIP prefix.
* `prefix`: The CIDR notation of the prefix (e.g., "192.0.2.0/24").
* `region`: The DigitalOcean region where the prefix is deployed.
* `advertised`: A boolean indicating whether the prefix is currently being advertised.
* `status`: The current status of the BYOIP prefix (e.g., "verified", "pending", "failed").
* `failure_reason`: The reason for failure if the status is "failed".
