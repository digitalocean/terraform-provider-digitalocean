package digitalocean

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceDigitalOceanDatabase() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanDatabaseCreate,
		Read:   resourceDigitalOceanDatabaseRead,
		Delete: resourceDigitalOceanDatabaseDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDigitalOceanDatabaseImport,
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
		},
	}
}

func resourceDigitalOceanDatabaseCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
	clusterID := d.Get("cluster_id").(string)

	opts := &godo.DatabaseCreateDBRequest{
		Name: d.Get("name").(string),
	}

	log.Printf("[DEBUG] Database create configuration: %#v", opts)
	db, _, err := client.Databases.CreateDB(context.Background(), clusterID, opts)
	if err != nil {
		return fmt.Errorf("Error creating Database: %s", err)
	}

	d.SetId(makeDatabaseID(clusterID, db.Name))
	log.Printf("[INFO] Database Name: %s", db.Name)

	return resourceDigitalOceanDatabaseRead(d, meta)
}

func resourceDigitalOceanDatabaseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	// Check if the database still exists
	_, resp, err := client.Databases.GetDB(context.Background(), clusterID, name)
	if err != nil {
		// If the database is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Database: %s", err)
	}

	return nil
}

func resourceDigitalOceanDatabaseDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	log.Printf("[INFO] Deleting Database: %s", d.Id())
	_, err := client.Databases.DeleteDB(context.Background(), clusterID, name)
	if err != nil {
		return fmt.Errorf("Error deleting Database: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceDigitalOceanDatabaseImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(makeDatabaseID(s[0], s[1]))
		d.Set("cluster_id", s[0])
		d.Set("name", s[1])
	}

	return []*schema.ResourceData{d}, nil
}

func makeDatabaseID(clusterID string, name string) string {
	return fmt.Sprintf("%s/database/%s", clusterID, name)
}
