---
page_title: "DigitalOcean: digitalocean_microvms"
subcategory: "MicroVMs"
---

# digitalocean_microvms

Get information on MicroVMs for use in other resources, with the ability to
filter and sort the results. If no filters are specified, all MicroVMs are
returned.

The `region`, `name`, and `tag_name` arguments are applied by the API and
combine with logical AND. The `filter` and `sort` blocks are applied to the
results afterwards.

Note: You can use the [`digitalocean_microvm`](microvm) data source to obtain
metadata about a single MicroVM if you already know its `id` or unique `name`.

-> **Note:** MicroVMs are in public preview.

## Example Usage

Get all MicroVMs in a region that carry a tag:

```hcl
data "digitalocean_microvms" "web" {
  region   = "nyc3"
  tag_name = "web"
}

output "names" {
  value = data.digitalocean_microvms.web.micro_vms[*].name
}
```

Get all running MicroVMs, sorted by creation time:

```hcl
data "digitalocean_microvms" "running" {
  filter {
    key    = "current_state"
    values = ["running"]
  }

  sort {
    key       = "created_at"
    direction = "desc"
  }
}
```

## Argument Reference

* `region` - (Optional) Only return MicroVMs in this region.
* `name` - (Optional) Only return MicroVMs with this exact name.
* `tag_name` - (Optional) Only return MicroVMs that carry this tag.
* `filter` - (Optional) Filter the results.
  The `filter` block is documented below.
* `sort` - (Optional) Sort the results.
  The `sort` block is documented below.

`filter` supports the following arguments:

* `key` - (Required) Filter the MicroVMs by this key. This may be one of `id`,
  `name`, `region`, `networking`, `http_protocol`, `current_state`,
  `failure_reason`, `auto_resume`, `ports`, `tags`, `urn`, or `created_at`.
* `values` - (Required) A list of values to match against the `key` field. Only
  retrieves MicroVMs where the `key` field takes on one or more of the values
  provided here.
* `match_by` - (Optional) One of `exact` (default), `re`, or `substring`. For
  string-typed fields, specify `re` to match by using the `values` as regular
  expressions, or specify `substring` to match by treating the `values` as
  substrings to find within the string field.
* `all` - (Optional) Set to `true` to require that a field match all of the
  `values` instead of just one or more of them. This is useful when matching
  against multi-valued fields such as lists or sets where you want to ensure
  that all of the `values` are present in the list or set.

`sort` supports the following arguments:

* `key` - (Required) Sort the MicroVMs by this key. This may be one of `id`,
  `name`, `region`, `networking`, `http_protocol`, `current_state`,
  `auto_resume`, `urn`, or `created_at`.
* `direction` - (Required) The sort direction. This may be either `asc` or `desc`.

## Attributes Reference

* `micro_vms` - A list of MicroVMs satisfying any `filter` and `sort` criteria.
  Each MicroVM has the same attributes as the
  [`digitalocean_microvm`](microvm) data source.
