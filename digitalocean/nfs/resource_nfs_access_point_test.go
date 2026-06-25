package nfs

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestExpandNfsAccessPointPolicy(t *testing.T) {
	set := schema.NewSet(schema.HashString, []interface{}{"nfs", "NFS4"})
	input := []interface{}{map[string]interface{}{
		"anonuid":                      123,
		"anongid":                      456,
		"protocols":                    set,
		"squash_config":                "root_squash",
		"identity_enforcement_enabled": true,
	}}

	policy := expandNfsAccessPointPolicy(input)
	if policy.Anonuid != uint64(123) {
		t.Fatalf("unexpected anonuid: %d", policy.Anonuid)
	}
	if policy.Anongid != uint64(456) {
		t.Fatalf("unexpected anongid: %d", policy.Anongid)
	}
	if policy.SquashConfig != "ROOT_SQUASH" {
		t.Fatalf("unexpected squash_config: %s", policy.SquashConfig)
	}
	if len(policy.Protocols) != 2 {
		t.Fatalf("unexpected protocol count: %d", len(policy.Protocols))
	}
	hasNFS := false
	hasNFS4 := false
	for _, p := range policy.Protocols {
		if p == "NFS" {
			hasNFS = true
		}
		if p == "NFS4" {
			hasNFS4 = true
		}
	}
	if !hasNFS || !hasNFS4 {
		t.Fatalf("unexpected protocols: %#v", policy.Protocols)
	}
	if !policy.IdentityEnforcementEnabled {
		t.Fatalf("expected identity_enforcement_enabled to be true")
	}
}

func TestFlattenNfsAccessPointPolicy(t *testing.T) {
	policy := nfsAccessPointPolicy{
		Anonuid:                    65534,
		Anongid:                    65534,
		Protocols:                  []string{"NFS4"},
		SquashConfig:               "ROOT_SQUASH",
		IdentityEnforcementEnabled: false,
	}

	flattened := flattenNfsAccessPointPolicy(policy)
	if len(flattened) != 1 {
		t.Fatalf("unexpected flattened len: %d", len(flattened))
	}
	if flattened[0]["anonuid"] != 65534 {
		t.Fatalf("unexpected anonuid: %#v", flattened[0]["anonuid"])
	}
	if flattened[0]["anongid"] != 65534 {
		t.Fatalf("unexpected anongid: %#v", flattened[0]["anongid"])
	}
	if flattened[0]["squash_config"] != "ROOT_SQUASH" {
		t.Fatalf("unexpected squash_config: %#v", flattened[0]["squash_config"])
	}
	if flattened[0]["identity_enforcement_enabled"] != false {
		t.Fatalf("unexpected identity_enforcement_enabled: %#v", flattened[0]["identity_enforcement_enabled"])
	}
}
