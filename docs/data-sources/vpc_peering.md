---
page_title: "DigitalOcean: digitalocean_vpc_peering"
---

# digitalocean_vpc_peering

> Note: VPC peering is currently in alpha. If you are not a member of the alpha group for this feature, you will not be able to use it until it has been more widely released. Please follow the official [DigitalOcean changelog](https://docs.digitalocean.com/release-notes/) for updates.

Retrieve information about a VPC Peering for use in other resources.

This data source provides all of the VPC Peering's properties as configured on your 
DigitalOcean account. This is useful if the VPC Peering in question is not managed by 
Terraform or you need to utilize any of the VPC Peering's data.

VPC Peerings may be looked up by `id` or `name`.

## Example Usage

### VPC Peering By Id

```hcl
data "digitalocean_vpc_peering" "example" {
  id = "example-id"
}
```

Reuse the data about a VPC Peering in other resources:

```hcl
data "digitalocean_vpc_peering" "example" {
  id = "example-id"
}

resource "digitalocean_droplet" "example" {
  name     = "example-01"
  size     = "s-1vcpu-1gb"
  image    = "ubuntu-18-04-x64"
  region   = "nyc3"
  vpc_uuid = data.digitalocean_vpc_peering.example.vpc_ids[0]
}
```

### VPC Peering By Name

```hcl
data "digitalocean_vpc_peering" "example" {
  name = "example-peering"
}
```

Reuse the data about a VPC Peering in other resources:

```hcl
data "digitalocean_vpc_peering" "example" {
  name = "example-peering"
}

resource "digitalocean_droplet" "example" {
  name     = "example-01"
  size     = "s-1vcpu-1gb"
  image    = "ubuntu-18-04-x64"
  region   = "nyc3"
  vpc_uuid = data.digitalocean_vpc_peering.example.vpc_ids[0]
}
```

## Argument Reference

The following arguments are supported and are mutually exclusive:

* `id` - The unique identifier of an existing VPC Peering.
* `name` - The name of an existing VPC Peering.

## Attributes Reference

* `id` - The unique identifier for the VPC Peering.
* `name` - The name of the VPC Peering.
* `vpc_ids` - The list of VPC IDs involved in the peering.
* `status` - The status of the VPC Peering.
* `created_at` - The date and time of when the VPC Peering was created.
