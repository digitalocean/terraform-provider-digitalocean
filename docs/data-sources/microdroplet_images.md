---
page_title: "DigitalOcean: digitalocean_microdroplet_images"
subcategory: "MicroDroplets"
---

# digitalocean_microdroplet_images

Fetches a list of [DigitalOcean MicroDroplet images](https://docs.digitalocean.com/products/microdroplets/) on the account. No server-side filters — narrow the results using `filter` and `sort`.

## Example Usage

```hcl
data "digitalocean_microdroplet_images" "all" {}

data "digitalocean_microdroplet_images" "available" {
  filter {
    key    = "status"
    values = ["IMAGE_AVAILABLE"]
  }
}

data "digitalocean_microdroplet_images" "newest_first" {
  sort {
    key       = "created_at"
    direction = "desc"
  }
}
```

## Argument Reference

* `filter` - (Optional) Repeatable client-side filter block:
  * `key` - (Required) Field to match.
  * `values` - (Required) List of values to match on `key`.
  * `all` - (Optional) Require every value to match. Defaults to `false`.
  * `match_by` - (Optional) `exact`, `re`, or `substring`. Defaults to `exact`.
* `sort` - (Optional) Repeatable client-side sort block:
  * `key` - (Required) Field to sort by.
  * `direction` - (Optional) `asc` (default) or `desc`.

## Attributes Reference

* `micro_droplet_images` - A list of MicroDroplet image records. Each record has the same attribute set as the `digitalocean_microdroplet_image` data source.
