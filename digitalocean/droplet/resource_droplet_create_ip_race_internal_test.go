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

func testDropletClient(t *testing.T, serverURL string) interface{} {
	t.Helper()
	client, err := (&config.Config{
		Token:       "test-token",
		APIEndpoint: serverURL,
	}).Client()
	if err != nil {
		t.Fatalf("error building client: %s", err)
	}
	return client
}

// Create should keep waiting if the droplet is active but has no public IP yet.
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
		// First GET is the status wait (active, no IP). Later GETs include the IP.
		if n >= 2 {
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

	d := ResourceDigitalOceanDroplet().TestResourceData()
	_ = d.Set("image", "debian-13-x64")
	_ = d.Set("name", "esc-25905-fixed")
	_ = d.Set("region", "ams3")
	_ = d.Set("size", "s-1vcpu-1gb")

	if diags := resourceDigitalOceanDropletCreate(context.Background(), d, testDropletClient(t, server.URL)); diags.HasError() {
		t.Fatalf("create returned error: %v", diags)
	}
	if got := d.Get("ipv4_address").(string); got != wantIP {
		t.Fatalf("ipv4_address = %q, expected %q (gets=%d)", got, wantIP, gets.Load())
	}
	if gets.Load() < 2 {
		t.Fatalf("expected create to poll through missing-IP window, gets=%d", gets.Load())
	}
}

// waitForDropletPublicIPv4 should error if a public IP never shows up.
func TestWaitForDropletPublicIPv4_TimesOut(t *testing.T) {
	const dropletID = 598498974

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

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
				"volume_ids": []
			}
		}`, dropletID)
	})

	d := ResourceDigitalOceanDroplet().TestResourceData()
	d.SetId(fmt.Sprintf("%d", dropletID))

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_, err := waitForDropletPublicIPv4(ctx, d, testDropletClient(t, server.URL))
	if err == nil {
		t.Fatal("expected timeout waiting for public IPv4")
	}
}

// private-only droplets should not wait for a public IP.
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

	d := ResourceDigitalOceanDroplet().TestResourceData()
	_ = d.Set("image", "debian-13-x64")
	_ = d.Set("name", "esc-25905-private")
	_ = d.Set("region", "ams3")
	_ = d.Set("size", "s-1vcpu-1gb")
	_ = d.Set("public_networking", false)

	if diags := resourceDigitalOceanDropletCreate(context.Background(), d, testDropletClient(t, server.URL)); diags.HasError() {
		t.Fatalf("create returned error: %v", diags)
	}
	if got := d.Get("ipv4_address").(string); got != "" {
		t.Fatalf("expected empty public ipv4_address for private droplet, got %q", got)
	}
	if got := d.Get("ipv4_address_private").(string); got != "10.108.0.20" {
		t.Fatalf("ipv4_address_private = %q, expected 10.108.0.20", got)
	}
}
