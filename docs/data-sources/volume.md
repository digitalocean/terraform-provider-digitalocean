---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_volume"
sidebar_current: "docs-do-datasource-volume"
description: |-
  Get information on a volume.
---

# digitalocean_volume

Get information on a volume for use in other resources. This data source provides
all of the volumes properties as configured on your DigitalOcean account. This is
useful if the volume in question is not managed by Terraform or you need to utilize
any of the volumes data.

An error is triggered if the provided volume name does not exist.

## Example Usage

Get the volume:

```hcl
data "digitalocean_volume" "example" {
  name   = "app-data"
  region = "nyc3"
}
```

Reuse the data about a volume to attach it to a Droplet:

```hcl
data "digitalocean_volume" "example" {
  name   = "app-data"
  region = "nyc3"
}

resource "digitalocean_droplet" "example" {
  name       = "foo"
  size       = "s-1vcpu-1gb"
  image      = "ubuntu-18-04-x64"
  region     = "nyc3"
}

resource "digitalocean_volume_attachment" "foobar" {
  droplet_id = digitalocean_droplet.example.id
  volume_id  = data.digitalocean_volume.example.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of block storage volume.
* `region` - (Optional) The region the block storage volume is provisioned in.

## Attributes Reference

The following attributes are exported:

* `id`: The ID of the block storage volume.
* `urn`: The uniform resource name for the storage volume.
* `size` - The size of the block storage volume in GiB.
* `description` - Text describing a block storage volume.
* `filesystem_type` - Filesystem type currently in-use on the block storage volume.
* `filesystem_label` - Filesystem label currently in-use on the block storage volume.
* `droplet_ids` - A list of associated Droplet ids.
* `tags` - A list of the tags associated to the Volume.