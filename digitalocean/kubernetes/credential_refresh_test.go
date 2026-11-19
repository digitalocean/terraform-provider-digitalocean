package kubernetes

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type credentialTransport func(*http.Request) (*http.Response, error)

func (f credentialTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestClusterCredentialRefresh(t *testing.T) {
	for _, resource := range []struct {
		name string
		make func() *schema.Resource
	}{{"resource", ResourceDigitalOceanKubernetesCluster}, {"data_source", DataSourceDigitalOceanKubernetesCluster}} {
		for _, tc := range []struct {
			name      string
			skip      bool
			expiry    string
			wantCalls int
		}{
			{"default_missing", false, "", 1},
			{"default_expired", false, "2000-01-01T00:00:00Z", 1},
			{"default_valid", false, time.Now().Add(time.Hour).UTC().Format(time.RFC3339), 0},
			{"skip_missing", true, "", 0},
			{"skip_expired", true, "2000-01-01T00:00:00Z", 0},
			{"skip_invalid", true, "invalid", 0},
		} {
			t.Run(resource.name+"/"+tc.name, func(t *testing.T) {
				d := schema.TestResourceDataRaw(t, resource.make().Schema, map[string]interface{}{
					"name": "test", "region": "nyc2", "version": "1.34.1-do.0", "skip_credentials": tc.skip,
				})
				if tc.expiry != "" {
					if err := d.Set("kube_config", []interface{}{map[string]interface{}{"expires_at": tc.expiry, "token": "old-test-token"}}); err != nil {
						t.Fatal(err)
					}
				}
				calls := 0
				client := godo.NewClient(&http.Client{Transport: credentialTransport(func(r *http.Request) (*http.Response, error) {
					calls++
					if r.Method != "GET" || r.URL.Path != "/v2/kubernetes/clusters/test-id/credentials" {
						t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
					}
					return &http.Response{StatusCode: 403, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(`{"message":"credential access denied"}`)), Request: r}, nil
				})})
				cluster := &godo.KubernetesCluster{ID: "test-id", Name: "live-name", RegionSlug: "nyc2", VersionSlug: "1.34.1-do.0", Endpoint: "https://live.example", Status: &godo.KubernetesClusterStatus{State: "running"}, MaintenancePolicy: &godo.KubernetesMaintenancePolicy{StartTime: "04:00", Day: godo.KubernetesMaintenancePolicyDay(1)}}
				diags := digitaloceanKubernetesClusterRead(client, cluster, d)
				if calls != tc.wantCalls {
					t.Fatalf("credential calls = %d, want %d", calls, tc.wantCalls)
				}
				if diags.HasError() != (tc.wantCalls > 0) {
					t.Fatalf("unexpected diagnostics: %v", diags)
				}
				if d.Get("name") != "live-name" || d.Get("endpoint") != "https://live.example" || d.Get("version") != "1.34.1-do.0" {
					t.Fatal("cluster metadata was not refreshed")
				}
				if tc.skip && len(d.Get("kube_config").([]interface{})) != 0 {
					t.Fatal("skipped credentials remain in refreshed state")
				}
			})
		}
	}
}
