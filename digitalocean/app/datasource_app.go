package app

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanApp() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanAppRead,
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"spec": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "A DigitalOcean App Platform Spec",
				Elem: &schema.Resource{
					Schema: appSpecSchema(false),
				},
			},
			"default_ingress": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The default URL to access the App",
			},
			"live_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active_deployment_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID the App's currently active deployment",
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The uniform resource identifier for the app",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of when the App was last updated",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The date and time of when the App was created",
			},
		},
	}
}

func dataSourceDigitalOceanAppRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	d.SetId(d.Get("app_id").(string))

	return resourceDigitalOceanAppRead(ctx, d, meta)
}
