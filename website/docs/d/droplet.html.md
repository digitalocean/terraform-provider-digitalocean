---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_droplet"
sidebar_current: "docs-do-datasource-droplet"
description: |-
  Get information on a Droplet.
---

# digitalocean_droplet

Get information on a Droplet for use in other resources. This data source provides
all of the Droplet's properties as configured on your DigitalOcean account. This
is useful if the Droplet in question is not managed by Terraform or you need to
utilize any of the Droplet's data.

**Note:** This data source returns a single Droplet. When specifying a `tag`, an
error is triggered if more than one Droplet is found.

## Example Usage

Get the Droplet by name:

```hcl
data "digitalocean_droplet" "example" {
  name = "web"
}

output "droplet_output" {
  value = data.digitalocean_droplet.example.ipv4_address
}
```

Get the Droplet by tag:

```hcl
data "digitalocean_droplet" "example" {
  tag = "web"
}
```

Get the Droplet by ID:

```hcl
data "digitalocean_droplet" "example" {
  id = digitalocean_kubernetes_cluster.example.node_pool[0].nodes[0].droplet_id
}
```

## Argument Reference

One of following the arguments must be provided:

* `id` - (Optional) The ID of the Droplet
* `name` - (Optional) The name of the Droplet.
* `tag` - (Optional) A tag applied to the Droplet.

## Attributes Reference

The following attributes are exported:

* `id`: The ID of the Droplet.
* `urn` - The uniform resource name of the Droplet
* `region` - The region the Droplet is running in.
* `image` - The Droplet image ID or slug.
* `size` - The unique slug that indentifies the type of Droplet.
* `disk` - The size of the Droplets disk in GB.
* `vcpus` - The number of the Droplets virtual CPUs.
* `memory` - The amount of the Droplets memory in MB.
* `price_hourly` - Droplet hourly price.
* `price_monthly` - Droplet monthly price.
* `status` - The status of the Droplet.
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
