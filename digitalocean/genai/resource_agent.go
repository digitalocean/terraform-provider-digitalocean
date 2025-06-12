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
		// Importer: &schema.ResourceImporter{
		// 	StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
		// 		// Validate UUID format
		// 		if !isValidUUID(d.Id()) {
		// 			return nil, fmt.Errorf("invalid agent UUID format: %s", d.Id())
		// 		}
		// 		return []*schema.ResourceData{d}, nil
		// 	},
		// },

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
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"chatbot_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Chatbot ID",
						},
					},
				},
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
				Elem:        DeploymentSchema(),
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the Agent was updated",
			},
			"functions": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of Functions",
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
				Description: "OpenAI API Key information",
				Elem:        OpenAiApiKeySchema(),
			},
			"provide_citations": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates if the agent should provide citations in responses",
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
				Computed:    true,
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

	d.Set("name", agent.Name)
	d.Set("instruction", agent.Instruction)
	d.Set("region", agent.Region)
	d.Set("project_id", agent.ProjectId)
	d.Set("description", agent.Description)
	d.Set("created_at", agent.CreatedAt.UTC().String())
	d.Set("updated_at", agent.UpdatedAt.UTC().String())
	d.Set("if_case", agent.IfCase)
	d.Set("max_tokens", agent.MaxTokens)
	d.Set("retrieval_method", agent.RetrievalMethod)
	d.Set("route_created_by", agent.RouteCreatedBy)
	d.Set("route_created_at", agent.RouteCreatedAt.UTC().String())
	d.Set("route_uuid", agent.RouteUuid)
	d.Set("route_name", agent.RouteName)
	d.Set("tags", agent.Tags)
	d.Set("temperature", agent.Temperature)
	d.Set("top_p", agent.TopP)
	d.Set("url", agent.Url)
	d.Set("user_id", agent.UserId)
	d.Set("model_uuid", agent.Model.Uuid)

	if err := d.Set("api_keys", flattenApiKeys(agent.ApiKeys)); err != nil {
		return diag.FromErr(err)
	}
	// if err := d.Set("chatbot_identifiers", flattenChatbotIdentifiers(agent.ChatbotIdentifiers)); err != nil {
	// 	return diag.FromErr(err)
	// }
	if err := d.Set("model", flattenModel([]*godo.Model{agent.Model})); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("anthropic_api_key", flattenAnthropicApiKey(agent.AnthropicApiKey)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("api_key_infos", flattenApiKeyInfos(agent.ApiKeyInfos)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("deployment", flattenDeployment(agent.Deployment)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("functions", flattenFunctions(agent.Functions)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("agent_guardrail", flattenAgentGuardrail(agent.Guardrails)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("chatbot", flattenChatbot(agent.ChatBot)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("knowledge_bases", flattenKnowledgeBases(agent.KnowledgeBases)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("open_ai_api_key", flattenOpenAiApiKey(agent.OpenAiApiKey)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("template", flattenTemplate(agent.Template)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("child_agents", flattenChildAgents(agent.ChildAgents)); err != nil {
		return diag.FromErr(err)
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

//
///
//
//
//
//
//
//	// fmt.Print(agent.KnowledgeBases)
// if agent.KnowledgeBases != nil {
// 	if err := d.Set("knowledge_bases", flattenKnowledgeBases(agent.KnowledgeBases)); err != nil {
// 		return diag.FromErr(err)
// 	}
// }
// modelSlice := []*godo.Model{agent.Model}
// 	if err := d.Set("model", flattenModel(modelSlice)); err != nil {
// 		return diag.FromErr(err)
// 	}
// }
//
//
//
//
//
//

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
