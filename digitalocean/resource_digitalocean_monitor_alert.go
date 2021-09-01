package digitalocean

import (
	"context"
	"log"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDigitalMonitorAlert() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanMonitorAlertCreate,
		ReadContext:   resourceDigitalOceanMonitorAlertRead,
		UpdateContext: resourceDigitalOceanMonitorAlertUpdate,
		DeleteContext: resourceDigitalOceanMonitorAlertDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
				Required: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					godo.DropletCPUUtilizationPercent,
					godo.DropletMemoryUtilizationPercent,
					godo.DropletDiskUtilizationPercent,
					godo.DropletPublicOutboundBandwidthRate,
					godo.DropletDiskReadRate,
					godo.DropletDiskWriteRate,
					godo.DropletOneMinuteLoadAverage,
					godo.DropletFiveMinuteLoadAverage,
					godo.DropletFifteenMinuteLoadAverage,
					// these are available as constants ...
					"v1/insights/droplet/public_inbound_bandwidth",
					"v1/insights/droplet/private_outbound_bandwidth",
					"v1/insights/droplet/private_inbound_bandwidth",
				}, false),
			},

			"compare": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(godo.GreaterThan),
					string(godo.LessThan),
				}, false),
				Description: "Description of the alert policy",
			},

			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Description of the alert policy",
			},

			"enabled": {
				Type:     schema.TypeBool,
				Default:  true,
				Required: true,
			},

			"value": {
				Type:         schema.TypeFloat,
				Required:     true,
				ValidateFunc: validation.FloatAtLeast(0),
			},

			"tags": {
				Type:        schema.TypeList,
				Computed:    false,
				Required:    false,
				Optional:    true,
				Description: "Tags for the monitoring alert",
			},

			"alerts": {
				Type:        schema.TypeList,
				Computed:    false,
				Required:    false,
				Description: "List with details how to notify about the alert. Support for Slack or email.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"slack": {
							Type:     schema.TypeList,
							Required: false,
							Computed: false,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"channel": {
										Type:         schema.TypeList,
										Computed:     false,
										Required:     true,
										Description:  "The Slack channel to send alerts to",
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"url": {
										Type:             schema.TypeString,
										Computed:         false,
										DiffSuppressFunc: CaseSensitive,
										Description:      "The webhook URL for Slack",
										ValidateFunc:     validation.StringIsNotEmpty,
									},
								},
							},
						},
						"email": {
							Type:         schema.TypeList,
							Computed:     false,
							Required:     false,
							ValidateFunc: validation.ListOfUniqueStrings,
						},
					},
				},
			},

			"entities": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				Description: "The droplets to apply the alert policy to",
			},

			"window": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"5m", "10m", "30m", "1h",
				}, false),
				MinItems: 1,
			},
		},
	}
}

func resourceDigitalOceanMonitorAlertCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()
	alertCreateRequest := &godo.AlertPolicyCreateRequest{
		Alerts:      expandAlerts(d.Get("alerts").(*schema.Set).List()),
		Type:        d.Get("type").(string),
		Description: d.Get("description").(string),
		Tags:        expandTags(d.Get("tags").(*schema.Set).List()),
		Compare:     d.Get("compare").(godo.AlertPolicyComp),
		Window:      d.Get("window").(string),
		Value:       d.Get("value").(float32),
		Entities:    expandEntities(d.Get("entities").(*schema.Set).List()),
	}

	// alertPolicy, resp, err
	alertPolicy, _, err := client.Monitoring.CreateAlertPolicy(ctx, alertCreateRequest)

	if err != nil {
		return diag.Errorf("Error creating Alert Policy: %s", err)
	}

	d.SetId(alertPolicy.UUID)
	log.Printf("[INFO] Alert Policy created, ID: %s", d.Id())

	return resourceDigitalOceanMonitorAlertRead(ctx, d, meta)
}

func expandAlerts(alerts godo.Alerts) (*godo.Alerts, error) {
	alertPolicyNotifications := &godo.Alerts{}
	alertMap = alerts[0].(map[string]interface{})

	expandedAlerts := make([]*godo.Alerts, 0, len(alerts))
	alertPolicyNotifications.Email = alert

	return alertPolicyNotifications, nil
}

func expandEntities(config []interface{}) []string {
	alertEntities := make([]string, len(config))

	for i, v := range config {
		alertEntities[i] = v.(string)
	}

	return alertEntities
}

func flattenEntities(entities []string) *schema.Set {
	// it seems there are many functions like this in different places in the code base.
	// maybe a utility library would be better
	result := schema.NewSet(schema.HashString, []interface{}{})

	for _, entity := range entities {
		result.Add(entity)
	}
	return result
}

func flattenAlerts(alerts []godo.Alerts) []map[string]interface{} {
	return nil
}

func resourceDigitalOceanMonitorAlertUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	// get the alert, and update it here
	alertPolicy, _, _ := client.Monitoring.GetAlertPolicy(context.Background(), d.Id())
	d.SetId(alertPolicy.UUID)

	updateRequest := &godo.AlertPolicyUpdateRequest{}

	if d.HasChange("alerts") {
		client.Monitoring.UpdateAlertPolicy(context.Background(), "", updateRequest)
	}

	// what other things to check in the update process?
	client.Monitoring.UpdateAlertPolicy(ctx, "alertPolicy", nil)

	return nil
}

func resourceDigitalOceanMonitorAlertRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	alert, resp, err := client.Monitoring.GetAlertPolicy(ctx, d.Id())

	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			log.Printf("[DEBUG] Alert (%s) was not found - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading Alert: %s", err)
	}

	d.SetId(alert.UUID)
	d.Set("description", alert.Description)
	d.Set("enabled", alert.Enabled)
	d.Set("compare", alert.Compare)
	// d.Set("alerts", flattenAlerts(alert.Alerts))
	d.Set("value", alert.Value)
	d.Set("window", alert.Window)
	d.Set("entities", flattenEntities(alert.Entities))
	d.Set("tags", flattenTags(alert.Tags))
	d.Set("type", alert.Type)

	return nil
}

func resourceDigitalOceanMonitorAlertDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	log.Printf("[INFO] Deleting the monitor alert")
	_, err := client.Monitoring.DeleteAlertPolicy(context.Background(), d.Id())
	if err != nil {
		return diag.Errorf("Error deleting monitor alert: %s", err)
	}
	d.SetId("")
	return nil
}
