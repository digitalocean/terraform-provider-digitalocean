package dedicatedinference

import (
	"context"
	"fmt"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDigitalOceanDedicatedInference() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanDedicatedInferenceRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The ID of the dedicated inference endpoint.",
				ValidateFunc: validation.NoZeroValues,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the dedicated inference endpoint.",
			},
			"region": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The region where the dedicated inference endpoint is deployed.",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current status of the dedicated inference endpoint.",
			},
			"vpc_uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The UUID of the VPC the dedicated inference endpoint is deployed in.",
			},
			"enable_public_endpoint": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether the public HTTPS endpoint is enabled.",
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
			"model_deployments": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of model deployments running on the dedicated inference endpoint.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"model_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique ID of the model.",
						},
						"model_slug": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The slug identifier for the model.",
						},
						"model_provider": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The provider of the model.",
						},
						"accelerators": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: "The GPU accelerators allocated for this model deployment.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"accelerator_slug": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The slug identifier for the GPU accelerator type.",
									},
									"scale": {
										Type:        schema.TypeInt,
										Computed:    true,
										Description: "The number of accelerator units allocated.",
									},
									"type": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: "The accelerator type.",
									},
								},
							},
						},
					},
				},
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

func dataSourceDigitalOceanDedicatedInferenceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	id := d.Get("id").(string)

	di, _, err := client.DedicatedInference.Get(ctx, id)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error reading dedicated inference endpoint (%s): %w", id, err))
	}

	d.SetId(di.ID)
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
