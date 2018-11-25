---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_certificate"
sidebar_current: "docs-do-datasource-certificate"
description: |-
  Get information on a certificate.
---

# digitalocean_certificate

Get information on a certificate. This data source provides the name, type, state,
domains, expiry date, and the sha1 fingerprint as configured on your DigitalOcean account.
This is useful if the certificate in question is not managed by Terraform or you need to utilize
any of the certificates data.

An error is triggered if the provided certificate name does not exist.

## Example Usage

Get the certificate:

```hcl
data "digitalocean_certificate" "example" {
  name = "example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of certificate.

## Attributes Reference

The following attributes are exported:

* `id`: The ID of the certificate.
* `type`: The type of the certificate.
* `state`: the current state of the certificate.
* `domains`: Domains for which the certificate was issued.
* `not_after`: The expiration date and time of the certificate.
* `sha1_fingerprint`: The SHA1 fingerprint of the certificate.
