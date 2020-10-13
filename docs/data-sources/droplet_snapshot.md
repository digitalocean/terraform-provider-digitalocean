---
page_title: "DigitalOcean: digitalocean_droplet_snapshot"
---

# digitalocean\_droplet\_snapshot

Droplet snapshots are saved instances of a Droplet. Use this data
source to retrieve the ID of a DigitalOcean Droplet snapshot for use in other
resources.

## Example Usage

Get the Droplet snapshot:

```hcl
data "digitalocean_droplet_snapshot" "web-snapshot" {
  name_regex  = "^web"
  region      = "nyc3"
  most_recent = true
}
```

## Argument Reference

* `name` - (Optional) The name of the Droplet snapshot.

* `name_regex` - (Optional) A regex string to apply to the Droplet snapshot list returned by DigitalOcean. This allows more advanced filtering not supported from the DigitalOcean API. This filtering is done locally on what DigitalOcean returns.

* `region` - (Optional) A "slug" representing a DigitalOcean region (e.g. `nyc1`). If set, only Droplet snapshots available in the region will be returned.

* `most_recent` - (Optional) If more than one result is returned, use the most recent Droplet snapshot.

~> **NOTE:** If more or less than a single match is returned by the search,
Terraform will fail. Ensure that your search is specific enough to return
a single Droplet snapshot ID only, or use `most_recent` to choose the most recent one.

## Attributes Reference

The following attributes are exported:

* `id` The ID of the Droplet snapshot.
* `created_at` - The date and time the Droplet snapshot was created.
* `min_disk_size` - The minimum size in gigabytes required for a Droplet to be created based on this Droplet snapshot.
* `regions` - A list of DigitalOcean region "slugs" indicating where the Droplet snapshot is available.
* `droplet_id` - The ID of the Droplet from which the Droplet snapshot originated.
* `size` - The billable size of the Droplet snapshot in gigabytes.