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

func ResourceDigitalOceanDatabaseAdvancedPostgreSQLConfig() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDatabaseAdvancedPostgreSQLConfigCreate,
		ReadContext:   resourceDigitalOceanDatabaseAdvancedPostgreSQLConfigRead,
		UpdateContext: resourceDigitalOceanDatabaseAdvancedPostgreSQLConfigUpdate,
		DeleteContext: resourceDigitalOceanDatabaseAdvancedPostgreSQLConfigDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDigitalOceanDatabaseAdvancedPostgreSQLConfigImport,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"pg_parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceDigitalOceanDatabaseAdvancedPostgreSQLConfigCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	if err := updateAdvancedPostgreSQLConfig(ctx, d, client); err != nil {
		return diag.Errorf("Error updating advanced PostgreSQL configuration: %s", err)
	}

	clusterID := d.Get("cluster_id").(string)
	d.SetId(makeDatabaseAdvancedPostgreSQLConfigID(clusterID))

	return resourceDigitalOceanDatabaseAdvancedPostgreSQLConfigRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseAdvancedPostgreSQLConfigUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	if err := updateAdvancedPostgreSQLConfig(ctx, d, client); err != nil {
		return diag.Errorf("Error updating advanced PostgreSQL configuration: %s", err)
	}

	return resourceDigitalOceanDatabaseAdvancedPostgreSQLConfigRead(ctx, d, meta)
}

func updateAdvancedPostgreSQLConfig(ctx context.Context, d *schema.ResourceData, client *godo.Client) error {
	clusterID := d.Get("cluster_id").(string)

	opts := &godo.AdvancedPostgresConfigUpdate{}

	if v, ok := d.GetOk("pg_parameters"); ok {
		opts.PGParameters = expandAdvancedPostgreSQLPGParameters(v.(map[string]interface{}))
	}

	log.Printf("[DEBUG] Advanced PostgreSQL configuration: %s", godo.Stringify(opts))

	if _, err := client.Databases.UpdateAdvancedPostgresSQLConfig(ctx, clusterID, opts); err != nil {
		return err
	}

	return nil
}

func resourceDigitalOceanDatabaseAdvancedPostgreSQLConfigRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	config, resp, err := client.Databases.GetAdvancedPostgresSQLConfig(ctx, clusterID)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving advanced PostgreSQL configuration: %s", err)
	}

	if _, ok := d.GetOk("pg_parameters"); ok {
		apiParams := advancedPostgreSQLPGParametersToMap(config.PGParameters)
		managed := d.Get("pg_parameters").(map[string]interface{})

		if err := d.Set("pg_parameters", flattenAdvancedPostgreSQLPGParametersForRead(managed, apiParams)); err != nil {
			return diag.Errorf("Error setting pg_parameters: %s", err)
		}
	}

	return nil
}

func resourceDigitalOceanDatabaseAdvancedPostgreSQLConfigDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId("")

	warn := []diag.Diagnostic{
		{
			Severity: diag.Warning,
			Summary:  "digitalocean_database_advanced_postgresql_config removed from state",
			Detail:   "Database configurations are only removed from state when destroyed. The remote configuration is not unset.",
		},
	}

	return warn
}

func resourceDigitalOceanDatabaseAdvancedPostgreSQLConfigImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	clusterID := d.Id()

	d.SetId(makeDatabaseAdvancedPostgreSQLConfigID(clusterID))
	d.Set("cluster_id", clusterID)

	return []*schema.ResourceData{d}, nil
}

func makeDatabaseAdvancedPostgreSQLConfigID(clusterID string) string {
	return fmt.Sprintf("%s/advanced-postgresql-config", clusterID)
}

func expandAdvancedPostgreSQLPGParameters(raw map[string]interface{}) map[string]string {
	if len(raw) == 0 {
		return nil
	}

	result := make(map[string]string, len(raw))
	for k, v := range raw {
		result[k] = v.(string)
	}

	return result
}

func advancedPostgreSQLPGParametersToMap(params []godo.AdvancedPostgresPGParameter) map[string]string {
	result := make(map[string]string, len(params))
	for _, param := range params {
		result[param.Name] = param.Value
	}

	return result
}

// flattenAdvancedPostgreSQLPGParametersForRead maps managed parameters to state.
// The advanced_pg GET endpoint may return an empty value for parameters that are
// set via PATCH; preserve the configured value in that case to avoid perpetual drift.
func flattenAdvancedPostgreSQLPGParametersForRead(managed map[string]interface{}, apiParams map[string]string) map[string]string {
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
