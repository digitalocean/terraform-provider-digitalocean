package functions

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDigitalOceanFunctionsTrigger() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanFunctionsTriggerRead,
		Schema: map[string]*schema.Schema{
			"namespace_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The ID of the namespace the trigger belongs to.",
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
				Description:  "The trigger's name.",
			},
			"function": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"scheduled_details": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cron": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"body": {
							Type:     schema.TypeString,
							Computed: true,
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
	}
}

func dataSourceDigitalOceanFunctionsTriggerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	namespaceID := d.Get("namespace_id").(string)
	name := d.Get("name").(string)

	trigger, _, err := functionsService(meta).GetTrigger(ctx, namespaceID, name)
	if err != nil {
		return diag.Errorf("Error retrieving Functions trigger %s in namespace %s: %s", name, namespaceID, err)
	}

	d.SetId(makeTriggerID(namespaceID, trigger.Name))
	if err := setFunctionsTriggerAttributes(d, namespaceID, trigger); err != nil {
		return diag.Errorf("Error setting Functions trigger attributes: %s", err)
	}
	return nil
}
