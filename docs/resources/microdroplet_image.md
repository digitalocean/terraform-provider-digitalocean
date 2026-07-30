---
page_title: "DigitalOcean: digitalocean_microdroplet_image"
subcategory: "MicroDroplets"
---

# digitalocean_microdroplet_image

Provides a [DigitalOcean MicroDroplet image](https://docs.digitalocean.com/products/microdroplets/) resource. Images are OCI references imported into DigitalOcean and consumed by `digitalocean_microdroplet` resources.

Both `name` and `source` are immutable — changing either recreates the image. There is no in-place update; images are create-only and delete-only.

## Example Usage

```hcl
resource "digitalocean_microdroplet_image" "example" {
  name   = "my-app-v1"
  source = "docker.io/myorg/my-app:v1"
}

resource "digitalocean_microdroplet" "example" {
  name   = "my-app"
  region = "nyc3"
  size   = "microdroplet-1"
  image  = digitalocean_microdroplet_image.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the MicroDroplet image. Forces recreation on change.
* `source` - (Required) OCI reference to import (e.g. `docker.io/library/nginx:latest` or a DOCR ref). Forces recreation on change.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The MicroDroplet image UUID.
* `urn` - The uniform resource name (URN) for the MicroDroplet image.
* `status` - Lifecycle status (`IMAGE_AVAILABLE` once import completes).
* `created_at` - RFC3339 timestamp of when the image was created.

## Import

A MicroDroplet image can be imported using its `id`, e.g.

```
terraform import digitalocean_microdroplet_image.example 4a2c9f10-3d18-4b6a-9e3d-2b7f8e0f1c11
```
