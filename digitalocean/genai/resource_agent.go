package genai

import (
	"context"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ResourceDigitalOceanAgent defines the DigitalOcean Agent resource.
func ResourceDigitalOceanAgent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanAgentCreate,
		ReadContext:   resourceDigitalOceanAgentRead,
		UpdateContext: resourceDigitalOceanAgentUpdate,
		DeleteContext: resourceDigitalOceanAgentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
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
				Elem:        AgentSchema(),
			},
			"child_agents": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "List of child agents",
				Elem:        AgentSchema(),
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
			"updated_at": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Timestamp when the Agent was updated",
			},
			"functions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of Functions",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
					},
				},
			},
			"agent_guardrail": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "AgentGuardrail represents a Guardrail attached to Gen AI Agent",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
				Computed:    true,
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
					},
				},
			},
			"provide_citations": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the agent should provide citations in responses",
			},
			"retrieval_method": {
				//it can only be RETRIEVAL_METHOD_UNKNOWN,RETRIEVAL_METHOD_REWRITE, RETRIEVAL_METHOD_STEP_BACK,RETRIEVAL_METHOD_SUB_QUERIES,RETRIEVAL_METHOD_NONE
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
		},
	}
}

func resourceDigitalOceanAgentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	agentRequest := &godo.AgentCreateRequest{
		AnthropicKeyUuid:  d.Get("anthropic_key_uuid").(string),
		Description:       d.Get("description").(string),
		Instruction:       d.Get("instruction").(string),
		KnowledgeBaseUuid: convertToStringSlice(d.Get("knowledge_base_uuid")),
		ModelUuid:         d.Get("model_uuid").(string),
		Name:              d.Get("name").(string),
		OpenAiKeyUuid:     d.Get("open_ai_key_uuid").(string),
		ProjectId:         d.Get("project_id").(string),
		Region:            d.Get("region").(string),
		Tags:              convertToStringSlice(d.Get("tags")),
	}

	agent, _, err := client.GenAI.CreateAgent(ctx, agentRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(agent.Uuid)
	return resourceDigitalOceanAgentRead(ctx, d, meta)
}

func resourceDigitalOceanAgentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	agent, _, err := client.GenAI.GetAgent(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", agent.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("instruction", agent.Instruction); err != nil {
		return diag.FromErr(err)
	}
	if agent.Model != nil {
		if err := d.Set("model_uuid", agent.Model.Uuid); err != nil {
			return diag.FromErr(err)
		}
		modelSlice := []*godo.Model{agent.Model}
		if err := d.Set("model", flattenModel(modelSlice)); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("region", agent.Region); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("project_id", agent.ProjectId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("description", agent.Description); err != nil {
		return diag.FromErr(err)
	}

	// if err := d.Set("chatbot_identifiers", agent.ChatbotIdentifiers); err != nil {
	// 	return diag.FromErr(err)
	// }
	if err := d.Set("created_at", agent.CreatedAt.UTC().String()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("updated_at", agent.UpdatedAt.UTC().String()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("if_case", agent.IfCase); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("k", agent.K); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("max_tokens", agent.MaxTokens); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("retrieval_method", agent.RetrievalMethod); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("route_created_by", agent.RouteCreatedBy); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("route_created_at", agent.RouteCreatedAt.UTC().String()); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("route_uuid", agent.RouteUuid); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("route_name", agent.RouteName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tags", agent.Tags); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("temperature", agent.Temperature); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("top_p", agent.TopP); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("url", agent.Url); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("user_id", agent.UserId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("api_keys", flattenApiKeys(agent.ApiKeys)); err != nil {
		return diag.FromErr(err)
	}

	if agent.AnthropicApiKey != nil {
		if err := d.Set("anthropic_api_key", flattenAnthropicApiKey(agent.AnthropicApiKey)); err != nil {
			return diag.FromErr(err)
		}
	}

	if agent.ApiKeyInfos != nil {
		if err := d.Set("api_key_infos", flattenApiKeyInfos(agent.ApiKeyInfos)); err != nil {
			return diag.FromErr(err)
		}
	}

	if agent.Deployment != nil {
		if err := d.Set("deployment", flattenDeployment(agent.Deployment)); err != nil {
			return diag.FromErr(err)
		}
	}

	if agent.Functions != nil {
		if err := d.Set("functions", flattenFunctions(agent.Functions)); err != nil {
			return diag.FromErr(err)
		}
	}

	if agent.Guardrails != nil {
		if err := d.Set("agent_guardrail", flattenAgentGuardrail(agent.Guardrails)); err != nil {
			return diag.FromErr(err)
		}
	}

	if agent.ChatBot != nil {
		if err := d.Set("chatbot", flattenChatbot(agent.ChatBot)); err != nil {
			return diag.FromErr(err)
		}
	}

	// fmt.Print(agent.KnowledgeBases)
	// if agent.KnowledgeBases != nil {
	// 	if err := d.Set("knowledge_bases", flattenKnowledgeBases(agent.KnowledgeBases)); err != nil {
	// 		return diag.FromErr(err)
	// 	}
	// }
	d.Set("knowledge_bases", flattenKnowledgeBases(agent.KnowledgeBases))

	if agent.OpenAiApiKey != nil {
		if err := d.Set("open_ai_api_key", flattenOpenAiApiKey(agent.OpenAiApiKey)); err != nil {
			return diag.FromErr(err)
		}
	}

	if agent.Template != nil {
		if err := d.Set("template", flattenTemplate(agent.Template)); err != nil {
			return diag.FromErr(err)
		}
	}
	if agent.ChildAgents != nil {
		if err := d.Set("child_agents", flattenChildAgents(agent.ChildAgents)); err != nil {
			return diag.FromErr(err)
		}
	}
	d.SetId(agent.Uuid)

	return nil
}

func resourceDigitalOceanAgentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	// If visibility has changed, update it separately.
	if d.HasChange("deployment") {
		old, new := d.GetChange("deployment")

		if len(new.([]interface{})) > 0 {
			oldDeployment, newDeployment := extractDeploymentVisibility(old, new)

			if oldDeployment != newDeployment {
				diags := resourceDigitalOceanAgentUpdateVisibility(ctx, d, meta)
				if diags.HasError() {
					return diags
				}
			}
		}
	}

	agentRequest := &godo.AgentUpdateRequest{}

	if d.HasChange("anthropic_key_uuid") {
		agentRequest.AnthropicKeyUuid = d.Get("anthropic_key_uuid").(string)
	}
	if d.HasChange("description") {
		agentRequest.Description = d.Get("description").(string)
	}
	if d.HasChange("instruction") {
		agentRequest.Instruction = d.Get("instruction").(string)
	}
	if d.HasChange("model_uuid") {
		agentRequest.ModelUuid = d.Get("model_uuid").(string)
	}
	if d.HasChange("name") {
		agentRequest.Name = d.Get("name").(string)
	}
	if d.HasChange("project_id") {
		agentRequest.ProjectId = d.Get("project_id").(string)
	}
	if d.HasChange("region") {
		agentRequest.Region = d.Get("region").(string)
	}
	if d.HasChange("k") {
		agentRequest.K = d.Get("k").(int)
	}
	if d.HasChange("max_tokens") {
		agentRequest.MaxTokens = d.Get("max_tokens").(int)
	}
	if d.HasChange("open_ai_key_uuid") {
		agentRequest.OpenAiKeyUuid = d.Get("open_ai_key_uuid").(string)
	}

	if d.HasChange("retrieval_method") {
		agentRequest.RetrievalMethod = d.Get("retrieval_method").(string)
	}
	if d.HasChange("tags") {
		agentRequest.Tags = convertToStringSlice(d.Get("tags"))
	}
	if d.HasChange("temperature") {
		agentRequest.Temperature = d.Get("temperature").(float64)
	}
	if d.HasChange("top_p") {
		agentRequest.TopP = d.Get("top_p").(float64)
	}

	agent, _, err := client.GenAI.UpdateAgent(ctx, d.Id(), agentRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(agent.Uuid)
	return resourceDigitalOceanAgentRead(ctx, d, meta)
}

func resourceDigitalOceanAgentUpdateVisibility(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	deployments := d.Get("deployment").([]interface{})
	if len(deployments) == 0 {
		return diag.Errorf("deployment block is empty")
	}

	deployment := deployments[0].(map[string]interface{})

	visibility, ok := deployment["visibility"].(string)
	if !ok {
		return diag.Errorf("visibility is not a string or is missing")
	}
	updateReq := &godo.AgentVisibilityUpdateRequest{
		Uuid:       d.Id(),
		Visibility: visibility,
	}

	agent, _, err := client.GenAI.UpdateAgentVisibility(ctx, d.Id(), updateReq)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(agent.Uuid)
	return resourceDigitalOceanAgentRead(ctx, d, meta)
}

func resourceDigitalOceanAgentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	_, _, err := client.GenAI.DeleteAgent(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func extractDeploymentVisibility(old, new interface{}) (string, string) {
	oldVisibility := ""
	newVisibility := ""

	if len(old.([]interface{})) > 0 {
		oldDeployment := old.([]interface{})[0].(map[string]interface{})
		oldVisibility = oldDeployment["visibility"].(string)
	}

	if len(new.([]interface{})) > 0 {
		newDeployment := new.([]interface{})[0].(map[string]interface{})
		newVisibility = newDeployment["visibility"].(string)
	}

	return oldVisibility, newVisibility
}

func convertToStringSlice(val interface{}) []string {
	if val == nil {
		return []string{}
	}

	interfaceSlice, ok := val.([]interface{})
	if !ok {
		return []string{}
	}

	var result []string
	for _, elem := range interfaceSlice {
		result = append(result, elem.(string))
	}
	return result
}

func flattenChildAgents(childAgents []*godo.Agent) []interface{} {
	if childAgents == nil {
		return nil
	}
	result := make([]interface{}, 0, len(childAgents))
	for _, child := range childAgents {
		// Build a map with only the fields you want to expose.
		m := map[string]interface{}{
			"agent_id":    child.Uuid,
			"name":        child.Name,
			"region":      child.Region,
			"project_id":  child.ProjectId,
			"description": child.Description,
		}
		result = append(result, m)
	}
	return result
}

func flattenAnthropicApiKey(apiKey *godo.AnthropicApiKeyInfo) []interface{} {
	if apiKey == nil {
		return nil
	}

	m := map[string]interface{}{
		"created_at": apiKey.CreatedAt.UTC().String(),
		"created_by": apiKey.CreatedBy,
		"deleted_at": apiKey.DeletedAt,
		"name":       apiKey.Name,
		"updated_at": apiKey.UpdatedAt.UTC().String(),
		"uuid":       apiKey.Uuid,
	}

	return []interface{}{m}
}

func flattenApiKeyInfos(apiKeyInfos []*godo.ApiKeyInfo) []interface{} {
	if apiKeyInfos == nil {
		return nil
	}

	result := make([]interface{}, 0, len(apiKeyInfos))
	for _, info := range apiKeyInfos {
		m := map[string]interface{}{
			"created_at": info.CreatedAt.UTC().String(),
			"created_by": info.CreatedBy,
			"deleted_at": info.DeletedAt.UTC().String(),
			"name":       info.Name,
			"secret_key": info.SecretKey,
			"uuid":       info.Uuid,
		}
		result = append(result, m)
	}

	return result
}

func flattenDeployment(deployment *godo.AgentDeployment) []interface{} {
	if deployment == nil {
		return nil
	}

	m := map[string]interface{}{
		"created_at": deployment.CreatedAt.UTC().String(),
		"name":       deployment.Name,
		"status":     deployment.Status,
		"updated_at": deployment.UpdatedAt.UTC().String(),
		"url":        deployment.Url,
		"uuid":       deployment.Uuid,
		"visibility": deployment.Visibility,
	}
	return []interface{}{m}
}

func flattenFunctions(functions []*godo.AgentFunction) []interface{} {
	if functions == nil {
		return nil
	}

	result := make([]interface{}, 0, len(functions))
	for _, fn := range functions {
		m := map[string]interface{}{
			"api_key":        fn.ApiKey,
			"created_at":     fn.CreatedAt.UTC().String(),
			"description":    fn.Description,
			"guardrail_uuid": fn.GuardrailUuid,
			"faasname":       fn.FaasName,
			"faasnamespace":  fn.FaasNamespace,
			"name":           fn.Name,
			"updated_at":     fn.UpdatedAt.UTC().String(),
			"url":            fn.Url,
			"uuid":           fn.Uuid,
		}
		result = append(result, m)
	}

	return result
}

func flattenAgentGuardrail(guardrails []*godo.AgentGuardrail) []interface{} {
	if guardrails == nil {
		return nil
	}

	result := make([]interface{}, 0, len(guardrails))
	for _, guardrail := range guardrails {
		m := map[string]interface{}{
			"agent_uuid":       guardrail.AgentUuid,
			"created_at":       guardrail.CreatedAt.UTC().String(),
			"default_response": guardrail.DefaultResponse,
			"description":      guardrail.Description,
			"guardrail_uuid":   guardrail.GuardrailUuid,
			"is_attached":      guardrail.IsAttached,
			"is_default":       guardrail.IsDefault,
			"name":             guardrail.Name,
			"priority":         guardrail.Priority,
			"type":             guardrail.Type,
			"updated_at":       guardrail.UpdatedAt.UTC().String(),
			"uuid":             guardrail.Uuid,
		}
		result = append(result, m)
	}

	return result
}

func flattenChatbot(chatbot *godo.ChatBot) []interface{} {
	if chatbot == nil {
		return nil
	}

	m := map[string]interface{}{
		"button_background_color": chatbot.ButtonBackgroundColor,
		"logo":                    chatbot.Logo,
		"name":                    chatbot.Name,
		"primary_color":           chatbot.PrimaryColor,
		"secondary_color":         chatbot.SecondaryColor,
		"starting_message":        chatbot.StartingMessage,
	}

	return []interface{}{m}
}

func flattenKnowledgeBases(config []*godo.KnowledgeBase) []interface{} {
	if config == nil {
		return []interface{}{}
	}

	result := make([]interface{}, 0, len(config))

	for _, kb := range config {
		k := map[string]interface{}{
			"uuid":                 kb.Uuid,
			"name":                 kb.Name,
			"created_at":           kb.CreatedAt.UTC().String(),
			"updated_at":           kb.UpdatedAt.UTC().String(),
			"tags":                 kb.Tags,
			"region":               kb.Region,
			"embedding_model_uuid": kb.EmbeddingModelUuid,
			"project_id":           kb.ProjectId,
			"database_id":          kb.DatabaseId,
			"added_to_agent_at":    kb.AddedToAgentAt.UTC().String(),
		}

		if kb.LastIndexingJob != nil {
			k["last_indexing_job"] = flattenLastIndexingJob(kb.LastIndexingJob) // Flatten and take the first element
		}

		result = append(result, k)
	}

	return result
}

func flattenModel(models []*godo.Model) []interface{} {
	if models == nil {
		return nil
	}

	result := make([]interface{}, 0, len(models))
	for _, model := range models {
		m := map[string]interface{}{
			"created_at":        model.CreatedAt.UTC().String(),
			"inference_name":    model.InferenceName,
			"inference_version": model.InferenceVersion,
			"is_foundational":   model.IsFoundational,
			"name":              model.Name,
			"parent_uuid":       model.ParentUuid,
			"provider":          model.Provider,
			"updated_at":        model.UpdatedAt.UTC().String(),
			"upload_complete":   model.UploadComplete,
			"url":               model.Url,
			"usecases":          model.Usecases,
		}

		if model.Version != nil {
			versionMap := map[string]interface{}{
				"major": model.Version.Major,
				"minor": model.Version.Minor,
				"patch": model.Version.Patch,
			}
			m["versions"] = []interface{}{versionMap}
		}

		if model.Agreement != nil {
			agreementMap := map[string]interface{}{
				"description": model.Agreement.Description,
				"name":        model.Agreement.Name,
				"url":         model.Agreement.Url,
				"uuid":        model.Agreement.Uuid,
			}
			m["agreement"] = []interface{}{agreementMap}
		}

		result = append(result, m)
	}

	return result
}

func flattenApiKeys(apiKeys []*godo.ApiKey) []interface{} {
	if apiKeys == nil {
		return nil
	}

	result := make([]interface{}, 0, len(apiKeys))
	for _, key := range apiKeys {
		m := map[string]interface{}{
			"api_key": key.ApiKey,
		}
		result = append(result, m)
	}

	return result
}

func flattenOpenAiApiKey(apiKey *godo.OpenAiApiKey) []interface{} {
	if apiKey == nil {
		return nil
	}

	m := map[string]interface{}{
		"created_at": apiKey.CreatedAt.UTC().String(),
		"created_by": apiKey.CreatedBy,
		"deleted_at": apiKey.DeletedAt.UTC().String(),
		"name":       apiKey.Name,
		"updated_at": apiKey.UpdatedAt.UTC().String(),
		"uuid":       apiKey.Uuid,
		"model":      flattenModel(apiKey.Models),
	}

	return []interface{}{m}
}

func flattenTemplate(template *godo.AgentTemplate) []interface{} {
	if template == nil {
		return nil
	}

	m := map[string]interface{}{
		"created_at":  template.CreatedAt.UTC().String(),
		"instruction": template.Instruction,
		"description": template.Description,
		"k":           template.K,
		"max_tokens":  template.MaxTokens,
		"name":        template.Name,
		"temperature": template.Temperature,
		"top_p":       template.TopP,
		"uuid":        template.Uuid,
		"updated_at":  template.UpdatedAt.UTC().String(),
	}

	return []interface{}{m}
}

func flattenLastIndexingJob(job *godo.LastIndexingJob) []interface{} {
	if job == nil {
		return nil
	}

	// Convert datasource_uuids from []string to []interface{} if needed.
	var datasourceUuids []interface{}
	if job.DataSourceUuids != nil {
		datasourceUuids = make([]interface{}, len(job.DataSourceUuids))
		for i, id := range job.DataSourceUuids {
			datasourceUuids[i] = id
		}
	}

	m := map[string]interface{}{
		"completed_datasources": job.CompletedDatasources,
		"created_at":            job.CreatedAt.UTC().String(),
		"datasource_uuids":      datasourceUuids,
		"finished_at":           job.FinishedAt.UTC().String(),
		"knowledge_uuid":        job.KnowledgeBaseUuid,
		"phase":                 job.Phase,
		"started_at":            job.StartedAt.UTC().String(),
		"tokens":                job.Tokens,
		"total_datasources":     job.TotalDatasources,
		"updated_at":            job.UpdatedAt.UTC().String(),
		"uuid":                  job.Uuid,
	}

	return []interface{}{m}
}

// func flattenKnowledgeBases(kbs []*godo.KnowledgeBase) []interface{} {
// 	if kbs == nil {
// 		return nil
// 	}
// 	kb := kbs[0] // Assuming you want to flatten only the first KnowledgeBase for simplicity
// 	// result := make([]interface{}, 0, len(kbs))
// 	// for _, kb := range kbs {
// 	m := map[string]interface{}{
// 		// "added_to_agent_at":    kb.AddedToAgentAt.UTC().String(),
// 		// "created_at":           kb.CreatedAt.UTC().String(),
// 		"database_id":          kb.DatabaseId,
// 		"embedding_model_uuid": kb.EmbeddingModelUuid,
// 		"is_public":            kb.IsPublic,
// 		"name":                 kb.Name,
// 		"project_id":           kb.ProjectId,
// 		"region":               kb.Region,
// 		"user_id":              kb.UserId,
// 		"uuid":                 kb.Uuid,
// 		// "tags":                 kb.Tags,
// 	}

// 	// Flatten tags (assumed to be a []string)
// 	// if kb.Tags != nil {
// 	// 	tags := make([]interface{}, len(kb.Tags))
// 	// 	for i, tag := range kb.Tags {
// 	// 		tags[i] = tag
// 	// 	}
// 	// 	m["tags"] = tags
// 	// }

// 	// Flatten last indexing job as a nested block if present.
// 	// if kb.LastIndexingJob != nil {
// 	// 	m["last_indexing_job"] = flattenLastIndexingJob(kb.LastIndexingJob)
// 	// }

// 	// result = append(result, m)
// 	// }
// 	return []interface{}{m}
// }

// func flattenModel(model *[]godo.Model) []interface{} {
// 	if model == nil {
// 		return nil
// 	}

// 	m := map[string]interface{}{
// 		"created_at":        model.CreatedAt.UTC().String(),
// 		"inference_name":    model.InferenceName,
// 		"inference_version": model.InferenceVersion,
// 		"is_foundational":   model.IsFoundational,
// 		"name":              model.Name,
// 		"parent_uuid":       model.ParentUuid,
// 		"provider":          model.Provider,
// 		"updated_at":        model.UpdatedAt.UTC().String(),
// 		"upload_complete":   model.UploadComplete,
// 		"url":               model.Url,
// 		"usecases":          model.Usecases,
// 		// "versions":          model.Versions, - double nesting
// 	}

// 	if model.Version != nil {
// 		versionMap := map[string]interface{}{
// 			"major": model.Version.Major,
// 			"minor": model.Version.Minor,
// 			"patch": model.Version.Patch,
// 		}
// 		m["versions"] = []interface{}{versionMap}
// 	}

// 	if model.Agreement != nil {
// 		agreementMap := map[string]interface{}{
// 			"description": model.Agreement.Description,
// 			"name":        model.Agreement.Name,
// 			"url":         model.Agreement.Url,
// 			"uuid":        model.Agreement.Uuid,
// 		}
// 		m["agreement"] = []interface{}{agreementMap}
// 	}

// 	return []interface{}{m}
// }

// agentRequest.Uuid = d.Get("id").(string)
// if d.HasChange("uuid") {
// 	agentRequest.Uuid = d.Get("uuid").(string) // Assuming it's a string
// }

// agentRequest := &godo.AgentUpdateRequest{
// 	AnthropicKeyUuid: d.Get("anthropic_key_uuid").(string),
// 	Description:      d.Get("description").(string),
// 	Name:             d.Get("name").(string),
// 	Instruction:      d.Get("instruction").(string),
// 	ModelUuid:        d.Get("model_uuid").(string),
// 	Region:           d.Get("region").(string),
// 	ProjectId:        d.Get("project_id").(string),
// 	K:                d.Get("k").(int),
// 	MaxTokens:        d.Get("max_tokens").(int),
// 	OpenAiKeyUuid:    d.Get("open_ai_key_uuid").(string),
// }
