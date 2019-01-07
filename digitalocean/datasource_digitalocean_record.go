package digitalocean

import (
	"context"
	"fmt"
	"strconv"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceDigitalOceanRecord() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanRecordRead,
		Schema: map[string]*schema.Schema{

			"domain": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "domain of the name record",
				ValidateFunc: validation.NoZeroValues,
			},
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "name of the record",
				ValidateFunc: validation.NoZeroValues,
			},
			// computed attributes
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "type of the name record",
			},
			"data": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "name record data",
			},
			"priority": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "priority of the name record",
			},
			"port": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "port of the name record",
			},
			"ttl": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "ttl of the name record",
			},
			"weight": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "weight of the name record",
			},
			"flags": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "flags of the name record",
			},
			"tag": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "tag of the name record",
			},
		},
	}
}

func dataSourceDigitalOceanRecordRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*CombinedConfig).godoClient()
	domain := d.Get("domain").(string)
	name := d.Get("name").(string)

	opts := &godo.ListOptions{}

	records, resp, err := client.Domains.Records(context.Background(), domain, opts)
	if err != nil {
		if resp.StatusCode == 404 {
			return fmt.Errorf("domain not found: %s", err)
		}
		return fmt.Errorf("Error retrieving domain: %s", err)
	}

	record, err := findRecordByName(records, name)
	if err != nil {
		return err
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

func findRecordByName(records []godo.DomainRecord, name string) (*godo.DomainRecord, error) {
	results := make([]godo.DomainRecord, 0)
	for _, v := range records {
		if v.Name == name {
			results = append(results, v)
		}
	}
	if len(results) == 1 {
		return &results[0], nil
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no records found with name %s", name)
	}
	return nil, fmt.Errorf("too many records found (found %d, expected 1)", len(results))
}
