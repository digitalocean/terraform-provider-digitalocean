package domain

import (
	"context"
	"fmt"
	"strconv"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func DataSourceDigitalOceanRecord() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDigitalOceanRecordRead,
		Schema: map[string]*schema.Schema{

			"domain": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "domain of the name record",
				ValidateFunc: validation.NoZeroValues,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  "name of the record",
				ValidateFunc: validation.NoZeroValues,
			},
			// computed attributes
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "type of the name record",
			},
			"data": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name record data",
			},
			"priority": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "priority of the name record",
			},
			"port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "port of the name record",
			},
			"ttl": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ttl of the name record",
			},
			"weight": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "weight of the name record",
			},
			"flags": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "flags of the name record",
			},
			"tag": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "tag of the name record",
			},
		},
	}
}

func dataSourceDigitalOceanRecordRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*config.CombinedConfig).GodoClient()
	domain := d.Get("domain").(string)
	name := d.Get("name").(string)

	record, err := findRecordByName(client, domain, name)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(record.ID))
	d.Set("type", record.Type)
	d.Set("name", record.Name)
	d.Set("data", record.Data)
	d.Set("priority", record.Priority)
	d.Set("port", record.Port)
	d.Set("ttl", record.TTL)
	d.Set("weight", record.Weight)
	d.Set("tag", record.Tag)
	d.Set("flags", record.Flags)

	return nil
}

func findRecordByName(client *godo.Client, domain, name string) (*godo.DomainRecord, error) {
	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	for {
		records, resp, err := client.Domains.Records(context.Background(), domain, opts)
		if err != nil {
			if resp != nil && resp.StatusCode == 404 {
				return nil, fmt.Errorf("domain not found: %s", err)
			}
			return nil, fmt.Errorf("error retrieving domain: %s", err)
		}

		for _, r := range records {
			if r.Name == name {
				return &r, nil
			}
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("error retrieving records: %s", err)
		}

		opts.Page = page + 1
	}

	return nil, fmt.Errorf("no records found with name %s", name)
}
