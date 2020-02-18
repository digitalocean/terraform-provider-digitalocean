package digitalocean

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceDigitalOceanRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanRegionsRead,
		Schema: map[string]*schema.Schema{
			"available": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"features": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"slugs": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceDigitalOceanRegionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	regionsBySlug, err := loadRegionsBySlug(client)
	if err != nil {
		return fmt.Errorf("Unable to load regions: %s", err)
	}

	available, filterByAvailable := d.GetOk("available")
	requiredFeatures, filterByRequiredFeatures := d.GetOk("features")

	filteredRegionSlugs := []string{}
	for _, region := range regionsBySlug {
		if filterByAvailable && region.Available != available.(bool) {
			continue
		}

		if filterByRequiredFeatures {
			match := false
			for _, requiredFeature := range requiredFeatures.([]interface{}) {
				for _, feature := range region.Features {
					if feature == requiredFeature.(string) {
						match = true
					}
				}
			}
			if !match {
				continue
			}
		}

		filteredRegionSlugs = append(filteredRegionSlugs, region.Slug)
	}

	d.SetId(resource.UniqueId())
	if err := d.Set("slugs", filteredRegionSlugs); err != nil {
		return fmt.Errorf("Unable to set `slugs` attribute: %s", err)
	}

	return nil
}

func loadRegionsBySlug(client *godo.Client) (map[string]godo.Region, error) {
	allRegions := map[string]godo.Region{}

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	for {
		regions, resp, err := client.Regions.List(context.Background(), opts)

		if err != nil {
			return nil, fmt.Errorf("Error retrieving regions: %s", err)
		}

		for _, region := range regions {
			allRegions[region.Slug] = region
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("Error retrieving regions: %s", err)
		}

		opts.Page = page + 1
	}

	return allRegions, nil
}
