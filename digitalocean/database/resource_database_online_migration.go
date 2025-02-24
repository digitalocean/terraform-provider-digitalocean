package database

import (
	"context"
	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"time"
)

func ResourceDigitalOceanDatabaseOnlineMigration() *schema.Resource {
	return &schema.Resource{
		ReadContext:   resourceDigitalOceanDatabaseOnlineMigrationStatus,
		CreateContext: resourceDigitalOceanDatabaseOnlineMigrationStart,
		UpdateContext: resourceDigitalOceanDatabaseOnlineMigrationStart,
		DeleteContext: resourceDigitalOceanDatabaseOnlineMigrationStop,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"disable_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Disables SSL encryption when connecting to the source database",
			},
			"ignore_dbs": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "The list of databases to be ignored during the migration",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the migration",
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the online migration",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time when the online migration was created",
			},
			"source": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The FQDN pointing to the database cluster's current primary node",
						},
						"port": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "The port on which the database cluster is listening",
						},
						"db_name": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The name of the default database",
						},
						"username": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The default user of the database",
						},
						"password": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "The port on which the database cluster is listening.",
						},
					},
				},
			},
		},
	}
}

func resourceDigitalOceanDatabaseOnlineMigrationStart(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	opts := &godo.DatabaseStartOnlineMigrationRequest{}

	if v, ok := d.GetOk("disable_ssl"); ok {
		opts.DisableSSL = v.(bool)
	}

	if v, ok := d.GetOk("ignore_dbs"); ok {
		var ignoreDBs []string
		for _, db := range v.(*schema.Set).List() {
			ignoreDBs = append(ignoreDBs, db.(string))
		}

		opts.IgnoreDBs = ignoreDBs
	}

	if v, ok := d.GetOk("source"); ok {
		opts.Source = expandDBOnlineMigrationSource(v.([]interface{}))
	}

	migrationID, onlineMigrationStatus := waitForOnlineMigration(ctx, client, d, clusterID, opts)
	if onlineMigrationStatus != nil {
		return onlineMigrationStatus
	}

	d.SetId(migrationID)

	return resourceDigitalOceanDatabaseOnlineMigrationStatus(ctx, d, meta)
}

// Polls for errors in migration for 60 seconds. Requests can pass the API precheck and returns 200 response but still fail the migration quickly.
// Should notify user in this scenario.
func waitForOnlineMigration(ctx context.Context, client *godo.Client, d *schema.ResourceData, clusterID string, opts *godo.DatabaseStartOnlineMigrationRequest) (string, diag.Diagnostics) {
	_, _, err := client.Databases.Get(ctx, clusterID)
	if err != nil {
		return "", diag.Errorf("Cluster does not exist: %s", clusterID)
	}

	time.Sleep(30 * time.Second)

	_, _, err = client.Databases.StartOnlineMigration(ctx, clusterID, opts)
	if err != nil {
		return "", diag.Errorf("Error here: %s", clusterID)
	}

	tickerInterval := 10 //10s
	timeoutSeconds := 90
	n := 0
	ticker := time.NewTicker(time.Duration(tickerInterval) * time.Second)

	for range ticker.C {
		if n*tickerInterval > timeoutSeconds {
			ticker.Stop()
			break
		}
		status, _, _ := client.Databases.GetOnlineMigrationStatus(ctx, clusterID)
		// if status is nil, online migration might not have kicked off yet.
		if status == nil {
			continue
		} else if status.Status == "error" {
			// try again, maybe database wasn't ready for connection.
			_, _, err = client.Databases.StartOnlineMigration(ctx, clusterID, opts)
			if err != nil {
				return "", diag.Errorf("Error starting online migration for cluster: %s", clusterID)
			}
		} else if status.Status == "syncing" || status.Status == "done" {
			// if status is syncing, online-migration was a success and can notify user
			ticker.Stop()
			return status.ID, nil
		}
		n++
	}
	// if status never reaches syncing after one minute, report failure.
	return "", diag.Errorf("Error starting online migration for cluster: %s", clusterID)
}

func expandDBOnlineMigrationSource(config []interface{}) *godo.DatabaseOnlineMigrationConfig {
	source := &godo.DatabaseOnlineMigrationConfig{}
	if len(config) == 0 || config[0] == nil {
		return source
	}
	configMap := config[0].(map[string]interface{})
	if v, ok := configMap["host"]; ok {
		source.Host = v.(string)
	}
	if v, ok := configMap["port"]; ok {
		source.Port = v.(int)
	}
	if v, ok := configMap["db_name"]; ok {
		source.DatabaseName = v.(string)
	}
	if v, ok := configMap["username"]; ok {
		source.Username = v.(string)
	}
	if v, ok := configMap["password"]; ok {
		source.Password = v.(string)
	}
	return source
}

func resourceDigitalOceanDatabaseOnlineMigrationStatus(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	onlineMigration, resp, err := client.Databases.GetOnlineMigrationStatus(ctx, clusterID)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving database online migration status: %s", err)
	}

	d.SetId(onlineMigration.ID)
	d.Set("status", onlineMigration.Status)
	d.Set("created_at", onlineMigration.CreatedAt)

	return nil
}

func resourceDigitalOceanDatabaseOnlineMigrationStop(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)
	migrationID := d.Get("id").(string)

	_, err := client.Databases.StopOnlineMigration(ctx, clusterID, migrationID)
	if err != nil {
		return diag.Errorf("Error stopping online migration: %s", err)
	}

	d.SetId("")
	return nil
}
