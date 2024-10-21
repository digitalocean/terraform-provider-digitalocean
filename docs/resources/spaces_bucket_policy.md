---
page_title: "DigitalOcean: digitalocean_spaces_bucket_policy"
subcategory: "Spaces Object Storage"
---

# digitalocean\_spaces\_bucket_policy

Provides a bucket policy resource for Spaces, DigitalOcean's object storage product.
The `digitalocean_spaces_bucket_policy` resource allows Terraform to attach bucket
policy to Spaces.

The [Spaces API](https://docs.digitalocean.com/reference/api/spaces-api/) was
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

### Limiting access to specific IP addresses

```hcl
resource "digitalocean_spaces_bucket" "foobar" {
  name   = "foobar"
  region = "nyc3"
}

resource "digitalocean_spaces_bucket_policy" "foobar" {
  region = digitalocean_spaces_bucket.foobar.region
  bucket = digitalocean_spaces_bucket.foobar.name
  policy = jsonencode({
    "Version" : "2012-10-17",
    "Statement" : [
      {
        "Sid" : "IPAllow",
        "Effect" : "Deny",
        "Principal" : "*",
        "Action" : "s3:*",
        "Resource" : [
          "arn:aws:s3:::${digitalocean_spaces_bucket.foobar.name}",
          "arn:aws:s3:::${digitalocean_spaces_bucket.foobar.name}/*"
        ],
        "Condition" : {
          "NotIpAddress" : {
            "aws:SourceIp" : "54.240.143.0/24"
          }
        }
      }
    ]
  })
}
```

!> **Warning:** Before using this policy, replace the 54.240.143.0/24 IP address range in this example with an appropriate value for your use case. Otherwise, you will lose the ability to access your bucket.

## Argument Reference

The following arguments are supported:

* `region` - (Required) The region where the bucket resides.
* `bucket` - (Required) The name of the bucket to which to apply the policy.
* `policy` - (Required) The text of the policy.

## Attributes Reference

No additional attributes are exported.

## Import

Bucket policies can be imported using the `region` and `bucket` attributes (delimited by a comma):

```
terraform import digitalocean_spaces_bucket_policy.foobar `region`,`bucket`
```
