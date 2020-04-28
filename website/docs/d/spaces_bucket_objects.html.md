---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_spaces_bucket_objects"
sidebar_current: "docs-do-datasource-spaces-bucket-objects"
description: |-
  Returns keys and metadata of Spaces objects
---

# digitalocean_spaces_bucket_objects

~> **NOTE on `max_keys`:** Retrieving very large numbers of keys can adversely affect Terraform's performance.

The bucket-objects data source returns keys (i.e., file names) and other metadata about objects in a Spaces bucket.

## Example Usage

The following example retrieves a list of all object keys in a Spaces bucket and creates corresponding Terraform object
data sources:

```hcl
data "digitalocean_spaces_bucket_objects" "my_objects" {
  bucket = "ourcorp"
  region = "nyc3"
}

data "digitalocean_spaces_bucket_object" "object_info" {
  count  = length(data.digitalocean_spaces_bucket_objects.my_objects.keys)
  key    = element(data.digitalocean_spaces_bucket_objects.my_objects.keys, count.index)
  bucket = data.digitalocean_spaces_bucket_objects.my_objects.bucket
  region = data.digitalocean_spaces_bucket_objects.my_objects.region
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) Lists object keys in this Spaces bucket
* `region` - (Required) The slug of the region where the bucket is stored.
* `prefix` - (Optional) Limits results to object keys with this prefix (Default: none)
* `delimiter` - (Optional) A character used to group keys (Default: none)
* `encoding_type` - (Optional) Encodes keys using this method (Default: none; besides none, only "url" can be used)
* `max_keys` - (Optional) Maximum object keys to return (Default: 1000)

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `keys` - List of strings representing object keys
* `common_prefixes` - List of any keys between `prefix` and the next occurrence of `delimiter` (i.e., similar to subdirectories of the `prefix` "directory"); the list is only returned when you specify `delimiter`
* `owners` - List of strings representing object owner IDs (see `fetch_owner` above)
