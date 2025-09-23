package genai

import (
	"context"
	"fmt"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanKnowledgeBase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanKnowledgeBaseCreate,
		ReadContext:   resourceDigitalOceanKnowledgeBaseRead,
		UpdateContext: resourceDigitalOceanKnowledgeBaseUpdate,
		DeleteContext: resourceDigitalOceanKnowledgeBaseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, 32),
				Description:  "The name of the knowledge base.",
			},
			"project_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The unique identifier of the project to which the knowledge base belongs.",
			},
			"region": {
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(val interface{}) string {
					// DO API V2 region slug is always lowercase
					return strings.ToLower(val.(string))
				},
				ValidateFunc: validation.NoZeroValues,
			},
			"vpc_uuid": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The unique identifier of the VPC to which the knowledge base belongs.",
			},
			"added_to_agent_at": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The time when the knowledge base was added to the agent.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The time when the knowledge base was created.",
			},
			"database_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The unique identifier of the DigitalOcean OpenSearch database this knowledge base will use",
			},
			"embedding_model_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique identifier of the embedding model",
			},
			"is_public": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether the knowledge base is public or private.",
			},
			"last_indexing_job": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The last indexing job for the knowledge base.",
				Elem:        LastIndexingJobSchema(),
			},
			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"datasources": {
				Type:        schema.TypeList,
				Required:    true,
				ForceNew:    true,
				Description: "Data sources for the knowledge base",
				Elem:        knowledgeBaseDatasourcesSchema(),
			},
		},
	}
}

func ResourceDigitalOceanKnowledgeBaseDataSource() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanKnowledgeBaseDatasourceCreate,
		ReadContext:   resourceDigitalOceanKnowledgeBaseDatasourceRead,
		DeleteContext: resourceDigitalOceanKnowledgeBaseDatasourceDelete,
		Schema: map[string]*schema.Schema{
			"knowledge_base_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "UUID of the Knowledge Base",
			},
			"spaces_data_source": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Elem:     spacesDataSourceSchema(),
			},
			"web_crawler_data_source": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				ForceNew: true,
				Elem:     webCrawlerDataSourceSchema(),
			},
		},
	}
}

func ResourceDigitalOceanAgentKnowledgeBaseAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanAgentKnowledgeBaseAttachmentCreate,
		ReadContext:   resourceDigitalOceanAgentKnowledgeBaseAttachmentRead,
		DeleteContext: resourceDigitalOceanAgentKnowledgeBaseAttachmentDelete,
		Schema: map[string]*schema.Schema{
			"agent_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A unique identifier for an agent.",
			},
			"knowledge_base_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A unique identifier for a knowledge base.",
			},
		},
	}
}
func resourceDigitalOceanKnowledgeBaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	// Build the create request from schema data
	createRequest := &godo.KnowledgeBaseCreateRequest{
		Name:               d.Get("name").(string),
		ProjectID:          d.Get("project_id").(string),
		Region:             d.Get("region").(string),
		EmbeddingModelUuid: d.Get("embedding_model_uuid").(string),
	}

	// Handle optional fields
	if vpcUUID, ok := d.GetOk("vpc_uuid"); ok {
		createRequest.VPCUuid = vpcUUID.(string)
	}

	if databaseID, ok := d.GetOk("database_id"); ok {
		createRequest.DatabaseID = databaseID.(string)
	}

	// Handle tags
	if tagsSet, ok := d.GetOk("tags"); ok {
		tags := make([]string, 0)
		for _, tag := range tagsSet.(*schema.Set).List() {
			tags = append(tags, tag.(string))
		}
		createRequest.Tags = tags
	}

	// Handle datasources
	if datasourcesRaw, ok := d.GetOk("datasources"); ok {
		datasources := datasourcesRaw.([]interface{})
		createRequest.DataSources = expandKnowledgeBaseDatasources(datasources)
	}

	// Make the API call
	kb, _, err := client.GenAI.CreateKnowledgeBase(ctx, createRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating knowledge base: %s", err))
	}

	d.SetId(kb.Uuid)

	// Read the created resource to populate all fields
	return resourceDigitalOceanKnowledgeBaseRead(ctx, d, meta)
}

func resourceDigitalOceanKnowledgeBaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	id := d.Id()

	kb, resp, _, err := client.GenAI.GetKnowledgeBase(ctx, id)
	if err != nil {
		if resp != "" && resp == "404" {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error retrieving Knowledge Base: %s", err))
	}

	// Set all attributes
	_ = d.Set("name", kb.Name)
	_ = d.Set("project_id", kb.ProjectId)
	_ = d.Set("region", kb.Region)
	_ = d.Set("created_at", kb.CreatedAt.UTC().String())
	_ = d.Set("database_id", kb.DatabaseId)
	_ = d.Set("embedding_model_uuid", kb.EmbeddingModelUuid)
	_ = d.Set("is_public", kb.IsPublic)
	if kb.AddedToAgentAt != nil {
		_ = d.Set("added_to_agent_at", kb.AddedToAgentAt.UTC().String())
	}

	// Get datasources separately using ListKnowledgebaseDataSources API
	datasources, _, err := client.GenAI.ListKnowledgeBaseDataSources(ctx, id, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving knowledge base datasources: %s", err))
	}

	// Flatten and set datasources if any exist
	if len(datasources) > 0 {
		flattenedDatasources := flattenKnowledgeBaseDataSources(datasources)
		if err := d.Set("datasources", flattenedDatasources); err != nil {
			return diag.FromErr(fmt.Errorf("error setting datasources: %s", err))
		}
	}

	// Set tags if they exist
	if kb.Tags != nil {
		if err := d.Set("tags", tag.FlattenTags(kb.Tags)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting tags: %s", err))
		}
	}

	return nil
}

func resourceDigitalOceanKnowledgeBaseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	id := d.Id()

	// Build the update request
	updateRequest := &godo.UpdateKnowledgeBaseRequest{
		KnowledgeBaseUUID: id,
	}
	hasChanges := false

	// Check what fields have changed and add them to the update request
	if d.HasChange("name") {
		updateRequest.Name = d.Get("name").(string)
		hasChanges = true
	}

	if d.HasChange("project_id") {
		updateRequest.ProjectID = d.Get("project_id").(string)
		hasChanges = true
	}

	if d.HasChange("embedding_model_uuid") {
		updateRequest.EmbeddingModelUuid = d.Get("embedding_model_uuid").(string)
		hasChanges = true
	}

	if d.HasChange("database_id") {
		updateRequest.DatabaseID = d.Get("database_id").(string)
		hasChanges = true
	}

	if d.HasChange("tags") {
		var tags []string
		if tagsSet, ok := d.GetOk("tags"); ok {
			for _, tag := range tagsSet.(*schema.Set).List() {
				tags = append(tags, tag.(string))
			}
		}
		updateRequest.Tags = tags
		hasChanges = true
	}

	if !hasChanges {
		return nil
	}

	_, _, err := client.GenAI.UpdateKnowledgeBase(ctx, id, updateRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating knowledge base: %s", err))
	}

	return resourceDigitalOceanKnowledgeBaseRead(ctx, d, meta)
}

func resourceDigitalOceanKnowledgeBaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	id := d.Id()

	_, resp, err := client.GenAI.DeleteKnowledgeBase(ctx, id)
	if err != nil {
		if resp != nil && resp.Response != nil && resp.StatusCode == 404 {
			return nil
		}
		return diag.FromErr(fmt.Errorf("error deleting Knowledge Base (%s): %s", id, err))
	}

	return nil
}

func resourceDigitalOceanKnowledgeBaseDatasourceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	kbUUID := d.Get("knowledge_base_uuid").(string)

	req := &godo.AddKnowledgeBaseDataSourceRequest{}

	// Handle spaces_data_source
	if spacesRaw, ok := d.GetOk("spaces_data_source"); ok {
		spacesList := spacesRaw.([]interface{})
		if len(spacesList) > 0 {
			req.SpacesDataSource = expandSpacesDataSource(spacesList)
		}
	}

	// Handle web_crawler_data_source
	if webCrawlerRaw, ok := d.GetOk("web_crawler_data_source"); ok {
		webCrawlerList := webCrawlerRaw.([]interface{})
		if len(webCrawlerList) > 0 {
			req.WebCrawlerDataSource = expandWebCrawlerDataSource(webCrawlerList)
		}
	}

	ds, _, err := client.GenAI.AddKnowledgeBaseDataSource(ctx, kbUUID, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating data source for Knowledge Base: %s", err))
	}

	d.SetId(ds.Uuid)
	return resourceDigitalOceanKnowledgeBaseDatasourceRead(ctx, d, meta)
}

func resourceDigitalOceanKnowledgeBaseDatasourceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	kbUUID := d.Get("knowledge_base_uuid").(string)
	dsUUID := d.Id()

	// List all datasources and find the one with dsUUID
	datasources, _, err := client.GenAI.ListKnowledgeBaseDataSources(ctx, kbUUID, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading data sources for Knowledge Base: %s", err))
	}

	for _, ds := range datasources {
		if ds.Uuid == dsUUID {
			flattened := flattenKnowledgeBaseDataSources([]godo.KnowledgeBaseDataSource{ds})
			if len(flattened) > 0 {
				for k, v := range flattened[0].(map[string]interface{}) {
					d.Set(k, v)
				}
			}
			return nil
		}
	}

	d.SetId("")
	return nil
}

func resourceDigitalOceanKnowledgeBaseDatasourceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	kbUUID := d.Get("knowledge_base_uuid").(string)
	dsUUID := d.Id()

	_, _, _, err := client.GenAI.DeleteKnowledgeBaseDataSource(ctx, kbUUID, dsUUID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting data source (%s) from Knowledge Base (%s): %s", dsUUID, kbUUID, err))
	}

	d.SetId("")
	return nil
}

func resourceDigitalOceanAgentKnowledgeBaseAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	agentUUID := d.Get("agent_uuid").(string)
	kbUUID := d.Get("knowledge_base_uuid").(string)

	agent, _, err := client.GenAI.AttachKnowledgeBaseToAgent(ctx, agentUUID, kbUUID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error attaching knowledge base to agent: %s", err))
	}

	d.SetId(fmt.Sprintf("%s-%s", agentUUID, kbUUID))

	flattenAgent, _ := FlattenDigitalOceanAgent(agent)
	for k, v := range flattenAgent {
		d.Set(k, v)
	}

	return resourceDigitalOceanAgentKnowledgeBaseAttachmentRead(ctx, d, meta)
}

func resourceDigitalOceanAgentKnowledgeBaseAttachmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	agentUUID := d.Get("agent_uuid").(string)

	agent, _, err := client.GenAI.GetAgent(ctx, agentUUID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading agent: %s", err))
	}

	flattenAgent, _ := FlattenDigitalOceanAgent(agent)
	for k, v := range flattenAgent {
		d.Set(k, v)
	}

	return nil
}

func resourceDigitalOceanAgentKnowledgeBaseAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	agentUUID := d.Get("agent_uuid").(string)
	kbUUID := d.Get("knowledge_base_uuid").(string)

	agent, _, err := client.GenAI.DetachKnowledgeBaseToAgent(ctx, agentUUID, kbUUID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error detaching knowledge base from agent: %s", err))
	}

	flattenAgent, _ := FlattenDigitalOceanAgent(agent)
	for k, v := range flattenAgent {
		d.Set(k, v)
	}

	d.SetId("")
	return nil
}
