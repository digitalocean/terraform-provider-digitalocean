---
page_title: "DigitalOcean: digitalocean_vpc"
---

# digitalocean_vpc

Provides a [DigitalOcean VPC](https://developers.digitalocean.com/documentation/v2/#vpcs) resource.

VPCs are virtual networks containing resources that can communicate with each
other in full isolation, using private IP addresses.

## Example Usage

```hcl
resource "digitalocean_vpc" "example" {
  name     = "example-project-network"
  region   = "nyc3"
  ip_range = "10.10.10.0/24"
}
```

### Resource Assignment

`digitalocean_droplet`, `digitalocean_kubernetes_cluster`,
`digitalocean_load_balancer`, and `digitalocean_database_cluster` resources
may be assigned to a VPC by referencing its `id`. For example:

```hcl
resource "digitalocean_vpc" "example" {
  name     = "example-project-network"
  region   = "nyc3"
}

resource "digitalocean_droplet" "example" {
  name     = "example-01"
  size     = "s-1vcpu-1gb"
  image    = "ubuntu-18-04-x64"
  region   = "nyc3"
  vpc_uuid = digitalocean_vpc.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A name for the VPC. Must be unique and contain alphanumeric characters, dashes, and periods only.
* `region` - (Required) The DigitalOcean region slug for the VPC's location.
* `description` - (Optional) A free-form text field up to a limit of 255 characters to describe the VPC.
* `ip_range` - (Optional) The range of IP addresses for the VPC in CIDR notation. Network ranges cannot overlap with other networks in the same account and must be in range of private addresses as defined in RFC1918. It may not be larger than `/16` or smaller than `/24`.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The unique identifier for the VPC.
* `urn` - The uniform resource name (URN) for the VPC.
* `default` - A boolean indicating whether or not the VPC is the default one for the region.
* `created_at` - The date and time of when the VPC was created.

## Import

A VPC can be imported using its `id`, e.g.

```
terraform import digitalocean_vpc.example 506f78a4-e098-11e5-ad9f-000f53306ae1
```
