package util

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestCaseSensitive(t *testing.T) {
	cases := []struct {
		Name     string
		Left     string
		Right    string
		Suppress bool
	}{
		{
			Name:     "empty",
			Left:     "",
			Right:    "",
			Suppress: true,
		},
		{
			Name:     "empty and text",
			Left:     "text",
			Right:    "",
			Suppress: false,
		},
		{
			Name:     "different text",
			Left:     "text",
			Right:    "different text",
			Suppress: false,
		},
		{
			Name:     "same text",
			Left:     "text",
			Right:    "text",
			Suppress: true,
		},
		{
			Name:     "same text different case",
			Left:     "text",
			Right:    "TeXT",
			Suppress: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			if CaseSensitive("test", tc.Left, tc.Right, nil) != tc.Suppress {
				t.Fatalf("Expected CaseSensitive to return %t for '%q' == '%q'", tc.Suppress, tc.Left, tc.Right)
			}
		})
	}
}

func TestFloat32Precision(t *testing.T) {
	cases := []struct {
		Name     string
		Left     string
		Right    string
		Suppress bool
	}{
		{
			Name:     "exact match",
			Left:     "0.001",
			Right:    "0.001",
			Suppress: true,
		},
		{
			Name:     "float32 expansion of 0.001",
			Left:     "0.0010000000474974513",
			Right:    "0.001",
			Suppress: true,
		},
		{
			Name:     "zero is distinct from small values",
			Left:     "0",
			Right:    "0.001",
			Suppress: false,
		},
		{
			Name:     "materially different values",
			Left:     "0.1",
			Right:    "0.01",
			Suppress: false,
		},
		{
			Name:     "zero equals zero",
			Left:     "0",
			Right:    "0",
			Suppress: true,
		},
		{
			Name:     "invalid values are not suppressed",
			Left:     "not-a-number",
			Right:    "0.001",
			Suppress: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			if Float32Precision("test", tc.Left, tc.Right, nil) != tc.Suppress {
				t.Fatalf("Expected Float32Precision to return %t for '%q' == '%q'", tc.Suppress, tc.Left, tc.Right)
			}
		})
	}
}

func TestNormalizeFloat32(t *testing.T) {
	cases := []struct {
		Name string
		In   float64
		Want float64
	}{
		{Name: "zero", In: 0, Want: 0},
		{Name: "0.001", In: 0.001, Want: 0.001},
		{Name: "0.01", In: 0.01, Want: 0.01},
		{Name: "0.1", In: 0.1, Want: 0.1},
		{Name: "exact binary", In: 0.0009765625, Want: 0.0009765625},
		{Name: "ten", In: 10, Want: 10},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			got := NormalizeFloat32(float32(tc.In))
			if got != tc.Want {
				t.Fatalf("NormalizeFloat32(%v) = %v, want %v", tc.In, got, tc.Want)
			}
		})
	}
}

func TestSetFloat32Attribute(t *testing.T) {
	r := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"long_query_time": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
		},
	}

	d := r.TestResourceData()
	api := float32(0.001)
	SetFloat32Attribute(d, "long_query_time", &api)

	got := d.Get("long_query_time").(float64)
	if got != 0.001 {
		t.Fatalf("expected 0.001 after normalize, got %v", got)
	}
}
