---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_droplets"
sidebar_current: "docs-do-datasource-droplets"
description: |-
  Retrieve information on Droplets.
---

# digitalocean_droplets

Get information on Droplets for use in other resources, with the ability to filter and sort the results.
If no filters are specified, all Droplets will be returned.

This data source is useful if the Droplets in question are not managed by Terraform or you need to
utilize any of the Droplets' data.

Note: You can use the [`digitalocean_droplet`](droplet) data source to obtain metadata
about a single Droplet if you already know the `id`, unique `name`, or unique `tag` to retrieve.

## Example Usage

Use the `filter` block with a `key` string and `values` list to filter images.

For example to find all Droplets with size `s-1vcpu-1gb`:

```hcl
data "digitalocean_droplets" "small" {
  filter {
    key = "size"
    values = ["s-1vcpu-1gb"]
  }
}
```

You can filter on multiple fields and sort the results as well:

```hcl
data "digitalocean_droplets" "small-with-backups" {
  filter {
    key = "size"
    values = ["s-1vcpu-1gb"]
  }
  filter {
    key = "backups"
    values = ["true"]
  }
  sort {
    key = "created_at"
    direction = "desc"
  }
}
```

## Argument Reference

* `filter` - (Optional) Filter the results.
  The `filter` block is documented below.

* `sort` - (Optional) Sort the results.
  The `sort` block is documented below.

`filter` supports the following arguments:

* `key` - (Required) Filter the Droplets by this key. This may be one of `backups`, `created_at`, `disk`, `id`,
  `image`, `ipv4_address`, `ipv4_address_private`, `ipv6`, `ipv6_address`, `ipv6_address_private`, `locked`,
  `memory`, `monitoring`, `name`, `price_hourly`, `price_monthly`, `private_networking`, `region`, `size`,
  `status`, `tags`, `urn`, `vcpus`, `volume_ids`, or `vpc_uuid`.

* `values` - (Required) A list of values to match against the `key` field. Only retrieves Droplets
  where the `key` field takes on one or more of the values provided here.

`sort` supports the following arguments:

* `key` - (Required) Sort the Droplets by this key. This may be one of `backups`, `created_at`, `disk`, `id`,
  `image`, `ipv4_address`, `ipv4_address_private`, `ipv6`, `ipv6_address`, `ipv6_address_private`, `locked`,
  `memory`, `monitoring`, `name`, `price_hourly`, `price_monthly`, `private_networking`, `region`, `size`,
  `status`, `urn`, `vcpus`, or `vpc_uuid`.

* `direction` - (Required) The sort direction. This may be either `asc` or `desc`.

## Attributes Reference

* `droplets` - A list of Droplets satisfying any `filter` and `sort` criteria. Each Droplet has the following attributes:  

  - `id` - The ID of the Droplet.
  - `urn` - The uniform resource name of the Droplet
  - `region` - The region the Droplet is running in.
  - `image` - The Droplet image ID or slug.
  - `size` - The unique slug that identifies the type of Droplet.
  - `disk` - The size of the Droplet's disk in GB.
  - `vcpus` - The number of the Droplet's virtual CPUs.
  - `memory` - The amount of the Droplet's memory in MB.
  - `price_hourly` - Droplet hourly price.
  - `price_monthly` - Droplet monthly price.
  - `status` - The status of the Droplet.
  - `locked` - Whether the Droplet is locked.
  - `ipv6_address` - The Droplet's public IPv6 address
  - `ipv6_address_private` - The Droplet's private IPv6 address
  - `ipv4_address` - The Droplet's public IPv4 address
  - `ipv4_address_private` - The Droplet's private IPv4 address
  - `backups` - Whether backups are enabled.
  - `ipv6` - Whether IPv6 is enabled.
  - `private_networking` - Whether private networks are enabled.
  - `monitoring` - Whether monitoring agent is installed.
  - `volume_ids` - List of the IDs of each volumes attached to the Droplet.
  - `tags` - A list of the tags associated to the Droplet.
  - `vpc_uuid` - The ID of the VPC where the Droplet is located.
