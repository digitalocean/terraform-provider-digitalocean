package datalist

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestExpandFilters(t *testing.T) {
	recordSchema := map[string]*schema.Schema{
		"fieldA": {
			Type: schema.TypeString,
		},
		"fieldB": {
			Type: schema.TypeString,
		},
	}

	rawFilters := []interface{}{
		map[string]interface{}{
			"key":    "fieldA",
			"values": []interface{}{"foo", "bar"},
		},
		map[string]interface{}{
			"key":    "fieldB",
			"values": []interface{}{"20", "40"},
		},
	}

	expandedFilters, err := expandFilters(recordSchema, rawFilters)
	if err != nil {
		t.Fatalf("expandFilters returned error: %s", err)
	}

	if len(rawFilters) != len(expandedFilters) {
		t.Fatalf("incorrect expected length of expanded filters")
	}
	if expandedFilters[0].key != "fieldA" ||
		len(expandedFilters[0].values) != 2 ||
		expandedFilters[0].values[0] != "foo" ||
		expandedFilters[0].values[1] != "bar" {
		t.Fatalf("incorrect expansion of the 1st expanded filters")
	}
	if expandedFilters[1].key != "fieldB" ||
		len(expandedFilters[1].values) != 2 ||
		expandedFilters[1].values[0] != "20" ||
		expandedFilters[1].values[1] != "40" {
		t.Fatalf("incorrect expansion of the 2nd expanded filters")
	}
}

func sizesTestSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"slug": {
			Type: schema.TypeString,
		},
		"available": {
			Type: schema.TypeBool,
		},
		"transfer": {
			Type: schema.TypeFloat,
		},
		"price_monthly": {
			Type: schema.TypeFloat,
		},
		"price_hourly": {
			Type: schema.TypeFloat,
		},
		"memory": {
			Type: schema.TypeInt,
		},
		"vcpus": {
			Type: schema.TypeInt,
		},
		"disk": {
			Type: schema.TypeInt,
		},
		"regions": {
			Type: schema.TypeList,
			Elem: &schema.Schema{Type: schema.TypeString},
		},
		"regions_set": {
			Type: schema.TypeSet,
			Elem: &schema.Schema{Type: schema.TypeString},
		},
	}
}
func sizesTestData() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"slug":          "s-1vcpu-1gb",
			"memory":        1024,
			"vcpus":         1,
			"disk":          25,
			"transfer":      1.0,
			"price_monthly": 5.0,
			"price_hourly":  0.007439999841153622,
			"regions":       []interface{}{"sgp1", "sgp2"},
			"regions_set":   schema.NewSet(schema.HashString, []interface{}{"sgp1", "sgp2"}),
			"available":     true,
		},
		{
			"slug":          "s-2vcpu-2gb",
			"memory":        2048,
			"vcpus":         2,
			"disk":          60,
			"transfer":      3.0,
			"price_monthly": 15.0,
			"price_hourly":  0.02232000045478344,
			"regions":       []interface{}{"nyc1", "nyc2"},
			"regions_set":   schema.NewSet(schema.HashString, []interface{}{"nyc1", "nyc2"}),
			"available":     false,
		},
		{
			"slug":          "s-4vcpu-8gb",
			"memory":        8192,
			"vcpus":         4,
			"disk":          160,
			"transfer":      5.0,
			"price_monthly": 40.0,
			"price_hourly":  0.05951999872922897,
			"regions":       []interface{}{"ams1", "ams2"},
			"regions_set":   schema.NewSet(schema.HashString, []interface{}{"ams1", "ams2"}),
			"available":     true,
		},
		{
			"slug":          "m-1vcpu-8gb",
			"memory":        8192,
			"vcpus":         1,
			"disk":          40,
			"transfer":      3.0,
			"price_monthly": 50.0,
			"price_hourly":  0.05952,
			"regions":       []interface{}{"nyc1", "ams1"},
			"regions_set":   schema.NewSet(schema.HashString, []interface{}{"nyc1", "ams1"}),
			"available":     false,
		},
	}
}

func TestApplyFilters(t *testing.T) {
	testCases := []struct {
		name         string
		filter       commonFilter
		expectations []string // Expectations are filled with the expected size slugs in order
	}{
		{
			"BySlug",
			commonFilter{
				"slug",
				[]interface{}{"s-1vcpu-1gb", "s-4vcpu-8gb"},
				false,
				"exact",
			},
			[]string{"s-1vcpu-1gb", "s-4vcpu-8gb"},
		},
		{
			"ByMemory",
			commonFilter{
				"memory",
				[]interface{}{1024, 8192},
				false,
				"exact",
			},
			[]string{"s-1vcpu-1gb", "s-4vcpu-8gb", "m-1vcpu-8gb"},
		},
		{
			"ByCPU",
			commonFilter{
				"vcpus",
				[]interface{}{1, 4},
				false,
				"exact",
			},
			[]string{"s-1vcpu-1gb", "s-4vcpu-8gb", "m-1vcpu-8gb"},
		},
		{
			"ByDisk",
			commonFilter{
				"disk",
				[]interface{}{25, 160},
				false,
				"exact",
			},
			[]string{"s-1vcpu-1gb", "s-4vcpu-8gb"},
		},
		{
			"ByTransfer",
			commonFilter{
				"transfer",
				[]interface{}{1.0, 5.0},
				false,
				"exact",
			},
			[]string{"s-1vcpu-1gb", "s-4vcpu-8gb"},
		},
		{
			"ByPriceMonthly",
			commonFilter{
				"price_monthly",
				[]interface{}{5.0, 40.0},
				false,
				"exact",
			},
			[]string{"s-1vcpu-1gb", "s-4vcpu-8gb"},
		},
		{
			"ByPriceHourly",
			commonFilter{
				"price_hourly",
				[]interface{}{0.00744, 0.05952},
				false,
				"exact",
			},
			[]string{"s-1vcpu-1gb", "s-4vcpu-8gb", "m-1vcpu-8gb"},
		},
		{
			"ByRegions",
			commonFilter{
				"regions",
				[]interface{}{"sgp1", "ams2"},
				false,
				"exact",
			},
			[]string{"s-1vcpu-1gb", "s-4vcpu-8gb"},
		},
		{
			"ByRegionsSet",
			commonFilter{
				"regions_set",
				[]interface{}{"sgp1", "ams2"},
				false,
				"exact",
			},
			[]string{"s-1vcpu-1gb", "s-4vcpu-8gb"},
		},
		{
			"ByAvailable",
			commonFilter{
				"available",
				[]interface{}{true},
				false,
				"exact",
			},
			[]string{"s-1vcpu-1gb", "s-4vcpu-8gb"},
		},
		{
			"ByRegionsSetWithAllValues",
			commonFilter{
				"regions_set",
				[]interface{}{"nyc1", "ams1"},
				true,
				"exact",
			},
			[]string{"m-1vcpu-8gb"},
		},
		{
			"AllBySlug",
			commonFilter{
				"slug",
				[]interface{}{"s-1vcpu-1gb", "s-4vcpu-8gb"},
				true,
				"exact",
			},
			nil,
		},
		{
			"BySlugWithRegularExpression",
			commonFilter{
				"slug",
				[]interface{}{regexp.MustCompile("8gb$")},
				false,
				"re",
			},
			[]string{"s-4vcpu-8gb", "m-1vcpu-8gb"},
		},
		{
			"ByRegionSetWithSubstring",
			commonFilter{
				"regions_set",
				[]interface{}{"nyc"},
				false,
				"substring",
			},
			[]string{"s-2vcpu-2gb", "m-1vcpu-8gb"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			sizes := applyFilters(sizesTestSchema(), sizesTestData(), []commonFilter{testCase.filter})
			var slugs []string
			for _, size := range sizes {
				slugs = append(slugs, size["slug"].(string))
			}
			assert.Equal(t, testCase.expectations, slugs)
		})
	}
}
