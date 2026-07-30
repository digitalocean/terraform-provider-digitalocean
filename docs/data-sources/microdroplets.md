---
page_title: "DigitalOcean: digitalocean_microdroplets"
subcategory: "MicroDroplets"
---

# digitalocean_microdroplets

Fetches a list of [DigitalOcean MicroDroplets](https://docs.digitalocean.com/products/microdroplets/) with optional `region` and `name` server-side filters, plus `filter` and `sort` blocks for client-side narrowing.

## Example Usage

```hcl
data "digitalocean_microdroplets" "all" {}

data "digitalocean_microdroplets" "nyc3" {
  region = "nyc3"
}

data "digitalocean_microdroplets" "by_name" {
  name = "my-microdroplet"
}

data "digitalocean_microdroplets" "running" {
  filter {
    key    = "current_state"
    values = ["running"]
  }
}
```

## Argument Reference

* `region` - (Optional) Server-side filter: only include MicroDroplets in this region. Conflicts with `name`.
* `name` - (Optional) Server-side filter: only include MicroDroplets whose name matches exactly. Conflicts with `region`.
* `filter` - (Optional) Repeatable client-side filter block:
  * `key` - (Required) Field to match.
  * `values` - (Required) List of values to match on `key`.
  * `all` - (Optional) Require every value to match. Defaults to `false`.
  * `match_by` - (Optional) `exact`, `re`, or `substring`. Defaults to `exact`.
* `sort` - (Optional) Repeatable client-side sort block:
  * `key` - (Required) Field to sort by.
  * `direction` - (Optional) `asc` (default) or `desc`.

## Attributes Reference

* `micro_droplets` - A list of MicroDroplet records. Each record has the same attribute set as the `digitalocean_microdroplet` data source.
