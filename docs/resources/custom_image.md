---
page_title: "DigitalOcean: digitalocean_custom_image"
subcategory: "Account"
---

# digitalocean\_custom\_image

Provides a resource which can be used to create a [custom image](https://www.digitalocean.com/docs/images/custom-images/)
from a URL. The URL must point to an image in one of the following file formats:

- Raw (.img) with an MBR or GPT partition table
- qcow2
- VHDX
- VDI
- VMDK

The image may be compressed using gzip or bzip2. See the DigitalOcean Custom
Image documentation for [additional requirements](https://www.digitalocean.com/docs/images/custom-images/#image-requirements).

## Example Usage

```hcl
resource "digitalocean_custom_image" "flatcar" {
  name    = "flatcar"
  url     = "https://stable.release.flatcar-linux.net/amd64-usr/2605.7.0/flatcar_production_digitalocean_image.bin.bz2"
  regions = ["nyc3"]
}

resource "digitalocean_droplet" "example" {
  image    = digitalocean_custom_image.flatcar.id
  name     = "example-01"
  region   = "nyc3"
  size     = "s-1vcpu-1gb"
  ssh_keys = [12345]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A name for the Custom Image.
* `url` - (Required) A URL from which the custom Linux virtual machine image may be retrieved.
* `regions` - (Required) A list of regions. (Currently only one is supported).
* `description` - An optional description for the image.
* `distribution` - An optional distribution name for the image. Valid values are documented [here](https://docs.digitalocean.com/reference/api/api-reference/#operation/create_custom_image)
* `tags` - A list of optional tags for the image.

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
