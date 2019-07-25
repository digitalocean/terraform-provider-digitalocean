---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_volume_snapshot"
sidebar_current: "docs-do-resource-volume-snapshot"
description: |-
  Provides a DigitalOcean volume snapshot resource.
---

# digitalocean\_volume\_snapshot

Provides a DigitalOcean Volume Snapshot which can be used to create a snapshot from an existing volume.

## Example Usage

```hcl
resource "digitalocean_volume" "foobar" {
  region      = "nyc1"
  name        = "baz"
  size        = 100
  description = "an example volume"
}

resource "digitalocean_volume_snapshot" "foobar" {
  name      = "foo"
  volume_id = digitalocean_volume.foobar.id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A name for the volume snapshot.
* `volume_id` - (Required) The ID of the volume from which the volume snapshot originated.

## Attributes Reference

The following attributes are exported:

* `id` The ID of the volume snapshot.
* `created_at` - The date and time the volume snapshot was created.
* `min_disk_size` - The minimum size in gigabytes required for a volume to be created based on this volume snapshot.
* `regions` - A list of DigitalOcean region "slugs" indicating where the volume snapshot is available.
* `size` - The billable size of the volume snapshot in gigabytes.


## Import

Volume Snapshots can be imported using the `snapshot id`, e.g.

```
terraform import digitalocean_volume_snapshot.snapshot 506f78a4-e098-11e5-ad9f-000f53306ae1
```
