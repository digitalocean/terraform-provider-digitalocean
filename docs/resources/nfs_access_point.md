---
page_title: "DigitalOcean: digitalocean_nfs_access_point"
subcategory: "NFS Storage"
---

# digitalocean\_nfs\_access\_point

Provides a DigitalOcean NFS access point for a Network File Storage share. Access points define export paths and access policies for mounting an NFS share from a VPC.

## Example Usage

```hcl
resource "digitalocean_vpc" "foobar" {
  name   = "example-vpc"
  region = "atl1"
}

resource "digitalocean_nfs" "foobar" {
  region           = "atl1"
  name             = "example-nfs"
  size             = 50
  vpc_id           = digitalocean_vpc.foobar.id
  performance_tier = "standard"
}

resource "digitalocean_nfs_access_point" "foobar" {
  name     = "example-access-point"
  share_id = digitalocean_nfs.foobar.id
  path     = "/data"
  vpc_id   = digitalocean_vpc.foobar.id

  access_policy {
    anonuid                      = 65534
    anongid                      = 65534
    protocols                    = ["NFS4"]
    squash_config                = "ROOT_SQUASH"
    identity_enforcement_enabled = false
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A name for the NFS access point.
* `share_id` - (Required) The ID of the NFS share.
* `path` - (Required) The export path for the access point.
* `vpc_id` - (Required) The ID of the VPC that can access this access point. The VPC must be attached to the NFS share.
* `access_policy` - (Required) Access policy configuration for the access point. See [Access Policy](#access-policy) below.

### Access Policy

The `access_policy` block supports the following:

* `anonuid` - (Optional) Anonymous UID mapped for NFS clients. Defaults to `65534`.
* `anongid` - (Optional) Anonymous GID mapped for NFS clients. Defaults to `65534`.
* `protocols` - (Optional) List of NFS protocols. Defaults to `["NFS4"]`.
* `squash_config` - (Optional) Squash configuration. Valid values are `NO_SQUASH`, `ROOT_SQUASH`, and `ALL_SQUASH`. Defaults to `ROOT_SQUASH`.
* `identity_enforcement_enabled` - (Optional) Whether identity enforcement is enabled. Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier for the NFS access point.
* `name` - Name of the NFS access point.
* `share_id` - The ID of the NFS share.
* `path` - The export path for the access point.
* `vpc_id` - The ID of the VPC associated with the access point.
* `status` - The status of the access point.
* `is_default` - Whether this is the default access point for the share.
* `created_at` - The date and time when the access point was created.
* `updated_at` - The date and time when the access point was last updated.
* `access_policy` - Access policy configuration. See [Access Policy](#access-policy) above.

## Import

NFS access points can be imported using the access point ID, e.g.

```
terraform import digitalocean_nfs_access_point.foobar 506f78a4-e098-11e5-ad9f-000f53306ae1
```
