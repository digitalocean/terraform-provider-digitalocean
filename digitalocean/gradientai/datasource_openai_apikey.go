package gradientai

import (
	"context"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanOpenAIApiKey() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanOpenAIApiKeyRead,
		Schema:      OpenAIApiKeySchemaRead(),
	}
}

func dataSourceDigitalOceanOpenAIApiKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	apiKeyID := d.Get("uuid").(string)

	apiKeyInfo, _, err := client.GradientAI.GetOpenAIAPIKey(ctx, apiKeyID)
	if err != nil {
		return diag.FromErr(err)
	}
	flattened, err := FlattenOpenAIApiKeyInfo(apiKeyInfo)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := util.SetResourceDataFromMap(d, flattened); err != nil {
		return diag.FromErr(err)
	}
	d.SetId(apiKeyInfo.Uuid)
	return nil
}

func dataSourceDigitalOceanAgentsByOpenAIApiKeyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	uuid := d.Get("uuid").(string)

	agents, _, err := client.GradientAI.ListAgentsByOpenAIAPIKey(ctx, uuid, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	flattened, err := FlattenDigitalOceanAgents(agents)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("agents", flattened); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(uuid)
	return nil
}
