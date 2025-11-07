package nfs

import (
	"context"
	"log"
	"strings"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceDigitalOceanNfsSnapshot() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanNfsSnapshotCreate,
		ReadContext:   resourceDigitalOceanNfsSnapshotRead,
		DeleteContext: resourceDigitalOceanNfsSnapshotDelete,
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

			"share_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},

			"region": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(val interface{}) string {
					// DO API V2 region slug is always lowercase
					return strings.ToLower(val.(string))
				},
			},

			"size": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDigitalOceanNfsSnapshotCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	name := d.Get("name").(string)
	shareID := d.Get("share_id").(string)
	region := d.Get("region").(string)

	snapshot, _, err := client.NfsActions.Snapshot(context.Background(), shareID, name, region)
	if err != nil {
		return diag.Errorf("Error creating Share Snapshot: %s", err)
	}

	d.SetId(snapshot.ResourceID)

	return resourceDigitalOceanNfsSnapshotRead(ctx, d, meta)
}

func resourceDigitalOceanNfsSnapshotRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	region := d.Get("region").(string)
	snapshot, resp, err := client.Nfs.GetSnapshot(context.Background(), d.Id(), region)
	if err != nil {
		// If the snapshot is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving share snapshot: %s", err)
	}

	d.Set("name", snapshot.Name)
	d.Set("share_id", snapshot.ShareID)
	d.Set("region", snapshot.Region)
	d.Set("size", snapshot.SizeGib)
	d.Set("created_at", snapshot.CreatedAt)

	return nil
}

func resourceDigitalOceanNfsSnapshotDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	region := d.Get("region").(string)

	log.Printf("[INFO] Deleting snapshot: %s", d.Id())
	_, err := client.Nfs.DeleteSnapshot(context.Background(), d.Id(), region)
	if err != nil {
		return diag.Errorf("Error deleting snapshot: %s", err)
	}

	d.SetId("")
	return nil
}
