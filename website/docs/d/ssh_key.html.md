---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_ssh_key"
sidebar_current: "docs-do-datasource-ssh-key"
description: |-
  Get information on a ssh key.
---

# digitalocean_ssh_key

Get information on a ssh key. This data source provides the name, public key, and fingerprint
as configured on your Digital Ocean account. This is useful if the ssh key 
in question is not managed by Terraform or you need to utilize any of the keys data.

An error is triggered if the provided ssh key name does not exist.

## Example Usage

Get the ssh key:

```hcl
data "digitalocean_ssh_key" "example" {
  name = "example"
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
