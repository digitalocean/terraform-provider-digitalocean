---
page_title: "DigitalOcean: digitalocean_gradientai"
subcategory: "GradientAI"
---

# digitalocean_gradientai_agent

Provides a resource to manage a DigitalOcean Gradient AI Agent. With this resource you can create, update, and delete agents, as well as update the agent's visibility status.

## Example Usage


```hcl
resource "digitalocean_gradientai_agent" "terraform-testing" {
  description = "Agent for testing update and delete functionality."
  instruction = "You are DigitalOcean's Solutions Architect Assistant, designed to help users find the perfect solution for their technical needs."
  model_uuid  = "d754f2d7-d1f0-11ef-bf8f-4e013e2ddde4"
  name        = "terraform-testing"
  project_id  = "84e1e297-ee40-41ac-95ff-1067cf2206e9"
  region      = "tor1"
  tags        = ["marketplace-agent-terraform"]

  deployment {
    visibility = "VISIBILITY_PRIVATE" # Valid values: VISIBILITY_PRIVATE, VISIBILITY_PUBLIC, VISIBILITY_PLAYGROUND, VISIBILITY_DISABLED, VISIBILITY_UNKNOWN
  }
}
```

> **Note:** The `visibility` argument must be set inside the `deployment` block, and must use one of the following values: `VISIBILITY_PRIVATE`, `VISIBILITY_PUBLIC`, `VISIBILITY_PLAYGROUND`, `VISIBILITY_DISABLED`, or `VISIBILITY_UNKNOWN`.

## Argument Reference

The following arguments are supported:

- **description** (Optional) - A description for the agent.
- **instruction** (Required) - The detailed instruction for the agent.
- **model_uuid** (Required) - The UUID of the agent's associated model.
- **name** (Required) - The name assigned to the agent.
- **project_id** (Required) - The project identifier for the agent.
- **region** (Required) - The region where the agent is deployed.
- **tags** (Optional) - A list of tags associated with the agent.
- **deployment** (Optional) - One or more deployment blocks. The `visibility` argument must be set inside this block. Valid values for `visibility` are: `VISIBILITY_PRIVATE`, `VISIBILITY_PUBLIC`, `VISIBILITY_PLAYGROUND`, `VISIBILITY_DISABLED`, or `VISIBILITY_UNKNOWN`.
- **anthropic_key_uuid** (Optional) - Anthropic API key UUID to use with Anthropic models.
- **knowledge_base_uuid** (Optional) - List of knowledge base UUIDs to attach to the agent.
- **open_ai_key_uuid** (Optional) - OpenAI API key UUID to use with OpenAI models.
- **workspace_uuid** (Optional) - Identifier for the workspace.
- **anthropic_api_key** (Optional) - Anthropic API Key information block.
- **api_key_infos** (Optional) - List of API Key Info blocks.
- **api_keys** (Optional) - List of API Key blocks.
- **chatbot_identifiers** (Optional) - List of chatbot identifiers.
- **deployment** (Optional) - List of deployment blocks.
- **functions** (Optional) - List of function blocks.

**agent_guardrail** (Optional) - List of agent guardrail blocks. Each block supports the following:
  - **agent_uuid** (Optional) - The UUID of the agent the guardrail is attached to.
  - **created_at** (Computed) - Timestamp when the guardrail was created.
  - **default_response** (Optional) - The default response if the guardrail is triggered.
  - **description** (Optional) - Description of the guardrail.
  - **guardrail_uuid** (Optional) - The unique identifier for the guardrail.
  - **is_attached** (Computed) - Whether the guardrail is currently attached to the agent (boolean).
  - **is_default** (Optional) - Whether this guardrail is the default (boolean).
  - **name** (Optional) - Name of the guardrail.
  - **priority** (Optional) - Priority of the guardrail (integer).
  - **type** (Optional) - The type of guardrail. Possible values:
    - "GUARDRAIL_TYPE_UNKNOWN" (default)
    - "GUARDRAIL_TYPE_JAILBREAK"
    - "GUARDRAIL_TYPE_SENSITIVE_DATA"
    - "GUARDRAIL_TYPE_CONTENT_MODERATION"
  - **updated_at** (Computed) - Timestamp when the guardrail was last updated.
  - **uuid** (Optional) - The unique identifier for the guardrail.

#### Example: Adding an Agent Guardrail

```hcl
resource "digitalocean_gradientai_agent" "with_guardrail" {
  # ... other required arguments ...
  deployment {
    visibility = "VISIBILITY_PRIVATE"
  }
  agent_guardrail {
    guardrail_uuid   = "12345678-1234-1234-1234-123456789012"
    name             = "My Guardrail"
    type             = "GUARDRAIL_TYPE_UNKNOWN"
    priority         = 123
    default_response = "example string."
    description      = "example string"
    is_default       = true
  }
}
```
# digitalocean_gradientai_function

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
- **workspace_uuid** - Identifier for the workspace.
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

# digitalocean_gradientai_function

We can pick up the agent id from the agent terraform resource and input, output schema have json values as currently there is no defined schema available.
Checkout the following API docs - https://docs.digitalocean.com/reference/api/digitalocean/#tag/GradientAI-Platform/operation/gradient_ai_attach_agent_function

```hcl

resource "digitalocean_gradientai_function" "check" {
  agent_id       = digitalocean_gradientai_agent.terraform-testing.id
  description    = "Adding a function route and this will also tell temperature"
  faas_name      = "default/testing"
  faas_namespace = "fn-b90faf52-2b42-49c2-9792-75edfbb6f397"
  function_name  = "terraform-tf-complete"
  input_schema   = <<EOF
  {
      "parameters": [
          {
          "in": "query",
          "name": "zipCode",
          "schema": {
              "type": "string"
          },
          "required": false,
          "description": "The ZIP code for which to fetch the weather"
          },
          {
          "name": "measurement",
          "schema": {
              "enum": [
              "F",
              "C"
              ],
              "type": "string"
          },
          "required": false,
          "description": "The measurement unit for temperature (F or C)",
          "in": "query"
          }
      ]
  }
  EOF

  output_schema = <<EOF
  {
      "properties": [
          {
          "name": "temperature",
          "type": "number",
          "description": "The temperature for the specified location"
          },
          {
          "name": "measurement",
          "type": "string",
          "description": "The measurement unit used for the temperature (F or C)"
          },
          {
          "name": "conditions",
          "type": "string",
          "description": "A description of the current weather conditions (Sunny, Cloudy, etc)"
          }
      ]
  }
  EOF
}
```

## Attributes Reference

After creation, the following attributes are exported:

- **agent_id** - The unique identifier of the agent.
- **description** - Description for the function
- **faas_name** - The name of the function in the DigitalOcean functions platform
- **faas_namespace** - The namespace of the function in the DigitalOcean functions platform
- **function_name** - The name for function to be assigned inside agent, two functions inside agent cannot have same name
- **input_schema** - The input schema associated with the function.
- **output_schema** - The output schema associated with the function.

**input_schema** and **output_schema** have a json input please check out this docs for more clarity - https://docs.digitalocean.com/reference/api/digitalocean/#tag/GradientAI-Platform/operation/gradient_ai_attach_agent_function

## Import

A DigitalOcean Gradient AI Agent can be imported using its UUID. For example:

```sh
terraform import digitalocean_gradientai_agent.terraform-testing 79292fb6-3627-11f0-bf8f-4e013e2ddde4
```

## Usage Notes

Changes to the agent's configuration, such as updating the instruction, description, or visibility, will trigger the corresponding update functions in the provider. This resource enables you to manage the complete lifecycle of a DigitalOcean Gradient AI Agent within your Terraform configuration.

---

# digitalocean_gradientai_knowledge_base

Provides a resource to manage a DigitalOcean Gradient AI Knowledge Base. With this resource you can create, update, and delete knowledge bases, as well as configure their data sources.

## Example Usage

```hcl
resource "digitalocean_gradientai_knowledge_base" "example" {
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
- **project_id** (Required) - The unique identifier of the project to which the - \*\*knowledge base belongs.
- **region** (Required) - The region where the knowledge base is deployed.
- **embedding_model_uuid** (Required) - The unique identifier of the - \*\*embedding model.
- **datasources** (Required) - One or more data source configurations - \*\*for the knowledge base.
  - **web_crawler_data_source** - Web crawler data source configuration:
    - **base_url** - The base URL for the web crawler to index.
    - **crawling_option** - The crawling option (e.g., "SCOPED").
    - **embed_media** - Whether to embed media content.
  - **spaces_data_source** - Spaces data source configuration:
    - **bucket_name** - The name of the Spaces bucket.
    - **item_path** - The path to items within the bucket.
    - **region** - The region of the Spaces bucket.
  - **file_upload_data_source** - File upload data source configuration.
- **vpc_uuid** (Optional) - The unique identifier of the VPC to which - \*\*the knowledge base belongs.
- **database_id** (Optional) - The unique identifier of the DigitalOcean - \*\*OpenSearch database this knowledge base will use.
- **is_public** (Optional) - Indicates whether the knowledge base is public or - \*\*private.
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

A DigitalOcean Gradient AI Knowledge Base can be imported using its UUID. For example:

```sh
terraform import digitalocean_gradientai_knowledge_base.example a1b2c3d4-5678-90ab-cdef-1234567890ab
```

## Usage Notes

- Changes to **datasources**, **embedding_model_uuid**, **spaces_data_source**, **web_crawler_data_source**, **agent_uuid** and **vpc_uuid** will force recreation of the knowledge base.
- To add additional data sources after creation, use the `digitalocean_gradientai_knowledge_base_data_source` resource.
- To attach a knowledge base to an agent, use the `digitalocean_gradientai_agent_knowledge_base_attachment` resource.

# digitalocean_gradientai_openai_api_key

Provides a resource to manage a DigitalOcean Gradient AI OpenAI API Key. With this resource you can create, update, and delete OpenAI API keys, as well as reference them in other Gradient AI resources (such as agents).

## Example Usage

```hcl
resource "digitalocean_gradientai_openai_api_key" "example" {
  api_key = "sk-proj-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  name    = "Production Key"
}

data "digitalocean_gradientai_openai_api_keys" "all" {}

output "all_openai_api_keys" {
  value = data.digitalocean_gradientai_openai_api_keys.all.openai_api_keys
}

data "digitalocean_gradientai_openai_api_key" "by_id" {
  uuid = digitalocean_gradientai_openai_api_key.example.uuid
}

output "openai_api_key_info" {
  value = data.digitalocean_gradientai_openai_api_key.by_id
}
```

## Argument Reference

The following arguments are supported:

- **api_key** (Required) - The OpenAI API key string.
- **name** (Required) - The name assigned to the API key.

## Attributes Reference

After creation, the following attributes are exported:

- **id** - The unique identifier of the OpenAI API key (same as `uuid`).
- **uuid** - The UUID of the OpenAI API key.
- **name** - The name of the API key.
- **created_at** - The timestamp when the API key was created.
- **updated_at** - The timestamp when the API key was last updated.
- **deleted_at** - The timestamp when the API key was deleted (if applicable).
- **created_by** - The user who created the API key.
- **models** - The list of models associated with the API key.

### List All OpenAI API Keys

```hcl
data "digitalocean_gradientai_openai_api_keys" "all" {}

output "all_openai_api_keys" {
  value = data.digitalocean_gradientai_openai_api_keys.all.openai_api_keys
}
```

### Get OpenAI API Key by UUID

```hcl
data "digitalocean_gradientai_openai_api_key" "by_id" {
  uuid = "your-openai-api-key-uuid"
}

output "openai_api_key_info" {
  value = data.digitalocean_gradientai_openai_api_key.by_id
}
```

### List Agents by OpenAI API Key

```hcl
data "digitalocean_gradientai_agents_by_openai_api_key" "by_key" {
  uuid = digitalocean_gradientai_openai_api_key.example.uuid
}

output "agents_by_openai_key" {
  value = data.digitalocean_gradientai_agents_by_openai_api_key.by_key.agents
}
```

## Update Behavior

When the **name** attribute is changed, the provider invokes the update API endpoint to adjust the OpenAI API key's configuration.

## Import

A DigitalOcean Gradient AI OpenAI API Key can be imported using its UUID. For example:

```sh
terraform import digitalocean_gradientai_openai_api_key.example a1b2c3d4-5678-90ab-cdef-1234567890ab
```

## Usage Notes

- The OpenAI API key resource can be referenced by agents and other Gradient AI resources.
- Deleting the API key resource in Terraform will remove it from your DigitalOcean account.

# digitalocean_gradientai_agent_route

Provides a resource to manage a DigitalOcean Gradient AI Agent Route. With this resource you can create, update, and delete agent routes to connect parent agents with child agents for routing functionality.

## Example Usage

```hcl

resource "digitalocean_gradientai_agent_route" "weather_route" {
  parent_agent_uuid = "b90e05b8-566f-11f0-bf8f-4e013e2ddde4"
  child_agent_uuid  = "01efac06-500e-11f0-bf8f-4e013e2ddde4"
  route_name        = "weather_route"
  if_case           = "use this to get weather information"
}
```

## Argument Reference

The following arguments are supported:

- **parent_agent_uuid** (Required) - The UUID of the parent agent that will route requests.
- **child_agent_uuid** (Required) - The UUID of the child agent that will handle routed requests.
- **route_name** (Optional) - The name assigned to the route for identification.
- **if_case** (Optional) - The condition or case description for when this route should be used.

## Attributes Reference

After creation, the following attributes are exported:

- **uuid** - The unique identifier of the agent route.
- **parent_agent_uuid** - The UUID of the parent agent.
- **child_agent_uuid** - The UUID of the child agent.
- **route_name** - The name of the route.
- **if_case** - The condition for using this route.

## Update Behavior

When the **route_name** or **if_case** attributes are changed, the provider invokes the update API endpoint to adjust the route's configuration. The **parent_agent_uuid** and **child_agent_uuid** cannot be changed after creation.

## Import

A DigitalOcean Gradient AI Agent Route can be imported using its UUID. For example:

```sh
terraform import digitalocean_gradientai_agent_route.weather_route 12345678-1234-1234-1234-123456789012
```

## Usage Notes

- Agent routes enable hierarchical agent structures where parent agents can route requests to appropriate child agents based on conditions.
- Both parent and child agents must exist before creating a route between them.
- Changes to **parent_agent_uuid** or **child_agent_uuid** will force recreation of the route.

---

# digitalocean_gradientai_indexing_job_cancel

Provides a resource to cancel running or pending indexing jobs for DigitalOcean Gradient AI Knowledge Bases. This resource is useful for managing long-running indexing operations that need to be stopped before completion.

## Example Usage

```hcl
# Cancel a specific indexing job
resource "digitalocean_gradientai_indexing_job_cancel" "cancel_job" {
  uuid = "f1e2d3c4-5678-90ab-cdef-1234567890ab"
}

# Cancel a job conditionally based on its status
data "digitalocean_gradientai_indexing_job" "monitor_job" {
  uuid = "f1e2d3c4-5678-90ab-cdef-1234567890ab"
}

resource "digitalocean_gradientai_indexing_job_cancel" "conditional_cancel" {
  count = data.digitalocean_gradientai_indexing_job.monitor_job.status == "running" && data.digitalocean_gradientai_indexing_job.monitor_job.phase == "processing" ? 1 : 0

  uuid = data.digitalocean_gradientai_indexing_job.monitor_job.uuid
}
```

## Argument Reference

The following arguments are supported:

- **uuid** (Required) - The unique identifier of the indexing job to cancel.

## Attributes Reference

After creation, the following attributes are exported:

- **id** - The unique identifier of the indexing job (same as uuid).
- **uuid** - The UUID of the indexing job.
- **status** - The status of the indexing job after cancellation.
- **knowledge_base_uuid** - The UUID of the knowledge base associated with this indexing job.
- **phase** - Current phase of the indexing job.
- **completed_datasources** - Number of data sources that were completed before cancellation.
- **total_datasources** - Total number of data sources in the indexing job.
- **tokens** - Number of tokens processed before cancellation.
- **total_items_failed** - Total number of items that failed during indexing.
- **total_items_indexed** - Total number of items that were successfully indexed.
- **total_items_skipped** - Total number of items that were skipped during indexing.
- **data_source_uuids** - List of data source UUIDs associated with this indexing job.
- **created_at** - When the indexing job was created (in RFC3339 format).
- **updated_at** - When the indexing job was last updated (in RFC3339 format).
- **started_at** - When the indexing job was started (in RFC3339 format).
- **finished_at** - When the indexing job was finished (in RFC3339 format).

## Behavior Notes

- **Immediate Effect**: Once this resource is created, the cancellation request is sent immediately to the DigitalOcean API.
- **Status Requirements**: Only indexing jobs with status "pending" or "running" can be cancelled. Attempting to cancel completed, failed, or already cancelled jobs will result in an error.

## Lifecycle Behavior

- **Creation**: Sends cancellation request to the specified indexing job.
- **Update**: This resource is immutable - changes to `uuid` will force recreation.
- **Deletion**: Removing this resource from Terraform configuration does not restore or restart the cancelled indexing job.

## Error Handling

The resource will fail with an error in the following scenarios:

- The indexing job UUID does not exist
- The indexing job is not in a cancellable state ("pending" or "running")
- Insufficient permissions to cancel the indexing job
- The associated knowledge base has been deleted

## Import

A DigitalOcean Gradient AI Indexing Job Cancel operation cannot be imported as it represents a one-time action rather than a persistent resource state.

## Usage Notes

- **One-time Operation**: This resource represents a cancellation action. Once the job is cancelled, the resource serves as a record of the cancellation.
- **Monitoring**: Use with data sources like `digitalocean_gradientai_indexing_job` to monitor job status before and after cancellation.
- **Cleanup**: Consider using lifecycle rules or conditional logic to only cancel jobs when specific conditions are met.
- **Auditing**: The `reason` field is useful for maintaining an audit trail of why indexing jobs were cancelled.

## Example: Conditional Cancellation Based on Status

```hcl
# Monitor job status and cancel if in certain state
data "digitalocean_gradientai_indexing_job" "long_running" {
  uuid = var.indexing_job_uuid
}

# Cancel if job is running and in processing phase
resource "digitalocean_gradientai_indexing_job_cancel" "cancel_processing" {
  count = (
    data.digitalocean_gradientai_indexing_job.long_running.status == "running" &&
    data.digitalocean_gradientai_indexing_job.long_running.phase == "processing"
  ) ? 1 : 0

  uuid = data.digitalocean_gradientai_indexing_job.long_running.uuid
}
```

---

# digitalocean_gradientai_custom_model

Provides a DigitalOcean Gradient AI Custom Model resource. This can be used to import a model from a supported source (HuggingFace, Spaces, etc.) so that it can be deployed on dedicated inference endpoints.

The resource blocks until the underlying import completes (status becomes `STATUS_READY`).

## Example Usage

### Import a public HuggingFace model

```hcl
resource "digitalocean_gradientai_custom_model" "example" {
  name                        = "qwen3-8b"
  description                 = "Qwen3 8B chat model"
  source_type                 = "SOURCE_TYPE_HUGGINGFACE"
  preferred_gpu_region        = "tor1"
  accept_terms_and_conditions = true
  tags                        = ["chat", "qwen"]

  source_ref {
    repo_id     = "Qwen/Qwen3-8B"
    access_type = "ACCESS_TYPE_PUBLIC"
  }
}
```

### Import a gated HuggingFace model

```hcl
resource "digitalocean_gradientai_custom_model" "gated" {
  name                        = "llama-3-8b"
  source_type                 = "SOURCE_TYPE_HUGGINGFACE"
  preferred_gpu_region        = "tor1"
  accept_terms_and_conditions = true

  source_ref {
    repo_id     = "meta-llama/Meta-Llama-3-8B"
    access_type = "ACCESS_TYPE_GATED"
    hf_token    = var.hugging_face_token
  }
}
```

### Import from a Spaces bucket

```hcl
resource "digitalocean_gradientai_custom_model" "spaces" {
  name                        = "my-finetune"
  source_type                 = "SOURCE_TYPE_SPACES_BUCKET"
  accept_terms_and_conditions = true

  source_ref {
    bucket = "my-model-artifacts"
    region = "nyc3"
    prefix = "models/my-finetune/"
  }
}
```

### Update description and tags

Once a model has been imported, its `description` and `tags` can be updated in place for any source type. Re-applying the configuration below against an existing model will issue a metadata-update API call without re-importing the artifacts.

```hcl
resource "digitalocean_gradientai_custom_model" "example" {
  name                        = "qwen3-8b"
  description                 = "Production Qwen3 8B chat model"
  source_type                 = "SOURCE_TYPE_HUGGINGFACE"
  preferred_gpu_region        = "tor1"
  accept_terms_and_conditions = true
  tags                        = ["chat", "qwen", "production"]

  source_ref {
    repo_id     = "Qwen/Qwen3-8B"
    access_type = "ACCESS_TYPE_PUBLIC"
  }
}
```

### Override Spaces-only metadata fields

For `SOURCE_TYPE_SPACES_BUCKET` imports, the API additionally accepts overrides for `license`, `parameters`, `input_modalities`, and `output_modalities`. Setting any of these attributes on a different source type (e.g. `SOURCE_TYPE_HUGGINGFACE`) is rejected by the provider at plan time, to prevent a permanent plan diff against API-reported values that the API would silently refuse to update.

```hcl
resource "digitalocean_gradientai_custom_model" "spaces" {
  name                        = "my-finetune"
  description                 = "Production fine-tune"
  source_type                 = "SOURCE_TYPE_SPACES_BUCKET"
  accept_terms_and_conditions = true
  tags                        = ["chat", "production"]

  license           = "apache-2.0"
  parameters        = "8B"
  input_modalities  = ["text"]
  output_modalities = ["text"]

  source_ref {
    bucket = "my-model-artifacts"
    region = "nyc3"
    prefix = "models/my-finetune/"
  }
}
```

## Argument Reference

The following arguments are supported:

- **name** (Required) - A human-readable name for the custom model.
- **source_type** (Required) - The source the model is being imported from. One of `SOURCE_TYPE_HUGGINGFACE`, `SOURCE_TYPE_SPACES_BUCKET`, `SOURCE_TYPE_SDK_UPLOAD`, `SOURCE_TYPE_FINE_TUNING`. The underlying API does not support changing source fields after import; changes in HCL will not take effect on the existing model — destroy and recreate the resource to use a different source.
- **source_ref** (Required) - A single nested block describing the source. The underlying API does not support changing any source field after import. Supports:
  - **repo_id** (Optional) - Repository identifier (required for `SOURCE_TYPE_HUGGINGFACE`).
  - **commit_sha** (Optional, Computed) - Commit SHA to pin for the import. If omitted, the API resolves the SHA and writes it back to state.
  - **access_type** (Optional) - One of `ACCESS_TYPE_PUBLIC`, `ACCESS_TYPE_PRIVATE`, `ACCESS_TYPE_GATED`.
  - **bucket** (Optional) - Spaces bucket name for `SOURCE_TYPE_SPACES_BUCKET` sources.
  - **region** (Optional) - Region of the source bucket.
  - **prefix** (Optional) - Key prefix inside the source bucket.
  - **hf_token** (Optional, Sensitive) - HuggingFace token used to access `ACCESS_TYPE_PRIVATE` or `ACCESS_TYPE_GATED` repositories. Write-only — the API never echoes it back and Terraform retains the value only in state.
- **description** (Optional) - Description of the custom model.
- **tags** (Optional) - Set of user-defined tags associated with the custom model.
- **license** (Optional, Computed) - License of the model. Defaults to the value reported by the importer. Caller-supplied overrides are accepted only for `SOURCE_TYPE_SPACES_BUCKET` imports; setting this attribute on any other source type is rejected at plan time.
- **parameters** (Optional, Computed) - Parameter-count summary for the model (e.g. `8B`). Defaults to the value reported by the importer. Caller-supplied overrides are accepted only for `SOURCE_TYPE_SPACES_BUCKET` imports; setting this attribute on any other source type is rejected at plan time.
- **input_modalities** (Optional, Computed) - List of input modalities supported by the model (e.g. `["text"]`). Defaults to the values reported by the importer. Caller-supplied overrides are accepted only for `SOURCE_TYPE_SPACES_BUCKET` imports; setting this attribute on any other source type is rejected at plan time.
- **output_modalities** (Optional, Computed) - List of output modalities produced by the model. Defaults to the values reported by the importer. Caller-supplied overrides are accepted only for `SOURCE_TYPE_SPACES_BUCKET` imports; setting this attribute on any other source type is rejected at plan time.
- **preferred_gpu_region** (Optional) - Preferred GPU region where the model artifacts should be staged. Set at import time only.
- **accept_terms_and_conditions** (Optional) - Whether the caller accepts the model provider's terms and conditions. Defaults to `false`. Write-only and consulted only at import time.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

- **id** / **uuid** - The UUID of the custom model.
- **status** - The current status of the custom model.
- **error_message** - Error message if the custom model import failed.
- **architecture** - Model architecture reported by the importer.
- **total_size_bytes** - Total size of the imported model artifacts in bytes (string-encoded int64).
- **file_count** - Number of files that make up the imported model.
- **context_length** - Maximum context length supported by the model.
- **cost_estimate_per_month** - Estimated monthly cost of running the model.
- **team_id** - ID of the team that owns the custom model.
- **storage_region** - Region where the custom model artifacts are stored.
- **created_at** - Timestamp when the custom model was created.
- **updated_at** - Timestamp when the custom model was last updated.
- **active_deployments** - A list of dedicated inference deployments referencing this model. Each entry exports `id`, `name`, `region_slug`, `state`, `created_at`, `updated_at`, and an `endpoints` block with `public_endpoint_fqdn` and `private_endpoint_fqdn`.

## Update Behavior

The metadata-update API accepts changes to `description` and `tags` for any custom model, and additionally accepts changes to `license`, `parameters`, `input_modalities`, and `output_modalities` for `SOURCE_TYPE_SPACES_BUCKET` imports only. All other arguments — `name`, `source_type`, every `source_ref` field, `preferred_gpu_region`, and `accept_terms_and_conditions` — are fixed at import time and ignored by subsequent applies.

- Editing **description** or **tags** triggers an in-place metadata update for any source type; the model is not re-imported.
- Editing **license**, **parameters**, **input_modalities**, or **output_modalities** is supported only when the model was imported with `source_type = "SOURCE_TYPE_SPACES_BUCKET"`. The provider rejects plans that set any of these attributes for other source types, since the API would silently ignore the PATCH and the resulting plan diff would never converge.
- Editing any other argument is a no-op against the API. Terraform will keep reconciling state against the values that were originally imported. To change a fixed field, destroy and recreate the resource.
- **commit_sha** is `Computed`: leave it unset to let the API resolve and record the SHA, or pin it explicitly at import time.

## Delete Behavior

`terraform destroy` issues a synchronous delete API call; the call returns only after the model has been removed. After a successful destroy:

- The resource is removed from Terraform state.
- Subsequent `GET` calls against the same UUID return `404`.
- Any dedicated inference deployments that referenced the model continue to exist independently and must be cleaned up separately.

## Importing Existing Models into Terraform

Custom models that already exist in your DigitalOcean account can be brought under Terraform management with `terraform import`:

```
terraform import digitalocean_gradientai_custom_model.example 1c2dca47-c1f3-4d6b-9d39-04a9b0c81f7a
```

Because the API never echoes write-only fields back, `hf_token` and `accept_terms_and_conditions` will not be populated by import; add them to your configuration to match the originally imported values.

## Timeouts

`digitalocean_gradientai_custom_model` provides the following [Timeouts](https://www.terraform.io/docs/language/resources/syntax.html#operation-timeouts) configuration options:

- **create** - (Default `60 minutes`) Used while waiting for the model import to reach `STATUS_READY`.
- **update** - (Default `5 minutes`) Used for metadata updates.
- **delete** - (Default `5 minutes`) Used for the delete API call.
