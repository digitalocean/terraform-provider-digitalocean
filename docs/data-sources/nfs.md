---
page_title: "DigitalOcean: digitalocean_nfs"
subcategory: "NFS Storage"
---

# digitalocean\_nfs

Get information about a DigitalOcean NFS share.

## Example Usage

Get the NFS share by name and region:

```hcl
data "digitalocean_nfs" "example" {
  name   = "example-nfs"
  region = "nyc1"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the NFS share.
* `region` - (Optional) The region where the NFS share is located.

## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier for the NFS share.
* `name` - Name of the NFS share.
* `region` - The region where the NFS share is located.
* `size` - The size of the NFS share in GiB.
* `status` - The current status of the NFS share.
* `host` - The host IP of the NFS server accessible from the associated VPC.
* `mount_path` - The mount path for accessing the NFS share.  
