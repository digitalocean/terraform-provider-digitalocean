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
		var flattenedKBs []interface{}
		for _, kb := range agent.KnowledgeBases {
			if kb != nil {
				flatKB, err := FlattenDigitalOceanKnowledgeBase(kb)
				if err != nil {
					return nil, err
				}
				flattenedKBs = append(flattenedKBs, flatKB)
			}
		}
		result["knowledge_bases"] = flattenedKBs
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

func FlattenDigitalOceanAgents(agents []*godo.Agent) ([]interface{}, error) {
	var result []interface{}
	for _, agent := range agents {
		if agent == nil {
			continue
		}
		flat, err := FlattenDigitalOceanAgent(agent)
		if err != nil {
			return nil, err
		}
		result = append(result, flat)
	}
	return result, nil
}

func getDigitalOceanOpenAIApiKeys(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.CombinedConfig).GodoClient()

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	var allOpenAIApiKeys []interface{}
	for {
		openAIApiKeys, resp, err := client.GenAI.ListOpenAIAPIKeys(context.Background(), opts)
		if err != nil {
			return nil, fmt.Errorf("error retrieving OpenAI API keys: %s", err)
		}

		for _, openAIApiKey := range openAIApiKeys {
			if openAIApiKey != nil {
				allOpenAIApiKeys = append(allOpenAIApiKeys, openAIApiKey)
			}
		}
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("error retrieving OpenAI API keys: %s", err)
		}
		opts.Page = page + 1
	}
	return allOpenAIApiKeys, nil
}

func flattenOpenAIApiKeyInfo(rawDomain, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	openAIApiKey, ok := rawDomain.(*godo.OpenAiApiKey)
	if !ok || openAIApiKey == nil {
		// Return nil without error to safely skip nil or wrong type entries
		return nil, nil
	}
	return FlattenOpenAIApiKeyInfo(openAIApiKey)
}

func FlattenOpenAIApiKeyInfo(openAIApiKey *godo.OpenAiApiKey) (map[string]interface{}, error) {
	if openAIApiKey == nil {
		return nil, nil
	}

	result := map[string]interface{}{
		"created_by": openAIApiKey.CreatedBy,
		"name":       openAIApiKey.Name,
		"uuid":       openAIApiKey.Uuid,
		"models":     flattenModel(openAIApiKey.Models),
	}

	if openAIApiKey.DeletedAt != nil {
		result["deleted_at"] = openAIApiKey.DeletedAt.UTC().String()
	} else {
		result["deleted_at"] = ""
	}
	if openAIApiKey.CreatedAt != nil {
		result["created_at"] = openAIApiKey.CreatedAt.UTC().String()
	} else {
		result["created_at"] = ""
	}
	if openAIApiKey.UpdatedAt != nil {
		result["updated_at"] = openAIApiKey.UpdatedAt.UTC().String()
	} else {
		result["updated_at"] = ""
	}

	return result, nil
}

func flattenDigitalOceanModel(rawDomain, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	model, ok := rawDomain.(*godo.Model)
	if !ok {
		return nil, nil
	}
	if model == nil {
		return nil, nil
	}

	return FlattenDigitalOceanModel(model)
}

func FlattenDigitalOceanModel(model *godo.Model) (map[string]interface{}, error) {
	if model == nil {
		return nil, nil
	}

	result := map[string]interface{}{
		"is_foundational": model.IsFoundational,
		"name":            model.Name,
		"parent_uuid":     model.ParentUuid,
		"upload_complete": model.UploadComplete,
		"url":             model.Url,
		"uuid":            model.Uuid,
	}

	// Handle timestamps
	if model.CreatedAt != nil {
		result["created_at"] = model.CreatedAt.UTC().String()
	}
	if model.UpdatedAt != nil {
		result["updated_at"] = model.UpdatedAt.UTC().String()
	}

	// Handle agreement
	if model.Agreement != nil {
		result["agreement"] = []interface{}{
			map[string]interface{}{
				"description": model.Agreement.Description,
				"name":        model.Agreement.Name,
				"url":         model.Agreement.Url,
				"uuid":        model.Agreement.Uuid,
			},
		}
	} else {
		result["agreement"] = []interface{}{}
	}

	// Handle version
	if model.Version != nil {
		result["version"] = []interface{}{
			map[string]interface{}{
				"major": model.Version.Major,
				"minor": model.Version.Minor,
				"patch": model.Version.Patch,
			},
		}
	} else {
		result["version"] = []interface{}{}
	}

	return result, nil
}

func getDigitalOceanModels(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.CombinedConfig).GodoClient()

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	var allModels []interface{}
	for {
		models, resp, err := client.GenAI.ListAvailableModels(context.Background(), opts)
		if err != nil {
			return nil, fmt.Errorf("error retrieving models : %s", err)
		}

		for i := range models {
			allModels = append(allModels, models[i])
		}
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("error retrieving models: %s", err)
		}

		opts.Page = page + 1
	}
	return allModels, nil
}

func flattenDigitalOceanRegion(rawDomain, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	region, ok := rawDomain.(*godo.DatacenterRegions)
	if !ok || region == nil {
		// Return nil without error to safely skip nil or wrong type entries
		return nil, nil
	}
	return FlattenDigitalOceanRegion(region)
}

func FlattenDigitalOceanRegion(region *godo.DatacenterRegions) (map[string]interface{}, error) {
	if region == nil {
		return nil, nil
	}

	result := map[string]interface{}{
		"region":               region.Region,
		"inference_url":        region.InferenceUrl,
		"serves_batch":         region.ServesBatch,
		"serves_inference":     region.ServesInference,
		"stream_inference_url": region.StreamInferenceUrl,
	}

	return result, nil
}

func getDigitalOceanRegions(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.CombinedConfig).GodoClient()

	var allRegions []interface{}
	servesInference := (*bool)(nil)
	servesBatch := (*bool)(nil)

	regions, _, err := client.GenAI.ListDatacenterRegions(context.Background(), servesInference, servesBatch)
	if err != nil {
		return nil, fmt.Errorf("error retrieving regions : %s", err)
	}

	for i := range regions {
		allRegions = append(allRegions, regions[i])
	}
	return allRegions, nil
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

func getDigitalOceanKnowledgeBases(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.CombinedConfig).GodoClient()

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	var allKnowledgeBases []interface{}
	for {
		knowledge_bases, resp, err := client.GenAI.ListKnowledgeBases(context.Background(), opts)
		if err != nil {
			return nil, fmt.Errorf("error retrieving knowledge bases : %s", err)
		}

		for i := range knowledge_bases {
			allKnowledgeBases = append(allKnowledgeBases, &knowledge_bases[i])
		}
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("error retrieving knowledge bases: %s", err)
		}

		opts.Page = page + 1
	}
	return allKnowledgeBases, nil
}

func flattenDigitalOceanKnowledgeBase(rawDomain, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	kb, ok := rawDomain.(*godo.KnowledgeBase)
	if !ok {
		return nil, fmt.Errorf("expected *godo.KnowledgeBase, got %T", rawDomain)
	}

	if kb == nil {
		return nil, fmt.Errorf("knowledgeBase is nil")
	}

	return FlattenDigitalOceanKnowledgeBase(kb)
}

func FlattenDigitalOceanKnowledgeBase(kb *godo.KnowledgeBase) (map[string]interface{}, error) {
	if kb == nil {
		return nil, fmt.Errorf("knowledgeBase is nil")
	}

	flattenedKnowledgeBase := map[string]interface{}{
		"uuid":                 kb.Uuid,
		"name":                 kb.Name,
		"project_id":           kb.ProjectId,
		"region":               kb.Region,
		"embedding_model_uuid": kb.EmbeddingModelUuid,
		"database_id":          kb.DatabaseId,
		"is_public":            kb.IsPublic,
		"user_id":              kb.UserId,
	}

	if kb.CreatedAt != nil {
		flattenedKnowledgeBase["created_at"] = kb.CreatedAt.UTC().String()
	}
	if kb.UpdatedAt != nil {
		flattenedKnowledgeBase["updated_at"] = kb.UpdatedAt.UTC().String()
	}
	if kb.AddedToAgentAt != nil {
		flattenedKnowledgeBase["added_to_agent_at"] = kb.AddedToAgentAt.UTC().String()
	}

	// Tags as []string
	if kb.Tags != nil {
		tags := make([]interface{}, len(kb.Tags))
		for i, tag := range kb.Tags {
			tags[i] = tag
		}
		flattenedKnowledgeBase["tags"] = tags
	} else {
		flattenedKnowledgeBase["tags"] = []interface{}{}
	}

	// Flatten last_indexing_job as a map (not []interface{})
	if kb.LastIndexingJob != nil {
		flatJob := flattenKnowledgeBaseLastIndexingJob(kb.LastIndexingJob)
		flattenedKnowledgeBase["last_indexing_job"] = flatJob
	} else {
		flattenedKnowledgeBase["last_indexing_job"] = []interface{}{}
	}

	return flattenedKnowledgeBase, nil
}

func flattenKnowledgeBaseLastIndexingJob(lastJob *godo.LastIndexingJob) []interface{} {
	if lastJob == nil {
		return []interface{}{}
	}

	jobMap := map[string]interface{}{
		"uuid":                  lastJob.Uuid,
		"knowledge_base_uuid":   lastJob.KnowledgeBaseUuid,
		"phase":                 lastJob.Phase,
		"completed_datasources": lastJob.CompletedDatasources,
		"total_datasources":     lastJob.TotalDatasources,
		"tokens":                lastJob.Tokens,
	}

	// Handle data source UUIDs
	if lastJob.DataSourceUuids != nil {
		dataSourceUuids := make([]interface{}, len(lastJob.DataSourceUuids))
		for i, uuid := range lastJob.DataSourceUuids {
			dataSourceUuids[i] = uuid
		}
		jobMap["data_source_uuids"] = dataSourceUuids
	} else {
		jobMap["data_source_uuids"] = []interface{}{}
	}

	// Handle timestamps
	if lastJob.CreatedAt != nil {
		jobMap["created_at"] = lastJob.CreatedAt.UTC().String()
	}
	if lastJob.UpdatedAt != nil {
		jobMap["updated_at"] = lastJob.UpdatedAt.UTC().String()
	}
	if lastJob.StartedAt != nil {
		jobMap["started_at"] = lastJob.StartedAt.UTC().String()
	}
	if lastJob.FinishedAt != nil {
		jobMap["finished_at"] = lastJob.FinishedAt.UTC().String()
	}

	return []interface{}{jobMap}
}

// flattenKnowledgeBaseIndexingJobs flattens a slice of LastIndexingJob structs
func flattenKnowledgeBaseIndexingJobs(jobs []godo.LastIndexingJob) []interface{} {
	if len(jobs) == 0 {
		return []interface{}{}
	}

	flattenedJobs := make([]interface{}, len(jobs))
	for i, job := range jobs {
		jobMap := map[string]interface{}{
			"uuid":                  job.Uuid,
			"knowledge_base_uuid":   job.KnowledgeBaseUuid,
			"phase":                 job.Phase,
			"status":                job.Status,
			"completed_datasources": job.CompletedDatasources,
			"total_datasources":     job.TotalDatasources,
			"tokens":                job.Tokens,
			"total_items_failed":    job.TotalItemsFailed,
			"total_items_indexed":   job.TotalItemsIndexed,
			"total_items_skipped":   job.TotalItemsSkipped,
		}

		// Handle data source UUIDs
		if job.DataSourceUuids != nil {
			dataSourceUuids := make([]interface{}, len(job.DataSourceUuids))
			for j, uuid := range job.DataSourceUuids {
				dataSourceUuids[j] = uuid
			}
			jobMap["data_source_uuids"] = dataSourceUuids
		} else {
			jobMap["data_source_uuids"] = []interface{}{}
		}

		// Handle timestamps
		if job.CreatedAt != nil {
			jobMap["created_at"] = job.CreatedAt.UTC().String()
		}
		if job.UpdatedAt != nil {
			jobMap["updated_at"] = job.UpdatedAt.UTC().String()
		}
		if job.StartedAt != nil {
			jobMap["started_at"] = job.StartedAt.UTC().String()
		}
		if job.FinishedAt != nil {
			jobMap["finished_at"] = job.FinishedAt.UTC().String()
		}

		flattenedJobs[i] = jobMap
	}

	return flattenedJobs
}

// flattenIndexedDataSources flattens a slice of IndexedDataSource structs
func flattenIndexedDataSources(dataSources []godo.IndexedDataSource) []interface{} {
	if len(dataSources) == 0 {
		return []interface{}{}
	}

	flattenedDataSources := make([]interface{}, len(dataSources))
	for i, ds := range dataSources {
		dsMap := map[string]interface{}{
			"data_source_uuid":    ds.DataSourceUuid,
			"status":              ds.Status,
			"error_msg":           ds.ErrorMsg,
			"error_details":       ds.ErrorDetails,
			"total_file_count":    ds.TotalFileCount,
			"indexed_file_count":  ds.IndexedFileCount,
			"indexed_item_count":  ds.IndexedItemCount,
			"failed_item_count":   ds.FailedItemCount,
			"skipped_item_count":  ds.SkippedItemCount,
			"removed_item_count":  ds.RemovedItemCount,
			"total_bytes":         ds.TotalBytes,
			"total_bytes_indexed": ds.TotalBytesIndexed,
		}

		// Handle timestamps
		if ds.StartedAt != nil {
			dsMap["started_at"] = ds.StartedAt.UTC().String()
		}
		if ds.CompletedAt != nil {
			dsMap["completed_at"] = ds.CompletedAt.UTC().String()
		}

		flattenedDataSources[i] = dsMap
	}

	return flattenedDataSources
}

// flattenIndexingJob flattens a single LastIndexingJob struct
func flattenIndexingJob(job *godo.LastIndexingJob) map[string]interface{} {
	if job == nil {
		return map[string]interface{}{}
	}

	jobMap := map[string]interface{}{
		"uuid":                  job.Uuid,
		"knowledge_base_uuid":   job.KnowledgeBaseUuid,
		"phase":                 job.Phase,
		"status":                job.Status,
		"completed_datasources": job.CompletedDatasources,
		"total_datasources":     job.TotalDatasources,
		"tokens":                job.Tokens,
		"total_items_failed":    job.TotalItemsFailed,
		"total_items_indexed":   job.TotalItemsIndexed,
		"total_items_skipped":   job.TotalItemsSkipped,
	}

	// Handle data source UUIDs
	if job.DataSourceUuids != nil {
		dataSourceUuids := make([]interface{}, len(job.DataSourceUuids))
		for j, uuid := range job.DataSourceUuids {
			dataSourceUuids[j] = uuid
		}
		jobMap["data_source_uuids"] = dataSourceUuids
	} else {
		jobMap["data_source_uuids"] = []interface{}{}
	}

	// Handle timestamps
	if job.CreatedAt != nil {
		jobMap["created_at"] = job.CreatedAt.UTC().String()
	}
	if job.UpdatedAt != nil {
		jobMap["updated_at"] = job.UpdatedAt.UTC().String()
	}
	if job.StartedAt != nil {
		jobMap["started_at"] = job.StartedAt.UTC().String()
	}
	if job.FinishedAt != nil {
		jobMap["finished_at"] = job.FinishedAt.UTC().String()
	}

	return jobMap
}

// flattenKnowledgeBaseFileUploadDataSource flattens a FileUploadDataSource struct
func flattenKnowledgeBaseFileUploadDataSource(fileUpload *godo.FileUploadDataSource) []interface{} {
	if fileUpload == nil {
		return []interface{}{}
	}

	fileUploadMap := map[string]interface{}{
		"original_file_name": fileUpload.OriginalFileName,
		"size":               fileUpload.Size,
		"stored_object_key":  fileUpload.StoredObjectKey,
	}

	return []interface{}{fileUploadMap}
}

// flattenKnowledgeBaseSpacesDataSource flattens a SpacesDataSource struct
func flattenKnowledgeBaseSpacesDataSource(spaces *godo.SpacesDataSource) []interface{} {
	if spaces == nil {
		return []interface{}{}
	}

	spacesMap := map[string]interface{}{
		"bucket_name": spaces.BucketName,
		"item_path":   spaces.ItemPath,
		"region":      spaces.Region,
	}

	return []interface{}{spacesMap}
}

// flattenKnowledgeBaseWebCrawlerDataSource flattens a WebCrawlerDataSource struct
func flattenKnowledgeBaseWebCrawlerDataSource(webCrawler *godo.WebCrawlerDataSource) []interface{} {
	if webCrawler == nil {
		return []interface{}{}
	}

	webCrawlerMap := map[string]interface{}{
		"base_url":        webCrawler.BaseUrl,
		"crawling_option": webCrawler.CrawlingOption,
		"embed_media":     webCrawler.EmbedMedia,
	}

	return []interface{}{webCrawlerMap}
}

// flattenKnowledgeBaseDataSources flattens a slice of KnowledgeBaseDataSource structs
func flattenKnowledgeBaseDataSources(dataSources []godo.KnowledgeBaseDataSource) []interface{} {
	if len(dataSources) == 0 {
		return []interface{}{}
	}

	flattenedDataSources := make([]interface{}, len(dataSources))
	for i, ds := range dataSources {
		dsMap := map[string]interface{}{
			"uuid": ds.Uuid,
		}

		// Handle timestamps for data source
		if ds.CreatedAt != nil {
			dsMap["created_at"] = ds.CreatedAt.UTC().String()
		}
		if ds.UpdatedAt != nil {
			dsMap["updated_at"] = ds.UpdatedAt.UTC().String()
		}

		// Handle nested data sources using the separate flatten functions
		dsMap["file_upload_data_source"] = flattenKnowledgeBaseFileUploadDataSource(ds.FileUploadDataSource)
		dsMap["spaces_data_source"] = flattenKnowledgeBaseSpacesDataSource(ds.SpacesDataSource)
		dsMap["web_crawler_data_source"] = flattenKnowledgeBaseWebCrawlerDataSource(ds.WebCrawlerDataSource)
		dsMap["last_indexing_job"] = flattenKnowledgeBaseLastIndexingJob(ds.LastIndexingJob)

		flattenedDataSources[i] = dsMap
	}

	return flattenedDataSources
}

// expandKnowledgeBaseDatasources converts Terraform schema data to slice of godo.KnowledgeBaseDatasource
func expandKnowledgeBaseDatasources(rawDatasources []interface{}) []godo.KnowledgeBaseDataSource {
	if len(rawDatasources) == 0 {
		return nil
	}

	var datasources []godo.KnowledgeBaseDataSource

	for _, rawDS := range rawDatasources {
		if rawDS == nil {
			continue
		}

		dsMap := rawDS.(map[string]interface{})
		ds := godo.KnowledgeBaseDataSource{}

		// Process nested datasources - only one should be set
		if fileUploadRaw, ok := dsMap["file_upload_data_source"].([]interface{}); ok && len(fileUploadRaw) > 0 {
			ds.FileUploadDataSource = expandFileUploadDataSource(fileUploadRaw)
		}

		if spacesRaw, ok := dsMap["spaces_data_source"].([]interface{}); ok && len(spacesRaw) > 0 {
			ds.SpacesDataSource = expandSpacesDataSource(spacesRaw)
		}

		if webCrawlerRaw, ok := dsMap["web_crawler_data_source"].([]interface{}); ok && len(webCrawlerRaw) > 0 {
			ds.WebCrawlerDataSource = expandWebCrawlerDataSource(webCrawlerRaw)
		}

		datasources = append(datasources, ds)
	}

	return datasources
}

// expandFileUploadDataSource converts Terraform schema data to godo.FileUploadDataSource
func expandFileUploadDataSource(rawFileUpload []interface{}) *godo.FileUploadDataSource {
	if len(rawFileUpload) == 0 || rawFileUpload[0] == nil {
		return nil
	}

	fileUploadMap := rawFileUpload[0].(map[string]interface{})
	fileUpload := &godo.FileUploadDataSource{}

	if originalFileName, ok := fileUploadMap["original_file_name"].(string); ok {
		fileUpload.OriginalFileName = originalFileName
	}

	if size, ok := fileUploadMap["size"].(string); ok {
		fileUpload.Size = size
	}

	if storedObjectKey, ok := fileUploadMap["stored_object_key"].(string); ok {
		fileUpload.StoredObjectKey = storedObjectKey
	}

	return fileUpload
}

// expandSpacesDataSource converts Terraform schema data to godo.SpacesDataSource
func expandSpacesDataSource(rawSpaces []interface{}) *godo.SpacesDataSource {
	if len(rawSpaces) == 0 || rawSpaces[0] == nil {
		return nil
	}

	spacesMap := rawSpaces[0].(map[string]interface{})
	spaces := &godo.SpacesDataSource{}

	if bucketName, ok := spacesMap["bucket_name"].(string); ok {
		spaces.BucketName = bucketName
	}

	if itemPath, ok := spacesMap["item_path"].(string); ok {
		spaces.ItemPath = itemPath
	}

	if region, ok := spacesMap["region"].(string); ok {
		spaces.Region = region
	}

	return spaces
}

// expandWebCrawlerDataSource converts Terraform schema data to godo.WebCrawlerDataSource
func expandWebCrawlerDataSource(rawWebCrawler []interface{}) *godo.WebCrawlerDataSource {
	if len(rawWebCrawler) == 0 || rawWebCrawler[0] == nil {
		return nil
	}

	webCrawlerMap := rawWebCrawler[0].(map[string]interface{})
	webCrawler := &godo.WebCrawlerDataSource{}

	if baseUrl, ok := webCrawlerMap["base_url"].(string); ok {
		webCrawler.BaseUrl = baseUrl
	}

	if crawlingOption, ok := webCrawlerMap["crawling_option"].(string); ok {
		webCrawler.CrawlingOption = crawlingOption
	}

	if embedMedia, ok := webCrawlerMap["embed_media"].(bool); ok {
		webCrawler.EmbedMedia = embedMedia
	}

	return webCrawler
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
