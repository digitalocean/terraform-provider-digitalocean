package gradientai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDigitalOceanGradientAIFunctionRoute() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanGradientAIFunctionRouteCreate,
		ReadContext:   resourceDigitalOceanGradientAIFunctionRouteRead,
		UpdateContext: resourceDigitalOceanGradientAIFunctionRouteUpdate,
		DeleteContext: resourceDigitalOceanGradientAIFunctionRouteDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"agent_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the GradientAI resource.",
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The region where the GradientAI resource will be created.",
			},
			"faas_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The model to use for the GradientAI resource.",
			},
			"faas_namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The current status of the GradientAI resource.",
			},
			"function_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The creation timestamp of the GradientAI resource.",
			},
			"function_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique identifier of the GradientAI function.",
			},
			// currently keeping input_schema and output_schema as strings due to no specific schema definition in documentation
			"input_schema": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The input schema of the GradientAI resource.",
			},
			"output_schema": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The output schema of the GradientAI resource.",
			},
		},
	}
}

func resourceDigitalOceanGradientAIFunctionRouteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	createRequest := &godo.FunctionRouteCreateRequest{
		AgentUuid:     d.Get("agent_id").(string),
		Description:   d.Get("description").(string),
		FaasName:      d.Get("faas_name").(string),
		FaasNamespace: d.Get("faas_namespace").(string),
		FunctionName:  d.Get("function_name").(string),
	}

	// Parse input_schema JSON string into FunctionInputSchema struct.
	inputSchemaStr := d.Get("input_schema").(string)
	var inputSchema godo.FunctionInputSchema
	if err := json.Unmarshal([]byte(inputSchemaStr), &inputSchema); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse input_schema: %s", err))
	}
	createRequest.InputSchema = inputSchema

	// Optionally validate output_schema JSON. Since it is a json.RawMessage, we'll check its validity.
	outputSchemaStr := d.Get("output_schema").(string)
	if outputSchemaStr != "" {
		var tmp json.RawMessage
		if err := json.Unmarshal([]byte(outputSchemaStr), &tmp); err != nil {
			return diag.FromErr(fmt.Errorf("failed to parse output_schema: %s", err))
		}
		createRequest.OutputSchema = []byte(outputSchemaStr)
	}

	agent, resp, err := client.GradientAI.CreateFunctionRoute(context.Background(), createRequest.AgentUuid, createRequest)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return diag.FromErr(fmt.Errorf("GradientAI resource not found: %s", err))
		}
		return diag.FromErr(fmt.Errorf("error creating GradientAI resource: %s", err))
	}

	d.SetId(agent.Uuid)
	d.Set("agent_id", agent.Uuid)
	return resourceDigitalOceanGradientAIFunctionRouteRead(ctx, d, meta)
}

func resourceDigitalOceanGradientAIFunctionRouteDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	agentUuid := d.Id()
	if agentUuid == "" {
		return diag.FromErr(fmt.Errorf("agent_uuid is required for deletion"))
	}

	functionUuid, ok := d.GetOk("function_uuid")
	fmt.Println("function_uuid is ", functionUuid.(string))
	if !ok || functionUuid.(string) == "" {
		return diag.FromErr(fmt.Errorf("function_uuid is required for deletion"))
	}

	_, resp, err := client.GradientAI.DeleteFunctionRoute(context.Background(), agentUuid, functionUuid.(string))
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error deleting GradientAI resource: %s", err))
	}

	return nil
}

func resourceDigitalOceanGradientAIFunctionRouteUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	agentUuid := d.Get("agent_id").(string)

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

	inputSchemaStr := d.Get("input_schema").(string)
	var inputSchema godo.FunctionInputSchema
	if err := json.Unmarshal([]byte(inputSchemaStr), &inputSchema); err != nil {
		return diag.FromErr(fmt.Errorf("failed to parse input_schema: %s", err))
	}
	updateRequest.InputSchema = inputSchema

	outputSchemaStr := d.Get("output_schema").(string)
	if outputSchemaStr != "" {
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

	agent, resp, err := client.GradientAI.UpdateFunctionRoute(context.Background(), agentUuid, updateRequest.FunctionUuid, updateRequest)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return diag.FromErr(fmt.Errorf("GradientAI resource not found: %s", err))
		}
		return diag.FromErr(fmt.Errorf("error updating GradientAI resource: %s", err))
	}

	d.SetId(agent.Uuid)
	return nil
}

func resourceDigitalOceanGradientAIFunctionRouteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	agentUuid := d.Get("agent_id").(string)
	if agentUuid == "" {
		return diag.FromErr(fmt.Errorf("agent_id is required for reading"))
	}

	agent, resp, err := client.GradientAI.GetAgent(context.Background(), agentUuid)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error reading GradientAI resource: %s", err))
	}

	functionName := d.Get("function_name").(string)
	var foundFn *godo.AgentFunction
	for _, fn := range agent.Functions {
		if fn.Name == functionName {
			foundFn = fn
			break
		}
	}

	if foundFn == nil {
		// If not found, the function route may have been deleted remotely.
		d.SetId("")
		return nil
	}

	d.SetId(agent.Uuid)
	d.Set("agent_uuid", agent.Uuid)
	d.Set("description", foundFn.Description)
	d.Set("faas_name", foundFn.FaasName)
	d.Set("faas_namespace", foundFn.FaasNamespace)
	d.Set("function_name", foundFn.Name)
	d.Set("function_uuid", foundFn.Uuid)
	return nil
}
