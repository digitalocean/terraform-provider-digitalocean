package microdroplet

import (
	"context"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// DataSourceDigitalOceanMicroDropletImage returns a data source that looks up
// a MicroDroplet image by `id` or by `name`. Name lookup requires paging the
// full image list since the API has no `ListByName` for images.
func DataSourceDigitalOceanMicroDropletImage() *schema.Resource {
	recordSchema := microDropletImageDataSourceSchema()
	recordSchema["id"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		Description:  "MicroDroplet image ID",
		ValidateFunc: validation.NoZeroValues,
		ExactlyOneOf: []string{"id", "name"},
	}
	recordSchema["name"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		Description:  "MicroDroplet image name",
		ValidateFunc: validation.NoZeroValues,
		ExactlyOneOf: []string{"id", "name"},
	}

	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanMicroDropletImageRead,
		Schema:      recordSchema,
	}
}

func dataSourceDigitalOceanMicroDropletImageRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	var found *godo.MicroDropletImage
	if id, ok := d.GetOk("id"); ok {
		img, _, err := client.MicroDropletImages.Get(ctx, id.(string))
		if err != nil {
			return diag.Errorf("Error retrieving MicroDroplet image: %s", err)
		}
		found = img
	} else if name, ok := d.GetOk("name"); ok {
		images, err := listAllMicroDropletImages(ctx, client)
		if err != nil {
			return diag.Errorf("Error listing MicroDroplet images: %s", err)
		}
		var matches []godo.MicroDropletImage
		for _, img := range images {
			if img.Name == name.(string) {
				matches = append(matches, img)
			}
		}
		switch len(matches) {
		case 0:
			return diag.Errorf("no MicroDroplet image found with name %s", name.(string))
		case 1:
			found = &matches[0]
		default:
			return diag.Errorf("too many MicroDroplet images found with name %s (found %d, expected 1)", name.(string), len(matches))
		}
	}

	d.SetId(found.ID)
	d.Set("name", found.Name)
	d.Set("source", found.Source)
	d.Set("status", string(found.Status))
	d.Set("created_at", found.Created)
	d.Set("urn", found.URN())
	return nil
}

func listAllMicroDropletImages(ctx context.Context, client *godo.Client) ([]godo.MicroDropletImage, error) {
	opts := &godo.ListOptions{Page: 1, PerPage: 200}
	var all []godo.MicroDropletImage
	for {
		batch, resp, err := client.MicroDropletImages.List(ctx, opts)
		if err != nil {
			return nil, err
		}
		all = append(all, batch...)
		if resp == nil || resp.Links == nil || resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}
		opts.Page = page + 1
	}
	return all, nil
}
