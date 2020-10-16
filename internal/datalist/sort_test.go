package datalist

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func sizesTestDataForSorts() []map[string]interface{} {
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
	}
}

func TestExpandSorts(t *testing.T) {
	rawSorts := []interface{}{
		map[string]interface{}{
			"key":       "fieldA",
			"direction": "asc",
		},
		map[string]interface{}{
			"key":       "fieldB",
			"direction": "desc",
		},
	}

	expandedSorts := expandSorts(rawSorts)

	if len(rawSorts) != len(expandedSorts) {
		t.Fatalf("incorrect expected length of expanded sorts")
	}
	if expandedSorts[0].key != "fieldA" ||
		expandedSorts[0].direction != "asc" {
		t.Fatalf("incorrect expansion of the 1st expanded sorts")
	}
	if expandedSorts[1].key != "fieldB" ||
		expandedSorts[1].direction != "desc" {
		t.Fatalf("incorrect expansion of the 2nd expanded sorts")
	}
}

func TestApplySorts(t *testing.T) {
	testCases := []struct {
		name        string
		key         string
		expectedAsc []string // Expected sizes if sorted ascendingly
	}{
		{"BySlug", "slug", []string{"s-1vcpu-1gb", "s-2vcpu-2gb", "s-4vcpu-8gb"}},
		{"ByMemory", "memory", []string{"s-1vcpu-1gb", "s-2vcpu-2gb", "s-4vcpu-8gb"}},
		{"ByCPU", "vcpus", []string{"s-1vcpu-1gb", "s-2vcpu-2gb", "s-4vcpu-8gb"}},
		{"ByDisk", "disk", []string{"s-1vcpu-1gb", "s-2vcpu-2gb", "s-4vcpu-8gb"}},
		{"ByTransfer", "transfer", []string{"s-1vcpu-1gb", "s-2vcpu-2gb", "s-4vcpu-8gb"}},
		{"ByPriceMonthly", "price_monthly", []string{"s-1vcpu-1gb", "s-2vcpu-2gb", "s-4vcpu-8gb"}},
		{"ByPriceHourly", "price_hourly", []string{"s-1vcpu-1gb", "s-2vcpu-2gb", "s-4vcpu-8gb"}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Test ascending order
			sizes := applySorts(sizesTestSchema(), sizesTestDataForSorts(), []commonSort{{testCase.key, "asc"}})
			if len(sizes) != len(testCase.expectedAsc) {
				t.Fatalf("Expecting %d size results, found %d size results instead", len(testCase.expectedAsc), len(sizes))
			}
			for i, expectedSlug := range testCase.expectedAsc {
				if sizes[i]["slug"] != expectedSlug {
					t.Fatalf("Expecting size index %d to be %s, found %s instead", i, expectedSlug, sizes[i]["slug"])
				}
			}

			// Test descending order
			sizes = applySorts(sizesTestSchema(), sizesTestDataForSorts(), []commonSort{{testCase.key, "desc"}})
			if len(sizes) != len(testCase.expectedAsc) {
				t.Fatalf("Expecting %d size results, found %d size results instead", len(testCase.expectedAsc), len(sizes))
			}
			for i, expectedSlug := range testCase.expectedAsc {
				if sizes[len(sizes)-i-1]["slug"] != expectedSlug {
					t.Fatalf("Expecting size index %d to be %s, found %s instead", i, expectedSlug, sizes[i]["slug"])
				}
			}
		})
	}
}

func TestApplySortsMultiple(t *testing.T) {
	testData := []map[string]interface{}{
		{
			"slug":          "s-1vcpu-1gb",
			"memory":        1024,
			"vcpus":         1,
			"disk":          25,
			"transfer":      1.0,
			"price_monthly": 5.0,
			"price_hourly":  0.007439999841153622,
			"regions":       []string{"sgp1", "sgp2"},
			"available":     true,
		},
		{
			"slug":          "1gb",
			"memory":        1024,
			"vcpus":         1,
			"disk":          30,
			"transfer":      2.0,
			"price_monthly": 10.0,
			"price_hourly":  0.01487999968230724,
			"regions":       []string{"sgp1", "sgp2"},
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
			"regions":       []string{"nyc1", "nyc2"},
			"available":     false,
		},
	}

	// Test ascending order
	sizes := applySorts(sizesTestSchema(), testData, []commonSort{
		{"memory", "desc"}, // Sort by memory descendingly first
		{"disk", "asc"},    // Then for sizes with same memory, sort by disk ascendingly
	})

	if len(sizes) != 3 {
		t.Fatalf("Expecting 3 size results, found %d size results instead", len(sizes))
	}

	// s-2vcpu-2gb 	(Memory = 2048)
	// s-1vcpu-1gb 	(Memory = 1024, Disk = 25)
	// 1gb			(Memory = 1024, Disk = 30)
	if sizes[0]["slug"] != "s-2vcpu-2gb" ||
		sizes[1]["slug"] != "s-1vcpu-1gb" ||
		sizes[2]["slug"] != "1gb" {
		t.Fatalf("Expecting sizes to be sorted by memory in descending order, then by disk in ascending order")
	}

}
