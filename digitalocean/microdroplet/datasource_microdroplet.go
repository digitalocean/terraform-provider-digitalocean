package microdroplet

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// DataSourceDigitalOceanMicroDroplet returns a data source that looks up a
// MicroDroplet by `id` or by `name`.
func DataSourceDigitalOceanMicroDroplet() *schema.Resource {
	recordSchema := microDropletDataSourceSchema()
	recordSchema["id"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		Description:  "MicroDroplet ID",
		ValidateFunc: validation.NoZeroValues,
		ExactlyOneOf: []string{"id", "name"},
	}
	recordSchema["name"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		Description:  "MicroDroplet name",
		ValidateFunc: validation.NoZeroValues,
		ExactlyOneOf: []string{"id", "name"},
	}

	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanMicroDropletRead,
		Schema:      recordSchema,
	}
}

func dataSourceDigitalOceanMicroDropletRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	var found *godo.MicroDroplet
	if id, ok := d.GetOk("id"); ok {
		m, _, err := client.MicroDroplets.Get(ctx, id.(string))
		if err != nil {
			return diag.Errorf("Error retrieving MicroDroplet: %s", err)
		}
		found = m
	} else if name, ok := d.GetOk("name"); ok {
		matches, err := listMicroDropletsByName(ctx, client, name.(string))
		if err != nil {
			return diag.Errorf("Error listing MicroDroplets: %s", err)
		}
		switch len(matches) {
		case 0:
			return diag.Errorf("no MicroDroplet found with name %s", name.(string))
		case 1:
			found = &matches[0]
		default:
			return diag.Errorf("too many MicroDroplets found with name %s (found %d, expected 1)", name.(string), len(matches))
		}
	}

	if err := setMicroDropletAttributes(d, found); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func listMicroDropletsByName(ctx context.Context, client *godo.Client, name string) ([]godo.MicroDroplet, error) {
	opts := &godo.ListOptions{Page: 1, PerPage: 200}
	var out []godo.MicroDroplet
	for {
		batch, resp, err := client.MicroDroplets.ListByName(ctx, name, opts)
		if err != nil {
			return nil, err
		}
		out = append(out, batch...)
		if resp == nil || resp.Links == nil || resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("error paging MicroDroplets: %w", err)
		}
		opts.Page = page + 1
	}
	return out, nil
}
