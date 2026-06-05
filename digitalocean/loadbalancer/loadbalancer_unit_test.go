package loadbalancer

import (
	"context"
	"reflect"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// newLBResourceData materializes a *schema.ResourceData for the
// digitalocean_loadbalancer resource and sets the provided attributes.
// It is the most ergonomic way to exercise helpers that read from the
// resource without standing up a full provider.
func newLBResourceData(t *testing.T, attrs map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := ResourceDigitalOceanLoadbalancer()
	d := r.TestResourceData()
	for k, v := range attrs {
		if err := d.Set(k, v); err != nil {
			t.Fatalf("d.Set(%q): %v", k, err)
		}
	}
	return d
}

func TestGLBDesiredDomainCertificates(t *testing.T) {
	cases := []struct {
		name    string
		domains []interface{}
		want    map[string]string
	}{
		{
			name:    "no domains returns empty map",
			domains: nil,
			want:    map[string]string{},
		},
		{
			name: "single custom domain with certificate is included",
			domains: []interface{}{
				map[string]interface{}{
					"name":             "example.com",
					"certificate_name": "cert-new",
					"is_managed":       false,
				},
			},
			want: map[string]string{"example.com": "cert-new"},
		},
		{
			name: "managed domain is excluded",
			domains: []interface{}{
				map[string]interface{}{
					"name":             "managed.example.com",
					"certificate_name": "any-cert",
					"is_managed":       true,
				},
			},
			want: map[string]string{},
		},
		{
			name: "domain without certificate_name is excluded",
			domains: []interface{}{
				map[string]interface{}{
					"name":       "no-cert.example.com",
					"is_managed": false,
				},
			},
			want: map[string]string{},
		},
		{
			name: "domain with empty name is excluded",
			domains: []interface{}{
				map[string]interface{}{
					"name":             "",
					"certificate_name": "cert-new",
				},
			},
			want: map[string]string{},
		},
		{
			name: "mixed set keeps only eligible domains",
			domains: []interface{}{
				map[string]interface{}{
					"name":             "a.example.com",
					"certificate_name": "cert-a",
				},
				map[string]interface{}{
					"name":             "b.example.com",
					"certificate_name": "cert-b",
				},
				map[string]interface{}{
					"name":             "managed.example.com",
					"certificate_name": "ignored",
					"is_managed":       true,
				},
				map[string]interface{}{
					"name": "no-cert.example.com",
				},
			},
			want: map[string]string{
				"a.example.com": "cert-a",
				"b.example.com": "cert-b",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			attrs := map[string]interface{}{}
			if tc.domains != nil {
				attrs["domains"] = tc.domains
			}
			d := newLBResourceData(t, attrs)

			got := glbDesiredDomainCertificates(d)
			if !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("glbDesiredDomainCertificates() = %#v, want %#v", got, tc.want)
			}
		})
	}
}

// TestWaitForGLBDomainCertificatesIfNeeded_NoOpPaths covers the cheap paths
// that must short-circuit before any API call is made. We deliberately pass a
// nil godo client; if any of these branches were to dereference the client
// the test would panic, surfacing the regression.
func TestWaitForGLBDomainCertificatesIfNeeded_NoOpPaths(t *testing.T) {
	cases := []struct {
		name  string
		attrs map[string]interface{}
	}{
		{
			name: "regional load balancer skips GLB polling",
			attrs: map[string]interface{}{
				"type": "REGIONAL",
				"domains": []interface{}{
					map[string]interface{}{
						"name":             "example.com",
						"certificate_name": "cert-new",
					},
				},
			},
		},
		{
			name: "global load balancer with no custom domains skips polling",
			attrs: map[string]interface{}{
				"type": "GLOBAL",
			},
		},
		{
			name: "global load balancer with only managed domains skips polling",
			attrs: map[string]interface{}{
				"type": "GLOBAL",
				"domains": []interface{}{
					map[string]interface{}{
						"name":       "managed.example.com",
						"is_managed": true,
					},
				},
			},
		},
		{
			name: "global load balancer with domain but no certificate_name skips polling",
			attrs: map[string]interface{}{
				"type": "GLOBAL",
				"domains": []interface{}{
					map[string]interface{}{
						"name": "no-cert.example.com",
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			d := newLBResourceData(t, tc.attrs)
			// nil client is safe only because each case must short-circuit.
			if err := waitForGLBDomainCertificatesIfNeeded(context.Background(), d, nil, "ignored-lb-id"); err != nil {
				t.Fatalf("expected nil error on no-op path, got %v", err)
			}
		})
	}
}

// TestWaitForGLBDomainCertificateBindings_EmptyDesired makes sure the
// underlying poller exits immediately when there is nothing to wait for,
// without touching the godo client.
func TestWaitForGLBDomainCertificateBindings_EmptyDesired(t *testing.T) {
	if err := waitForGLBDomainCertificateBindings(context.Background(), nil, "ignored-lb-id", nil); err != nil {
		t.Fatalf("expected nil error for empty desired set, got %v", err)
	}
	if err := waitForGLBDomainCertificateBindings(context.Background(), nil, "ignored-lb-id", map[string]string{}); err != nil {
		t.Fatalf("expected nil error for empty desired map, got %v", err)
	}
}
