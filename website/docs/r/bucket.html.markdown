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
for S3 usage with Spaces without much tweaking.  Due to this, Spaces mirrors S3's
authentication framework and requests to Spaces require a key pair similar to
Amazon's Access ID and Secret Key.

For more information, See [An Introduction to DigitalOcean Spaces](https://www.digitalocean.com/community/tutorials/an-introduction-to-digitalocean-spaces)

## Example Usage

```hcl
# Create a new bucket
resource "digitalocean_bucket" "foobar" {
  name = "foobar"
  region = "nyc3"
  access_key = "${var.do_access_key}"
  secret_key = "${var.do_secret_key}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the bucket
* `region` - (Required) The region where the bucket resides
* `access_key` - (Required) Spaces Access ID
* `secret_key` - (Required) Spaces Secret Key

## Attributes Reference

The following attributes are exported:

* `name` - The name of the bucket
* `region` - The name of the region


## Import

Buckets can be imported using the `name`, e.g.

```
terraform import digitalocean_bucket.foobar `name`
```
