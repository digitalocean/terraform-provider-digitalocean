package genai

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
)

func getDigitalOceanAgents(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.CombinedConfig).GodoClient()

	only_deployed := false
	if v, ok := extra["only_deployed"]; ok {
		if b, ok := v.(bool); ok {
			only_deployed = b
		} else {
			return nil, fmt.Errorf("only deployed can only be a boolean value")
		}
	}

	opts := &godo.ListOptions{
		Page:     1,
		PerPage:  200,
		Deployed: only_deployed,
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

	return FlattenDigitalOceanAgent(agent)
}

func FlattenDigitalOceanAgent(agent *godo.Agent) (map[string]interface{}, error) {
	if agent == nil {
		return nil, fmt.Errorf("agent is nil")
	}
	result := map[string]interface{}{
		"description":      agent.Description,
		"if_case":          agent.IfCase,
		"instruction":      agent.Instruction,
		"k":                agent.K,
		"max_tokens":       agent.MaxTokens,
		"name":             agent.Name,
		"project_id":       agent.ProjectId,
		"region":           agent.Region,
		"retrieval_method": agent.RetrievalMethod,
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
	if agent.CreatedAt != nil {
		result["created_at"] = agent.CreatedAt.UTC().String()
	}
	if agent.UpdatedAt != nil {
		result["updated_at"] = agent.UpdatedAt.UTC().String()
	}
	if agent.RouteCreatedAt != nil {
		result["route_created_at"] = agent.RouteCreatedAt.UTC().String()
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
	if agent.Deployment != nil {
		result["deployment"] = flattenDeployment(agent.Deployment)
	} else {
		result["deployment"] = []interface{}{}
	}

	if agent.ChatbotIdentifiers != nil {
		result["chatbot_identifiers"] = flattenChatbotIdentifiers(agent.ChatbotIdentifiers)
	} else {
		result["chatbot_identifiers"] = []interface{}{}
	}
	if agent.ParentAgents != nil {
		result["parent_agents"] = flattenRelatedAgents(agent.ParentAgents)
	} else {
		result["parent_agents"] = []interface{}{}
	}
	if agent.ChildAgents != nil {
		result["child_agents"] = flattenRelatedAgents(agent.ChildAgents)
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

	if agent.Tags != nil {
		result["tags"] = flattenAgentTags(agent.Tags)
	} else {
		result["tags"] = []interface{}{}
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

func flattenRelatedAgents(relatedAgents []*godo.Agent) []interface{} {
	if relatedAgents == nil {
		return []interface{}{}
	}
	result := make([]interface{}, 0, len(relatedAgents))
	for _, agent := range relatedAgents {
		m := map[string]interface{}{
			"agent_id":    agent.Uuid,
			"name":        agent.Name,
			"region":      agent.Region,
			"project_id":  agent.ProjectId,
			"description": agent.Description,
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
		"created_by": apiKey.CreatedBy,
		"deleted_at": apiKey.DeletedAt,
		"name":       apiKey.Name,
		"uuid":       apiKey.Uuid,
	}
	if apiKey.CreatedAt != nil {
		m["created_at"] = apiKey.CreatedAt.UTC().String()
	}
	if apiKey.UpdatedAt != nil {
		m["updated_at"] = apiKey.UpdatedAt.UTC().String()
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
			"created_by": info.CreatedBy,
			"deleted_at": info.DeletedAt,
			"name":       info.Name,
			"secret_key": info.SecretKey,
			"uuid":       info.Uuid,
		}
		if info.CreatedAt != nil {
			m["created_at"] = info.CreatedAt.UTC().String()
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
		"name":       deployment.Name,
		"status":     deployment.Status,
		"url":        deployment.Url,
		"uuid":       deployment.Uuid,
		"visibility": deployment.Visibility,
	}
	if deployment.CreatedAt != nil {
		m["created_at"] = deployment.CreatedAt.UTC().String()
	}
	if deployment.UpdatedAt != nil {
		m["updated_at"] = deployment.UpdatedAt.UTC().String()
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
			"description":    fn.Description,
			"guardrail_uuid": fn.GuardrailUuid,
			"faasname":       fn.FaasName,
			"faasnamespace":  fn.FaasNamespace,
			"name":           fn.Name,
			"url":            fn.Url,
			"uuid":           fn.Uuid,
		}
		if fn.CreatedAt != nil {
			m["created_at"] = fn.CreatedAt.UTC().String()
		}
		if fn.UpdatedAt != nil {
			m["updated_at"] = fn.UpdatedAt.UTC().String()
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
			"default_response": guardrail.DefaultResponse,
			"description":      guardrail.Description,
			"guardrail_uuid":   guardrail.GuardrailUuid,
			"is_attached":      guardrail.IsAttached,
			"is_default":       guardrail.IsDefault,
			"name":             guardrail.Name,
			"priority":         guardrail.Priority,
			"type":             guardrail.Type,
			"uuid":             guardrail.Uuid,
		}
		if guardrail.CreatedAt != nil {
			m["created_at"] = guardrail.CreatedAt.UTC().String()
		}
		if guardrail.UpdatedAt != nil {
			m["updated_at"] = guardrail.UpdatedAt.UTC().String()
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
			"tags":                 kb.Tags,
			"region":               kb.Region,
			"embedding_model_uuid": kb.EmbeddingModelUuid,
			"project_id":           kb.ProjectId,
			"database_id":          kb.DatabaseId,
		}
		if kb.CreatedAt != nil {
			k["created_at"] = kb.CreatedAt.UTC().String()
		}
		if kb.UpdatedAt != nil {
			k["updated_at"] = kb.UpdatedAt.UTC().String()
		}
		if kb.AddedToAgentAt != nil {
			k["added_to_agent_at"] = kb.AddedToAgentAt.UTC().String()
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
			"inference_name":    model.InferenceName,
			"inference_version": model.InferenceVersion,
			"is_foundational":   model.IsFoundational,
			"name":              model.Name,
			"parent_uuid":       model.ParentUuid,
			"provider":          model.Provider,
			"upload_complete":   model.UploadComplete,
			"url":               model.Url,
			"usecases":          model.Usecases,
		}

		if model.CreatedAt != nil {
			m["created_at"] = model.CreatedAt.UTC().String()
		}
		if model.UpdatedAt != nil {
			m["updated_at"] = model.UpdatedAt.UTC().String()
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
		"created_by": apiKey.CreatedBy,
		"name":       apiKey.Name,
		"uuid":       apiKey.Uuid,
		"model":      flattenModel(apiKey.Models),
	}

	if apiKey.CreatedAt != nil {
		m["created_at"] = apiKey.CreatedAt.UTC().String()
	}
	if apiKey.UpdatedAt != nil {
		m["updated_at"] = apiKey.UpdatedAt.UTC().String()
	}
	if apiKey.DeletedAt != nil {
		m["deleted_at"] = apiKey.DeletedAt.UTC().String()
	}

	return []interface{}{m}
}

func flattenTemplate(template *godo.AgentTemplate) []interface{} {
	if template == nil {
		return []interface{}{}
	}

	m := map[string]interface{}{
		"instruction": template.Instruction,
		"description": template.Description,
		"k":           template.K,
		"max_tokens":  template.MaxTokens,
		"name":        template.Name,
		"temperature": template.Temperature,
		"top_p":       template.TopP,
		"uuid":        template.Uuid,
	}
	if template.CreatedAt != nil {
		m["created_at"] = template.CreatedAt.UTC().String()
	}

	if template.UpdatedAt != nil {
		m["updated_at"] = template.UpdatedAt.UTC().String()
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
		"datasource_uuids":      datasourceUuids,
		"knowledge_base_uuid":   job.KnowledgeBaseUuid,
		"phase":                 job.Phase,
		"tokens":                job.Tokens,
		"total_datasources":     job.TotalDatasources,
		"uuid":                  job.Uuid,
	}
	if job.CreatedAt != nil {
		m["created_at"] = job.CreatedAt.UTC().String()
	}
	if job.FinishedAt != nil {
		m["finished_at"] = job.FinishedAt.UTC().String()
	}
	if job.UpdatedAt != nil {
		m["updated_at"] = job.UpdatedAt.UTC().String()
	}
	if job.StartedAt != nil {
		m["started_at"] = job.StartedAt.UTC().String()
	}

	return []interface{}{m}
}

func getDigitalOceanAgentVersions(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.CombinedConfig).GodoClient()

	agentID, ok := extra["agent_id"].(string)
	if !ok {
		return nil, fmt.Errorf("agent_id is not defined or not a string")
	}

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	var allAgentVersions []interface{}
	for {
		agentVersions, resp, err := client.GenAI.ListAgentVersions(context.Background(), agentID, opts)
		if err != nil {
			return nil, fmt.Errorf("error retrieving agents : %s", err)
		}

		for _, agent := range agentVersions {
			allAgentVersions = append(allAgentVersions, agent)
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
	return allAgentVersions, nil

}

func flattenDigitalOceanAgentVersion(rawDomain, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	agentVersions, ok := rawDomain.(*godo.AgentVersion)
	if !ok {
		return nil, fmt.Errorf("expected *godo.AgentVersion, got %T", rawDomain)
	}

	if agentVersions == nil {
		return nil, fmt.Errorf("agent versions are nil")
	}

	result := map[string]interface{}{
		"can_rollback":      agentVersions.CanRollback,
		"created_by_email":  agentVersions.CreatedByEmail,
		"currently_applied": agentVersions.CurrentlyApplied,
		"description":       agentVersions.Description,
		"id":                agentVersions.ID,
		"instruction":       agentVersions.Instruction,
		"k":                 agentVersions.K,
		"max_tokens":        agentVersions.MaxTokens,
		"name":              agentVersions.Name,
		"provide_citations": agentVersions.ProvideCitations,
		"retrieval_method":  agentVersions.RetrievalMethod,
		"temperature":       agentVersions.Temperature,
		"top_p":             agentVersions.TopP,
		"trigger_action":    agentVersions.TriggerAction,
		"agent_uuid":        agentVersions.AgentUuid,
		"version_hash":      agentVersions.VersionHash,
	}
	if agentVersions.CreatedAt != nil {
		result["created_at"] = agentVersions.CreatedAt.UTC().String()
	}

	if agentVersions.AttachedChildAgents != nil {
		result["attached_child_agents"] = flattenAttachedChildAgents(agentVersions.AttachedChildAgents)
	} else {
		result["attached_child_agents"] = []interface{}{}
	}

	if agentVersions.AttachedFunctions != nil {
		result["attached_functions"] = flattenAttachedFunctionsSchema(agentVersions.AttachedFunctions)
	} else {
		result["attached_functions"] = []interface{}{}
	}

	if agentVersions.AttachedGuardrails != nil {
		result["attached_guardrails"] = flattenAttachedGuardRails(agentVersions.AttachedGuardrails)
	} else {
		result["attached_guardrails"] = []interface{}{}
	}

	if agentVersions.AttachedKnowledgeBases != nil {
		result["attached_knowledge_bases"] = flattenAttachedKnowledgeBases(agentVersions.AttachedKnowledgeBases)
	} else {
		result["attached_knowledge_bases"] = []interface{}{}
	}

	if agentVersions.Tags != nil {
		result["tags"] = flattenAgentTags(agentVersions.Tags)
	} else {
		result["tags"] = []interface{}{}
	}

	return result, nil

}

func flattenAttachedChildAgents(childAgents []*godo.AttachedChildAgent) []interface{} {
	if childAgents == nil {
		return []interface{}{}
	}
	result := make([]interface{}, 0, len(childAgents))
	for _, child := range childAgents {
		m := map[string]interface{}{
			"agent_name":       child.AgentName,
			"child_agent_uuid": child.ChildAgentUuid,
			"if_case":          child.IfCase,
			"is_deleted":       child.IsDeleted,
			"route_name":       child.RouteName,
		}
		result = append(result, m)
	}
	return result
}

func flattenAttachedFunctionsSchema(functions []*godo.AgentFunction) []interface{} {
	if functions == nil {
		return []interface{}{}
	}

	result := make([]interface{}, 0, len(functions))
	for _, fn := range functions {
		m := map[string]interface{}{
			"description":    fn.Description,
			"faas_name":      fn.FaasName,
			"faas_namespace": fn.FaasNamespace,
			"is_deleted":     fn.IsDeleted,
			"name":           fn.Name,
		}
		result = append(result, m)
	}

	return result
}
func flattenAttachedGuardRails(guardrails []*godo.AgentGuardrail) []interface{} {
	if guardrails == nil {
		return []interface{}{}
	}

	result := make([]interface{}, 0, len(guardrails))
	for _, guardrail := range guardrails {
		m := map[string]interface{}{
			"is_deleted": guardrail.IsDeleted,
			"name":       guardrail.Name,
			"priority":   guardrail.Priority,
			"uuid":       guardrail.Uuid,
		}
		result = append(result, m)
	}

	return result
}
func flattenAttachedKnowledgeBases(kbs []*godo.KnowledgeBase) []interface{} {
	if kbs == nil {
		return []interface{}{}
	}

	result := make([]interface{}, 0, len(kbs))

	for _, kb := range kbs {
		k := map[string]interface{}{
			"is_deleted": kb.IsDeleted,
			"name":       kb.Name,
			"uuid":       kb.Uuid,
		}
		result = append(result, k)
	}

	return result
}

func flattenAgentTags(tags []string) []interface{} {
	if tags == nil {
		return []interface{}{}
	}

	result := make([]interface{}, 0, len(tags))
	for _, tag := range tags {
		result = append(result, tag)
	}

	return result
}
