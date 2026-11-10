---
page_title: "DigitalOcean: digitalocean_action_gateway_endpoint"
subcategory: "Action Gateway"
---

# digitalocean_action_gateway_endpoint

Get information about an [Action Gateway](https://docs.digitalocean.com/products/managed-agents/action-gateway/)
endpoint and the version its head resolves to, by `id` or by `name`.

## Example Usage

### By name

```hcl
data "digitalocean_action_gateway_endpoint" "support" {
  name = "support-agent"
}

output "mcp_url" {
  value = data.digitalocean_action_gateway_endpoint.support.mcp_url
}

output "pinned_mcp_url" {
  value = data.digitalocean_action_gateway_endpoint.support.version_mcp_url
}
```

### By name, filtered by actor

```hcl
data "digitalocean_action_gateway_endpoint" "support" {
  name     = "support-agent"
  actor_id = "user-123"
}
```

### By ID

```hcl
data "digitalocean_action_gateway_endpoint" "support" {
  id = digitalocean_action_gateway_endpoint.support.id
}
```

## Argument Reference

The following arguments are supported. One of `id` or `name` is required.

* `id` - (Optional) The ID of the endpoint. Conflicts with `name`.
* `name` - (Optional) The name of the endpoint. Conflicts with `id`. Exactly one endpoint must carry the name; when several do, the lookup lists their IDs instead of picking one.
* `actor_id` - (Optional) An actor ID to filter a `name` lookup by. Endpoints bound to a different actor are ignored.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `mcp_url` - The MCP URL that follows the endpoint's current head version.
* `version_mcp_url` - The MCP URL pinned to the current head version.
* `is_default` - Whether this is the default endpoint that a generic `/mcp` request resolves to for your user.
* `created_at` - When the endpoint was created.
* `updated_at` - When the endpoint was last modified.
* `version` - The head version the endpoint resolves to, with `actor_id`, `tools`, `pinned_tools`, `policy`, `config`, `network`, `id`, `version_number`, and `promoted_from_version_number` as described on the [`digitalocean_action_gateway_endpoint` resource](action_gateway_endpoint) page. An endpoint with no live version exports no `version` block.
