---
page_title: "DigitalOcean: digitalocean_action_gateway_provider_credential"
subcategory: "Action Gateway"
---

# digitalocean_action_gateway_provider_credential

Get information about an [Action Gateway](https://docs.digitalocean.com/products/managed-agents/action-gateway/)
provider credential, such as one registered in the control panel, to pass its
ID to a [`digitalocean_action_gateway_connection`](../resources/action_gateway_connection).
Look a credential up by its `id`, or by its `tool_provider` and `display_name`.

The lookup finds only credentials visible to the DigitalOcean user who owns
Terraform's API token: that user's own and those shared with the team. The API
never returns a credential's API key, client ID, or client secret, so this data
source does not export them.

## Example Usage

```hcl
data "digitalocean_action_gateway_provider_credential" "github" {
  tool_provider = "github"
  display_name  = "Production GitHub app"
}

resource "digitalocean_action_gateway_connection" "github" {
  tool_provider = data.digitalocean_action_gateway_provider_credential.github.tool_provider
  user_id       = "user-123"
  credential_id = data.digitalocean_action_gateway_provider_credential.github.id
}
```

## Argument Reference

One of the following must be provided:

* `id` - (Optional) The ID of the provider credential.
* `tool_provider` and `display_name` - (Optional) The slug of the credential's tool provider and the credential's label. Both must be set. The label is matched ignoring case.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the provider credential.
* `tool_provider` - The slug of the tool provider.
* `display_name` - The label of the credential.
* `kind` - The kind of credential, such as `oauth_client` or `api_key`.
* `status` - The status of the credential: `active` or `disabled`.
* `base_url` - The HTTPS origin of the tool provider's API, when one is configured.
* `reusable_by_team` - Whether everyone on the team can see and use the credential.
* `scope` - Who can use the credential: `user` for only its creator, or `team`.
* `redirect_url` - The redirect URI registered with an OAuth client. Empty when DigitalOcean's callback is used.
* `owning_user_id` - The DigitalOcean user ID of the credential's creator. Empty for a credential shared with the team.
* `created_at` - When the credential was created.
* `updated_at` - When the credential was last modified.
