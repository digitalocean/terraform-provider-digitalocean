---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_droplet"
sidebar_current: "docs-do-datasource-droplet"
description: |-
  Get information on a droplet.
---

# digitalocean_droplet

Get information on a droplet for use in other resources. This data source provides all of
the droplets properties as configured on your Digital Ocean account.
This is useful if the droplet in question is not managed by Terraform or you need to utilize
any of the droplets data.

An error is triggered if the provided droplet name does not exist.

## Example Usage

Get the droplet:

```hcl
data "digitalocean_droplet" "example" {
  name = "web"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of droplet.

## Attributes Reference

The following attributes are exported:

* `id`: The ID of the droplet.
* `region` - The region the droplet is running in.
* `image` - The Droplet image ID or slug.
* `size` - The unique slug that indentifies the type of Droplet.
* `disk` - The size of the Droplets disk in GB.
* `vcpus` - The number of the Droplets virtual CPUs.
* `memory` - The amount of the Droplets memory in MB.
* `price_hourly` - Droplet hourly price.
* `price_monthly` - Droplet monthly price.
* `status` - The status of the droplet.
* `locked` - Whether the Droplet is locked.
* `ipv6_address` - The Droplets public IPv6 address
* `ipv6_address_private` - The Droplets private IPv6 address
* `ipv4_address` - The Droplets public IPv4 address
* `ipv4_address_private` - The Droplets private IPv4 address
* `backups` - Whether backups are enabled.
* `ipv6` - Whether IPv6 is enabled.
* `private_networking` - Whether private networks are enabled.
* `monitoring` - Whether monitoring agent is installed.
* `volume_ids` - List of the IDs of each volumes attached to the Droplet.
* `tags` - A list of the tags associated to the Droplet. 
