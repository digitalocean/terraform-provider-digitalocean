package util

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
