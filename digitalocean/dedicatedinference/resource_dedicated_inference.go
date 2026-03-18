package dedicatedinference

import (
	"context"
	"fmt"
	"net/http"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanDedicatedInference() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDedicatedInferenceCreate,
		ReadContext:   resourceDigitalOceanDedicatedInferenceRead,
		UpdateContext: resourceDigitalOceanDedicatedInferenceUpdate,
		DeleteContext: resourceDigitalOceanDedicatedInferenceDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "A human-readable name for the dedicated inference endpoint.",
				ValidateFunc: validation.NoZeroValues,
			},
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The region slug where the dedicated inference endpoint will be deployed.",
				ValidateFunc: validation.NoZeroValues,
			},
			"enable_public_endpoint": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether to enable a public HTTPS endpoint for the dedicated inference endpoint.",
			},
			"vpc_uuid": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "The UUID of the VPC to deploy the dedicated inference endpoint into.",
				ValidateFunc: validation.IsUUID,
			},
			"model_deployments": {
				Type:        schema.TypeList,
				Required:    true,
				MinItems:    1,
				Description: "The list of model deployments to run on the dedicated inference endpoint.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"model_slug": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The slug identifier for the model to deploy.",
							ValidateFunc: validation.NoZeroValues,
						},
						"model_provider": {
							Type:         schema.TypeString,
							Required:     true,
							Description:  "The provider of the model (e.g. 'digitalocean', 'huggingface').",
							ValidateFunc: validation.NoZeroValues,
						},
						"model_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "The unique ID of the model.",
						},
						"accelerators": {
							Type:        schema.TypeList,
							Required:    true,
							MinItems:    1,
							Description: "The GPU accelerators to allocate for this model deployment.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"accelerator_slug": {
										Type:         schema.TypeString,
										Required:     true,
										Description:  "The slug identifier for the GPU accelerator type.",
										ValidateFunc: validation.NoZeroValues,
									},
									"scale": {
										Type:         schema.TypeInt,
										Required:     true,
										Description:  "The number of accelerator units to allocate.",
										ValidateFunc: validation.IntAtLeast(1),
									},
									"type": {
										Type:         schema.TypeString,
										Required:     true,
										Description:  "The accelerator type.",
										ValidateFunc: validation.NoZeroValues,
									},
								},
							},
						},
					},
				},
			},
			"hugging_face_token": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "A HuggingFace token for accessing gated models.",
			},

			// Computed fields
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the dedicated inference endpoint.",
			},
			"public_endpoint_fqdn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The fully-qualified domain name of the public endpoint, if enabled.",
			},
			"private_endpoint_fqdn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The fully-qualified domain name of the private endpoint.",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time when the dedicated inference endpoint was created.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time when the dedicated inference endpoint was last updated.",
			},
		},
	}
}

func resourceDigitalOceanDedicatedInferenceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	spec := buildDedicatedInferenceSpec(d)
	req := &godo.DedicatedInferenceCreateRequest{
		Spec: spec,
	}

	if v, ok := d.GetOk("hugging_face_token"); ok {
		req.Secrets = &godo.DedicatedInferenceSecrets{
			HuggingFaceToken: v.(string),
		}
	}

	di, _, _, err := client.DedicatedInference.Create(ctx, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating dedicated inference endpoint: %w", err))
	}

	d.SetId(di.ID)
	return resourceDigitalOceanDedicatedInferenceRead(ctx, d, meta)
}

func resourceDigitalOceanDedicatedInferenceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	di, resp, err := client.DedicatedInference.Get(ctx, d.Id())
	if err != nil {
		if resp != nil && resp.Response != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.FromErr(fmt.Errorf("error reading dedicated inference endpoint (%s): %w", d.Id(), err))
	}

	d.Set("name", di.Name)
	d.Set("region", di.Region)
	d.Set("status", di.Status)

	if di.Endpoints != nil {
		d.Set("public_endpoint_fqdn", di.Endpoints.PublicEndpointFQDN)
		d.Set("private_endpoint_fqdn", di.Endpoints.PrivateEndpointFQDN)
	}

	if !di.CreatedAt.IsZero() {
		d.Set("created_at", di.CreatedAt.UTC().String())
	}
	if !di.UpdatedAt.IsZero() {
		d.Set("updated_at", di.UpdatedAt.UTC().String())
	}

	if di.DeploymentSpec != nil {
		d.Set("enable_public_endpoint", di.DeploymentSpec.EnablePublicEndpoint)

		if di.DeploymentSpec.VPCConfig != nil {
			d.Set("vpc_uuid", di.DeploymentSpec.VPCConfig.VPCUUID)
		}

		if err := d.Set("model_deployments", flattenModelDeployments(di.DeploymentSpec.ModelDeployments)); err != nil {
			return diag.FromErr(fmt.Errorf("error setting model_deployments: %w", err))
		}
	}

	return nil
}

func resourceDigitalOceanDedicatedInferenceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	if !d.HasChanges("name", "enable_public_endpoint", "model_deployments", "hugging_face_token") {
		return nil
	}

	spec := buildDedicatedInferenceSpec(d)
	req := &godo.DedicatedInferenceUpdateRequest{
		Spec: spec,
	}

	if v, ok := d.GetOk("hugging_face_token"); ok {
		req.Secrets = &godo.DedicatedInferenceSecrets{
			HuggingFaceToken: v.(string),
		}
	}

	_, _, err := client.DedicatedInference.Update(ctx, d.Id(), req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error updating dedicated inference endpoint (%s): %w", d.Id(), err))
	}

	return resourceDigitalOceanDedicatedInferenceRead(ctx, d, meta)
}

func resourceDigitalOceanDedicatedInferenceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	resp, err := client.DedicatedInference.Delete(ctx, d.Id())
	if err != nil {
		if resp != nil && resp.Response != nil && resp.StatusCode == http.StatusNotFound {
			return nil
		}
		return diag.FromErr(fmt.Errorf("error deleting dedicated inference endpoint (%s): %w", d.Id(), err))
	}

	d.SetId("")
	return nil
}

func buildDedicatedInferenceSpec(d *schema.ResourceData) *godo.DedicatedInferenceSpecRequest {
	spec := &godo.DedicatedInferenceSpecRequest{
		Name:                 d.Get("name").(string),
		Region:               d.Get("region").(string),
		EnablePublicEndpoint: d.Get("enable_public_endpoint").(bool),
	}

	if v, ok := d.GetOk("vpc_uuid"); ok {
		spec.VPC = &godo.DedicatedInferenceVPCRequest{
			UUID: v.(string),
		}
	}

	if v, ok := d.GetOk("model_deployments"); ok {
		spec.ModelDeployments = expandModelDeployments(v.([]interface{}))
	}

	return spec
}

func expandModelDeployments(raw []interface{}) []*godo.DedicatedInferenceModelRequest {
	deployments := make([]*godo.DedicatedInferenceModelRequest, 0, len(raw))
	for _, item := range raw {
		m := item.(map[string]interface{})
		deployment := &godo.DedicatedInferenceModelRequest{
			ModelSlug:     m["model_slug"].(string),
			ModelProvider: m["model_provider"].(string),
		}
		if v, ok := m["model_id"].(string); ok && v != "" {
			deployment.ModelID = v
		}
		if v, ok := m["accelerators"].([]interface{}); ok {
			deployment.Accelerators = expandAccelerators(v)
		}
		deployments = append(deployments, deployment)
	}
	return deployments
}

func expandAccelerators(raw []interface{}) []*godo.DedicatedInferenceAcceleratorRequest {
	accelerators := make([]*godo.DedicatedInferenceAcceleratorRequest, 0, len(raw))
	for _, item := range raw {
		a := item.(map[string]interface{})
		accelerators = append(accelerators, &godo.DedicatedInferenceAcceleratorRequest{
			AcceleratorSlug: a["accelerator_slug"].(string),
			Scale:           uint64(a["scale"].(int)),
			Type:            a["type"].(string),
		})
	}
	return accelerators
}

func flattenModelDeployments(deployments []*godo.DedicatedInferenceModelDeployment) []interface{} {
	result := make([]interface{}, 0, len(deployments))
	for _, d := range deployments {
		m := map[string]interface{}{
			"model_id":       d.ModelID,
			"model_slug":     d.ModelSlug,
			"model_provider": d.ModelProvider,
			"accelerators":   flattenAccelerators(d.Accelerators),
		}
		result = append(result, m)
	}
	return result
}

func flattenAccelerators(accelerators []*godo.DedicatedInferenceAccelerator) []interface{} {
	result := make([]interface{}, 0, len(accelerators))
	for _, a := range accelerators {
		result = append(result, map[string]interface{}{
			"accelerator_slug": a.AcceleratorSlug,
			"scale":            int(a.Scale),
			"type":             a.Type,
		})
	}
	return result
}
