---
page_title: "DigitalOcean: digitalocean_firewall"
---

# digitalocean_firewall

Get information on a DigitalOcean Firewall.

## Example Usage

Get the firewall:

```hcl
data "digitalocean_firewall" "example" {
  firewall_id = "1df48973-6eef-4214-854f-fa7726e7e583"
}

output "example_firewall_name" {
  value = data.digitalocean_firewall.example.name
}
```

## Argument Reference

* `firewall_id` - (Required) The ID of the firewall to retrieve information
  about.

## Attributes Reference

The following attributes are exported:

* `id` - A unique ID that can be used to identify and reference a Firewall.
* `status` - A status string indicating the current state of the Firewall.
  This can be "waiting", "succeeded", or "failed".
* `created_at` - A time value given in ISO8601 combined date and time format
  that represents when the Firewall was created.
* `pending_changes` - A set of object containing the fields, `droplet_id`,
  `removing`, and `status`.  It is provided to detail exactly which Droplets
  are having their security policies updated.  When empty, all changes
  have been successfully applied.
* `name` - The name of the Firewall.
* `droplet_ids` - The list of the IDs of the Droplets assigned to
  the Firewall.
* `tags` - The names of the Tags assigned to the Firewall.
* `inbound_rules` - The inbound access rule block for the Firewall.
* `outbound_rules` - The outbound access rule block for the Firewall.

`inbound_rule` supports the following:

* `protocol` - The type of traffic to be allowed.
  This may be one of "tcp", "udp", or "icmp".
* `port_range` - The ports on which traffic will be allowed
  specified as a string containing a single port, a range (e.g. "8000-9000"),
  or "1-65535" to open all ports for a protocol. Required for when protocol is
  `tcp` or `udp`.
* `source_addresses` - An array of strings containing the IPv4
  addresses, IPv6 addresses, IPv4 CIDRs, and/or IPv6 CIDRs from which the
  inbound traffic will be accepted.
* `source_droplet_ids` - An array containing the IDs of
  the Droplets from which the inbound traffic will be accepted.
* `source_tags` - A set of names of Tags corresponding to group of
  Droplets from which the inbound traffic will be accepted.
* `source_load_balancer_uids` - An array containing the IDs
  of the Load Balancers from which the inbound traffic will be accepted.

`outbound_rule` supports the following:

* `protocol` - The type of traffic to be allowed.
  This may be one of "tcp", "udp", or "icmp".
* `port_range` - The ports on which traffic will be allowed
  specified as a string containing a single port, a range (e.g. "8000-9000"),
  or "1-65535" to open all ports for a protocol. Required for when protocol is
  `tcp` or `udp`.
* `destination_addresses` - An array of strings containing the IPv4
  addresses, IPv6 addresses, IPv4 CIDRs, and/or IPv6 CIDRs to which the
  outbound traffic will be allowed.
* `destination_droplet_ids` - An array containing the IDs of
  the Droplets to which the outbound traffic will be allowed.
* `destination_tags` - An array containing the names of Tags
  corresponding to groups of Droplets to which the outbound traffic will
  be allowed.
  traffic.
* `destination_load_balancer_uids` - An array containing the IDs
  of the Load Balancers to which the outbound traffic will be allowed.
