---
page_title: "DigitalOcean: digitalocean_nfs_snapshot"
subcategory: "NFS Storage"
---

# digitalocean\_nfs\_snapshot

Get information about a DigitalOcean NFS snapshot.

## Example Usage

Get the NFS snapshot by ID:

```hcl
data "digitalocean_nfs_snapshot" "example" {
  id = "506f78a4-e098-11e5-ad9f-000f53306ae1"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the NFS snapshot.
* `region` - (Required) The region where the NFS snapshot is located.

## Attributes Reference

See the NFS Snapshot Resource for details on the returned attributes — these are identical.  
