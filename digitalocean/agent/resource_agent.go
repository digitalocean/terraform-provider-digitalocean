package agent

import (
	"context"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDigitalOceanAgent() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanAgentCreate,
		ReadContext:   resourceDigitalOceanAgentRead,
		UpdateContext: resourceDigitalOceanAgentUpdate,
		DeleteContext: resourceDigitalOceanAgentDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Agent",
			},
			"instruction": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Instruction for the Agent",
			},
			"model_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Model UUID of the Agent",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Project ID of the Agent",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Region where the Agent is deployed",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Description where the Agent is deployed",
			},
		},
	}
}

func resourceDigitalOceanAgentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	agentRequest := &godo.AgentCreateRequest{
		Instruction: d.Get("instruction").(string),
		Name:        d.Get("name").(string),
		ModelUuid:   d.Get("model_uuid").(string),
		Region:      d.Get("region").(string),
		ProjectId:   d.Get("project_id").(string),
		Description: d.Get("description").(string),
	}

	agent, _, err := client.GenAI.CreateAgent(ctx, agentRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(agent.Uuid)
	return resourceDigitalOceanAgentRead(ctx, d, meta)
}

func resourceDigitalOceanAgentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	agent, _, err := client.GenAI.GetAgent(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("name", agent.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("instruction", agent.Instruction); err != nil {
		return diag.FromErr(err)
	}
	if agent.Model != nil {
		if err := d.Set("model_uuid", agent.Model.Uuid); err != nil {
			return diag.FromErr(err)
		}
	}
	if err := d.Set("region", agent.Region); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("project_id", agent.ProjectId); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDigitalOceanAgentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	agentRequest := &godo.AgentUpdateRequest{
		Name:        d.Get("name").(string),
		Instruction: d.Get("instruction").(string),
		ModelUuid:   d.Get("model_uuid").(string),
		Region:      d.Get("region").(string),
		ProjectId:   d.Get("project_id").(string),
	}

	agent, _, err := client.GenAI.UpdateAgent(ctx, d.Id(), agentRequest)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(agent.Uuid)
	return resourceDigitalOceanAgentRead(ctx, d, meta)
}

func resourceDigitalOceanAgentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	_, _, err := client.GenAI.DeleteAgent(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func convertToStringSlice(input interface{}) []string {
	if input == nil {
		return nil
	}

	list := input.([]interface{})
	result := make([]string, len(list))
	for i, v := range list {
		result[i] = v.(string)
	}
	return result
}
