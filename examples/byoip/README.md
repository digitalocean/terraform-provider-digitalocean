# BYOIP (Bring Your Own IP) Example

This example demonstrates how to use BYOIP prefixes with DigitalOcean resources.

## Prerequisites

- A BYOIP prefix must be created and verified outside of Terraform (via the DigitalOcean control panel or API)
- You need the UUID of your BYOIP prefix

## What This Example Does

This example shows how to:
1. Query details about your BYOIP prefix
2. List IP addresses that are already assigned from the BYOIP prefix
3. Output information for monitoring and auditing purposes

**Note:** This example is for **querying existing resources only**. BYOIP prefixes and IP 
allocations are managed outside of Terraform through the DigitalOcean API or control panel.

## Usage

1. Query your existing BYOIP prefix:

```hcl
data "digitalocean_byoip_prefix" "example" {
  uuid = var.byoip_prefix_uuid
}
```

2. List IP addresses already assigned from the BYOIP prefix:

```hcl
data "digitalocean_byoip_addresses" "example" {
  byoip_prefix_uuid = data.digitalocean_byoip_prefix.example.uuid
}
```

3. Use the data for reporting or monitoring:

```hcl
output "byoip_summary" {
  value = {
    prefix         = data.digitalocean_byoip_prefix.example.prefix
    region         = data.digitalocean_byoip_prefix.example.region
    status         = data.digitalocean_byoip_prefix.example.status
    assigned_ips   = length(data.digitalocean_byoip_addresses.example.addresses)
  }
}
```

## Notes

- BYOIP prefixes are created outside of Terraform and are read-only in Terraform
- The prefix must be in "ACTIVE" status before IPs can be assigned to resources.
- IP addresses listed by `digitalocean_byoip_addresses` are already assigned to resources.
- To allocate new IPs from the BYOIP prefix, use `digitalocean_reserved_ip` resource.
