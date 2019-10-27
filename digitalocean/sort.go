package digitalocean

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

var (
	sortKeys = []string{"asc", "desc"}
)

type commonSort struct {
	key       string
	direction string
}

func sortSchema(allowedKeys []string) *schema.Schema {
	return &schema.Schema{
		Type: schema.TypeSet,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.StringInSlice(allowedKeys, false),
				},
				"direction": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringInSlice(sortKeys, false),
				},
			},
		},
		Optional:    true,
		Description: "One or more key/direction pairs to sort droplet size results.",
	}
}

func expandSorts(rawSorts []interface{}) []commonSort {
	expandedSorts := make([]commonSort, len(rawSorts))
	for i, rawSort := range rawSorts {
		f := rawSort.(map[string]interface{})

		expandedSort := commonSort{
			key:       f["key"].(string),
			direction: f["direction"].(string),
		}

		expandedSorts[i] = expandedSort
	}
	return expandedSorts
}
