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
			"droplet_id": {
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

	resourceId, _ := strconv.Atoi(d.Get("droplet_id").(string))
	action, _, err := client.DropletActions.Snapshot(context.Background(), resourceId, d.Get("name").(string))
	if err != nil {
		return fmt.Errorf("Error creating Droplet Snapshot: %s", err)
	}

	if err = waitForAction(client, action); err != nil {
		return fmt.Errorf(
			"Error waiting for Droplet snapshot (%v) to finish: %s", resourceId, err)
	}

	v, err := findSnapshotInSnapshotList(context.Background(), client, *action.StartedAt)

	if err != nil {
		return fmt.Errorf("Error retriving Droplet Snapshot: %s", err)
	}

	d.SetId(v.ID)
	d.Set("name", v.Name)
	d.Set("droplet_id", v.ID)
	d.Set("regions", v.Regions)
	d.Set("created_at", v.Created)
	d.Set("min_disk_size", v.MinDiskSize)

	return resourceDigitalOceanDropletSnapshotRead(d, meta)
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
	d.Set("droplet_id", snapshot.ResourceID)
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

func findSnapshotInSnapshotList(ctx context.Context, client *godo.Client, actionTime godo.Timestamp) (godo.Snapshot, error) {
	opt := &godo.ListOptions{}
	for {
		snapshots, resp, err := client.Snapshots.List(ctx, opt)
		if err != nil {
			return godo.Snapshot{}, err
		}

		// check the current pages droplets for our snapshot
		for _, d := range snapshots {
			createdTime, _ := time.Parse("2006-01-02T15:04:05Z", d.Created)
			checkTime := godo.Timestamp{Time: createdTime}
			if checkTime == actionTime {
				return d, nil
			}
		}

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return godo.Snapshot{}, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}
	return godo.Snapshot{}, fmt.Errorf("Error Could not locate the Droplet Snaphot")
}
