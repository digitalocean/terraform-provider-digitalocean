package digitalocean

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var (
	dataSourceDigitalOceanRegionsFilterKeys = []string{
		"slug",
		"name",
		"available",
		"features",
		"sizes",
	}

	dataSourceDigitalOceanRegionsSortKeys = []string{
		"slug",
		"name",
		"available",
	}
)

func dataSourceDigitalOceanRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanRegionsRead,
		Schema: map[string]*schema.Schema{
			"filter": filterSchema(dataSourceDigitalOceanRegionsFilterKeys),
			"sort":   sortSchema(dataSourceDigitalOceanRegionsSortKeys),

			"regions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
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
				},
			},
		},
	}
}

func dataSourceDigitalOceanRegionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()

	regions, err := getDigitalOceanRegions(client)
	if err != nil {
		return fmt.Errorf("Unable to load regions: %s", err)
	}

	if v, ok := d.GetOk("filter"); ok {
		filters := expandFilters(v.(*schema.Set).List())
		regions = filterDigitalOceanRegions(regions, filters)
	}

	if v, ok := d.GetOk("sort"); ok {
		sorts := expandSorts(v.([]interface{}))
		regions = sortDigitalOceanRegions(regions, sorts)
	}

	d.SetId(resource.UniqueId())

	flattenedRegions := []map[string]interface{}{}
	for _, region := range regions {
		flattenedRegions = append(flattenedRegions, flattenRegion(region))
	}

	if err := d.Set("regions", flattenedRegions); err != nil {
		return fmt.Errorf("Unable to set `regions` attribute: %s", err)
	}

	return nil
}

func filterDigitalOceanRegions(regions []godo.Region, filters []commonFilter) []godo.Region {
	for _, f := range filters {
		// Handle multiple filters by applying them in order
		var filteredRegions []godo.Region

		// Define the filter strategy based on the filter key
		filterFunc := func(region godo.Region) bool {
			result := false

			for _, filterValue := range f.values {
				switch f.key {
				case "slug":
					result = result || strings.EqualFold(filterValue, region.Slug)

				case "name":
					result = result || strings.EqualFold(filterValue, region.Name)

				case "available":
					if available, err := strconv.ParseBool(filterValue); err == nil {
						result = result || available == region.Available
					}

				case "features":
					for _, feature := range region.Features {
						result = result || strings.EqualFold(filterValue, feature)
					}

				case "sizes":
					for _, size := range region.Sizes {
						result = result || strings.EqualFold(filterValue, size)
					}

				default:
				}
			}

			return result
		}

		for _, region := range regions {
			if filterFunc(region) {
				filteredRegions = append(filteredRegions, region)
			}
		}

		regions = filteredRegions
	}

	return regions
}

func sortDigitalOceanRegions(regions []godo.Region, sorts []commonSort) []godo.Region {
	sort.Slice(regions, func(_i, _j int) bool {
		for _, s := range sorts {
			// Handle multiple sorts by applying them in order

			i := _i
			j := _j
			if strings.EqualFold(s.direction, "desc") {
				// If the direction is desc, reverse index to compare
				i = _j
				j = _i
			}

			switch s.key {
			case "slug":
				if regions[i].Slug != regions[j].Slug {
					return strings.Compare(regions[i].Slug, regions[j].Slug) < 0
				}

			case "name":
				if regions[i].Name != regions[j].Name {
					return strings.Compare(regions[i].Name, regions[j].Name) < 0
				}
			case "available":
				if regions[i].Available != regions[j].Available {
					return !regions[i].Available && regions[j].Available
				}
			}
		}

		return true
	})

	return regions
}
