---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_domain"
sidebar_current: "docs-do-resource-domain"
description: |-
  Provides a DigitalOcean domain resource.
---

# digitalocean\_domain

Provides a DigitalOcean domain resource.

## Example Usage

```hcl
# Create a new domain
resource "digitalocean_domain" "default" {
  name       = "example.com"
  ip_address = digitalocean_droplet.foo.ipv4_address
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the domain
* `ip_address` - (Optional) The IP address of the domain. If specified, this IP
   is used to created an initial A record for the domain.

## Attributes Reference

The following attributes are exported:

* `id` - The name of the domain
* `urn` - The uniform resource name of the domain

## Import

Domains can be imported using the `domain name`, e.g.

```
terraform import digitalocean_domain.mydomain mytestdomain.com
```
