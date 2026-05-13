---
page_title: "DigitalOcean: digitalocean_nfs"
subcategory: "NFS Storage"
---

# digitalocean\_nfs

Provides a DigitalOcean NFS share which can be mounted to Droplets to provide shared storage.

## Example Usage

```hcl
resource "digitalocean_vpc" "example" {
  name   = "example-vpc"
  region = "nyc1"
}

resource "digitalocean_nfs" "example" {
  region           = "nyc1"
  name             = "example-nfs"
  size             = 50
  vpc_id           = digitalocean_vpc.example.id
  performance_tier = "high"
}
```

## Example Usage - Moving Share Between VPCs

To move an NFS share from one VPC to another using the Reassign API:

```hcl
resource "digitalocean_vpc" "source" {
  name   = "source-vpc"
  region = "nyc1"
}

resource "digitalocean_vpc" "destination" {
  name   = "destination-vpc"
  region = "nyc1"
}

resource "digitalocean_nfs" "example" {
  region           = "nyc1"
  name             = "example-nfs"
  size             = 50
  vpc_id           = digitalocean_vpc.source.id
  performance_tier = "high"
}

# Attach to source VPC
resource "digitalocean_nfs_attachment" "source" {
  share_id = digitalocean_nfs.example.id
  vpc_id   = digitalocean_vpc.source.id
  region   = "nyc1"
}

# Reassign to destination VPC - uses the efficient Reassign API
resource "digitalocean_nfs_attachment" "destination" {
  share_id = digitalocean_nfs.example.id
  vpc_id   = digitalocean_vpc.destination.id
  region   = "nyc1"

  depends_on = [digitalocean_nfs_attachment.source]
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required) The region where the NFS share will be created.
* `name` - (Required) A name for the NFS share. Must be lowercase and composed only of numbers, letters, and "-", up to a limit of 64 characters. The name must begin with a letter.
* `size` - (Required) The size of the NFS share in GiB. Minimum size is 50 GiB.
* `vpc_id` - (Required) The ID of the VPC where the NFS share will be created.
* `performance_tier` - (Optional) The performance tier for the NFS share. Can be `standard` or `high`. Defaults to `high`. Changing this will cause the performance tier to be switched.
> **Note:** You cannot downgrade the performance tier from `high` to `standard` after creation. Upgrades from `standard` to `high` are allowed.

## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier for the NFS share.
* `name` - Name of the NFS share.
* `region` - The region where the NFS share is created.
* `size` - The size of the NFS share in GiB.
* `performance_tier` - The performance tier of the NFS share (`standard` or `high`).
* `vpc_id` - The ID of the VPC where the NFS share is located.
* `status` - The current status of the NFS share.
* `created_at` - The date and time when the NFS share was created.
* `host` - The host IP of the NFS server accessible from the associated VPC.
* `mount_path` - The mount path for accessing the NFS share.

## Notes

Multiple NFS shares can now be attached to the same VPC, providing greater flexibility for storage management.

## Import

NFS shares can be imported using the `share id` and the `region`, e.g.

```
terraform import digitalocean_nfs.foobar 506f78a4-e098-11e5-ad9f-000f53306ae1,atl1
```
