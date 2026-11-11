---
page_title: "DigitalOcean: digitalocean_action_gateway_toolbelt"
subcategory: "Action Gateway"
---

# digitalocean_action_gateway_toolbelt

Get information about the latest or a specific version of an
[Action Gateway](https://docs.digitalocean.com/products/managed-agents/action-gateway/)
toolbelt.

## Example Usage

```hcl
data "digitalocean_action_gateway_toolbelt" "search" {
  name = "search-toolbelt"
}

output "search_toolbelt" {
  value = "toolbelt:${data.digitalocean_action_gateway_toolbelt.search.reference}"
}
```

## Argument Reference

* `name` - (Required) The name of the toolbelt.
* `version` - (Optional) The version to read. Defaults to the latest version. A deleted toolbelt can only be read by an explicit version.

## Attributes Reference

The following attributes are exported:

* `id` - The version-pinned reference to the toolbelt.
* `display_name` - The human-readable label of the toolbelt.
* `description` - The description of the toolbelt.
* `tools` - The members of the version, each pinned as `<tool_slug>@<version>`.
* `reference` - The version-pinned reference to the toolbelt, such as `search-toolbelt@2`.
* `reference_latest` - The reference to the latest version of the toolbelt.
* `tool_count` - The number of tools in the toolbelt.
* `status` - The status of the toolbelt.
* `created_at` - When the version was created.
* `updated_at` - When the toolbelt was last modified.
