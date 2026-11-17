---
page_title: "DigitalOcean: digitalocean_action_gateway_providers"
subcategory: "Action Gateway"
---

# digitalocean_action_gateway_providers

Get information on the [Action Gateway](https://docs.digitalocean.com/products/managed-agents/action-gateway/)
tool providers, with the ability to filter and sort the results. Use it to
find a provider's slug, the scopes it offers, and the connection parameters a
[`digitalocean_action_gateway_connection`](../resources/action_gateway_connection)
needs. If no filters are specified, all providers are returned.

## Example Usage

```hcl
data "digitalocean_action_gateway_providers" "jira" {
  filter {
    key    = "name"
    values = ["jira"]
  }
}

output "jira_parameters" {
  value = data.digitalocean_action_gateway_providers.jira.providers[0].connection_parameters[*].key
}
```

## Argument Reference

* `filter` - (Optional) Filter the results.
  The `filter` block is documented below.

* `sort` - (Optional) Sort the results.
  The `sort` block is documented below.

`filter` supports the following arguments:

* `key` - (Required) Filter the providers by this key. This may be one of `name`, `display_name`, `description`, `auth_types`, and `scopes`.

* `values` - (Required) A list of values to match against the `key` field. Only retrieves providers
  where the `key` field takes on one or more of the values provided here.

* `match_by` - (Optional) One of `exact` (default), `re`, or `substring`. For string-typed fields, specify `re` to
  match by using the `values` as regular expressions, or specify `substring` to match by treating the `values` as
  substrings to find within the string field.

* `all` - (Optional) Set to `true` to require that a field match all of the `values` instead of just one or more of
  them. This is useful when matching against multi-valued fields such as lists or sets where you want to ensure
  that all of the `values` are present in the list or set.

`sort` supports the following arguments:

* `key` - (Required) Sort the providers by this key. This may be one of `name`, `display_name`, and `description`.

* `direction` - (Required) The sort direction. This may be either `asc` or `desc`.

## Attributes Reference

* `providers` - A list of tool providers satisfying any `filter` and `sort` criteria. Each provider has the following attributes:
  - `name` - The provider slug, used as `tool_provider` on connections.
  - `display_name` - The human-readable name of the provider.
  - `description` - The description of the provider.
  - `auth_types` - The credential kinds the provider accepts.
  - `scopes` - The OAuth scopes the provider offers.
  - `connection_parameters` - The parameters a connection to the provider accepts.
    - `key` - The parameter's key in `connection_parameters`.
    - `label` - The parameter's label.
    - `description` - The parameter's description.
    - `input_kind` - The kind of input the parameter takes.
    - `required` - Whether the parameter is required.
    - `allowed_values` - The values the parameter accepts, when restricted.
