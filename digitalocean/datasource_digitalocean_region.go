package digitalocean

import (
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceDigitalOceanRegion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanRegionRead,
		Schema: map[string]*schema.Schema{
			"slug": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sizes": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"features": {
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"available": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceDigitalOceanRegionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	regions, err := getDigitalOceanRegions(client)
	if err != nil {
		return fmt.Errorf("Unable to load regions: %s", err)
	}

	slug, ok := d.GetOk("slug")
	if !ok || slug == "" {
		return fmt.Errorf("`slug` property must be specified")
	}

	var regionForSlug *godo.Region
	for _, region := range regions {
		if region.Slug == slug.(string) {
			regionForSlug = &region
			break
		}
	}

	if regionForSlug == nil {
		return fmt.Errorf("Region does not exist: %s", slug)
	}

	d.SetId(resource.UniqueId())
	if err := setResourceDataFromMap(d, flattenRegion(*regionForSlug)); err != nil {
		return err
	}

	return nil
}
