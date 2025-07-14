---
page_title: "DigitalOcean: digitalocean_genai"
subcategory: "GenAI"
---

# digitalocean_genai_agent

Provides a resource to manage a DigitalOcean GenAI Agent. With this resource you can create, update, and delete agents, as well as update the agent's visibility status.

## Example Usage

```hcl
resource "digitalocean_genai_agent" "terraform-testing" {
  description = "Agent for testing update and delete functionality."
  instruction = "You are DigitalOcean's Solutions Architect Assistant, designed to help users find the perfect solution for their technical needs."
  model_uuid  = "d754f2d7-d1f0-11ef-bf8f-4e013e2ddde4"
  name        = "terraform-testing"
  project_id  = "84e1e297-ee40-41ac-95ff-1067cf2206e9"
  region      = "tor1"
  tags        = ["marketplace-agent-terraform"]
}
```

## Argument Reference

The following arguments are supported:

- **description** (Optional) - A description for the agent.
- **instruction** (Required) - The detailed instruction for the agent.
- **model_uuid** (Required) - The UUID of the agent's associated model.
- **name** (Required) - The name assigned to the agent.
- **project_id** (Required) - The project identifier for the agent.
- **region** (Required) - The region where the agent is deployed.
- **tags** (Optional) - A list of tags associated with the agent.
- **visibility** (Optional) - The visibility of the agent (e.g., "public" or "private").
- **anthropic_key_uuid** (Optional) - Anthropic API key UUID to use with Anthropic models.
- **knowledge_base_uuid** (Optional) - List of knowledge base UUIDs to attach to the agent.
- **open_ai_key_uuid** (Optional) - OpenAI API key UUID to use with OpenAI models.
- **anthropic_api_key** (Optional) - Anthropic API Key information block.
- **api_key_infos** (Optional) - List of API Key Info blocks.
- **api_keys** (Optional) - List of API Key blocks.
- **chatbot_identifiers** (Optional) - List of chatbot identifiers.
- **deployment** (Optional) - List of deployment blocks.
- **functions** (Optional) - List of function blocks.
- **agent_guardrail** (Optional) - List of agent guardrail blocks.
- **chatbot** (Optional) - Chatbot configuration block.
- **if_case** (Optional) - If case condition.
- **k** (Optional) - K value.
- **knowledge_bases** (Optional, Computed) - List of knowledge base blocks.
- **max_tokens** (Optional) - Maximum tokens allowed.
- **model** (Optional, Computed) - Model block.
- **open_ai_api_key** (Optional) - OpenAI API Key information block.
- **provide_citations** (Optional) - Whether the agent should provide citations.
- **retrieval_method** (Optional) - Retrieval method used.
- **route_created_by** (Optional) - User who created the route.
- **route_created_at** (Optional) - Timestamp when the route was created.
- **route_uuid** (Optional) - Route UUID.
- **route_name** (Optional) - Route name.
- **template** (Optional) - Agent template block.
- **temperature** (Optional) - Temperature setting.
- **top_p** (Optional) - Top-p sampling parameter.
- **url** (Optional) - URL for the agent.
- **user_id** (Optional) - User ID linked with the agent.

## Attributes Reference

After creation, the following attributes are exported:

- **agent_id** - The unique identifier of the agent.
- **created_at** - The timestamp when the agent was created.
- **updated_at** - The timestamp when the agent was last updated.
- **instruction** - The instruction used with the agent.
- **model_uuid** - The UUID of the agent's model.
- **name** - The name of the agent.
- **project_id** - The project associated with the agent.
- **region** - The region where the agent is deployed.
- **description** - The agent's description.
- **visibility** - The agent's visibility status.
- **tags** - The list of tags assigned to the agent.
- **if_case** - A condition parameter for agent behavior.
- **k** - An integer representing the "k" value.
- **max_tokens** - Maximum tokens allowed.
- **retrieval_method** - The retrieval method used.
- **route_created_at** - Timestamp for when the agent route was created.
- **route_created_by** - Who created the route.
- **route_uuid** - The unique identifier for the route.
- **route_name** - The name of the route.
- **temperature** - The temperature setting of the agent.
- **top_p** - The top-p sampling parameter.
- **url** - The URL associated with the agent.
- **user_id** - The user ID linked with the agent.
- **anthropic_key_uuid** - Anthropic API key UUID.
- **knowledge_base_uuid** - List of knowledge base UUIDs.
- **open_ai_key_uuid** - OpenAI API key UUID.
- **anthropic_api_key** - Anthropic API Key information.
- **api_key_infos** - List of API Key Info blocks.
- **api_keys** - List of API Key blocks.
- **chatbot_identifiers** - List of chatbot identifiers.
- **deployment** - List of deployment blocks.
- **functions** - List of function blocks.
- **agent_guardrail** - List of agent guardrail blocks.
- **chatbot** - Chatbot configuration block.
- **knowledge_bases** - List of knowledge base blocks.
- **model** - Model block.
- **open_ai_api_key** - OpenAI API Key information block.
- **provide_citations** - Whether the agent provides citations.
- **template** - Agent template block.

## Update Behavior

When the **visibility**, **description**, **instruction**, **k**, **max_tokens**, **model_uuid**, **name**, **open_ai_key_uuid**, **project_id**, **retrieval_method**, **region**, **tags**, **temperature**, or **top_p** attribute is changed, the provider invokes the update API endpoint to adjust the agent's configuration.

## Import

A DigitalOcean GenAI Agent can be imported using its UUID. For example:

```sh
terraform import digitalocean_genai_agent.terraform-testing 79292fb6-3627-11f0-bf8f-4e013e2ddde4
```

## Usage Notes

Changes to the agent's configuration, such as updating the instruction, description, or visibility, will trigger the corresponding update functions in the provider. This resource enables you to manage the complete lifecycle of a DigitalOcean GenAI Agent within your Terraform configuration.

---

# digitalocean_genai_knowledge_base

Provides a resource to manage a DigitalOcean GenAI Knowledge Base. With this resource you can create, update, and delete knowledge bases, as well as configure their data sources.

## Example Usage

```hcl
resource "digitalocean_genai_knowledge_base" "example" {
  name                 = "terraform-kb-example"
  project_id           = "84e1e297-ee40-41ac-95ff-1067cf2206e9"
  region               = "tor1"
  embedding_model_uuid = "22653204-79ed-11ef-bf8f-4e013e2ddde4"
  tags                 = ["documentation", "terraform-managed"]

  datasources {
    web_crawler_data_source {
      base_url        = "https://docs.digitalocean.com/products/kubernetes/"
      crawling_option = "SCOPED"
      embed_media     = true
    }
  }
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required) - The name of the knowledge base (between 2-32 characters).
- **project_id** (Required) - The unique identifier of the project to which the - **knowledge base belongs.
- **region** (Required) - The region where the knowledge base is deployed.
- **embedding_model_uuid** (Required) - The unique identifier of the - **embedding model.
- **datasources** (Required) - One or more data source configurations - **for the knowledge base.
  - **web_crawler_data_source** - Web crawler data source configuration:
    - **base_url** - The base URL for the web crawler to index.
    - **crawling_option** - The crawling option (e.g., "SCOPED").
    - **embed_media** - Whether to embed media content.
  - **spaces_data_source** - Spaces data source configuration:
    - **bucket_name** - The name of the Spaces bucket.
    - **item_path** - The path to items within the bucket.
    - **region** - The region of the Spaces bucket.
  - **file_upload_data_source** - File upload data source configuration.
- **vpc_uuid** (Optional) - The unique identifier of the VPC to which - **the knowledge base belongs.
- **database_id** (Optional) - The unique identifier of the DigitalOcean - **OpenSearch database this knowledge base will use.
- **is_public** (Optional) - Indicates whether the knowledge base is public or - **private.
- **tags** (Optional) - A list of tags associated with the knowledge base.

## Attributes Reference

After creation, the following attributes are exported:

- **id** - The unique identifier of the knowledge base.
- **name** - The name of the knowledge base.
- **project_id** - The project associated with the knowledge base.
- **region** - The region where the knowledge base is deployed.
- **vpc_uuid** - The VPC identifier associated with the knowledge base.
- **embedding_model_uuid** - The UUID of the embedding model used by the knowledge base.
- **datasources** - The data sources configured for the knowledge base.
- **created_at** - The timestamp when the knowledge base was created.
- **updated_at** - The timestamp when the knowledge base was last updated.
- **added_to_agent_at** - The timestamp when the knowledge base was added to an agent (if applicable).
- **database_id** - The OpenSearch database identifier.
- **is_public** - Whether the knowledge base is public.
- **last_indexing_job** - Information about the last indexing job:
  - **status** - The status of the indexing job.
  - **created_at** - When the indexing job was created.
  - **finished_at** - When the indexing job finished (if completed).
- **tags** - The list of tags assigned to the knowledge base.

## Update Behavior

When the **database_id**, **embedding_model_uuid**, **name**, **project_id**, **tags** or **uuid** attributes are changed, the provider invokes the update API endpoint to adjust the knowledge base's configuration.

## Import

A DigitalOcean GenAI Knowledge Base can be imported using its UUID. For example:

```sh
terraform import digitalocean_genai_knowledge_base.example a1b2c3d4-5678-90ab-cdef-1234567890ab
```

## Usage Notes

- Changes to **datasources**, **embedding_model_uuid**, **spaces_data_source**, **web_crawler_data_source**, **agent_uuid** and **vpc_uuid** will force recreation of the knowledge base.
- To add additional data sources after creation, use the `digitalocean_genai_knowledge_base_data_source` resource.
- To attach a knowledge base to an agent, use the `digitalocean_genai_agent_knowledge_base_attachment` resource.
