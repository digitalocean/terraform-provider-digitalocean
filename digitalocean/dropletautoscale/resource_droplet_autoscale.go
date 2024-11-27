package dropletautoscale

import (
	"context"
	"net/http"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanDropletAutoscale() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDropletAutoscaleCreate,
		ReadContext:   resourceDigitalOceanDropletAutoscaleRead,
		UpdateContext: resourceDigitalOceanDropletAutoscaleUpdate,
		DeleteContext: resourceDigitalOceanDropletAutoscaleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the Droplet autoscale pool",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Droplet autoscale pool",
			},
			"config": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_instances": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Min number of members",
							ValidateFunc: validation.All(validation.NoZeroValues),
						},
						"max_instances": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Max number of members",
							ValidateFunc: validation.All(validation.NoZeroValues),
						},
						"target_cpu_utilization": {
							Type:         schema.TypeFloat,
							Optional:     true,
							Description:  "CPU target threshold",
							ValidateFunc: validation.All(validation.FloatBetween(0, 1)),
						},
						"target_memory_utilization": {
							Type:         schema.TypeFloat,
							Optional:     true,
							Description:  "Memory target threshold",
							ValidateFunc: validation.All(validation.FloatBetween(0, 1)),
						},
						"cooldown_minutes": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Cooldown duration",
							ValidateFunc: validation.All(validation.NoZeroValues),
						},
						"target_number_instances": {
							Type:         schema.TypeInt,
							Optional:     true,
							Description:  "Target number of members",
							ValidateFunc: validation.All(validation.NoZeroValues),
						},
					},
				},
			},
			"droplet_template": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Droplet size",
						},
						"region": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Droplet region",
						},
						"image": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Droplet image",
						},
						"tags": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Optional:    true,
							Description: "Droplet tags",
						},
						"ssh_keys": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Required:    true,
							Description: "Droplet SSH keys",
						},
						"vpc_uuid": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Droplet VPC UUID",
						},
						"with_droplet_agent": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Enable droplet agent",
						},
						"project_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Droplet project ID",
						},
						"ipv6": {
							Type:        schema.TypeBool,
							Optional:    true,
							Description: "Enable droplet IPv6",
						},
						"user_data": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Droplet user data",
						},
					},
				},
			},
			"current_utilization": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"memory": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Average Memory utilization",
						},
						"cpu": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Average CPU utilization",
						},
					},
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Droplet autoscale pool status",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Droplet autoscale pool create timestamp",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Droplet autoscale pool update timestamp",
			},
		},
	}
}

func expandConfig(config []interface{}) *godo.DropletAutoscaleConfiguration {
	if len(config) > 0 {
		poolConfig := config[0].(map[string]interface{})
		return &godo.DropletAutoscaleConfiguration{
			MinInstances:            poolConfig["min_instances"].(uint64),
			MaxInstances:            poolConfig["max_instances"].(uint64),
			TargetCPUUtilization:    poolConfig["target_cpu_utilization"].(float64),
			TargetMemoryUtilization: poolConfig["target_memory_utilization"].(float64),
			CooldownMinutes:         poolConfig["cooldown_minutes"].(uint32),
			TargetNumberInstances:   poolConfig["target_number_instances"].(uint64),
		}
	}
	return nil
}

func expandTemplate(template []interface{}) *godo.DropletAutoscaleResourceTemplate {
	if len(template) > 0 {
		poolTemplate := template[0].(map[string]interface{})

		var tags []string
		if v, ok := poolTemplate["tags"]; ok {
			for _, tag := range v.(*schema.Set).List() {
				tags = append(tags, tag.(string))
			}
		}

		var sshKeys []string
		if v, ok := poolTemplate["ssh_keys"]; ok {
			for _, key := range v.(*schema.Set).List() {
				sshKeys = append(sshKeys, key.(string))
			}
		}

		return &godo.DropletAutoscaleResourceTemplate{
			Size:             poolTemplate["size"].(string),
			Region:           poolTemplate["region"].(string),
			Image:            poolTemplate["image"].(string),
			Tags:             tags,
			SSHKeys:          sshKeys,
			VpcUUID:          poolTemplate["vpc_uuid"].(string),
			WithDropletAgent: poolTemplate["with_droplet_agent"].(bool),
			ProjectID:        poolTemplate["project_id"].(string),
			IPV6:             poolTemplate["ipv6"].(bool),
			UserData:         poolTemplate["user_data"].(string),
		}
	}
	return nil
}

func resourceDigitalOceanDropletAutoscaleCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	pool, _, err := client.DropletAutoscale.Create(context.Background(), &godo.DropletAutoscalePoolRequest{
		Name:            d.Get("name").(string),
		Config:          expandConfig(d.Get("config").([]interface{})),
		DropletTemplate: expandTemplate(d.Get("droplet_template").([]interface{})),
	})
	if err != nil {
		return diag.Errorf("Error creating Droplet autoscale pool: %v", err)
	}

	d.SetId(pool.ID)
	return resourceDigitalOceanDropletAutoscaleRead(ctx, d, meta)
}

func resourceDigitalOceanDropletAutoscaleRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	pool, resp, err := client.DropletAutoscale.Get(context.Background(), d.Id())
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error retrieving Droplet autoscale pool: %v", err)
	}

	d.Set("name", pool.Name)
	d.Set("config", flattenConfig(pool.Config))
	d.Set("droplet_template", flattenTemplate(pool.DropletTemplate))
	d.Set("current_utilization", flattenUtilization(pool.CurrentUtilization))
	d.Set("status", pool.Status)
	d.Set("created_at", pool.CreatedAt.UTC().String())
	d.Set("updated_at", pool.UpdatedAt.UTC().String())

	return nil
}

func resourceDigitalOceanDropletAutoscaleUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	pool, _, err := client.DropletAutoscale.Update(context.Background(), d.Id(), &godo.DropletAutoscalePoolRequest{
		Name:            d.Get("name").(string),
		Config:          expandConfig(d.Get("config").([]interface{})),
		DropletTemplate: expandTemplate(d.Get("droplet_template").([]interface{})),
	})
	if err != nil {
		return diag.Errorf("Error updating Droplet autoscale pool: %v", err)
	}

	d.SetId(pool.ID)
	return resourceDigitalOceanDropletAutoscaleRead(ctx, d, meta)
}

func resourceDigitalOceanDropletAutoscaleDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	_, err := client.DropletAutoscale.DeleteDangerous(context.Background(), d.Id())
	if err != nil {
		return diag.Errorf("Error updating Droplet autoscale pool: %v", err)
	}

	d.SetId("")
	return nil
}
