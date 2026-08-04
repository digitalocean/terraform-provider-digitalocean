package database

import (
	"context"
	"fmt"
	"log"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanDatabaseAdvancedMySQLConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDatabaseAdvancedMySQLConfigCreate,
		ReadContext:   resourceDigitalOceanDatabaseAdvancedMySQLConfigRead,
		UpdateContext: resourceDigitalOceanDatabaseAdvancedMySQLConfigUpdate,
		DeleteContext: resourceDigitalOceanDatabaseAdvancedMySQLConfigDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDigitalOceanDatabaseAdvancedMySQLConfigImport,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"mysql_parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceDigitalOceanDatabaseAdvancedMySQLConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	if err := updateAdvancedMySQLConfig(ctx, d, client); err != nil {
		return diag.Errorf("Error updating advanced MySQL configuration: %s", err)
	}

	clusterID := d.Get("cluster_id").(string)
	d.SetId(makeDatabaseAdvancedMySQLConfigID(clusterID))

	return resourceDigitalOceanDatabaseAdvancedMySQLConfigRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseAdvancedMySQLConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	if err := updateAdvancedMySQLConfig(ctx, d, client); err != nil {
		return diag.Errorf("Error updating advanced MySQL configuration: %s", err)
	}

	return resourceDigitalOceanDatabaseAdvancedMySQLConfigRead(ctx, d, meta)
}

func updateAdvancedMySQLConfig(ctx context.Context, d *schema.ResourceData, client *godo.Client) error {
	clusterID := d.Get("cluster_id").(string)

	opts := &godo.AdvancedMySQLConfigUpdate{}

	if v, ok := d.GetOk("mysql_parameters"); ok {
		opts.MySQLParameters = expandAdvancedMySQLParameters(v.(map[string]interface{}))
	}

	log.Printf("[DEBUG] Advanced MySQL configuration: %s", godo.Stringify(opts))

	if _, err := client.Databases.UpdateAdvancedMySQLConfig(ctx, clusterID, opts); err != nil {
		return err
	}

	return nil
}

func resourceDigitalOceanDatabaseAdvancedMySQLConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	config, resp, err := client.Databases.GetAdvancedMySQLConfig(ctx, clusterID)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving advanced MySQL configuration: %s", err)
	}

	if _, ok := d.GetOk("mysql_parameters"); ok {
		apiParams := advancedMySQLParametersToMap(config.MySQLParameters)
		managed := d.Get("mysql_parameters").(map[string]interface{})

		if err := d.Set("mysql_parameters", flattenAdvancedMySQLParametersForRead(managed, apiParams)); err != nil {
			return diag.Errorf("Error setting mysql_parameters: %s", err)
		}
	}

	return nil
}

func resourceDigitalOceanDatabaseAdvancedMySQLConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	warn := []diag.Diagnostic{
		{
			Severity: diag.Warning,
			Summary:  "digitalocean_database_advanced_mysql_config removed from state",
			Detail:   "Database configurations are only removed from state when destroyed. The remote configuration is not unset.",
		},
	}

	return warn
}

func resourceDigitalOceanDatabaseAdvancedMySQLConfigImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	clusterID := d.Id()

	d.SetId(makeDatabaseAdvancedMySQLConfigID(clusterID))
	d.Set("cluster_id", clusterID)

	return []*schema.ResourceData{d}, nil
}

func makeDatabaseAdvancedMySQLConfigID(clusterID string) string {
	return fmt.Sprintf("%s/advanced-mysql-config", clusterID)
}

func expandAdvancedMySQLParameters(raw map[string]interface{}) map[string]string {
	if len(raw) == 0 {
		return nil
	}

	result := make(map[string]string, len(raw))
	for k, v := range raw {
		result[k] = v.(string)
	}

	return result
}

func advancedMySQLParametersToMap(params []godo.AdvancedMySQLParameter) map[string]string {
	result := make(map[string]string, len(params))
	for _, param := range params {
		result[param.Name] = param.Value
	}

	return result
}

// flattenAdvancedMySQLParametersForRead maps managed parameters to state.
// The advanced_mysql GET endpoint may return an empty value for parameters that are
// set via PATCH; preserve the configured value in that case to avoid perpetual drift.
func flattenAdvancedMySQLParametersForRead(managed map[string]interface{}, apiParams map[string]string) map[string]string {
	updated := make(map[string]string, len(managed))

	for k, configVal := range managed {
		configStr := configVal.(string)

		if v, exists := apiParams[k]; exists && v != "" {
			updated[k] = v
			continue
		}

		updated[k] = configStr
	}

	return updated
}
