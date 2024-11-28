package dropletautoscale

import (
	"context"
	"fmt"

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
			"list_member_opts": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Pagination options for listing Droplet autoscale pool members",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"page": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Page offset",
						},
						"per_page": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Per-page count",
						},
					},
				},
			},
			"list_history_opts": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Pagination options for listing Droplet autoscale pool history events",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"page": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Page offset",
						},
						"per_page": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Per-page count",
						},
					},
				},
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
			"members": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"droplet_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Droplet ID",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Droplet create timestamp",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Droplet update timestamp",
						},
						"health_status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Droplet health status",
						},
						"unhealthy_reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Droplet unhealthy reason",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Droplet state",
						},
						"current_utilization": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"memory": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Droplet Memory utilization",
									},
									"cpu": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: "Droplet CPU utilization",
									},
								},
							},
						},
					},
				},
			},
			"history_events": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"history_event_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event ID",
						},
						"current_instance_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Event reported current member count",
						},
						"desired_instance_count": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Event reported target member count",
						},
						"reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event description",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event status",
						},
						"error_reason": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event error reason",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event create timestamp",
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Event update timestamp",
						},
					},
				},
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
			if dropletAutoscalePoolList[i] == name {
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

	if memberOpts, ok := d.GetOk("list_member_opts"); ok {
		opts := expandPaginationOpts(memberOpts.([]interface{}))
		members, _, err := client.DropletAutoscale.ListMembers(context.Background(), foundDropletAutoscalePool.ID, opts)
		if err != nil {
			return diag.Errorf("Error listing Droplet autoscale pool members: %v", err)
		}
		d.Set("members", flattenMembers(members))
	}

	if historyEventOpts, ok := d.GetOk("list_history_opts"); ok {
		opts := expandPaginationOpts(historyEventOpts.([]interface{}))
		events, _, err := client.DropletAutoscale.ListHistory(context.Background(), foundDropletAutoscalePool.ID, opts)
		if err != nil {
			return diag.Errorf("Error listing Droplet autoscale pool history events: %v", err)
		}
		d.Set("history_events", flatterHistoryEvents(events))
	}

	return nil
}

func expandPaginationOpts(opts []interface{}) *godo.ListOptions {
	if len(opts) > 0 {
		optsMap := opts[0].(map[string]interface{})
		return &godo.ListOptions{
			Page:    optsMap["page"].(int),
			PerPage: optsMap["per_page"].(int),
		}
	}
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

func flattenMembers(members []*godo.DropletAutoscaleResource) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(members))
	for _, member := range members {
		r := make(map[string]interface{})
		r["droplet_id"] = fmt.Sprint(member.DropletID)
		r["created_at"] = member.CreatedAt.UTC().String()
		r["updated_at"] = member.UpdatedAt.UTC().String()
		r["health_status"] = member.HealthStatus
		r["unhealthy_reason"] = member.UnhealthyReason
		r["status"] = member.Status
		r["current_utilization"] = flattenUtilization(member.CurrentUtilization)
		result = append(result, r)
	}
	return result
}

func flatterHistoryEvents(events []*godo.DropletAutoscaleHistoryEvent) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(events))
	for _, event := range events {
		r := make(map[string]interface{})
		r["history_event_id"] = event.HistoryEventID
		r["current_instance_count"] = int(event.CurrentInstanceCount)
		r["desired_instance_count"] = int(event.DesiredInstanceCount)
		r["reason"] = event.Reason
		r["status"] = event.Status
		r["error_reason"] = event.ErrorReason
		r["created_at"] = event.CreatedAt.UTC().String()
		r["updated_at"] = event.UpdatedAt.UTC().String()
		result = append(result, r)
	}
	return result
}
