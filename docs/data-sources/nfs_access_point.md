---
page_title: "DigitalOcean: digitalocean_nfs_access_point"
subcategory: "NFS Storage"
---

# digitalocean\_nfs\_access\_point

Get information about a DigitalOcean NFS access point.

## Example Usage

Get the NFS access point by ID:

```hcl
data "digitalocean_nfs_access_point" "example" {
  id = "506f78a4-e098-11e5-ad9f-000f53306ae1"
}
```

Get the NFS access point by name and share ID:

```hcl
data "digitalocean_nfs_access_point" "example" {
  name     = "example-access-point"
  share_id = digitalocean_nfs.foobar.id
  vpc_id   = digitalocean_vpc.foobar.id
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional) The ID of the NFS access point. Conflicts with `name`, `share_id`, and `vpc_id`.
* `name` - (Optional) The name of the NFS access point. Must be used with `share_id`. Conflicts with `id`.
* `share_id` - (Optional) The ID of the NFS share. Must be used with `name`. Conflicts with `id`.
* `vpc_id` - (Optional) When looking up by `name` and `share_id`, optionally filter to the access point attached to this VPC. Conflicts with `id`.

## Attributes Reference

See the NFS Access Point Resource for details on the returned attributes — these are identical.
