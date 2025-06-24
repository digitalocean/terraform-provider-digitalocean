---
page_title: "DigitalOcean: digitalocean_partner_attachment"
subcategory: "Networking"
---

# digitalocean_partner_attachment

Provides a [DigitalOcean Partner Attachment](#digitalocean_partner_attachment) resource.

Partner Attachments enable private connectivity between VPC networks across different cloud providers via a supported service provider.

## Example Usage

```hcl
resource "digitalocean_partner_attachment" "foobar" {
  name                         = "example-partner-attachment"
  connection_bandwidth_in_mbps = 1000
  region                       = "nyc"
  naas_provider                = "MEGAPORT"
  vpc_ids = [
    digitalocean_vpc.vpc1.id,
    digitalocean_vpc.vpc2.id
  ]
  bgp {
    local_router_ip = "169.254.100.1/29"
    peer_router_asn = 133937
    peer_router_ip  = "169.254.100.6/29"
    auth_key        = "BGPAu7hK3y!"
  }
}
```

## Argument Reference

The following arguments are supported and are mutually exclusive:

* `name` - (Required) The name of an existing Partner Attachment.
* `connection_bandwidth_in_mbps` - (Required) The bandwidth in megabits per second of the connection.
* `region` - (Required) The region where the Partner Attachment is located.
* `naas_provider` - (Required) The network as a service provider for the Partner Attachment.
* `vpc_ids` - (Required) The list of VPC IDs involved in the Partner Attachment.
* `bgp` - (Optional) The BGP configuration for the Partner Attachment.
    * `local_router_ip` - The local router IP address in CIDR notation.
    * `peer_router_asn` - The peer autonomous system number.
    * `peer_router_ip` - The peer router IP address in CIDR notation.
    * `auth_key` - The authentication key for the BGP session.

## Attributes Reference

* `id` - The unique identifier of an existing Partner Attachment.
* `state` - The state of the Partner Attachment.
* `created_at` - The date and time of when the Partner Attachment was created.
