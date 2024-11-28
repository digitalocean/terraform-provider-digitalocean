package dropletautoscale

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
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
			"members": {
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
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
				Type:     schema.TypeList,
				Computed: true,
				Optional: true,
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

func expandConfig(config []interface{}) *godo.DropletAutoscaleConfiguration {
	if len(config) > 0 {
		poolConfig := config[0].(map[string]interface{})
		return &godo.DropletAutoscaleConfiguration{
			MinInstances:            uint64(poolConfig["min_instances"].(int)),
			MaxInstances:            uint64(poolConfig["max_instances"].(int)),
			TargetCPUUtilization:    poolConfig["target_cpu_utilization"].(float64),
			TargetMemoryUtilization: poolConfig["target_memory_utilization"].(float64),
			CooldownMinutes:         uint32(poolConfig["cooldown_minutes"].(int)),
			TargetNumberInstances:   uint64(poolConfig["target_number_instances"].(int)),
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

	stateConf := &retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{"provisioning"},
		Target:     []string{"active"},
		Refresh:    dropletAutoscaleRefreshFunc(client, d.Id()),
		MinTimeout: 15 * time.Second,
		Timeout:    15 * time.Minute,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("Error waiting for Droplet autoscale pool (%s) to become active: %v", pool.Name, err)
	}

	return resourceDigitalOceanDropletAutoscaleRead(ctx, d, meta)
}

func resourceDigitalOceanDropletAutoscaleRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	pool, _, err := client.DropletAutoscale.Get(context.Background(), d.Id())
	if err != nil {
		if strings.Contains(err.Error(), fmt.Sprintf("autoscale group with id %s not found", d.Id())) {
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

	if memberOpts, ok := d.GetOk("list_member_opts"); ok {
		opts := expandPaginationOpts(memberOpts.([]interface{}))
		members, _, err := client.DropletAutoscale.ListMembers(context.Background(), d.Id(), opts)
		if err != nil {
			return diag.Errorf("Error listing Droplet autoscale pool members: %v", err)
		}
		d.Set("members", flattenMembers(members))
	} else {
		d.Set("members", []map[string]interface{}{})
	}

	if historyEventOpts, ok := d.GetOk("list_history_opts"); ok {
		opts := expandPaginationOpts(historyEventOpts.([]interface{}))
		events, _, err := client.DropletAutoscale.ListHistory(context.Background(), d.Id(), opts)
		if err != nil {
			return diag.Errorf("Error listing Droplet autoscale pool history events: %v", err)
		}
		d.Set("history_events", flatterHistoryEvents(events))
	} else {
		d.Set("history_events", []map[string]interface{}{})
	}

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

	stateConf := &retry.StateChangeConf{
		Delay:      5 * time.Second,
		Pending:    []string{http.StatusText(http.StatusOK)},
		Target:     []string{http.StatusText(http.StatusNotFound)},
		Refresh:    dropletAutoscaleRefreshFunc(client, d.Id()),
		MinTimeout: 5 * time.Second,
		Timeout:    1 * time.Minute,
	}
	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("Error waiting for Droplet autoscale pool (%s) to become be deleted: %v", d.Get("name"), err)
	}

	d.SetId("")
	return nil
}

func dropletAutoscaleRefreshFunc(client *godo.Client, poolID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		// Check autoscale pool status
		pool, _, err := client.DropletAutoscale.Get(context.Background(), poolID)
		if err != nil {
			if strings.Contains(err.Error(), fmt.Sprintf("autoscale group with id %s not found", poolID)) {
				return pool, http.StatusText(http.StatusNotFound), nil
			}
			return nil, "", fmt.Errorf("Error retrieving Droplet autoscale pool: %v", err)
		}
		if pool.Status != "active" {
			return pool, pool.Status, nil
		}
		members := make([]*godo.DropletAutoscaleResource, 0)
		opts := &godo.ListOptions{
			Page:    1,
			PerPage: 100,
		}
		// Paginate through autoscale pool members and validate status
		for {
			m, resp, err := client.DropletAutoscale.ListMembers(context.Background(), poolID, opts)
			if err != nil {
				return nil, "", fmt.Errorf("Error listing Droplet autoscale pool members: %v", err)
			}
			members = append(members, m...)
			if resp.Links.IsLastPage() {
				break
			}
			page, err := resp.Links.CurrentPage()
			if err != nil {
				break
			}
			opts.Page = page + 1
		}
		// Scan through the list to find a non-active provision state
		for i := range members {
			if members[i].Status != "active" {
				return members, members[i].Status, nil
			}
		}
		return members, "active", nil
	}
}
