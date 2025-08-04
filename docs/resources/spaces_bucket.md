---
page_title: "DigitalOcean: digitalocean_spaces_bucket"
subcategory: "Spaces Object Storage"
---

# digitalocean\_spaces\_bucket

Provides a bucket resource for Spaces, DigitalOcean's object storage product.

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

### Create a New Bucket

```hcl
resource "digitalocean_spaces_bucket" "foobar" {
  name   = "foobar"
  region = "nyc3"
}
```

### Create a New Bucket With CORS Rules

```hcl
resource "digitalocean_spaces_bucket" "foobar" {
  name   = "foobar"
  region = "nyc3"

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["GET"]
    allowed_origins = ["*"]
    max_age_seconds = 3000
  }

  cors_rule {
    allowed_headers = ["*"]
    allowed_methods = ["PUT", "POST", "DELETE"]
    allowed_origins = ["https://www.example.com"]
    max_age_seconds = 3000
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the bucket
* `region` - The region where the bucket resides (Defaults to `nyc3`)
* `acl` - Canned ACL applied on bucket creation: `private` or `public-read` (Defaults to `private`)
* `cors_rule` - (Optional) A rule of Cross-Origin Resource Sharing (documented below).
* `lifecycle_rule` - (Optional) A configuration of object lifecycle management (documented below).
* `versioning` - (Optional) A state of versioning (documented below)
* `force_destroy` - Unless `true`, the bucket will only be destroyed if empty (Defaults to `false`)

The `cors_rule` object supports the following:

* `allowed_headers` - (Optional) A list of headers that will be included in the CORS preflight request's `Access-Control-Request-Headers`. A header may contain one wildcard (e.g. `x-amz-*`).
* `allowed_methods` - (Required) A list of HTTP methods (e.g. `GET`) which are allowed from the specified origin.
* `allowed_origins` - (Required) A list of hosts from which requests using the specified methods are allowed. A host may contain one wildcard (e.g. http://*.example.com).
* `max_age_seconds` - (Optional) The time in seconds that browser can cache the response for a preflight request.

The `lifecycle_rule` object supports the following:

* `id` - (Optional) Unique identifier for the rule.
* `prefix` - (Optional) Object key prefix identifying one or more objects to which the rule applies.
* `enabled` - (Required) Specifies lifecycle rule status.
* `abort_incomplete_multipart_upload_days` (Optional) Specifies the number of days after initiating a multipart
   upload when the multipart upload must be completed or else Spaces will abort the upload.
* `expiration` - (Optional) Specifies a time period after which applicable objects expire (documented below).
* `noncurrent_version_expiration` - (Optional) Specifies when non-current object versions expire (documented below).

At least one of `expiration` or `noncurrent_version_expiration` must be specified.

The `expiration` object supports the following:

* `date` - (Optional) Specifies the date/time after which you want applicable objects to expire. The argument uses
  RFC3339 format, e.g. "2020-03-22T15:03:55Z" or parts thereof e.g. "2019-02-28".
* `days` - (Optional) Specifies the number of days after object creation when the applicable objects will expire.
* `expired_object_delete_marker` - (Optional) On a versioned bucket (versioning-enabled or versioning-suspended
  bucket), setting this to true directs Spaces to delete expired object delete markers.

The `noncurrent_version_expiration` object supports the following:

* `days` - (Required) Specifies the number of days after which an object's non-current versions expire.

The `versioning` object supports the following:

* `enabled` - (Optional) Enable versioning. Once you version-enable a bucket, it can never return to an unversioned
state. You can, however, suspend versioning on that bucket.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the bucket
* `urn` - The uniform resource name for the bucket
* `region` - The name of the region
* `bucket_domain_name` - The FQDN of the bucket (e.g. bucket-name.nyc3.digitaloceanspaces.com)
* `endpoint` - The FQDN of the bucket without the bucket name (e.g. nyc3.digitaloceanspaces.com)

## Import

Buckets can be imported using the `region` and `name` attributes (delimited by a comma):

```
terraform import digitalocean_spaces_bucket.foobar `region`,`name`
```
