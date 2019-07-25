---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_tag"
sidebar_current: "docs-do-resource-tag"
description: |-
  Provides a DigitalOcean Tag resource.
---

# digitalocean\_tag

Provides a DigitalOcean Tag resource. A Tag is a label that can be applied to a
Droplet resource in order to better organize or facilitate the lookups and
actions on it. Tags created with this resource can be referenced in your Droplet
configuration via their ID or name.

## Example Usage

```hcl
# Create a new tag
resource "digitalocean_tag" "foobar" {
  name = "foobar"
}

# Create a new Droplet in nyc3 with the foobar tag
resource "digitalocean_droplet" "web" {
  image  = "ubuntu-18-04-x64"
  name   = "web-1"
  region = "nyc3"
  size   = "s-1vcpu-1gb"
  tags   = [digitalocean_tag.foobar.id]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the tag

## Attributes Reference

The following attributes are exported:

* `id` - The id of the tag
* `name` - The name of the tag


## Import

Tags can be imported using the `name`, e.g.

```
terraform import digitalocean_tag.mytag tagname
```