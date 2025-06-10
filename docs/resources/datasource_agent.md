---
page_title: "DigitalOcean: digitalocean_agent"
subcategory: "Agents"
---

# digitalocean_agent

Provides a data source that retrieves details about an existing DigitalOcean Agent. Use this data source to query an agent by its unique identifier.

## Example Usage

```hcl
data "digitalocean_agent" "example" {
  agent_id = "79292fb6-3627-11f0-bf8f-4e013e2ddde4"
}

output "agent_detail" {
  value = data.digitalocean_agent.example
}
```

## Argument Reference

The following argument is supported:

- **agent_id** - (Required) The unique identifier of the agent to retrieve.

## Attributes Reference

The following attributes are exported:

- **uuid** - The unique identifier of the agent.
- **name** - The name assigned to the agent.
- **instruction** - The instruction configured for the agent.
- **model_uuid** - The UUID of the agent's associated model.
- **project_id** - The project identifier linked with the agent.
- **region** - The region where the agent is deployed.
- **description** - A description for the agent.
- **visibility** - The visibility of the agent (e.g., public or private).
- **created_at** - The timestamp when the agent was created (in RFC3339 format).
- **updated_at** - The timestamp when the agent was last updated (in RFC3339 format).
- **if_case** - A conditional parameter used for agent behavior.
- **k** - An integer representing the "k" value.
- **max_tokens** - The maximum number of tokens allowed.
- **retrieval_method** - The method used for data retrieval.
- **route_created_at** - The timestamp when the agent route was created.
- **route_created_by** - Information about who created the route.
- **route_uuid** - The unique identifier for the route.
- **route_name** - The name of the route.
- **tags** - A list of tags associated with the agent.
- **temperature** - The temperature setting of the agent.
- **top_p** - The top-p sampling parameter.
- **url** - A URL associated with the agent.
- **user_id** - The user ID linked with the agent.

## Usage Notes

This data source can be used to dynamically fetch the details of an existing agent into your Terraform configuration. You may reference exported attributes in other resources or outputs.

For example, to reference the agent's name:

```hcl
output "agent_name" {
  value = data.digitalocean_agent.example.name
}
```

This data source is useful for integrating agent details into your workflow or for performing validations against current configurations.

---