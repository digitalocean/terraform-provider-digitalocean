---
page_title: "DigitalOcean: digitalocean_partner_network_connect"
subcategory: "Networking"
---

# digitalocean_partner_network_connect

-> Partner Network Connect is currently in private preview. If you are not a member of the private preview group for this feature, you will not be able to use it until is has been more widely released. Please follow the official [DigitalOcean changelog](https://docs.digitalocean.com/release-notes/) for updates.

Retrieve information about a Partner Network Connect for use in other resources.

This data source provides all of the Partner Network Connect's properties as configured on your
DigitalOcean account. This is useful if the Partner Network Connect in question is not managed by
Terraform or you need to utilize any of the Partner Network Connect's data.

Partner Network Connect may be looked up by `id` or `name`.

## Example Usage

### Partner Network Connect By Id

```hcl
data "digitalocean_partner_network_connect" "example" {
  id = "example-id"
}
```

### Partner Network Connect By Name

```hcl
data "digitalocean_partner_network_connect" "example" {
  name = "example-pia"
}
```

## Argument Reference

The following arguments are supported and are mutually exclusive:

* `id` - The unique identifier of an existing Partner Network Connect.
* `name` - The name of an existing Partner Network Connect.

## Attributes Reference

* `id` - The unique identifier of an existing Partner Network Connect.
* `name` - The name of the Partner Network Connect.
* `connection_bandwidth_in_mbps` - The bandwidth in megabits per second of the connection.
* `region` - The region where the Partner Network Connect is located.
* `naas_provider` - The network as a service provider for the Partner Network Connect.
* `vpc_ids` - The list of VPC IDs involved in the partner network connect.
* `bgp` - The BGP configuration for the Partner Network Connect.
    * `local_router_ip` - The local router IP address in CIDR notation.
    * `peer_router_asn` - The peer autonomous system number.
    * `peer_router_ip` - The peer router IP address in CIDR notation.
    * `auth_key` - The authentication key for the BGP session.
* `state` - The state of the Partner Network Connect.
* `created_at` - The date and time of when the Partner Network Connect was created.
