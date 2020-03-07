---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_images"
sidebar_current: "docs-do-datasource-images"
description: |-
  Retrieve metadata about DigitalOcean images (public and private).
---

# digitalocean_images

Get information on images for use in other resources (e.g. creating a Droplet
based on snapshot), with the ability to filter and sort the results. If no filters are specified,
all images will be returned.

This data source provides all of the image properties as configured on your DigitalOcean account.
This is useful if the image in question is not managed by Terraform or you need to utilize any
of the image's data.

Note: You can use the `digitalocean_image` data source to obtain metadata about a single
image if you already know the `slug`, unique `name`, or `id` to retrieve.

## Example Usage

Use the `filter` block with a `key` string and `values` list to filter images.

For example to find all Ubuntu images:

```hcl
data "digitalocean_images" "ubuntu" {
  filter {
    key = "distribution"
    values = ["Ubuntu"]
  }
} 
```

You can filter on multiple fields and sort the results as well:

```hcl
data "digitalocean_images" "available" {
  filter {
    key = "distribution"
    values = ["Ubuntu"]
  }
  filter {
    key = "regions"
    values = ["nyc3"]
  }
  sort {
    key = "created"
    direction = "desc"
  }
}
```

## Argument Reference

* `filter` - (Optional) Filter the results.
  The `filter` block is documented below.
* `sort` - (Optional) Sort the results.
  The `sort` block is documented below.

`filter` supports the following:
* `key` - (Required) Filter the images by this key. This may be one of `distribution`, `error_message`,
  `id`, `image`, `min_disk_size`, `name`, `private`, `regions`, `size_gigabytes`, `slug`, `status`,
  `tags`, or `type`.

* `values` - (Required) A list of values to match against the `key` field. Only retrieves images
  where the `key` field takes on one or more of the values provided here.

`sort` supports the following:

* `key` - (Required) Sort the images by this key. This may be one of `distribution`, `error_message`, `id`,
   `image`, `min_disk_size`, `name`, `private`, `size_gigabytes`, `slug`, `status`, or `type`.
* `direction` - (Required) The sort direction. This may be either `asc` or `desc`.

## Attributes Reference

* `images` - A set of images satisfying any `filter` and `sort` criteria. Each image has the following attributes:  
  - `slug`: Unique text identifier of the image.
  - `id`: The ID of the image.
  - `name`: The name of the image.
  - `type`: Type of the image.
  - `distribution` - The name of the distribution of the OS of the image.
  - `min_disk_size`: The minimum 'disk' required for the image.
  - `size_gigabytes`: The size of the image in GB.
  - `private` - Is image a public image or not. Public images represent
    Linux distributions or One-Click Applications, while non-public images represent
    snapshots and backups and are only available within your account.
  - `regions`: A set of the regions that the image is available in.
  - `tags`: A set of tags applied to the image 
  - `created`: When the image was created
  - `status`: Current status of the image
  - `error_message`: Any applicable error message pertaining to the image
  - `image` - The id of the image (legacy parameter).

