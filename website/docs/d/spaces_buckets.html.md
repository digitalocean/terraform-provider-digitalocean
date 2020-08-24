---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_spaces_buckets"
sidebar_current: "docs-do-datasource-spaces-buckets"
description: |-
  Retrieve information on Spaces buckets.
---

# digitalocean_spaces_buckets

Get information on Spaces buckets for use in other resources, with the ability to filter and sort the results.
If no filters are specified, all Spaces buckets will be returned.

Note: You can use the [`digitalocean_spaces_bucket`](spaces_bucket) data source to
obtain metadata about a single bucket if you already know its `name` and `region`.

## Example Usage

Use the `filter` block with a `key` string and `values` list to filter buckets.

Get all buckets in a region:

```hcl
data "digitalocean_spaces_buckets" "nyc3" {
  filter {
    key = "region"
    values = ["nyc3"]
  }
}
```
You can sort the results as well:

```hcl
data "digitalocean_spaces_buckets" "nyc3" {
  filter {
    key = "region"
    values = ["nyc3"]
  }
  sort {
    key = "name"
    direction = "desc"
  }
}
```

## Argument Reference

* `filter` - (Optional) Filter the results.
  The `filter` block is documented below.

* `sort` - (Optional) Sort the results.
  The `sort` block is documented below.

`filter` supports the following arguments:

* `key` - (Required) Filter the images by this key. This may be one of `bucket_domain_name`, `name`, `region`, or `urn`.

* `values` - (Required) A list of values to match against the `key` field. Only retrieves images
  where the `key` field takes on one or more of the values provided here.

`sort` supports the following arguments:

* `key` - (Required) Sort the images by this key. This may be one of `bucket_domain_name`, `name`, `region`, or `urn`.
* `direction` - (Required) The sort direction. This may be either `asc` or `desc`.

## Attributes Reference

* `buckets` - A list of Spaces buckets satisfying any `filter` and `sort` criteria. Each bucket has the following attributes:  

  - `name` - The name of the Spaces bucket
  - `region` - The slug of the region where the bucket is stored.
  - `urn` - The uniform resource name of the bucket
  - `bucket_domain_name` - The FQDN of the bucket (e.g. bucket-name.nyc3.digitaloceanspaces.com)
