---
page_title: "DigitalOcean: digitalocean_nfs_attachment"
subcategory: "NFS Storage"
---

# digitalocean\_nfs\_attachment

Manages attaching an NFS share to a VPC. A share can be attached to multiple VPCs by creating one `digitalocean_nfs_attachment` resource per VPC.

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
  performance_tier = "high"
}

resource "digitalocean_nfs_attachment" "foobar" {
  share_id = digitalocean_nfs.foobar.id
  vpc_id   = digitalocean_vpc.foobar.id
  region   = "atl1"
}
```

## Example Usage - Multiple VPCs

Attach the same NFS share to additional VPCs one at a time:

```hcl
resource "digitalocean_vpc" "primary" {
  name   = "primary-vpc"
  region = "atl1"
}

resource "digitalocean_vpc" "secondary" {
  name   = "secondary-vpc"
  region = "atl1"
}

resource "digitalocean_nfs" "example" {
  region           = "atl1"
  name             = "example-nfs"
  size             = 50
  vpc_id           = digitalocean_vpc.primary.id
  performance_tier = "high"
}

resource "digitalocean_nfs_attachment" "primary" {
  share_id = digitalocean_nfs.example.id
  vpc_id   = digitalocean_vpc.primary.id
  region   = "atl1"
}

resource "digitalocean_nfs_attachment" "secondary" {
  share_id = digitalocean_nfs.example.id
  vpc_id   = digitalocean_vpc.secondary.id
  region   = "atl1"

  depends_on = [digitalocean_nfs_attachment.primary]
}
```

Deleting an attachment resource detaches the share from that VPC only. Other VPC attachments remain in place.

## Argument Reference

The following arguments are supported:

* `share_id` - (Required) The ID of the NFS share to attach.  
* `vpc_id` - (Required) The ID of the VPC to attach the NFS share to.  
* `region` - (Required) The region of the NFS share.

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
