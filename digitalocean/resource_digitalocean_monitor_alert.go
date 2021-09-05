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
				Optional: true,
			},

			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
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
				Description: "The comparison operator to use for value",
			},

			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Description of the alert policy",
			},

			"enabled": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},

			"value": {
				Type:         schema.TypeFloat,
				Required:     true,
				ValidateFunc: validation.FloatAtLeast(0),
			},

			"tags": tagsSchema(),

			"alerts": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List with details how to notify about the alert. Support for Slack or email.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"slack": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"channel": {
										Type:             schema.TypeString,
										Required:         true,
										DiffSuppressFunc: CaseSensitive,
										Description:      "The Slack channel to send alerts to",
										ValidateFunc:     validation.StringIsNotEmpty,
									},
									"url": {
										Type:             schema.TypeString,
										Required:         true,
										DiffSuppressFunc: CaseSensitive,
										Description:      "The webhook URL for Slack",
										ValidateFunc:     validation.StringIsNotEmpty,
									},
								},
							},
						},
						"email": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"entities": {
				Type:        schema.TypeSet,
				Required:    true,
				MinItems:    1,
				Description: "The droplets to apply the alert policy to",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"window": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"5m", "10m", "30m", "1h",
				}, false),
			},
		},
	}
}

func resourceDigitalOceanMonitorAlertCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	alertCreateRequest := &godo.AlertPolicyCreateRequest{
		Type:        d.Get("type").(string),
		Description: d.Get("description").(string),
		Tags:        expandTags(d.Get("tags").(*schema.Set).List()),
		Compare:     godo.AlertPolicyComp(d.Get("compare").(string)),
		Window:      d.Get("window").(string),
		Value:       float32(d.Get("value").(float64)),
		Entities:    expandEntities(d.Get("entities").(*schema.Set).List()),
	}

	alertCreateRequest.Alerts = expandAlerts(d.Get("alerts").([]interface{}))

	// alertPolicy, resp, err
	alertPolicy, _, err := client.Monitoring.CreateAlertPolicy(ctx, alertCreateRequest)

	if err != nil {
		return diag.Errorf("Error creating Alert Policy: %s", err)
	}

	d.SetId(alertPolicy.UUID)
	log.Printf("[INFO] Alert Policy created, ID: %s", d.Id())

	return resourceDigitalOceanMonitorAlertRead(ctx, d, meta)
}

func expandAlerts(config []interface{}) godo.Alerts {
	alertConfig := config[0].(map[string]interface{})
	alerts := godo.Alerts{
		Slack: expandSlack(alertConfig["slack"].([]interface{})),
		Email: expandEmail(alertConfig["email"].([]interface{})),
	}
	return alerts
}

func flattenAlerts(alerts godo.Alerts) map[string]interface{} {
	result := make(map[string]interface{})

	result["email"] = flattenEmail(alerts.Email)
	result["slack"] = flattenSlack(alerts.Slack)

	return result
}

func expandSlack(slackChannels []interface{}) []godo.SlackDetails {
	expandedSlackChannels := make([]godo.SlackDetails, 0, len(slackChannels))

	for _, slackChannel := range slackChannels {
		slack := slackChannel.(map[string]interface{})
		n := godo.SlackDetails{
			Channel: slack["channel"].(string),
			URL:     slack["url"].(string),
		}

		expandedSlackChannels = append(expandedSlackChannels, n)
	}

	return expandedSlackChannels
}

func flattenSlack(slackChannels []godo.SlackDetails) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(slackChannels))

	for _, slackChannel := range slackChannels {
		item := make(map[string]interface{})
		item["url"] = slackChannel.URL
		item["channel"] = slackChannel.Channel
		result = append(result, item)
	}

	return result
}

func expandEmail(config []interface{}) []string {
	emailList := make([]string, len(config))

	for i, v := range config {
		emailList[i] = v.(string)
	}

	return emailList
}

func flattenEmail(emails []string) []string {
	if emails == nil {
		return nil
	}

	flattenedEmails := make([]string, len(emails))
	for i, v := range emails {
		flattenedEmails[i] = v
	}

	return flattenedEmails
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

func resourceDigitalOceanMonitorAlertUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	// get the alert, and update it here
	alertPolicy, _, _ := client.Monitoring.GetAlertPolicy(context.Background(), d.Id())
	d.SetId(alertPolicy.UUID)

	updateRequest := &godo.AlertPolicyUpdateRequest{}

	// TODO
	// if d.HasChange("alerts") {
	// 	// Check the individual things to see if there are changes
	// 	updateRequest := godo.AlertDestinationUpdateRequest{
	// 		SlackWebhooks: expandSlack(d.Get("slack").([]godo.AppAlertSlackWebhook)),
	// 	}
	// 	client.Monitoring.UpdateAlertPolicy(context.Background(), "", updateRequest)
	// }

	if d.HasChange("tags") {
		updateRequest.Tags = d.Get("tags").([]string)
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
	d.Set("alerts", flattenAlerts(alert.Alerts))
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
