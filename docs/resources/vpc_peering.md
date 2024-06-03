---
page_title: "DigitalOcean: digitalocean_vpc_peering"
---

# digitalocean_vpc_peering

-> VPC peering is currently in alpha. If you are not a member of the alpha group for this feature, you will not be able to use it until it has been more widely released. Please follow the official [DigitalOcean changelog](https://docs.digitalocean.com/release-notes/) for updates.

Provides a [DigitalOcean VPC Peering](#digitalocean_vpc_peering) resource.

VPC Peerings are used to connect two VPC networks allowing resources in each 
VPC to communicate with each other privately.

## Example Usage

```hcl
resource "digitalocean_vpc_peering" "example" {
  name = "example-peering"
  vpc_ids = [
    digitalocean_vpc.vpc1.id,
    digitalocean_vpc.vpc2.id
  ]
}
```

### Resource Assignement

You can use the VPC Peering resource to allow communication between resources
in different VPCs. For example:

```hcl
resource "digitalocean_vpc" "vpc1" {
  name   = "vpc1"
  region = "nyc3"
}

resource "digitalocean_vpc" "vpc2" {
  name   = "vpc2"
  region = "nyc3"
}

resource "digitalocean_vpc_peering" "example" {
  name = "example-peering"
  vpc_ids = [
    digitalocean_vpc.vpc1.id,
    digitalocean_vpc.vpc2.id
  ]
}

resource "digitalocean_droplet" "example1" {
  name     = "example1"
  size     = "s-1vcpu-1gb"
  image    = "ubuntu-18-04-x64"
  region   = "nyc3"
  vpc_uuid = digitalocean_vpc.vpc1.id
}

resource "digitalocean_droplet" "example2" {
  name     = "example2"
  size     = "s-1vcpu-1gb"
  image    = "ubuntu-18-04-x64"
  region   = "nyc3"
  vpc_uuid = digitalocean_vpc.vpc2.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A name for the VPC Peering. Must be unique and contain alphanumeric characters, dashes, and periods only.
* `vpc_ids` - (Required) A set of two VPC IDs to be peered.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The unique identifier for the VPC Peering.
* `status` - The status of the VPC Peering.
* `created_at` - The date and time of when the VPC Peering was created.

## Import

A VPC Peering can be imported using its `id`, e.g.

```
terraform import digitalocean_vpc_peering.example 771ad360-c017-4b4e-a34e-73934f5f0190
```
