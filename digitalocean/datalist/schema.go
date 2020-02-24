package datalist

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

type DataListResourceConfig struct {
	RecordSchema        map[string]*schema.Schema
	FilterKeys          []string
	SortKeys            []string
	ResultAttributeName string
	FlattenRecord       func(record interface{}) map[string]interface{}
	GetRecords          func(meta interface{}) ([]interface{}, error)
}

func NewDataListResource(config *DataListResourceConfig) *schema.Resource {
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

func dataListResourceRead(config *DataListResourceConfig) schema.ReadFunc {
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
