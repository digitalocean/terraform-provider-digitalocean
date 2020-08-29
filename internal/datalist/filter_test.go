package datalist

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func TestExpandFilters(t *testing.T) {
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

	expandedFilters := expandFilters(rawFilters)

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
		{"BySlug", commonFilter{"slug", []string{"s-1vcpu-1gb", "s-4vcpu-8gb"}, false}, []string{"s-1vcpu-1gb", "s-4vcpu-8gb"}},
		{"ByMemory", commonFilter{"memory", []string{"1024", "8192"}, false}, []string{"s-1vcpu-1gb", "s-4vcpu-8gb", "m-1vcpu-8gb"}},
		{"ByCPU", commonFilter{"vcpus", []string{"1", "4"}, false}, []string{"s-1vcpu-1gb", "s-4vcpu-8gb", "m-1vcpu-8gb"}},
		{"ByDisk", commonFilter{"disk", []string{"25", "160"}, false}, []string{"s-1vcpu-1gb", "s-4vcpu-8gb"}},
		{"ByTransfer", commonFilter{"transfer", []string{"1.0", "5.0"}, false}, []string{"s-1vcpu-1gb", "s-4vcpu-8gb"}},
		{"ByPriceMonthly", commonFilter{"price_monthly", []string{"5.0", "40.0"}, false}, []string{"s-1vcpu-1gb", "s-4vcpu-8gb"}},
		{"ByPriceHourly", commonFilter{"price_hourly", []string{"0.00744", "0.05952"}, false}, []string{"s-1vcpu-1gb", "s-4vcpu-8gb", "m-1vcpu-8gb"}},
		{"ByRegions", commonFilter{"regions", []string{"sgp1", "ams2"}, false}, []string{"s-1vcpu-1gb", "s-4vcpu-8gb"}},
		{"ByRegionsSet", commonFilter{"regions_set", []string{"sgp1", "ams2"}, false}, []string{"s-1vcpu-1gb", "s-4vcpu-8gb"}},
		{"ByAvailable", commonFilter{"available", []string{"true"}, false}, []string{"s-1vcpu-1gb", "s-4vcpu-8gb"}},
		{"AllByRegionsSet", commonFilter{"regions_set", []string{"nyc1", "ams1"}, true}, []string{"m-1vcpu-8gb"}},
		{"AllBySlug", commonFilter{"slug", []string{"s-1vcpu-1gb", "s-4vcpu-8gb"}, true}, nil},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			sizes := applyFilters(sizesTestSchema(), sizesTestData(), []commonFilter{testCase.filter})
			if len(sizes) != len(testCase.expectations) {
				t.Fatalf("Expecting %d size results, found %d size results instead", len(testCase.expectations), len(sizes))
			}
			for i, expectedSlug := range testCase.expectations {
				if sizes[i]["slug"] != expectedSlug {
					t.Fatalf("Expecting size index %d to be %s, found %s instead", i, expectedSlug, sizes[i]["slug"])
				}
			}
		})
	}
}
