---
page_title: "DigitalOcean: digitalocean_partner_network_connect"
subcategory: "Networking"
---

# digitalocean_partner_network_connect

-> Partner Network Connect is currently in private preview. If you are not a member of the private preview group for this feature, you will not be able to use it until is has been more widely released. Please follow the official [DigitalOcean changelog](https://docs.digitalocean.com/release-notes/) for updates.

Provides a [DigitalOcean Partner Network Connect](#digitalocean_partner_network_connect) resource.

Partner Network Connects enable private connectivity between VPC networks across different cloud providers via a supported service provider.

## Example Usage

```hcl
resource "digitalocean_partner_network_connect" "foobar" {
  name                         = "example-partner-network-connect"
  connection_bandwidth_in_mbps = 100
  region                       = "nyc"
  naas_provider                = "MEGAPORT"
  vpc_ids = [
    digitalocean_vpc.vpc1.id,
    digitalocean_vpc.vpc2.id
  ]
  bgp {
    local_router_ip = "169.254.0.1/29"
    peer_router_asn = 133937
    peer_router_ip  = "169.254.0.6/29"
    auth_key        = "BGPAu7hK3y!"
  }
}
```

## Argument Reference

The following arguments are supported and are mutually exclusive:

* `name` - (Required) The name of an existing Partner Network Connect.
* `connection_bandwidth_in_mbps` - (Required) The bandwidth in megabits per second of the connection.
* `region` - (Required) The region where the Partner Network Connect is located.
* `naas_provider` - (Required) The network as a service provider for the Partner Network Connect.
* `vpc_ids` - (Required) The list of VPC IDs involved in the Partner Network Connect.
* `bgp` - (Optional) The BGP configuration for the Partner Network Connect.
    * `local_router_ip` - The local router IP address in CIDR notation.
    * `peer_router_asn` - The peer autonomous system number.
    * `peer_router_ip` - The peer router IP address in CIDR notation.
    * `auth_key` - The authentication key for the BGP session.

## Attributes Reference

* `id` - The unique identifier of an existing Partner Network Connect.
* `state` - The state of the Partner Network Connect.
* `created_at` - The date and time of when the Partner Network Connect was created.
