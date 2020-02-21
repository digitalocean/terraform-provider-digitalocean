package digitalocean

import (
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func dataSourceDigitalOceanRegion() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanRegionRead,
		Schema: map[string]*schema.Schema{
			"slug": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
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

	slug := d.Get("slug").(string)

	var regionForSlug *godo.Region
	for _, region := range regions {
		if region.Slug == slug {
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
