package snapshot

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanDropletSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanDropletSnapshotCreate,
		ReadContext:   resourceDigitalOceanDropletSnapshotRead,
		DeleteContext: resourceDigitalOceanDropletSnapshotDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func resourceDigitalOceanDropletSnapshotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	resourceId, _ := strconv.Atoi(d.Get("droplet_id").(string))
	action, _, err := client.DropletActions.Snapshot(context.Background(), resourceId, d.Get("name").(string))
	if err != nil {
		return diag.Errorf("Error creating Droplet Snapshot: %s", err)
	}

	if err = util.WaitForAction(client, action); err != nil {
		return diag.Errorf(
			"Error waiting for Droplet snapshot (%v) to finish: %s", resourceId, err)
	}

	snapshot, err := findSnapshotInSnapshotList(context.Background(), client, *action)

	if err != nil {
		return diag.Errorf("Error retrieving Droplet Snapshot: %s", err)
	}

	d.SetId(strconv.Itoa(snapshot.ID))
	if err = d.Set("name", snapshot.Name); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("droplet_id", fmt.Sprintf("%d", snapshot.ID)); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("regions", snapshot.Regions); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("created_at", snapshot.Created); err != nil {
		return diag.FromErr(err)
	}
	if err = d.Set("min_disk_size", snapshot.MinDiskSize); err != nil {
		return diag.FromErr(err)
	}

	return resourceDigitalOceanDropletSnapshotRead(ctx, d, meta)
}

func resourceDigitalOceanDropletSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	snapshot, resp, err := client.Snapshots.Get(context.Background(), d.Id())
	if err != nil {
		// If the snapshot is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error retrieving Droplet snapshot: %s", err)
	}

	d.Set("name", snapshot.Name)
	d.Set("droplet_id", snapshot.ResourceID)
	d.Set("regions", snapshot.Regions)
	d.Set("size", snapshot.SizeGigaBytes)
	d.Set("created_at", snapshot.Created)
	d.Set("min_disk_size", snapshot.MinDiskSize)

	return nil
}

func resourceDigitalOceanDropletSnapshotDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	log.Printf("[INFO] Deleting snapshot: %s", d.Id())
	_, err := client.Snapshots.Delete(context.Background(), d.Id())
	if err != nil {
		return diag.Errorf("Error deleting snapshot: %s", err)
	}

	d.SetId("")
	return nil
}

func findSnapshotInSnapshotList(ctx context.Context, client *godo.Client, action godo.Action) (godo.Image, error) {
	opt := &godo.ListOptions{PerPage: 200}
	for {
		snapshots, resp, err := client.Droplets.Snapshots(ctx, action.ResourceID, opt)
		if err != nil {
			return godo.Image{}, err
		}

		// check the current page for our snapshot
		for _, s := range snapshots {
			createdTime, _ := time.Parse("2006-01-02T15:04:05Z", s.Created)
			checkTime := &godo.Timestamp{Time: createdTime}
			if *checkTime == *action.StartedAt {
				return s, nil
			}
		}

		// if we are at the last page, break out the for loop
		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return godo.Image{}, err
		}

		// set the page we want for the next request
		opt.Page = page + 1
	}
	return godo.Image{}, fmt.Errorf("Error Could not locate the Droplet Snapshot")
}
