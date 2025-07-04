---
page_title: "DigitalOcean: digitalocean_genai_agent"
subcategory: "GenAI"
---

# digitalocean_genai_agent

Provides a data source that retrieves details about an existing DigitalOcean GenAI Agent. Use this data source to query an agent by its unique identifier.

## Example Usage

```hcl
data "digitalocean_genai_agent" "example" {
  agent_id = "79292fb6-3627-11f0-bf8f-4e013e2ddde4"
}

output "agent_detail" {
  value = data.digitalocean_genai_agent.example
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
- **knowledge_base_uuid**– List of knowledge base UUIDs attached to the agent.
- **open_ai_key_uuid** – OpenAI API key UUID to use with OpenAI models.
- **anthropic_api_key** – Anthropic API Key information block.
- **api_key_infos** – List of API Key Info blocks.
- **api_keys** – List of API Key blocks.
- **chatbot_identifiers**– List of chatbot identifiers.
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

This data source is useful for integrating agent details into your workflow or for performing validations against current configurations.


# digitalocean_genai_knowledge_base

Provides a data source that retrieves details about an existing DigitalOcean GenAI Knowledge Base. Use this data source to query a knowledge base by its unique identifier (UUID).

## Example Usage

```hcl
data "digitalocean_genai_knowledge_base" "example" {
  uuid = "a1b2c3d4-5678-90ab-cdef-1234567890ab"
}

output "kb_details" {
  value = data.digitalocean_genai_knowledge_base.example
}
```

## Argument Reference

The following argument is supported:

- **uuid** (Required) – The unique identifier of the knowledge base to retrieve.

## Attributes Reference

All fields below are exported and may be referenced:

- **id** - The unique identifier of the knowledge base (same as uuid).
- **name** - The name assigned to the knowledge base.
- **project_id** - The unique identifier of the project to which the knowledge base belongs.
- **region** - The region where the knowledge base is deployed.
- **vpc_uuid** - The unique identifier of the VPC to which the knowledge base belongs (if applicable).
- **created_at** - The timestamp when the knowledge base was created (in RFC3339 format).
- **updated_at** - The timestamp when the knowledge base was last updated (in RFC3339 format).
- **added_to_agent_at** - The timestamp when the knowledge base was added to an agent (if applicable).
- **database_id** - The unique identifier of the DigitalOcean OpenSearch database used by this knowledge base.
- **embedding_model_uuid** - The unique identifier of the embedding model used by the knowledge base.
- **is_public** - Indicates whether the knowledge base is public or private.
- **tags** - A list of tags associated with the knowledge base.
- **datasources** - A list of data sources configured for the knowledge base, each containing:
  - **web_crawler_data_source** - Details of web crawler data sources:
    - **base_url** - The base URL for the web crawler to index.
    - **crawling_option** - The crawling option (e.g., "SCOPED").
    - **embed_media** - Whether to embed media content.
  - **spaces_data_source** - Details of Spaces data sources:
    - **bucket_name** - The name of the Spaces bucket.
    - **item_path** - The path to items within the bucket.
    - **region** - The region of the Spaces bucket.
  - **file_upload_data_source** - Details of file upload data sources.
- **last_indexing_job** - Information about the last indexing job for the knowledge base

## Usage Notes

This data source can be used to dynamically fetch the details of an existing knowledge base into your Terraform configuration. You may reference exported attributes in other resources or outputs.

For example, to create an agent with an existing knowledge base:

```hcl
data "digitalocean_genai_knowledge_base" "existing" {
  uuid = "a1b2c3d4-5678-90ab-cdef-1234567890ab"
}

resource "digitalocean_genai_agent_knowledge_base_attachment" "example" {
  agent_uuid          = digitalocean_genai_agent.example.id
  knowledge_base_uuid = data.digitalocean_genai_knowledge_base.existing.id
}
```
