---
page_title: "DigitalOcean: digitalocean_action_gateway_mcp_server"
subcategory: "Action Gateway"
---

# digitalocean_action_gateway_mcp_server

Provides a DigitalOcean [Action Gateway](https://docs.digitalocean.com/products/managed-agents/action-gateway/)
MCP server resource. It registers your own remote MCP server with your team,
discovers its tools, and enables the ones you choose. Enabled tools are added
to your team's catalog with slugs of the form `<server_ref>_<tool>`.

## Example Usage

### API key sent to the server

```hcl
resource "digitalocean_action_gateway_mcp_server" "docs" {
  server_ref            = "docs"
  endpoint              = "https://mcp.example.com/mcp"
  description           = "Internal documentation search."
  credential_ref_source = "secret"
  api_key               = var.docs_mcp_api_key

  enabled_tools = ["docs_search", "docs_get_page"]
}
```

### Per-user OAuth

```hcl
resource "digitalocean_action_gateway_mcp_server" "tickets" {
  server_ref            = "tickets"
  endpoint              = "https://mcp.tickets.example.com/mcp"
  credential_ref_source = "connection"
  oauth_client_id       = var.tickets_oauth_client_id
  oauth_client_secret   = var.tickets_oauth_client_secret
  oauth_authorize_url   = "https://mcp.tickets.example.com/oauth/authorize"
  oauth_token_url       = "https://mcp.tickets.example.com/oauth/token"
  oauth_scopes          = ["tickets:read"]

  sync_user_id  = "user-123"
  enabled_tools = ["tickets_search"]
}
```

## Argument Reference

The following arguments are supported. Only `description`, `enabled_tools`,
and `sync_user_id` can change in place; changing any other argument registers
the server again.

* `server_ref` - (Required) The identifier for the server, unique within your team. It must start with a lowercase letter and contain at most 64 characters from `[a-z0-9-]`.
* `endpoint` - (Required) The HTTPS URL of the server's MCP endpoint. Its host must resolve only to public IP addresses.
* `transport` - (Optional) The MCP transport. Only `streamable_http`, the default, is supported.
* `credential_ref_source` - (Optional) How requests to the server authenticate: `none` (the default), `secret`, or `connection`.
* `credential_ref` - (Optional) For `secret`: a reference to a secret your team stores with DigitalOcean. Conflicts with `api_key`.
* `api_key` - (Optional, Sensitive) For `secret`: the key or token itself, sent in the `Authorization` header. It is write-only, so Terraform cannot detect changes made outside of Terraform. Conflicts with `credential_ref`.
* `oauth_client_id` - (Optional) For `connection`: the client ID of your team's OAuth client for the server. Each user of the server then authorizes individually.
* `oauth_client_secret` - (Optional, Sensitive) For `connection`: the client secret of your team's OAuth client. It is write-only.
* `oauth_authorize_url` - (Optional) For `connection`: the OAuth authorization endpoint, on the endpoint's host or another host under the same registrable domain.
* `oauth_token_url` - (Optional) For `connection`: the OAuth token endpoint, under the same rules as `oauth_authorize_url`.
* `oauth_scopes` - (Optional) For `connection`: the scopes to request from each user, 1 to 64 of them.
* `oauth_authorization_ttl_seconds` - (Optional) How long, in seconds, a user's authorization may be reused before they must consent again. `0` means 30 days.
* `description` - (Optional) The description shown on the server's catalog card, at most 1024 characters.
* `enabled_tools` - (Optional) The slugs of the discovered tools to enable, such as `docs_search`. Every other tool of the server is disabled. When omitted, Terraform does not manage which tools are enabled. Changing this discovers the server's tools again before enabling them.
* `sync_user_id` - (Optional) For `connection`: the actor ID of a user who has authorized the server. Discovering and verifying the tools of such a server requires it. Adding or changing it discovers the server's tools again.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The `server_ref` of the server.
* `credential_ref` - The reference to the secret the server uses, including one DigitalOcean created to store `api_key`.
* `sync_status` - The status of the most recent tool discovery.
* `sync_error` - The error from the most recent tool discovery, if any.
* `last_synced_at` - When the server's tools were last discovered.
* `protocol_version` - The MCP protocol version the server reported.
* `tool_count` - The number of tools discovered on the server.
* `discovered_tools` - The tools discovered on the server.
  - `name` - The tool's name on the server.
  - `tool_slug` - The tool's catalog slug, used in `enabled_tools`.
  - `description` - The tool's description.
  - `enabled` - Whether the tool is enabled.
  - `quarantined` - Whether the tool is quarantined.
  - `quarantine_reason` - Why the tool is quarantined.
* `created_at` - When the server was registered.
* `updated_at` - When the server was last modified.

When the server is created, Terraform discovers its tools and then enables
`enabled_tools`. If discovery cannot complete, for example because
`sync_user_id` has not authorized the server yet, Terraform reports a warning
instead of failing. `enabled_tools` then stays empty, and the next
`terraform apply` discovers and enables the configured tools again. A server
registered without `sync_user_id` discovers its tools once you set it.

Setting `enabled_tools` to an empty list does not disable the server's tools.
To disable every tool, disable them in the control panel or remove the server.

## Import

An MCP server can be imported using its `server_ref`, e.g.

```
terraform import digitalocean_action_gateway_mcp_server.docs docs
```

The write-only arguments `api_key`, `oauth_client_id`, `oauth_client_secret`,
and `oauth_token_url` are not imported.
