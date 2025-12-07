package nfs

import (
	"context"
	"fmt"
	"strings"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/tag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDigitalOceanNfs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanNfsRead,
		Schema: map[string]*schema.Schema{

			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "name of the share",
				ValidateFunc: validation.NoZeroValues,
			},
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "the region that the share is created in",
				StateFunc: func(val interface{}) string {
					// DO API V2 region slug is always lowercase
					return strings.ToLower(val.(string))
				},
				ValidateFunc: validation.NoZeroValues,
			},
			"size": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "the size of the share in gigabytes",
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mount_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": tag.TagsDataSourceSchema(),
		},
	}
}

func dataSourceDigitalOceanNfsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()

	name := d.Get("name").(string)

	region := d.Get("region").(string)

	opts := &godo.ListOptions{
		PerPage: 200,
	}

	if s, ok := d.GetOk("region"); ok {
		region = s.(string)
	}

	sharesList := []*godo.Nfs{}

	for {
		shares, resp, err := client.Nfs.List(context.Background(), opts, region)

		if err != nil {
			return diag.Errorf("Error retrieving shares: %s", err)
		}

		sharesList = append(sharesList, shares...)

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return diag.Errorf("Error retrieving shares: %s", err)
		}

		opts.Page = page + 1
	}

	share, err := findShareByName(sharesList, name)

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(share.ID)
	d.Set("name", share.Name)
	d.Set("region", share.Region)
	d.Set("size", share.SizeGib)
	d.Set("status", share.Status)
	d.Set("mount_path", share.MountPath)

	return nil
}

func findShareByName(shares []*godo.Nfs, name string) (*godo.Nfs, error) {
	results := make([]*godo.Nfs, 0)
	for _, s := range shares {
		if s != nil && s.Name == name {
			results = append(results, s)
		}
	}
	if len(results) == 1 {
		return results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no shares found with name %s", name)
	}
	return nil, fmt.Errorf("too many shares found with name %s (found %d, expected 1)", name, len(results))
}
