package digitalocean

import (
	"context"
	"log"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDigitalOceanUptimeAlert() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanUptimeAlertCreate,
		ReadContext:   resourceDigitalOceanUptimeAlertRead,
		UpdateContext: resourceDigitalOceanUptimeAlertUpdate,
		DeleteContext: resourceDigitalOceanUptimeAlertDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "A human-friendly display name for the alert.",
				Required:    true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"check_id": {
				Type:        schema.TypeString,
				Description: "A unique identifier for a check.",
				Required:    true,
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The type of health check to perform. Enum: 'latency' 'down' 'down_global' 'ssl_expiry'",
				ValidateFunc: validation.StringInSlice([]string{
					"latency",
					"down",
					"down_global",
					"ssl_expiry",
				}, false),
				Required: true,
			},
			"threshold": {
				Type:        schema.TypeInt,
				Description: "The threshold at which the alert will enter a trigger state. The specific threshold is dependent on the alert type.",
				Optional:    true,
			},
			"comparison": {
				Type:        schema.TypeString,
				Description: "The comparison operator used against the alert's threshold. Enum: 'greater_than' 'less_than",
				ValidateFunc: validation.StringInSlice([]string{
					"greater_than",
					"less_than",
				}, false),
				Optional: true,
			},
			"period": {
				Type:        schema.TypeString,
				Description: "Period of time the threshold must be exceeded to trigger the alert. Enum '2m' '3m' '5m' '10m' '15m' '30m' '1h'",
				ValidateFunc: validation.StringInSlice([]string{
					"2m",
					"3m",
					"5m",
					"10m",
					"15m",
					"30m",
					"1hr",
				}, false),
				Optional: true,
			},
			"notifications": {
				Type:        schema.TypeList,
				Required:    true,
				Description: "The notification settings for a trigger alert.",
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
							Type:        schema.TypeList,
							Optional:    true,
							Description: "List of email addresses to sent notifications to",
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func resourceDigitalOceanUptimeAlertCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	checkID := d.Get("check_id").(string)

	opts := &godo.CreateUptimeAlertRequest{
		Name:          d.Get("name").(string),
		Type:          d.Get("type").(string),
		Notifications: expandNotifications(d.Get("notifications").([]interface{})),
		Comparison:    d.Get("comparison").(string),
		Threshold:     d.Get("threshold").(int),
		Period:        d.Get("period").(string),
	}

	log.Printf("[DEBUG] Uptime alert create configuration: %#v", opts)
	alert, _, err := client.UptimeChecks.CreateAlert(ctx, checkID, opts)
	if err != nil {
		return diag.Errorf("Error creating Uptime Alert: %s", err)
	}

	d.SetId(alert.ID)
	log.Printf("[INFO] Uptime Alert name: %s", alert.Name)

	return resourceDigitalOceanUptimeAlertRead(ctx, d, meta)
}

func expandNotifications(config []interface{}) *godo.Notifications {
	alertConfig := config[0].(map[string]interface{})
	alerts := &godo.Notifications{
		Slack: expandSlack(alertConfig["slack"].([]interface{})),
		Email: expandEmail(alertConfig["email"].([]interface{})),
	}
	return alerts
}

func resourceDigitalOceanUptimeAlertUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	checkID := d.Get("check_id").(string)

	opts := &godo.UpdateUptimeAlertRequest{
		Name:          d.Get("name").(string),
		Type:          d.Get("type").(string),
		Notifications: expandNotifications(d.Get("notifications").([]interface{})),
	}

	if v, ok := d.GetOk("comparison"); ok {
		opts.Comparison = v.(string)
	}
	if v, ok := d.GetOk("threshold"); ok {
		opts.Threshold = v.(int)
	}
	if v, ok := d.GetOk("period"); ok {
		opts.Period = v.(string)
	}

	log.Printf("[DEBUG] Uptime alert update configuration: %#v", opts)

	alert, _, err := client.UptimeChecks.UpdateAlert(ctx, checkID, d.Id(), opts)
	if err != nil {
		return diag.Errorf("Error updating Alert: %s", err)
	}

	log.Printf("[INFO] Uptime Alert name: %s", alert.Name)

	return resourceDigitalOceanUptimeAlertRead(ctx, d, meta)
}

func resourceDigitalOceanUptimeAlertDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	checkID := d.Get("check_id").(string)

	log.Printf("[INFO] Deleting uptime alert: %s", d.Id())

	// Delete the uptime alert
	_, err := client.UptimeChecks.DeleteAlert(ctx, checkID, d.Id())

	if err != nil {
		return diag.Errorf("Error deleting uptime alerts: %s", err)
	}

	return nil
}

func resourceDigitalOceanUptimeAlertRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	checkID := d.Get("check_id").(string)

	alert, resp, err := client.UptimeChecks.GetAlert(ctx, checkID, d.Id())
	if err != nil {
		// If the check is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving check: %s", err)
	}

	d.SetId(alert.ID)
	d.Set("name", alert.Name)
	d.Set("type", alert.Type)
	d.Set("threshold", alert.Threshold)
	d.Set("notifications", flattenNotifications(alert.Notifications))
	d.Set("comparison", alert.Comparison)
	d.Set("period", alert.Period)

	return nil
}

func flattenNotifications(alerts *godo.Notifications) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"email": flattenEmail(alerts.Email),
			"slack": flattenSlack(alerts.Slack),
		},
	}
}
