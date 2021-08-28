package digitalocean

import (
	"context"

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
			// TODO: sort this list
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"uuid": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"v1/insights/droplet/load_1", "v1/insights/droplet/load_5", "v1/insights/droplet/load_15",
					"v1/insights/droplet/memory_utilization_percent", "v1/insights/droplet/disk_utilization_percent",
					"v1/insights/droplet/cpu", "v1/insights/droplet/disk_read", "v1/insights/droplet/disk_write",
					"v1/insights/droplet/public_outbound_bandwidth", "v1/insights/droplet/public_inbound_bandwidth",
					"v1/insights/droplet/private_outbound_bandwidth", "v1/insights/droplet/private_inbound_bandwidth",
				}, false),
			},

			"compare": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"GreaterThan", "LessThan",
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

			// {
			// 	"alerts": {
			// 	  "email": [
			// 		"bob@exmaple.com"
			// 	  ],
			// 	  "slack": [
			// 		{
			// 		  "channel": "Production Alerts",
			// 		  "url": "https://hooks.slack.com/services/T1234567/AAAAAAAA/ZZZZZZ"
			// 		}
			// 	  ]
			// 	},
			// "alerts": {
			// 	Type:        schema.TypeList,
			// 	Computed:    false,
			// 	Required:    false,
			// 	Description: "List with details how to notify about the alert. Support for Slack or email.",
			// },

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
		Type:        d.Get("type").(string),
		Description: d.Get("description").(string),
		Tags:        expandTags(d.Get("tags")).([]string),
		Compare:     d.Get("compare").(string),
	}

	_, err, _ = client.Monitoring.CreateAlertPolicy(ctx, alertCreateRequest)

	return nil
}

// func expandAlerts(alerts []interface{}) ([]godo.Alerts, error) {
// 	//

// 	expandedAlerts := make([]godo.Alerts, len(alerts))
// 	for i, s := range alerts {
// 		alert := s.(string)

// 		var expandedSshKey godo.DropletCreateSSHKey
// 		if id, err := strconv.Atoi(sshKey); err == nil {
// 			expandedSshKey.ID = id
// 		} else {
// 			expandedSshKey.Fingerprint = sshKey
// 		}

// 		expandedSshKeys[i] = expandedSshKey
// 	}

// 	return expandedSshKeys, nil
// }

func resourceDigitalOceanMonitorAlertUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	client.Monitoring.UpdateAlertPolicy(ctx, d)

	return nil
}

func resourceDigitalOceanMonitorAlertRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()
	// get a single monitoring alert by 
	client.Monitoring.GetAlertPolicy(ctx, , meta)

	return nil
}

func resourceDigitalOceanMonitorAlertDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	client.Monitoring.CreateAlertPolicy(ctx, d)

	return nil

	// if err != nil {
	// 	return diag.FromErr(err)
	// } else {
	// 	return nil
	// }
}
