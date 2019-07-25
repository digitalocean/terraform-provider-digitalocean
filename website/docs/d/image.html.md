---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_image"
sidebar_current: "docs-do-datasource-image"
description: |-
  Get information on an snapshot.
---

# digitalocean_image

Get information on an images for use in other resources (e.g. creating a Droplet
based on snapshot). This data source provides all of the image properties as
configured on your DigitalOcean account. This is useful if the image in question
is not managed by Terraform or you need to utilize any of the image's data.

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
data "digitalocean_image" "example" {
  name = "example-1.0.0"
}

resource "digitalocean_droplet" "example" {
  image  = data.digitalocean_image.example.id
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
* `private` - Is image a public image or not. Public images represent
  Linux distributions or One-Click Applications, while non-public images represent
  snapshots and backups and are only available within your account.
* `regions`: The regions that the image is available in.
* `type`: Type of the image.

