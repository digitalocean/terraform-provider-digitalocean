---
page_title: "DigitalOcean: digitalocean_vpc_nat_gateway"
subcategory: "Networking"
---

# digitalocean\_vpc\_nat\_gateway

Get information on a VPC NAT Gateway for use with other managed resources  This datasource provides all the VPC
NAT Gateway properties as configured on the DigitalOcean account. This is useful if the VPC NAT Gateway in question
is not managed by Terraform, or any of the relevant data would need to be referenced in other managed resources.

## Example Usage

Get the VPC NAT Gateway by name:

```hcl
data "digitalocean_vpc_nat_gateway" "my-imported-vpc-nat-gateway" {
  name = digitalocean_vpc_nat_gateway.my-existing-vpc-nat-gateway.name
}
```

Get the VPC NAT Gateway by ID:

```hcl
data "digitalocean_vpc_nat_gateway" "my-imported-vpc-nat-gateway" {
  id = digitalocean_vpc_nat_gateway.my-existing-vpc-nat-gateway.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of VPC NAT Gateway.
* `id` - (Optional) The ID of VPC NAT Gateway.

## Attributes Reference

See the [VPC NAT Gateway Resource](../resources/vpc_nat_gateway.md) for details on the
returned attributes - they are identical.
