---
page_title: "DigitalOcean: digitalocean_tag"
---

# digitalocean_tag

Get information on a tag. This data source provides the name as configured on
your DigitalOcean account. This is useful if the tag name in question is not
managed by Terraform or you need validate if the tag exists in the account.

An error is triggered if the provided tag name does not exist.

## Example Usage

Get the tag:

```hcl
data "digitalocean_tag" "example" {
  name = "example"
}

resource "digitalocean_droplet" "example" {
  image  = "ubuntu-18-04-x64"
  name   = "example-1"
  region = "nyc2"
  size   = "s-1vcpu-1gb"
  tags   = [data.digitalocean_tag.example.name]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the tag.

## Attributes Reference

The following attributes are exported:

* `id`: The ID of the tag.
* `total_resource_count` - A count of the total number of resources that the tag is applied to.
* `droplets_count` - A count of the Droplets the tag is applied to.
* `images_count` - A count of the images that the tag is applied to.
* `volumes_count` - A count of the volumes that the tag is applied to.
* `volume_snapshots_count` - A count of the volume snapshots that the tag is applied to.
* `databases_count` - A count of the database clusters that the tag is applied to.
