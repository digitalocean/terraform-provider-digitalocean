terraform {
  required_providers {
    digitalocean = {
      source = "digitalocean/digitalocean"
    }
  }
}

provider "digitalocean" {
  # Set DIGITALOCEAN_TOKEN in your environment. Connections belong to the
  # DigitalOcean user who owns this token.
}

variable "actor_id" {
  type        = string
  description = "Actor ID of the sessions that use these connections."
  default     = "user-123"
}

# The GitHub provider's slug, scopes, and connection parameters.
data "digitalocean_action_gateway_providers" "github" {
  filter {
    key    = "name"
    values = ["github"]
  }
}

# A GitHub connection through DigitalOcean's OAuth application. It stays
# pending until the account holder opens the connect URL and approves access.
resource "digitalocean_action_gateway_connection" "github" {
  tool_provider = "github"
  user_id       = var.actor_id
}

# A reusable set of web search tools for sessions.
resource "digitalocean_action_gateway_toolbelt" "search" {
  name         = "search-toolbelt"
  display_name = "Search Tools"
  tools        = ["exa_web_search", "exa_web_fetch"]
}

output "github_scopes" {
  value = data.digitalocean_action_gateway_providers.github.providers[0].scopes
}

output "github_connection_status" {
  value = digitalocean_action_gateway_connection.github.status
}

output "github_connect_url" {
  value     = try(digitalocean_action_gateway_connection.github.authorization[0].connect_url, null)
  sensitive = true
}

output "search_toolbelt" {
  value = "toolbelt:${digitalocean_action_gateway_toolbelt.search.reference}"
}
