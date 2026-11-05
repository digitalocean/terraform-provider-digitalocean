package dropletautoscale

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestExpandTemplatePublicNetworking(t *testing.T) {
	baseTemplate := map[string]interface{}{
		"size":               "s-1vcpu-1gb",
		"region":             "nyc3",
		"image":              "ubuntu-24-04-x64",
		"tags":               []interface{}{},
		"ssh_keys":           []interface{}{"fingerprint"},
		"vpc_uuid":           "",
		"with_droplet_agent": true,
		"project_id":         "",
		"ipv6":               false,
		"user_data":          "",
	}

	t.Run("omitted leaves PublicNetworking nil so the API default applies", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, ResourceDigitalOceanDropletAutoscale().Schema, map[string]interface{}{
			"name":             "pool",
			"droplet_template": []interface{}{baseTemplate},
		})

		got := expandTemplate(d)
		if got == nil {
			t.Fatal("expected template, got nil")
		}
		if got.PublicNetworking != nil {
			t.Fatalf("PublicNetworking = %v, want nil", *got.PublicNetworking)
		}
	})

	t.Run("explicit false is propagated", func(t *testing.T) {
		tmpl := copyMap(baseTemplate)
		tmpl["public_networking"] = false
		d := schema.TestResourceDataRaw(t, ResourceDigitalOceanDropletAutoscale().Schema, map[string]interface{}{
			"name":             "pool",
			"droplet_template": []interface{}{tmpl},
		})

		got := expandTemplate(d)
		if got == nil {
			t.Fatal("expected template, got nil")
		}
		if got.PublicNetworking == nil {
			t.Fatal("PublicNetworking is nil, want false")
		}
		if *got.PublicNetworking {
			t.Fatal("PublicNetworking = true, want false")
		}
	})

	t.Run("explicit true is propagated", func(t *testing.T) {
		tmpl := copyMap(baseTemplate)
		tmpl["public_networking"] = true
		d := schema.TestResourceDataRaw(t, ResourceDigitalOceanDropletAutoscale().Schema, map[string]interface{}{
			"name":             "pool",
			"droplet_template": []interface{}{tmpl},
		})

		got := expandTemplate(d)
		if got == nil {
			t.Fatal("expected template, got nil")
		}
		if got.PublicNetworking == nil {
			t.Fatal("PublicNetworking is nil, want true")
		}
		if !*got.PublicNetworking {
			t.Fatal("PublicNetworking = false, want true")
		}
	})
}

func copyMap(in map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
