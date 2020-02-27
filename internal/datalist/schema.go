package datalist

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

// This is the configuration for a "data list" resource. It represents the schema and operations
// needed to create the data list resource.
type ResourceConfig struct {
	// The schema for a single instance of the resource.
	RecordSchema map[string]*schema.Schema

	// A string slice with the attribute keys on which the filter attribute can operate.
	// The filter attribute operates on all field types except for schema.TypeMap.
	FilterKeys []string

	// A string slice with the attribute keys on which the sort attribute can operate.
	// The sort attribute only operates on the schema.TypeString, schema.TypeBool,
	// schema.TypeInt, and schema.TypeFloat field types.
	SortKeys []string

	// The name of the attribute in the resource through which to expose results.
	ResultAttributeName string

	// Given a record returned from the GetRecords function, flatten the record to a
	// map acceptable to the Set method on schema.ResourceData.
	FlattenRecord func(record interface{}) map[string]interface{}

	// Return all of the records on which the data list resource should operate.
	// The `meta` argument is the same meta argument passed into the resource's Read
	// function.
	GetRecords func(meta interface{}) ([]interface{}, error)
}

// Returns a new "data list" resource given the specified configuration. This
// is a resource with `filter` and `sort` attributes that can select a subset
// of records from a list of records for a particular type of resource.
func NewResource(config *ResourceConfig) *schema.Resource {
	resultSchema := map[string]*schema.Schema{}
	for key, value := range config.RecordSchema {
		resultSchema[key] = &schema.Schema{
			Type:     value.Type,
			Elem:     value.Elem,
			Computed: true,
		}
	}

	return &schema.Resource{
		Read: dataListResourceRead(config),
		Schema: map[string]*schema.Schema{
			"filter": filterSchema(config.FilterKeys),
			"sort":   sortSchema(config.SortKeys),
			config.ResultAttributeName: {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: resultSchema,
				},
			},
		},
	}
}

func dataListResourceRead(config *ResourceConfig) schema.ReadFunc {
	return func(d *schema.ResourceData, meta interface{}) error {
		records, err := config.GetRecords(meta)
		if err != nil {
			return fmt.Errorf("Unable to load records: %s", err)
		}

		flattenedRecords := make([]map[string]interface{}, len(records))
		for i, record := range records {
			flattenedRecords[i] = config.FlattenRecord(record)
		}

		if v, ok := d.GetOk("filter"); ok {
			filters := expandFilters(v.(*schema.Set).List())
			flattenedRecords = applyFilters(config.RecordSchema, flattenedRecords, filters)
		}

		if v, ok := d.GetOk("sort"); ok {
			sorts := expandSorts(v.([]interface{}))
			flattenedRecords = applySorts(config.RecordSchema, flattenedRecords, sorts)
		}

		d.SetId(resource.UniqueId())

		if err := d.Set(config.ResultAttributeName, flattenedRecords); err != nil {
			return fmt.Errorf("unable to set `%s` attribute: %s", config.ResultAttributeName, err)
		}

		return nil
	}
}
