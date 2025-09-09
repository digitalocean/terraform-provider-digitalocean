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
## Example Usage: Fetching a Knowledge Base

```hcl
data "digitalocean_genai_knowledge_base" "example" {
  uuid = "a1b2c3d4-5678-90ab-cdef-1234567890ab"
}

output "kb_details" {
  value = data.digitalocean_genai_knowledge_base.example
}
```

## Example Usage: Fetching Knowledge Base Data Sources

```hcl
data "digitalocean_genai_knowledge_base_data_sources" "example" {
  knowledge_base_uuid = "a1b2c3d4-5678-90ab-cdef-1234567890ab"
}

output "kb_datasources" {
  value = data.digitalocean_genai_knowledge_base_data_sources.example.datasources
}
```

# digitalocean_genai_agent_versions

Provides a data source that retrieves all versions of an existing DigitalOcean GenAI Agent. Use this data source to query an agent by its unique identifier.

## Example Usage

```hcl
data "digitalocean_genai_agent_versions" "example" {
  agent_id = "79292fb6-3627-11f0-bf8f-4e013e2ddde4"
}

output "agent_detail" {
  value = data.digitalocean_genai_agent_versions.example
}
```

# digitalocean_genai_openai_api_keys

Provides a data source that lists all OpenAI API keys in your DigitalOcean account.

### Example Usage

```hcl
data "digitalocean_genai_openai_api_keys" "all" {}

output "all_openai_api_keys" {
  value = data.digitalocean_genai_openai_api_keys.all.openai_api_keys
}
```

### Attributes Reference

- **openai_api_keys** – List of OpenAI API keys.

---

## digitalocean_genai_openai_api_key

Provides a data source that retrieves a single OpenAI API key by UUID.

### Example Usage

```hcl
data "digitalocean_genai_openai_api_key" "by_id" {
  uuid = "your-openai-api-key-uuid"
}

output "openai_api_key_info" {
  value = data.digitalocean_genai_openai_api_key.by_id
}
```

### Argument Reference

- **uuid** (Required) – The UUID of the OpenAI API key.

### Attributes Reference

- **id** - The unique identifier of the OpenAI API key (same as uuid).
- **uuid** - The UUID of the OpenAI API key.
- **name** - The name of the API key.
- **created_at** - The timestamp when the API key was created.
- **updated_at** - The timestamp when the API key was last updated.
- **deleted_at** - The timestamp when the API key was deleted (if applicable).
- **created_by** - The user who created the API key.
- **models** - The list of models associated with the API key.

---

### digitalocean_genai_agents_by_openai_api_key

Provides a data source that lists all agents associated with a specific OpenAI API key.

### Example Usage

```hcl
data "digitalocean_genai_agents_by_openai_api_key" "by_key" {
  uuid = "your-openai-api-key-uuid"
}

output "agents_by_openai_key" {
  value = data.digitalocean_genai_agents_by_openai_api_key.by_key.agents
}
```

### Argument Reference

- **uuid** (Required) – The UUID of the OpenAI API key.

### Attributes Reference

- **agents** – List of agents associated with the OpenAI API key.

---

## Usage Notes

These data sources can be used to dynamically fetch details of existing GenAI resources into your Terraform configuration. You may reference exported attributes in other resources or outputs.

---

# digitalocean_genai_models

Provides a data source that lists all available GenAI models in DigitalOcean.

## Example Usage

```hcl
data "digitalocean_genai_models" "available_models" {}

output "all_models" {
  value = data.digitalocean_genai_models.available_models.models
}

output "model_names" {
  description = "Names of available models"
  value = [for model in data.digitalocean_genai_models.available_models.models : model.name]
}
```

## Attributes Reference

- **models** – List of available GenAI models. Each model contains:
  - **id** - The human-readable unique identifier of the model
  - **uuid** - The UUID of the model
  - **name** - The name of the model
  - **is_foundational** - Whether the model is a foundational model
  - **parent_uuid** - The UUID of the parent model (if applicable)
  - **upload_complete** - Whether the model upload is complete
  - **url** - The URL of the model
  - **created_at** - When the model was created
  - **updated_at** - When the model was last updated
  - **agreement** - License agreement information for the model:
    - **description** - Description of the agreement
    - **name** - Name of the agreement
    - **url** - URL to the full license text
    - **uuid** - UUID of the agreement
  - **version** - Version information of the model:
    - **major** - Major version number
    - **minor** - Minor version number
    - **patch** - Patch version number

## Usage Notes

This data source can be used to discover available GenAI models for use with agents or other GenAI resources. 


---

# digitalocean_genai_regions

Provides a data source that lists all available GenAI regions in DigitalOcean.

## Example Usage

```hcl
data "digitalocean_genai_regions" "available_regions" {}

output "all_regions" {
  value = data.digitalocean_genai_regions.available_regions.regions
}

output "region_names" {
  description = "Names of available regions"
  value = [for region in data.digitalocean_genai_regions.available_regions.regions : region.region]
}
```

## Attributes Reference

- **regions** – List of available GenAI regions. Each region contains:
  - **region** - The region identifier (e.g., "tor1")
  - **inference_url** - The inference URL for the region
  - **serves_batch** - Whether the region supports batch processing
  - **serves_inference** - Whether the region supports inference requests
  - **stream_inference_url** - The streaming inference URL for the region

## Usage Notes

This data source can be used to discover available regions for deploying GenAI resources like agents or knowledge bases. 