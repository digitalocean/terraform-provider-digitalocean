# BYOIP (Bring Your Own IP) Example

This example demonstrates how to use BYOIP prefixes with DigitalOcean resources.

## Prerequisites

- A BYOIP prefix can be created using the `digitalocean_byoip_prefix` resource, or outside of Terraform (via the DigitalOcean control panel or API)
- You need the UUID of your BYOIP prefix

## What This Example Does

This example shows how to:
1. Create a BYOIP prefix resource (optional)
2. Query details about your BYOIP prefix
3. List IP addresses that are already assigned from the BYOIP prefix
4. Output information for monitoring and auditing purposes

## Usage

1. (Optional) Create a new BYOIP prefix:

See [DigitalOcean's BYOIP provisioning guide](https://docs.digitalocean.com/products/networking/reserved-ips/how-to/provision-byoip/) for instructions on provisioning BYOIP Prefix.

```hcl
resource "digitalocean_byoip_prefix" "example" {
  prefix      = "192.0.2.0/24"
  signature   = var.prefix_signature  # Required: cryptographic signature proving ownership
  region      = "nyc3"
  advertised  = false  # Optional: defaults to false
}
```

2. Query your existing BYOIP prefix:

```hcl
data "digitalocean_byoip_prefix" "example" {
  uuid = var.byoip_prefix_uuid
}
```

3. List IP addresses already assigned from the BYOIP prefix:

```hcl
data "digitalocean_byoip_prefix_resources" "example" {
  byoip_prefix_uuid = data.digitalocean_byoip_prefix.example.uuid
}
```

4. Use the data for reporting or monitoring:

```hcl
output "byoip_summary" {
  value = {
    prefix         = data.digitalocean_byoip_prefix.example.prefix
    region         = data.digitalocean_byoip_prefix.example.region
    status         = data.digitalocean_byoip_prefix.example.status
    assigned_ips   = length(data.digitalocean_byoip_prefix_resources.example.addresses)
  }
}
```

## Notes

- BYOIP prefixes can be created using the `digitalocean_byoip_prefix` resource or managed outside of Terraform
- The prefix must be in "ACTIVE" status before IPs can be assigned to resources.
- IP addresses listed by `digitalocean_byoip_prefix_resources` are already assigned to resources.
- To allocate new IPs from the BYOIP prefix, use `digitalocean_reserved_ip_assignment` resource.
