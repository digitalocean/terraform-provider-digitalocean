package dropletautoscale

import (
	"context"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDigitalOceanDropletAutoscale() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanDropletAutoscaleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "ID of the Droplet autoscale pool",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Name of the Droplet autoscale pool",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_instances": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Min number of members",
						},
						"max_instances": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Max number of members",
						},
						"target_cpu_utilization": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "CPU target threshold",
						},
						"target_memory_utilization": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "Memory target threshold",
						},
						"cooldown_minutes": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Cooldown duration",
						},
						"target_number_instances": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Target number of members",
						},
					},
				},
			},
			"droplet_template": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"size": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Droplet size",
						},
						"region": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Droplet region",
						},
						"image": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Droplet image",
						},
						"tags": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Droplet tags",
						},
						"ssh_keys": {
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
							Description: "Droplet SSH keys",
						},
						"vpc_uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Droplet VPC UUID",
						},
						"with_droplet_agent": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Enable droplet agent",
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Droplet project ID",
						},
						"ipv6": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Enable droplet IPv6",
						},
						"user_data": {
							Type:        schema.TypeString,
							Computed:    true,
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

func dataSourceDigitalOceanDropletAutoscaleRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	var foundDropletAutoscalePool *godo.DropletAutoscalePool
	if id, ok := d.GetOk("id"); ok {
		pool, _, err := client.DropletAutoscale.Get(context.Background(), id.(string))
		if err != nil {
			return diag.Errorf("Error retrieving Droplet autoscale pool: %v", err)
		}
		foundDropletAutoscalePool = pool
	} else if name, ok := d.GetOk("name"); ok {
		dropletAutoscalePoolList := make([]*godo.DropletAutoscalePool, 0)
		opts := &godo.ListOptions{
			Page:    1,
			PerPage: 100,
		}
		// Paginate through all active resources
		for {
			pools, resp, err := client.DropletAutoscale.List(context.Background(), opts)
			if err != nil {
				return diag.Errorf("Error listing Droplet autoscale pools: %v", err)
			}
			dropletAutoscalePoolList = append(dropletAutoscalePoolList, pools...)
			if resp.Links.IsLastPage() {
				break
			}
			page, err := resp.Links.CurrentPage()
			if err != nil {
				break
			}
			opts.Page = page + 1
		}
		// Scan through the list to find a resource name match
		for i := range dropletAutoscalePoolList {
			if dropletAutoscalePoolList[i].Name == name {
				foundDropletAutoscalePool = dropletAutoscalePoolList[i]
				break
			}
		}
	} else {
		return diag.Errorf("Need to specify either a name or an id to look up the Droplet autoscale pool")
	}
	if foundDropletAutoscalePool == nil {
		return diag.Errorf("Droplet autoscale pool not found")
	}

	d.SetId(foundDropletAutoscalePool.ID)
	d.Set("name", foundDropletAutoscalePool.Name)
	d.Set("config", flattenConfig(foundDropletAutoscalePool.Config))
	d.Set("droplet_template", flattenTemplate(foundDropletAutoscalePool.DropletTemplate))
	d.Set("current_utilization", flattenUtilization(foundDropletAutoscalePool.CurrentUtilization))
	d.Set("status", foundDropletAutoscalePool.Status)
	d.Set("created_at", foundDropletAutoscalePool.CreatedAt.UTC().String())
	d.Set("updated_at", foundDropletAutoscalePool.UpdatedAt.UTC().String())

	return nil
}

func flattenConfig(config *godo.DropletAutoscaleConfiguration) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)
	if config != nil {
		r := make(map[string]interface{})
		r["min_instances"] = config.MinInstances
		r["max_instances"] = config.MaxInstances
		r["target_cpu_utilization"] = config.TargetCPUUtilization
		r["target_memory_utilization"] = config.TargetMemoryUtilization
		r["cooldown_minutes"] = config.CooldownMinutes
		r["target_number_instances"] = config.TargetNumberInstances
		result = append(result, r)
	}
	return result
}

func flattenTemplate(template *godo.DropletAutoscaleResourceTemplate) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)
	if template != nil {
		r := make(map[string]interface{})
		r["size"] = template.Size
		r["region"] = template.Region
		r["image"] = template.Image
		r["vpc_uuid"] = template.VpcUUID
		r["with_droplet_agent"] = template.WithDropletAgent
		r["project_id"] = template.ProjectID
		r["ipv6"] = template.IPV6
		r["user_data"] = template.UserData

		tagSet := schema.NewSet(schema.HashString, []interface{}{})
		for _, tag := range template.Tags {
			tagSet.Add(tag)
		}
		r["tags"] = tagSet

		keySet := schema.NewSet(schema.HashString, []interface{}{})
		for _, key := range template.SSHKeys {
			keySet.Add(key)
		}
		r["ssh_keys"] = keySet
		result = append(result, r)
	}
	return result
}

func flattenUtilization(util *godo.DropletAutoscaleResourceUtilization) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, 1)
	if util != nil {
		r := make(map[string]interface{})
		r["memory"] = util.Memory
		r["cpu"] = util.CPU
		result = append(result, r)
	}
	return result
}
