---
page_title: "DigitalOcean: digitalocean_partner_attachment"
subcategory: "Networking"
---

# digitalocean_partner_attachment

Retrieve information about a Partner Attachment for use in other resources.

This data source provides all of the Partner Attachment's properties as configured on your
DigitalOcean account. This is useful if the Partner Attachment in question is not managed by
Terraform or you need to utilize any of the Partner Attachment's data.

Partner Attachment may be looked up by `id` or `name`.

## Example Usage

### Partner Attachment By Id

```hcl
data "digitalocean_partner_attachment" "example" {
  id = "example-id"
}
```

### Partner Attachment By Name

```hcl
data "digitalocean_partner_attachment" "example" {
  name = "example-pia"
}
```

## Argument Reference

The following arguments are supported and are mutually exclusive:

* `id` - The unique identifier of an existing Partner Attachment.
* `name` - The name of an existing Partner Attachment.

## Attributes Reference

* `id` - The unique identifier of an existing Partner Attachment.
* `name` - The name of the Partner Attachment.
* `connection_bandwidth_in_mbps` - The bandwidth in megabits per second of the connection.
* `region` - The region where the Partner Attachment is located.
* `naas_provider` - The network as a service provider for the Partner Attachment.
* `vpc_ids` - The list of VPC IDs involved in the Partner Attachment.
* `bgp` - The BGP configuration for the Partner Attachment.
    * `local_router_ip` - The local router IP address in CIDR notation.
    * `peer_router_asn` - The peer autonomous system number.
    * `peer_router_ip` - The peer router IP address in CIDR notation.
    * `auth_key` - The authentication key for the BGP session.
* `state` - The state of the Partner Attachment.
* `created_at` - The date and time of when the Partner Attachment was created.
* `redundancy_zone` - The redundancy zone of the Partner Attachment.
* `parent_uuid` - The UUID of the parent Partner Attachment, if applicable.
* `children` - The UUIDs of the child Partner Attachment, if applicable.
