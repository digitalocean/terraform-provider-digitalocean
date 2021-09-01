package digitalocean

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	mongoDBEngineSlug = "mongodb"
	mysqlDBEngineSlug = "mysql"
	redisDBEngineSlug = "redis"
)

func resourceDigitalOceanDatabaseCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDatabaseClusterCreate,
		ReadContext:   resourceDigitalOceanDatabaseClusterRead,
		UpdateContext: resourceDigitalOceanDatabaseClusterUpdate,
		DeleteContext: resourceDigitalOceanDatabaseClusterDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Type: schema.TypeString,
				// TODO: Finalize transition to being required.
				// In practice, this is already required. The transitionVersionToRequired
				// CustomizeDiffFunc is used to provide users with a better hint in the error message.
				// Required: true,
				Optional: true,
				ForceNew: true,
				// Redis clusters are being force upgraded from version 5 to 6.
				// Prevent attempting to recreate clusters specifying 5 in their config.
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return d.Get("engine") == redisDBEngineSlug && old == "6" && new == "5"
				},
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

			"private_network_uuid": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ValidateFunc: validation.NoZeroValues,
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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: customdiff.All(
			transitionVersionToRequired(),
			validateExclusiveAttributes(),
		),
	}
}

func transitionVersionToRequired() schema.CustomizeDiffFunc {
	return schema.CustomizeDiffFunc(func(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
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
	})
}

func validateExclusiveAttributes() schema.CustomizeDiffFunc {
	return schema.CustomizeDiffFunc(func(ctx context.Context, diff *schema.ResourceDiff, v interface{}) error {
		engine := diff.Get("engine")
		_, hasEvictionPolicy := diff.GetOk("eviction_policy")
		_, hasSqlMode := diff.GetOk("sql_mode")

		if hasSqlMode && engine != mysqlDBEngineSlug {
			return fmt.Errorf("sql_mode is only supported for MySQL Database Clusters")
		}

		if hasEvictionPolicy && engine != redisDBEngineSlug {
			return fmt.Errorf("eviction_policy is only supported for Redis Database Clusters")
		}

		return nil
	})
}

func resourceDigitalOceanDatabaseClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if v, ok := d.GetOk("private_network_uuid"); ok {
		opts.PrivateNetworkUUID = v.(string)
	}

	log.Printf("[DEBUG] database cluster create configuration: %#v", opts)
	database, _, err := client.Databases.Create(context.Background(), opts)
	if err != nil {
		return diag.Errorf("Error creating database cluster: %s", err)
	}

	// MongoDB clusters only return the password in response to the initial POST.
	// We need to set it here before any subsequent GETs.
	if database.EngineSlug == mongoDBEngineSlug {
		err = setDatabaseConnectionInfo(database, d)
		if err != nil {
			return diag.Errorf("Error setting connection info for database cluster: %s", err)
		}
	}

	d.SetId(database.ID)
	log.Printf("[INFO] database cluster Name: %s", database.Name)

	database, err = waitForDatabaseCluster(client, d, "online")
	if err != nil {
		d.SetId("")
		return diag.Errorf("Error creating database cluster: %s", err)
	}

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

			return diag.Errorf("Error adding maintenance window for database cluster: %s", err)
		}
	}

	if policy, ok := d.GetOk("eviction_policy"); ok {
		_, err := client.Databases.SetEvictionPolicy(context.Background(), d.Id(), policy.(string))
		if err != nil {
			return diag.Errorf("Error adding eviction policy for database cluster: %s", err)
		}
	}

	if mode, ok := d.GetOk("sql_mode"); ok {
		_, err := client.Databases.SetSQLMode(context.Background(), d.Id(), mode.(string))
		if err != nil {
			return diag.Errorf("Error adding SQL mode for database cluster: %s", err)
		}
	}

	return resourceDigitalOceanDatabaseClusterRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	if d.HasChanges("size", "node_count") {
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

			return diag.Errorf("Error resizing database cluster: %s", err)
		}

		_, err = waitForDatabaseCluster(client, d, "online")
		if err != nil {
			return diag.Errorf("Error resizing database cluster: %s", err)
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

			return diag.Errorf("Error migrating database cluster: %s", err)
		}

		_, err = waitForDatabaseCluster(client, d, "online")
		if err != nil {
			return diag.Errorf("Error migrating database cluster: %s", err)
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

			return diag.Errorf("Error updating maintenance window for database cluster: %s", err)
		}
	}

	if d.HasChange("eviction_policy") {
		if policy, ok := d.GetOk("eviction_policy"); ok {
			_, err := client.Databases.SetEvictionPolicy(context.Background(), d.Id(), policy.(string))
			if err != nil {
				return diag.Errorf("Error updating eviction policy for database cluster: %s", err)
			}
		} else {
			// If the eviction policy is completely removed from the config, set to noeviction
			_, err := client.Databases.SetEvictionPolicy(context.Background(), d.Id(), godo.EvictionPolicyNoEviction)
			if err != nil {
				return diag.Errorf("Error updating eviction policy for database cluster: %s", err)
			}
		}
	}

	if d.HasChange("sql_mode") {
		_, err := client.Databases.SetSQLMode(context.Background(), d.Id(), d.Get("sql_mode").(string))
		if err != nil {
			return diag.Errorf("Error updating SQL mode for database cluster: %s", err)
		}
	}

	if d.HasChange("tags") {
		err := setTags(client, d, godo.DatabaseResourceType)
		if err != nil {
			return diag.Errorf("Error updating tags: %s", err)
		}
	}

	return resourceDigitalOceanDatabaseClusterRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	database, resp, err := client.Databases.Get(context.Background(), d.Id())
	if err != nil {
		// If the database is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving database cluster: %s", err)
	}

	d.Set("name", database.Name)
	d.Set("engine", database.EngineSlug)
	d.Set("version", database.VersionSlug)
	d.Set("size", database.SizeSlug)
	d.Set("region", database.RegionSlug)
	d.Set("node_count", database.NumNodes)
	d.Set("tags", flattenTags(database.Tags))

	if _, ok := d.GetOk("maintenance_window"); ok {
		if err := d.Set("maintenance_window", flattenMaintWindowOpts(*database.MaintenanceWindow)); err != nil {
			return diag.Errorf("[DEBUG] Error setting maintenance_window - error: %#v", err)
		}
	}

	if _, ok := d.GetOk("eviction_policy"); ok {
		policy, _, err := client.Databases.GetEvictionPolicy(context.Background(), d.Id())
		if err != nil {
			return diag.Errorf("Error retrieving eviction policy for database cluster: %s", err)
		}

		d.Set("eviction_policy", policy)
	}

	if _, ok := d.GetOk("sql_mode"); ok {
		mode, _, err := client.Databases.GetSQLMode(context.Background(), d.Id())
		if err != nil {
			return diag.Errorf("Error retrieving SQL mode for database cluster: %s", err)
		}

		d.Set("sql_mode", mode)
	}

	// Computed values
	err = setDatabaseConnectionInfo(database, d)
	if err != nil {
		return diag.Errorf("Error setting connection info for database cluster: %s", err)
	}
	d.Set("urn", database.URN())
	d.Set("private_network_uuid", database.PrivateNetworkUUID)

	return nil
}

func resourceDigitalOceanDatabaseClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	log.Printf("[INFO] Deleting database cluster: %s", d.Id())
	_, err := client.Databases.Delete(context.Background(), d.Id())
	if err != nil {
		return diag.Errorf("Error deleting database cluster: %s", err)
	}

	d.SetId("")
	return nil
}

func waitForDatabaseCluster(client *godo.Client, d *schema.ResourceData, status string) (*godo.Database, error) {
	var (
		tickerInterval = 15 * time.Second
		timeoutSeconds = d.Timeout(schema.TimeoutDelete).Seconds()
		timeout        = int(timeoutSeconds / tickerInterval.Seconds())
		n              = 0
		ticker         = time.NewTicker(tickerInterval)
	)

	for range ticker.C {
		database, _, err := client.Databases.Get(context.Background(), d.Id())
		if err != nil {
			ticker.Stop()
			return nil, fmt.Errorf("Error trying to read database cluster state: %s", err)
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

	return nil, fmt.Errorf("Timeout waiting to database cluster to become %s", status)
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

func setDatabaseConnectionInfo(database *godo.Database, d *schema.ResourceData) error {
	if database.Connection != nil {
		d.Set("host", database.Connection.Host)
		d.Set("port", database.Connection.Port)
		d.Set("uri", database.Connection.URI)
		d.Set("database", database.Connection.Database)
		d.Set("user", database.Connection.User)
		if database.EngineSlug == mongoDBEngineSlug {
			if database.Connection.Password != "" {
				d.Set("password", database.Connection.Password)
			}
			uri, err := buildMongoDBConnectionURI(database.Connection, d)
			if err != nil {
				return err
			}
			d.Set("uri", uri)
		} else {
			d.Set("password", database.Connection.Password)
			d.Set("uri", database.Connection.URI)
		}
	}

	if database.PrivateConnection != nil {
		d.Set("private_host", database.PrivateConnection.Host)
		if database.EngineSlug == mongoDBEngineSlug {
			uri, err := buildMongoDBConnectionURI(database.PrivateConnection, d)
			if err != nil {
				return err
			}
			d.Set("private_uri", uri)
		} else {
			d.Set("private_uri", database.PrivateConnection.URI)
		}
	}

	return nil
}

// MongoDB clusters only return their password in response to the initial POST.
// The host for the cluster is not known until it becomes available. In order to
// build a usable connection URI, we must save the password and then add it to
// the URL returned latter.
func buildMongoDBConnectionURI(conn *godo.DatabaseConnection, d *schema.ResourceData) (string, error) {
	password := d.Get("password")
	uri, err := url.Parse(conn.URI)
	if err != nil {
		return "", err
	}

	userInfo := url.UserPassword(conn.User, password.(string))
	uri.User = userInfo

	return uri.String(), nil
}
