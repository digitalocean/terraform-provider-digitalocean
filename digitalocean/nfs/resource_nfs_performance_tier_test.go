package nfs

import "testing"

func TestNormalizeNfsPerformanceTier(t *testing.T) {
	tests := map[string]string{
		"standard":                  "standard",
		"high":                      "high",
		"PERFORMANCE_TIER_STANDARD": "standard",
		"PERFORMANCE_TIER_HIGH":     "high",
	}

	for input, want := range tests {
		if got := normalizeNfsPerformanceTier(input); got != want {
			t.Fatalf("normalizeNfsPerformanceTier(%q) = %q, want %q", input, got, want)
		}
	}
}
