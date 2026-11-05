---
page_title: "DigitalOcean: digitalocean_dedicated_inferences"
subcategory: "Dedicated Inference"
---

# digitalocean\_dedicated\_inferences

Returns a list of dedicated inference endpoints in your DigitalOcean account,
with the ability to filter and sort the results. If no filters are specified, all
endpoints will be returned.

## Example Usage

```hcl
data "digitalocean_dedicated_inferences" "all" {}

output "all_endpoints" {
  value = data.digitalocean_dedicated_inferences.all.dedicated_inferences
}
```

### Filter by name

```hcl
data "digitalocean_dedicated_inferences" "filtered" {
  filter {
    key    = "name"
    values = ["my-inference"]
  }
}
```

### Filter by region

```hcl
data "digitalocean_dedicated_inferences" "by_region" {
  filter {
    key    = "region"
    values = ["tor1"]
  }

  sort {
    key       = "name"
    direction = "asc"
  }
}
```

## Argument Reference

* `filter` - (Optional) Filter the results. The `filter` block is documented below.
* `sort` - (Optional) Sort the results. The `sort` block is documented below.

---

`filter` supports the following arguments:

* `key` - (Required) Filter the dedicated inference endpoints by this key. This may be one of `id`, `name`, `region`, `status`, `vpc_uuid`, `public_endpoint_fqdn`, `private_endpoint_fqdn`, `created_at`, `updated_at`.
* `values` - (Required) A list of values to match against the `key` field.
* `match_by` - (Optional) One of `exact` (default), `re`, or `substring`. For string-typed fields, the match mode controls how the filter is applied.
* `all` - (Optional) Set to `true` to require that a field match all of the `values` instead of just one.

`sort` supports the following arguments:

* `key` - (Required) Sort the dedicated inference endpoints by this key. This may be one of the keys listed in `filter`.
* `direction` - (Optional) The sort direction. This may be either `asc` or `desc`.

## Attributes Reference

* `dedicated_inferences` - A list of dedicated inference endpoints satisfying any `filter` and `sort` criteria. Each element contains the following attributes:
  - `id` - The unique ID of the dedicated inference endpoint.
  - `name` - The name of the dedicated inference endpoint.
  - `region` - The region where the dedicated inference endpoint is deployed.
  - `status` - The current status of the dedicated inference endpoint.
  - `vpc_uuid` - The UUID of the VPC the dedicated inference endpoint is deployed in.
  - `public_endpoint_fqdn` - The fully-qualified domain name of the public endpoint, if enabled.
  - `private_endpoint_fqdn` - The fully-qualified domain name of the private endpoint.
  - `created_at` - The date and time when the dedicated inference endpoint was created.
  - `updated_at` - The date and time when the dedicated inference endpoint was last updated.
