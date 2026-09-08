package droplet

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
)

// TestResourceDigitalOceanDropletCreate_WaitsThroughMissingIPThenSucceeds covers
// ESC-25905: create must not finish on status=active alone when public networking
// is enabled. It keeps polling until a public IPv4 is readable.
func TestResourceDigitalOceanDropletCreate_WaitsThroughMissingIPThenSucceeds(t *testing.T) {
	const dropletID = 598498975
	const wantIP = "161.35.84.149"

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	var gets atomic.Int32

	mux.HandleFunc("/v2/droplets", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, `{
			"droplet": {
				"id": %d,
				"name": "esc-25905-fixed",
				"status": "new",
				"locked": false,
				"created_at": "2026-09-07T15:00:00Z",
				"memory": 1024,
				"vcpus": 1,
				"disk": 25,
				"region": {"slug": "ams3"},
				"size": {"slug": "s-1vcpu-1gb", "price_hourly": 0.00893, "price_monthly": 6},
				"networks": {"v4": [], "v6": []},
				"features": [],
				"tags": [],
				"volume_ids": []
			}
		}`, dropletID)
	})

	mux.HandleFunc(fmt.Sprintf("/v2/droplets/%d", dropletID), func(w http.ResponseWriter, r *http.Request) {
		n := gets.Add(1)
		w.Header().Set("Content-Type", "application/json")

		networks := map[string]interface{}{
			"v4": []interface{}{},
			"v6": []interface{}{},
		}
		// First active reads still omit the public IP (the race window).
		if n >= 3 {
			networks["v4"] = []map[string]string{{
				"ip_address": wantIP,
				"netmask":    "255.255.240.0",
				"gateway":    "161.35.80.1",
				"type":       "public",
			}}
		}

		payload := map[string]interface{}{
			"droplet": map[string]interface{}{
				"id":         dropletID,
				"name":       "esc-25905-fixed",
				"status":     "active",
				"locked":     false,
				"created_at": "2026-09-07T15:00:00Z",
				"memory":     1024,
				"vcpus":      1,
				"disk":       25,
				"region":     map[string]string{"slug": "ams3"},
				"size": map[string]interface{}{
					"slug":          "s-1vcpu-1gb",
					"price_hourly":  0.00893,
					"price_monthly": 6,
				},
				"networks":   networks,
				"features":   []string{},
				"tags":       []string{},
				"volume_ids": []string{},
			},
		}
		_ = json.NewEncoder(w).Encode(payload)
	})

	cfg := &config.Config{
		Token:       "test-token",
		APIEndpoint: server.URL,
	}
	client, err := cfg.Client()
	if err != nil {
		t.Fatalf("error building client: %s", err)
	}

	d := ResourceDigitalOceanDroplet().TestResourceData()
	_ = d.Set("image", "debian-13-x64")
	_ = d.Set("name", "esc-25905-fixed")
	_ = d.Set("region", "ams3")
	_ = d.Set("size", "s-1vcpu-1gb")

	if diags := resourceDigitalOceanDropletCreate(context.Background(), d, client); diags.HasError() {
		t.Fatalf("create returned error: %v", diags)
	}
	if got := d.Get("ipv4_address").(string); got != wantIP {
		t.Fatalf("ipv4_address = %q, expected %q (gets=%d)", got, wantIP, gets.Load())
	}
	if gets.Load() < 3 {
		t.Fatalf("expected create to poll through missing-IP window, gets=%d", gets.Load())
	}
}

// TestResourceDigitalOceanDropletCreate_ActiveWithoutPublicIPv4TimesOut ensures
// create does not succeed with an empty public IP when public networking is on.
func TestResourceDigitalOceanDropletCreate_ActiveWithoutPublicIPv4TimesOut(t *testing.T) {
	const dropletID = 598498974

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/v2/droplets", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, `{
			"droplet": {
				"id": %d,
				"name": "esc-25905",
				"status": "new",
				"locked": false,
				"created_at": "2026-09-07T15:00:00Z",
				"memory": 1024,
				"vcpus": 1,
				"disk": 25,
				"region": {"slug": "ams3"},
				"size": {"slug": "s-1vcpu-1gb", "price_hourly": 0.00893, "price_monthly": 6},
				"networks": {"v4": [], "v6": []},
				"features": ["private_networking"],
				"tags": [],
				"volume_ids": []
			}
		}`, dropletID)
	})

	mux.HandleFunc(fmt.Sprintf("/v2/droplets/%d", dropletID), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"droplet": {
				"id": %d,
				"name": "esc-25905",
				"status": "active",
				"locked": false,
				"created_at": "2026-09-07T15:00:00Z",
				"memory": 1024,
				"vcpus": 1,
				"disk": 25,
				"region": {"slug": "ams3"},
				"size": {"slug": "s-1vcpu-1gb", "price_hourly": 0.00893, "price_monthly": 6},
				"networks": {"v4": [], "v6": []},
				"features": ["private_networking"],
				"tags": [],
				"volume_ids": [],
				"vpc_uuid": "00000000-0000-4000-8000-000000000001"
			}
		}`, dropletID)
	})

	cfg := &config.Config{
		Token:       "test-token",
		APIEndpoint: server.URL,
	}
	client, err := cfg.Client()
	if err != nil {
		t.Fatalf("error building client: %s", err)
	}

	d := ResourceDigitalOceanDroplet().TestResourceData()
	_ = d.Set("image", "debian-13-x64")
	_ = d.Set("name", "esc-25905")
	_ = d.Set("region", "ams3")
	_ = d.Set("size", "s-1vcpu-1gb")

	// Bound the wait so the test does not use the full create timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	diags := resourceDigitalOceanDropletCreate(ctx, d, client)
	if !diags.HasError() {
		t.Fatalf("expected create to error while public IPv4 stays missing, got success with ipv4=%q", d.Get("ipv4_address"))
	}
}

// TestResourceDigitalOceanDropletCreate_PrivateNetworkingSkipsPublicIPWait ensures
// public_networking=false remains non-breaking: create succeeds without a public IP.
func TestResourceDigitalOceanDropletCreate_PrivateNetworkingSkipsPublicIPWait(t *testing.T) {
	const dropletID = 598498976

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc("/v2/droplets", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, `{
			"droplet": {
				"id": %d,
				"name": "esc-25905-private",
				"status": "new",
				"locked": false,
				"created_at": "2026-09-07T15:00:00Z",
				"memory": 1024,
				"vcpus": 1,
				"disk": 25,
				"region": {"slug": "ams3"},
				"size": {"slug": "s-1vcpu-1gb", "price_hourly": 0.00893, "price_monthly": 6},
				"networks": {"v4": [], "v6": []},
				"features": ["private_networking"],
				"tags": [],
				"volume_ids": []
			}
		}`, dropletID)
	})

	mux.HandleFunc(fmt.Sprintf("/v2/droplets/%d", dropletID), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{
			"droplet": {
				"id": %d,
				"name": "esc-25905-private",
				"status": "active",
				"locked": false,
				"created_at": "2026-09-07T15:00:00Z",
				"memory": 1024,
				"vcpus": 1,
				"disk": 25,
				"region": {"slug": "ams3"},
				"size": {"slug": "s-1vcpu-1gb", "price_hourly": 0.00893, "price_monthly": 6},
				"networks": {
					"v4": [{"ip_address": "10.108.0.20", "netmask": "255.255.240.0", "gateway": "10.108.0.1", "type": "private"}],
					"v6": []
				},
				"features": ["private_networking"],
				"tags": [],
				"volume_ids": [],
				"vpc_uuid": "00000000-0000-4000-8000-000000000001"
			}
		}`, dropletID)
	})

	cfg := &config.Config{
		Token:       "test-token",
		APIEndpoint: server.URL,
	}
	client, err := cfg.Client()
	if err != nil {
		t.Fatalf("error building client: %s", err)
	}

	d := ResourceDigitalOceanDroplet().TestResourceData()
	_ = d.Set("image", "debian-13-x64")
	_ = d.Set("name", "esc-25905-private")
	_ = d.Set("region", "ams3")
	_ = d.Set("size", "s-1vcpu-1gb")
	_ = d.Set("public_networking", false)

	if diags := resourceDigitalOceanDropletCreate(context.Background(), d, client); diags.HasError() {
		t.Fatalf("create returned error: %v", diags)
	}
	if got := d.Get("ipv4_address").(string); got != "" {
		t.Fatalf("expected empty public ipv4_address for private droplet, got %q", got)
	}
	if got := d.Get("ipv4_address_private").(string); got != "10.108.0.20" {
		t.Fatalf("ipv4_address_private = %q, expected 10.108.0.20", got)
	}
}
