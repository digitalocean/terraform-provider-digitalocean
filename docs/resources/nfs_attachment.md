---
page_title: "DigitalOcean: digitalocean_nfs_attachment"
subcategory: "NFS Storage"
---

# digitalocean\_nfs\_attachment

Manages attaching a NFS share to a vpc.

## Example Usage

```hcl
resource "digitalocean_vpc" "foobar" {
  name   = "example-vpc"
  region = "atl1"
}

resource "digitalocean_nfs" "foobar" {
  region = "atl1"
  name   = "example-nfs"
  size   = 50
  vpc_id = digitalocean_vpc.foobar.id
}

resource "digitalocean_nfs_attachment" "foobar" {
  share_id   = digitalocean_nfs.foobar.id
  vpc_id = digitalocean_vpc.foobar.id
}

```

## Argument Reference

The following arguments are supported:

* `share_id` - (Required) The ID of the NFS share to attach.  
* `vpc_id` - (Required) The ID of the vpc to attach the NFS share to.  

## Attributes Reference

The following attributes are exported:

* `id` - The unique identifier for the NFS attachment.  
* `share_id` - The ID of the NFS share.  
* `vpc_id` - The ID of the vpc.  

## Import

NFS attachments can be imported using the `share_id` and `vpc_id` separated by a comma, e.g.  

```
terraform import digitalocean_nfs_attachment.foobar 506f78a4-e098-11e5-ad9f-000f53306ae1,d1ebc5a4-e098-11e5-ad9f-000f53306ae1
```