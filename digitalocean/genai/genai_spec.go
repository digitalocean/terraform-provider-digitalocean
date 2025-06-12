package genai

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func LastIndexingJobSchema() *schema.Resource {
	lastIndexingSchema := map[string]*schema.Schema{
		"completed_datasources": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Number of completed datasources in the last indexing job",
		},
		"created_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Created At timestamp for the last indexing job",
		},
		"datasource_uuids": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Datasource UUIDs for the last indexing job",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"finished_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Timestamp when the last indexing job finished",
		},
		"knowledge_uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "UUID	of the Knowledge Base for the last indexing job",
		},
		"phase": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Phase of the last indexing job",
		},
		"started_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Timestamp when the last indexing job started",
		},
		"tokens": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Number of tokens processed in the last indexing job",
		},
		"total_datasources": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Total number of datasources in the last indexing job",
		},
		"updated_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Timestamp when the last indexing job updated",
		},
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "UUID	of the last indexing job",
		},
	}
	return &schema.Resource{
		Schema: lastIndexingSchema,
	}
}

func AgreementSchema() *schema.Resource {
	agreementSchema := map[string]*schema.Schema{
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description of the agreement",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of the agreement",
		},
		"url": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "URL of the agreement",
		},
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "UUID of the agreement",
		},
	}
	return &schema.Resource{
		Schema: agreementSchema,
	}
}

func AnthropicApiKeySchema() *schema.Resource {
	anthropicApiKeySchema := map[string]*schema.Schema{
		"created_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Timestamp when the API Key was created",
		},
		"created_by": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Created By user ID for the API Key",
		},
		"deleted_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Deleted At timestamp for the API Key",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of the API Key",
		},
		"updated_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Updated At timestamp for the API Key",
		},
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "API Key value",
		},
	}
	return &schema.Resource{
		Schema: anthropicApiKeySchema,
	}
}

// create a flatten agent function that flattens child,parent,agent
func AgentSchema() *schema.Resource { //map[string]*schema.Schema - didn't work
	agentSchema := map[string]*schema.Schema{
		"anthropic_api_key": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Anthropic API Key information",
			Elem:        AnthropicApiKeySchema(),
		},
		"api_key_infos": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of API Key Infos",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"created_at": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "API Key value",
					},
					"created_by": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Created By user ID for the API Key",
					},
					"deleted_at": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Deleted At timestamp for the API Key",
					},
					"name": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Name of the API Key",
					},
					"secret_key": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Updated At timestamp for the API Key",
					},
					"uuid": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "API Key value",
					},
				},
			},
		},
		"api_keys": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of API Keys",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"api_key": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "API Key value",
					},
				},
			},
		},
		"chatbot": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "ChatBot configuration",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"button_background_color": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Background color for the chatbot button",
					},
					"logo": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Logo for the chatbot",
					},
					"name": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Name of the chatbot",
					},
					"primary_color": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Primary color for the chatbot",
					},
					"secondary_color": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Secondary color for the chatbot",
					},
					"starting_message": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Starting message for the chatbot",
					},
				},
			},
		},
		"chatbot_identifiers": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of Chatbot Identifiers",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"deployment": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of API Key Infos",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"created_at": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "API Key value",
					},
					"name": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Name of the API Key",
					},
					"status": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Status of the Deployment",
					},
					"updated_at": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Updated At timestamp for the Agent",
					},
					"url": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Url of the Deployment",
					},
					"uuid": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "API Key value",
					},
					"visibility": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Visibility of the Deployment",
					},
				},
			},
		},
		"name": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Name of the Agent",
		},
		"instruction": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Instruction for the Agent",
		},
		"model_uuid": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Model UUID of the Agent",
		},
		"project_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Project ID of the Agent",
		},
		"region": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "Region where the Agent is deployed",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description for the Agent",
		},
		"agent_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "ID of the child agent",
		},
	}
	return &schema.Resource{
		Schema: agentSchema,
	}
}

func TemplateSchema() *schema.Resource {
	templateSchem := map[string]*schema.Schema{
		"created_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Created At timestamp for the Knowledge Base",
		},
		"instruction": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Instruction for the Agent",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description of the Agent Template",
		},
		"k": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "K value for the Agent Template",
		},
		"knowledge_bases": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of Knowledge Bases",
			Elem:        KnowledgeBaseSchema(),
		},
		"max_tokens": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Maximum tokens allowed",
		},
		"model": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Model of the Agent Template",
			Elem:        ModelSchema(),
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of the Agent Template",
		},
		"temperature": {
			Type:        schema.TypeFloat,
			Optional:    true,
			Description: "Agent temperature setting",
		},
		"top_p": {
			Type:        schema.TypeFloat,
			Optional:    true,
			Description: "Top P sampling parameter",
		},
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "uuid of the Agent Template",
		},
		"updated_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Updated At timestamp for the Agent Template",
		},
	}
	return &schema.Resource{
		Schema: templateSchem,
	}

}

func ChatbotSchema() *schema.Resource {
	chatbotSchema := map[string]*schema.Schema{
		"button_background_color": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Background color for the chatbot button",
		},
		"logo": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Logo for the chatbot",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of the chatbot",
		},
		"primary_color": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Primary color for the chatbot",
		},
		"secondary_color": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Secondary color for the chatbot",
		},
		"starting_message": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Starting message for the chatbot",
		},
	}

	return &schema.Resource{
		Schema: chatbotSchema,
	}
}

func FunctionsSchema() *schema.Resource {
	functionsSchema := map[string]*schema.Schema{
		"api_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "API Key value",
		},
		"created_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Created At timestamp for the Function",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description of the Function",
		},
		"guardrail_uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Guardrail UUID for the Function",
		},
		"faasname": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of function",
		},
		"faasnamespace": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Namespace of function",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of function",
		},
		"updated_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Updated At timestamp for the Agent",
		},
		"url": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Url of the Deployment",
		},
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "API Key value",
		},
	}

	return &schema.Resource{
		Schema: functionsSchema,
	}
}

func DeploymentSchema() *schema.Resource {
	deploymentSchema := map[string]*schema.Schema{
		"created_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "API Key value",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of the API Key",
		},
		"status": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Status of the Deployment",
		},
		"updated_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Updated At timestamp for the Agent",
		},
		"url": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Url of the Deployment",
		},
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "API Key value",
		},
		"visibility": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Visibility of the Deployment",
		},
	}

	return &schema.Resource{
		Schema: deploymentSchema,
	}

}

func OpenAiApiKeySchema() *schema.Resource {
	openAiApiKeySchema := map[string]*schema.Schema{
		"created_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Timestamp when the API Key was created",
		},
		"created_by": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Created By user ID for the API Key",
		},
		"deleted_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Deleted At timestamp for the API Key",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of the API Key",
		},
		"updated_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Updated At timestamp for the API Key",
		},
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "API Key value",
		},
	}
	return &schema.Resource{
		Schema: openAiApiKeySchema,
	}
}

func ApiKeysSchema() *schema.Resource {
	apiKeysSchema := map[string]*schema.Schema{
		"created_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "API Key value",
		},
		"created_by": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Created By user ID for the API Key",
		},
		"deleted_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Deleted At timestamp for the API Key",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of the API Key",
		},
		"secret_key": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Updated At timestamp for the API Key",
		},
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "API Key value",
		},
	}

	return &schema.Resource{
		Schema: apiKeysSchema,
	}
}

func AgentGuardrailSchema() *schema.Resource {
	agentGuardRailSchema := map[string]*schema.Schema{
		"agent_uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Agent UUID for the Guardrail",
		},
		"created_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Created At timestamp for the Guardrail",
		},
		"default_response": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Default response for the Guardrail",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description of the Guardrail",
		},
		"guardrail_uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Guardrail UUID",
		},
		"is_attached": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Indicates if the Guardrail is attached",
		},
		"is_default": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Indicates if the Guardrail is default",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of Guardrail",
		},
		"priority": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Priority of the Guardrail",
		},
		"type": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Type of the Guardrail",
		},
		"updated_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Updated At timestamp for the Guardrail",
		},
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Guardrail UUID",
		},
	}

	return &schema.Resource{
		Schema: agentGuardRailSchema,
	}
}

func ModelSchema() *schema.Resource {
	modelSchema := map[string]*schema.Schema{
		"agreement": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Agreement information for the model",
			Elem:        AgreementSchema(),
		},
		"created_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Created At timestamp for the Knowledge Base",
		},
		"inference_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Inference name of the model",
		},
		"inference_version": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Infernce version of the model",
		},
		"is_foundational": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Indicates if the Model Base is foundational",
		},
		"name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Name of the Knowledge Base",
		},
		"parent_uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Parent UUID of the Model",
		},
		"provider": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Provider of the Model",
		},
		"updated_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Timestamp when the Knowledge Base was updated",
		},
		"upload_complete": {
			Type:        schema.TypeBool,
			Optional:    true,
			Description: "Indicates if the Model upload is complete",
		},
		"url": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "URL of the Model",
		},
		"usecases": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of Usecases for the Model",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"versions": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "URL of the Model",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"major": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Major version of the model",
					},
					"minor": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Minor version of the model",
					},
					"patch": {
						Type:        schema.TypeInt,
						Optional:    true,
						Description: "Patch version of the model",
					},
				},
			},
		},
	}

	return &schema.Resource{
		Schema: modelSchema,
	}
}

func KnowledgeBaseSchema() *schema.Resource {
	knowledgeBaseSchema := map[string]*schema.Schema{
		"added_to_agent_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Timestamp when the Knowledge Base was added to the Agent",
		},
		"created_at": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Created At timestamp for the Knowledge Base",
		},
		"database_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Database ID of the Knowledge Base",
		},
		"embedding_model_uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Embedding model UUID for the Knowledge Base",
		},
		"is_public": {
			Type:        schema.TypeBool,
			Computed:    true,
			Description: "Indicates if the Knowledge Base is public",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the Knowledge Base",
		},
		"project_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Project ID of the Knowledge Base",
		},
		"region": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Region of the Knowledge Base",
		},
		"user_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "User ID of the Knowledge Base",
		},
		"uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "UUID of the Knowledge Base",
		},
		"tags": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of tags",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"last_indexing_job": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Last indexing job for the Knowledge Base",
			Elem:        LastIndexingJobSchema(),
		},
	}
	return &schema.Resource{
		Schema: knowledgeBaseSchema,
	}
}

func AgentSchemaRead() *schema.Resource {
	agentSchema := map[string]*schema.Schema{
		"agent_id": {
			Type:        schema.TypeString,
			Required:    true,
			Description: "ID of the Agent to retrieve",
		},
		"name": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Name of the Agent",
		},
		"instruction": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Instruction for the Agent",
		},
		"model_uuid": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Model UUID of the Agent",
		},
		"project_id": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Project ID of the Agent",
		},
		"region": {
			Type:        schema.TypeString,
			Computed:    true,
			Description: "Region where the Agent is deployed",
		},
		"description": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Description for the Agent",
		},
		"anthropic_key_uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional Anthropic API key ID to use with Anthropic models",
		},
		"knowledge_base_uuid": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "Ids of the knowledge base(s) to attach to the agent",
			Elem:        &schema.Schema{Type: schema.TypeString}, //it was TypeList which gave error before
		},
		"open_ai_key_uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Optional OpenAI API key ID to use with OpenAI models",
		},
		"anthropic_api_key": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Anthropic API Key information",
			Elem:        AnthropicApiKeySchema(),
		},
		"api_key_infos": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of API Key Infos",
			Elem:        ApiKeysSchema(),
		},
		"api_keys": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of API Keys",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"api_key": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "API Key value",
					},
				},
			},
		},
		"chatbot_identifiers": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of Chatbot Identifiers",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"created_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Timestamp when the Agent was created",
		},
		"parent_agents": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "List of child agents",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Name of the Agent",
					},
					"instruction": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Instruction for the Agent",
					},
					"model_uuid": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Model UUID of the Agent",
					},
					"project_id": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Project ID of the Agent",
					},
					"region": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Region where the Agent is deployed",
					},
					"description": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Description for the Agent",
					},
					"agent_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "ID of the child agent",
					},
				},
			},
		},
		"child_agents": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "List of child agents",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"name": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Name of the Agent",
					},
					"instruction": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Instruction for the Agent",
					},
					"model_uuid": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Model UUID of the Agent",
					},
					"project_id": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Project ID of the Agent",
					},
					"region": {
						Type:        schema.TypeString,
						Required:    true,
						Description: "Region where the Agent is deployed",
					},
					"description": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "Description for the Agent",
					},
					"agent_id": {
						Type:        schema.TypeString,
						Computed:    true,
						Description: "ID of the child agent",
					},
				},
			},
		},
		"deployment": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of API Key Infos",
			Elem:        DeploymentSchema(),
		},
		"updated_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Timestamp when the Agent was updated",
		},
		"functions": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of API Key Infos",
			Elem:        FunctionsSchema(),
		},
		"agent_guardrail": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "AgentGuardrail represents a Guardrail attached to Gen AI Agent",
			Elem:        AgentGuardrailSchema(),
		},
		"chatbot": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "ChatBot configuration",
			Elem:        ChatbotSchema(),
		},
		"if_case": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "If case condition",
		},
		"k": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "K value",
		},
		"knowledge_bases": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of Knowledge Bases",
			Elem:        KnowledgeBaseSchema(),
		},
		"max_tokens": {
			Type:        schema.TypeInt,
			Optional:    true,
			Description: "Maximum tokens allowed",
		},
		"model": {
			Type:        schema.TypeList,
			Optional:    true,
			Computed:    true,
			Description: "Model of the Agent",
			Elem:        ModelSchema(),
		},
		"open_ai_api_key": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "OpenAI API Key information",
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"api_key": {
						Type:        schema.TypeString,
						Optional:    true,
						Description: "OpenAI API Key",
					},
				},
			},
		},
		"retrieval_method": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Retrieval method used",
		},
		"route_created_by": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "User who created the route",
		},
		"route_created_at": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Timestamp when the route was created",
		},
		"route_uuid": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Route UUID",
		},
		"route_name": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "Route name",
		},
		"tags": {
			Type:        schema.TypeList,
			Optional:    true,
			Description: "List of Tags",
			Elem:        &schema.Schema{Type: schema.TypeString},
		},
		"template": {
			Type:        schema.TypeList,
			Optional:    true,
			MaxItems:    1,
			Description: "Agent Template",
			Elem:        TemplateSchema(),
		},
		"temperature": {
			Type:        schema.TypeFloat,
			Optional:    true,
			Description: "Agent temperature setting",
		},
		"top_p": {
			Type:        schema.TypeFloat,
			Optional:    true,
			Description: "Top P sampling parameter",
		},
		"url": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "URL for the Agent",
		},
		"user_id": {
			Type:        schema.TypeString,
			Optional:    true,
			Description: "User ID linked with the Agent",
		},
	}
	return &schema.Resource{
		Schema: agentSchema,
	}
}
