---
page_title: "DigitalOcean: digitalocean_action_gateway_connection"
subcategory: "Action Gateway"
---

# digitalocean_action_gateway_connection

Get information about an [Action Gateway](https://docs.digitalocean.com/products/managed-agents/action-gateway/)
connection, such as one created in the control panel. Look a connection up by
its `id`, or by its `tool_provider` and `user_id`.

The lookup finds only connections visible to the DigitalOcean user who owns
Terraform's API token.

## Example Usage

```hcl
data "digitalocean_action_gateway_connection" "github" {
  tool_provider = "github"
  user_id       = "user-123"
}

output "github_connection_status" {
  value = data.digitalocean_action_gateway_connection.github.status
}
```

## Argument Reference

One of the following must be provided:

* `id` - (Optional) The ID of the connection.
* `tool_provider` and `user_id` - (Optional) The slug of the connection's tool provider and the actor ID it acts for. Both must be set, and `user_id` must match exactly.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the connection.
* `tool_provider` - The slug of the tool provider.
* `tool_provider_display_name` - The human-readable name of the tool provider.
* `user_id` - The actor ID the connection acts for.
* `status` - The status of the connection: `pending`, `active`, or `expired`.
* `scopes` - The OAuth scopes granted, or requested while pending.
* `connection_parameters` - The provider's connection parameters.
* `credential_kind` - The kind of credential the connection uses: `digitalocean_oauth`, `private_oauth`, or `team_api_key`.
* `credential_id` - The ID of the team provider credential the connection uses, empty for `digitalocean_oauth`.
* `owning_user_id` - The DigitalOcean user ID of the user who created the connection.
* `vpc` - The private network the connection's calls use, when configured.
  - `vpc_uuid` - The UUID of the VPC.
  - `region` - The region of the VPC.
  - `destination` - The private hosts the connection's calls may reach, each with `host`, `port`, and `allowed_ip_cidrs`.
* `authorization` - The authorization the account holder still has to complete. Set only while the connection is pending.
  - `status` - The status of the authorization.
  - `connect_url` - (Sensitive) The URL the account holder opens to authorize the provider.
  - `verification_code` - (Sensitive) The code the account holder confirms during authorization.
  - `expires_at` - When the authorization link expires.
* `granted_at` - When the account holder completed authorization.
* `created_at` - When the connection was created.
* `updated_at` - When the connection was last modified.
