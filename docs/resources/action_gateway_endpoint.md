---
page_title: "DigitalOcean: digitalocean_action_gateway_endpoint"
subcategory: "Action Gateway"
---

# digitalocean_action_gateway_endpoint

Provides a DigitalOcean [Action Gateway](https://docs.digitalocean.com/products/managed-agents/action-gateway/)
endpoint resource. An endpoint is a named lineage of immutable versions, and a
version is the tool selection, permission policy, and configuration that
Gateway applies to a request. This resource owns an endpoint and publishes the
version its `version` block describes.

Changing the `version` block publishes a new head version in place. The
endpoint keeps its ID, the previous version stays live and reachable through its
own pinned URL, and nothing is deleted. A revert is another apply with the
earlier configuration.

## Example Usage

```hcl
resource "digitalocean_action_gateway_endpoint" "support" {
  name = "support-agent"

  insights {
    logs   = true
    traces = false
  }

  version {
    actor_id = "user-123"
    tools    = ["jira_get_issue", "exa_web_search", "toolbelt:read-only@1"]

    policy {
      default_action = "ask"

      rules {
        tool   = "jira_get_issue"
        action = "allow"
      }
    }

    config = jsonencode({
      customInstructions = "Answer in the customer's language."
    })
  }
}

output "mcp_url" {
  value = digitalocean_action_gateway_endpoint.support.mcp_url
}

output "pinned_mcp_url" {
  value = digitalocean_action_gateway_endpoint.support.version_mcp_url
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the endpoint, unique among your team's active endpoints. It must contain 1 to 64 characters from `[A-Za-z0-9._-]`. Changing the name renames the endpoint in place.
* `insights` - (Optional) Insights consent, which applies to every version of the endpoint. See [Insights](#insights) below.
* `version` - (Required) The version the configuration describes. See [Version](#version) below. Changing it publishes a new head version.

### Insights

* `metrics` - (Optional) Whether the endpoint's metrics are sent to Insights.
* `logs` - (Optional) Whether the endpoint's logs are sent to Insights. Your team can require this to be on for every endpoint; turning it off is then refused.
* `traces` - (Optional) Whether the endpoint's traces are sent to Insights.

Each signal is optional and computed, so an omitted signal keeps the value the
endpoint already stores instead of resetting it.

### Version

* `actor_id` - (Optional) The actor the version's tool calls act as. It must contain 1 to 64 characters from `[A-Za-z0-9._-]`. No operation changes an endpoint's actor, so changing this forces a new endpoint.
* `tools` - (Optional) The tools and toolbelts the version enables, each optionally pinned as `<tool_slug>@<version>` or `toolbelt:<name>@<version>`. **Omit the argument to enable every tool, and set an empty list to enable none.** The API distinguishes the two, and so does this resource.
* `policy` - (Optional) The tool permission policy. Omitting it asks before every call. See [Policy](#policy) below.
* `config` - (Optional) A JSON object of options Gateway applies to the version, such as `preloadTools`, `customInstructions`, or `outputViews`. It is validated as JSON, normalized for comparison, and stored as supplied.
* `network` - (Optional) A VPC the version's tool calls are routed through. See [Network](#network) below.

### Policy

* `default_action` - (Optional) The action for calls that no rule matches: `allow`, `ask`, or `deny`. Empty means `ask`.
* `rules` - (Optional) Rules for specific tools and toolbelts, in the order the configuration lists them. Each rule supports:
  * `tool` - (Required) A tool slug, optionally pinned as `<tool_slug>@<version>`, or `toolbelt:<name>@<version>`.
  * `action` - (Required) One of `allow`, `ask`, or `deny`.
  * `match` - (Optional) Argument values a call must match for the rule to apply, keyed by argument name.

### Network

* `vpc_uuid` - (Required) The UUID of a VPC owned by your team.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The ID of the endpoint.
* `mcp_url` - The MCP URL that follows the endpoint's current head version. Its host can change when a new head version moves the endpoint to a VPC in another region, so do not treat it as a stable hostname.
* `version_mcp_url` - The MCP URL pinned to the current head version. It is stable only while Terraform owns the head.
* `is_default` - Whether this is the default endpoint that a generic `/mcp` request resolves to for your user.
* `created_at` - When the endpoint was created.
* `updated_at` - When the endpoint was last modified.

The `version` block additionally exports:

* `pinned_tools` - The tools and toolbelts the head version enables, each pinned by the API as `<tool_slug>@<version>` or `toolbelt:<name>@<version>`. Empty when every tool is enabled.
* `policy.0.default_action` - The resolved default action, which is `ask` when the configuration leaves it empty.
* `network.0.region` - The region of the VPC the version is attached to.
* `id` - The ID of the head version, `<endpoint UUID>.v<version number>`.
* `version_number` - The number of the head version.
* `promoted_from_version_number` - The version number the head was promoted from, when it was created by a promotion rather than a publish.

## Drift

`terraform plan` reads the endpoint's head version into the `version` block. A
version published from the console, PyDo, another workspace, or a promotion
therefore appears as a diff, and the next apply converges by publishing the
configured version as a new head. The resource never adopts an out-of-band
configuration silently, and it never deletes a version it did not publish.

## Limit on live versions

The resource never prunes versions. Each applied configuration change adds one
live version, an endpoint keeps at most 100 live versions, and the next publish
beyond that limit is refused with `endpoint_version_limit_exceeded`. Deleting
history is an out-of-band action, because a version identity is never reused.

## Promotion

Promotion is not modeled, because a promotion creates a new head whose identity
the caller cannot predict. Reverting the configuration to an earlier one
publishes it again as a new head, which has the same effect; only the
`promoted_from_version_number` provenance marker is out of reach.

## Import

An endpoint can be imported using its UUID, e.g.

```
terraform import digitalocean_action_gateway_endpoint.support 1b4d4e0f-9a7c-4d5b-8f3c-2a6e6e8f0a13
```

An imported endpoint plans no change when its head version matches the
configuration, and one publish when it does not. An import reports `tools` and
`pinned_tools` in the API's pinned spelling, such as `exa_web_search@1`, so a
configuration that matches an imported head names the pins too. An imported
**default** endpoint is managed like any other, but deleting it makes the next
generic `/mcp` request create a new default endpoint with a new ID, so a client
pinned to the deleted endpoint's own URL stops working.
