package genai

import (
	"context"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanKnowledgeBase() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanKnowledgeBasesRead,
		Schema:      KnowledgeBaseSchemaRead(),
	}
}

func DataSourceDigitalOceanKnowledgeBaseDatasources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanKnowledgeBaseDatasourcesRead,
		Schema: map[string]*schema.Schema{
			"knowledge_base_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "UUID of the Knowledge Base",
			},
			"datasources": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of data sources for the Knowledge Base",
				Elem:        knowledgeBaseDatasourcesSchema(),
			},
		},
	}
}

func dataSourceDigitalOceanKnowledgeBasesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	kbIDRaw, ok := d.GetOk("uuid")
	if !ok || kbIDRaw == nil {
		return diag.Errorf("uuid must be provided")
	}
	kbID := kbIDRaw.(string)

	kb, _, _, err := client.GenAI.GetKnowledgeBase(ctx, kbID)
	if err != nil {
		return diag.FromErr(err)
	}

	flattened, err := FlattenDigitalOceanKnowledgeBase(kb)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := util.SetResourceDataFromMap(d, flattened); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(kb.Uuid)
	return nil
}

func dataSourceDigitalOceanKnowledgeBaseDatasourcesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	kbUUIDRaw, ok := d.GetOk("knowledge_base_uuid")
	if !ok || kbUUIDRaw == nil {
		return diag.Errorf("knowledge_base_uuid must be provided")
	}
	kbUUID := kbUUIDRaw.(string)

	// Call the API to list data sources for the KB
	datasources, _, err := client.GenAI.ListKnowledgeBaseDataSources(ctx, kbUUID, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	// Flatten and set datasources
	flattened := flattenKnowledgeBaseDataSources(datasources)
	if err := d.Set("datasources", flattened); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(kbUUID)
	return nil
}
