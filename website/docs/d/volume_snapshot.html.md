---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_volume_snapshot"
sidebar_current: "docs-do-datasource-volume-snapshot"
description: |-
  Get information about a DigitalOcean volume snapshot.
---

# digitalocean\_volume\_snapshot

Volume snapshots are saved instances of a block storage volume. Use this data
source to retrieve the ID of a DigitalOcean volume snapshot for use in other
resources.

## Example Usage

```
data "digitalocean_volume_snapshot" "snapshot" {
  name_regex  = "^web"
  region= "nyc2"
  most_recent = true
}
```

## Argument Reference

* `name` - (Optional) The name of the volume snapshot.

* `name_regex` - (Optional) A regex string to apply to the volume snapshot list returned by DigitalOcean. This allows more advanced filtering not supported from the DigitalOcean API. This filtering is done locally on what DigitalOcean returns.

* `region` - (Optional) A "slug" representing a DigitalOcean region (e.g. `nyc1`). If set, only volume snapshots available in the region will be returned.

* `most_recent` - (Optional) If more than one result is returned, use the most recent volume snapshot.

~> **NOTE:** If more or less than a single match is returned by the search,
Terraform will fail. Ensure that your search is specific enough to return
a single volume snapshot ID only, or use `most_recent` to choose the most recent one.

## Attributes Reference

The following attributes are exported:

* `id` The ID of the volume snapshot.
* `created_at` - The date and time the volume snapshot was created.
* `min_disk_size` - The minimum size in gigabytes required for a volume to be created based on this volume snapshot.
* `regions` - A list of DigitalOcean region "slugs" indicating where the volume snapshot is available.
* `volume_id` - The ID of the volume from which the volume snapshot originated.
* `size` - The billable size of the volume snapshot in gigabytes.