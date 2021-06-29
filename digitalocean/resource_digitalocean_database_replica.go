package digitalocean

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDigitalOceanDatabaseReplica() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDatabaseReplicaCreate,
		ReadContext:   resourceDigitalOceanDatabaseReplicaRead,
		DeleteContext: resourceDigitalOceanDatabaseReplicaDelete,
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

func resourceDigitalOceanDatabaseReplicaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()
	clusterId := d.Get("cluster_id").(string)

	opts := &godo.DatabaseCreateReplicaRequest{
		Name:   d.Get("name").(string),
		Region: d.Get("region").(string),
		Size:   d.Get("size").(string),
		Tags:   expandTags(d.Get("tags").(*schema.Set).List()),
	}

	if v, ok := d.GetOk("private_network_uuid"); ok {
		opts.PrivateNetworkUUID = v.(string)
	}

	log.Printf("[DEBUG] DatabaseReplica create configuration: %#v", opts)
	replica, _, err := client.Databases.CreateReplica(context.Background(), clusterId, opts)
	if err != nil {
		return diag.Errorf("Error creating DatabaseReplica: %s", err)
	}

	replica, err = waitForDatabaseReplica(client, clusterId, "online", replica.Name)
	if err != nil {
		return diag.Errorf("Error creating DatabaseReplica: %s", err)
	}

	d.SetId(makeReplicaId(clusterId, replica.Name))
	log.Printf("[INFO] DatabaseReplica Name: %s", replica.Name)

	return resourceDigitalOceanDatabaseReplicaRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseReplicaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

		return diag.Errorf("Error retrieving DatabaseReplica: %s", err)
	}

	d.Set("region", replica.Region)
	d.Set("tags", flattenTags(replica.Tags))

	// Computed values
	d.Set("host", replica.Connection.Host)
	d.Set("private_host", replica.PrivateConnection.Host)
	d.Set("port", replica.Connection.Port)
	d.Set("uri", replica.Connection.URI)
	d.Set("private_uri", replica.PrivateConnection.URI)
	d.Set("database", replica.Connection.Database)
	d.Set("user", replica.Connection.User)
	d.Set("password", replica.Connection.Password)
	d.Set("private_network_uuid", replica.PrivateNetworkUUID)

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

func resourceDigitalOceanDatabaseReplicaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()
	clusterId := d.Get("cluster_id").(string)
	name := d.Get("name").(string)

	log.Printf("[INFO] Deleting DatabaseReplica: %s", d.Id())
	_, err := client.Databases.DeleteReplica(context.Background(), clusterId, name)
	if err != nil {
		return diag.Errorf("Error deleting DatabaseReplica: %s", err)
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
