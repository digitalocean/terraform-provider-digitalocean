package digitalocean

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceDigitalOceanDatabaseReplica() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanDatabaseReplicaCreate,
		Read:   resourceDigitalOceanDatabaseReplicaRead,
		Delete: resourceDigitalOceanDatabaseReplicaDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDigitalOceanDatabaseReplicaImport,
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

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"size": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"uri": {
				Type:     schema.TypeString,
				Computed: true,
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
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validateTag,
				},
				Set: HashStringIgnoreCase,
			},
		},
	}
}

func resourceDigitalOceanDatabaseReplicaCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
	clusterId := d.Get("cluster_id").(string)

	opts := &godo.DatabaseCreateReplicaRequest{
		Name:   d.Get("name").(string),
		Region: d.Get("region").(string),
		Size:   d.Get("size").(string),
		Tags:   expandTags(d.Get("tags").(*schema.Set).List()),
	}

	log.Printf("[DEBUG] DatabaseReplica create configuration: %#v", opts)
	replica, _, err := client.Databases.CreateReplica(context.Background(), clusterId, opts)
	if err != nil {
		return fmt.Errorf("Error creating DatabaseReplica: %s", err)
	}

	replica, err = waitForDatabaseReplica(client, clusterId, "online", replica.Name)
	if err != nil {
		return fmt.Errorf("Error creating DatabaseReplica: %s", err)
	}

	d.SetId(makeReplicaId(clusterId, replica.Name))
	log.Printf("[INFO] DatabaseReplica Name: %s", replica.Name)

	return resourceDigitalOceanDatabaseReplicaRead(d, meta)
}

func resourceDigitalOceanDatabaseReplicaRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)
	replica, resp, err := client.Databases.GetReplica(context.Background(), clusterId, name)
	if err != nil {
		// If the database is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving DatabaseReplica: %s", err)
	}

	d.Set("region", replica.Region)
	d.Set("tags", replica.Tags)

	// Computed values
	d.Set("host", replica.Connection.Host)
	d.Set("port", replica.Connection.Port)
	d.Set("uri", replica.Connection.URI)
	d.Set("database", replica.Connection.Database)
	d.Set("user", replica.Connection.User)
	d.Set("password", replica.Connection.Password)

	return nil
}

func resourceDigitalOceanDatabaseReplicaImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(makeReplicaId(s[0], s[1]))
		d.Set("cluster_id", s[0])
		d.Set("name", s[1])
	}

	return []*schema.ResourceData{d}, nil
}

func resourceDigitalOceanDatabaseReplicaDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	log.Printf("[INFO] Deleting DatabaseReplica: %s", d.Id())
	_, err := client.Databases.DeleteReplica(context.Background(), clusterId, name)
	if err != nil {
		return fmt.Errorf("Error deleting DatabaseReplica: %s", err)
	}

	d.SetId("")
	return nil
}

func makeReplicaId(clusterId string, replicaName string) string {
	return fmt.Sprintf("%s/replicas/%s", clusterId, replicaName)
}

func waitForDatabaseReplica(client *godo.Client, cluster_id, status, name string) (*godo.DatabaseReplica, error) {
	ticker := time.NewTicker(15 * time.Second)
	timeout := 120
	n := 0

	for range ticker.C {
		replica, _, err := client.Databases.GetReplica(context.Background(), cluster_id, name)
		if err != nil {
			ticker.Stop()
			return nil, fmt.Errorf("Error trying to read DatabaseReplica state: %s", err)
		}

		if replica.Status == status {
			ticker.Stop()
			return replica, nil
		}

		if n > timeout {
			ticker.Stop()
			break
		}

		n++
	}

	return nil, fmt.Errorf("Timeout waiting to DatabaseReplica to become %s", status)
}
