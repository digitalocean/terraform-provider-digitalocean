---
page_title: "DigitalOcean: digitalocean_spaces_bucket_cors_configuration"
subcategory: "Spaces Object Storage"
---

# digitalocean\_spaces\_cors_configuration

Provides a CORS configuration resource for Spaces, DigitalOcean's object storage product.
The `digitalocean_spaces_bucket_cors_configuration` resource allows Terraform to to attach CORS configuration to Spaces.

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

```hcl
provider "digitalocean" {
  token = var.digitalocean_token

  spaces_access_id  = var.access_id
  spaces_secret_key = var.secret_key
}

resource "digitalocean_spaces_bucket" "static-assets" {
  # ...
}
```

For more information, See [An Introduction to DigitalOcean Spaces](https://www.digitalocean.com/community/tutorials/an-introduction-to-digitalocean-spaces)

## Example Usage

### Create a Key in a Spaces Bucket

```hcl
resource "digitalocean_spaces_bucket" "foobar" {
  name   = "foobar"
  region = "nyc3"
}

resource "digitalocean_spaces_bucket_cors_configuration" "test" {
  bucket = digitalocean_spaces_bucket.foobar.id
  region = "nyc3"

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST"]
    allowed_origins = ["https://s3-website-test.hashicorp.com"]
    expose_headers  = ["ETag"]
    max_age_seconds = 3000
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The name of the bucket to which to apply the CORS configuration.
* `region` - (Required) The region where the bucket resides.
* `cors_rule` - (Required) Set of origins and methods (cross-origin access that you want to allow). See below. You can configure up to 100 rules.

`cors_rule` supports the following:

* `allowed_headers` - (Optional) Set of Headers that are specified in the Access-Control-Request-Headers header.
* `allowed_methods` - (Required) Set of HTTP methods that you allow the origin to execute. Valid values are GET, PUT, HEAD, POST, and DELETE.
* `allowed_origins` - (Required) Set of origins you want customers to be able to access the bucket from.
* `expose_headers` - (Optional) Set of headers in the response that you want customers to be able to access from their applications (for example, from a JavaScript XMLHttpRequest object).
* `id` - (Optional) Unique identifier for the rule. The value cannot be longer than 255 characters.
* `max_age_seconds` - (Optional) Time in seconds that your browser is to cache the preflight response for the specified resource.

## Attributes Reference

No additional attributes are exported.

## Import

Bucket policies can be imported using the `region` and `bucket` attributes (delimited by a comma):

```
terraform import digitalocean_spaces_bucket_cors_configuration.foobar `region`,`bucket`
```
