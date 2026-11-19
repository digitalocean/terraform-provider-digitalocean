package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"slices"
	"strings"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	testAppID                   = "6a7c6e1e-3b5d-4a8e-9f62-1d0c7b2a9e41"
	testSecretCiphertext        = "EV[1:zqiRIeaaYK/NqctZDYzy6t0pTrtRDeq8:wqGpZRrsKN5nPhWQrS479cfBiXT0WQ==]"
	testChangedSecretCiphertext = "EV[1:Y2hhbmdlZG5vbmNlZm9ydGVzdA==:Y2hhbmdlZGNpcGhlcnRleHQ=]"
)

func TestAppSecretEnvsAreNotPlannedAfterRefresh(t *testing.T) {
	client := newFakeAppsAPI(t, appWithEncryptedSecretEnvs(testSecretCiphertext))
	appConfig := appConfigWithSecretEnvs("v1")
	state := refreshedApp(t, client, createdApp(t, client, appConfig))

	plan := planApp(t, client, state, appConfig)

	for name, envSet := range map[string]string{
		"app-wide":       "spec.0.env",
		"service":        "spec.0.service.0.env",
		"second service": "spec.0.service.1.env",
		"static site":    "spec.0.static_site.0.env",
		"worker":         "spec.0.worker.0.env",
		"job":            "spec.0.job.0.env",
		"function":       "spec.0.function.0.env",
	} {
		t.Run(name+" secret is left unchanged", func(t *testing.T) {
			if changes := plannedChanges(plan, envSet); len(changes) > 0 {
				t.Errorf("plan changes %v, expected no changes", changes)
			}
		})
	}
}

func TestAppRejectedUpdateKeepsSecretEnvRotationPlanned(t *testing.T) {
	client := newFakeAppsAPI(t, appWithEncryptedSecretEnvs(testSecretCiphertext))
	appConfig := appConfigWithSecretEnvs("v1")
	state := refreshedApp(t, client, createdApp(t, client, appConfig))
	rotatedConfig := appConfigWithSecretEnvs("v2")

	rotation := planApp(t, client, state, rotatedConfig)

	afterUpdate, diags := ResourceDigitalOceanApp().Apply(context.Background(), state, rotation, client)
	if !diags.HasError() {
		t.Fatalf("update succeeded, expected the fake API to reject it")
	}

	plan := planApp(t, client, refreshedApp(t, client, afterUpdate), rotatedConfig)
	if changes := plannedChanges(plan, "spec.0.service.0.env"); len(changes) == 0 {
		t.Errorf("plan has no changes after the rejected update, expected the rotated secret")
	}
}

func TestAppSecretEnvChangedOutsideTerraformIsPlanned(t *testing.T) {
	appConfig := appConfigWithSecretEnvs("v1")
	initialClient := newFakeAppsAPI(t, appWithEncryptedSecretEnvs(testSecretCiphertext))
	state := createdApp(t, initialClient, appConfig)

	changedClient := newFakeAppsAPI(t, appWithEncryptedSecretEnvs(testChangedSecretCiphertext))
	plan := planApp(t, changedClient, refreshedApp(t, changedClient, state), appConfig)

	if changes := plannedChanges(plan, "spec.0.service.0.env"); len(changes) == 0 {
		t.Errorf("plan has no changes, expected to restore the secret changed outside Terraform")
	}
}

func appConfigWithSecretEnvs(version string) map[string]any {
	secretEnv := func(block string) []any {
		return []any{map[string]any{
			"key":   "TOKEN",
			"value": block + "-plaintext-" + version,
			"scope": "RUN_TIME",
			"type":  "SECRET",
		}}
	}
	components := func(names ...string) []any {
		var components []any
		for _, name := range names {
			components = append(components, map[string]any{
				"name": name,
				"env":  secretEnv(name),
			})
		}
		return components
	}

	return map[string]any{
		"spec": []any{map[string]any{
			"name":        "secret-envs",
			"env":         secretEnv("app"),
			"service":     components("web", "admin"),
			"static_site": components("site"),
			"worker":      components("queue"),
			"job":         components("migrate"),
			"function":    components("api"),
		}},
	}
}

func appWithEncryptedSecretEnvs(ciphertext string) *godo.App {
	secretEnvs := func() []*godo.AppVariableDefinition {
		return []*godo.AppVariableDefinition{{
			Key:   "TOKEN",
			Value: ciphertext,
			Scope: godo.AppVariableScope_RunTime,
			Type:  godo.AppVariableType_Secret,
		}}
	}

	return &godo.App{
		ID:               testAppID,
		ActiveDeployment: &godo.Deployment{ID: "b3fd7c84-5a1f-4b0e-8a39-6c2f1e7d9a50"},
		Spec: &godo.AppSpec{
			Name:        "secret-envs",
			Envs:        secretEnvs(),
			Services:    []*godo.AppServiceSpec{{Name: "web", Envs: secretEnvs()}, {Name: "admin", Envs: secretEnvs()}},
			StaticSites: []*godo.AppStaticSiteSpec{{Name: "site", Envs: secretEnvs()}},
			Workers:     []*godo.AppWorkerSpec{{Name: "queue", Envs: secretEnvs()}},
			Jobs:        []*godo.AppJobSpec{{Name: "migrate", Envs: secretEnvs()}},
			Functions:   []*godo.AppFunctionsSpec{{Name: "api", Envs: secretEnvs()}},
		},
	}
}

// newFakeAppsAPI serves app and rejects every update.
func newFakeAppsAPI(t *testing.T, app *godo.App) *config.CombinedConfig {
	t.Helper()

	getBody, err := json.Marshal(map[string]*godo.App{"app": app})
	if err != nil {
		t.Fatalf("error encoding app: %s", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v2/apps/"+app.ID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(getBody))
	})
	mux.HandleFunc("PUT /v2/apps/"+app.ID, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprint(w, `{"id":"unprocessable_entity","message":"update rejected"}`)
	})

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client, err := (&config.Config{Token: "test-token", APIEndpoint: server.URL}).Client()
	if err != nil {
		t.Fatalf("error building client: %s", err)
	}
	return client
}

func createdApp(t *testing.T, client *config.CombinedConfig, raw map[string]any) *terraform.InstanceState {
	t.Helper()
	data := schema.TestResourceDataRaw(t, ResourceDigitalOceanApp().Schema, raw)
	data.SetId(testAppID)
	readApp(t, data, client)
	return data.State()
}

func refreshedApp(t *testing.T, client *config.CombinedConfig, state *terraform.InstanceState) *terraform.InstanceState {
	t.Helper()
	data := ResourceDigitalOceanApp().Data(state)
	readApp(t, data, client)
	return data.State()
}

func readApp(t *testing.T, data *schema.ResourceData, client *config.CombinedConfig) {
	t.Helper()
	if diags := resourceDigitalOceanAppRead(context.Background(), data, client); diags.HasError() {
		t.Fatalf("read returned error: %v", diags)
	}
}

func planApp(t *testing.T, client *config.CombinedConfig, state *terraform.InstanceState, raw map[string]any) *terraform.InstanceDiff {
	t.Helper()
	plan, err := ResourceDigitalOceanApp().SimpleDiff(context.Background(), state, terraform.NewResourceConfigRaw(raw), client)
	if err != nil {
		t.Fatalf("plan returned error: %s", err)
	}
	return plan
}

func plannedChanges(plan *terraform.InstanceDiff, path string) []string {
	var changes []string
	for attribute := range plan.Attributes {
		if strings.HasPrefix(attribute, path+".") {
			changes = append(changes, attribute)
		}
	}
	slices.Sort(changes)
	return changes
}
