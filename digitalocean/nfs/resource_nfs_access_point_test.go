package nfs

import (
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestResourceDigitalOceanNfsAccessPointInternalValidate(t *testing.T) {
	if err := ResourceDigitalOceanNfsAccessPoint().InternalValidate(nil, true); err != nil {
		t.Fatalf("resource schema invalid: %s", err)
	}
}

func TestDataSourceDigitalOceanNfsAccessPointInternalValidate(t *testing.T) {
	if err := DataSourceDigitalOceanNfsAccessPoint().InternalValidate(nil, false); err != nil {
		t.Fatalf("data source schema invalid: %s", err)
	}
}

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
	if policy.SquashConfig != godo.NfsSquashConfigRootSquash {
		t.Fatalf("unexpected squash_config: %s", policy.SquashConfig)
	}
	if len(policy.Protocols) != 2 {
		t.Fatalf("unexpected protocol count: %d", len(policy.Protocols))
	}
	hasNFS := false
	hasNFS4 := false
	for _, p := range policy.Protocols {
		if p == godo.NfsAccessPolicyProtocolNFS {
			hasNFS = true
		}
		if p == godo.NfsAccessPolicyProtocolNFS4 {
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
	policy := godo.NfsAccessPolicy{
		Anonuid:                    65534,
		Anongid:                    65534,
		Protocols:                  []godo.NfsAccessPolicyProtocol{godo.NfsAccessPolicyProtocolNFS4},
		SquashConfig:               godo.NfsSquashConfigRootSquash,
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

func TestExpandNfsAccessPointPolicy_defaultsWithoutProtocols(t *testing.T) {
	input := []interface{}{map[string]interface{}{
		"anonuid":                      65534,
		"anongid":                      65534,
		"protocols":                    schema.NewSet(schema.HashString, nil),
		"squash_config":                "ROOT_SQUASH",
		"identity_enforcement_enabled": false,
	}}

	policy := expandNfsAccessPointPolicy(input)
	if len(policy.Protocols) != 1 || policy.Protocols[0] != godo.NfsAccessPolicyProtocolNFS4 {
		t.Fatalf("expected default protocol NFS4, got: %#v", policy.Protocols)
	}
}

func TestFlattenNfsAccessPointPolicy_emptyProtocols(t *testing.T) {
	policy := godo.NfsAccessPolicy{
		Anonuid:                    65534,
		Anongid:                    65534,
		Protocols:                  nil,
		SquashConfig:               godo.NfsSquashConfigRootSquash,
		IdentityEnforcementEnabled: false,
	}

	flattened := flattenNfsAccessPointPolicy(policy)
	protocols := flattened[0]["protocols"].(*schema.Set)
	if protocols.Len() != 1 || !protocols.Contains("NFS4") {
		t.Fatalf("expected default protocol NFS4, got: %#v", protocols.List())
	}
}
