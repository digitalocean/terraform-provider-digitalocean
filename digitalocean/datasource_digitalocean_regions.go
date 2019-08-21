package digitalocean

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDigitalOceanRegions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanRegionsRead,
		Schema: map[string]*schema.Schema{
			"regions": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slug": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sizes": {
							Elem:     &schema.Schema{Type: schema.TypeString},
							Computed: true,
						},
						"features": {
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

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	regionList := []godo.Region{}
	for {
		regions, resp, err := client.Regions.List(context.Background(), opts)

		if err != nil {
			return fmt.Errorf("Error retrieving ssh keys: %s", err)
		}

		for _, region := range regions {
			regionList = append(regionList, region)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return fmt.Errorf("Error retrieving ssh keys: %s", err)
		}

		opts.Page = page + 1
	}

	d.SetId(fmt.Sprintf("%d", hashcode.String(fmt.Sprintf("%+v\n", opts))))
	d.Set("regions", regionsList)

	return nil
}
