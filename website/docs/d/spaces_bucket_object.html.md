---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_spaces_bucket_object"
sidebar_current: "docs-do-datasource-spaces-bucket-object"
description: |-
  Provides metadata and optionally content of a Spaces object
---

# digitalocean_spaces_bucket_object

The Spaces object data source allows access to the metadata and
_optionally_ (see below) content of an object stored inside a Spaces bucket.

~> **Note:** The content of an object (`body` field) is available only for objects which have a human-readable
`Content-Type` (`text/*` and `application/json`). This is to prevent printing unsafe characters and potentially
downloading large amount of data which would be thrown away in favor of metadata.

## Example Usage

The following example retrieves a text object (which must have a `Content-Type`
value starting with `text/`) and uses it as the `user_data` for a Droplet:

```hcl
data "digitalocean_spaces_bucket_object" "bootstrap_script" {
  bucket = "ourcorp-deploy-config"
  region = "nyc3"
  key    = "droplet-bootstrap-script.sh"
}

resource "digitalocean_droplet" "web" {
  image     = "ubuntu-18-04-x64"
  name      = "web-1"
  region    = "nyc2"
  size      = "s-1vcpu-1gb"
  user_data = data.digitalocean_spaces_bucket_object.bootstrap_script.body
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The name of the bucket to read the object from.
* `region` - (Required) The slug of the region where the bucket is stored.
* `key` - (Required) The full path to the object inside the bucket
* `version_id` - (Optional) Specific version ID of the object returned (defaults to latest version)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `body` - Object data (see **limitations above** to understand cases in which this field is actually available)
* `cache_control` - Specifies caching behavior along the request/reply chain.
* `content_disposition` - Specifies presentational information for the object.
* `content_encoding` - Specifies what content encodings have been applied to the object and thus what decoding mechanisms must be applied to obtain the media-type referenced by the Content-Type header field.
* `content_language` - The language the content is in.
* `content_length` - Size of the body in bytes.
* `content_type` - A standard MIME type describing the format of the object data.
* `etag` - [ETag](https://en.wikipedia.org/wiki/HTTP_ETag) generated for the object (an MD5 sum of the object content in case it's not encrypted)
* `expiration` - If the object expiration is configured (see [object lifecycle management](http://docs.aws.amazon.com/AmazonS3/latest/dev/object-lifecycle-mgmt.html)), the field includes this header. It includes the expiry-date and rule-id key value pairs providing object expiration information. The value of the rule-id is URL encoded.
* `expires` - The date and time at which the object is no longer cacheable.
* `last_modified` - Last modified date of the object in RFC1123 format (e.g. `Mon, 02 Jan 2006 15:04:05 MST`)
* `metadata` - A map of metadata stored with the object in Spaces
* `version_id` - The latest version ID of the object returned.
* `website_redirect_location` - If the bucket is configured as a website, redirects requests for this object to another object in the same bucket or to an external URL. Spaces stores the value of this header in the object metadata.

-> **Note:** Terraform ignores all leading `/`s in the object's `key` and treats multiple `/`s in the rest of the
object's `key` as a single `/`, so values of `/index.html` and `index.html` correspond to the same Spaces object
as do `first//second///third//` and `first/second/third/`.
