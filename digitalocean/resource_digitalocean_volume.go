package digitalocean

import (
	"context"
	"fmt"
	"log"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/schema"
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"droplet_ids": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Computed: true,
			},

			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true, // Update-ability Coming Soon ™
			},

			"filesystem_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},

		// CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {

		// 	// if the new size of the volume is smaller than the old one force a new resource since
		// 	// only expanding the volume is allowed
		// 	oldSize, newSize := diff.GetChange("size")
		// 	if newSize < oldSize {
		// 		diff.ForceNew("size")
		// 	}
		// },
	}
}

func resourceDigitalOceanVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*godo.Client)

	opts := &godo.VolumeCreateRequest{
		Region:         d.Get("region").(string),
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		SizeGigaBytes:  int64(d.Get("size").(int)),
		FilesystemType: d.Get("filesystem_type").(string),
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

	d.Set("size", int(volume.SizeGigaBytes))

	dids := make([]interface{}, 0, len(volume.DropletIDs))
	for _, did := range volume.DropletIDs {
		dids = append(dids, did)
	}
	d.Set("droplet_ids", schema.NewSet(
		func(dropletID interface{}) int { return dropletID.(int) },
		dids,
	))

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
	rs.Set("description", volume.Description)
	rs.Set("size", int(volume.SizeGigaBytes))
	rs.Set("filesystem_type", volume.FilesystemType)

	dids := make([]interface{}, 0, len(volume.DropletIDs))
	for _, did := range volume.DropletIDs {
		dids = append(dids, did)
	}
	rs.Set("droplet_ids", schema.NewSet(
		func(dropletID interface{}) int { return dropletID.(int) },
		dids,
	))

	return []*schema.ResourceData{rs}, nil
}
