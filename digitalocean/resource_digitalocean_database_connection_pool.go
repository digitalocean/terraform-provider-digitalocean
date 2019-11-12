package digitalocean

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceDigitalOceanDatabaseConnectionPool() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanDatabaseConnectionPoolCreate,
		Read:   resourceDigitalOceanDatabaseConnectionPoolRead,
		Delete: resourceDigitalOceanDatabaseConnectionPoolDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"db_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"user": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"mode": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"size": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"database": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"pool_host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"pool_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"pool_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"pool_database": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"pool_user": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"pool_password": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDigitalOceanDatabaseConnectionPoolCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	dbID := d.Get("db_id").(string)
	opts := &godo.DatabaseCreatePoolRequest{
		Name:     d.Get("name").(string),
		User:     d.Get("user").(string),
		Mode:     d.Get("mode").(string),
		Size:     d.Get("size").(int),
		Database: d.Get("database").(string),
	}

	log.Printf("[DEBUG] DatabaseConnectionPool create configuration: %#v", opts)
	pool, _, err := client.Databases.CreatePool(context.Background(), dbID, opts)
	if err != nil {
		return fmt.Errorf("Error creating DatabaseConnectionPool: %s", err)
	}

	d.SetId(createConnectionPoolID(dbID, pool.Name))
	log.Printf("[INFO] DatabaseConnectionPool Name: %s", pool.Name)

	return resourceDigitalOceanDatabaseConnectionPoolRead(d, meta)
}

func resourceDigitalOceanDatabaseConnectionPoolRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
	dbID, poolID := splitConnectionPoolID(d.Id())

	pool, resp, err := client.Databases.GetPool(context.Background(), dbID, poolID)
	if err != nil {
		// If the pool is somehow already destroyed, mark as
		// successfully gone
		if resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving DatabaseConnectionPool: %s", err)
	}

	d.SetId(createConnectionPoolID(dbID, pool.Name))
	d.Set("db_id", dbID)
	d.Set("name", pool.Name)
	d.Set("user", pool.User)
	d.Set("mode", pool.Mode)
	d.Set("size", pool.Size)
	d.Set("database", pool.Database)

	// Computed values
	d.Set("pool_host", pool.Connection.Host)
	d.Set("pool_port", pool.Connection.Port)
	d.Set("pool_uri", pool.Connection.URI)
	d.Set("pool_database", pool.Connection.Database)
	d.Set("pool_user", pool.Connection.User)
	d.Set("pool_password", pool.Connection.Password)

	return nil
}

func resourceDigitalOceanDatabaseConnectionPoolDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
	dbID, poolID := splitConnectionPoolID(d.Id())

	log.Printf("[INFO] Deleting DatabaseConnectionPool: %s", poolID)
	_, err := client.Databases.DeletePool(context.Background(), dbID, poolID)
	if err != nil {
		return fmt.Errorf("Error deleting DatabaseConnectionPool: %s", err)
	}

	d.SetId("")
	return nil
}

func createConnectionPoolID(dbID string, poolID string) string {
	return fmt.Sprintf("%s/%s", dbID, poolID)
}

func splitConnectionPoolID(id string) (string, string) {
	splitID := strings.Split(id, "/")
	dbID := splitID[0]
	poolID := splitID[1]

	return dbID, poolID
}
