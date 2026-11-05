package vpcnatgateway

import (
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestFlattenEgresses(t *testing.T) {
	t.Parallel()

	t.Run("nil", func(t *testing.T) {
		got := flattenEgresses(nil)
		if len(got) != 0 {
			t.Fatalf("expected empty result, got %#v", got)
		}
	})

	t.Run("empty public gateways", func(t *testing.T) {
		got := flattenEgresses(&godo.Egresses{PublicGateways: nil})
		if len(got) != 0 {
			t.Fatalf("expected empty result, got %#v", got)
		}
	})

	t.Run("single gateway", func(t *testing.T) {
		got := flattenEgresses(&godo.Egresses{
			PublicGateways: []*godo.PublicGateway{
				{IPv4: "203.0.113.10"},
			},
		})
		assertSingleEgressWithIPs(t, got, []string{"203.0.113.10"})
	})

	t.Run("multiple gateways one egress entry", func(t *testing.T) {
		got := flattenEgresses(&godo.Egresses{
			PublicGateways: []*godo.PublicGateway{
				{IPv4: "203.0.113.10"},
				{IP: "203.0.113.11"},
			},
		})
		assertSingleEgressWithIPs(t, got, []string{"203.0.113.10", "203.0.113.11"})
	})

	t.Run("prefers IPv4 over IP", func(t *testing.T) {
		got := flattenEgresses(&godo.Egresses{
			PublicGateways: []*godo.PublicGateway{
				{IPv4: "203.0.113.10", IP: "203.0.113.99"},
			},
		})
		assertSingleEgressWithIPs(t, got, []string{"203.0.113.10"})
	})
}

func assertSingleEgressWithIPs(t *testing.T, got []map[string]interface{}, wantIPs []string) {
	t.Helper()
	if len(got) != 1 {
		t.Fatalf("expected exactly 1 egresses entry, got %d (%#v)", len(got), got)
	}
	set, ok := got[0]["public_gateways"].(*schema.Set)
	if !ok {
		t.Fatalf("expected public_gateways to be *schema.Set, got %T", got[0]["public_gateways"])
	}
	if set.Len() != len(wantIPs) {
		t.Fatalf("expected %d public_gateways, got %d", len(wantIPs), set.Len())
	}
	gotIPs := make(map[string]struct{}, set.Len())
	for _, raw := range set.List() {
		m := raw.(map[string]interface{})
		gotIPs[m["ipv4"].(string)] = struct{}{}
	}
	for _, ip := range wantIPs {
		if _, ok := gotIPs[ip]; !ok {
			t.Fatalf("missing ipv4 %q in %#v", ip, gotIPs)
		}
	}
}

func TestExpandEgresses(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		if got := expandEgresses(nil); got != nil {
			t.Fatalf("expected nil, got %#v", got)
		}
		if got := expandEgresses([]interface{}{}); got != nil {
			t.Fatalf("expected nil, got %#v", got)
		}
	})

	t.Run("skips empty ipv4", func(t *testing.T) {
		gatewaySet := schema.NewSet(schema.HashResource(egressPublicGatewaysSchemaResource()), []interface{}{
			map[string]interface{}{"ipv4": ""},
		})
		got := expandEgresses([]interface{}{
			map[string]interface{}{"public_gateways": gatewaySet},
		})
		if got != nil {
			t.Fatalf("expected nil when ipv4 empty, got %#v", got)
		}
	})

	t.Run("byoip", func(t *testing.T) {
		gatewaySet := schema.NewSet(schema.HashResource(egressPublicGatewaysSchemaResource()), []interface{}{
			map[string]interface{}{"ipv4": "203.0.113.10"},
		})
		got := expandEgresses([]interface{}{
			map[string]interface{}{"public_gateways": gatewaySet},
		})
		if got == nil || len(got.PublicGateways) != 1 {
			t.Fatalf("expected 1 public gateway, got %#v", got)
		}
		if got.PublicGateways[0].IP != "203.0.113.10" {
			t.Fatalf("expected IP 203.0.113.10, got %#v", got.PublicGateways[0])
		}
	})
}
