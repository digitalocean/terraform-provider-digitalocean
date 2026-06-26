package gradientai

import (
	"context"
	"fmt"
	"net/http"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// agentGuardrailItem is a single guardrail entry in the attach request body.
type agentGuardrailItem struct {
	GuardrailUuid string `json:"guardrail_uuid"`
	Priority      int    `json:"priority"`
}

// agentGuardrailAttachBody is the body for the guardrail attach request.
type agentGuardrailAttachBody struct {
	Guardrails []agentGuardrailItem `json:"guardrails"`
}

func ResourceDigitalOceanAgentGuardrailAttachment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanAgentGuardrailAttachmentCreate,
		ReadContext:   resourceDigitalOceanAgentGuardrailAttachmentRead,
		DeleteContext: resourceDigitalOceanAgentGuardrailAttachmentDelete,
		Schema: map[string]*schema.Schema{
			"agent_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A unique identifier for an agent.",
			},
			"guardrail_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "A unique identifier for a guardrail to attach to the agent.",
			},
			"priority": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     0,
				Description: "The priority of the guardrail for the agent. Lower numbers are evaluated first.",
			},
		},
	}
}

func resourceDigitalOceanAgentGuardrailAttachmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	agentUUID := d.Get("agent_uuid").(string)
	guardrailUUID := d.Get("guardrail_uuid").(string)
	priority := d.Get("priority").(int)

	// godo has no guardrail attach method, so use the client's raw request primitives.
	path := fmt.Sprintf("/v2/gen-ai/agents/%s/guardrails", agentUUID)
	body := &agentGuardrailAttachBody{
		Guardrails: []agentGuardrailItem{{GuardrailUuid: guardrailUUID, Priority: priority}},
	}
	req, err := client.NewRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return diag.FromErr(err)
	}
	if _, err := client.Do(ctx, req, nil); err != nil {
		return diag.FromErr(fmt.Errorf("error attaching guardrail to agent: %s", err))
	}

	d.SetId(fmt.Sprintf("%s-%s", agentUUID, guardrailUUID))

	return resourceDigitalOceanAgentGuardrailAttachmentRead(ctx, d, meta)
}

func resourceDigitalOceanAgentGuardrailAttachmentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	agentUUID := d.Get("agent_uuid").(string)
	guardrailUUID := d.Get("guardrail_uuid").(string)

	agent, _, err := client.GradientAI.GetAgent(ctx, agentUUID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading agent: %s", err))
	}

	found := false
	for _, g := range agent.Guardrails {
		if g == nil {
			continue
		}
		if g.GuardrailUuid == guardrailUUID || g.Uuid == guardrailUUID {
			found = true
			_ = d.Set("priority", g.Priority)
			break
		}
	}

	if !found {
		d.SetId("")
		return nil
	}

	_ = d.Set("agent_uuid", agentUUID)
	_ = d.Set("guardrail_uuid", guardrailUUID)

	return nil
}

func resourceDigitalOceanAgentGuardrailAttachmentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	agentUUID := d.Get("agent_uuid").(string)
	guardrailUUID := d.Get("guardrail_uuid").(string)

	path := fmt.Sprintf("/v2/gen-ai/agents/%s/guardrails/%s", agentUUID, guardrailUUID)
	req, err := client.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	resp, err := client.Do(ctx, req, nil)
	if err != nil {
		if resp != nil && resp.Response != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error detaching guardrail from agent: %s", err))
	}

	d.SetId("")
	return nil
}
