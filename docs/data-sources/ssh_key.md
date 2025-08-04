---
page_title: "DigitalOcean: digitalocean_ssh_key"
subcategory: "Account"
---

# digitalocean_ssh_key

Get information on a ssh key. This data source provides the name, public key,
and fingerprint as configured on your DigitalOcean account. This is useful if
the ssh key in question is not managed by Terraform or you need to utilize any
of the keys data.

An error is triggered if the provided ssh key name does not exist.

## Example Usage

Get the ssh key:

```hcl
data "digitalocean_ssh_key" "example" {
  name = "example"
}

resource "digitalocean_droplet" "example" {
  image    = "ubuntu-18-04-x64"
  name     = "example-1"
  region   = "nyc2"
  size     = "s-1vcpu-1gb"
  ssh_keys = [data.digitalocean_ssh_key.example.id]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the ssh key.

## Attributes Reference

The following attributes are exported:

* `id`: The ID of the ssh key.
* `public_key`: The public key of the ssh key.
* `fingerprint`: The fingerprint of the public key of the ssh key.
