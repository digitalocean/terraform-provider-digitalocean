package digitalocean

import (
	"context"
	"log"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDigitalOceanUptimeCheck() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanUptimeCheckCreate,
		ReadContext:   resourceDigitalOceanUptimeCheckRead,
		UpdateContext: resourceDigitalOceanUptimeCheckUpdate,
		DeleteContext: resourceDigitalOceanUptimeCheckDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Description: "A human-friendly display name for the check.",
				Required:    true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
				ForceNew: true,
			},
			"regions": {
				Type:        schema.TypeSet,
				Description: "An array containing the selected regions to perform healthchecks from.",
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"type": {
				Type:        schema.TypeString,
				Description: "The type of health check to perform. Enum: 'ping' 'http' 'https'",
				ValidateFunc: validation.StringInSlice([]string{
					"ping",
					"http",
					"https",
				}, false),
				Default:  "https",
				Optional: true,
			},
			"target": {
				Type:        schema.TypeString,
				Description: "The endpoint to perform healthchecks on.",
				Required:    true,
			},
			"enabled": {
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
				Description: "A boolean value indicating whether the check is enabled/disabled.",
			},
		},
	}
}

func resourceDigitalOceanUptimeCheckCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	opts := &godo.CreateUptimeCheckRequest{
		Name:   d.Get("name").(string),
		Target: d.Get("target").(string),
	}

	if v, ok := d.GetOk("type"); ok {
		opts.Type = v.(string)
	}
	if v, ok := d.GetOk("enabled"); ok {
		opts.Enabled = v.(bool)
	}
	if v, ok := d.GetOk("regions"); ok {
		expandedRegions := expandRegions(v.(*schema.Set).List())
		opts.Regions = expandedRegions
	}

	log.Printf("[DEBUG] Uptime check create configuration: %#v", opts)
	check, _, err := client.UptimeChecks.Create(ctx, opts)
	if err != nil {
		return diag.Errorf("Error creating Check: %s", err)
	}

	d.SetId(check.ID)
	log.Printf("[INFO] Uptime Check name: %s", check.Name)

	return resourceDigitalOceanUptimeCheckRead(ctx, d, meta)
}

func expandRegions(regions []interface{}) []string {
	var expandedRegions []string
	for _, r := range regions {
		expandedRegions = append(expandedRegions, r.(string))
	}
	return expandedRegions
}

func resourceDigitalOceanUptimeCheckUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	id := d.Id()

	opts := &godo.UpdateUptimeCheckRequest{
		Name:   d.Get("name").(string),
		Target: d.Get("target").(string),
	}

	if v, ok := d.GetOk("type"); ok {
		opts.Type = v.(string)
	}
	if v, ok := d.GetOk("enabled"); ok {
		opts.Enabled = v.(bool)
	}
	if v, ok := d.GetOk("regions"); ok {
		expandedRegions := expandRegions(v.(*schema.Set).List())
		opts.Regions = expandedRegions
	}

	log.Printf("[DEBUG] Uptime Check update configuration: %#v", opts)

	_, _, err := client.UptimeChecks.Update(ctx, id, opts)
	if err != nil {
		return diag.Errorf("Error updating uptime check: %s", err)
	}

	return resourceDigitalOceanUptimeCheckRead(ctx, d, meta)
}

func resourceDigitalOceanUptimeCheckDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	log.Printf("[INFO] Deleting uptime check: %s", d.Id())

	// Delete the uptime check
	_, err := client.UptimeChecks.Delete(ctx, d.Id())

	if err != nil {
		return diag.Errorf("Error deleting uptime checks: %s", err)
	}

	return nil
}

func resourceDigitalOceanUptimeCheckRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	check, resp, err := client.UptimeChecks.Get(context.Background(), d.Id())
	if err != nil {
		// If the check is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving check: %s", err)
	}

	d.SetId(check.ID)
	d.Set("name", check.Name)
	d.Set("target", check.Target)
	d.Set("enabled", check.Enabled)

	if err := d.Set("regions", flattenRegions(check.Regions)); err != nil {
		return diag.Errorf("[DEBUG] Error setting Uptime Check's regions - error: %#v", err)
	}

	if v := check.Type; v != "" {
		d.Set("type", v)
	}

	return nil
}

func flattenRegions(list []string) *schema.Set {
	flatSet := schema.NewSet(schema.HashString, []interface{}{})
	for _, v := range list {
		flatSet.Add(v)
	}
	return flatSet
}
