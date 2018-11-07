package digitalocean

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceDigitalOceanDropletSnapshot() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanDropletSnapshotCreate,
		Read:   resourceDigitalOceanDropletSnapshotRead,
		Delete: resourceDigitalOceanDropletSnapshotDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"regions": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"size": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"min_disk_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceDigitalOceanDropletSnapshotCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	resourceId, _ := strconv.Atoi(d.Get("resource_id").(string))
	action, _, err := client.DropletActions.Snapshot(context.Background(), resourceId, d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("Error creating Droplet Snapshot: %s", err)
	}

	for ok := true; ok; ok = checkActionProgress(action.Status) {
		action, _, _ = client.Actions.Get(context.Background(), action.ID)
		time.Sleep(10)
	}
	opt := godo.ListOptions{Page: 1, PerPage: 200}
	snapshotList, _, _ := client.Droplets.Snapshots(context.Background(), action.ResourceID, &opt)

	for _, v := range snapshotList {
		createdTime, _ := time.Parse("2006-01-02T15:04:05Z", v.Created)
		checkTime := godo.Timestamp{createdTime}

		if checkTime == *action.StartedAt {
			d.SetId(strconv.Itoa(v.ID))
			d.Set("name", v.Name)
			d.Set("resource_id", strconv.Itoa(v.ID))
			d.Set("regions", v.Regions)
			d.Set("created_at", v.Created)
			d.Set("min_disk_size", v.MinDiskSize)
		}
	}
	return resourceDigitalOceanDropletSnapshotRead(d, meta)
}

func checkActionProgress(actionProgress string) bool {
	if actionProgress == "in-progress" {
		return true
	}
	return false
}

func resourceDigitalOceanDropletSnapshotRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	snapshot, resp, err := client.Snapshots.Get(context.Background(), d.Id())
	if err != nil {
		// If the snapshot is somehow already destroyed, mark as
		// successfully gone
		if resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error retrieving Droplet snapshot: %s", err)
	}

	d.Set("name", snapshot.Name)
	d.Set("resource_id", snapshot.ResourceID)
	d.Set("regions", snapshot.Regions)
	d.Set("size", snapshot.SizeGigaBytes)
	d.Set("created_at", snapshot.Created)
	d.Set("min_disk_size", snapshot.MinDiskSize)

	return nil
}

func resourceDigitalOceanDropletSnapshotDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	log.Printf("[INFO] Deleting snaphot: %s", d.Id())
	_, err := client.Snapshots.Delete(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting snapshot: %s", err)
	}

	d.SetId("")
	return nil
}
