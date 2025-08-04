---
page_title: "DigitalOcean: digitalocean_vpc_nat_gateway"
subcategory: "Networking"
---

# digitalocean\_vpc\_nat\_gateway

Provides a DigitalOcean VPC NAT Gateway resource. This can be used to create, modify, 
read and delete VPC NAT Gateways.

NOTE: VPC NAT Gateway is currently in Private Preview.

## Example Usage

```hcl
resource "digitalocean_vpc" "my-vpc" {
  name   = "terraform-example"
  region = "nyc3"
}

resource "digitalocean_vpc_nat_gateway" "my-vpc-nat-gateway" {
  name   = "terraform-example"
  type   = "PUBLIC"
  region = "nyc3"
  size   = "1"
  vpcs {
    vpc_uuid = digitalocean_vpc.my-vpc.id
  }
  udp_timeout_seconds  = 30
  icmp_timeout_seconds = 30
  tcp_timeout_seconds  = 30
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the VPC NAT Gateway.
* `type` - (Required) The type of the VPC NAT Gateway.
* `region` - (Required) The region for the VPC NAT Gateway.
* `size` - (Required) The size of the VPC NAT Gateway.
* `vpcs` - (Required) The ingress VPC configuration of the VPC NAT Gateway, the supported arguments are
documented below.
* `udp_timeout_seconds` - The egress timeout value for UDP connections of the VPC NAT Gateway.
* `icmp_timeout_seconds` - The egress timeout value for ICMP connections of the VPC NAT Gateway.
* `tcp_timeout_seconds` - The egress timeout value for TCP connections of the VPC NAT Gateway.

`vpcs` supports the following attributes:

* `vpc_uuid` - The ID of the ingress VPC
* `gateway_ip` - (Read-only) The private IP of the VPC NAT Gateway
* `default_gateway` - Boolean flag indicating if this should be the default gateway in this VPC

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the VPC NAT Gateway.
* `egresses` - Embeds the list of public egresses assigned to the VPC NAT Gateway: resolves as list of
`public_gateways` embedding the reserved `ipv4` addresses.
* `created_at` - Created at timestamp for the VPC NAT Gateway.
* `updated_at` - Updated at timestamp for the VPC NAT Gateway.

## Import

VPC NAT Gateways can be imported using their `id`, e.g.

```
terraform import digitalocean_vpc_nat_gateway.my-vpc-nat-gateway-id 38e66834-d741-47ec-88e7-c70cbdcz0445
```
