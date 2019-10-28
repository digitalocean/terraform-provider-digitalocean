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

func resourceDigitalOceanDatabaseUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanDatabaseUserCreate,
		Read:   resourceDigitalOceanDatabaseUserRead,
		Delete: resourceDigitalOceanDatabaseUserDelete,
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

func resourceDigitalOceanDatabaseUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
	clusterID := d.Get("cluster_id").(string)

	opts := &godo.DatabaseCreateUserRequest{
		Name: d.Get("name").(string),
	}

	log.Printf("[DEBUG] Database User create configuration: %#v", opts)
	user, _, err := client.Databases.CreateUser(context.Background(), clusterID, opts)
	if err != nil {
		return fmt.Errorf("Error creating Database User: %s", err)
	}

	d.SetId(makeDatabaseUserID(clusterID, user.Name))
	log.Printf("[INFO] Database User Name: %s", user.Name)

	return resourceDigitalOceanDatabaseUserRead(d, meta)
}

func resourceDigitalOceanDatabaseUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
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

		return fmt.Errorf("Error retrieving Database User: %s", err)
	}

	d.Set("role", user.Role)
	d.Set("password", user.Password)

	return nil
}

func resourceDigitalOceanDatabaseUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
	clusterID := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	log.Printf("[INFO] Deleting Database User: %s", d.Id())
	_, err := client.Databases.DeleteUser(context.Background(), clusterID, name)
	if err != nil {
		return fmt.Errorf("Error deleting Database User: %s", err)
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
	}

	return []*schema.ResourceData{d}, nil
}

func makeDatabaseUserID(clusterID string, name string) string {
	return fmt.Sprintf("%s/user/%s", clusterID, name)
}
