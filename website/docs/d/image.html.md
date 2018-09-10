---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_image"
sidebar_current: "docs-do-datasource-image"
description: |-
  Get information on an snapshot.
---

# digitalocean_image

Get information on an snapshot images. The aim of this datasource is to enable
you to build Droplets based on snapshot names.

An error is triggered if zero or more than one result is returned by the query.

## Example Usage

Get the data about a snapshot:

```hcl
data "digitalocean_image" "example1" {
  name = "example-1.0.0"
}
```

Reuse the data about a snapshot to create a Droplet:

```hcl
data "digitalocean_image" "example1" {
  name = "example-1.0.0"
}
resource "digitalocean_droplet" "example1" {
  image  = "${data.digitalocean_image.example1.image}"
  name   = "example-1"
  region = "nyc2"
  size   = "s-1vcpu-1gb"
}
```

Get the data about an official image:

```hcl
data "digitalocean_image" "example2" {
  slug = "ubuntu-18-04-x64"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the private image.
* `slug` - (Optional) The slug of the official image.

## Attributes Reference

The following attributes are exported:

* `id`: The ID of the image.
* `image` - The id of the image.
* `distribution` - The name of the distribution of the OS of the image.
* `min_disk_size`: The minimum 'disk' required for the image.
* `private` - Is image a public image or not. Public images represents
  Linux distributions or Application, while non-public images represent
  snapshots and backups and are only available within your account.
* `regions`: The regions that the image is available in.
* `type`: Type of the image.

