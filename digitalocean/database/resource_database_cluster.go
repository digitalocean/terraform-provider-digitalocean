package database

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	mysqlDBEngineSlug = "mysql"
	redisDBEngineSlug = "redis"
)

func ResourceDigitalOceanDatabaseCluster() *schema.Resource {
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
				// When Redis clusters are forced to upgrade, this prevents attempting
				// to recreate clusters specifying the previous version in their config.
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					remoteVersion, _ := strconv.Atoi(old)
					configVersion, _ := strconv.Atoi(new)

					return d.Get("engine") == redisDBEngineSlug && remoteVersion > configVersion
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
							// Prevent a diff when seconds in response, e.g: "13:00" -> "13:00:00"
							DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
								newSplit := strings.Split(new, ":")
								oldSplit := strings.Split(old, ":")
								if len(newSplit) == 3 {
									new = strings.Join(newSplit[:2], ":")
								}
								if len(oldSplit) == 3 {
									old = strings.Join(oldSplit[:2], ":")
								}
								return old == new
							},
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

			"project_id": {
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

			"ui_host": {
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

			"ui_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"uri": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"ui_uri": {
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

			"ui_database": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"user": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ui_user": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"ui_password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tag.TagsSchema(),

			"backup_restore": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"database_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"backup_created_at": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},

			"storage_size_mib": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
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
	client := meta.(*config.CombinedConfig).GodoClient()

	opts := &godo.DatabaseCreateRequest{
		Name:       d.Get("name").(string),
		EngineSlug: d.Get("engine").(string),
		Version:    d.Get("version").(string),
		SizeSlug:   d.Get("size").(string),
		Region:     d.Get("region").(string),
		NumNodes:   d.Get("node_count").(int),
		Tags:       tag.ExpandTags(d.Get("tags").(*schema.Set).List()),
	}

	if v, ok := d.GetOk("private_network_uuid"); ok {
		opts.PrivateNetworkUUID = v.(string)
	}

	if v, ok := d.GetOk("project_id"); ok {
		opts.ProjectID = v.(string)
	}

	if v, ok := d.GetOk("backup_restore"); ok {
		opts.BackupRestore = expandBackupRestore(v.([]interface{}))
	}

	if v, ok := d.GetOk("storage_size_mib"); ok {
		v, err := strconv.ParseUint(v.(string), 10, 64)
		if err == nil {
			opts.StorageSizeMib = v
		}
	}

	log.Printf("[DEBUG] database cluster create configuration: %#v", opts)
	database, _, err := client.Databases.Create(context.Background(), opts)
	if err != nil {
		return diag.Errorf("Error creating database cluster: %s", err)
	}

	err = setDatabaseConnectionInfo(database, d)
	if err != nil {
		return diag.Errorf("Error setting connection info for database cluster: %s", err)
	}

	d.SetId(database.ID)
	log.Printf("[INFO] database cluster Name: %s", database.Name)

	_, err = waitForDatabaseCluster(client, d, "online")
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
	client := meta.(*config.CombinedConfig).GodoClient()

	if d.HasChanges("size", "node_count", "storage_size_mib") {
		opts := &godo.DatabaseResizeRequest{
			SizeSlug: d.Get("size").(string),
			NumNodes: d.Get("node_count").(int),
		}

		// only include the storage_size_mib in the resize request if it has changed
		// this avoids invalid values when plans sizes are increasing that require higher base levels of storage
		// excluding this parameter will utilize default base storage levels for the given plan size
		if v, ok := d.GetOk("storage_size_mib"); ok && d.HasChange("storage_size_mib") {
			v, err := strconv.ParseUint(v.(string), 10, 64)
			if err == nil {
				opts.StorageSizeMib = v
			}
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

	if d.HasChange("version") {
		upgradeVersionReq := &godo.UpgradeVersionRequest{Version: d.Get("version").(string)}
		_, err := client.Databases.UpgradeMajorVersion(context.Background(), d.Id(), upgradeVersionReq)
		if err != nil {
			return diag.Errorf("Error upgrading version for database cluster: %s", err)
		}
	}

	if d.HasChange("tags") {
		err := tag.SetTags(client, d, godo.DatabaseResourceType)
		if err != nil {
			return diag.Errorf("Error updating tags: %s", err)
		}
	}

	return resourceDigitalOceanDatabaseClusterRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

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
	d.Set("storage_size_mib", strconv.FormatUint(database.StorageSizeMib, 10))
	d.Set("tags", tag.FlattenTags(database.Tags))

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

	uiErr := setUIConnectionInfo(database, d)
	if uiErr != nil {
		return diag.Errorf("Error setting ui connection info for database cluster: %s", err)
	}

	d.Set("urn", database.URN())
	d.Set("private_network_uuid", database.PrivateNetworkUUID)
	d.Set("project_id", database.ProjectID)

	return nil
}

func resourceDigitalOceanDatabaseClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

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
		timeoutSeconds = d.Timeout(schema.TimeoutCreate).Seconds()
		timeout        = int(timeoutSeconds / tickerInterval.Seconds())
		n              = 0
		ticker         = time.NewTicker(tickerInterval)
	)

	for range ticker.C {
		database, resp, err := client.Databases.Get(context.Background(), d.Id())
		if resp.StatusCode == 404 {
			continue
		}

		if err != nil {
			ticker.Stop()
			return nil, fmt.Errorf("Error trying to read database cluster state: %s", err)
		}

		if database.Status == status {
			ticker.Stop()
			return database, nil
		}

		if n >= timeout {
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

		uri, err := buildDBConnectionURI(database.Connection, d)
		if err != nil {
			return err
		}

		d.Set("uri", uri)
		d.Set("database", database.Connection.Database)
		d.Set("user", database.Connection.User)
		if database.Connection.Password != "" {
			d.Set("password", database.Connection.Password)
		}
	}

	if database.PrivateConnection != nil {
		d.Set("private_host", database.PrivateConnection.Host)

		privateUri, err := buildDBPrivateURI(database.PrivateConnection, d)
		if err != nil {
			return err
		}

		d.Set("private_uri", privateUri)
	}

	return nil
}

func setUIConnectionInfo(database *godo.Database, d *schema.ResourceData) error {
	if database.UIConnection != nil {
		d.Set("ui_host", database.UIConnection.Host)
		d.Set("ui_port", database.UIConnection.Port)
		d.Set("ui_uri", database.UIConnection.URI)
		d.Set("ui_database", database.UIConnection.Database)
		d.Set("ui_user", database.UIConnection.User)
		d.Set("ui_password", database.UIConnection.Password)
	}

	return nil
}

// buildDBConnectionURI constructs a connection URI using the password stored in state.
//
// MongoDB clusters only return their password in response to the initial POST.
// The host for the cluster is not known until it becomes available. In order to
// build a usable connection URI, we must save the password and then add it to
// the URL returned latter.
//
// This also protects against the password being removed from the URI if the user
// switches to using a read-only token. All database engines redact the password
// in that case
func buildDBConnectionURI(conn *godo.DatabaseConnection, d *schema.ResourceData) (string, error) {
	password := d.Get("password")
	uri, err := url.Parse(conn.URI)
	if err != nil {
		return "", err
	}

	userInfo := url.UserPassword(conn.User, password.(string))
	uri.User = userInfo

	return uri.String(), nil
}

func buildDBPrivateURI(conn *godo.DatabaseConnection, d *schema.ResourceData) (string, error) {
	return buildDBConnectionURI(conn, d)
}

func expandBackupRestore(config []interface{}) *godo.DatabaseBackupRestore {
	backupRestoreConfig := config[0].(map[string]interface{})

	backupRestore := &godo.DatabaseBackupRestore{
		DatabaseName:    backupRestoreConfig["database_name"].(string),
		BackupCreatedAt: backupRestoreConfig["backup_created_at"].(string),
	}

	return backupRestore
}
