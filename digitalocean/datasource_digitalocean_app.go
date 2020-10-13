package digitalocean

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceDigitalOceanApp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanAppRead,
		Schema: map[string]*schema.Schema{
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"spec": {
				Type:        schema.TypeList,
				Computed:    true,
				MaxItems:    1,
				Description: "A DigitalOcean App Platform Spec",
				Elem: &schema.Resource{
					Schema: appSpecSchema(),
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

func dataSourceDigitalOceanAppRead(d *schema.ResourceData, meta interface{}) error {
	d.SetId(d.Get("app_id").(string))

	return resourceDigitalOceanAppRead(d, meta)
}
