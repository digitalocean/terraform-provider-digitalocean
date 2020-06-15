package digitalocean

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceDigitalOceanTag() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanTagRead,
		Schema: map[string]*schema.Schema{

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "name of the tag",
				ValidateFunc: validateTag,
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

func dataSourceDigitalOceanTagRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	name := d.Get("name").(string)

	tag, resp, err := client.Tags.Get(context.Background(), name)
	if err != nil {
		if resp != nil && resp.StatusCode == 404 {
			return fmt.Errorf("tag not found: %s", err)
		}
		return fmt.Errorf("Error retrieving tag: %s", err)
	}

	d.SetId(tag.Name)
	d.Set("name", tag.Name)
	d.Set("total_resource_count", tag.Resources.Count)
	d.Set("droplets_count", tag.Resources.Droplets.Count)
	d.Set("images_count", tag.Resources.Images.Count)
	d.Set("volumes_count", tag.Resources.Volumes.Count)
	d.Set("volume_snapshots_count", tag.Resources.VolumeSnapshots.Count)
	d.Set("databases_count", tag.Resources.Databases.Count)

	return nil
}
