---
page_title: "DigitalOcean: digitalocean_dedicated_inference_accelerators"
subcategory: "Dedicated Inference"
---

# digitalocean\_dedicated\_inference\_accelerators

Returns a list of accelerators (GPUs) attached to a dedicated inference endpoint,
with the ability to filter and sort the results.

## Example Usage

```hcl
data "digitalocean_dedicated_inference_accelerators" "example" {
  dedicated_inference_id = digitalocean_dedicated_inference.example.id
}

output "accelerators" {
  value = data.digitalocean_dedicated_inference_accelerators.example.accelerators
}
```

### Filter by slug

```hcl
data "digitalocean_dedicated_inference_accelerators" "filtered" {
  dedicated_inference_id = digitalocean_dedicated_inference.example.id

  filter {
    key    = "slug"
    values = ["gpu-h100x1-80gb"]
  }
}
```

## Argument Reference

* `dedicated_inference_id` - (Required) The ID of the dedicated inference endpoint to list accelerators for.
* `filter` - (Optional) Filter the results. The `filter` block is documented below.
* `sort` - (Optional) Sort the results. The `sort` block is documented below.

---

`filter` supports the following arguments:

* `key` - (Required) Filter the accelerators by this key. This may be one of `id`, `name`, `slug`, `status`, `created_at`.
* `values` - (Required) A list of values to match against the `key` field.
* `match_by` - (Optional) One of `exact` (default), `re`, or `substring`.
* `all` - (Optional) Set to `true` to require that a field match all of the `values` instead of just one.

`sort` supports the following arguments:

* `key` - (Required) Sort the accelerators by this key. This may be one of the keys listed in `filter`.
* `direction` - (Optional) The sort direction. This may be either `asc` or `desc`.

## Attributes Reference

* `accelerators` - A list of accelerators satisfying any `filter` and `sort` criteria. Each element contains:
  - `id` - The unique ID of the accelerator.
  - `name` - The name of the accelerator.
  - `slug` - The slug identifier for the accelerator type.
  - `status` - The current status of the accelerator.
  - `created_at` - The date and time when the accelerator was created.
