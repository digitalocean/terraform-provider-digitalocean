package util

import "testing"

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
