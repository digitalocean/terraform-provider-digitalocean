package kubernetes

import (
	"context"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDigitalOceanKubernetesCluster() *schema.Resource {
	dsNodePoolSchema := nodePoolSchema(false)

	for _, k := range dsNodePoolSchema {
		k.Computed = true
		k.Required = false
		k.Default = nil
		k.ValidateFunc = nil
	}

	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanKubernetesClusterRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"kubeconfig_expire_seconds": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntAtLeast(0),
			},

			controlPlaneFirewallField: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"allowed_addresses": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			routingAgentField: {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},

			amdGpuDevicePluginField: {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},

			amdGpuDeviceMetricsExporterPlugin: {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"surge_upgrade": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"ha": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"vpc_uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"cluster_subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"service_subnet": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ipv4_address": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tag.TagsSchema(),

			"maintenance_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"day": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"duration": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"cluster_autoscaler_configuration": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scale_down_utilization_threshold": {
							Type:     schema.TypeFloat,
							Optional: true,
						},
						"scale_down_unneeded_time": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"expanders": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"node_pool": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: dsNodePoolSchema,
				},
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"kube_config": kubernetesConfigSchema(),

			"auto_upgrade": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDigitalOceanKubernetesClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	clusters, resp, err := client.Kubernetes.List(context.Background(), &godo.ListOptions{})
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return diag.Errorf("No clusters found")
		}

		return diag.Errorf("Error listing Kubernetes clusters: %s", err)
	}

	// select the correct cluster
	for _, c := range clusters {
		if c.Name == d.Get("name").(string) {
			d.SetId(c.ID)

			return digitaloceanKubernetesClusterRead(client, c, d)
		}
	}

	return diag.Errorf("Unable to find cluster with name: %s", d.Get("name").(string))
}
