package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/tag"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanDatabaseReplica() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDatabaseReplicaCreate,
		ReadContext:   resourceDigitalOceanDatabaseReplicaRead,
		UpdateContext: resourceDigitalOceanDatabaseReplicaUpdate,
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

			"uuid": {
				Type:        schema.TypeString,
				Computed:    true,
				ForceNew:    true,
				Description: "The unique universal identifier for the database replica.",
			},

			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"size": {
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

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: tag.ValidateTag,
				},
				Set: util.HashStringIgnoreCase,
			},

			"storage_size_mib": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceDigitalOceanDatabaseReplicaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterId := d.Get("cluster_id").(string)

	opts := &godo.DatabaseCreateReplicaRequest{
		Name:   d.Get("name").(string),
		Region: d.Get("region").(string),
		Size:   d.Get("size").(string),
		Tags:   tag.ExpandTags(d.Get("tags").(*schema.Set).List()),
	}

	if v, ok := d.GetOk("private_network_uuid"); ok {
		opts.PrivateNetworkUUID = v.(string)
	}

	if v, ok := d.GetOk("storage_size_mib"); ok {
		v, err := strconv.ParseUint(v.(string), 10, 64)
		if err == nil {
			opts.StorageSizeMib = v
		}
	}

	log.Printf("[DEBUG] DatabaseReplica create configuration: %#v", opts)

	var replicaCluster *godo.DatabaseReplica

	// Retry requests that fail w. Failed Precondition (412). New DBs can be marked ready while
	// first backup is still being created.
	err := resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		rc, resp, err := client.Databases.CreateReplica(context.Background(), clusterId, opts)
		if err != nil {
			if resp.StatusCode == 412 {
				return resource.RetryableError(err)
			} else {
				return resource.NonRetryableError(fmt.Errorf("Error creating DatabaseReplica: %s", err))
			}
		}
		replicaCluster = rc

		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	err = setReplicaConnectionInfo(replicaCluster, d)
	if err != nil {
		return diag.Errorf("Error building connection URI: %s", err)
	}

	replica, err := waitForDatabaseReplica(client, clusterId, "online", replicaCluster.Name)
	if err != nil {
		return diag.Errorf("Error creating DatabaseReplica: %s", err)
	}

	// Terraform requires a unique ID for each resource,
	// this concatc the parent cluster's ID and the replica's
	// name to form a replica's ID for Terraform state. This is
	// before the replica's ID was exposed in the DO API.
	d.SetId(makeReplicaId(clusterId, replica.Name))
	// the replica ID is now exposed in the DO API. It can be referenced
	// via the uuid in order to not change Terraform's
	// internal ID for existing resources.
	d.Set("uuid", replica.ID)
	log.Printf("[INFO] DatabaseReplica Name: %s", replica.Name)

	return resourceDigitalOceanDatabaseReplicaRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseReplicaRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
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

	d.Set("size", replica.Size)
	d.Set("region", replica.Region)
	d.Set("tags", tag.FlattenTags(replica.Tags))

	// Computed values
	d.Set("uuid", replica.ID)
	d.Set("private_network_uuid", replica.PrivateNetworkUUID)
	d.Set("storage_size_mib", strconv.FormatUint(replica.StorageSizeMib, 10))

	err = setReplicaConnectionInfo(replica, d)
	if err != nil {
		return diag.Errorf("Error building connection URI: %s", err)
	}

	return nil
}

func setReplicaConnectionInfo(replica *godo.DatabaseReplica, d *schema.ResourceData) error {
	if replica.Connection != nil {
		d.Set("host", replica.Connection.Host)
		d.Set("port", replica.Connection.Port)
		d.Set("database", replica.Connection.Database)
		d.Set("user", replica.Connection.User)

		if replica.Connection.Password != "" {
			d.Set("password", replica.Connection.Password)
		}

		uri, err := buildDBConnectionURI(replica.Connection, d)
		if err != nil {
			return err
		}

		d.Set("uri", uri)
	}

	if replica.PrivateConnection != nil {
		d.Set("private_host", replica.PrivateConnection.Host)

		privateURI, err := buildDBConnectionURI(replica.PrivateConnection, d)
		if err != nil {
			return err
		}

		d.Set("private_uri", privateURI)
	}

	return nil
}

func resourceDigitalOceanDatabaseReplicaUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	clusterID := d.Get("cluster_id").(string)
	replicaID := d.Get("uuid").(string)
	replicaName := d.Get("name").(string)

	if d.HasChanges("size", "storage_size_mib") {
		opts := &godo.DatabaseResizeRequest{
			SizeSlug: d.Get("size").(string),
			NumNodes: 1, // Read-only replicas only support a single node configuration.
		}

		if v, ok := d.GetOk("storage_size_mib"); ok {
			v, err := strconv.ParseUint(v.(string), 10, 64)
			if err == nil {
				opts.StorageSizeMib = v
			}
		}

		resp, err := client.Databases.Resize(context.Background(), replicaID, opts)
		if err != nil {
			if resp != nil && resp.StatusCode == 404 {
				d.SetId("")
				return nil
			}

			return diag.Errorf("Error resizing database replica: %s", err)
		}

		_, err = waitForDatabaseReplica(client, clusterID, "online", replicaName)
		if err != nil {
			return diag.Errorf("Error resizing database replica: %s", err)
		}
	}

	return resourceDigitalOceanDatabaseReplicaRead(ctx, d, meta)
}

func resourceDigitalOceanDatabaseReplicaImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if strings.Contains(d.Id(), ",") {
		s := strings.Split(d.Id(), ",")
		d.SetId(makeReplicaId(s[0], s[1]))
		d.Set("cluster_id", s[0])
		d.Set("name", s[1])
	} else {
		return nil, errors.New("must use the ID of the source database cluster and the replica name joined with a comma (e.g. `id,name`)")
	}

	return []*schema.ResourceData{d}, nil
}

func resourceDigitalOceanDatabaseReplicaDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
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
		replica, resp, err := client.Databases.GetReplica(context.Background(), cluster_id, name)
		if resp.StatusCode == 404 {
			continue
		}

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
