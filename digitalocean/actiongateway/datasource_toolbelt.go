package actiongateway

import (
	"context"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// DataSourceDigitalOceanActionGatewayToolbelt returns a data source that
// looks up the latest or a specific version of an Action Gateway toolbelt.
func DataSourceDigitalOceanActionGatewayToolbelt() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanActionGatewayToolbeltRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "Name of the toolbelt",
				ValidateFunc: validation.NoZeroValues,
			},
			"version": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "Version to read. Defaults to the latest version.",
				ValidateFunc: validation.NoZeroValues,
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tools": {
				Type:        schema.TypeSet,
				Computed:    true,
				Description: "Members of the version, each pinned as <tool_slug>@<version>",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"reference": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reference_latest": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tool_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDigitalOceanActionGatewayToolbeltRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	name := d.Get("name").(string)
	opt := &godo.ActionGatewayGetToolbeltOptions{
		Version:  d.Get("version").(string),
		PageSize: 1,
	}
	found, _, err := client.ActionGateway.Toolbelts.Get(ctx, name, opt)
	if err != nil {
		return diag.Errorf("Error retrieving Action Gateway toolbelt (%s): %s", name, err)
	}
	if found.Toolbelt == nil {
		return diag.Errorf("Error retrieving Action Gateway toolbelt (%s): response contained no toolbelt", name)
	}

	d.SetId(found.Toolbelt.Reference)
	for k, v := range flattenToolbelt(found.Toolbelt) {
		if err := d.Set(k, v); err != nil {
			return diag.Errorf("error setting %s: %s", k, err)
		}
	}
	return nil
}
