---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_certificate"
sidebar_current: "docs-do-resource-certificate"
description: |-
  Provides a DigitalOcean Certificate resource.
---

# digitalocean\_certificate

Provides a DigitalOcean Certificate resource that allows you to manage
certificates for configuring TLS termination in Load Balancers.
Certificates created with this resource can be referenced in your
Load Balancer configuration via their ID. The certificate can either
be a custom one provided by you or automatically generated one with
Let's Encrypt.

## Example Usage

#### Custom Certificate

```hcl
resource "digitalocean_certificate" "cert" {
  name              = "custom-terraform-example"
  type              = "custom"
  private_key       = file("/Users/terraform/certs/privkey.pem")
  leaf_certificate  = file("/Users/terraform/certs/cert.pem")
  certificate_chain = file("/Users/terraform/certs/fullchain.pem")
}
```

#### Let's Encrypt Certificate

```hcl
resource "digitalocean_certificate" "cert" {
  name    = "le-terraform-example"
  type    = "lets_encrypt"
  domains = ["example.com"]
}
```

#### Use with Other Resources

Both custom and Let's Encrypt certificates can be used with other resources
including the `digitalocean_loadbalancer` and `digitalocean_cdn` resources.

```hcl
resource "digitalocean_certificate" "cert" {
  name    = "le-terraform-example"
  type    = "lets_encrypt"
  domains = ["example.com"]
}

# Create a new Load Balancer with TLS termination
resource "digitalocean_loadbalancer" "public" {
  name        = "secure-loadbalancer-1"
  region      = "nyc3"
  droplet_tag = "backend"

  forwarding_rule {
    entry_port     = 443
    entry_protocol = "https"

    target_port     = 80
    target_protocol = "http"

    certificate_id = digitalocean_certificate.cert.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the certificate for identification.
* `type` - (Optional) The type of certificate to provision. Can be either
`custom` or `lets_encrypt`. Defaults to `custom`.
* `private_key` - (Optional) The contents of a PEM-formatted private-key
corresponding to the SSL certificate. Only valid when type is `custom`.
* `leaf_certificate` - (Optional) The contents of a PEM-formatted public
TLS certificate. Only valid when type is `custom`.
* `certificate_chain` - (Optional) The full PEM-formatted trust chain
between the certificate authority's certificate and your domain's TLS
certificate. Only valid when type is `custom`.
* `domains` - (Optional) List of fully qualified domain names (FQDNs) for
which the certificate will be issued. The domains must be managed using
DigitalOcean's DNS. Only valid when type is `lets_encrypt`.


## Attributes Reference

The following attributes are exported:

* `id` - The unique ID of the certificate
* `name` - The name of the certificate
* `not_after` - The expiration date of the certificate
* `sha1_fingerprint` - The SHA-1 fingerprint of the certificate


## Import

Certificates can be imported using the certificate `id`, e.g.

```
terraform import digitalocean_certificate.mycertificate 892071a0-bb95-49bc-8021-3afd67a210bf
```
