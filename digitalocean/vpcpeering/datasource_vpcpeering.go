package vpcpeering

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDigitalOceanVPCPeering() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanVPCPeeringRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The ID of the VPC Peering",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Description:  "The name of the VPC Peering",
				ValidateFunc: validation.NoZeroValues,
				ExactlyOneOf: []string{"id", "name"},
			},
			"vpc_ids": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: "The list of VPCs to be peered",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDigitalOceanVPCPeeringRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	var foundVPCPeering *godo.VPCPeering

	if id, ok := d.GetOk("id"); ok {
		vpcPeering, _, err := client.VPCs.GetVPCPeering(context.Background(), id.(string))
		if err != nil {
			return diag.Errorf("Error retrieving VPC Peering: %s", err)
		}

		foundVPCPeering = vpcPeering
	} else if name, ok := d.GetOk("name"); ok {
		vpcPeerings, err := listVPCPeerings(client)
		if err != nil {
			return diag.Errorf("Error retrieving VPC Peering: %s", err)
		}

		vpcPeering, err := findVPCPeeringByName(vpcPeerings, name.(string))
		if err != nil {
			return diag.Errorf("Error retrieving VPC Peering: %s", err)
		}

		foundVPCPeering = vpcPeering
	}

	d.SetId(foundVPCPeering.ID)
	d.Set("name", foundVPCPeering.Name)
	d.Set("vpc_ids", foundVPCPeering.VPCIDs)
	d.Set("status", foundVPCPeering.Status)
	d.Set("created_at", foundVPCPeering.CreatedAt.UTC().String())

	return nil
}

func listVPCPeerings(client *godo.Client) ([]*godo.VPCPeering, error) {
	peeringsList := []*godo.VPCPeering{}
	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	for {
		peerings, resp, err := client.VPCs.ListVPCPeerings(context.Background(), opts)

		if err != nil {
			return peeringsList, fmt.Errorf("error retrieving VPC Peerings: %s", err)
		}

		peeringsList = append(peeringsList, peerings...)

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return peeringsList, fmt.Errorf("error retrieving VPC Peerings: %s", err)
		}

		opts.Page = page + 1
	}

	return peeringsList, nil
}

func findVPCPeeringByName(vpcPeerings []*godo.VPCPeering, name string) (*godo.VPCPeering, error) {
	results := make([]*godo.VPCPeering, 0)
	for _, v := range vpcPeerings {
		if v.Name == name {
			results = append(results, v)
		}
	}
	if len(results) == 1 {
		return results[0], nil
	} else if len(results) == 0 {
		return nil, fmt.Errorf("no VPC Peerings found with name %s", name)
	}

	return nil, fmt.Errorf("too many VPC Peerings found with name %s (found %d, expected 1)", name, len(results))
}
