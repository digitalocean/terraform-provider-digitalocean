---
page_title: "DigitalOcean: digitalocean_microdroplet_image"
subcategory: "MicroDroplets"
---

# digitalocean_microdroplet_image

Fetches a single [DigitalOcean MicroDroplet image](https://docs.digitalocean.com/products/microdroplets/) by `id` or by `name`.

Name lookups page through the full image list on the account since the API has no filtered list endpoint for images.

## Example Usage

```hcl
data "digitalocean_microdroplet_image" "byid" {
  id = "4a2c9f10-3d18-4b6a-9e3d-2b7f8e0f1c11"
}

data "digitalocean_microdroplet_image" "byname" {
  name = "my-app-v1"
}
```

## Argument Reference

Exactly one of the following must be provided:

* `id` - (Optional) The MicroDroplet image UUID.
* `name` - (Optional) The MicroDroplet image name. Errors if zero or more than one match.

## Attributes Reference

* `source` - The OCI reference used when the image was imported.
* `status` - Lifecycle status.
* `created_at` - RFC3339 timestamp of when the image was created.
* `urn` - The uniform resource name (URN) for the MicroDroplet image.
