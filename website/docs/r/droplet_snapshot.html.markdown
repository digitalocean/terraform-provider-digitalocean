---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_droplet_snapshot"
sidebar_current: "docs-do-resource-droplet-snapshot"
description: |-
  Provides a DigitalOcean droplet snapshot resource.
---

# digitalocean\_droplet\_snapshot

Provides a DigitalOcean Droplet Snapshot which can be used to create a snapshot from an existing droplet.

## Example Usage

```hcl
resource "digitalocean_droplet" "foobar" {
	name      = "foo-%d"
	size      = "512mb"
	image     = "centos-7-x64"
	region    = "nyc3"
	user_data = "foobar"
}

resource "digitalocean_droplet_snapshot" "foobar" {
	droplet_id = "${digitalocean_droplet.foobar.id}"
	name = "snapshotone-foo"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A name for the droplet snapshot.
* `droplet_id` - (Required) The ID of the droplet from which the snapshot originated.

## Attributes Reference

The following attributes are exported:

* `id` The ID of the droplet snapshot.
* `created_at` - The date and time the droplet snapshot was created.
* `min_disk_size` - The minimum size in gigabytes required for a droplet to be created based on this droplet snapshot.
* `regions` - A list of DigitalOcean region "slugs" indicating where the droplet snapshot is available.
* `size` - The billable size of the droplet snapshot in gigabytes.


## Import

Droplet Snapshots can be imported using the `snapshot id`, e.g.

```
terraform import digitalocean_droplet_snapshot.snapshot 506f78a4-e098-11e5-ad9f-000f53306ae1
```
