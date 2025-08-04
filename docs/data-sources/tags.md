---
page_title: "DigitalOcean: digitalocean_tags"
subcategory: "Account"
---

# digitalocean_tags

Returns a list of tags in your DigitalOcean account, with the ability to
filter and sort the results. If no filters are specified, all tags will be
returned.

## Example Usage

```hcl
data "digitalocean_tags" "list" {
  sort {
    key       = "total_resource_count"
    direction = "asc"
  }
}

output "sorted_tags" {
  value = data.digitalocean_tags.list.tags
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) Filter the results.
  The `filter` block is documented below.
* `sort` - (Optional) Sort the results.
  The `sort` block is documented below.

`filter` supports the following arguments:

* `key` - (Required) Filter the tags by this key. This may be one of `name`, `total_resource_count`,  `droplets_count`, `images_count`, `volumes_count`, `volume_snapshots_count`, or `databases_count`.
* `values` - (Required) Only retrieves tags which keys has value that matches
  one of the values provided here.

* `match_by` - (Optional) One of `exact` (default), `re`, or `substring`. For string-typed fields, specify `re` to
  match by using the `values` as regular expressions, or specify `substring` to match by treating the `values` as
  substrings to find within the string field.
  
* `all` - (Optional) Set to `true` to require that a field match all of the `values` instead of just one or more of
  them. This is useful when matching against multi-valued fields such as lists or sets where you want to ensure
  that all of the `values` are present in the list or set.

`sort` supports the following arguments:

* `key` - (Required) Sort the tags by this key. This may be one of `name`, `total_resource_count`,  `droplets_count`, `images_count`, `volumes_count`, `volume_snapshots_count`, or `databases_count`.
* `direction` - (Required) The sort direction. This may be either `asc` or `desc`.

## Attributes Reference

The following attributes are exported for each tag:

* `name` - The name of the tag.
* `total_resource_count` - A count of the total number of resources that the tag is applied to.
* `droplets_count` - A count of the Droplets the tag is applied to.
* `images_count` - A count of the images that the tag is applied to.
* `volumes_count` - A count of the volumes that the tag is applied to.
* `volume_snapshots_count` - A count of the volume snapshots that the tag is applied to.
* `databases_count` - A count of the database clusters that the tag is applied to.
