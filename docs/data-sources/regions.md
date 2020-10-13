---
page_title: "DigitalOcean: digitalocean_regions"
---

# digitalocean_regions

Retrieve information about all supported DigitalOcean regions, with the ability to
filter and sort the results. If no filters are specified, all regions will be returned.

Note: You can use the [`digitalocean_region`](region) data source
to obtain metadata about a single region if you already know the `slug` to retrieve.

## Example Usage

Use the `filter` block with a `key` string and `values` list to filter regions.

For example to find all available regions:

```hcl
data "digitalocean_regions" "available" {
  filter {
    key    = "available"
    values = ["true"]
  }
}
```

You can filter on multiple fields and sort the results as well:

```hcl
data "digitalocean_regions" "available" {
  filter {
    key    = "available"
    values = ["true"]
  }
  filter {
    key    = "features"
    values = ["private_networking"]
  }
  sort {
    key       = "name"
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

* `key` - (Required) Filter the regions by this key. This may be one of `slug`,
  `name`, `available`, `features`, or `sizes`.

* `values` - (Required) A list of values to match against the `key` field. Only retrieves regions
  where the `key` field takes on one or more of the values provided here.

* `match_by` - (Optional) One of `exact` (default), `re`, or `substring`. For string-typed fields, specify `re` to
  match by using the `values` as regular expressions, or specify `substring` to match by treating the `values` as
  substrings to find within the string field.
  
* `all` - (Optional) Set to `true` to require that a field match all of the `values` instead of just one or more of
  them. This is useful when matching against multi-valued fields such as lists or sets where you want to ensure
  that all of the `values` are present in the list or set.

`sort` supports the following arguments:

* `key` - (Required) Sort the regions by this key. This may be one of `slug`,
  `name`, or `available`.
* `direction` - (Required) The sort direction. This may be either `asc` or `desc`.

## Attributes Reference

* `regions` - A set of regions satisfying any `filter` and `sort` criteria. Each region has the following attributes:  
  - `slug` - A human-readable string that is used as a unique identifier for each region.
  - `name` - The display name of the region.
  - `available` - A boolean value that represents whether new Droplets can be created in this region.
  - `sizes` - A set of identifying slugs for the Droplet sizes available in this region.
  - `features` - A set of features available in this region.
