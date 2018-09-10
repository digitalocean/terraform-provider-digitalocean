---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_floating_ip"
sidebar_current: "docs-do-datasource-floating-ip"
description: |-
  Get information on a floating ip.
---

# digitalocean_floating_ip

Get information on a floating ip. This data source provides the region and droplet id
as configured on your DigitalOcean account. This is useful if the floating ip
in question is not managed by Terraform or you need to find the droplet the ip is
attached to.

An error is triggered if the provided floating ip does not exist.

## Example Usage

Get the floating ip:

```hcl
variable "public_ip" {}

data "digitalocean_floating_ip" "example" {
  ip_address = "${var.public_ip}"
}
```

## Argument Reference

The following arguments are supported:

* `ip_address` - (Required) The allocated ip address of the specific floating ip to retrieve.

## Attributes Reference

The following attributes are exported:

* `region`: The region that the floating ip is reserved to.
* `droplet_id`: The droplet id that the floating ip has been assigned to.
