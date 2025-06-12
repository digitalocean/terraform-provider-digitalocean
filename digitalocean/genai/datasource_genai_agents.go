package genai

import (
	"context"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceDigitalOceanAgentList defines the data source for listing agents.
func DataSourceDigitalOceanAgentList() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanAgentListRead,
		Schema: map[string]*schema.Schema{
			"agents": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of agents",
				Elem: &schema.Resource{
					Schema: AgentSchemaRead().Schema,
				},
			},
		},
	}
}

func dataSourceDigitalOceanAgentListRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	agents, _, err := client.GenAI.ListAgents(ctx, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var results []map[string]interface{}
	for _, agent := range agents {
		result := map[string]interface{}{
			"anthropic_api_key": flattenAnthropicApiKey(agent.AnthropicApiKey),
			"api_key_infos":     flattenApiKeyInfos(agent.ApiKeyInfos),
			"api_keys":          flattenApiKeys(agent.ApiKeys),
			"chatbot":           flattenChatbot(agent.ChatBot),
			// "chatbot_identifiers:": flattenChatbotIdentifiers(agent.ChatbotIdentifiers),
			"created_at":      agent.CreatedAt.UTC().String(),
			"child_agents":    flattenChildAgents(agent.ChildAgents),
			"deployment":      flattenDeployment(agent.Deployment),
			"description":     agent.Description,
			"updated_at":      agent.UpdatedAt.UTC().String(),
			"functions":       flattenFunctions(agent.Functions),
			"agent_guardrail": flattenAgentGuardrail(agent.Guardrails),
			"if_case":         agent.IfCase,
			"instruction":     agent.Instruction,
			"k":               agent.K,
			"knowledge_bases": flattenKnowledgeBases(agent.KnowledgeBases),
			"max_tokens":      agent.MaxTokens,
			"name":            agent.Name,
			"open_ai_api_key": flattenOpenAiApiKey(agent.OpenAiApiKey),
			//ParentAgents
			"project_id":       agent.ProjectId,
			"region":           agent.Region,
			"retrieval_method": agent.RetrievalMethod,
			"route_created_at": agent.RouteCreatedAt.UTC().String(),
			"route_created_by": agent.RouteCreatedBy,
			"route_uuid":       agent.RouteUuid,
			"route_name":       agent.RouteName,
			"tags":             agent.Tags,
			"template":         flattenTemplate(agent.Template),
			"temperature":      agent.Temperature,
			"top_p":            agent.TopP,
			"url":              agent.Url,
			"user_id":          agent.UserId,
			"agent_id":         agent.Uuid,
		}
		results = append(results, result)
	}

	if err := d.Set("agents", results); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("agent-list")
	return nil
}
