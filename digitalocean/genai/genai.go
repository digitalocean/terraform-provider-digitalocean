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
