package genai

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
)

func getDigitalOceanAgents(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.CombinedConfig).GodoClient()

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	var allAgents []interface{}
	for {
		agents, resp, err := client.GenAI.ListAgents(context.Background(), opts)
		if err != nil {
			return nil, fmt.Errorf("error retrieving agents : %s", err)
		}

		for _, agent := range agents {
			allAgents = append(allAgents, agent)
		}
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("error retrieving agents: %s", err)
		}

		opts.Page = page + 1

	}
	return allAgents, nil

}

func flattenDigitalOceanAgent(rawDomain, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	agent, ok := rawDomain.(*godo.Agent)
	if !ok {
		return nil, fmt.Errorf("expected *godo.Agent, got %T", rawDomain)
	}

	if agent == nil {
		return nil, fmt.Errorf("agent is nil")
	}

	flattenedAgent := map[string]interface{}{
		"created_at":       agent.CreatedAt.UTC().String(),
		"description":      agent.Description,
		"updated_at":       agent.UpdatedAt.UTC().String(),
		"if_case":          agent.IfCase,
		"instruction":      agent.Instruction,
		"k":                agent.K,
		"max_tokens":       agent.MaxTokens,
		"name":             agent.Name,
		"project_id":       agent.ProjectId,
		"region":           agent.Region,
		"retrieval_method": agent.RetrievalMethod,
		"route_created_at": agent.RouteCreatedAt.UTC().String(),
		"route_created_by": agent.RouteCreatedBy,
		"route_uuid":       agent.RouteUuid,
		"route_name":       agent.RouteName,
		"tags":             agent.Tags,
		"temperature":      agent.Temperature,
		"top_p":            agent.TopP,
		"url":              agent.Url,
		"user_id":          agent.UserId,
		"agent_id":         agent.Uuid,
	}

	if agent.Model != nil {
		if agent.Model.Uuid != "" {
			flattenedAgent["model_uuid"] = agent.Model.Uuid
		}
		modelSlice := []*godo.Model{agent.Model}
		flattenedAgent["model"] = flattenModel(modelSlice)
	} else {
		flattenedAgent["model"] = []interface{}{}
	}
	if agent.AnthropicApiKey != nil {
		flattenedAgent["anthropic_api_key"] = flattenAnthropicApiKey(agent.AnthropicApiKey)
	} else {
		flattenedAgent["anthropic_api_key"] = []interface{}{}
	}

	if agent.ApiKeyInfos != nil {
		flattenedAgent["api_key_infos"] = flattenApiKeyInfos(agent.ApiKeyInfos)
	} else {
		flattenedAgent["api_key_infos"] = []interface{}{}
	}

	if agent.ApiKeys != nil {
		flattenedAgent["api_keys"] = flattenApiKeys(agent.ApiKeys)
	} else {
		flattenedAgent["api_keys"] = []interface{}{}
	}

	if agent.ChatBot != nil {
		flattenedAgent["chatbot"] = flattenChatbot(agent.ChatBot)
	} else {
		flattenedAgent["chatbot"] = []interface{}{}
	}

	if agent.ChatbotIdentifiers != nil {
		flattenedAgent["chatbot_identifiers"] = flattenChatbotIdentifiers(agent.ChatbotIdentifiers)
	} else {
		flattenedAgent["chatbot_identifiers"] = []interface{}{}
	}
	if agent.ParentAgents != nil {
		flattenedParents := make([]interface{}, 0, len(agent.ParentAgents))
		for _, parent := range agent.ParentAgents {
			if parent != nil {
				flatParent, err := FlattenDigitalOceanAgent(parent)
				if err != nil {
					return nil, err
				}
				flattenedParents = append(flattenedParents, flatParent)
			}
		}
		flattenedAgent["parent_agents"] = flattenedParents
	} else {
		flattenedAgent["parent_agents"] = []interface{}{}
	}
	if agent.ChildAgents != nil {
		flattenedChilds := make([]interface{}, 0, len(agent.ChildAgents))
		for _, child := range agent.ChildAgents {
			if child != nil {
				flatParent, err := FlattenDigitalOceanAgent(child)
				if err != nil {
					return nil, err
				}
				flattenedChilds = append(flattenedChilds, flatParent)
			}
		}
		flattenedAgent["child_agents"] = flattenedChilds
	} else {
		flattenedAgent["child_agents"] = []interface{}{}
	}
	if agent.Guardrails != nil {
		flattenedAgent["agent_guardrail"] = flattenAgentGuardrail(agent.Guardrails)
	} else {
		flattenedAgent["agent_guardrail"] = []interface{}{}
	}

	if agent.KnowledgeBases != nil {
		flattenedAgent["knowledge_bases"] = flattenKnowledgeBases(agent.KnowledgeBases)
	} else {
		flattenedAgent["knowledge_bases"] = []interface{}{}
	}

	if agent.Template != nil {
		flattenedAgent["template"] = flattenTemplate(agent.Template)
	} else {
		flattenedAgent["template"] = []interface{}{}
	}

	return flattenedAgent, nil
}

func FlattenDigitalOceanAgent(agent *godo.Agent) (map[string]interface{}, error) {
	if agent == nil {
		return nil, fmt.Errorf("agent is nil")
	}
	result := map[string]interface{}{
		"created_at":       agent.CreatedAt.UTC().String(),
		"description":      agent.Description,
		"updated_at":       agent.UpdatedAt.UTC().String(),
		"if_case":          agent.IfCase,
		"instruction":      agent.Instruction,
		"k":                agent.K,
		"max_tokens":       agent.MaxTokens,
		"name":             agent.Name,
		"project_id":       agent.ProjectId,
		"region":           agent.Region,
		"retrieval_method": agent.RetrievalMethod,
		"route_created_at": agent.RouteCreatedAt.UTC().String(),
		"route_created_by": agent.RouteCreatedBy,
		"route_uuid":       agent.RouteUuid,
		"route_name":       agent.RouteName,
		"tags":             agent.Tags,
		"temperature":      agent.Temperature,
		"top_p":            agent.TopP,
		"url":              agent.Url,
		"user_id":          agent.UserId,
		"agent_id":         agent.Uuid,
	}

	if agent.Model != nil {
		if agent.Model.Uuid != "" {
			result["model_uuid"] = agent.Model.Uuid
		}
		modelSlice := []*godo.Model{agent.Model}
		result["model"] = flattenModel(modelSlice)
	} else {
		result["model"] = []interface{}{}
	}
	if agent.AnthropicApiKey != nil {
		result["anthropic_api_key"] = flattenAnthropicApiKey(agent.AnthropicApiKey)
	} else {
		result["anthropic_api_key"] = []interface{}{}
	}

	if agent.ApiKeyInfos != nil {
		result["api_key_infos"] = flattenApiKeyInfos(agent.ApiKeyInfos)
	} else {
		result["api_key_infos"] = []interface{}{}
	}

	if agent.ApiKeys != nil {
		result["api_keys"] = flattenApiKeys(agent.ApiKeys)
	} else {
		result["api_keys"] = []interface{}{}
	}

	if agent.ChatBot != nil {
		result["chatbot"] = flattenChatbot(agent.ChatBot)
	} else {
		result["chatbot"] = []interface{}{}
	}

	if agent.ChatbotIdentifiers != nil {
		result["chatbot_identifiers"] = flattenChatbotIdentifiers(agent.ChatbotIdentifiers)
	} else {
		result["chatbot_identifiers"] = []interface{}{}
	}
	if agent.ParentAgents != nil {
		flattenedParents := make([]interface{}, 0, len(agent.ParentAgents))
		for _, parent := range agent.ParentAgents {
			if parent != nil {
				flatParent, err := FlattenDigitalOceanAgent(parent)
				if err != nil {
					return nil, err
				}
				flattenedParents = append(flattenedParents, flatParent)
			}
		}
		result["parent_agents"] = flattenedParents
	} else {
		result["parent_agents"] = []interface{}{}
	}
	if agent.ChildAgents != nil {
		flattenedChilds := make([]interface{}, 0, len(agent.ChildAgents))
		for _, child := range agent.ChildAgents {
			if child != nil {
				flatParent, err := FlattenDigitalOceanAgent(child)
				if err != nil {
					return nil, err
				}
				flattenedChilds = append(flattenedChilds, flatParent)
			}
		}
		result["child_agents"] = flattenedChilds
	} else {
		result["child_agents"] = []interface{}{}
	}
	if agent.Guardrails != nil {
		result["agent_guardrail"] = flattenAgentGuardrail(agent.Guardrails)
	} else {
		result["agent_guardrail"] = []interface{}{}
	}

	if agent.KnowledgeBases != nil {
		result["knowledge_bases"] = flattenKnowledgeBases(agent.KnowledgeBases)
	} else {
		result["knowledge_bases"] = []interface{}{}
	}

	if agent.Template != nil {
		result["template"] = flattenTemplate(agent.Template)
	} else {
		result["template"] = []interface{}{}
	}

	return result, nil
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
		return []interface{}{}
	}
	result := make([]interface{}, 0, len(childAgents))
	for _, child := range childAgents {
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
		return []interface{}{}
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
		return []interface{}{}
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
		return []interface{}{}
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
		return []interface{}{}
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
		return []interface{}{}
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
		return []interface{}{}
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
			k["last_indexing_job"] = flattenLastIndexingJob(kb.LastIndexingJob)
		}

		result = append(result, k)
	}

	return result
}

func flattenModel(models []*godo.Model) []interface{} {
	if models == nil {
		return []interface{}{}
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
		return []interface{}{}
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

func flattenChatbotIdentifiers(chatbotIdentifiers []*godo.AgentChatbotIdentifier) []interface{} {
	if chatbotIdentifiers == nil {
		return []interface{}{}
	}

	result := make([]interface{}, 0, len(chatbotIdentifiers))
	for _, identifier := range chatbotIdentifiers {
		if identifier != nil {
			m := map[string]interface{}{
				"chatbot_id": identifier.AgentChatbotIdentifier,
			}
			result = append(result, m)
		}
	}

	return result
}

func flattenOpenAiApiKey(apiKey *godo.OpenAiApiKey) []interface{} {
	if apiKey == nil {
		return []interface{}{}
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
		return []interface{}{}
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
		return []interface{}{}
	}

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
