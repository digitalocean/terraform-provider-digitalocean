---
page_title: "DigitalOcean: digitalocean_custom_image"
---

# digitalocean\_custom\_image

Provides a resource which can be used to create a custom Image from a URL

## Example Usage

```hcl
resource "digitalocean_custom_image" "flatcar" {
  name   = "flatcar"
  url = "https://stable.release.flatcar-linux.net/amd64-usr/2605.7.0/flatcar_production_digitalocean_image.bin.bz2"
  regions = ["nyc3"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A name for the Custom Image.
* `url` - (Required) A URL from which the custom Linux virtual machine image may be retrieved.
* `regions` - (Required) A list of regions. (Currently only one is supported)

## Attributes Reference

The following attributes are exported:

* `image_id` A unique number that can be used to identify and reference a specific image.
* `type` Describes the kind of image.
* `slug` A uniquely identifying string for each image.
* `public` Indicates whether the image in question is public or not.
* `min_disk_size` The minimum disk size in GB required for a Droplet to use this image.
* `size_gigabytes` The size of the image in gigabytes.
* `created_at` A time value given in ISO8601 combined date and time format that represents when the image was created.
* `status` A status string indicating the state of a custom image.
