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

- **agent_id** (Required) – The unique identifier of the agent to retrieve.

## Attributes Reference

All fields below are exported and may be referenced:

- **uuid** – The unique identifier of the agent.
- **name** – The name assigned to the agent.
- **instruction** – The instruction configured for the agent.
- **model_uuid** – The UUID of the agent's associated model.
- **project_id** – The project identifier linked with the agent.
- **region** – The region where the agent is deployed.
- **description** – A description for the agent.
- **tags** – A list of tags associated with the agent.
- **visibility** – The visibility of the agent (e.g., public or private).
- **anthropic_key_uuid** – Anthropic API key UUID to use with Anthropic models.
- **knowledge_base_uuid** – List of knowledge base UUIDs attached to the agent.
- **open_ai_key_uuid** – OpenAI API key UUID to use with OpenAI models.
- **anthropic_api_key** – Anthropic API Key information block.
- **api_key_infos** – List of API Key Info blocks.
- **api_keys** – List of API Key blocks.
- **chatbot_identifiers** – List of chatbot identifiers.
- **deployment** – List of deployment blocks.
- **functions** – List of function blocks.
- **agent_guardrail** – List of agent guardrail blocks.
- **chatbot** – Chatbot configuration block.
- **if_case** – If case condition.
- **k** – K value.
- **knowledge_bases** – List of knowledge base blocks.
- **max_tokens** – Maximum tokens allowed.
- **model** – Model block.
- **open_ai_api_key** – OpenAI API Key information block.
- **provide_citations** – Whether the agent should provide citations.
- **retrieval_method** – Retrieval method used.
- **route_created_by** – User who created the route.
- **route_created_at** – Timestamp when the route was created.
- **route_uuid** – Route UUID.
- **route_name** – Route name.
- **template** – Agent template block.
- **temperature** – Temperature setting.
- **top_p** – Top-p sampling parameter.
- **url** – URL for the agent.
- **user_id** – User ID linked with the agent.
- **created_at** – The timestamp when the agent was created (in RFC3339 format).
- **updated_at** – The timestamp when the agent was last updated (in RFC3339 format).

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