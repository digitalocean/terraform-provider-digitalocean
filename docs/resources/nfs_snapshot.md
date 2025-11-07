---
page_title: "DigitalOcean: digitalocean_nfs_snapshot"
subcategory: "NFS Storage"
---

# digitalocean\_nfs\_snapshot

Provides a DigitalOcean NFS snapshot which can be used to create new NFS shares.

## Example Usage

```hcl
resource "digitalocean_vpc" "foobar" {
  name   = "example-vpc"
  region = "nyc1"
}

resource "digitalocean_nfs" "foobar" {
  region = "nyc1"
  name   = "example-nfs"
  size   = 50
  vpc_id = digitalocean_vpc.foobar.id
}

resource "digitalocean_nfs_snapshot" "foobar" {
  name     = "example-snapshot"
  share_id = digitalocean_nfs.foobar.id
  region   = "nyc1"
}

```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A name for the NFS snapshot. Must be lowercase and composed only of numbers, letters, and "-", up to a limit of 64 characters.
* `share_id` - (Required) The ID of the NFS share to snapshot.
* `region` - (Required) The region where the NFS snapshot will be created.

## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier for the NFS snapshot.
* `name` - Name of the NFS snapshot.
* `share_id` - The ID of the NFS share.
* `region` - The region where the NFS snapshot is stored.
* `size` - The size of the snapshot in GiB.
* `created_at` - The date and time when the snapshot was created.

## Import

NFS snapshots can be imported using the snapshot ID, e.g.  

```
terraform import digitalocean_nfs_snapshot.foobar 506f78a4-e098-11e5-ad9f-000f53306ae1
```