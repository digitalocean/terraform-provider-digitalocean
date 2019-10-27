package digitalocean

import "testing"

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
