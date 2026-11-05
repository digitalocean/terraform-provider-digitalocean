package microvm

import (
	"reflect"
	"strings"
	"testing"

	"github.com/digitalocean/godo"
)

func TestExpandEgressPolicy(t *testing.T) {
	block := func(action string, rules ...map[string]interface{}) []interface{} {
		rs := make([]interface{}, len(rules))
		for i, r := range rules {
			rs[i] = r
		}
		return []interface{}{map[string]interface{}{"default_action": action, "rule": rs}}
	}

	tests := []struct {
		name    string
		raw     interface{}
		want    *godo.MicroVMEgressPolicy
		wantErr string
	}{
		{name: "absent block leaves egress unrestricted", raw: []interface{}{}, want: nil},
		{
			name: "domain and ip rules",
			raw: block("deny",
				map[string]interface{}{"domain": "*.example.com", "ip": ""},
				map[string]interface{}{"domain": "", "ip": "203.0.113.5"},
			),
			want: &godo.MicroVMEgressPolicy{
				DefaultAction: godo.MicroVMEgressActionDeny,
				Rules: []godo.MicroVMEgressRule{
					{Domain: "*.example.com"},
					{IP: "203.0.113.5"},
				},
			},
		},
		{
			name:    "rule with both domain and ip",
			raw:     block("deny", map[string]interface{}{"domain": "a.example.com", "ip": "203.0.113.5"}),
			wantErr: "exactly one of domain or ip",
		},
		{
			name:    "rule with neither",
			raw:     block("deny", map[string]interface{}{"domain": "", "ip": ""}),
			wantErr: "exactly one of domain or ip",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := expandEgressPolicy(tt.raw)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestFlattenEgressPolicy_roundTrip(t *testing.T) {
	if got := flattenEgressPolicy(nil); got != nil {
		t.Fatalf("expected nil for a nil policy, got %v", got)
	}

	policy := &godo.MicroVMEgressPolicy{
		DefaultAction: godo.MicroVMEgressActionAllow,
		Rules:         []godo.MicroVMEgressRule{{Domain: "api.example.com"}, {IP: "203.0.113.5"}},
	}
	got, err := expandEgressPolicy(flattenEgressPolicy(policy))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(got, policy) {
		t.Errorf("round trip: got %+v, want %+v", got, policy)
	}
}

func TestValidateEgressDomain(t *testing.T) {
	for _, ok := range []string{"api.example.com", "*.example.com"} {
		if _, errs := validateEgressDomain(ok, "domain"); len(errs) != 0 {
			t.Errorf("%q should be valid, got %v", ok, errs)
		}
	}
	for _, bad := range []string{"", "10.0.0.0/8", "example.com:443", "https://example.com", "example.com/path"} {
		if _, errs := validateEgressDomain(bad, "domain"); len(errs) == 0 {
			t.Errorf("%q should be rejected", bad)
		}
	}
}
