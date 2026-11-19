---
page_title: "DigitalOcean: digitalocean_action_gateway_provider_credential"
subcategory: "Action Gateway"
---

# digitalocean_action_gateway_provider_credential

Provides a DigitalOcean [Action Gateway](https://docs.digitalocean.com/products/managed-agents/action-gateway/)
provider credential resource. A provider credential is an OAuth client or an
API key that your team registers for a tool provider, such as GitHub, Jira, or
Linear, to use instead of DigitalOcean's shared OAuth application. A
[`digitalocean_action_gateway_connection`](action_gateway_connection) connects
through it by its `credential_id`.

~> **Warning:** Destroying a provider credential also deletes every connection
that uses it, including connections Terraform does not manage, and their users
lose access to the provider. The API cannot restore a deleted credential. Changing
`tool_provider`, `kind`, `client_id`, or `redirect_url`, or turning
`reusable_by_team` off, replaces the credential, which destroys it first. Set
[`prevent_destroy`](https://developer.hashicorp.com/terraform/language/meta-arguments/lifecycle#prevent_destroy)
on credentials you cannot afford to lose, as the examples below do. Rotating a
secret or renaming a credential never replaces it.

-> **Note:** The API never returns `api_key` or `client_secret`. Terraform
keeps the values you configure in its state, marked sensitive, so protect the
state accordingly.

## Example Usage

### OAuth client

```hcl
variable "github_client_secret" {
  type      = string
  sensitive = true
}

resource "digitalocean_action_gateway_provider_credential" "github" {
  tool_provider = "github"
  kind          = "oauth_client"
  display_name  = "Production GitHub app"
  client_id     = "Iv1.0123456789abcdef"
  client_secret = var.github_client_secret

  lifecycle {
    prevent_destroy = true
  }
}

resource "digitalocean_action_gateway_connection" "github" {
  tool_provider = digitalocean_action_gateway_provider_credential.github.tool_provider
  user_id       = "user-123"
  credential_id = digitalocean_action_gateway_provider_credential.github.id
}
```

### API key

A connection through an API-key credential is active immediately, without an
authorization step.

```hcl
variable "linear_api_key" {
  type      = string
  sensitive = true
}

resource "digitalocean_action_gateway_provider_credential" "linear" {
  tool_provider = "linear"
  kind          = "api_key"
  display_name  = "Production key"
  api_key       = var.linear_api_key

  lifecycle {
    prevent_destroy = true
  }
}

resource "digitalocean_action_gateway_connection" "linear" {
  tool_provider = digitalocean_action_gateway_provider_credential.linear.tool_provider
  user_id       = "support-bot"
  credential_id = digitalocean_action_gateway_provider_credential.linear.id
}
```

### API key with credential parameters

```hcl
resource "digitalocean_action_gateway_provider_credential" "example" {
  tool_provider = "example-provider"
  kind          = "api_key"
  display_name  = "Production key"
  api_key       = var.example_api_key

  credential_parameters = {
    subdomain = "example"
  }
}
```

## Argument Reference

The following arguments are supported:

* `tool_provider` - (Required) The slug of the tool provider the credential is for, such as `github` or `jira`. The [`digitalocean_action_gateway_providers`](../data-sources/action_gateway_providers) data source lists the available providers and, in `auth_types`, the kinds of credential each accepts. Changing this replaces the credential.
* `kind` - (Required) The kind of credential: `oauth_client` for an OAuth client you created at the tool provider, or `api_key` for a key or token. Changing this replaces the credential.
* `display_name` - (Required) A label of 1 to 255 characters, unique among your team's credentials for the tool provider, ignoring case.
* `api_key` - (Optional, Sensitive) The key or token sent to the tool provider. Required for `kind = "api_key"` and not accepted otherwise. Changing it rotates the key in place: the credential keeps its ID, and its connections use the new key on their next call.
* `client_id` - (Optional) The client ID of your OAuth client. Required for `kind = "oauth_client"` and not accepted otherwise. A credential is bound to the client ID it was registered with, so changing this replaces the credential.
* `client_secret` - (Optional, Sensitive) The client secret of your OAuth client. Required for `kind = "oauth_client"` and not accepted otherwise. Changing it rotates the secret in place.
* `base_url` - (Optional) The HTTPS origin of the tool provider's API, such as `https://api.example.com`. Accepted only by tool providers that allow one. Conflicts with `credential_parameters`. Removing it from the configuration leaves the stored origin in place.
* `credential_parameters` - (Optional) A map of values for the credential parameters the tool provider declares, which the [providers API](https://docs.digitalocean.com/reference/api/) lists as `credential_parameters`. The API validates them and derives `base_url` from them. Accepted only for `kind = "api_key"`. Conflicts with `base_url`.
* `reusable_by_team` - (Optional) Whether everyone on your team can see and use the credential instead of only the DigitalOcean user who owns Terraform's API token. Defaults to `false`. Sharing a credential may require the team owner role. Sharing cannot be undone, so changing this from `true` to `false` replaces the credential.
* `redirect_url` - (Optional) The HTTPS redirect URI registered with your OAuth client, for `kind = "oauth_client"` only. Custom OAuth redirects must be enabled for your account. Omit it to use DigitalOcean's callback. Changing this replaces the credential.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The ID of the provider credential, used as `credential_id` on connections.
* `status` - The status of the credential: `active` or `disabled`. Connections cannot use a disabled credential.
* `scope` - Who can use the credential: `user` for only its creator, or `team`.
* `owning_user_id` - The DigitalOcean user ID of the credential's creator. Empty for a credential shared with the team.
* `created_at` - When the credential was created.
* `updated_at` - When the credential was last modified.

## Import

A provider credential can be imported using its `id`, e.g.

```
terraform import digitalocean_action_gateway_provider_credential.github example-credential-id
```

The API does not return a credential's secret or its OAuth client ID, so an
imported credential has neither in state, and the first plan after an import
shows them being set. Applying it sends the configured `api_key`, or
`client_id` and `client_secret`, as an in-place update that keeps the
credential and its connections. For an OAuth client the API accepts that update
only when `client_id` is the ID the credential was registered with, so a
credential imported into the wrong configuration fails to apply instead of
being overwritten.

Import reads `tool_provider`, `kind`, and `redirect_url` from the API. If one
differs from your configuration, the next plan replaces the credential. Run
`terraform plan` after importing and make sure it shows an in-place update, not
a replacement, before you apply.
