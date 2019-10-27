package digitalocean

import (
	"context"
	"fmt"

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
)

func dataSourceDigitalOceanSizes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanSizeRead,
		Schema: map[string]*schema.Schema{
			"filter": filterSchema(filterKeysDigitalOceanSizes),
			"sort":   sortSchema(filterKeysDigitalOceanSizes),

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
	sizes, err := getAllDigitalOceanSizes(client)
	if err != nil {
		return err
	}

	if v, ok := d.GetOk("filter"); ok {
		filters := expandFilters(v.(*schema.Set).List())
		fmt.Printf("%+v", filters)
	}

	if v, ok := d.GetOk("sort"); ok {
		sorts := expandSorts(v.(*schema.Set).List())
		fmt.Printf("%+v", sorts)
	}

	d.SetId(resource.UniqueId())
	if err := d.Set("sizes", flattenDigitalOceanSizes(sizes)); err != nil {
		return fmt.Errorf("Error setting `sizes`: %+v", err)
	}

	return nil
}

func getAllDigitalOceanSizes(client *godo.Client) ([]godo.Size, error) {
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
