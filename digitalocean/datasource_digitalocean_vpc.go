package digitalocean

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceDigitalOceanVPC() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanVPCRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name", "region"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name", "region"},
			},
			"region": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name", "region"},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_range": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDigitalOceanVPCRead(ctx context.Context, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
	var foundVPC *godo.VPC

	if id, ok := d.GetOk("id"); ok {
		vpc, _, err := client.VPCs.Get(context.Background(), id.(string))
		if err != nil {
			return fmt.Errorf("Error retrieving VPC: %s", err)
		}

		foundVPC = vpc
	} else if slug, ok := d.GetOk("region"); ok {
		vpcs, err := listVPCs(client)
		if err != nil {
			return fmt.Errorf("Error retrieving VPC: %s", err)
		}

		vpc, err := findRegionDefaultVPC(vpcs, slug.(string))
		if err != nil {
			return fmt.Errorf("Error retrieving VPC: %s", err)
		}

		foundVPC = vpc
	} else if name, ok := d.GetOk("name"); ok {
		vpcs, err := listVPCs(client)
		if err != nil {
			return fmt.Errorf("Error retrieving VPC: %s", err)
		}

		vpc, err := findVPCByName(vpcs, name.(string))
		if err != nil {
			return fmt.Errorf("Error retrieving VPC: %s", err)
		}

		foundVPC = vpc
	}

	d.SetId(foundVPC.ID)
	d.Set("name", foundVPC.Name)
	d.Set("region", foundVPC.RegionSlug)
	d.Set("description", foundVPC.Description)
	d.Set("ip_range", foundVPC.IPRange)
	d.Set("urn", foundVPC.URN)
	d.Set("default", foundVPC.Default)
	d.Set("created_at", foundVPC.CreatedAt.UTC().String())

	return nil
}

func listVPCs(client *godo.Client) ([]*godo.VPC, error) {
	vpcList := []*godo.VPC{}
	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	for {
		vpcs, resp, err := client.VPCs.List(context.Background(), opts)

		if err != nil {
			return vpcList, fmt.Errorf("Error retrieving VPCs: %s", err)
		}

		for _, v := range vpcs {
			vpcList = append(vpcList, v)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return vpcList, fmt.Errorf("Error retrieving VPCs: %s", err)
		}

		opts.Page = page + 1
	}

	return vpcList, nil
}

func findRegionDefaultVPC(vpcs []*godo.VPC, region string) (*godo.VPC, error) {
	results := make([]*godo.VPC, 0)
	for _, v := range vpcs {
		if v.RegionSlug == region && v.Default {
			results = append(results, v)
		}
	}
	if len(results) == 1 {
		return results[0], nil
	}

	return nil, fmt.Errorf("unable to find default VPC in %s region", region)
}

func findVPCByName(vpcs []*godo.VPC, name string) (*godo.VPC, error) {
	results := make([]*godo.VPC, 0)
	for _, v := range vpcs {
		if v.Name == name {
			results = append(results, v)
		}
	}
	if len(results) == 1 {
		return results[0], nil
	} else if len(results) == 0 {
		return nil, fmt.Errorf("no VPCs found with name %s", name)
	}

	return nil, fmt.Errorf("too many VPCs found with name %s (found %d, expected 1)", name, len(results))
}
