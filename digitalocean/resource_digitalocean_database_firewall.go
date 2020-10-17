package digitalocean

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDigitalOceanDatabaseFirewall() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDatabaseFirewallCreate,
		ReadContext:   resourceDigitalOceanDatabaseFirewallRead,
		UpdateContext: resourceDigitalOceanDatabaseFirewallUpdate,
		DeleteContext: resourceDigitalOceanDatabaseFirewallDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDigitalOceanDatabaseFirewallImport,
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"rule": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"ip_addr",
								"droplet",
								"k8s",
								"tag",
							}, false),
						},

						"value": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.NoZeroValues,
						},

						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceDigitalOceanDatabaseFirewallCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
	clusterID := d.Get("cluster_id").(string)

	rules := buildDatabaseFirewallRequest(d.Get("rule").(*schema.Set).List())

	_, err := client.Databases.UpdateFirewallRules(context.TODO(), clusterID, &rules)
	if err != nil {
		return fmt.Errorf("Error creating DatabaseFirewall: %s", err)
	}

	d.SetId(resource.PrefixedUniqueId(clusterID + "-"))

	return resourceDigitalOceanDatabaseFirewallRead(d, meta)
}

func resourceDigitalOceanDatabaseFirewallRead(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
	clusterID := d.Get("cluster_id").(string)

	rules, _, err := client.Databases.GetFirewallRules(context.TODO(), clusterID)
	if err != nil {
		return fmt.Errorf("Error retrieving DatabaseFirewall: %s", err)
	}

	return d.Set("rule", flattenDatabaseFirewallRules(rules))
}

func resourceDigitalOceanDatabaseFirewallUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
	clusterID := d.Get("cluster_id").(string)

	rules := buildDatabaseFirewallRequest(d.Get("rule").(*schema.Set).List())

	_, err := client.Databases.UpdateFirewallRules(context.TODO(), clusterID, &rules)
	if err != nil {
		return fmt.Errorf("Error updating DatabaseFirewall: %s", err)
	}

	return resourceDigitalOceanDatabaseFirewallRead(d, meta)
}

func resourceDigitalOceanDatabaseFirewallDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
	clusterID := d.Get("cluster_id").(string)

	log.Printf("[INFO] Deleting DatabaseFirewall: %s", d.Id())
	req := godo.DatabaseUpdateFirewallRulesRequest{
		Rules: []*godo.DatabaseFirewallRule{},
	}

	_, err := client.Databases.UpdateFirewallRules(context.TODO(), clusterID, &req)
	if err != nil {
		return fmt.Errorf("Error deleting DatabaseFirewall: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceDigitalOceanDatabaseFirewallImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	clusterID := d.Id()
	d.Set("cluster_id", clusterID)
	d.SetId(resource.PrefixedUniqueId(clusterID + "-"))

	return []*schema.ResourceData{d}, nil
}

func buildDatabaseFirewallRequest(rules []interface{}) godo.DatabaseUpdateFirewallRulesRequest {
	expandedRules := make([]*godo.DatabaseFirewallRule, 0, len(rules))
	for _, rawRule := range rules {
		rule := rawRule.(map[string]interface{})

		r := godo.DatabaseFirewallRule{
			Type:  rule["type"].(string),
			Value: rule["value"].(string),
		}

		if rule["uuid"].(string) != "" {
			r.UUID = rule["uuid"].(string)
		}

		expandedRules = append(expandedRules, &r)
	}

	return godo.DatabaseUpdateFirewallRulesRequest{
		Rules: expandedRules,
	}
}

func flattenDatabaseFirewallRules(rules []godo.DatabaseFirewallRule) []interface{} {
	if rules == nil {
		return nil
	}

	flattenedRules := make([]interface{}, len(rules))
	for i, rule := range rules {
		rawRule := map[string]interface{}{
			"uuid":       rule.UUID,
			"type":       rule.Type,
			"value":      rule.Value,
			"created_at": rule.CreatedAt.Format(time.RFC3339),
		}

		flattenedRules[i] = rawRule
	}

	return flattenedRules
}
