package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/internal/mutexkv"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var mutexKV = mutexkv.NewMutexKV()

func ResourceDigitalOceanDatabaseUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDatabaseUserCreate,
		ReadContext:   resourceDigitalOceanDatabaseUserRead,
		UpdateContext: resourceDigitalOceanDatabaseUserUpdate,
		DeleteContext: resourceDigitalOceanDatabaseUserDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDigitalOceanDatabaseUserImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"cluster_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"mysql_auth_plugin": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					godo.SQLAuthPluginNative,
					godo.SQLAuthPluginCachingSHA2,
				}, false),
				// Prevent diffs when default is used and not specificed in the config.
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old == godo.SQLAuthPluginCachingSHA2 && new == ""
				},
			},

			// Computed Properties
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceDigitalOceanDatabaseUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)

	opts := &godo.DatabaseCreateUserRequest{
		Name: d.Get("name").(string),
	}

	if v, ok := d.GetOk("mysql_auth_plugin"); ok {
		opts.MySQLSettings = &godo.DatabaseMySQLUserSettings{
			AuthPlugin: v.(string),
		}
	}

	// Prevent parallel creation of users for same cluster.
	key := fmt.Sprintf("digitalocean_database_cluster/%s/users", clusterID)
	mutexKV.Lock(key)
	defer mutexKV.Unlock(key)

	log.Printf("[DEBUG] Database User create configuration: %#v", opts)
	user, _, err := client.Databases.CreateUser(context.Background(), clusterID, opts)
	if err != nil {
		return diag.Errorf("Error creating Database User: %s", err)
	}

	d.SetId(makeDatabaseUserID(clusterID, user.Name))
	log.Printf("[INFO] Database User Name: %s", user.Name)

	// MongoDB clusters only return the password in response to the initial POST.
	// So we need to set it here before any subsequent GETs.
	d.Set("password", user.Password)

	return resourceDigitalOceanDatabaseUserRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	// Check if the database user still exists
	user, resp, err := client.Databases.GetUser(context.Background(), clusterID, name)
	if err != nil {
		// If the database user is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving Database User: %s", err)
	}

	d.Set("role", user.Role)
	// This will be blank for MongoDB clusters. Don't overwrite the password set on create.
	if user.Password != "" {
		d.Set("password", user.Password)
	}

	if user.MySQLSettings != nil {
		d.Set("mysql_auth_plugin", user.MySQLSettings.AuthPlugin)
	}

	return nil
}

func resourceDigitalOceanDatabaseUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	if d.HasChange("mysql_auth_plugin") {
		authReq := &godo.DatabaseResetUserAuthRequest{}
		if d.Get("mysql_auth_plugin").(string) != "" {
			authReq.MySQLSettings = &godo.DatabaseMySQLUserSettings{
				AuthPlugin: d.Get("mysql_auth_plugin").(string),
			}
		} else {
			// If blank, restore default value.
			authReq.MySQLSettings = &godo.DatabaseMySQLUserSettings{
				AuthPlugin: godo.SQLAuthPluginCachingSHA2,
			}
		}

		_, _, err := client.Databases.ResetUserAuth(context.Background(), d.Get("cluster_id").(string), d.Get("name").(string), authReq)
		if err != nil {
			return diag.Errorf("Error updating mysql_auth_plugin for DatabaseUser: %s", err)
		}
	}

	return resourceDigitalOceanDatabaseUserRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	// Prevent parallel deletion of users for same cluster.
	key := fmt.Sprintf("digitalocean_database_cluster/%s/users", clusterID)
	mutexKV.Lock(key)
	defer mutexKV.Unlock(key)

	log.Printf("[INFO] Deleting Database User: %s", d.Id())
	_, err := client.Databases.DeleteUser(context.Background(), clusterID, name)
	if err != nil {
		return diag.Errorf("Error deleting Database User: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceDigitalOceanDatabaseUserImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(makeDatabaseUserID(s[0], s[1]))
		d.Set("cluster_id", s[0])
		d.Set("name", s[1])
	} else {
		return nil, errors.New("must use the ID of the source database cluster and the name of the user joined with a comma (e.g. `id,name`)")
	}

	return []*schema.ResourceData{d}, nil
}

func makeDatabaseUserID(clusterID string, name string) string {
	return fmt.Sprintf("%s/user/%s", clusterID, name)
}
