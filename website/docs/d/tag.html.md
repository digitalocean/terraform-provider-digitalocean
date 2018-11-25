---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_tag"
sidebar_current: "docs-do-datasource-tag"
description: |-
  Get information on a tag.
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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the tag.

## Attributes Reference

The following attributes are exported:

* `id`: The ID of the tag.
