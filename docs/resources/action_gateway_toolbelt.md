---
page_title: "DigitalOcean: digitalocean_action_gateway_toolbelt"
subcategory: "Action Gateway"
---

# digitalocean_action_gateway_toolbelt

Provides a DigitalOcean [Action Gateway](https://docs.digitalocean.com/products/managed-agents/action-gateway/)
toolbelt resource. A toolbelt is a named, versioned set of tools that your team
can attach to sessions as `toolbelt:<name>@<version>`. A toolbelt contains no
account authorization; use [`digitalocean_action_gateway_connection`](action_gateway_connection)
to authorize tool providers.

## Example Usage

```hcl
resource "digitalocean_action_gateway_toolbelt" "search" {
  name         = "search-toolbelt"
  display_name = "Search Tools"
  description  = "Tools for searching and fetching public web pages."
  tools        = ["exa_web_search", "exa_web_fetch"]
}

output "search_toolbelt" {
  value = "toolbelt:${digitalocean_action_gateway_toolbelt.search.reference}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the toolbelt, unique among your team's active toolbelts. It must start with a lowercase letter and contain at most 64 characters from `[a-z0-9_-]`. `search` is reserved. Changing this forces a new toolbelt.
* `display_name` - (Optional) A human-readable label for the toolbelt, at most 128 bytes. Changing this forces a new toolbelt.
* `description` - (Optional) A description of the toolbelt, at most 255 bytes. Changing this forces a new toolbelt.
* `tools` - (Optional) The catalog tool slugs in the toolbelt, such as `exa_web_search`. Each may be pinned as `<tool_slug>@<version>`, and a pin must equal the tool's current released version. An unpinned tool is pinned to its current released version when it is added. Changing the tools creates a new version of the toolbelt.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The name of the toolbelt.
* `version` - The current version of the toolbelt.
* `reference` - The version-pinned reference to the toolbelt, such as `search-toolbelt@2`. Sessions reference it as `toolbelt:<reference>`.
* `reference_latest` - The reference to the latest version of the toolbelt.
* `pinned_tools` - The members of the current version, each pinned as `<tool_slug>@<version>`.
* `tool_count` - The number of tools in the toolbelt.
* `status` - The status of the toolbelt.
* `created_at` - When the toolbelt was created.
* `updated_at` - When the toolbelt was last modified.

## Import

A toolbelt can be imported using its `name`, e.g.

```
terraform import digitalocean_action_gateway_toolbelt.search search-toolbelt
```

After an import, `tools` lists every member pinned as `<tool_slug>@<version>`.
