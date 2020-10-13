---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_volume_attachment"
sidebar_current: "docs-do-resource-volume-attachment"
description: |-
  Provides a DigitalOcean volume attachment resource.
---

# digitalocean\_volume\_attachment

Manages attaching a Volume to a Droplet.

~> **NOTE:** Volumes can be attached either directly on the `digitalocean_droplet` resource, or using the `digitalocean_volume_attachment` resource - but the two cannot be used together. If both are used against the same Droplet, the volume attachments will constantly drift.


## Example Usage

```hcl
resource "digitalocean_volume" "foobar" {
  region                  = "nyc1"
  name                    = "baz"
  size                    = 100
  initial_filesystem_type = "ext4"
  description             = "an example volume"
}

resource "digitalocean_droplet" "foobar" {
  name   = "baz"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-18-04-x64"
  region = "nyc1"
}

resource "digitalocean_volume_attachment" "foobar" {
  droplet_id = digitalocean_droplet.foobar.id
  volume_id  = digitalocean_volume.foobar.id
}
```

## Argument Reference

The following arguments are supported:

* `droplet_id` - (Required) ID of the Droplet to attach the volume to.
* `volume_id` - (Required) ID of the Volume to be attached to the Droplet.

## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier for the volume attachment.

