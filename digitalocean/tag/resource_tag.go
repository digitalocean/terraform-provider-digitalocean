package tag

import (
	"context"
	"log"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceDigitalOceanTag() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDigitalOceanTagCreate,
		ReadContext:   resourceDigitalOceanTagRead,
		DeleteContext: resourceDigitalOceanTagDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: ValidateTag,
			},
			"total_resource_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"droplets_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"images_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"volumes_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"volume_snapshots_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"databases_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceDigitalOceanTagCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	// Build up our creation options
	opts := &godo.TagCreateRequest{
		Name: d.Get("name").(string),
	}

	log.Printf("[DEBUG] Tag create configuration: %#v", opts)
	tag, _, err := client.Tags.Create(context.Background(), opts)
	if err != nil {
		return diag.Errorf("Error creating tag: %s", err)
	}

	d.SetId(tag.Name)
	log.Printf("[INFO] Tag: %s", tag.Name)

	return resourceDigitalOceanTagRead(ctx, d, meta)
}

func resourceDigitalOceanTagRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	tag, resp, err := client.Tags.Get(context.Background(), d.Id())
	if err != nil {
		// If the tag is somehow already destroyed, mark as
		// successfully gone
		if resp != nil && resp.StatusCode == 404 {
			d.SetId("")
			return nil
		}

		return diag.Errorf("Error retrieving tag: %s", err)
	}

	d.Set("name", tag.Name)
	d.Set("total_resource_count", tag.Resources.Count)
	d.Set("droplets_count", tag.Resources.Droplets.Count)
	d.Set("images_count", tag.Resources.Images.Count)
	d.Set("volumes_count", tag.Resources.Volumes.Count)
	d.Set("volume_snapshots_count", tag.Resources.VolumeSnapshots.Count)
	d.Set("databases_count", tag.Resources.Databases.Count)

	return nil
}

func resourceDigitalOceanTagDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	log.Printf("[INFO] Deleting tag: %s", d.Id())
	_, err := client.Tags.Delete(context.Background(), d.Id())
	if err != nil {
		return diag.Errorf("Error deleting tag: %s", err)
	}

	d.SetId("")
	return nil
}
