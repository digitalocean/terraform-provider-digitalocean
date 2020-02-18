package digitalocean

import (
	"fmt"

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

	regionsBySlug, err := getDigitalOceanRegions(client)
	if err != nil {
		return fmt.Errorf("Unable to load regions: %s", err)
	}

	slug, ok := d.GetOk("slug")
	if !ok || slug == "" {
		return fmt.Errorf("`slug` property must be specified")
	}

	region, ok := regionsBySlug[slug.(string)]
	if !ok {
		return fmt.Errorf("Region does not exist: %s", slug)
	}

	d.SetId(resource.UniqueId())

	if err := d.Set("slug", region.Slug); err != nil {
		return fmt.Errorf("Unable to set `slug` attribute: %s", err)
	}

	if err := d.Set("name", region.Name); err != nil {
		return fmt.Errorf("Unable to set `name` attribute: %s", err)
	}

	if err := d.Set("available", region.Available); err != nil {
		return fmt.Errorf("Unable to set `available` attribute: %s", err)
	}

	sizesSet := schema.NewSet(schema.HashString, []interface{}{})
	for _, size := range region.Sizes {
		sizesSet.Add(size)
	}
	if err := d.Set("sizes", sizesSet); err != nil {
		return fmt.Errorf("Unable to set `sizes` attribute: %s", err)
	}

	featuresSet := schema.NewSet(schema.HashString, []interface{}{})
	for _, feature := range region.Features {
		featuresSet.Add(feature)
	}
	if err := d.Set("features", featuresSet); err != nil {
		return fmt.Errorf("Unable to set `features` attribute: %s", err)
	}

	return nil
}
