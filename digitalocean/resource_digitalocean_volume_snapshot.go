package digitalocean

import (
	"context"
	"log"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDigitalOceanVolumeSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanVolumeSnapshotCreate,
		ReadContext:   resourceDigitalOceanVolumeSnapshotRead,
		UpdateContext: resourceDigitalOceanVolumeSnapshotUpdate,
		DeleteContext: resourceDigitalOceanVolumeSnapshotDelete,
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

			"volume_id": {
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

			"tags": tagsSchema(),
		},
	}
}

func resourceDigitalOceanVolumeSnapshotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	opts := &godo.SnapshotCreateRequest{
		Name:     d.Get("name").(string),
		VolumeID: d.Get("volume_id").(string),
		Tags:     expandTags(d.Get("tags").(*schema.Set).List()),
	}

	log.Printf("[DEBUG] Volume Snapshot create configuration: %#v", opts)
	snapshot, _, err := client.Storage.CreateSnapshot(context.Background(), opts)
	if err != nil {
		return diag.Errorf("Error creating Volume Snapshot: %s", err)
	}

	d.SetId(snapshot.ID)
	log.Printf("[INFO] Volume Snapshot name: %s", snapshot.Name)

	return resourceDigitalOceanVolumeSnapshotRead(ctx, d, meta)
}

func resourceDigitalOceanVolumeSnapshotUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	if d.HasChange("tags") {
		err := setTags(client, d, godo.VolumeSnapshotResourceType)
		if err != nil {
			return diag.Errorf("Error updating tags: %s", err)
		}
	}

	return resourceDigitalOceanVolumeSnapshotRead(ctx, d, meta)
}

func resourceDigitalOceanVolumeSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	snapshot, resp, err := client.Snapshots.Get(context.Background(), d.Id())
	if err != nil {
		// If the snapshot is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving volume snapshot: %s", err)
	}

	d.Set("name", snapshot.Name)
	d.Set("volume_id", snapshot.ResourceID)
	d.Set("regions", snapshot.Regions)
	d.Set("size", snapshot.SizeGigaBytes)
	d.Set("created_at", snapshot.Created)
	d.Set("min_disk_size", snapshot.MinDiskSize)
	d.Set("tags", flattenTags(snapshot.Tags))

	return nil
}

func resourceDigitalOceanVolumeSnapshotDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	log.Printf("[INFO] Deleting snaphot: %s", d.Id())
	_, err := client.Snapshots.Delete(context.Background(), d.Id())
	if err != nil {
		return diag.Errorf("Error deleting snapshot: %s", err)
	}

	d.SetId("")
	return nil
}
