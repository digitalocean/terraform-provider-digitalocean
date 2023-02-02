package util

import "testing"

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
