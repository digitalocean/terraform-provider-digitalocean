---
page_title: "DigitalOcean: digitalocean_action_gateway_connection"
subcategory: "Action Gateway"
---

# digitalocean_action_gateway_connection

Provides a DigitalOcean [Action Gateway](https://docs.digitalocean.com/products/managed-agents/action-gateway/)
connection resource. A connection authorizes a tool provider, such as GitHub,
Jira, or Linear, to act on an account for one actor ID. Sessions whose actor ID
matches the connection's `user_id` use it for that provider's tools.

-> **Note:** A connection belongs to the DigitalOcean user who owns the API
token Terraform uses, within that token's team. Other team members cannot use
it, even with the same `user_id`. To let team members create their own
connections from a shared credential, register the provider credential in the
control panel, make it reusable by the team, and pass its ID as
`credential_id`.

An OAuth connection starts `pending`. The account holder completes it by
opening `authorization.0.connect_url`, which is exported as a sensitive value
while the connection is pending. A connection through a team API-key
credential is active immediately.

## Example Usage

### OAuth through DigitalOcean's OAuth application

```hcl
resource "digitalocean_action_gateway_connection" "jira" {
  tool_provider = "jira"
  user_id       = "user-123"
  scopes        = ["write:jira-work"]

  connection_parameters = {
    site_url = "https://example.atlassian.net"
  }
}

output "jira_connect_url" {
  value     = try(digitalocean_action_gateway_connection.jira.authorization[0].connect_url, null)
  sensitive = true
}
```

### Team provider credential

```hcl
resource "digitalocean_action_gateway_connection" "linear" {
  tool_provider = "linear"
  user_id       = "support-bot"
  credential_id = "8d6f2b1a-3c4e-4d5f-9a8b-7c6d5e4f3a2b"
}
```

### Private network

```hcl
resource "digitalocean_vpc" "example" {
  name   = "example-vpc"
  region = "ric1"
}

resource "digitalocean_action_gateway_connection" "private_jira" {
  tool_provider = "jira"
  user_id       = "user-123"
  credential_id = "8d6f2b1a-3c4e-4d5f-9a8b-7c6d5e4f3a2b"

  connection_parameters = {
    site_url = "https://jira.internal.example.com"
  }

  vpc {
    vpc_uuid = digitalocean_vpc.example.id

    destination {
      host             = "jira.internal.example.com"
      port             = 443
      allowed_ip_cidrs = ["10.116.0.0/20"]
    }
  }
}
```

## Argument Reference

The following arguments are supported. The API cannot change a connection, so
changing any argument revokes the connection and creates a new one.

* `tool_provider` - (Required) The slug of the tool provider, such as `github` or `jira`. The [`digitalocean_action_gateway_providers`](../data-sources/action_gateway_providers) data source lists the available providers.
* `user_id` - (Required) The actor ID the connection acts for: 1 to 64 characters from `[A-Za-z0-9._-]`. Use the actor ID of the sessions that need this account.
* `scopes` - (Optional) The OAuth scopes to request. When omitted or empty, every scope the provider offers is requested. Ignored for API-key credentials.
* `connection_parameters` - (Optional) A map of values for the provider's connection parameters, such as `site_url` for Jira.
* `credential_id` - (Optional) The ID of a team provider credential, such as your own OAuth application or an API key, to connect through. When omitted, DigitalOcean's shared OAuth application is used.
* `vpc` - (Optional) Routes the connection's calls through a VPC. VPC networking must be enabled for your team.
  - `vpc_uuid` - (Required) The UUID of a VPC owned by your team.
  - `destination` - (Required) One or more private hosts the connection's calls may reach.
    - `host` - (Required) A host name or IP address that resolves only to private addresses.
    - `port` - (Required) The destination port.
    - `allowed_ip_cidrs` - (Optional) CIDR ranges that every address the host resolves to must fall within.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The ID of the connection.
* `status` - The status of the connection: `pending`, `active`, or `expired`.
* `tool_provider_display_name` - The human-readable name of the tool provider.
* `credential_kind` - The kind of credential the connection uses: `digitalocean_oauth`, `private_oauth`, or `team_api_key`.
* `owning_user_id` - The DigitalOcean user ID of the user who created the connection.
* `vpc.0.region` - The region of the connection's VPC.
* `authorization` - The authorization the account holder still has to complete. Set only while the connection is pending.
  - `status` - The status of the authorization.
  - `connect_url` - (Sensitive) The URL the account holder opens to authorize the provider. Share it only with the intended account holder.
  - `verification_code` - (Sensitive) The code the account holder confirms during authorization.
  - `expires_at` - When the authorization link expires.
* `granted_at` - When the account holder completed authorization.
* `created_at` - When the connection was created.
* `updated_at` - When the connection was last modified.

Each DigitalOcean user has one connection per tool provider and `user_id`. If
one that is not revoked already exists when the resource is created, creation
fails with that connection's ID so you can [import](#import) it instead.
Terraform does not take over a connection it did not create. For the same
reason, `create_before_destroy` cannot replace a connection: the new connection
would be the one it replaces.

If an OAuth connection's authorization link lapses before the account holder
completes it, the connection becomes `expired`. Replace the resource, for
example with `terraform apply -replace`, to start a new authorization.

## Import

A connection can be imported using its `id`, e.g.

```
terraform import digitalocean_action_gateway_connection.jira 39f9f768-a7f9-4c91-9557-3328c53dd054
```

Import reads `scopes` and `connection_parameters` from the API. If the
connection was granted more scopes than your configuration lists, or the API
stored a connection parameter in a different form, the next plan replaces the
connection. Run `terraform plan` after importing and update the configuration
to match the imported values before you apply.
