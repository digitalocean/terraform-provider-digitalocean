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
		Read:   dataSourceDigitalOceanRegionsRead,
		Schema: map[string]*schema.Schema{},
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

	regions := make([]string, len(regionList))
	for _, region := range regionList {
		regions = append(regions, region.Slug)
	}

	d.SetId(fmt.Sprintf("%d", hashcode.String(fmt.Sprintf("%+v\n", opts))))
	d.Set("regions", regions)

	return nil
}
