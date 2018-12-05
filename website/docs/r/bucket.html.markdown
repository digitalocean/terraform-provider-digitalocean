---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_bucket"
sidebar_current: "docs-do-resource-bucket"
description: |-
  Provides a DigitalOcean Spaces Bucket resource.
---

# digitalocean\_bucket

Provides a DigitalOcean Spaces bucket resource. Spaces is DigitalOcean's object
storage product.  It has been designed to operate almost exactly like Amazon's
S3 service (even using their terminology). This allows users to reuse code written
for S3 usage with Spaces without much tweaking. Spaces mirrors S3's authentication
framework and requests to Spaces require a key pair similar to Amazon's Access ID
and Secret Key.  Due to these similarities, this functionality uses the AWS Go SDK to make these calls.
The authentication requirement can be met by either setting the `DO_ACCESS_KEY_ID` and `DO_SECRET_ACCESS_KEY`
environment variables or the `access_id` and `secret_key` arguments to the access ID and secret you generate via the Digital Ocean control panel.

For more information, See [An Introduction to DigitalOcean Spaces](https://www.digitalocean.com/community/tutorials/an-introduction-to-digitalocean-spaces)

## Example Usage

```hcl
# Create a new bucket
resource "digitalocean_bucket" "foobar" {
  name = "foobar"
  region = "nyc3"
}
```

## Argument Reference

The following arguments are supported:

* `access_id` - (Required) The access key ID used for Spaces API operations (Defaults to the value of the `DO_ACCESS_KEY_ID` environment variable)
* `secret_key` - (Required) The secret access key used for Spaces API operations (Defaults to the value of the `DO_SECRET_ACCESS_KEY` environment variable)
* `name` - (Required) The name of the bucket
* `region` - (Optional) The region where the bucket resides (Defaults to `nyc3`)
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
