---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_record"
sidebar_current: "docs-do-datasource-record"
description: |-
  Get information on a DNS record.
---

# digitalocean_record

Get information on a DNS record. This data source provides the name, TTL, and zone
file as configured on your Digital Ocean account. This is useful if the record
in question is not managed by Terraform.

An error is triggered if the provided domain name or record are not managed with
your Digital Ocean account.

## Example Usage

Get data from a DNS record:

```hcl
data "digitalocean_record" "example" {
  domain  = "example.com"
  name    = "test"
}

output "record_type" {
  value = "${data.digitalocean_record.example.type}"
}
output "record_ttl" {
  value = "${data.digitalocean_record.example.ttl}"
}
```

```
  $ terraform apply

data.digitalocean_record.example: Refreshing state...

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

record_ttl = 3600
record_type = A
```

## Argument Reference

The following arguments are supported:

* `name` - The name of the record.
* `domain` - The domain name of the record.

## Attributes Reference

The following attributes are exported:

* `name`: See Argument Reference above.
* `type`:	The type of the DNS record. For example: A, CNAME, TXT, ...
* `id`:	The host name, alias, or service being defined by the record.
* `data`:	Variable data depending on record type. For example, the "data" value for an A record would be the IPv4 address to which the domain will be mapped. For a CAA record, it would contain the domain name of the CA being granted permission to issue certificates.
* `priority`:	The priority for SRV and MX records.
* `port`:	The port for SRV records.
* `ttl`: This value is the time to live for the record, in seconds. This defines the time frame that clients can cache queried information before a refresh should be requested.
* `weight`:	The weight for SRV records.
* `flags`: An unsigned integer between 0-255 used for CAA records.
* `tag`: The parameter tag for CAA records. Valid values are "issue", "wildissue", or "iodef"
