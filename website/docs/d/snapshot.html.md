---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_snapshot"
sidebar_current: "docs-do-datasource-snapshot"
description: |-
  Get information about a DigitalOcean snapshot.
---

# digitalocean\_snapshot

Snapshots are saved instances of either a Droplet or a block storage volume. Use this data source to retrieve the ID of a DigitalOcean snapshot for use in other resources.

## Example Usage

```
data "digitalocean_snapshot" "snapshot" {
  most_recent = true
  name_regex  = "^web"
  region_filter = "nyc2"
  resource_type = "droplet"
}
```

## Argument Reference

* `resource_type` - (Required) The type of DigitalOcean resource from which the snapshot originated. This currently must be either `droplet` or `volume`.

* `most_recent` - (Optional) If more than one result is returned, use the most
recent snapshot.

* `name_regex` - (Optional) A regex string to apply to the snapshot list returned by DigitalOcean. This allows more advanced filtering not supported from the DigitalOcean API. This filtering is done locally on what DigitalOcean returns.

* `region_filter` - (Optional) A "slug" representing a DigitalOcean region (e.g. `nyc1`). If set, only snapshots available in the region will be returned.

~> **NOTE:** If more or less than a single match is returned by the search,
Terraform will fail. Ensure that your search is specific enough to return
a single snapshot ID only, or use `most_recent` to choose the most recent one.

## Attributes Reference

`id` is set to the ID of the found snapshot. In addition, the following attributes are exported:

* `created_at` - The date and time the image was created.
* `min_disk_size` - The minimum size in gigabytes required for a volume or Droplet to be created based on this snapshot.
* `name` - The name of the snapshot.
* `regions` - A list of DigitalOcean region "slugs" indicating where the snapshot is available.
* `resource_id` - The ID of the resource from which the snapshot originated.
* `size_gigabytes` - The billable size of the snapshot in gigabytes.
