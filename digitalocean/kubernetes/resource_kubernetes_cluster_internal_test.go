package kubernetes

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const testClusterID = "1c2f814e-33f2-4284-9104-6799104ceb0c"

// startSSORequiredClusterAPI serves an SSO-required cluster and counts credentials requests.
func startSSORequiredClusterAPI(t *testing.T, credentialsGets *atomic.Int32) *config.CombinedConfig {
	t.Helper()

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	mux.HandleFunc("/v2/kubernetes/clusters/"+testClusterID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"kubernetes_cluster": {
			"id": %q,
			"name": "sso-cluster",
			"region": "ams3",
			"status": {"state": "running"},
			"created_at": "2026-03-20T00:00:00Z",
			"updated_at": "2026-03-20T00:00:00Z",
			"maintenance_policy": {"start_time": "00:00", "day": "any"},
			"sso": {"enabled": true, "required": true}
		}}`, testClusterID)
	})

	mux.HandleFunc("/v2/kubernetes/clusters/"+testClusterID+"/credentials", func(w http.ResponseWriter, r *http.Request) {
		credentialsGets.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, `{"id": "unprocessable_entity", "message": "cluster requires SSO, but a token kubeconfig was requested"}`)
	})

	cfg := &config.Config{
		Token:            "test-token",
		APIEndpoint:      server.URL,
		HTTPRetryMax:     1,
		HTTPRetryWaitMin: 0.01,
		HTTPRetryWaitMax: 0.01,
	}
	client, err := cfg.Client()
	if err != nil {
		t.Fatalf("error building client: %s", err)
	}

	return client
}

func TestKubernetesClusterReadSSORequired(t *testing.T) {
	t.Parallel()

	for name, kubeConfig := range map[string][]interface{}{
		"no stored credentials": nil,
		"valid stored credentials": {map[string]interface{}{
			"token":      "old-token",
			"expires_at": time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
		}},
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			var credentialsGets atomic.Int32
			client := startSSORequiredClusterAPI(t, &credentialsGets)

			d := schema.TestResourceDataRaw(t, ResourceDigitalOceanKubernetesCluster().Schema, map[string]interface{}{})
			d.SetId(testClusterID)
			if err := d.Set("kube_config", kubeConfig); err != nil {
				t.Fatalf("error setting kube_config: %s", err)
			}

			if diags := resourceDigitalOceanKubernetesClusterRead(context.Background(), d, client); diags.HasError() {
				t.Fatalf("read returned error: %v", diags)
			}

			if n := credentialsGets.Load(); n != 0 {
				t.Errorf("credentials requests = %d, expected 0", n)
			}
			if kc := d.Get("kube_config").([]interface{}); len(kc) != 0 {
				t.Errorf("kube_config = %v, expected it to be empty", kc)
			}
		})
	}
}
