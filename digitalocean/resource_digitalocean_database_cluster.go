package digitalocean

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceDigitalOceanDatabaseCluster() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanDatabaseClusterCreate,
		Read:   resourceDigitalOceanDatabaseClusterRead,
		Update: resourceDigitalOceanDatabaseClusterUpdate,
		Delete: resourceDigitalOceanDatabaseClusterDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"engine": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"version": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"size": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"region": {
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(val interface{}) string {
					// DO API V2 region slug is always lowercase
					return strings.ToLower(val.(string))
				},
				ValidateFunc: validation.NoZeroValues,
			},

			"node_count": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"maintenance_window": {
				Type:     schema.TypeList,
				Optional: true,
				MinItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"day": {
							Type:     schema.TypeString,
							Required: true,
						},
						"hour": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},

			"eviction_policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"sql_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"private_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"uri": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"private_uri": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"database": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"user": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},

		CustomizeDiff: customdiff.All(
			schema.CustomizeDiffFunc(func(diff *schema.ResourceDiff, v interface{}) error {
				engine := diff.Get("engine")
				_, hasVersion := diff.GetOk("version")
				old, _ := diff.GetChange("version")

				if !hasVersion {
					if old != "" {
						return fmt.Errorf(`The argument "version" is now required. Set the %v version to the value saved to state: %v`, engine, old)
					}

					return fmt.Errorf(`The argument "version" is required, but no definition was found.`)
				}

				return nil
			}),
			schema.CustomizeDiffFunc(func(diff *schema.ResourceDiff, v interface{}) error {
				engine := diff.Get("engine")
				_, hasEvictionPolicy := diff.GetOk("eviction_policy")
				_, hasSqlMode := diff.GetOk("sql_mode")

				if hasSqlMode && engine != "mysql" {
					return fmt.Errorf("sql_mode is only supported for MySQL Database Clusters")
				}

				if hasEvictionPolicy && engine != "redis" {
					return fmt.Errorf("eviction_policy is only supported for Redis Database Clusters")
				}

				return nil
			}),
		),
	}
}

func resourceDigitalOceanDatabaseClusterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	opts := &godo.DatabaseCreateRequest{
		Name:       d.Get("name").(string),
		EngineSlug: d.Get("engine").(string),
		Version:    d.Get("version").(string),
		SizeSlug:   d.Get("size").(string),
		Region:     d.Get("region").(string),
		NumNodes:   d.Get("node_count").(int),
		Tags:       expandTags(d.Get("tags").(*schema.Set).List()),
	}

	log.Printf("[DEBUG] DatabaseCluster create configuration: %#v", opts)
	database, _, err := client.Databases.Create(context.Background(), opts)
	if err != nil {
		return fmt.Errorf("Error creating DatabaseCluster: %s", err)
	}

	database, err = waitForDatabaseCluster(client, database.ID, "online")
	if err != nil {
		return fmt.Errorf("Error creating DatabaseCluster: %s", err)
	}

	d.SetId(database.ID)
	log.Printf("[INFO] DatabaseCluster Name: %s", database.Name)

	if v, ok := d.GetOk("maintenance_window"); ok {
		opts := expandMaintWindowOpts(v.([]interface{}))

		resp, err := client.Databases.UpdateMaintenance(context.Background(), d.Id(), opts)
		if err != nil {
			// If the database is somehow already destroyed, mark as
			// successfully gone
			if resp != nil && resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}

			return fmt.Errorf("Error adding maintenance window for DatabaseCluster: %s", err)
		}
	}

	if policy, ok := d.GetOk("eviction_policy"); ok {
		_, err := client.Databases.SetEvictionPolicy(context.Background(), d.Id(), policy.(string))
		if err != nil {
			return fmt.Errorf("Error adding eviction policy for DatabaseCluster: %s", err)
		}
	}

	if mode, ok := d.GetOk("sql_mode"); ok {
		_, err := client.Databases.SetSQLMode(context.Background(), d.Id(), mode.(string))
		if err != nil {
			return fmt.Errorf("Error adding SQL mode for DatabaseCluster: %s", err)
		}
	}

	return resourceDigitalOceanDatabaseClusterRead(d, meta)
}

func resourceDigitalOceanDatabaseClusterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	if d.HasChange("size") || d.HasChange("node_count") {
		opts := &godo.DatabaseResizeRequest{
			SizeSlug: d.Get("size").(string),
			NumNodes: d.Get("node_count").(int),
		}

		resp, err := client.Databases.Resize(context.Background(), d.Id(), opts)
		if err != nil {
			// If the database is somehow already destroyed, mark as
			// successfully gone
			if resp != nil && resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}

			return fmt.Errorf("Error resizing DatabaseCluster: %s", err)
		}

		_, err = waitForDatabaseCluster(client, d.Id(), "online")
		if err != nil {
			return fmt.Errorf("Error resizing DatabaseCluster: %s", err)
		}
	}

	if d.HasChange("region") {
		opts := &godo.DatabaseMigrateRequest{
			Region: d.Get("region").(string),
		}

		resp, err := client.Databases.Migrate(context.Background(), d.Id(), opts)
		if err != nil {
			// If the database is somehow already destroyed, mark as
			// successfully gone
			if resp != nil && resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}

			return fmt.Errorf("Error migrating DatabaseCluster: %s", err)
		}

		_, err = waitForDatabaseCluster(client, d.Id(), "online")
		if err != nil {
			return fmt.Errorf("Error migrating DatabaseCluster: %s", err)
		}
	}

	if d.HasChange("maintenance_window") {
		opts := expandMaintWindowOpts(d.Get("maintenance_window").([]interface{}))

		resp, err := client.Databases.UpdateMaintenance(context.Background(), d.Id(), opts)
		if err != nil {
			// If the database is somehow already destroyed, mark as
			// successfully gone
			if resp != nil && resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}

			return fmt.Errorf("Error updating maintenance window for DatabaseCluster: %s", err)
		}
	}

	if d.HasChange("eviction_policy") {
		if policy, ok := d.GetOk("eviction_policy"); ok {
			_, err := client.Databases.SetEvictionPolicy(context.Background(), d.Id(), policy.(string))
			if err != nil {
				return fmt.Errorf("Error updating eviction policy for DatabaseCluster: %s", err)
			}
		} else {
			// If the eviction policy is completely removed from the config, set to noeviction
			_, err := client.Databases.SetEvictionPolicy(context.Background(), d.Id(), godo.EvictionPolicyNoEviction)
			if err != nil {
				return fmt.Errorf("Error updating eviction policy for DatabaseCluster: %s", err)
			}
		}
	}

	if d.HasChange("sql_mode") {
		_, err := client.Databases.SetSQLMode(context.Background(), d.Id(), d.Get("sql_mode").(string))
		if err != nil {
			return fmt.Errorf("Error updating SQL mode for DatabaseCluster: %s", err)
		}
	}

	if d.HasChange("tags") {
		err := setTags(client, d, godo.DatabaseResourceType)
		if err != nil {
			return fmt.Errorf("Error updating tags: %s", err)
		}
	}

	return resourceDigitalOceanDatabaseClusterRead(d, meta)
}

func resourceDigitalOceanDatabaseClusterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	database, resp, err := client.Databases.Get(context.Background(), d.Id())
	if err != nil {
		// If the database is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving DatabaseCluster: %s", err)
	}

	d.Set("name", database.Name)
	d.Set("engine", database.EngineSlug)
	d.Set("version", database.VersionSlug)
	d.Set("size", database.SizeSlug)
	d.Set("region", database.RegionSlug)
	d.Set("node_count", database.NumNodes)
	d.Set("tags", database.Tags)

	if _, ok := d.GetOk("maintenance_window"); ok {
		if err := d.Set("maintenance_window", flattenMaintWindowOpts(*database.MaintenanceWindow)); err != nil {
			return fmt.Errorf("[DEBUG] Error setting maintenance_window - error: %#v", err)
		}
	}

	if _, ok := d.GetOk("eviction_policy"); ok {
		policy, _, err := client.Databases.GetEvictionPolicy(context.Background(), d.Id())
		if err != nil {
			return fmt.Errorf("Error retrieving eviction policy for DatabaseCluster: %s", err)
		}

		d.Set("eviction_policy", policy)
	}

	if _, ok := d.GetOk("sql_mode"); ok {
		mode, _, err := client.Databases.GetSQLMode(context.Background(), d.Id())
		if err != nil {
			return fmt.Errorf("Error retrieving SQL mode for DatabaseCluster: %s", err)
		}

		d.Set("sql_mode", mode)
	}

	// Computed values
	d.Set("host", database.Connection.Host)
	d.Set("private_host", database.PrivateConnection.Host)
	d.Set("port", database.Connection.Port)
	d.Set("uri", database.Connection.URI)
	d.Set("private_uri", database.PrivateConnection.URI)
	d.Set("database", database.Connection.Database)
	d.Set("user", database.Connection.User)
	d.Set("password", database.Connection.Password)
	d.Set("urn", database.URN())

	return nil
}

func resourceDigitalOceanDatabaseClusterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	log.Printf("[INFO] Deleting DatabaseCluster: %s", d.Id())
	_, err := client.Databases.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting DatabaseCluster: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForDatabaseCluster(client *godo.Client, id string, status string) (*godo.Database, error) {
	ticker := time.NewTicker(15 * time.Second)
	timeout := 120
	n := 0

	for range ticker.C {
		database, _, err := client.Databases.Get(context.Background(), id)
		if err != nil {
			ticker.Stop()
			return nil, fmt.Errorf("Error trying to read DatabaseCluster state: %s", err)
		}

		if database.Status == status {
			ticker.Stop()
			return database, nil
		}

		if n > timeout {
			ticker.Stop()
			break
		}

		n++
	}

	return nil, fmt.Errorf("Timeout waiting to DatabaseCluster to become %s", status)
}

func expandMaintWindowOpts(config []interface{}) *godo.DatabaseUpdateMaintenanceRequest {
	maintWindowOpts := &godo.DatabaseUpdateMaintenanceRequest{}
	configMap := config[0].(map[string]interface{})

	if v, ok := configMap["day"]; ok {
		maintWindowOpts.Day = v.(string)
	}

	if v, ok := configMap["hour"]; ok {
		maintWindowOpts.Hour = v.(string)
	}

	return maintWindowOpts
}

func flattenMaintWindowOpts(opts godo.DatabaseMaintenanceWindow) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)
	item := make(map[string]interface{})

	item["day"] = opts.Day
	item["hour"] = opts.Hour
	result = append(result, item)

	return result
}
