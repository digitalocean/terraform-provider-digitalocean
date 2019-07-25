---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_spaces_bucket"
sidebar_current: "docs-do-resource-spaces-bucket"
description: |-
  Provides a DigitalOcean Spaces Bucket resource.
---

# digitalocean\_spaces\_bucket

Provides a bucket resource for Spaces, DigitalOcean's object storage product.

The [Spaces API](https://developers.digitalocean.com/documentation/spaces/) was
designed to be interoperable with Amazon's AWS S3 API. This allows users to
interact with the service while using the tools they already know. Spaces
mirrors S3's authentication framework and requests to Spaces require a key pair
similar to Amazon's Access ID and Secret Key.

The authentication requirement can be met by either setting the
`SPACES_ACCESS_KEY_ID` and `SPACES_SECRET_ACCESS_KEY` environment variables or
the provider's `spaces_access_id` and `spaces_secret_key` arguments to the
access ID and secret you generate via the DigitalOcean control panel. For
example:

```
provider "digitalocean" {
  token             = var.digitalocean_token

  spaces_access_id  = var.access_id
  spaces_secret_key = var.secret_key
}

resource "digitalocean_spaces_bucket" "static-assets" {
  # ...
}
```

For more information, See [An Introduction to DigitalOcean Spaces](https://www.digitalocean.com/community/tutorials/an-introduction-to-digitalocean-spaces)

## Example Usage

```hcl
# Create a new bucket
resource "digitalocean_spaces_bucket" "foobar" {
  name   = "foobar"
  region = "nyc3"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the bucket
* `region` - The region where the bucket resides (Defaults to `nyc3`)
* `acl` - Canned ACL applied on bucket creation (`private` or `public-read`)
* `force_destroy` - Unless `true`, the bucket will only be destroyed if empty (Defaults to `false`)

## Attributes Reference

The following attributes are exported:

* `name` - The name of the bucket
* `urn` - The uniform resource name for the bucket
* `region` - The name of the region
* `bucket_domain_name` - The FQDN of the bucket (e.g. bucket-name.nyc3.digitaloceanspaces.com)

## Import

Buckets can be imported using the `region` and `name` attributes (delimited by a comma):

```
terraform import digitalocean_spaces_bucket.foobar `region`,`name`
```
