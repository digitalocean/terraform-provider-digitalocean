package genai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDigitalOceanGenAI() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanGenAICreate,
		ReadContext:   resourceDigitalOceanAgentRead,
		UpdateContext: resourceDigitalOceanGenAIUpdate,
		DeleteContext: resourceDigitalOceanGenAIDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"agent_uuid": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the GenAI resource.",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The region where the GenAI resource will be created.",
			},
			"faas_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The model to use for the GenAI resource.",
			},
			"faas_namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The current status of the GenAI resource.",
			},
			"function_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The creation timestamp of the GenAI resource.",
			},
			"function_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the GenAI function.",
			},
			"input_schema": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The input schema of the GenAI resource.",
			},
			"output_schema": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The output schema of the GenAI resource.",
			},
		},
	}
}

func resourceDigitalOceanGenAICreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	createRequest := &godo.FunctionRouteCreateRequest{}

	if a, err := d.Get("agent_uuid").(string); err {
		createRequest.AgentUuid = a
	} else {
		return diag.FromErr(fmt.Errorf("agent_uuid is required"))
	}

	createRequest.Description = d.Get("description").(string)
	createRequest.FaasName = d.Get("faas_name").(string)
	createRequest.FaasNamespace = d.Get("faas_namespace").(string)

	// Parse input_schema JSON string into FunctionInputSchema struct.
	inputSchemaStr := d.Get("input_schema").(string)
	var inputSchema godo.FunctionInputSchema
	if err := json.Unmarshal([]byte(inputSchemaStr), &inputSchema); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse input_schema: %s", err))
	}
	createRequest.InputSchema = inputSchema

	// Optionally validate output_schema JSON. Since it is a json.RawMessage, we'll check its validity.
	outputSchemaStr := d.Get("output_schema").(string)
	var tmp json.RawMessage
	if err := json.Unmarshal([]byte(outputSchemaStr), &tmp); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse output_schema: %s", err))
	}
	createRequest.OutputSchema = []byte(outputSchemaStr)

	createRequest.FunctionName = d.Get("function_name").(string)
	agent, resp, err := client.GenAI.CreateFunctionRoute(context.Background(), createRequest.AgentUuid, createRequest)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return diag.FromErr(fmt.Errorf("GenAI resource not found: %s", err))
		}
		return diag.FromErr(fmt.Errorf("error creating GenAI resource: %s", err))
	}

	for _, function := range agent.Functions {
		if function.Name == createRequest.FunctionName {
			fmt.Println("Function found:", function.Name)
			d.Set("function_uuid", function.Uuid)
			break
		}
	}
	d.Set("agent_uuid", agent.Uuid)
	d.SetId(agent.Uuid)
	return resourceDigitalOceanAgentRead(ctx, d, meta)
}

func resourceDigitalOceanGenAIDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	agentUuid, ok := d.GetOk("agent_uuid")
	if !ok || agentUuid.(string) == "" {
		return diag.FromErr(fmt.Errorf("agent_uuid is required for deletion"))
	}

	functionUuid, ok := d.GetOk("function_uuid")
	fmt.Println("function_uuid is ", functionUuid.(string))
	if !ok || functionUuid.(string) == "" {
		return diag.FromErr(fmt.Errorf("function_uuid is required for deletion"))
	}

	_, resp, err := client.GenAI.DeleteFunctionRoute(context.Background(), agentUuid.(string), functionUuid.(string))
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error deleting GenAI resource: %s", err))
	}

	return nil
}

func resourceDigitalOceanGenAIUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	agentUuid := d.Get("agent_uuid").(string)

	if agentUuid == "" {
		return diag.FromErr(fmt.Errorf("agent_uuid is required for deletion"))
	}

	updateRequest := &godo.FunctionRouteUpdateRequest{}
	if d.HasChange("description") {
		updateRequest.Description = d.Get("description").(string)
	}

	if d.HasChange("faas_name") {
		updateRequest.FaasName = d.Get("faas_name").(string)
	}

	if d.HasChange("faas_namespace") {
		updateRequest.FaasNamespace = d.Get("faas_namespace").(string)
	}

	if d.HasChange("input_schema") {
		inputSchemaStr := d.Get("input_schema").(string)
		var inputSchema godo.FunctionInputSchema
		if err := json.Unmarshal([]byte(inputSchemaStr), &inputSchema); err != nil {
			return diag.FromErr(fmt.Errorf("failed to parse input_schema: %s", err))
		}
		updateRequest.InputSchema = inputSchema
	}

	if d.HasChange("output_schema") {
		outputSchemaStr := d.Get("output_schema").(string)
		var tmp json.RawMessage
		if err := json.Unmarshal([]byte(outputSchemaStr), &tmp); err != nil {
			return diag.FromErr(fmt.Errorf("failed to parse output_schema: %s", err))
		}
		updateRequest.OutputSchema = []byte(outputSchemaStr)
	}

	if d.HasChange("function_name") {
		updateRequest.FunctionName = d.Get("function_name").(string)
	}

	updateRequest.FunctionUuid = d.Get("function_uuid").(string)

	agent, resp, err := client.GenAI.UpdateFunctionRoute(context.Background(), agentUuid, updateRequest.FunctionUuid, updateRequest)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return diag.FromErr(fmt.Errorf("GenAI resource not found: %s", err))
		}
		return diag.FromErr(fmt.Errorf("error updating GenAI resource: %s", err))
	}

	d.SetId(agent.Uuid)
	return nil
}
