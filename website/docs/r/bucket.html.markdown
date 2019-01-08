---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_bucket"
sidebar_current: "docs-do-resource-bucket"
description: |-
  Provides a DigitalOcean Spaces Bucket resource.
---

# digitalocean\_bucket

Provides a bucket resource for Spaces, DigitalOcean's object storage product.

The [Spaces API](https://developers.digitalocean.com/documentation/spaces/) was
designed to be interoperable with Amazon's AWS S3 API. This allows users to
interact with the service while using the tools they already know. Spaces
mirrors S3's authentication framework and requests to Spaces require a key pair
similar to Amazon's Access ID and Secret Key.

The authentication requirement can be met by either setting the
`SPACES_ACCESS_KEY_ID` and `SPACES_SECRET_ACCESS_KEY` environment variables or
the provider's `access_id` and `secret_key` arguments to the access ID and
secret you generate via the DigitalOcean control panel. For example:

```
provider "digitalocean" {
  token      = "${var.digitalocean_token}"

  access_id  = "${var.access_id}"
  secret_key = "${var.secret_key}"
}

resource "digitalocean_bucket" "static-assets" {
  # ...
}
```

For more information, See [An Introduction to DigitalOcean Spaces](https://www.digitalocean.com/community/tutorials/an-introduction-to-digitalocean-spaces)

## Example Usage

```hcl
# Create a new bucket
resource "digitalocean_bucket" "foobar" {
  name   = "foobar"
  region = "nyc3"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the bucket
* `region` - The region where the bucket resides (Defaults to `nyc3`)
* `acl` - Canned ACL applied on bucket creation (`private` or `public-read`)

## Attributes Reference

The following attributes are exported:

* `name` - The name of the bucket
* `region` - The name of the region


## Import

Buckets can be imported using the `region` and `name` attributes (delimited by a comma):

```
terraform import digitalocean_bucket.foobar `region`,`name`
```
