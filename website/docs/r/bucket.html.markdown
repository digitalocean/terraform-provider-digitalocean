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
The authentication requirement can be met by installing `awscli` and adding a
DigitalOcean profile to your credentials file (usually found in`~/.aws/credentials`).
This should look like:

```
[digitalocean-spaces]
aws_access_key_id = QAZWSXRFVTGBYHNUJMIK
aws_secret_access_key = 1QAZ2WSX3EDC4RFV5TGB6YHN7UJM8IK9OL0P1QAZ2WS
```

For more information, See [An Introduction to DigitalOcean Spaces](https://www.digitalocean.com/community/tutorials/an-introduction-to-digitalocean-spaces)

## Example Usage

```hcl
# Create a new bucket
resource "digitalocean_bucket" "foobar" {
  name = "foobar"
  region = "nyc3"
  profile = "digitalocean-spaces"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the bucket
* `region` - (Required) The region where the bucket resides
* `profile` - (Required) Spaces Access Profile (defined in your AWS Credentials file)
* `acl` - Canned ACL applied on bucket creation (`private` or `public-read`)

## Attributes Reference

The following attributes are exported:

* `name` - The name of the bucket
* `region` - The name of the region


## Import

Buckets can be imported using the `name`, e.g.

```
terraform import digitalocean_bucket.foobar `name`
```
