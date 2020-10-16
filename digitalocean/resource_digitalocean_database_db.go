package digitalocean

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDigitalOceanDatabaseDB() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanDatabaseDBCreate,
		Read:   resourceDigitalOceanDatabaseDBRead,
		Delete: resourceDigitalOceanDatabaseDBDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDigitalOceanDatabaseDBImport,
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

func resourceDigitalOceanDatabaseDBCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
	clusterID := d.Get("cluster_id").(string)

	opts := &godo.DatabaseCreateDBRequest{
		Name: d.Get("name").(string),
	}

	log.Printf("[DEBUG] Database DB create configuration: %#v", opts)
	db, _, err := client.Databases.CreateDB(context.Background(), clusterID, opts)
	if err != nil {
		return fmt.Errorf("Error creating Database DB: %s", err)
	}

	d.SetId(makeDatabaseDBID(clusterID, db.Name))
	log.Printf("[INFO] Database DB Name: %s", db.Name)

	return resourceDigitalOceanDatabaseDBRead(d, meta)
}

func resourceDigitalOceanDatabaseDBRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	// Check if the database DB still exists
	_, resp, err := client.Databases.GetDB(context.Background(), clusterID, name)
	if err != nil {
		// If the database DB is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Database DB: %s", err)
	}

	return nil
}

func resourceDigitalOceanDatabaseDBDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	log.Printf("[INFO] Deleting Database DB: %s", d.Id())
	_, err := client.Databases.DeleteDB(context.Background(), clusterID, name)
	if err != nil {
		return fmt.Errorf("Error deleting Database DB: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceDigitalOceanDatabaseDBImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(makeDatabaseDBID(s[0], s[1]))
		d.Set("cluster_id", s[0])
		d.Set("name", s[1])
	}

	return []*schema.ResourceData{d}, nil
}

func makeDatabaseDBID(clusterID string, name string) string {
	return fmt.Sprintf("%s/database/%s", clusterID, name)
}
