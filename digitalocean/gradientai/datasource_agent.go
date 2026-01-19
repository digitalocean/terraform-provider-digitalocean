package gradientai

import (
	"context"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanAgent() *schema.Resource {

	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanAgentRead,
		Schema:      AgentSchemaRead(),
	}
}

func dataSourceDigitalOceanAgentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	agentID := d.Get("agent_id").(string)

	agent, _, err := client.GradientAI.GetAgent(ctx, agentID)
	if err != nil {
		return diag.FromErr(err)
	}

	flattened, err := FlattenDigitalOceanAgent(agent)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := util.SetResourceDataFromMap(d, flattened); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(agent.Uuid)
	return nil
}
