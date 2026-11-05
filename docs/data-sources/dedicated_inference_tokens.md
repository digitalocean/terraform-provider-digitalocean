---
page_title: "DigitalOcean: digitalocean_dedicated_inference_tokens"
subcategory: "Dedicated Inference"
---

# digitalocean\_dedicated\_inference\_tokens

Returns a list of API tokens for a dedicated inference endpoint, with the ability
to filter and sort the results.

~> **Note:** Token values (secrets) are not returned by this data source. Only
token metadata (ID, name, creation time) is available.

## Example Usage

```hcl
data "digitalocean_dedicated_inference_tokens" "example" {
  dedicated_inference_id = digitalocean_dedicated_inference.example.id
}

output "tokens" {
  value = data.digitalocean_dedicated_inference_tokens.example.tokens
}
```

### Filter by name

```hcl
data "digitalocean_dedicated_inference_tokens" "filtered" {
  dedicated_inference_id = digitalocean_dedicated_inference.example.id

  filter {
    key    = "name"
    values = ["my-token"]
  }
}
```

## Argument Reference

* `dedicated_inference_id` - (Required) The ID of the dedicated inference endpoint to list tokens for.
* `filter` - (Optional) Filter the results. The `filter` block is documented below.
* `sort` - (Optional) Sort the results. The `sort` block is documented below.

---

`filter` supports the following arguments:

* `key` - (Required) Filter the tokens by this key. This may be one of `id`, `name`, `created_at`.
* `values` - (Required) A list of values to match against the `key` field.
* `match_by` - (Optional) One of `exact` (default), `re`, or `substring`.
* `all` - (Optional) Set to `true` to require that a field match all of the `values` instead of just one.

`sort` supports the following arguments:

* `key` - (Required) Sort the tokens by this key. This may be one of the keys listed in `filter`.
* `direction` - (Optional) The sort direction. This may be either `asc` or `desc`.

## Attributes Reference

* `tokens` - A list of tokens satisfying any `filter` and `sort` criteria. Each element contains:
  - `id` - The unique ID of the token.
  - `name` - The name of the token.
  - `created_at` - The date and time when the token was created.
