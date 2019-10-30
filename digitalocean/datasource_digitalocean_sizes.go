package digitalocean

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

var (
	filterKeysDigitalOceanSizes = []string{
		"slug",
		"regions",
		"memory",
		"vcpus",
		"disk",
		"transfer",
		"price_monthly",
		"price_hourly",
		"available",
	}

	sortKeysDigitalOceanSizes = []string{
		"slug",
		"memory",
		"vcpus",
		"disk",
		"transfer",
		"price_monthly",
		"price_hourly",
	}
)

func dataSourceDigitalOceanSizes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanSizeRead,
		Schema: map[string]*schema.Schema{
			"filter": filterSchema(filterKeysDigitalOceanSizes),
			"sort":   sortSchema(sortKeysDigitalOceanSizes),

			// Computed properties
			"sizes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"slug": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "A human-readable string that is used to uniquely identify each size.",
						},
						"available": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "This represents whether new Droplets can be created with this size.",
						},
						"transfer": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The amount of transfer bandwidth that is available for Droplets created in this size. This only counts traffic on the public interface. The value is given in terabytes.",
						},
						"price_monthly": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The monthly cost of Droplets created in this size if they are kept for an entire month. The value is measured in US dollars.",
						},
						"price_hourly": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: "The hourly cost of Droplets created in this size as measured hourly. The value is measured in US dollars.",
						},
						"memory": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The amount of RAM allocated to Droplets created of this size. The value is measured in megabytes.",
						},
						"vcpus": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The number of CPUs allocated to Droplets of this size.",
						},
						"disk": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The amount of disk space set aside for Droplets of this size. The value is measured in gigabytes.",
						},
						"regions": {
							Type:        schema.TypeSet,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: "List of region slugs where Droplets can be created in this size.",
						},
					},
				},
				Description: "List of filtered digital ocean sizes.",
			},
		},
	}
}

func dataSourceDigitalOceanSizeRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
	sizes, err := getDigitalOceanSizes(client)
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("filter"); ok {
		filters := expandFilters(v.(*schema.Set).List())
		sizes = filterDigitalOceanSizes(sizes, filters)
	}

	if v, ok := d.GetOk("sort"); ok {
		sorts := expandSorts(v.([]interface{}))
		sizes = sortDigitalOceanSizes(sizes, sorts)
	}

	d.SetId(resource.UniqueId())
	if err := d.Set("sizes", flattenDigitalOceanSizes(sizes)); err != nil {
		return fmt.Errorf("Error setting `sizes`: %+v", err)
	}

	return nil
}

func getDigitalOceanSizes(client *godo.Client) ([]godo.Size, error) {
	sizes := []godo.Size{}

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	for {
		partialSizes, resp, err := client.Sizes.List(context.Background(), opts)
		if err != nil {
			return nil, fmt.Errorf("Error retrieving sizes: %s", err)
		}

		for _, partialSize := range partialSizes {
			sizes = append(sizes, partialSize)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("Error retrieving sizes: %s", err)
		}

		opts.Page = page + 1
	}

	return sizes, nil
}

func filterDigitalOceanSizes(sizes []godo.Size, filters []commonFilter) []godo.Size {
	for _, f := range filters {
		// Handle multiple filters by applying them in order
		var filteredSizes []godo.Size

		// Define the filter strategy based on the filter key
		filterFunc := func(size godo.Size) bool {
			result := false

			for _, filterValue := range f.values {
				switch f.key {
				case "slug":
					result = result || strings.EqualFold(filterValue, size.Slug)
				case "regions":
					for _, region := range size.Regions {
						result = result || strings.EqualFold(filterValue, region)
					}
				case "memory":
					if memory, err := strconv.Atoi(filterValue); err == nil {
						result = result || memory == size.Memory
					}
				case "vcpus":
					if vcpus, err := strconv.Atoi(filterValue); err == nil {
						result = result || vcpus == size.Vcpus
					}
				case "disk":
					if disk, err := strconv.Atoi(filterValue); err == nil {
						result = result || disk == size.Disk
					}
				case "transfer":
					if transfer, err := strconv.ParseFloat(filterValue, 64); err == nil {
						result = result || fmt.Sprintf("%.5f", transfer) == fmt.Sprintf("%.5f", size.Transfer)
					}
				case "price_monthly":
					if priceMonthly, err := strconv.ParseFloat(filterValue, 64); err == nil {
						result = result || fmt.Sprintf("%.5f", priceMonthly) == fmt.Sprintf("%.5f", size.PriceMonthly)
					}
				case "price_hourly":
					if priceHourly, err := strconv.ParseFloat(filterValue, 64); err == nil {
						result = result || fmt.Sprintf("%.5f", priceHourly) == fmt.Sprintf("%.5f", size.PriceHourly)
					}
				case "available":
					if available, err := strconv.ParseBool(filterValue); err == nil {
						result = result || available == size.Available
					}
				default:
				}
			}

			return result
		}

		for _, size := range sizes {
			if filterFunc(size) {
				filteredSizes = append(filteredSizes, size)
			}
		}

		sizes = filteredSizes
	}

	return sizes
}

func sortDigitalOceanSizes(sizes []godo.Size, sorts []commonSort) []godo.Size {
	sort.Slice(sizes, func(_i, _j int) bool {
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
				if sizes[i].Slug != sizes[j].Slug {
					return sizes[i].Slug < sizes[j].Slug
				}
			case "memory":
				if sizes[i].Memory != sizes[j].Memory {
					return sizes[i].Memory < sizes[j].Memory
				}
			case "vcpus":
				if sizes[i].Vcpus != sizes[j].Vcpus {
					return sizes[i].Vcpus < sizes[j].Vcpus
				}
			case "disk":
				if sizes[i].Disk != sizes[j].Disk {
					return sizes[i].Disk < sizes[j].Disk
				}
			case "transfer":
				if sizes[i].Transfer != sizes[j].Transfer {
					return sizes[i].Transfer < sizes[j].Transfer
				}
			case "price_monthly":
				if sizes[i].PriceMonthly != sizes[j].PriceMonthly {
					return sizes[i].PriceMonthly < sizes[j].PriceMonthly
				}
			case "price_hourly":
				if sizes[i].PriceHourly != sizes[j].PriceHourly {
					return sizes[i].PriceHourly < sizes[j].PriceHourly
				}
			}
		}

		return true
	})

	return sizes
}

func flattenDigitalOceanSizes(sizes []godo.Size) []interface{} {
	flattenedSizes := make([]interface{}, len(sizes))

	for i, s := range sizes {
		flattenedSize := make(map[string]interface{})
		flattenedSize["slug"] = s.Slug
		flattenedSize["available"] = s.Available
		flattenedSize["transfer"] = s.Transfer
		flattenedSize["price_monthly"] = s.PriceMonthly
		flattenedSize["price_hourly"] = s.PriceHourly
		flattenedSize["memory"] = s.Memory
		flattenedSize["vcpus"] = s.Vcpus
		flattenedSize["disk"] = s.Disk

		flattenedRegions := schema.NewSet(schema.HashString, []interface{}{})
		for _, r := range s.Regions {
			flattenedRegions.Add(r)
		}
		flattenedSize["regions"] = flattenedRegions

		flattenedSizes[i] = flattenedSize
	}

	return flattenedSizes
}
