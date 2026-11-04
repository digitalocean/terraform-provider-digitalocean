package functions

import (
	"context"
	"log"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanFunctionsTrigger() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanFunctionsTriggerCreate,
		ReadContext:   resourceDigitalOceanFunctionsTriggerRead,
		UpdateContext: resourceDigitalOceanFunctionsTriggerUpdate,
		DeleteContext: resourceDigitalOceanFunctionsTriggerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDigitalOceanFunctionsTriggerImport,
		},

		Schema: map[string]*schema.Schema{
			"namespace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The ID of the namespace the trigger belongs to.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The trigger's name, unique within the namespace.",
			},
			"function": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The name of the function the trigger invokes, including its package if any (e.g. sample/hello).",
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      triggerTypeScheduled,
				ValidateFunc: validation.StringInSlice([]string{triggerTypeScheduled}, false),
				Description:  "The trigger type. Only SCHEDULED is currently supported.",
			},
			"enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the trigger is enabled.",
			},
			"scheduled_details": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cron": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotWhiteSpace,
							Description:  "A standard cron expression for when the trigger fires, in UTC.",
						},
						"body": {
							Type:             schema.TypeString,
							Optional:         true,
							ValidateFunc:     validateTriggerBody,
							DiffSuppressFunc: suppressEquivalentJSON,
							Description:      "A JSON object passed to the function as its parameters each time the trigger fires.",
						},
					},
				},
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_run_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"next_run_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		// The update endpoint omits empty bodies, so removing a body in place
		// would silently leave the old one on the trigger. Replace instead.
		CustomizeDiff: customdiff.ForceNewIfChange("scheduled_details.0.body", func(_ context.Context, old, new, _ interface{}) bool {
			return old.(string) != "" && new.(string) == ""
		}),
	}
}

func expandScheduledDetails(raw []interface{}) (*godo.TriggerScheduledDetails, error) {
	if len(raw) == 0 || raw[0] == nil {
		return nil, nil
	}
	m := raw[0].(map[string]interface{})
	body, err := expandTriggerBody(m["body"].(string))
	if err != nil {
		return nil, err
	}
	return &godo.TriggerScheduledDetails{
		Cron: m["cron"].(string),
		Body: body,
	}, nil
}

func resourceDigitalOceanFunctionsTriggerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	namespaceID := d.Get("namespace_id").(string)
	details, err := expandScheduledDetails(d.Get("scheduled_details").([]interface{}))
	if err != nil {
		return diag.FromErr(err)
	}

	req := &godo.FunctionsTriggerCreateRequest{
		Name:             d.Get("name").(string),
		Type:             d.Get("type").(string),
		Function:         d.Get("function").(string),
		IsEnabled:        d.Get("enabled").(bool),
		ScheduledDetails: details,
	}

	log.Printf("[DEBUG] Functions trigger create request: %#v", req)
	trigger, _, err := functionsService(meta).CreateTrigger(ctx, namespaceID, req)
	if err != nil {
		return diag.Errorf("Error creating Functions trigger: %s", err)
	}

	d.SetId(makeTriggerID(namespaceID, trigger.Name))
	log.Printf("[INFO] Functions trigger created: %s", d.Id())

	return resourceDigitalOceanFunctionsTriggerRead(ctx, d, meta)
}

func resourceDigitalOceanFunctionsTriggerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	namespaceID, name, err := parseTriggerID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	trigger, resp, err := functionsService(meta).GetTrigger(ctx, namespaceID, name)
	if err != nil {
		if isNotFound(resp, err) {
			log.Printf("[WARN] Functions trigger %s not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading Functions trigger %s: %s", d.Id(), err)
	}

	if err := setFunctionsTriggerAttributes(d, namespaceID, trigger); err != nil {
		return diag.Errorf("Error setting Functions trigger attributes: %s", err)
	}

	return nil
}

func resourceDigitalOceanFunctionsTriggerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	namespaceID, name, err := parseTriggerID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChanges("enabled", "scheduled_details") {
		details, err := expandScheduledDetails(d.Get("scheduled_details").([]interface{}))
		if err != nil {
			return diag.FromErr(err)
		}

		req := &godo.FunctionsTriggerUpdateRequest{
			IsEnabled:        godo.PtrTo(d.Get("enabled").(bool)),
			ScheduledDetails: details,
		}

		log.Printf("[DEBUG] Functions trigger update request for %s: %#v", d.Id(), req)
		if _, _, err := functionsService(meta).UpdateTrigger(ctx, namespaceID, name, req); err != nil {
			return diag.Errorf("Error updating Functions trigger %s: %s", d.Id(), err)
		}
	}

	return resourceDigitalOceanFunctionsTriggerRead(ctx, d, meta)
}

func resourceDigitalOceanFunctionsTriggerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	namespaceID, name, err := parseTriggerID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting Functions trigger: %s", d.Id())
	resp, err := functionsService(meta).DeleteTrigger(ctx, namespaceID, name)
	if err != nil && !isNotFound(resp, err) {
		return diag.Errorf("Error deleting Functions trigger %s: %s", d.Id(), err)
	}

	d.SetId("")
	return nil
}

func resourceDigitalOceanFunctionsTriggerImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	namespaceID, name, err := parseTriggerID(d.Id())
	if err != nil {
		return nil, err
	}
	d.SetId(makeTriggerID(namespaceID, name))
	return []*schema.ResourceData{d}, nil
}

func setFunctionsTriggerAttributes(d *schema.ResourceData, namespaceID string, trigger *godo.FunctionsTrigger) error {
	var details []interface{}
	if trigger.ScheduledDetails != nil {
		body, err := flattenTriggerBody(trigger.ScheduledDetails.Body)
		if err != nil {
			return err
		}
		details = []interface{}{map[string]interface{}{
			"cron": trigger.ScheduledDetails.Cron,
			"body": body,
		}}
	}

	var lastRun, nextRun string
	if trigger.ScheduledRuns != nil {
		lastRun = formatTime(trigger.ScheduledRuns.LastRunAt)
		nextRun = formatTime(trigger.ScheduledRuns.NextRunAt)
	}

	attrs := map[string]interface{}{
		"namespace_id":      namespaceID,
		"name":              trigger.Name,
		"function":          trigger.Function,
		"type":              trigger.Type,
		"enabled":           trigger.IsEnabled,
		"scheduled_details": details,
		"created_at":        formatTime(trigger.CreatedAt),
		"updated_at":        formatTime(trigger.UpdatedAt),
		"last_run_at":       lastRun,
		"next_run_at":       nextRun,
	}
	for k, v := range attrs {
		if err := d.Set(k, v); err != nil {
			return err
		}
	}
	return nil
}
