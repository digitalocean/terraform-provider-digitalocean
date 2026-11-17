package kubernetes

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	testClusterID   = "1c2f814e-33f2-4284-9104-6799104ceb0c"
	testClusterName = "sso-cluster"
)

// fakeKubernetesAPI serves the cluster endpoints used by the cluster resource and data source.
type fakeKubernetesAPI struct {
	mu          sync.Mutex
	ssoRequired bool
	// credentialsStatus is the status returned by the credentials endpoint; 0 means 200.
	credentialsStatus int
	credentialsGets   int
}

func (api *fakeKubernetesAPI) clusterJSON() string {
	return fmt.Sprintf(`{
		"id": %q,
		"name": %q,
		"region": "ams3",
		"version": "1.35.1-do.0",
		"endpoint": "https://%s.k8s.ondigitalocean.com",
		"status": {"state": "running"},
		"created_at": "2026-03-20T00:00:00Z",
		"updated_at": "2026-03-20T00:00:00Z",
		"maintenance_policy": {"start_time": "00:00", "day": "any"},
		"sso": {"enabled": true, "required": %t},
		"node_pools": [{"id": "pool-1", "name": "default", "size": "s-1vcpu-2gb", "count": 1, "tags": [%q]}]
	}`, testClusterID, testClusterName, testClusterID, api.ssoRequired, DigitaloceanKubernetesDefaultNodePoolTag)
}

func (api *fakeKubernetesAPI) start(t *testing.T) *config.CombinedConfig {
	t.Helper()

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	mux.HandleFunc("/v2/kubernetes/clusters", func(w http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		defer api.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodPost {
			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, `{"kubernetes_cluster": %s}`, api.clusterJSON())
			return
		}
		fmt.Fprintf(w, `{"kubernetes_clusters": [%s], "meta": {"total": 1}}`, api.clusterJSON())
	})

	mux.HandleFunc("/v2/kubernetes/clusters/"+testClusterID, func(w http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		defer api.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"kubernetes_cluster": %s}`, api.clusterJSON())
	})

	mux.HandleFunc("/v2/kubernetes/clusters/"+testClusterID+"/credentials", func(w http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		defer api.mu.Unlock()

		api.credentialsGets++
		w.Header().Set("Content-Type", "application/json")
		if api.ssoRequired {
			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprint(w, `{"id": "unprocessable_entity", "message": "cluster requires SSO, but a token kubeconfig was requested"}`)
			return
		}
		if api.credentialsStatus != 0 {
			w.WriteHeader(api.credentialsStatus)
			fmt.Fprint(w, `{"id": "unprocessable_entity", "message": "unprocessable"}`)
			return
		}
		fmt.Fprintf(w, `{
			"server": "https://%s.k8s.ondigitalocean.com",
			"certificate_authority_data": "Y2EtZGF0YQ==",
			"token": "fresh-token",
			"expires_at": %q
		}`, testClusterID, time.Now().Add(time.Hour).UTC().Format(time.RFC3339))
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

func newTestClusterData(t *testing.T, kubeConfig []interface{}) *schema.ResourceData {
	t.Helper()

	raw := map[string]interface{}{
		"name":    testClusterName,
		"region":  "ams3",
		"version": "1.35.1-do.0",
		"node_pool": []interface{}{map[string]interface{}{
			"name":       "default",
			"size":       "s-1vcpu-2gb",
			"node_count": 1,
		}},
		"sso": []interface{}{map[string]interface{}{
			"enabled":  true,
			"required": true,
		}},
	}
	d := schema.TestResourceDataRaw(t, ResourceDigitalOceanKubernetesCluster().Schema, raw)
	d.SetId(testClusterID)
	if kubeConfig != nil {
		if err := d.Set("kube_config", kubeConfig); err != nil {
			t.Fatalf("error setting kube_config: %s", err)
		}
	}
	return d
}

func storedKubeConfig(expiresAt time.Time) []interface{} {
	return []interface{}{map[string]interface{}{
		"host":       "https://old.k8s.ondigitalocean.com",
		"token":      "old-token",
		"expires_at": expiresAt.UTC().Format(time.RFC3339),
	}}
}

func assertNoKubeConfig(t *testing.T, d *schema.ResourceData) {
	t.Helper()

	if kc := d.Get("kube_config").([]interface{}); len(kc) != 0 {
		t.Errorf("kube_config = %v, expected it to be empty", kc)
	}
}

func TestKubernetesClusterReadSSORequired(t *testing.T) {
	t.Parallel()

	for name, kubeConfig := range map[string][]interface{}{
		"no stored credentials":      nil,
		"expired stored credentials": storedKubeConfig(time.Now().Add(-time.Hour)),
		"valid stored credentials":   storedKubeConfig(time.Now().Add(time.Hour)),
	} {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			api := &fakeKubernetesAPI{ssoRequired: true}
			client := api.start(t)
			d := newTestClusterData(t, kubeConfig)

			if diags := resourceDigitalOceanKubernetesClusterRead(context.Background(), d, client); diags.HasError() {
				t.Fatalf("read returned error: %v", diags)
			}

			if api.credentialsGets != 0 {
				t.Errorf("credentials requests = %d, expected 0", api.credentialsGets)
			}
			if d.Get("sso.0.required") != true {
				t.Error("expected sso.0.required to be true")
			}
			assertNoKubeConfig(t, d)
		})
	}
}

func TestKubernetesClusterReadRefreshesExpiredCredentials(t *testing.T) {
	t.Parallel()

	api := &fakeKubernetesAPI{}
	client := api.start(t)
	d := newTestClusterData(t, storedKubeConfig(time.Now().Add(-time.Hour)))

	if diags := resourceDigitalOceanKubernetesClusterRead(context.Background(), d, client); diags.HasError() {
		t.Fatalf("read returned error: %v", diags)
	}

	if api.credentialsGets != 1 {
		t.Errorf("credentials requests = %d, expected 1", api.credentialsGets)
	}
	if token := d.Get("kube_config.0.token"); token != "fresh-token" {
		t.Errorf("kube_config.0.token = %q, expected %q", token, "fresh-token")
	}
}

func TestKubernetesClusterReadKeepsValidCredentials(t *testing.T) {
	t.Parallel()

	api := &fakeKubernetesAPI{}
	client := api.start(t)
	d := newTestClusterData(t, storedKubeConfig(time.Now().Add(time.Hour)))

	if diags := resourceDigitalOceanKubernetesClusterRead(context.Background(), d, client); diags.HasError() {
		t.Fatalf("read returned error: %v", diags)
	}

	if api.credentialsGets != 0 {
		t.Errorf("credentials requests = %d, expected 0", api.credentialsGets)
	}
	if token := d.Get("kube_config.0.token"); token != "old-token" {
		t.Errorf("kube_config.0.token = %q, expected %q", token, "old-token")
	}
}

func TestKubernetesClusterReadCredentialsErrorWithoutSSORequired(t *testing.T) {
	t.Parallel()

	api := &fakeKubernetesAPI{credentialsStatus: http.StatusUnprocessableEntity}
	client := api.start(t)
	d := newTestClusterData(t, nil)

	diags := resourceDigitalOceanKubernetesClusterRead(context.Background(), d, client)
	if !diags.HasError() {
		t.Fatal("expected read to return an error")
	}
	if !strings.Contains(diags[0].Summary, "Unable to fetch Kubernetes credentials") {
		t.Errorf("error = %q, expected a credentials error", diags[0].Summary)
	}
}

func TestKubernetesClusterImportSSORequired(t *testing.T) {
	t.Parallel()

	api := &fakeKubernetesAPI{ssoRequired: true}
	client := api.start(t)
	d := schema.TestResourceDataRaw(t, ResourceDigitalOceanKubernetesCluster().Schema, map[string]interface{}{})
	d.SetId(testClusterID)

	imported, err := resourceDigitalOceanKubernetesClusterImportState(context.Background(), d, client)
	if err != nil {
		t.Fatalf("import returned error: %s", err)
	}
	if diags := resourceDigitalOceanKubernetesClusterRead(context.Background(), imported[0], client); diags.HasError() {
		t.Fatalf("read after import returned error: %v", diags)
	}

	if api.credentialsGets != 0 {
		t.Errorf("credentials requests = %d, expected 0", api.credentialsGets)
	}
	if name := imported[0].Get("name"); name != testClusterName {
		t.Errorf("name = %q, expected %q", name, testClusterName)
	}
	assertNoKubeConfig(t, imported[0])
}

func TestKubernetesClusterDataSourceSSORequired(t *testing.T) {
	t.Parallel()

	api := &fakeKubernetesAPI{ssoRequired: true}
	client := api.start(t)
	d := schema.TestResourceDataRaw(t, DataSourceDigitalOceanKubernetesCluster().Schema, map[string]interface{}{
		"name": testClusterName,
	})

	if diags := dataSourceDigitalOceanKubernetesClusterRead(context.Background(), d, client); diags.HasError() {
		t.Fatalf("data source read returned error: %v", diags)
	}

	if api.credentialsGets != 0 {
		t.Errorf("credentials requests = %d, expected 0", api.credentialsGets)
	}
	if d.Id() != testClusterID {
		t.Errorf("id = %q, expected %q", d.Id(), testClusterID)
	}
	assertNoKubeConfig(t, d)
}

func TestKubernetesClusterCreateSSORequired(t *testing.T) {
	t.Parallel()

	api := &fakeKubernetesAPI{ssoRequired: true}
	client := api.start(t)
	d := newTestClusterData(t, nil)
	d.SetId("")

	if diags := resourceDigitalOceanKubernetesClusterCreate(context.Background(), d, client); diags.HasError() {
		t.Fatalf("create returned error: %v", diags)
	}

	if d.Id() != testClusterID {
		t.Errorf("id = %q, expected %q", d.Id(), testClusterID)
	}
	if api.credentialsGets != 0 {
		t.Errorf("credentials requests = %d, expected 0", api.credentialsGets)
	}
	assertNoKubeConfig(t, d)
}
