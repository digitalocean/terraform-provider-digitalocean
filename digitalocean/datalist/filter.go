package datalist

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

type commonFilter struct {
	key    string
	values []string
}

func filterSchema(allowedKeys []string) *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeSet,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(allowedKeys, false),
				},
				"values": {
					Type:     schema.TypeList,
					Required: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},
			},
		},
		Optional:    true,
		Description: "One or more key/value pairs to filter droplet size results.",
	}
}

func expandFilters(rawFilters []interface{}) []commonFilter {
	expandedFilters := make([]commonFilter, len(rawFilters))
	for i, rawFilter := range rawFilters {
		f := rawFilter.(map[string]interface{})

		expandedFilter := commonFilter{
			key:    f["key"].(string),
			values: expandFilterValues(f["values"].([]interface{})),
		}

		expandedFilters[i] = expandedFilter
	}
	return expandedFilters
}

func expandFilterValues(rawFilterValues []interface{}) []string {
	expandedFilterValues := make([]string, len(rawFilterValues))
	for i, v := range rawFilterValues {
		expandedFilterValues[i] = v.(string)
	}

	return expandedFilterValues
}

func applyFilters(recordSchema map[string]*schema.Schema, records []map[string]interface{}, filters []commonFilter) []map[string]interface{} {
	for _, f := range filters {
		// Handle multiple filters by applying them in order
		var filteredRecords []map[string]interface{}

		filterFunc := func(record map[string]interface{}) bool {
			result := false

			for _, filterValue := range f.values {
				result = result || valueMatches(recordSchema[f.key], record[f.key], filterValue)
			}

			return result
		}

		for _, record := range records {
			if filterFunc(record) {
				filteredRecords = append(filteredRecords, record)
			}
		}

		records = filteredRecords
	}

	return records
}
