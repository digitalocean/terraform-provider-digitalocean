package domain

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func recordsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Type:        schema.TypeInt,
			Description: "ID of the record",
		},
		"domain": {
			Type:        schema.TypeString,
			Description: "domain of the name record",
		},
		"name": {
			Type:        schema.TypeString,
			Description: "name of the record",
		},
		"type": {
			Type:        schema.TypeString,
			Description: "type of the name record",
		},
		"value": {
			Type:        schema.TypeString,
			Description: "name record data",
		},
		"priority": {
			Type:        schema.TypeInt,
			Description: "priority of the name record",
		},
		"port": {
			Type:        schema.TypeInt,
			Description: "port of the name record",
		},
		"ttl": {
			Type:        schema.TypeInt,
			Description: "ttl of the name record",
		},
		"weight": {
			Type:        schema.TypeInt,
			Description: "weight of the name record",
		},
		"flags": {
			Type:        schema.TypeInt,
			Description: "flags of the name record",
		},
		"tag": {
			Type:        schema.TypeString,
			Description: "tag of the name record",
		},
	}
}

func getDigitalOceanRecords(meta interface{}, extra map[string]interface{}) ([]interface{}, error) {
	client := meta.(*config.CombinedConfig).GodoClient()

	domain, ok := extra["domain"].(string)
	if !ok {
		return nil, fmt.Errorf("unable to find `domain` key from query data")
	}

	var allRecords []interface{}

	opts := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	for {
		records, resp, err := client.Domains.Records(context.Background(), domain, opts)
		if err != nil {
			return nil, fmt.Errorf("Error retrieving records: %s", err)
		}

		for _, record := range records {
			allRecords = append(allRecords, record)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("Error retrieving projects: %s", err)
		}

		opts.Page = page + 1
	}

	return allRecords, nil
}

func flattenDigitalOceanRecord(rawRecord interface{}, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error) {
	domain, ok := extra["domain"].(string)
	if !ok {
		return nil, fmt.Errorf("unable to find `domain` key from query data")
	}

	record, ok := rawRecord.(godo.DomainRecord)
	if !ok {
		return nil, fmt.Errorf("unable to convert to godo.DomainRecord")
	}

	flattenedRecord := map[string]interface{}{
		"id":       record.ID,
		"domain":   domain,
		"name":     record.Name,
		"type":     record.Type,
		"value":    record.Data,
		"priority": record.Priority,
		"port":     record.Port,
		"ttl":      record.TTL,
		"weight":   record.Weight,
		"flags":    record.Flags,
		"tag":      record.Tag,
	}

	return flattenedRecord, nil
}
