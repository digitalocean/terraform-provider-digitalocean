---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_cdn"
sidebar_current: "docs-do-resource-cdn"
description: |-
  Provides a DigitalOcean CDN Endpoint resource.
---

# digitalocean\_cdn

Provides a DigitalOcean CDN Endpoint resource for use with Spaces.

## Example Usage

#### Basic Example

```hcl
# Create a new Spaces Bucket
resource "digitalocean_spaces_bucket" "mybucket" {
  name   = "example"
  region = "sfo2"
  acl    = "public-read"
}

# Add a CDN endpoint to the Spaces Bucket
resource "digitalocean_cdn" "mycdn" {
  origin = digitalocean_spaces_bucket.mybucket.bucket_domain_name
}

# Output the endpoint for the CDN resource
output "fqdn" {
  value = digitalocean_cdn.mycdn.endpoint
}
```

#### Custom Sub-Domain Example

```hcl
# Create a new Spaces Bucket
resource "digitalocean_spaces_bucket" "mybucket" {
  name   = "example"
  region = "sfo2"
  acl    = "public-read"
}

# Create a DigitalOcean managed Let's Encrypt Certificate
resource "digitalocean_certificate" "cert" {
  name    = "cdn-cert"
  type    = "lets_encrypt"
  domains = ["static.example.com"]
}

# Add a CDN endpoint with a custom sub-domain to the Spaces Bucket
resource "digitalocean_cdn" "mycdn" {
  origin         = digitalocean_spaces_bucket.mybucket.bucket_domain_name
  custom_domain  = "static.example.com"
  certificate_id = digitalocean_certificate.cert.id
}
```

## Argument Reference

The following arguments are supported:

* `origin` - (Required) The fully qualified domain name, (FQDN) for a Space.
* `ttl` - (Optional) The time to live for the CDN Endpoint, in seconds. Default is 3600 seconds.
* `certificate_id`- (Optional) The ID of a DigitalOcean managed TLS certificate used for SSL when a custom subdomain is provided.
* `custom_domain` - (Optional) The fully qualified domain name (FQDN) of the custom subdomain used with the CDN Endpoint.

## Attributes Reference

The following attributes are exported:

* `id` - A unique ID that can be used to identify and reference a CDN Endpoint.
* `origin` - The fully qualified domain name, (FQDN) of a space referenced by the CDN Endpoint.
* `endpoint` - The fully qualified domain name (FQDN) from which the CDN-backed content is served.
* `created_at` - The date and time when the CDN Endpoint was created.
* `ttl` - The time to live for the CDN Endpoint, in seconds.
* `certificate_id`- The ID of a DigitalOcean managed TLS certificate used for SSL when a custom subdomain is provided.
* `custom_domain` - The fully qualified domain name (FQDN) of the custom subdomain used with the CDN Endpoint.


## Import

CDN Endpoints can be imported using the CDN `id`, e.g.

```
terraform import digitalocean_cdn.mycdn fb06ad00-351f-45c8-b5eb-13523c438661
```
