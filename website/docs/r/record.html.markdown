---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_record"
sidebar_current: "docs-do-resource-record"
description: |-
  Provides a DigitalOcean DNS record resource.
---

# digitalocean\_record

Provides a DigitalOcean DNS record resource.

## Example Usage

```hcl
resource "digitalocean_domain" "default" {
  name = "example.com"
}

# Add a record to the domain
resource "digitalocean_record" "www" {
  domain = digitalocean_domain.default.name
  type   = "A"
  name   = "www"
  value  = "192.168.0.11"
}

# Output the FQDN for the record
output "fqdn" {
  value = digitalocean_record.www.fqdn
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required) The type of record. Must be one of `A`, `AAAA`, `CAA`, `CNAME`, `MX`, `NS`, `TXT`, or `SRV`.
* `domain` - (Required) The domain to add the record to.
* `value` - (Required) The value of the record.
* `name` - (Required) The name of the record.
* `port` - (Optional) The port of the record. Only valid when type is `SRV`.  Must be between 1 and 65535.
* `priority` - (Optional) The priority of the record. Only valid when type is `MX` or `SRV`. Must be between 0 and 65535.
* `weight` - (Optional) The weight of the record. Only valid when type is `SRV`.  Must be between 0 and 65535.
* `ttl` - (Optional) The time to live for the record, in seconds. Must be at least 0.
* `flags` - (Optional) The flags of the record. Only valid when type is `CAA`. Must be between 0 and 255.
* `tag` - (Optional) The tag of the record. Only valid when type is `CAA`. Must be one of `issue`, `issuewild`, or `iodef`.

## Attributes Reference

The following attributes are exported:

* `id` - The record ID
* `fqdn` - The FQDN of the record

## Import

Records can be imported using the domain name and record `id` when joined with a comma. See the following example:

```
terraform import digitalocean_record.example_record example.com,12345678
```
