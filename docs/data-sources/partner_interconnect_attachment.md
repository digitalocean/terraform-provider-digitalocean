---
page_title: "DigitalOcean: digitalocean_partner_interconnect_attachment"
subcategory: "Networking"
---

# digitalocean_partner_interconnect_attachment

-> Partner Interconnect Attachment is currently in private preview. If you are not a member of the private preview group for this feature, you will not be able to use it until is has been more widely released. Please follow the official [DigitalOcean changelog](https://docs.digitalocean.com/release-notes/) for updates.

Retrieve information about a Partner Interconnect Attachment for use in nother resources.

This data source provides all of the Partner Interconnect Attachment's properties as configured on your
DigitalOcean account. This is useful if the Partner Interconnect Attachment in question is not managed by
Terraform or you need to utilize any of the Partner Interconnect Attachment's data.

Partner Interconnect Attachments may be looked up by `id` or `name`.

## Example Usage

### Partner Interconnect Attachment By Id

```hcl
data "digitalocean_partner_interconnect_attachment" "example" {
  id = "example-id"
}
```

### Partner Interconnect Attachment By Name

```hcl
data "digitalocean_partner_interconnect_attachment" "example" {
  name = "example-pia"
}
```

## Argument Reference

The following arguments are supported and are mutually exclusive:

* `id` - The unique identifier of an existing Partner Interconnect Attachment.
* `name` - The name of an existing Partner Interconnect Attachment.

## Attributes Reference

* `id` - The unique identifier of an existing Partner Interconnect Attachment.
* `name` - The name of the Partner Interconnect Attachment.
* `connection_bandwidth_in_mbps` - The bandwidth in megabits per second of the connection.
* `region` - The region where the Partner Interconnect Attachment is located.
* `naas_provider` - The network as a service provider for the Partner Interconnect Attachment.
* `vpc_ids` - The list of VPC IDs involved in the partner interconnect.
* `bgp` - The BGP configuration for the Partner Interconnect Attachment.
    * `local_router_ip` - The local router IP address in CIDR notation.
    * `peer_router_asn` - The peer autonomous system number.
    * `peer_router_ip` - The peer router IP address in CIDR notation.
    * `auth_key` - The authentication key for the BGP session.
* `state` - The state of the Partner Interconnect Attachment.
* `created_at` - The date and time of when the Partner Interconnect Attachment was created.
