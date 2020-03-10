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

One of the following arguments must be provided:

* `id` - The id of the image
* `name` - The name of the image.
* `slug` - The slug of the official image.

If `name` is specified, you may also specify:

* `source` - (Optional) Restrict the search to one of the following categories of images:
  - `all` - All images (whether public or private)
  - `applications` - One-click applications
  - `distribution` - Distributions
  - `user` - (Default) User (private) images

## Attributes Reference

The following attributes are exported:

* `slug`: Unique text identifier of the image.
* `id`: The ID of the image.
* `name`: The name of the image.
* `type`: Type of the image.
* `distribution` - The name of the distribution of the OS of the image.
* `min_disk_size`: The minimum 'disk' required for the image.
* `size_gigabytes`: The size of the image in GB.
* `private` - Is image a public image or not. Public images represent
  Linux distributions or One-Click Applications, while non-public images represent
  snapshots and backups and are only available within your account.
* `regions`: A set of the regions that the image is available in.
* `tags`: A set of tags applied to the image 
* `created`: When the image was created
* `status`: Current status of the image
* `error_message`: Any applicable error message pertaining to the image
* `image` - The id of the image (legacy parameter).

