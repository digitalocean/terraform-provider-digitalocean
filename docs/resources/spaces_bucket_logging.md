---
page_title: "DigitalOcean: digitalocean_spaces_bucket_logging"
subcategory: "Spaces Object Storage"
---

# digitalocean\_spaces\_bucket\_logging

Provides a bucket logging resource for Spaces, DigitalOcean's object storage product.
The `digitalocean_spaces_bucket_logging` resource allows Terraform to configure access
logging for Spaces buckets. For more information, see:
[How to Configure Spaces Access Logs](https://docs.digitalocean.com/products/spaces/how-to/access-logs/)

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


## Example Usage


```hcl
resource "digitalocean_spaces_bucket" "assets" {
  name   = "assets"
  region = "nyc3"
}

resource "digitalocean_spaces_bucket" "logs" {
  name   = "logs"
  region = "nyc3"
}

resource "digitalocean_spaces_bucket_logging" "example" {
  region = "%s"
  bucket = digitalocean_spaces_bucket.assets.id

  target_bucket = digitalocean_spaces_bucket.logs.id
  target_prefix = "access-logs/"
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Required) The region where the bucket resides.
* `bucket` - (Required) The name of the bucket which will be logged.
* `target_bucket` - (Required) The name of the bucket which will store the logs.
* `target_prefix` - (Required) The prefix for the log files.

## Attributes Reference

No additional attributes are exported.

## Import

Spaces bucket logging can be imported using the `region` and `bucket` attributes (delimited by a comma):

```
terraform import digitalocean_spaces_bucket_logging.example `region`,`bucket`
```
