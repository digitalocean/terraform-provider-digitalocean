package digitalocean

import (
	"fmt"

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
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of region slugs where Droplets can be created in this size.",
			},
		},
	}
}

func dataSourceDigitalOceanSizeRead(d *schema.ResourceData, meta interface{}) error {
	if v, ok := d.GetOk("filter"); ok {
		filters := expandFilters(v.(*schema.Set).List())
		fmt.Printf("%+v", filters)
	}

	if v, ok := d.GetOk("sort"); ok {
		sorts := expandSorts(v.(*schema.Set).List())
		fmt.Printf("%+v", sorts)
	}

	return nil
}
