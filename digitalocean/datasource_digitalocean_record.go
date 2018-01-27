package digitalocean

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceDigitalOceanRecord() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDigitalOceanRecordRead,
		Schema: map[string]*schema.Schema{

			"domain": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "domain of the name record",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "name of the record",
			},
			// computed attributes
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "type of the name record",
			},
			"id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "id of the name record",
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
	client := meta.(*godo.Client)
	domain := d.Get("domain").(string)

	opts := &godo.ListOptions{}

	records, _, err := client.Domains.Records(context.Background(), domain, opts)
	if err != nil {
		d.SetId("")
		return err
	}
	record, err := findRecordByName(records, d.Get("name").(string))

	if err != nil {
		return err
	}

	d.SetId(record.Name)
	d.Set("id", record.ID)
	d.Set("type", record.Type)
	d.Set("name", record.Name)
	d.Set("data", record.Data)
	d.Set("priority", record.Priority)
	d.Set("port", record.Port)
	d.Set("ttl", record.TTL)
	d.Set("weight", record.Weight)

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
