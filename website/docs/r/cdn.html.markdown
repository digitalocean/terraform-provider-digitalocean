---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_cdn"
sidebar_current: "docs-do-resource-cdn"
description: |-
  Provides a DigitalOcean CDN Endpoint resource.
---

# digitalocean\_cdn

Provides a DigitalOcean CDN Endpoint resource.

## Example Usage

```hcl
# Create a new Spaces Bucket
resource "digitalocean_spaces_bucket" "bucket" {
	name = "example.sfo2.digitaloceanspaces.com"
	region = "sfo2"
	acl = "public-read"
}

# Add a CDN endpoint to the Spaces Bucket
resource "digitalocean_cdn" "cdn" {
  origin = "${digitalocean_spaces_bucket.bucket.bucket_domain_name}"
}

# Output the endpoint for the CDN reosurce
output "fqdn" {
  value = "${digitalocean_cdn.cdn.endpoint}"
}
```

## Argument Reference

The following arguments are supported:

* 

## Attributes Reference

The following attributes are exported:

* 

## Import


