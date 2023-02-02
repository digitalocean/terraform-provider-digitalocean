package size

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/internal/datalist"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceDigitalOceanSizes() *schema.Resource {
	dataListConfig := &datalist.ResourceConfig{
		RecordSchema: map[string]*schema.Schema{
			"slug": {
				Type:        schema.TypeString,
				Description: "A human-readable string that is used to uniquely identify each size.",
			},
			"available": {
				Type:        schema.TypeBool,
				Description: "This represents whether new Droplets can be created with this size.",
			},
			"transfer": {
				Type:        schema.TypeFloat,
				Description: "The amount of transfer bandwidth that is available for Droplets created in this size. This only counts traffic on the public interface. The value is given in terabytes.",
			},
			"price_monthly": {
				Type:        schema.TypeFloat,
				Description: "The monthly cost of Droplets created in this size if they are kept for an entire month. The value is measured in US dollars.",
			},
			"price_hourly": {
				Type:        schema.TypeFloat,
				Description: "The hourly cost of Droplets created in this size as measured hourly. The value is measured in US dollars.",
			},
			"memory": {
				Type:        schema.TypeInt,
				Description: "The amount of RAM allocated to Droplets created of this size. The value is measured in megabytes.",
			},
			"vcpus": {
				Type:        schema.TypeInt,
				Description: "The number of CPUs allocated to Droplets of this size.",
			},
			"disk": {
				Type:        schema.TypeInt,
				Description: "The amount of disk space set aside for Droplets of this size. The value is measured in gigabytes.",
			},
			"regions": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of region slugs where Droplets can be created in this size.",
			},
		},
		ResultAttributeName: "sizes",
		FlattenRecord:       flattenDigitalOceanSize,
		GetRecords:          getDigitalOceanSizes,
	}

	return datalist.NewResource(dataListConfig)
}

func getDigitalOceanSizes(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.CombinedConfig).GodoClient()

	sizes := []interface{}{}

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

func flattenDigitalOceanSize(size, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	s := size.(godo.Size)

	flattenedSize := map[string]interface{}{}
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

	return flattenedSize, nil
}
