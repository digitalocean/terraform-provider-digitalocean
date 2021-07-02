package digitalocean

import (
	"context"
	"fmt"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceDigitalOceanVolume() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanVolumeRead,
		Schema: map[string]*schema.Schema{

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "name of the volume",
				ValidateFunc: validation.NoZeroValues,
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "the region that the volume is provisioned in",
				StateFunc: func(val interface{}) string {
					// DO API V2 region slug is always lowercase
					return strings.ToLower(val.(string))
				},
				ValidateFunc: validation.NoZeroValues,
			},
			// computed attributes
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the iniform resource name for the volume",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "volume description",
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the size of the volume in gigabytes",
			},
			"filesystem_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the type of filesystem currently in-use on the volume",
			},
			"filesystem_label": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "the label currently applied to the filesystem",
			},
			"droplet_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeInt},
				Computed:    true,
				Description: "list of droplet ids the volume is attached to",
			},
			"tags": tagsDataSourceSchema(),
		},
	}
}

func dataSourceDigitalOceanVolumeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*CombinedConfig).godoClient()

	name := d.Get("name").(string)

	opts := &godo.ListVolumeParams{
		Name: name,
		ListOptions: &godo.ListOptions{
			Page:    1,
			PerPage: 200,
		},
	}

	if v, ok := d.GetOk("region"); ok {
		opts.Region = v.(string)
	}

	volumeList := []godo.Volume{}

	for {
		volumes, resp, err := client.Storage.ListVolumes(context.Background(), opts)

		if err != nil {
			return diag.Errorf("Error retrieving volumes: %s", err)
		}

		volumeList = append(volumeList, volumes...)

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return diag.Errorf("Error retrieving load balancers: %s", err)
		}

		opts.ListOptions.Page = page + 1
	}

	volume, err := findVolumeByName(volumeList, name)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(volume.ID)
	d.Set("name", volume.Name)
	d.Set("urn", volume.URN())
	d.Set("region", volume.Region.Slug)
	d.Set("size", int(volume.SizeGigaBytes))
	d.Set("tags", flattenTags(volume.Tags))

	if v := volume.Description; v != "" {
		d.Set("description", v)
	}

	if v := volume.FilesystemType; v != "" {
		d.Set("filesystem_type", v)
	}
	if v := volume.FilesystemLabel; v != "" {
		d.Set("filesystem_label", v)
	}

	if err = d.Set("droplet_ids", flattenDigitalOceanVolumeDropletIds(volume.DropletIDs)); err != nil {
		return diag.Errorf("[DEBUG] Error setting droplet_ids: %#v", err)
	}

	return nil
}

func findVolumeByName(volumes []godo.Volume, name string) (*godo.Volume, error) {
	results := make([]godo.Volume, 0)
	for _, v := range volumes {
		if v.Name == name {
			results = append(results, v)
		}
	}
	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no volumes found with name %s", name)
	}
	return nil, fmt.Errorf("too many volumes found with name %s (found %d, expected 1)", name, len(results))
}
