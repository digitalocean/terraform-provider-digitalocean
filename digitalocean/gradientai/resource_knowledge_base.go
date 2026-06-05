package gradientai

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// knowledgeBaseRootResponse mirrors the godo (unexported) knowledgebaseRoot so the
// provider can issue raw GenAI requests via client.NewRequest/client.Do without
// modifying the vendored godo SDK.
type knowledgeBaseRootResponse struct {
	KnowledgeBase *godo.KnowledgeBase `json:"knowledge_base"`
}

// agentRootResponse mirrors the godo (unexported) agent root for raw requests.
type agentRootResponse struct {
	Agent *godo.Agent `json:"agent"`
}

// agentKnowledgeBaseAttachBody is the JSON body required by
// POST /v2/gen-ai/agents/{agent_uuid}/knowledge_bases/{kb_uuid}.
type agentKnowledgeBaseAttachBody struct {
	AgentUuid         string `json:"agent_uuid"`
	KnowledgeBaseUuid string `json:"knowledge_base_uuid"`
}

func ResourceDigitalOceanKnowledgeBase() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanKnowledgeBaseCreate,
		ReadContext:   resourceDigitalOceanKnowledgeBaseRead,
		UpdateContext: resourceDigitalOceanKnowledgeBaseUpdate,
		DeleteContext: resourceDigitalOceanKnowledgeBaseDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
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
				Optional:    true,
				ForceNew:    true,
				Description: "Data sources for the knowledge base. Omit for an empty knowledge base; add data later with digitalocean_gradientai_knowledge_base_data_source.",
				Elem:        knowledgeBaseDatasourcesSchema(),
			},
			"wait_for_database": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "When true (default), waits for the knowledge base's managed database to become ONLINE before completing creation. This is required for an agent to attach the knowledge base (e.g. inline knowledge_base_uuid on agent create, which fails while the database is still provisioning).",
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},
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
			"wait_for_database": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     true,
				Description: "When true (default), waits for the knowledge base's managed database to become ONLINE before attaching. The attach fails while the database is still provisioning. Indexing does not need to be complete.",
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

	// Handle tags
	if tagsSet, ok := d.GetOk("tags"); ok {
		tags := make([]string, 0)
		for _, tag := range tagsSet.(*schema.Set).List() {
			tags = append(tags, tag.(string))
		}
		createRequest.Tags = tags
	}

	// Data sources are optional at create time (empty KB); use [] so JSON is [], not null.
	createRequest.DataSources = []godo.KnowledgeBaseDataSource{}
	if datasourcesRaw, ok := d.GetOk("datasources"); ok {
		if dsList := datasourcesRaw.([]interface{}); len(dsList) > 0 {
			createRequest.DataSources = expandKnowledgeBaseDatasources(dsList)
		}
	}

	var kb *godo.KnowledgeBase
	if len(createRequest.DataSources) == 0 {
		// godo's CreateKnowledgeBase rejects an empty datasource list client-side, but the
		// API supports it. Issue the request directly through the godo client's exported
		// primitives so the vendored SDK does not need to change.
		var err error
		kb, err = createKnowledgeBaseRaw(ctx, client, createRequest)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error creating knowledge base: %s", err))
		}
	} else {
		var err error
		kb, _, err = client.GradientAI.CreateKnowledgeBase(ctx, createRequest)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error creating knowledge base: %s", err))
		}
	}

	d.SetId(kb.Uuid)

	// The managed database provisions asynchronously (~4 min). Block until it is ONLINE so
	// the knowledge base is actually usable (e.g. attachable to an agent) once created.
	if d.Get("wait_for_database").(bool) {
		if err := waitKnowledgeBaseDatabaseOnline(ctx, client, kb.Uuid, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.FromErr(err)
		}
	}

	// Read the created resource to populate all fields
	return resourceDigitalOceanKnowledgeBaseRead(ctx, d, meta)
}

// createKnowledgeBaseRaw POSTs a knowledge base create request directly via the godo
// client, bypassing godo's client-side "at least one datasource" guard. It mirrors the
// validation godo performs for the fields the API requires.
func createKnowledgeBaseRaw(ctx context.Context, client *godo.Client, createRequest *godo.KnowledgeBaseCreateRequest) (*godo.KnowledgeBase, error) {
	if createRequest.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if strings.Contains(createRequest.Name, " ") {
		return nil, fmt.Errorf("name cannot contain spaces")
	}
	if createRequest.Region == "" {
		createRequest.Region = "tor1"
	}
	if createRequest.EmbeddingModelUuid == "" {
		return nil, fmt.Errorf("embedding_model_uuid is required")
	}
	if createRequest.ProjectID == "" {
		return nil, fmt.Errorf("project_id is required")
	}
	if createRequest.DataSources == nil {
		createRequest.DataSources = []godo.KnowledgeBaseDataSource{}
	}

	req, err := client.NewRequest(ctx, http.MethodPost, godo.KnowledgeBasePath, createRequest)
	if err != nil {
		return nil, err
	}
	root := new(knowledgeBaseRootResponse)
	if _, err := client.Do(ctx, req, root); err != nil {
		return nil, err
	}
	if root.KnowledgeBase == nil {
		return nil, fmt.Errorf("knowledge base create returned an empty response")
	}
	return root.KnowledgeBase, nil
}

func resourceDigitalOceanKnowledgeBaseRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	id := d.Id()

	kb, resp, _, err := client.GradientAI.GetKnowledgeBase(ctx, id)
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
	_ = d.Set("created_at", kb.CreatedAt)
	_ = d.Set("database_id", kb.DatabaseId)
	_ = d.Set("embedding_model_uuid", kb.EmbeddingModelUuid)
	_ = d.Set("is_public", kb.IsPublic)
	_ = d.Set("added_to_agent_at", kb.AddedToAgentAt)

	// Get datasources separately using ListKnowledgebaseDataSources API
	datasources, _, err := client.GradientAI.ListKnowledgeBaseDataSources(ctx, id, nil)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving knowledge base datasources: %s", err))
	}

	// Flatten and set datasources (including empty)
	flattenedDatasources := flattenKnowledgeBaseDataSources(datasources)
	if err := d.Set("datasources", flattenedDatasources); err != nil {
		return diag.FromErr(fmt.Errorf("error setting datasources: %s", err))
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

	_, _, err := client.GradientAI.UpdateKnowledgeBase(ctx, id, updateRequest)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating knowledge base: %s", err))
	}

	return resourceDigitalOceanKnowledgeBaseRead(ctx, d, meta)
}

func resourceDigitalOceanKnowledgeBaseDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	id := d.Id()

	_, resp, err := client.GradientAI.DeleteKnowledgeBase(ctx, id)
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

	ds, _, err := client.GradientAI.AddKnowledgeBaseDataSource(ctx, kbUUID, req)
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
	datasources, _, err := client.GradientAI.ListKnowledgeBaseDataSources(ctx, kbUUID, nil)
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

	_, _, _, err := client.GradientAI.DeleteKnowledgeBaseDataSource(ctx, kbUUID, dsUUID)
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

	if d.Get("wait_for_database").(bool) {
		if err := waitKnowledgeBaseDatabaseOnline(ctx, client, kbUUID, d.Timeout(schema.TimeoutCreate)); err != nil {
			return diag.FromErr(err)
		}
	}

	agent, err := attachKnowledgeBaseToAgentRaw(ctx, client, agentUUID, kbUUID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error attaching knowledge base to agent: %s", err))
	}

	d.SetId(fmt.Sprintf("%s-%s", agentUUID, kbUUID))

	if agent != nil {
		flattenAgent, _ := FlattenDigitalOceanAgent(agent)
		for k, v := range flattenAgent {
			d.Set(k, v)
		}
	}

	return resourceDigitalOceanAgentKnowledgeBaseAttachmentRead(ctx, d, meta)
}

// attachKnowledgeBaseToAgentRaw issues the KB->agent attach POST with the body the API
// requires ({"agent_uuid","knowledge_base_uuid"}), via the godo client's exported
// primitives so the vendored godo SDK (whose AttachKnowledgeBaseToAgent sends no body)
// does not need to change.
func attachKnowledgeBaseToAgentRaw(ctx context.Context, client *godo.Client, agentUUID, kbUUID string) (*godo.Agent, error) {
	path := fmt.Sprintf(godo.AgentKnowledgeBasePath, agentUUID, kbUUID)
	body := &agentKnowledgeBaseAttachBody{
		AgentUuid:         agentUUID,
		KnowledgeBaseUuid: kbUUID,
	}
	req, err := client.NewRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	root := new(agentRootResponse)
	if _, err := client.Do(ctx, req, root); err != nil {
		return nil, err
	}
	return root.Agent, nil
}

func resourceDigitalOceanAgentKnowledgeBaseAttachmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	agentUUID := d.Get("agent_uuid").(string)

	agent, _, err := client.GradientAI.GetAgent(ctx, agentUUID)
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

	agent, _, err := client.GradientAI.DetachKnowledgeBaseToAgent(ctx, agentUUID, kbUUID)
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

// waitKnowledgeBaseDatabaseOnline blocks until the knowledge base's managed database
// reports ONLINE. Empirically this (not indexing completion) is the precondition for
// attaching a KB to an agent: while the database is still provisioning the attach
// returns 400 "failed to update agent config cache after linking knowledge base".
// Uses the provider's idiomatic retry.RetryContext pattern (see
// dedicatedinference.waitForDedicatedInferenceReady).
func waitKnowledgeBaseDatabaseOnline(ctx context.Context, client *godo.Client, kbUUID string, timeout time.Duration) error {
	return retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		_, dbStatus, _, err := client.GradientAI.GetKnowledgeBase(ctx, kbUUID)
		if err != nil {
			return retry.NonRetryableError(fmt.Errorf("get knowledge base (%s) while waiting for database: %w", kbUUID, err))
		}
		switch strings.ToUpper(dbStatus) {
		case "ONLINE":
			return nil
		case "DECOMMISSIONED", "UNHEALTHY":
			return retry.NonRetryableError(fmt.Errorf("knowledge base (%s) database entered terminal state %q", kbUUID, dbStatus))
		default:
			return retry.RetryableError(fmt.Errorf("knowledge base (%s) database is %q, waiting for ONLINE", kbUUID, dbStatus))
		}
	})
}
