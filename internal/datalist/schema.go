package datalist

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// This is the configuration for a "data list" resource. It represents the schema and operations
// needed to create the data list resource.
type ResourceConfig struct {
	// The schema for a single instance of the resource.
	RecordSchema map[string]*schema.Schema

	// The name of the attribute in the resource through which to expose results.
	ResultAttributeName string

	// Given a record returned from the GetRecords function, flatten the record to a
	// map acceptable to the Set method on schema.ResourceData.
	FlattenRecord func(record, meta interface{}, extra map[string]interface{}) (map[string]interface{}, error)

	// Return all of the records on which the data list resource should operate.
	// The `meta` argument is the same meta argument passed into the resource's Read
	// function.
	GetRecords func(meta interface{}, extra map[string]interface{}) ([]interface{}, error)

	// Extra parameters to expose on the datasource alongside `filter` and `sort`.
	ExtraQuerySchema map[string]*schema.Schema
}

// Returns a new "data list" resource given the specified configuration. This
// is a resource with `filter` and `sort` attributes that can select a subset
// of records from a list of records for a particular type of resource.
func NewResource(config *ResourceConfig) *schema.Resource {
	err := validateResourceConfig(config)
	if err != nil {
		// Panic if the resource config is invalid since this will prevent the resource
		// from operating.
		log.Panicf("datalist.NewResource: invalid resource configuration: %v", err)
	}

	recordSchema := map[string]*schema.Schema{}
	for attributeName, attributeSchema := range config.RecordSchema {
		newAttributeSchema := &schema.Schema{}
		*newAttributeSchema = *attributeSchema
		newAttributeSchema.Computed = true
		newAttributeSchema.Required = false
		newAttributeSchema.Optional = false
		recordSchema[attributeName] = newAttributeSchema
	}

	filterKeys := computeFilterKeys(recordSchema)
	sortKeys := computeSortKeys(recordSchema)

	datasourceSchema := map[string]*schema.Schema{
		"filter": filterSchema(filterKeys),
		"sort":   sortSchema(sortKeys),
		config.ResultAttributeName: {
			Type:     schema.TypeList,
			Computed: true,
			Elem: &schema.Resource{
				Schema: recordSchema,
			},
		},
	}

	for key, value := range config.ExtraQuerySchema {
		datasourceSchema[key] = value
	}

	return &schema.Resource{
		ReadContext: dataListResourceRead(config),
		Schema:      datasourceSchema,
	}
}

func dataListResourceRead(config *ResourceConfig) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		extra := map[string]interface{}{}
		for key, _ := range config.ExtraQuerySchema {
			extra[key] = d.Get(key)
		}

		records, err := config.GetRecords(meta, extra)
		if err != nil {
			return diag.Errorf("Unable to load records: %s", err)
		}

		flattenedRecords := make([]map[string]interface{}, len(records))
		for i, record := range records {
			flattenedRecord, err := config.FlattenRecord(record, meta, extra)
			if err != nil {
				return diag.FromErr(err)
			}
			flattenedRecords[i] = flattenedRecord
		}

		if v, ok := d.GetOk("filter"); ok {
			filters, err := expandFilters(config.RecordSchema, v.(*schema.Set).List())
			if err != nil {
				return diag.FromErr(err)
			}
			flattenedRecords = applyFilters(config.RecordSchema, flattenedRecords, filters)
		}

		if v, ok := d.GetOk("sort"); ok {
			sorts := expandSorts(v.([]interface{}))
			flattenedRecords = applySorts(config.RecordSchema, flattenedRecords, sorts)
		}

		d.SetId(resource.UniqueId())

		if err := d.Set(config.ResultAttributeName, flattenedRecords); err != nil {
			return diag.Errorf("unable to set `%s` attribute: %s", config.ResultAttributeName, err)
		}

		return nil
	}
}

// Compute the set of filter keys for the resource.
func computeFilterKeys(recordSchema map[string]*schema.Schema) []string {
	var filterKeys []string

	for key, schemaForKey := range recordSchema {
		if schemaForKey.Type != schema.TypeMap {
			filterKeys = append(filterKeys, key)
		}
	}

	return filterKeys
}

// Compute the set of sort keys for the source.
func computeSortKeys(recordSchema map[string]*schema.Schema) []string {
	var sortKeys []string

	for key, schemaForKey := range recordSchema {
		supported := false
		switch schemaForKey.Type {
		case schema.TypeString, schema.TypeBool, schema.TypeInt, schema.TypeFloat:
			supported = true
		}

		if supported {
			sortKeys = append(sortKeys, key)
		}
	}

	return sortKeys

}

// Validate a ResourceConfig to ensure it conforms to this package's assumptions.
func validateResourceConfig(config *ResourceConfig) error {
	// Ensure that ResultAttributeName exists.
	if config.ResultAttributeName == "" {
		return fmt.Errorf("ResultAttributeName must be specified")
	}

	return nil
}
