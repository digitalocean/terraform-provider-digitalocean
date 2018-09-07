package digitalocean

import (
	"context"
	"fmt"
	"log"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceDigitalOceanVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourceDigitalOceanVolumeCreate,
		Read:   resourceDigitalOceanVolumeRead,
		Update: resourceDigitalOceanVolumeUpdate,
		Delete: resourceDigitalOceanVolumeDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDigitalOceanVolumeImport,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true, // Update-ability Coming Soon ™
				ValidateFunc: validation.NoZeroValues,
			},

			"snapshot_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"initial_filesystem_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ext4",
					"xfs",
				}, false),
			},

			"initial_filesystem_label": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"droplet_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Computed: true,
			},

			"filesystem_type": {
				Type:     schema.TypeString,
				Optional: true, // Backward compatibility for existing resources.
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ext4",
					"xfs",
				}, false),
				ConflictsWith: []string{"initial_filesystem_type"},
				Deprecated:    "This fields functionality has been replaced by `initial_filesystem_type`. The property will still remain as a computed attribute representing the current volumes filesystem type.",
			},

			"filesystem_label": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},

		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {

			// if the new size of the volume is smaller than the old one return an error since
			// only expanding the volume is allowed
			oldSize, newSize := diff.GetChange("size")
			log.Printf("QQQQQ %v %v", oldSize, newSize)
			if newSize.(int) < oldSize.(int) {
				return fmt.Errorf("volumes `size` can only be expanded and not shrunk")
			}

			return nil
		},
	}
}

func resourceDigitalOceanVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	opts := &godo.VolumeCreateRequest{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
	}

	if v, ok := d.GetOk("region"); ok {
		opts.Region = v.(string)
	}
	if v, ok := d.GetOk("size"); ok {
		opts.SizeGigaBytes = int64(v.(int))
	}
	if v, ok := d.GetOk("snapshot_id"); ok {
		opts.SnapshotID = v.(string)
	}
	if v, ok := d.GetOk("initial_filesystem_type"); ok {
		opts.FilesystemType = v.(string)
	} else if v, ok := d.GetOk("filesystem_type"); ok {
		// backward compatibility
		opts.FilesystemType = v.(string)
	}
	if v, ok := d.GetOk("initial_filesystem_label"); ok {
		opts.FilesystemLabel = v.(string)
	}

	log.Printf("[DEBUG] Volume create configuration: %#v", opts)
	volume, _, err := client.Storage.CreateVolume(context.Background(), opts)
	if err != nil {
		return fmt.Errorf("Error creating Volume: %s", err)
	}

	d.SetId(volume.ID)
	log.Printf("[INFO] Volume name: %s", volume.Name)

	return resourceDigitalOceanVolumeRead(d, meta)
}

func resourceDigitalOceanVolumeUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	id := d.Id()
	region := d.Get("region").(string)

	if d.HasChange("size") {
		size := d.Get("size").(int)

		log.Printf("[DEBUG] Volume resize configuration: %v", size)
		action, _, err := client.StorageActions.Resize(context.Background(), id, size, region)
		if err != nil {
			return fmt.Errorf("Error resizing volume (%s): %s", id, err)
		}

		log.Printf("[DEBUG] Volume resize action id: %d", action.ID)
		if err = waitForAction(client, action); err != nil {
			return fmt.Errorf(
				"Error waiting for resize volume (%s) to finish: %s", id, err)
		}
	}

	return resourceDigitalOceanVolumeRead(d, meta)
}

func resourceDigitalOceanVolumeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	volume, resp, err := client.Storage.GetVolume(context.Background(), d.Id())
	if err != nil {
		// If the volume is somehow already destroyed, mark as
		// successfully gone
		if resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving volume: %s", err)
	}

	d.Set("region", volume.Region.Slug)
	d.Set("size", int(volume.SizeGigaBytes))

	if v := volume.FilesystemType; v != "" {
		d.Set("filesystem_type", v)
	}
	if v := volume.FilesystemLabel; v != "" {
		d.Set("filesystem_label", v)
	}

	if err = d.Set("droplet_ids", flattenDigitalOceanVolumeDropletIds(volume.DropletIDs)); err != nil {
		return fmt.Errorf("[DEBUG] Error setting droplet_ids: %#v", err)
	}

	return nil
}

func resourceDigitalOceanVolumeDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	log.Printf("[INFO] Deleting volume: %s", d.Id())
	_, err := client.Storage.DeleteVolume(context.Background(), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting volume: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceDigitalOceanVolumeImport(rs *schema.ResourceData, v interface{}) ([]*schema.ResourceData, error) {
	client := v.(*godo.Client)
	volume, _, err := client.Storage.GetVolume(context.Background(), rs.Id())
	if err != nil {
		return nil, err
	}

	rs.Set("name", volume.Name)
	rs.Set("region", volume.Region.Slug)
	rs.Set("size", int(volume.SizeGigaBytes))

	if v := volume.Description; v != "" {
		rs.Set("description", v)
	}
	if v := volume.FilesystemType; v != "" {
		rs.Set("filesystem_type", v)
	}
	if v := volume.FilesystemLabel; v != "" {
		rs.Set("filesystem_label", v)
	}

	if err = rs.Set("droplet_ids", flattenDigitalOceanVolumeDropletIds(volume.DropletIDs)); err != nil {
		return nil, fmt.Errorf("[DEBUG] Error setting droplet_ids: %#v", err)
	}

	return []*schema.ResourceData{rs}, nil
}

// Seperate validation function to support common cumputed
func validateDigitalOceanVolumeSchema(d *schema.ResourceData) error {
	_, hasRegion := d.GetOk("region")
	_, hasSize := d.GetOk("size")
	_, hasSnapshotId := d.GetOk("snapshot_id")
	if !hasSnapshotId {
		if !hasRegion {
			return fmt.Errorf("`region` must be assigned when not specifying a `snapshot_id`")
		}
		if !hasSize {
			return fmt.Errorf("`size` must be assigned when not specifying a `snapshot_id`")
		}
	}

	return nil
}

func flattenDigitalOceanVolumeDropletIds(droplets []int) *schema.Set {
	flattenedDroplets := schema.NewSet(schema.HashInt, []interface{}{})
	for _, v := range droplets {
		flattenedDroplets.Add(v)
	}

	return flattenedDroplets
}
