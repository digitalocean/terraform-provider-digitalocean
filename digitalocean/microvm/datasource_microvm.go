package microvm

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// DataSourceDigitalOceanMicroVM returns a data source that looks up a
// MicroVM by `id` or by `name`.
func DataSourceDigitalOceanMicroVM() *schema.Resource {
	recordSchema := microVMDataSourceSchema()
	recordSchema["id"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		Description:  "MicroVM ID",
		ValidateFunc: validation.NoZeroValues,
		ExactlyOneOf: []string{"id", "name"},
	}
	recordSchema["name"] = &schema.Schema{
		Type:         schema.TypeString,
		Optional:     true,
		Computed:     true,
		Description:  "MicroVM name",
		ValidateFunc: validation.NoZeroValues,
		ExactlyOneOf: []string{"id", "name"},
	}

	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanMicroVMRead,
		Schema:      recordSchema,
	}
}

func dataSourceDigitalOceanMicroVMRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	var found *godo.MicroVM
	if id, ok := d.GetOk("id"); ok {
		m, _, err := client.MicroVMs.Get(ctx, id.(string))
		if err != nil {
			return diag.Errorf("Error retrieving MicroVM: %s", err)
		}
		found = m
	} else if name, ok := d.GetOk("name"); ok {
		matches, err := listMicroVMsByName(ctx, client, name.(string))
		if err != nil {
			return diag.Errorf("Error listing MicroVMs: %s", err)
		}
		switch len(matches) {
		case 0:
			return diag.Errorf("no MicroVM found with name %s", name.(string))
		case 1:
			found = &matches[0]
		default:
			return diag.Errorf("too many MicroVMs found with name %s (found %d, expected 1)", name.(string), len(matches))
		}
	}

	if err := setMicroVMAttributes(d, found); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func listMicroVMsByName(ctx context.Context, client *godo.Client, name string) ([]godo.MicroVM, error) {
	opts := &godo.ListOptions{Page: 1, PerPage: 200}
	var out []godo.MicroVM
	for {
		batch, resp, err := client.MicroVMs.ListByName(ctx, name, opts)
		if err != nil {
			return nil, err
		}
		out = append(out, batch...)
		if resp == nil || resp.Links == nil || resp.Links.IsLastPage() {
			break
		}
		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("error paging MicroVMs: %w", err)
		}
		opts.Page = page + 1
	}
	return out, nil
}
