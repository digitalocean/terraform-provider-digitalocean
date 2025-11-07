---
page_title: "DigitalOcean: digitalocean_nfs"
subcategory: "NFS Storage"
---

# digitalocean\_nfs

Get information about a DigitalOcean NFS share.

## Example Usage

Get the NFS share by ID and region:

```hcl
data "digitalocean_nfs" "example" {
  id = "506f78a4-e098-11e5-ad9f-000f53306ae1"
  region = "nyc1"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The id of the NFS share.
* `region` - (Required) The region where the NFS share is located.

## Attributes Reference

See the NFS Resource for details on the returned attributes — these are identical.  
