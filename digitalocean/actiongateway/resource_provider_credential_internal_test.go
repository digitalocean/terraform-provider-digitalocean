package actiongateway

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// fakeProviderCredentialAPI keeps provider credentials in memory and applies
// the rules the public API documents: write-only secrets, an OAuth client ID
// that every secret or base URL change must repeat, sharing that cannot be
// undone, and a delete that takes the credential's connections with it.
type fakeProviderCredentialAPI struct {
	mu          sync.Mutex
	next        int
	credentials map[string]*fakeProviderCredential
	// connections counts the connections that use each credential.
	connections map[string]int
	// report overrides fields of every credential the API returns, to stand
	// in for a value the API stores in another spelling.
	report  map[string]interface{}
	creates []map[string]interface{}
	patches []map[string]interface{}
	deletes []string
}

type fakeProviderCredential struct {
	id           string
	provider     string
	displayName  string
	kind         string
	baseURL      string
	redirectURL  string
	reusable     bool
	apiKey       string
	clientID     string
	clientSecret string
	revision     int
}

func newFakeProviderCredentialAPI() *fakeProviderCredentialAPI {
	return &fakeProviderCredentialAPI{
		credentials: map[string]*fakeProviderCredential{},
		connections: map[string]int{},
	}
}

func (api *fakeProviderCredentialAPI) render(cred *fakeProviderCredential, status string) map[string]interface{} {
	out := map[string]interface{}{
		"id":               cred.id,
		"provider":         cred.provider,
		"display_name":     cred.displayName,
		"kind":             cred.kind,
		"base_url":         cred.baseURL,
		"status":           status,
		"created_at":       "2026-11-01T12:00:00Z",
		"updated_at":       fmt.Sprintf("2026-11-01T12:%02d:00Z", cred.revision),
		"reusable_by_team": cred.reusable,
		"redirect_url":     cred.redirectURL,
		"scope":            "user",
	}
	if cred.reusable {
		out["scope"] = "team"
	} else {
		out["owning_user_id"] = "2f8e1c3a-6b4d-4f7e-9a1c-5d3b7e9f1a2c"
	}
	for k, v := range api.report {
		out[k] = v
	}
	return out
}

func (api *fakeProviderCredentialAPI) respond(t *testing.T, w http.ResponseWriter, cred *fakeProviderCredential) {
	writeJSON(t, w, http.StatusOK, map[string]interface{}{
		"credential":          api.render(cred, "PROVIDER_CREDENTIAL_STATUS_ACTIVE"),
		"revoked_connections": 0,
	})
}

func writeBadRequest(t *testing.T, w http.ResponseWriter, message string) {
	writeJSON(t, w, http.StatusBadRequest, map[string]string{"id": "bad_request", "message": message})
}

// add registers a credential the way another client would.
func (api *fakeProviderCredentialAPI) add(cred fakeProviderCredential) string {
	api.mu.Lock()
	defer api.mu.Unlock()
	api.next++
	cred.id = fmt.Sprintf("cred-%d", api.next)
	api.credentials[cred.id] = &cred
	return cred.id
}

func (api *fakeProviderCredentialAPI) get(id string) fakeProviderCredential {
	api.mu.Lock()
	defer api.mu.Unlock()
	if cred, ok := api.credentials[id]; ok {
		return *cred
	}
	return fakeProviderCredential{}
}

func (api *fakeProviderCredentialAPI) start(t *testing.T) *config.CombinedConfig {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v2/action-gateway/provider-credentials", func(w http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		defer api.mu.Unlock()
		body := decodeJSON(t, r)
		api.creates = append(api.creates, body)

		cred := &fakeProviderCredential{}
		cred.provider, _ = body["provider"].(string)
		cred.displayName, _ = body["display_name"].(string)
		cred.kind, _ = body["kind"].(string)
		cred.apiKey, _ = body["api_key"].(string)
		cred.clientID, _ = body["client_id"].(string)
		cred.clientSecret, _ = body["client_secret"].(string)
		cred.redirectURL, _ = body["redirect_url"].(string)
		cred.reusable, _ = body["reusable_by_team"].(bool)
		baseURL, _ := body["base_url"].(string)
		params, _ := body["credential_parameters"].(map[string]interface{})

		switch cred.kind {
		case "PROVIDER_CREDENTIAL_KIND_API_KEY":
			if cred.apiKey == "" || cred.clientID != "" || cred.clientSecret != "" || cred.redirectURL != "" {
				writeBadRequest(t, w, "an API-key credential requires api_key and accepts no OAuth client fields")
				return
			}
		case "PROVIDER_CREDENTIAL_KIND_OAUTH_CLIENT":
			if cred.clientID == "" || cred.clientSecret == "" || cred.apiKey != "" || params != nil {
				writeBadRequest(t, w, "an OAuth client credential requires client_id and client_secret")
				return
			}
		default:
			writeBadRequest(t, w, "unsupported kind")
			return
		}
		if baseURL != "" && params != nil {
			writeBadRequest(t, w, "base_url cannot be sent alongside credential_parameters")
			return
		}
		for _, other := range api.credentials {
			if other.provider == cred.provider && strings.EqualFold(other.displayName, cred.displayName) {
				writeJSON(t, w, http.StatusConflict, map[string]string{"id": "conflict", "message": "display name already in use"})
				return
			}
			if cred.clientID != "" && other.provider == cred.provider && other.clientID == cred.clientID {
				writeJSON(t, w, http.StatusConflict, map[string]string{"id": "conflict", "message": "client ID already registered"})
				return
			}
		}
		cred.baseURL = fakeBaseURL(baseURL, params)

		api.next++
		cred.id = fmt.Sprintf("cred-%d", api.next)
		api.credentials[cred.id] = cred
		api.respond(t, w, cred)
	})
	mux.HandleFunc("GET /v2/action-gateway/provider-credentials", func(w http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		defer api.mu.Unlock()
		provider := r.URL.Query().Get("provider")
		creds := []interface{}{}
		for _, cred := range api.credentials {
			if provider == "" || strings.EqualFold(cred.provider, provider) {
				creds = append(creds, api.render(cred, "PROVIDER_CREDENTIAL_STATUS_ACTIVE"))
			}
		}
		writeJSON(t, w, http.StatusOK, map[string]interface{}{"credentials": creds})
	})
	mux.HandleFunc("GET /v2/action-gateway/provider-credentials/{id}", func(w http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		defer api.mu.Unlock()
		cred, ok := api.credentials[r.PathValue("id")]
		if !ok {
			writeNotFound(t, w)
			return
		}
		api.respond(t, w, cred)
	})
	mux.HandleFunc("PATCH /v2/action-gateway/provider-credentials/{id}", func(w http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		defer api.mu.Unlock()
		cred, ok := api.credentials[r.PathValue("id")]
		if !ok {
			writeNotFound(t, w)
			return
		}
		body := decodeJSON(t, r)
		api.patches = append(api.patches, body)
		if len(body) == 0 {
			writeBadRequest(t, w, "at least one field is required")
			return
		}

		oauthClient := cred.kind == "PROVIDER_CREDENTIAL_KIND_OAUTH_CLIENT"
		clientID, hasClientID := body["client_id"].(string)
		_, hasSecret := body["client_secret"]
		_, hasBaseURL := body["base_url"]
		_, hasAPIKey := body["api_key"]
		if oauthClient {
			if hasAPIKey {
				writeBadRequest(t, w, "api_key is only valid for API-key credentials")
				return
			}
			if (hasSecret || hasBaseURL) && !hasClientID {
				writeBadRequest(t, w, "client_id is required to change an OAuth client")
				return
			}
			if hasClientID && clientID != cred.clientID {
				writeBadRequest(t, w, "client_id does not match the registered OAuth client")
				return
			}
		} else if hasClientID || hasSecret {
			writeBadRequest(t, w, "OAuth client fields are only valid for OAuth client credentials")
			return
		}
		if shared, ok := body["reusable_by_team"].(bool); ok {
			if !shared {
				writeBadRequest(t, w, "sharing cannot be undone")
				return
			}
			cred.reusable = true
		}

		if name, ok := body["display_name"].(string); ok {
			cred.displayName = name
		}
		if key, ok := body["api_key"].(string); ok {
			cred.apiKey = key
		}
		if secret, ok := body["client_secret"].(string); ok {
			cred.clientSecret = secret
		}
		baseURL, _ := body["base_url"].(string)
		params, _ := body["credential_parameters"].(map[string]interface{})
		if hasBaseURL || params != nil {
			cred.baseURL = fakeBaseURL(baseURL, params)
		}
		cred.revision++
		api.respond(t, w, cred)
	})
	mux.HandleFunc("DELETE /v2/action-gateway/provider-credentials/{id}", func(w http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		defer api.mu.Unlock()
		id := r.PathValue("id")
		cred, ok := api.credentials[id]
		if !ok {
			writeNotFound(t, w)
			return
		}
		api.deletes = append(api.deletes, id)
		delete(api.credentials, id)
		writeJSON(t, w, http.StatusOK, map[string]interface{}{
			"credential":          api.render(cred, "PROVIDER_CREDENTIAL_STATUS_DELETED"),
			"revoked_connections": api.connections[id],
		})
	})
	return newTestMeta(t, mux)
}

// fakeBaseURL stores a base URL as its normalized origin, or derives one from
// a subdomain parameter.
func fakeBaseURL(baseURL string, params map[string]interface{}) string {
	if subdomain, ok := params["subdomain"].(string); ok {
		return "https://" + subdomain + ".example.com"
	}
	return strings.ToLower(strings.TrimSuffix(baseURL, "/"))
}

func providerCredentialRawConfig(r *schema.Resource, config map[string]interface{}) cty.Value {
	attrs := map[string]cty.Value{}
	for name, value := range config {
		switch v := value.(type) {
		case string:
			attrs[name] = cty.StringVal(v)
		case bool:
			attrs[name] = cty.BoolVal(v)
		case map[string]interface{}:
			entries := map[string]cty.Value{}
			for k, item := range v {
				entries[k] = cty.StringVal(item.(string))
			}
			attrs[name] = cty.MapVal(entries)
		}
	}
	return testRawConfig(r, attrs)
}

// planProviderCredential plans one configuration with the raw configuration
// Terraform would send, which the resource's plan-time checks read.
func planProviderCredential(t *testing.T, r *schema.Resource, meta *config.CombinedConfig, state *terraform.InstanceState, config map[string]interface{}) *terraform.InstanceDiff {
	t.Helper()
	diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), meta)
	if err != nil {
		t.Fatalf("error planning: %s", err)
	}
	return diff
}

// applyProviderCredential plans and applies one configuration, the way
// Terraform calls the provider, and returns the state and the diagnostics.
func applyProviderCredential(t *testing.T, r *schema.Resource, meta *config.CombinedConfig, state *terraform.InstanceState, config map[string]interface{}) (*terraform.InstanceState, diag.Diagnostics) {
	t.Helper()
	diff := planProviderCredential(t, r, meta, state, config)
	if diff == nil || diff.Empty() {
		return state, nil
	}
	if state != nil && diff.RequiresNew() {
		t.Fatalf("expected an in-place change, got a replacement:\n%s", describeDiff(diff))
	}
	diff.RawConfig = providerCredentialRawConfig(r, config)
	return r.Apply(context.Background(), state, diff, meta)
}

func mustApplyProviderCredential(t *testing.T, r *schema.Resource, meta *config.CombinedConfig, state *terraform.InstanceState, config map[string]interface{}) *terraform.InstanceState {
	t.Helper()
	next, diags := applyProviderCredential(t, r, meta, state, config)
	requireNoErrors(t, diags)
	return next
}

func requireEmptyPlan(t *testing.T, r *schema.Resource, meta *config.CombinedConfig, state *terraform.InstanceState, config map[string]interface{}) {
	t.Helper()
	refreshed, diags := r.RefreshWithoutUpgrade(context.Background(), state, meta)
	requireNoErrors(t, diags)
	if plan := planProviderCredential(t, r, meta, refreshed, config); plan != nil && !plan.Empty() {
		t.Errorf("expected an empty plan, got:\n%s", describeDiff(plan))
	}
}

func importProviderCredential(t *testing.T, r *schema.Resource, meta *config.CombinedConfig, id string) *terraform.InstanceState {
	t.Helper()
	ctx := context.Background()
	imported, err := r.Importer.StateContext(ctx, r.Data(&terraform.InstanceState{ID: id}), meta)
	if err != nil {
		t.Fatalf("error importing: %s", err)
	}
	state, diags := r.RefreshWithoutUpgrade(ctx, imported[0].State(), meta)
	requireNoErrors(t, diags)
	return state
}

func apiKeyCredentialConfig() map[string]interface{} {
	return map[string]interface{}{
		"tool_provider": "linear",
		"kind":          "api_key",
		"display_name":  "Production key",
		"api_key":       "lin_api_one",
	}
}

func oauthClientCredentialConfig() map[string]interface{} {
	return map[string]interface{}{
		"tool_provider": "github",
		"kind":          "oauth_client",
		"display_name":  "Production OAuth app",
		"client_id":     "Iv1.0123456789abcdef",
		"client_secret": "secret-one",
	}
}

func withValues(base map[string]interface{}, overrides map[string]interface{}) map[string]interface{} {
	out := make(map[string]interface{}, len(base)+len(overrides))
	for k, v := range base {
		out[k] = v
	}
	for k, v := range overrides {
		if v == nil {
			delete(out, k)
		} else {
			out[k] = v
		}
	}
	return out
}

func TestProviderCredentialCreateAPIKey(t *testing.T) {
	t.Parallel()

	api := newFakeProviderCredentialAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayProviderCredential()
	config := apiKeyCredentialConfig()

	state := mustApplyProviderCredential(t, r, meta, nil, config)

	want := map[string]interface{}{
		"provider":     "linear",
		"kind":         "PROVIDER_CREDENTIAL_KIND_API_KEY",
		"display_name": "Production key",
		"api_key":      "lin_api_one",
	}
	if got := api.creates[0]; !reflect.DeepEqual(got, want) {
		t.Errorf("create request = %v, expected %v", got, want)
	}

	d := r.Data(state)
	checks := map[string]interface{}{
		"tool_provider":    "linear",
		"kind":             "api_key",
		"display_name":     "Production key",
		"api_key":          "lin_api_one",
		"status":           "active",
		"scope":            "user",
		"reusable_by_team": false,
		"owning_user_id":   "2f8e1c3a-6b4d-4f7e-9a1c-5d3b7e9f1a2c",
		"created_at":       "2026-11-01T12:00:00Z",
		"base_url":         "",
	}
	for key, want := range checks {
		if got := d.Get(key); got != want {
			t.Errorf("%s = %v, expected %v", key, got, want)
		}
	}
	if state.ID != "cred-1" {
		t.Errorf("id = %q, expected the ID the API returned", state.ID)
	}

	requireEmptyPlan(t, r, meta, state, config)
}

func TestProviderCredentialCreateOAuthClient(t *testing.T) {
	t.Parallel()

	api := newFakeProviderCredentialAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayProviderCredential()
	config := withValues(oauthClientCredentialConfig(), map[string]interface{}{
		"redirect_url":     "https://app.example.com/oauth/callback",
		"reusable_by_team": true,
	})

	state := mustApplyProviderCredential(t, r, meta, nil, config)

	want := map[string]interface{}{
		"provider":         "github",
		"kind":             "PROVIDER_CREDENTIAL_KIND_OAUTH_CLIENT",
		"display_name":     "Production OAuth app",
		"client_id":        "Iv1.0123456789abcdef",
		"client_secret":    "secret-one",
		"redirect_url":     "https://app.example.com/oauth/callback",
		"reusable_by_team": true,
	}
	if got := api.creates[0]; !reflect.DeepEqual(got, want) {
		t.Errorf("create request = %v, expected %v", got, want)
	}

	d := r.Data(state)
	if got := d.Get("scope"); got != "team" {
		t.Errorf("scope = %v, expected team for a shared credential", got)
	}
	if got := d.Get("owning_user_id"); got != "" {
		t.Errorf("owning_user_id = %v, expected none for a shared credential", got)
	}

	requireEmptyPlan(t, r, meta, state, config)
}

func TestProviderCredentialRotatesAPIKeyInPlace(t *testing.T) {
	t.Parallel()

	api := newFakeProviderCredentialAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayProviderCredential()

	state := mustApplyProviderCredential(t, r, meta, nil, apiKeyCredentialConfig())
	rotated := withValues(apiKeyCredentialConfig(), map[string]interface{}{"api_key": "lin_api_two"})
	next := mustApplyProviderCredential(t, r, meta, state, rotated)

	if next.ID != state.ID {
		t.Errorf("id = %q, expected the credential to keep %q", next.ID, state.ID)
	}
	if got, want := api.patches, []map[string]interface{}{{"api_key": "lin_api_two"}}; !reflect.DeepEqual(got, want) {
		t.Errorf("update requests = %v, expected %v", got, want)
	}
	if got := api.get(state.ID).apiKey; got != "lin_api_two" {
		t.Errorf("stored key = %q, expected the rotated key", got)
	}
	if len(api.deletes) != 0 {
		t.Errorf("deletes = %v, expected a rotation to delete nothing", api.deletes)
	}
	if got := r.Data(next).Get("updated_at"); got != "2026-11-01T12:01:00Z" {
		t.Errorf("updated_at = %v, expected the time of the rotation", got)
	}

	requireEmptyPlan(t, r, meta, next, rotated)
}

func TestProviderCredentialRotatesClientSecretWithClientID(t *testing.T) {
	t.Parallel()

	api := newFakeProviderCredentialAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayProviderCredential()

	state := mustApplyProviderCredential(t, r, meta, nil, oauthClientCredentialConfig())
	rotated := withValues(oauthClientCredentialConfig(), map[string]interface{}{"client_secret": "secret-two"})
	next := mustApplyProviderCredential(t, r, meta, state, rotated)

	want := []map[string]interface{}{{"client_id": "Iv1.0123456789abcdef", "client_secret": "secret-two"}}
	if got := api.patches; !reflect.DeepEqual(got, want) {
		t.Errorf("update requests = %v, expected %v", got, want)
	}
	if next.ID != state.ID || len(api.deletes) != 0 {
		t.Errorf("id = %q, deletes = %v, expected the credential to be kept", next.ID, api.deletes)
	}
	if got := api.get(state.ID).clientSecret; got != "secret-two" {
		t.Errorf("stored secret = %q, expected the rotated secret", got)
	}
}

func TestProviderCredentialRenameSendsOnlyDisplayName(t *testing.T) {
	t.Parallel()

	api := newFakeProviderCredentialAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayProviderCredential()

	state := mustApplyProviderCredential(t, r, meta, nil, oauthClientCredentialConfig())
	renamed := withValues(oauthClientCredentialConfig(), map[string]interface{}{"display_name": "Staging OAuth app"})
	next := mustApplyProviderCredential(t, r, meta, state, renamed)

	if got, want := api.patches, []map[string]interface{}{{"display_name": "Staging OAuth app"}}; !reflect.DeepEqual(got, want) {
		t.Errorf("update requests = %v, expected %v", got, want)
	}
	if got := r.Data(next).Get("display_name"); got != "Staging OAuth app" {
		t.Errorf("display_name = %v", got)
	}
}

// TestProviderCredentialReplacementArguments pins which changes replace the
// credential. A replacement deletes every connection that uses it, so no other
// change may plan one.
func TestProviderCredentialReplacementArguments(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		base        map[string]interface{}
		change      map[string]interface{}
		replacesFor string
	}{
		{name: "another tool provider", base: apiKeyCredentialConfig(), change: map[string]interface{}{"tool_provider": "jira"}, replacesFor: "tool_provider"},
		{name: "another OAuth client", base: oauthClientCredentialConfig(), change: map[string]interface{}{"client_id": "Iv1.fedcba9876543210"}, replacesFor: "client_id"},
		{name: "a redirect URL", base: oauthClientCredentialConfig(), change: map[string]interface{}{"redirect_url": "https://app.example.com/oauth/callback"}, replacesFor: "redirect_url"},
		{
			name:        "another kind",
			base:        apiKeyCredentialConfig(),
			change:      map[string]interface{}{"kind": "oauth_client", "api_key": nil, "client_id": "abc", "client_secret": "def"},
			replacesFor: "kind",
		},
		{name: "a rotated API key", base: apiKeyCredentialConfig(), change: map[string]interface{}{"api_key": "lin_api_two"}},
		{name: "a rotated client secret", base: oauthClientCredentialConfig(), change: map[string]interface{}{"client_secret": "secret-two"}},
		{name: "a new display name", base: apiKeyCredentialConfig(), change: map[string]interface{}{"display_name": "Renamed"}},
		{name: "sharing with the team", base: apiKeyCredentialConfig(), change: map[string]interface{}{"reusable_by_team": true}},
		{name: "a base URL", base: apiKeyCredentialConfig(), change: map[string]interface{}{"base_url": "https://api.example.com"}},
		{name: "credential parameters", base: apiKeyCredentialConfig(), change: map[string]interface{}{"credential_parameters": map[string]interface{}{"subdomain": "acme"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			api := newFakeProviderCredentialAPI()
			meta := api.start(t)
			r := ResourceDigitalOceanActionGatewayProviderCredential()

			state := mustApplyProviderCredential(t, r, meta, nil, tt.base)
			plan := planProviderCredential(t, r, meta, state, withValues(tt.base, tt.change))
			if plan == nil || plan.Empty() {
				t.Fatal("expected the change to be planned")
			}

			if tt.replacesFor == "" {
				if plan.RequiresNew() {
					t.Errorf("expected an in-place change, got a replacement:\n%s", describeDiff(plan))
				}
				return
			}
			if attr := plan.Attributes[tt.replacesFor]; attr == nil || !attr.RequiresNew {
				t.Errorf("expected %s to force a replacement, got:\n%s", tt.replacesFor, describeDiff(plan))
			}
		})
	}
}

func TestProviderCredentialSharingIsOneWay(t *testing.T) {
	t.Parallel()

	api := newFakeProviderCredentialAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayProviderCredential()

	state := mustApplyProviderCredential(t, r, meta, nil, apiKeyCredentialConfig())
	shared := withValues(apiKeyCredentialConfig(), map[string]interface{}{"reusable_by_team": true})
	next := mustApplyProviderCredential(t, r, meta, state, shared)

	if got, want := api.patches, []map[string]interface{}{{"reusable_by_team": true}}; !reflect.DeepEqual(got, want) {
		t.Errorf("update requests = %v, expected %v", got, want)
	}
	d := r.Data(next)
	if got := d.Get("scope"); got != "team" {
		t.Errorf("scope = %v, expected team after sharing", got)
	}
	if got := d.Get("owning_user_id"); got != "" {
		t.Errorf("owning_user_id = %v, expected none after sharing", got)
	}
	requireEmptyPlan(t, r, meta, next, shared)

	plan := planProviderCredential(t, r, meta, next, apiKeyCredentialConfig())
	if attr := plan.Attributes["reusable_by_team"]; attr == nil || !attr.RequiresNew {
		t.Errorf("expected turning sharing off to force a replacement, got:\n%s", describeDiff(plan))
	}
}

// TestProviderCredentialImportUpdatesInPlace covers adopting a credential
// registered by hand. The API returns neither the client ID nor the secret, so
// the first apply sends both, and it must not replace the credential.
func TestProviderCredentialImportUpdatesInPlace(t *testing.T) {
	t.Parallel()

	api := newFakeProviderCredentialAPI()
	id := api.add(fakeProviderCredential{
		provider:     "github",
		displayName:  "Production OAuth app",
		kind:         "PROVIDER_CREDENTIAL_KIND_OAUTH_CLIENT",
		clientID:     "Iv1.0123456789abcdef",
		clientSecret: "secret-registered-by-hand",
		redirectURL:  "https://app.example.com/oauth/callback",
	})
	api.connections[id] = 3
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayProviderCredential()

	state := importProviderCredential(t, r, meta, id)
	d := r.Data(state)
	for key, want := range map[string]string{
		"tool_provider": "github",
		"kind":          "oauth_client",
		"display_name":  "Production OAuth app",
		"redirect_url":  "https://app.example.com/oauth/callback",
		"client_id":     "",
		"client_secret": "",
	} {
		if got := d.Get(key); got != want {
			t.Errorf("%s = %v, expected %q after import", key, got, want)
		}
	}

	config := withValues(oauthClientCredentialConfig(), map[string]interface{}{"redirect_url": "https://app.example.com/oauth/callback"})
	next := mustApplyProviderCredential(t, r, meta, state, config)

	want := []map[string]interface{}{{"client_id": "Iv1.0123456789abcdef", "client_secret": "secret-one"}}
	if got := api.patches; !reflect.DeepEqual(got, want) {
		t.Errorf("update requests = %v, expected %v", got, want)
	}
	if len(api.deletes) != 0 {
		t.Errorf("deletes = %v, expected an imported credential to be kept", api.deletes)
	}
	if next.ID != id {
		t.Errorf("id = %q, expected %q", next.ID, id)
	}
	requireEmptyPlan(t, r, meta, next, config)
}

// TestProviderCredentialImportRejectsAnotherClientID covers importing a
// credential into a configuration for a different OAuth client: the API
// refuses the client ID, nothing is deleted, and the failed values are not
// recorded as applied.
func TestProviderCredentialImportRejectsAnotherClientID(t *testing.T) {
	t.Parallel()

	api := newFakeProviderCredentialAPI()
	id := api.add(fakeProviderCredential{
		provider:     "github",
		displayName:  "Production OAuth app",
		kind:         "PROVIDER_CREDENTIAL_KIND_OAUTH_CLIENT",
		clientID:     "Iv1.registered",
		clientSecret: "secret-registered-by-hand",
	})
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayProviderCredential()

	state := importProviderCredential(t, r, meta, id)
	next, diags := applyProviderCredential(t, r, meta, state, oauthClientCredentialConfig())
	if !diags.HasError() {
		t.Fatal("expected the API to refuse a client ID other than the registered one")
	}
	if !strings.Contains(diags[0].Summary, "client_id does not match") {
		t.Errorf("error = %q, expected the API's reason", diags[0].Summary)
	}
	if len(api.deletes) != 0 {
		t.Errorf("deletes = %v, expected nothing to be deleted", api.deletes)
	}
	if got := api.get(id).clientSecret; got != "secret-registered-by-hand" {
		t.Errorf("stored secret = %q, expected it to be untouched", got)
	}

	d := r.Data(next)
	if got := d.Get("client_id"); got != "" {
		t.Errorf("client_id = %v, expected the refused ID not to be recorded", got)
	}
	plan := planProviderCredential(t, r, meta, next, oauthClientCredentialConfig())
	if plan == nil || plan.Empty() {
		t.Fatal("expected the refused change to be planned again")
	}
	if plan.RequiresNew() {
		t.Errorf("expected the retry to stay in place, got a replacement:\n%s", describeDiff(plan))
	}
}

func TestProviderCredentialFailedRotationIsPlannedAgain(t *testing.T) {
	t.Parallel()

	api := newFakeProviderCredentialAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayProviderCredential()
	state := mustApplyProviderCredential(t, r, meta, nil, apiKeyCredentialConfig())

	failing := http.NewServeMux()
	failing.HandleFunc("PATCH /v2/action-gateway/provider-credentials/{id}", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, http.StatusForbidden, map[string]string{"id": "forbidden", "message": "You are not authorized to perform this operation"})
	})
	rotated := withValues(apiKeyCredentialConfig(), map[string]interface{}{"api_key": "lin_api_two"})
	next, diags := applyProviderCredential(t, r, newTestMeta(t, failing), state, rotated)
	if !diags.HasError() {
		t.Fatal("expected the rotation to fail")
	}

	if got := r.Data(next).Get("api_key"); got != "lin_api_one" {
		t.Errorf("api_key in state = %v, expected the key that is still in use", got)
	}
	if plan := planProviderCredential(t, r, meta, next, rotated); plan == nil || plan.Attributes["api_key"] == nil {
		t.Error("expected the failed rotation to be planned again")
	}
}

// TestProviderCredentialRefreshKeepsRegistrationArguments covers an API that
// reports an immutable value in another spelling. A refresh must not turn that
// into a replacement.
func TestProviderCredentialRefreshKeepsRegistrationArguments(t *testing.T) {
	t.Parallel()

	api := newFakeProviderCredentialAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayProviderCredential()
	config := withValues(oauthClientCredentialConfig(), map[string]interface{}{"redirect_url": "https://app.example.com/oauth/callback"})
	state := mustApplyProviderCredential(t, r, meta, nil, config)

	api.mu.Lock()
	api.report = map[string]interface{}{
		"provider":     "GitHub",
		"kind":         "PROVIDER_CREDENTIAL_KIND_UNSPECIFIED",
		"redirect_url": "https://app.example.com/oauth/callback/",
	}
	api.mu.Unlock()

	requireEmptyPlan(t, r, meta, state, config)
}

func TestProviderCredentialBaseURL(t *testing.T) {
	t.Parallel()

	api := newFakeProviderCredentialAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayProviderCredential()

	// The API stores the origin normalized, which must not show as a change.
	config := withValues(apiKeyCredentialConfig(), map[string]interface{}{"base_url": "https://API.example.com/"})
	state := mustApplyProviderCredential(t, r, meta, nil, config)
	if got := api.creates[0]["base_url"]; got != "https://API.example.com/" {
		t.Errorf("requested base_url = %v", got)
	}
	if got := r.Data(state).Get("base_url"); got != "https://api.example.com" {
		t.Errorf("base_url = %v, expected the stored origin", got)
	}
	requireEmptyPlan(t, r, meta, state, config)

	moved := withValues(apiKeyCredentialConfig(), map[string]interface{}{"base_url": "https://eu.example.com"})
	next := mustApplyProviderCredential(t, r, meta, state, moved)
	if got, want := api.patches, []map[string]interface{}{{"base_url": "https://eu.example.com"}}; !reflect.DeepEqual(got, want) {
		t.Errorf("update requests = %v, expected %v", got, want)
	}
	requireEmptyPlan(t, r, meta, next, moved)
}

func TestProviderCredentialOAuthBaseURLChangeCarriesClientID(t *testing.T) {
	t.Parallel()

	api := newFakeProviderCredentialAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayProviderCredential()

	config := withValues(oauthClientCredentialConfig(), map[string]interface{}{"base_url": "https://github.example.com"})
	state := mustApplyProviderCredential(t, r, meta, nil, config)
	moved := withValues(oauthClientCredentialConfig(), map[string]interface{}{"base_url": "https://git.example.com"})
	mustApplyProviderCredential(t, r, meta, state, moved)

	want := []map[string]interface{}{{"client_id": "Iv1.0123456789abcdef", "base_url": "https://git.example.com"}}
	if got := api.patches; !reflect.DeepEqual(got, want) {
		t.Errorf("update requests = %v, expected %v", got, want)
	}
}

func TestProviderCredentialParametersDeriveBaseURL(t *testing.T) {
	t.Parallel()

	api := newFakeProviderCredentialAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayProviderCredential()

	config := withValues(apiKeyCredentialConfig(), map[string]interface{}{
		"credential_parameters": map[string]interface{}{"subdomain": "acme"},
	})
	state := mustApplyProviderCredential(t, r, meta, nil, config)
	if _, ok := api.creates[0]["base_url"]; ok {
		t.Errorf("create request = %v, expected no base_url alongside credential_parameters", api.creates[0])
	}
	if got := r.Data(state).Get("base_url"); got != "https://acme.example.com" {
		t.Errorf("base_url = %v, expected the origin the API derived", got)
	}
	requireEmptyPlan(t, r, meta, state, config)

	moved := withValues(apiKeyCredentialConfig(), map[string]interface{}{
		"credential_parameters": map[string]interface{}{"subdomain": "acme-eu"},
	})
	diff := planProviderCredential(t, r, meta, state, moved)
	diff.RawConfig = providerCredentialRawConfig(r, moved)
	next, diags := r.Apply(context.Background(), state, diff, meta)
	requireNoErrors(t, diags)

	want := []map[string]interface{}{{"credential_parameters": map[string]interface{}{"subdomain": "acme-eu"}}}
	if got := api.patches; !reflect.DeepEqual(got, want) {
		t.Errorf("update requests = %v, expected %v", got, want)
	}
	if got := r.Data(next).Get("base_url"); got != "https://acme-eu.example.com" {
		t.Errorf("base_url = %v, expected the newly derived origin", got)
	}
	requireEmptyPlan(t, r, meta, next, moved)
}

func TestProviderCredentialReadRemovesMissing(t *testing.T) {
	t.Parallel()

	api := newFakeProviderCredentialAPI()
	meta := api.start(t)

	d := schema.TestResourceDataRaw(t, ResourceDigitalOceanActionGatewayProviderCredential().Schema, apiKeyCredentialConfig())
	d.SetId("cred-gone")
	requireNoErrors(t, resourceDigitalOceanActionGatewayProviderCredentialRead(context.Background(), d, meta))
	if d.Id() != "" {
		t.Errorf("id = %q, expected a missing credential to leave state", d.Id())
	}
}

func TestProviderCredentialReadKeepsStateOnEmptyResponse(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v2/action-gateway/provider-credentials/cred-1", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, http.StatusOK, map[string]interface{}{})
	})
	meta := newTestMeta(t, mux)

	d := schema.TestResourceDataRaw(t, ResourceDigitalOceanActionGatewayProviderCredential().Schema, apiKeyCredentialConfig())
	d.SetId("cred-1")
	diags := resourceDigitalOceanActionGatewayProviderCredentialRead(context.Background(), d, meta)
	if !diags.HasError() {
		t.Fatal("expected an error for a response without a credential")
	}
	if d.Id() != "cred-1" {
		t.Errorf("id = %q, expected the credential to stay in state", d.Id())
	}
}

func TestProviderCredentialImportMissing(t *testing.T) {
	t.Parallel()

	api := newFakeProviderCredentialAPI()
	meta := api.start(t)

	r := ResourceDigitalOceanActionGatewayProviderCredential()
	if _, err := r.Importer.StateContext(context.Background(), r.Data(&terraform.InstanceState{ID: "cred-gone"}), meta); err == nil {
		t.Fatal("expected an error importing a credential that does not exist")
	}
}

func TestProviderCredentialDelete(t *testing.T) {
	t.Parallel()

	api := newFakeProviderCredentialAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayProviderCredential()

	unused := mustApplyProviderCredential(t, r, meta, nil, apiKeyCredentialConfig())
	d := r.Data(unused)
	diags := resourceDigitalOceanActionGatewayProviderCredentialDelete(context.Background(), d, meta)
	if len(diags) != 0 {
		t.Errorf("diagnostics = %v, expected none for a credential without connections", diags)
	}
	if d.Id() != "" {
		t.Errorf("id = %q, expected it to be cleared", d.Id())
	}

	inUse := mustApplyProviderCredential(t, r, meta, nil, apiKeyCredentialConfig())
	api.mu.Lock()
	api.connections[inUse.ID] = 2
	api.mu.Unlock()
	d = r.Data(inUse)
	diags = resourceDigitalOceanActionGatewayProviderCredentialDelete(context.Background(), d, meta)
	requireNoErrors(t, diags)
	if len(diags) != 1 || diags[0].Severity != diag.Warning || !strings.Contains(diags[0].Detail, "2 connection(s)") {
		t.Errorf("diagnostics = %v, expected a warning naming the 2 deleted connections", diags)
	}
	if d.Id() != "" {
		t.Errorf("id = %q, expected it to be cleared", d.Id())
	}

	// Deleting a credential that is already gone succeeds.
	d = r.Data(inUse)
	requireNoErrors(t, resourceDigitalOceanActionGatewayProviderCredentialDelete(context.Background(), d, meta))
	if d.Id() != "" {
		t.Errorf("id = %q, expected it to be cleared", d.Id())
	}
}

func TestValidateProviderCredentialConfig(t *testing.T) {
	t.Parallel()

	r := ResourceDigitalOceanActionGatewayProviderCredential()
	tests := []struct {
		name    string
		attrs   map[string]cty.Value
		wantErr string
	}{
		{
			name:  "API key",
			attrs: map[string]cty.Value{"kind": cty.StringVal("api_key"), "api_key": cty.StringVal("key")},
		},
		{
			name: "API key with parameters",
			attrs: map[string]cty.Value{
				"kind":                  cty.StringVal("api_key"),
				"api_key":               cty.StringVal("key"),
				"credential_parameters": cty.MapVal(map[string]cty.Value{"subdomain": cty.StringVal("acme")}),
			},
		},
		{
			name:    "API key without a key",
			attrs:   map[string]cty.Value{"kind": cty.StringVal("api_key")},
			wantErr: `kind "api_key" requires api_key`,
		},
		{
			name:    "API key with a client ID",
			attrs:   map[string]cty.Value{"kind": cty.StringVal("api_key"), "api_key": cty.StringVal("key"), "client_id": cty.StringVal("abc")},
			wantErr: `client_id cannot be set for kind "api_key"`,
		},
		{
			name:    "API key with a redirect URL",
			attrs:   map[string]cty.Value{"kind": cty.StringVal("api_key"), "api_key": cty.StringVal("key"), "redirect_url": cty.StringVal("https://app.example.com/cb")},
			wantErr: `redirect_url cannot be set for kind "api_key"`,
		},
		{
			name:  "OAuth client",
			attrs: map[string]cty.Value{"kind": cty.StringVal("oauth_client"), "client_id": cty.StringVal("abc"), "client_secret": cty.StringVal("def")},
		},
		{
			name:    "OAuth client without a secret",
			attrs:   map[string]cty.Value{"kind": cty.StringVal("oauth_client"), "client_id": cty.StringVal("abc")},
			wantErr: `kind "oauth_client" requires client_secret`,
		},
		{
			name:    "OAuth client without a client ID",
			attrs:   map[string]cty.Value{"kind": cty.StringVal("oauth_client"), "client_secret": cty.StringVal("def")},
			wantErr: `kind "oauth_client" requires client_id`,
		},
		{
			name: "OAuth client with an API key",
			attrs: map[string]cty.Value{
				"kind": cty.StringVal("oauth_client"), "client_id": cty.StringVal("abc"), "client_secret": cty.StringVal("def"), "api_key": cty.StringVal("key"),
			},
			wantErr: `api_key cannot be set for kind "oauth_client"`,
		},
		{
			name: "OAuth client with parameters",
			attrs: map[string]cty.Value{
				"kind": cty.StringVal("oauth_client"), "client_id": cty.StringVal("abc"), "client_secret": cty.StringVal("def"),
				"credential_parameters": cty.MapVal(map[string]cty.Value{"subdomain": cty.StringVal("acme")}),
			},
			wantErr: `credential_parameters cannot be set for kind "oauth_client"`,
		},
		{
			name: "OAuth client with empty parameters",
			attrs: map[string]cty.Value{
				"kind": cty.StringVal("oauth_client"), "client_id": cty.StringVal("abc"), "client_secret": cty.StringVal("def"),
				"credential_parameters": cty.MapValEmpty(cty.String),
			},
		},
		{
			name:  "secret unknown until apply",
			attrs: map[string]cty.Value{"kind": cty.StringVal("oauth_client"), "client_id": cty.StringVal("abc"), "client_secret": cty.UnknownVal(cty.String)},
		},
		{
			name:  "kind unknown until apply",
			attrs: map[string]cty.Value{"kind": cty.UnknownVal(cty.String)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateProviderCredentialConfig(testRawConfig(r, tt.attrs))
			switch {
			case tt.wantErr == "" && err != nil:
				t.Errorf("unexpected error: %s", err)
			case tt.wantErr != "" && (err == nil || err.Error() != tt.wantErr):
				t.Errorf("error = %v, expected %q", err, tt.wantErr)
			}
		})
	}

	if err := validateProviderCredentialConfig(cty.NilVal); err != nil {
		t.Errorf("unexpected error without a configuration: %s", err)
	}
}

func TestNormalizeOrigin(t *testing.T) {
	t.Parallel()

	same := [][2]string{
		{"https://api.example.com", "https://api.example.com/"},
		{"https://API.Example.com", "https://api.example.com"},
		{"https://api.example.com:443", "https://api.example.com"},
		{"HTTPS://api.example.com", "https://api.example.com"},
	}
	for _, pair := range same {
		if !suppressSameOrigin("base_url", pair[0], pair[1], nil) {
			t.Errorf("expected %q and %q to be the same origin", pair[0], pair[1])
		}
	}

	different := [][2]string{
		{"https://api.example.com", "https://eu.example.com"},
		{"https://api.example.com", "https://api.example.com:8443"},
		{"https://api.example.com", "https://api.example.com/v2"},
		{"https://api.example.com", ""},
		{"", "https://api.example.com"},
	}
	for _, pair := range different {
		if suppressSameOrigin("base_url", pair[0], pair[1], nil) {
			t.Errorf("expected %q and %q to differ", pair[0], pair[1])
		}
	}
}

func TestProviderCredentialDataSource(t *testing.T) {
	t.Parallel()

	api := newFakeProviderCredentialAPI()
	id := api.add(fakeProviderCredential{
		provider:    "linear",
		displayName: "Production key",
		kind:        "PROVIDER_CREDENTIAL_KIND_API_KEY",
		apiKey:      "lin_api_one",
		baseURL:     "https://api.linear.app",
		reusable:    true,
	})
	api.add(fakeProviderCredential{provider: "linear", displayName: "Staging key", kind: "PROVIDER_CREDENTIAL_KIND_API_KEY"})
	api.add(fakeProviderCredential{provider: "github", displayName: "Production key", kind: "PROVIDER_CREDENTIAL_KIND_OAUTH_CLIENT"})
	meta := api.start(t)

	lookups := map[string]map[string]interface{}{
		"by ID":                         {"id": id},
		"by display name":               {"tool_provider": "linear", "display_name": "Production key"},
		"by display name ignoring case": {"tool_provider": "linear", "display_name": "production KEY"},
	}
	for name, lookup := range lookups {
		t.Run(name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, DataSourceDigitalOceanActionGatewayProviderCredential().Schema, lookup)
			requireNoErrors(t, dataSourceDigitalOceanActionGatewayProviderCredentialRead(context.Background(), d, meta))

			if d.Id() != id {
				t.Errorf("id = %q, expected %q", d.Id(), id)
			}
			checks := map[string]interface{}{
				"tool_provider":    "linear",
				"display_name":     "Production key",
				"kind":             "api_key",
				"status":           "active",
				"scope":            "team",
				"reusable_by_team": true,
				"base_url":         "https://api.linear.app",
			}
			for key, want := range checks {
				if got := d.Get(key); got != want {
					t.Errorf("%s = %v, expected %v", key, got, want)
				}
			}
		})
	}

	t.Run("no match", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, DataSourceDigitalOceanActionGatewayProviderCredential().Schema, map[string]interface{}{
			"tool_provider": "jira",
			"display_name":  "Production key",
		})
		diags := dataSourceDigitalOceanActionGatewayProviderCredentialRead(context.Background(), d, meta)
		if !diags.HasError() || !strings.Contains(diags[0].Summary, "no Action Gateway provider credential") {
			t.Errorf("diagnostics = %v, expected a not-found error", diags)
		}
	})
}

func TestProviderCredentialDataSourceHidesSecrets(t *testing.T) {
	t.Parallel()

	s := DataSourceDigitalOceanActionGatewayProviderCredential().Schema
	for _, key := range providerCredentialWriteOnlyArguments {
		if _, ok := s[key]; ok {
			t.Errorf("the data source publishes %s, which the API never returns", key)
		}
	}
	if err := schema.InternalMap(s).InternalValidate(nil); err != nil {
		t.Errorf("invalid data source schema: %s", err)
	}
	if err := ResourceDigitalOceanActionGatewayProviderCredential().InternalValidate(nil, true); err != nil {
		t.Errorf("invalid resource schema: %s", err)
	}
}
