package uptime

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const testCheckID = "9c6e0217-7ed3-4d3b-babe-67fca97bef23"

func newUptimeTestServer(t *testing.T, status int, body string) *config.CombinedConfig {
	t.Helper()

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	mux.HandleFunc("/v2/uptime/checks/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		fmt.Fprint(w, body)
	})

	cfg := &config.Config{
		Token:       "test-token",
		APIEndpoint: server.URL,
	}
	client, err := cfg.Client()
	if err != nil {
		t.Fatalf("error building client: %s", err)
	}
	return client
}

func TestUptimeReadClearsStateWhenGone(t *testing.T) {
	cases := map[string]struct {
		status int
		body   string
	}{
		"404": {
			status: http.StatusNotFound,
			body:   `{"id":"not_found","message":"The resource you were accessing could not be found."}`,
		},
		// The Uptime API answers a deleted check with this 403.
		"403 for a deleted check": {
			status: http.StatusForbidden,
			body:   `{"id":"forbidden","message":"you are not authorized to access this resource"}`,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			client := newUptimeTestServer(t, tc.status, tc.body)

			check := schema.TestResourceDataRaw(t, ResourceDigitalOceanUptimeCheck().Schema, map[string]interface{}{
				"name":   "check",
				"target": "https://www.digitalocean.com",
			})
			check.SetId(testCheckID)
			if diags := resourceDigitalOceanUptimeCheckRead(context.Background(), check, client); diags.HasError() {
				t.Fatalf("check read: unexpected error: %v", diags)
			}
			if check.Id() != "" {
				t.Fatalf("check read: expected the ID to be cleared, got %q", check.Id())
			}

			alert := schema.TestResourceDataRaw(t, ResourceDigitalOceanUptimeAlert().Schema, map[string]interface{}{
				"check_id": testCheckID,
				"name":     "alert",
				"type":     "latency",
				"period":   "2m",
			})
			alert.SetId("5a4981aa-9653-4bd1-bef5-d6bff52042e4")
			if diags := resourceDigitalOceanUptimeAlertRead(context.Background(), alert, client); diags.HasError() {
				t.Fatalf("alert read: unexpected error: %v", diags)
			}
			if alert.Id() != "" {
				t.Fatalf("alert read: expected the ID to be cleared, got %q", alert.Id())
			}
		})
	}
}

func TestUptimeReadKeepsOtherForbiddenErrors(t *testing.T) {
	client := newUptimeTestServer(t, http.StatusForbidden,
		`{"id":"forbidden","message":"You are missing the required permission uptime:read."}`)

	check := schema.TestResourceDataRaw(t, ResourceDigitalOceanUptimeCheck().Schema, map[string]interface{}{
		"name":   "check",
		"target": "https://www.digitalocean.com",
	})
	check.SetId(testCheckID)
	if diags := resourceDigitalOceanUptimeCheckRead(context.Background(), check, client); !diags.HasError() {
		t.Fatal("expected an error for a 403 that is not a missing check")
	}
	if check.Id() != testCheckID {
		t.Fatalf("expected the ID to be kept, got %q", check.Id())
	}
}
