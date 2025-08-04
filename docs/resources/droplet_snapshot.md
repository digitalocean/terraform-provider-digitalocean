---
page_title: "DigitalOcean: digitalocean_droplet_snapshot"
subcategory: "Backups & Snapshots"
---

# digitalocean\_droplet\_snapshot

Provides a resource which can be used to create a snapshot from an existing DigitalOcean Droplet.

## Example Usage

```hcl
resource "digitalocean_droplet" "web" {
  name   = "web-01"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
}

resource "digitalocean_droplet_snapshot" "web-snapshot" {
  droplet_id = digitalocean_droplet.web.id
  name       = "web-snapshot-01"
}


resource "digitalocean_droplet" "from-snapshot" {
  image  = digitalocean_droplet_snapshot.web-snapshot.id
  name   = "web-02"
  region = "nyc3"
  size   = "s-2vcpu-4gb"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A name for the Droplet snapshot.
* `droplet_id` - (Required) The ID of the Droplet from which the snapshot will be taken.

## Attributes Reference

The following attributes are exported:

* `id` The ID of the Droplet snapshot.
* `created_at` - The date and time the Droplet snapshot was created.
* `min_disk_size` - The minimum size in gigabytes required for a Droplet to be created based on this snapshot.
* `regions` - A list of DigitalOcean region "slugs" indicating where the droplet snapshot is available.
* `size` - The billable size of the Droplet snapshot in gigabytes.


## Import

Droplet Snapshots can be imported using the `snapshot id`, e.g.

```
terraform import digitalocean_droplet_snapshot.mysnapshot 123456
```
