package actiongateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testConnectionID = "39f9f768-a7f9-4c91-9557-3328c53dd054"

func newTestMeta(t *testing.T, mux *http.ServeMux) *config.CombinedConfig {
	t.Helper()

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	cfg := &config.Config{
		Token:            "test-token",
		APIEndpoint:      server.URL,
		HTTPRetryMax:     1,
		HTTPRetryWaitMin: 0.01,
		HTTPRetryWaitMax: 0.01,
	}
	meta, err := cfg.Client()
	if err != nil {
		t.Fatalf("error building client: %s", err)
	}
	return meta
}

func writeJSON(t *testing.T, w http.ResponseWriter, status int, body interface{}) {
	t.Helper()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		t.Errorf("error encoding response: %s", err)
	}
}

func decodeJSON(t *testing.T, r *http.Request) map[string]interface{} {
	t.Helper()
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		t.Errorf("error decoding request: %s", err)
	}
	return body
}

func writeNotFound(t *testing.T, w http.ResponseWriter) {
	writeJSON(t, w, http.StatusNotFound, map[string]string{"id": "not_found", "message": "The resource you requested could not be found."})
}

func requireNoErrors(t *testing.T, diags diag.Diagnostics) {
	t.Helper()
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
}

func pendingJiraConnection(status string) map[string]interface{} {
	return map[string]interface{}{
		"id":                    testConnectionID,
		"provider":              "jira",
		"provider_display_name": "Jira",
		"user_id":               "alice",
		"status":                status,
		"credential_kind":       "digitalocean_oauth",
		"credential_id":         "",
		"owning_user_id":        "2f8e1c3a-6b4d-4f7e-9a1c-5d3b7e9f1a2c",
		// The API normalizes the site URL and grants a broader scope set than
		// requested when an existing connection already covered the request.
		"connection_parameters": map[string]interface{}{"site_url": "https://example.atlassian.net/"},
		"oauth":                 map[string]interface{}{"scopes": []string{"read:jira-work", "write:jira-work"}, "granted_at": nil},
		"scopes":                []string{"read:jira-work", "write:jira-work"},
		"created_at":            "2026-07-14T00:50:00Z",
		"updated_at":            "2026-07-14T00:50:00Z",
	}
}

// serveConnectionList serves conns as the only page of the connection list.
func serveConnectionList(t *testing.T, mux *http.ServeMux, conns ...map[string]interface{}) {
	if conns == nil {
		conns = []map[string]interface{}{}
	}
	mux.HandleFunc("GET /v2/action-gateway/connections", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, http.StatusOK, map[string]interface{}{
			"connections": conns,
			"pagination":  map[string]int{"page": 1, "per_page": 100, "total": len(conns)},
		})
	})
}

func pendingAuthorization() map[string]interface{} {
	return map[string]interface{}{
		"status":            "requires_authorization",
		"connect_url":       "https://cloud.digitalocean.com/security/connectlinks/confirm?token=abc",
		"verification_code": "97rwxzcy",
		"expires_at":        "2026-07-14T01:00:00Z",
	}
}

func TestConnectionCreatePendingOAuth(t *testing.T) {
	t.Parallel()

	var createBody map[string]interface{}
	mux := http.NewServeMux()
	serveConnectionList(t, mux)
	mux.HandleFunc("POST /v2/action-gateway/connections", func(w http.ResponseWriter, r *http.Request) {
		createBody = decodeJSON(t, r)
		writeJSON(t, w, http.StatusOK, map[string]interface{}{
			"connection":    pendingJiraConnection("pending"),
			"authorization": pendingAuthorization(),
		})
	})
	mux.HandleFunc("GET /v2/action-gateway/connections/"+testConnectionID, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, http.StatusOK, map[string]interface{}{
			"connection":    pendingJiraConnection("pending"),
			"authorization": pendingAuthorization(),
		})
	})
	meta := newTestMeta(t, mux)

	d := schema.TestResourceDataRaw(t, ResourceDigitalOceanActionGatewayConnection().Schema, map[string]interface{}{
		"tool_provider":         "jira",
		"user_id":               "alice",
		"scopes":                []interface{}{"write:jira-work"},
		"connection_parameters": map[string]interface{}{"site_url": "https://example.atlassian.net"},
	})
	requireNoErrors(t, resourceDigitalOceanActionGatewayConnectionCreate(context.Background(), d, meta))

	if d.Id() != testConnectionID {
		t.Errorf("id = %q, expected %q", d.Id(), testConnectionID)
	}
	if _, ok := createBody["credential"]; ok {
		t.Errorf("expected credential to be omitted so the shared OAuth application is used, got %v", createBody["credential"])
	}
	if got := createBody["scopes"]; !reflect.DeepEqual(got, []interface{}{"write:jira-work"}) {
		t.Errorf("requested scopes = %v", got)
	}

	checks := map[string]string{
		"status":                         "pending",
		"tool_provider_display_name":     "Jira",
		"credential_kind":                "digitalocean_oauth",
		"authorization.0.status":         "requires_authorization",
		"authorization.0.connect_url":    "https://cloud.digitalocean.com/security/connectlinks/confirm?token=abc",
		"authorization.0.expires_at":     "2026-07-14T01:00:00Z",
		"connection_parameters.site_url": "https://example.atlassian.net",
		"created_at":                     "2026-07-14T00:50:00Z",
	}
	for key, want := range checks {
		if got := d.Get(key); got != want {
			t.Errorf("%s = %v, expected %q", key, got, want)
		}
	}
	if got := expandStringSet(d.Get("scopes")); !reflect.DeepEqual(got, []string{"write:jira-work"}) {
		t.Errorf("scopes = %v, expected the configured scopes to be kept", got)
	}
}

func describeDiff(diff *terraform.InstanceDiff) string {
	var lines []string
	for k, v := range diff.Attributes {
		line := k + ": " + strconv.Quote(v.Old) + " -> " + strconv.Quote(v.New)
		if v.NewComputed {
			line += " (computed)"
		}
		if v.RequiresNew {
			line += " (forces replacement)"
		}
		lines = append(lines, line)
	}
	sort.Strings(lines)
	return strings.Join(lines, "\n")
}

// testRawConfig builds the configuration value Terraform sends the provider,
// which is how the provider tells an empty argument apart from an omitted one.
func testRawConfig(r *schema.Resource, attrs map[string]cty.Value) cty.Value {
	values := map[string]cty.Value{}
	for name, ty := range r.CoreConfigSchema().ImpliedType().AttributeTypes() {
		values[name] = cty.NullVal(ty)
	}
	for name, v := range attrs {
		values[name] = v
	}
	return cty.ObjectVal(values)
}

func TestConnectionPlanAfterApplyIsEmpty(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	serveConnectionList(t, mux)
	mux.HandleFunc("POST /v2/action-gateway/connections", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, http.StatusOK, map[string]interface{}{"connection": pendingJiraConnection("pending"), "authorization": pendingAuthorization()})
	})
	mux.HandleFunc("GET /v2/action-gateway/connections/"+testConnectionID, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, http.StatusOK, map[string]interface{}{"connection": pendingJiraConnection("pending"), "authorization": pendingAuthorization()})
	})
	meta := newTestMeta(t, mux)

	tests := []struct {
		name       string
		config     map[string]interface{}
		raw        map[string]cty.Value
		wantScopes []string
	}{
		{
			name:       "omitted scopes",
			config:     map[string]interface{}{"tool_provider": "jira", "user_id": "alice"},
			raw:        map[string]cty.Value{"tool_provider": cty.StringVal("jira"), "user_id": cty.StringVal("alice")},
			wantScopes: []string{"read:jira-work", "write:jira-work"},
		},
		{
			name:       "empty scopes",
			config:     map[string]interface{}{"tool_provider": "jira", "user_id": "alice", "scopes": []interface{}{}},
			raw:        map[string]cty.Value{"tool_provider": cty.StringVal("jira"), "user_id": cty.StringVal("alice"), "scopes": cty.SetValEmpty(cty.String)},
			wantScopes: []string{},
		},
		{
			name:       "narrower scopes",
			config:     map[string]interface{}{"tool_provider": "jira", "user_id": "alice", "scopes": []interface{}{"write:jira-work"}},
			raw:        map[string]cty.Value{"tool_provider": cty.StringVal("jira"), "user_id": cty.StringVal("alice"), "scopes": cty.SetVal([]cty.Value{cty.StringVal("write:jira-work")})},
			wantScopes: []string{"write:jira-work"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			r := ResourceDigitalOceanActionGatewayConnection()
			cfg := terraform.NewResourceConfigRaw(tt.config)

			diff, err := r.Diff(ctx, nil, cfg, meta)
			if err != nil {
				t.Fatalf("error planning create: %s", err)
			}
			diff.RawConfig = testRawConfig(r, tt.raw)
			state, diags := r.Apply(ctx, nil, diff, meta)
			requireNoErrors(t, diags)

			state, diags = r.RefreshWithoutUpgrade(ctx, state, meta)
			requireNoErrors(t, diags)
			if got := expandStringSet(r.Data(state).Get("scopes")); !reflect.DeepEqual(got, tt.wantScopes) {
				t.Errorf("scopes = %v, expected %v", got, tt.wantScopes)
			}

			next, err := r.Diff(ctx, state, cfg, meta)
			if err != nil {
				t.Fatalf("error planning after apply: %s", err)
			}
			if next != nil && !next.Empty() {
				t.Errorf("expected an empty plan after apply, got:\n%s", describeDiff(next))
			}
		})
	}
}

func TestConnectionImportFillsRequestArguments(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v2/action-gateway/connections/"+testConnectionID, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, http.StatusOK, map[string]interface{}{"connection": pendingJiraConnection("active")})
	})
	meta := newTestMeta(t, mux)

	ctx := context.Background()
	r := ResourceDigitalOceanActionGatewayConnection()
	imported, err := r.Importer.StateContext(ctx, r.Data(&terraform.InstanceState{ID: testConnectionID}), meta)
	if err != nil {
		t.Fatalf("error importing: %s", err)
	}
	state, diags := r.RefreshWithoutUpgrade(ctx, imported[0].State(), meta)
	requireNoErrors(t, diags)

	d := r.Data(state)
	if got := expandStringSet(d.Get("scopes")); !reflect.DeepEqual(got, []string{"read:jira-work", "write:jira-work"}) {
		t.Errorf("scopes = %v, expected the granted scopes", got)
	}
	if got := d.Get("connection_parameters.site_url"); got != "https://example.atlassian.net/" {
		t.Errorf("connection_parameters.site_url = %v", got)
	}

	cfg := terraform.NewResourceConfigRaw(map[string]interface{}{
		"tool_provider":         "jira",
		"user_id":               "alice",
		"connection_parameters": map[string]interface{}{"site_url": "https://example.atlassian.net/"},
	})
	next, err := r.Diff(ctx, state, cfg, meta)
	if err != nil {
		t.Fatalf("error planning after import: %s", err)
	}
	if next != nil && !next.Empty() {
		t.Errorf("expected an empty plan after import, got:\n%s", describeDiff(next))
	}
}

func TestConnectionImportMissing(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v2/action-gateway/connections/"+testConnectionID, func(w http.ResponseWriter, r *http.Request) {
		writeNotFound(t, w)
	})
	meta := newTestMeta(t, mux)

	r := ResourceDigitalOceanActionGatewayConnection()
	if _, err := r.Importer.StateContext(context.Background(), r.Data(&terraform.InstanceState{ID: testConnectionID}), meta); err == nil {
		t.Fatal("expected an error importing a connection that does not exist")
	}
}

func TestConnectionCreateTeamCredentialOverVPC(t *testing.T) {
	t.Parallel()

	var createBody map[string]interface{}
	connection := map[string]interface{}{
		"id":              testConnectionID,
		"provider":        "jira",
		"user_id":         "svc-bot",
		"status":          "active",
		"credential_kind": "team_api_key",
		"credential_id":   "8d6f2b1a-3c4e-4d5f-9a8b-7c6d5e4f3a2b",
		"scopes":          []string{},
		"network": map[string]interface{}{"vpc": map[string]interface{}{
			"vpc_uuid": "5a4981aa-9653-4bd1-bef5-d6bff52042e4",
			"region":   "ric1",
			"destinations": []interface{}{map[string]interface{}{
				"host":             "jira.internal.example.com",
				"port":             443,
				"allowed_ip_cidrs": []string{"10.116.0.0/20"},
			}},
		}},
	}
	mux := http.NewServeMux()
	serveConnectionList(t, mux)
	mux.HandleFunc("POST /v2/action-gateway/connections", func(w http.ResponseWriter, r *http.Request) {
		createBody = decodeJSON(t, r)
		writeJSON(t, w, http.StatusOK, map[string]interface{}{"connection": connection})
	})
	mux.HandleFunc("GET /v2/action-gateway/connections/"+testConnectionID, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, http.StatusOK, map[string]interface{}{"connection": connection})
	})
	meta := newTestMeta(t, mux)

	d := schema.TestResourceDataRaw(t, ResourceDigitalOceanActionGatewayConnection().Schema, map[string]interface{}{
		"tool_provider": "jira",
		"user_id":       "svc-bot",
		"credential_id": "8d6f2b1a-3c4e-4d5f-9a8b-7c6d5e4f3a2b",
		"vpc": []interface{}{map[string]interface{}{
			"vpc_uuid": "5a4981aa-9653-4bd1-bef5-d6bff52042e4",
			"destination": []interface{}{map[string]interface{}{
				"host":             "jira.internal.example.com",
				"port":             443,
				"allowed_ip_cidrs": []interface{}{"10.116.0.0/20"},
			}},
		}},
	})
	requireNoErrors(t, resourceDigitalOceanActionGatewayConnectionCreate(context.Background(), d, meta))

	wantCredential := map[string]interface{}{"team_credential": map[string]interface{}{"credential_id": "8d6f2b1a-3c4e-4d5f-9a8b-7c6d5e4f3a2b"}}
	if !reflect.DeepEqual(createBody["credential"], wantCredential) {
		t.Errorf("credential = %v, expected %v", createBody["credential"], wantCredential)
	}
	wantNetwork := map[string]interface{}{"vpc": map[string]interface{}{
		"vpc_uuid": "5a4981aa-9653-4bd1-bef5-d6bff52042e4",
		"destinations": []interface{}{map[string]interface{}{
			"host":             "jira.internal.example.com",
			"port":             float64(443),
			"allowed_ip_cidrs": []interface{}{"10.116.0.0/20"},
		}},
	}}
	if !reflect.DeepEqual(createBody["network"], wantNetwork) {
		t.Errorf("network = %v, expected %v", createBody["network"], wantNetwork)
	}

	if got := d.Get("vpc.0.region"); got != "ric1" {
		t.Errorf("vpc.0.region = %v, expected ric1", got)
	}
	if got := d.Get("status"); got != "active" {
		t.Errorf("status = %v, expected active", got)
	}
	if got := d.Get("authorization").([]interface{}); len(got) != 0 {
		t.Errorf("authorization = %v, expected none for an active connection", got)
	}
}

func TestConnectionReadRemovesRevoked(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v2/action-gateway/connections/"+testConnectionID, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, http.StatusOK, map[string]interface{}{"connection": pendingJiraConnection("revoked")})
	})
	meta := newTestMeta(t, mux)

	d := schema.TestResourceDataRaw(t, ResourceDigitalOceanActionGatewayConnection().Schema, map[string]interface{}{})
	d.SetId(testConnectionID)
	requireNoErrors(t, resourceDigitalOceanActionGatewayConnectionRead(context.Background(), d, meta))

	if d.Id() != "" {
		t.Errorf("id = %q, expected a revoked connection to be removed from state", d.Id())
	}
}

func TestConnectionDeleteAlreadyGone(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /v2/action-gateway/connections/"+testConnectionID, func(w http.ResponseWriter, r *http.Request) {
		writeNotFound(t, w)
	})
	meta := newTestMeta(t, mux)

	d := schema.TestResourceDataRaw(t, ResourceDigitalOceanActionGatewayConnection().Schema, map[string]interface{}{})
	d.SetId(testConnectionID)
	requireNoErrors(t, resourceDigitalOceanActionGatewayConnectionDelete(context.Background(), d, meta))

	if d.Id() != "" {
		t.Errorf("id = %q, expected it to be cleared", d.Id())
	}
}

func TestFindConnectionsMatchesExactUserID(t *testing.T) {
	t.Parallel()

	var pages []string
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v2/action-gateway/connections", func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("provider") != "github" || q.Get("user_id") != "alice" {
			t.Errorf("unexpected filters %v", q)
		}
		pages = append(pages, q.Get("page"))
		var conns []map[string]interface{}
		if q.Get("page") == "1" {
			conns = []map[string]interface{}{
				{"id": "a", "provider": "github", "user_id": "alice-2"},
				{"id": "b", "provider": "github", "user_id": "alice"},
			}
		} else {
			conns = []map[string]interface{}{
				{"id": "c", "provider": "github", "user_id": "malice"},
				{"id": "d", "provider": "github", "user_id": "alice", "status": "revoked"},
			}
		}
		writeJSON(t, w, http.StatusOK, map[string]interface{}{
			"connections": conns,
			"pagination":  map[string]int{"page": 1, "per_page": 100, "total": 101},
		})
	})
	meta := newTestMeta(t, mux)

	found, err := findConnections(context.Background(), meta.GodoClient(), "github", "alice")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if len(found) != 1 || found[0].ID != "b" {
		t.Errorf("found %+v, expected only connection b", found)
	}
	if !reflect.DeepEqual(pages, []string{"1", "2"}) {
		t.Errorf("requested pages %v, expected [1 2]", pages)
	}
}

func TestConnectionCreateRefusesExistingConnection(t *testing.T) {
	t.Parallel()

	otherOwner := pendingJiraConnection("active")
	otherOwner["id"] = "7c1e5a2b-4d3f-4e6a-8b9c-0d1e2f3a4b5c"
	otherOwner["owning_user_id"] = "9a8b7c6d-5e4f-4a3b-2c1d-0e9f8a7b6c5d"

	tests := []struct {
		name     string
		existing []map[string]interface{}
		wantErr  bool
	}{
		{name: "no existing connection"},
		{name: "another user's connection", existing: []map[string]interface{}{otherOwner}},
		{name: "revoked connection restarted", existing: []map[string]interface{}{pendingJiraConnection("revoked")}},
		{name: "active connection returned", existing: []map[string]interface{}{pendingJiraConnection("active")}, wantErr: true},
		{name: "pending connection returned", existing: []map[string]interface{}{pendingJiraConnection("pending")}, wantErr: true},
		{name: "expired connection restarted", existing: []map[string]interface{}{pendingJiraConnection("expired")}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mux := http.NewServeMux()
			serveConnectionList(t, mux, tt.existing...)
			mux.HandleFunc("POST /v2/action-gateway/connections", func(w http.ResponseWriter, r *http.Request) {
				writeJSON(t, w, http.StatusOK, map[string]interface{}{"connection": pendingJiraConnection("pending"), "authorization": pendingAuthorization()})
			})
			mux.HandleFunc("GET /v2/action-gateway/connections/"+testConnectionID, func(w http.ResponseWriter, r *http.Request) {
				writeJSON(t, w, http.StatusOK, map[string]interface{}{"connection": pendingJiraConnection("pending"), "authorization": pendingAuthorization()})
			})
			meta := newTestMeta(t, mux)

			d := schema.TestResourceDataRaw(t, ResourceDigitalOceanActionGatewayConnection().Schema, map[string]interface{}{
				"tool_provider": "jira",
				"user_id":       "alice",
			})
			diags := resourceDigitalOceanActionGatewayConnectionCreate(context.Background(), d, meta)

			if !tt.wantErr {
				requireNoErrors(t, diags)
				if d.Id() != testConnectionID {
					t.Errorf("id = %q, expected %q", d.Id(), testConnectionID)
				}
				return
			}
			if !diags.HasError() || !strings.Contains(diags[0].Summary, "terraform import digitalocean_action_gateway_connection.<name> "+testConnectionID) {
				t.Fatalf("diagnostics = %v, expected an error suggesting an import", diags)
			}
			if d.Id() != "" {
				t.Errorf("id = %q, expected Terraform not to take ownership of the existing connection", d.Id())
			}
		})
	}
}

func TestConnectionDataSourceSkipsRevoked(t *testing.T) {
	t.Parallel()

	const revokedID = "0b5d2c4e-8f1a-4c3b-9d7e-6a5f4e3d2c1b"
	revoked := pendingJiraConnection("revoked")
	revoked["id"] = revokedID

	mux := http.NewServeMux()
	serveConnectionList(t, mux, revoked, pendingJiraConnection("active"))
	mux.HandleFunc("GET /v2/action-gateway/connections/"+testConnectionID, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, http.StatusOK, map[string]interface{}{"connection": pendingJiraConnection("active")})
	})
	mux.HandleFunc("GET /v2/action-gateway/connections/"+revokedID, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, http.StatusOK, map[string]interface{}{"connection": revoked})
	})
	meta := newTestMeta(t, mux)
	r := DataSourceDigitalOceanActionGatewayConnection()

	t.Run("by tool provider and user ID", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{"tool_provider": "jira", "user_id": "alice"})
		requireNoErrors(t, dataSourceDigitalOceanActionGatewayConnectionRead(context.Background(), d, meta))
		if d.Id() != testConnectionID {
			t.Errorf("id = %q, expected the connection that is not revoked", d.Id())
		}
	})
	t.Run("by ID", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{"id": revokedID})
		diags := dataSourceDigitalOceanActionGatewayConnectionRead(context.Background(), d, meta)
		if !diags.HasError() || !strings.Contains(diags[0].Summary, "revoked") {
			t.Errorf("diagnostics = %v, expected an error for a revoked connection", diags)
		}
	})
}

func TestReconcileToolbeltTools(t *testing.T) {
	tests := []struct {
		name       string
		configured []string
		pinned     []string
		want       []string
	}{
		{
			name:       "unpinned configuration keeps its form",
			configured: []string{"exa_web_search"},
			pinned:     []string{"exa_web_search@1"},
			want:       []string{"exa_web_search"},
		},
		{
			name:       "pinned configuration keeps the pin",
			configured: []string{"exa_web_search@1"},
			pinned:     []string{"exa_web_search@1"},
			want:       []string{"exa_web_search@1"},
		},
		{
			name:       "members added outside Terraform are reported pinned",
			configured: []string{"exa_web_search"},
			pinned:     []string{"exa_web_fetch@3", "exa_web_search@1"},
			want:       []string{"exa_web_fetch@3", "exa_web_search"},
		},
		{
			name:   "import reports pinned members",
			pinned: []string{"exa_web_search@1"},
			want:   []string{"exa_web_search@1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reconcileToolbeltTools(tt.configured, tt.pinned); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToolbeltMembershipChanges(t *testing.T) {
	tests := []struct {
		name       string
		pinned     []string
		desired    []string
		wantRemove []string
		wantAdd    []string
	}{
		{
			name:    "add a tool",
			pinned:  []string{"exa_web_search@1"},
			desired: []string{"exa_web_search", "exa_web_fetch"},
			wantAdd: []string{"exa_web_fetch"},
		},
		{
			name:       "remove a tool by slug",
			pinned:     []string{"exa_web_fetch@2", "exa_web_search@1"},
			desired:    []string{"exa_web_search"},
			wantRemove: []string{"exa_web_fetch"},
		},
		{
			name:       "change a pin",
			pinned:     []string{"exa_web_search@1"},
			desired:    []string{"exa_web_search@2"},
			wantRemove: []string{"exa_web_search"},
			wantAdd:    []string{"exa_web_search@2"},
		},
		{
			name:    "pinning to the current version changes nothing",
			pinned:  []string{"exa_web_search@1"},
			desired: []string{"exa_web_search@1"},
		},
		{
			name:    "unpinning changes nothing",
			pinned:  []string{"exa_web_search@1"},
			desired: []string{"exa_web_search"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			remove, add := toolbeltMembershipChanges(tt.pinned, tt.desired)
			if !reflect.DeepEqual(remove, tt.wantRemove) {
				t.Errorf("remove = %v, want %v", remove, tt.wantRemove)
			}
			if !reflect.DeepEqual(add, tt.wantAdd) {
				t.Errorf("add = %v, want %v", add, tt.wantAdd)
			}
		})
	}
}

// fakeToolbeltAPI serves one toolbelt whose membership changes, and version
// increments, on each add or remove request.
type fakeToolbeltAPI struct {
	mu      sync.Mutex
	version int
	tools   map[string]string
	calls   []string
}

func (api *fakeToolbeltAPI) toolbelt() map[string]interface{} {
	var pinned []string
	for slug, version := range api.tools {
		pinned = append(pinned, slug+"@"+version)
	}
	version := strconv.Itoa(api.version)
	return map[string]interface{}{
		"name":             "search-toolbelt",
		"version":          version,
		"display_name":     "Search Tools",
		"tools":            pinned,
		"status":           "active",
		"reference":        "search-toolbelt@" + version,
		"reference_latest": "search-toolbelt",
		"tool_count":       len(pinned),
	}
}

func (api *fakeToolbeltAPI) start(t *testing.T) *config.CombinedConfig {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v2/action-gateway/toolbelts", func(w http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		defer api.mu.Unlock()
		body := decodeJSON(t, r)
		api.version = 1
		for _, ref := range body["tools"].([]interface{}) {
			slug, _ := splitToolPin(ref.(string))
			api.tools[slug] = "1"
		}
		writeJSON(t, w, http.StatusOK, map[string]interface{}{"toolbelt": api.toolbelt()})
	})
	mux.HandleFunc("GET /v2/action-gateway/toolbelts/search-toolbelt", func(w http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		defer api.mu.Unlock()
		writeJSON(t, w, http.StatusOK, map[string]interface{}{"toolbelt": api.toolbelt()})
	})
	mux.HandleFunc("POST /v2/action-gateway/toolbelts/search-toolbelt/tools/{op}", func(w http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		defer api.mu.Unlock()
		body := decodeJSON(t, r)
		var refs []string
		for _, ref := range body["tools"].([]interface{}) {
			refs = append(refs, ref.(string))
			slug, version := splitToolPin(ref.(string))
			if r.PathValue("op") == "remove" {
				delete(api.tools, slug)
			} else {
				if version == "" {
					version = "1"
				}
				api.tools[slug] = version
			}
		}
		api.calls = append(api.calls, r.PathValue("op")+":"+strings.Join(refs, ","))
		api.version++
		writeJSON(t, w, http.StatusOK, map[string]interface{}{"toolbelt": api.toolbelt()})
	})
	return newTestMeta(t, mux)
}

func TestToolbeltUpdateAddsAndRemovesTools(t *testing.T) {
	t.Parallel()

	api := &fakeToolbeltAPI{tools: map[string]string{}}
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayToolbelt()

	d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
		"name":         "search-toolbelt",
		"display_name": "Search Tools",
		"tools":        []interface{}{"exa_web_search", "exa_web_fetch"},
	})
	requireNoErrors(t, resourceDigitalOceanActionGatewayToolbeltCreate(context.Background(), d, meta))

	if got := expandStringSet(d.Get("tools")); !reflect.DeepEqual(got, []string{"exa_web_fetch", "exa_web_search"}) {
		t.Errorf("tools = %v, expected the configured unpinned slugs", got)
	}
	if got := expandStringSet(d.Get("pinned_tools")); !reflect.DeepEqual(got, []string{"exa_web_fetch@1", "exa_web_search@1"}) {
		t.Errorf("pinned_tools = %v", got)
	}
	if got := d.Get("reference"); got != "search-toolbelt@1" {
		t.Errorf("reference = %v, expected search-toolbelt@1", got)
	}

	state := d.State()
	cfg := terraform.NewResourceConfigRaw(map[string]interface{}{
		"name":         "search-toolbelt",
		"display_name": "Search Tools",
		"tools":        []interface{}{"exa_web_search", "tavily_search@1"},
	})
	diff, err := r.Diff(context.Background(), state, cfg, meta)
	if err != nil {
		t.Fatalf("error planning update: %s", err)
	}
	if diff.RequiresNew() {
		t.Fatalf("changing tools should update in place, got %v", diff)
	}
	if attr := diff.Attributes["reference"]; attr == nil || !attr.NewComputed {
		t.Errorf("expected reference to be unknown until apply, got %v", attr)
	}

	newState, diags := r.Apply(context.Background(), state, diff, meta)
	requireNoErrors(t, diags)
	updated := r.Data(newState)

	if want := []string{"remove:exa_web_fetch", "add:tavily_search@1"}; !reflect.DeepEqual(api.calls, want) {
		t.Errorf("calls = %v, expected %v", api.calls, want)
	}
	if got := expandStringSet(updated.Get("tools")); !reflect.DeepEqual(got, []string{"exa_web_search", "tavily_search@1"}) {
		t.Errorf("tools = %v", got)
	}
	if got := updated.Get("reference"); got != "search-toolbelt@3" {
		t.Errorf("reference = %v, expected search-toolbelt@3", got)
	}
}

// fakeMCPServerAPI serves one MCP server. A requested resync completes after
// pendingPolls reads of the server.
type fakeMCPServerAPI struct {
	mu           sync.Mutex
	server       map[string]interface{}
	tools        []map[string]interface{}
	resyncs      int
	resyncUsers  []string
	syncing      bool
	pendingPolls int
	authorize    bool
	created      bool
	enabled      [][]string
}

func (api *fakeMCPServerAPI) start(t *testing.T) *config.CombinedConfig {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v2/action-gateway/mcp-servers", func(w http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		defer api.mu.Unlock()
		body := decodeJSON(t, r)
		if _, ok := body["api_key"]; !ok {
			t.Errorf("expected api_key in request, got %v", body)
		}
		api.created = true
		api.server = map[string]interface{}{
			"serverRef":                       body["serverRef"],
			"endpoint":                        body["endpoint"],
			"transport":                       "streamable_http",
			"credentialRefSource":             body["credentialRefSource"],
			"credentialRef":                   "action-infra/tool-credentials/mcp/docs/3f5b",
			"description":                     body["description"],
			"syncStatus":                      "pending",
			"lastSyncedAt":                    "",
			"oauth_authorization_ttl_seconds": "0",
			"oauth_scopes":                    []string{},
		}
		writeJSON(t, w, http.StatusOK, map[string]interface{}{"mcpServer": api.server})
	})
	mux.HandleFunc("GET /v2/action-gateway/mcp-servers/docs", func(w http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		defer api.mu.Unlock()
		if api.syncing {
			if api.pendingPolls == 0 {
				api.syncing = false
				api.server["lastSyncedAt"] = fmt.Sprintf("2026-09-20T15:%02d:00Z", api.resyncs)
				api.server["syncStatus"] = "synced"
				api.server["toolCount"] = len(api.tools)
			} else {
				api.pendingPolls--
			}
		}
		writeJSON(t, w, http.StatusOK, map[string]interface{}{"mcpServer": api.server})
	})
	mux.HandleFunc("POST /v2/action-gateway/mcp-servers/docs/resync", func(w http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		defer api.mu.Unlock()
		body := decodeJSON(t, r)
		userID, _ := body["user_id"].(string)
		api.resyncs++
		api.resyncUsers = append(api.resyncUsers, userID)
		if api.authorize {
			writeJSON(t, w, http.StatusOK, map[string]interface{}{
				"mcpServer":     api.server,
				"authorization": pendingAuthorization(),
			})
			return
		}
		api.syncing = true
		writeJSON(t, w, http.StatusAccepted, map[string]interface{}{"mcpServer": api.server, "pending": true})
	})
	mux.HandleFunc("GET /v2/action-gateway/mcp-servers/docs/tools", func(w http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		defer api.mu.Unlock()
		tools := api.tools
		if api.server["lastSyncedAt"] == "" {
			tools = nil
		}
		writeJSON(t, w, http.StatusOK, map[string]interface{}{"tools": tools})
	})
	mux.HandleFunc("PUT /v2/action-gateway/mcp-servers/docs/tools", func(w http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		defer api.mu.Unlock()
		body := decodeJSON(t, r)
		var slugs []string
		wanted := map[string]bool{}
		for _, s := range body["enabledToolSlugs"].([]interface{}) {
			slugs = append(slugs, s.(string))
			wanted[s.(string)] = true
		}
		api.enabled = append(api.enabled, slugs)
		for _, tool := range api.tools {
			tool["enabled"] = wanted[tool["toolSlug"].(string)]
		}
		writeJSON(t, w, http.StatusOK, map[string]interface{}{"tools": api.tools})
	})
	return newTestMeta(t, mux)
}

func newFakeMCPServerAPI() *fakeMCPServerAPI {
	return &fakeMCPServerAPI{tools: []map[string]interface{}{
		{"name": "search", "toolSlug": "docs_search", "description": "Search the docs", "enabled": false},
		{"name": "get_page", "toolSlug": "docs_get_page", "description": "Fetch one page", "enabled": false},
	}}
}

func newTestMCPServerData(t *testing.T) *schema.ResourceData {
	return schema.TestResourceDataRaw(t, ResourceDigitalOceanActionGatewayMCPServer().Schema, map[string]interface{}{
		"server_ref":            "docs",
		"endpoint":              "https://mcp.example.com/mcp",
		"credential_ref_source": "secret",
		"api_key":               "sk_live_example",
		"description":           "Internal documentation search.",
		"enabled_tools":         []interface{}{"docs_search"},
	})
}

func TestMCPServerCreateDiscoversAndEnablesTools(t *testing.T) {
	t.Parallel()

	api := newFakeMCPServerAPI()
	api.pendingPolls = 1
	meta := api.start(t)
	d := newTestMCPServerData(t)

	diags := resourceDigitalOceanActionGatewayMCPServerCreate(context.Background(), d, meta)
	requireNoErrors(t, diags)
	if len(diags) != 0 {
		t.Errorf("unexpected diagnostics: %v", diags)
	}

	if want := [][]string{{"docs_search"}}; !reflect.DeepEqual(api.enabled, want) {
		t.Errorf("enabled tool requests = %v, expected %v", api.enabled, want)
	}
	if got := expandStringSet(d.Get("enabled_tools")); !reflect.DeepEqual(got, []string{"docs_search"}) {
		t.Errorf("enabled_tools = %v", got)
	}
	if got := d.Get("discovered_tools.#"); got != 2 {
		t.Errorf("discovered_tools.# = %v, expected 2", got)
	}
	if got := d.Get("last_synced_at"); got != "2026-09-20T15:01:00Z" {
		t.Errorf("last_synced_at = %v", got)
	}
	if got := d.Get("credential_ref"); got != "action-infra/tool-credentials/mcp/docs/3f5b" {
		t.Errorf("credential_ref = %v", got)
	}
	if got := d.Get("api_key"); got != "sk_live_example" {
		t.Errorf("api_key = %v, expected the write-only key to be kept from configuration", got)
	}
}

func TestMCPServerCreateWarnsWhenAuthorizationRequired(t *testing.T) {
	t.Parallel()

	api := newFakeMCPServerAPI()
	api.authorize = true
	meta := api.start(t)
	d := newTestMCPServerData(t)

	diags := resourceDigitalOceanActionGatewayMCPServerCreate(context.Background(), d, meta)
	requireNoErrors(t, diags)
	if len(diags) != 1 || diags[0].Severity != diag.Warning {
		t.Fatalf("diagnostics = %v, expected one warning", diags)
	}
	if strings.Contains(diags[0].Detail, "connectlinks") {
		t.Errorf("warning should not expose the consent link: %s", diags[0].Detail)
	}
	if d.Id() != "docs" {
		t.Errorf("id = %q, expected the server to stay in state", d.Id())
	}
	if len(api.enabled) != 0 {
		t.Errorf("expected no enabled tool requests, got %v", api.enabled)
	}
	if got := d.Get("enabled_tools").(*schema.Set).Len(); got != 0 {
		t.Errorf("enabled_tools has %d items, expected none so the next plan retries", got)
	}
}

func TestMCPServerCreateWarnsWhenDiscoveryOutlastsTimeout(t *testing.T) {
	t.Parallel()

	api := newFakeMCPServerAPI()
	api.pendingPolls = 1 << 30
	meta := api.start(t)

	ctx := context.Background()
	r := ResourceDigitalOceanActionGatewayMCPServer()
	cfg := terraform.NewResourceConfigRaw(map[string]interface{}{
		"server_ref":            "docs",
		"endpoint":              "https://mcp.example.com/mcp",
		"credential_ref_source": "secret",
		"api_key":               "sk_live_example",
		"enabled_tools":         []interface{}{"docs_search"},
		"timeouts":              map[string]interface{}{"create": "4s"},
	})
	diff, err := r.Diff(ctx, nil, cfg, meta)
	if err != nil {
		t.Fatalf("error planning create: %s", err)
	}

	state, diags := r.Apply(ctx, nil, diff, meta)
	requireNoErrors(t, diags)
	if len(diags) != 1 || diags[0].Severity != diag.Warning {
		t.Fatalf("diagnostics = %v, expected one warning", diags)
	}
	if state == nil || state.ID != "docs" || state.Tainted {
		t.Fatalf("state = %v, expected the server to be kept and not tainted", state)
	}
	if got := state.Attributes["sync_status"]; got != "pending" {
		t.Errorf("sync_status = %q, expected the server to be read after discovery timed out", got)
	}
	if len(api.enabled) != 0 {
		t.Errorf("expected no enabled tool requests, got %v", api.enabled)
	}
}

func TestMCPServerUpdateRediscoversWhenSyncUserSet(t *testing.T) {
	t.Parallel()

	base := map[string]interface{}{
		"server_ref":            "docs",
		"endpoint":              "https://mcp.example.com/mcp",
		"credential_ref_source": "connection",
		"oauth_client_id":       "client",
		"oauth_client_secret":   "secret",
		"oauth_authorize_url":   "https://mcp.example.com/authorize",
		"oauth_token_url":       "https://mcp.example.com/token",
		"oauth_scopes":          []interface{}{"docs:read"},
	}
	rawBase := map[string]cty.Value{
		"server_ref":            cty.StringVal("docs"),
		"endpoint":              cty.StringVal("https://mcp.example.com/mcp"),
		"credential_ref_source": cty.StringVal("connection"),
		"oauth_client_id":       cty.StringVal("client"),
		"oauth_client_secret":   cty.StringVal("secret"),
		"oauth_authorize_url":   cty.StringVal("https://mcp.example.com/authorize"),
		"oauth_token_url":       cty.StringVal("https://mcp.example.com/token"),
		"oauth_scopes":          cty.SetVal([]cty.Value{cty.StringVal("docs:read")}),
		"sync_user_id":          cty.StringVal("alice"),
	}
	with := func(m map[string]interface{}, extra map[string]interface{}) map[string]interface{} {
		out := map[string]interface{}{}
		for k, v := range m {
			out[k] = v
		}
		for k, v := range extra {
			out[k] = v
		}
		return out
	}

	tests := []struct {
		name        string
		prior       map[string]interface{}
		config      map[string]interface{}
		raw         map[string]cty.Value
		wantEnabled [][]string
	}{
		{
			name:   "sync user added, enabled tools unmanaged",
			prior:  base,
			config: with(base, map[string]interface{}{"sync_user_id": "alice"}),
			raw:    rawBase,
		},
		{
			name:        "sync user changed, enabled tools managed",
			prior:       with(base, map[string]interface{}{"sync_user_id": "bob", "enabled_tools": []interface{}{"docs_search"}}),
			config:      with(base, map[string]interface{}{"sync_user_id": "alice", "enabled_tools": []interface{}{"docs_search"}}),
			raw:         map[string]cty.Value{"enabled_tools": cty.SetVal([]cty.Value{cty.StringVal("docs_search")})},
			wantEnabled: [][]string{{"docs_search"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			api := newFakeMCPServerAPI()
			api.tools[0]["enabled"] = true
			api.server = map[string]interface{}{
				"serverRef":           "docs",
				"endpoint":            "https://mcp.example.com/mcp",
				"transport":           "streamable_http",
				"credentialRefSource": "connection",
				"oauth_authorize_url": "https://mcp.example.com/authorize",
				"oauth_scopes":        []string{"docs:read"},
				"syncStatus":          "synced",
				"lastSyncedAt":        "2026-09-20T14:00:00Z",
			}
			meta := api.start(t)
			ctx := context.Background()
			r := ResourceDigitalOceanActionGatewayMCPServer()

			prior := schema.TestResourceDataRaw(t, r.Schema, tt.prior)
			prior.SetId("docs")
			requireNoErrors(t, resourceDigitalOceanActionGatewayMCPServerRead(ctx, prior, meta))
			state := prior.State()

			raw := map[string]cty.Value{}
			for k, v := range rawBase {
				raw[k] = v
			}
			for k, v := range tt.raw {
				raw[k] = v
			}
			state.RawConfig = testRawConfig(r, raw)

			diff, err := r.Diff(ctx, state, terraform.NewResourceConfigRaw(tt.config), meta)
			if err != nil {
				t.Fatalf("error planning: %s", err)
			}
			if diff == nil || diff.RequiresNew() {
				t.Fatalf("expected an in-place update, got:\n%s", describeDiff(diff))
			}
			if a := diff.Attributes["last_synced_at"]; a == nil || !a.NewComputed {
				t.Errorf("expected last_synced_at to be unknown in the plan, got:\n%s", describeDiff(diff))
			}
			_, enabledComputed := diff.Attributes["enabled_tools.#"]
			if want := tt.wantEnabled == nil; enabledComputed != want {
				t.Errorf("enabled_tools unknown in the plan = %t, expected %t:\n%s", enabledComputed, want, describeDiff(diff))
			}

			newState, diags := r.Apply(ctx, state, diff, meta)
			requireNoErrors(t, diags)
			if !reflect.DeepEqual(api.resyncUsers, []string{"alice"}) {
				t.Errorf("resync users = %v, expected [alice]", api.resyncUsers)
			}
			if !reflect.DeepEqual(api.enabled, tt.wantEnabled) {
				t.Errorf("enabled tool requests = %v, expected %v", api.enabled, tt.wantEnabled)
			}
			if got := newState.Attributes["last_synced_at"]; got != "2026-09-20T15:01:00Z" {
				t.Errorf("last_synced_at = %q, expected the rediscovery's time", got)
			}
		})
	}
}

func TestMCPServerEnabledToolsRequireSyncUser(t *testing.T) {
	t.Parallel()

	connection := map[string]interface{}{
		"server_ref":            "docs",
		"endpoint":              "https://mcp.example.com/mcp",
		"credential_ref_source": "connection",
		"oauth_client_id":       "client",
		"oauth_client_secret":   "secret",
		"oauth_authorize_url":   "https://mcp.example.com/authorize",
		"oauth_token_url":       "https://mcp.example.com/token",
		"oauth_scopes":          []interface{}{"docs:read"},
	}
	rawConnection := map[string]cty.Value{
		"server_ref":            cty.StringVal("docs"),
		"endpoint":              cty.StringVal("https://mcp.example.com/mcp"),
		"credential_ref_source": cty.StringVal("connection"),
		"oauth_client_id":       cty.StringVal("client"),
		"oauth_client_secret":   cty.StringVal("secret"),
		"oauth_authorize_url":   cty.StringVal("https://mcp.example.com/authorize"),
		"oauth_token_url":       cty.StringVal("https://mcp.example.com/token"),
		"oauth_scopes":          cty.SetVal([]cty.Value{cty.StringVal("docs:read")}),
		"enabled_tools":         cty.SetVal([]cty.Value{cty.StringVal("docs_search")}),
	}
	withTools := map[string]interface{}{"enabled_tools": []interface{}{"docs_search"}}
	for k, v := range connection {
		withTools[k] = v
	}

	t.Run("create", func(t *testing.T) {
		api := newFakeMCPServerAPI()
		meta := api.start(t)
		d := schema.TestResourceDataRaw(t, ResourceDigitalOceanActionGatewayMCPServer().Schema, withTools)

		diags := resourceDigitalOceanActionGatewayMCPServerCreate(context.Background(), d, meta)
		if !diags.HasError() || !strings.Contains(diags[0].Summary, "requires sync_user_id") {
			t.Fatalf("diagnostics = %v, expected an error requiring sync_user_id", diags)
		}
		if api.created || d.Id() != "" {
			t.Errorf("expected the server not to be registered")
		}
	})

	plans := []struct {
		name  string
		prior map[string]interface{}
	}{
		{name: "create plan"},
		{name: "update plan", prior: connection},
	}
	for _, tt := range plans {
		t.Run(tt.name, func(t *testing.T) {
			api := newFakeMCPServerAPI()
			api.server = map[string]interface{}{
				"serverRef":           "docs",
				"endpoint":            "https://mcp.example.com/mcp",
				"transport":           "streamable_http",
				"credentialRefSource": "connection",
				"oauth_authorize_url": "https://mcp.example.com/authorize",
				"oauth_scopes":        []string{"docs:read"},
				"syncStatus":          "pending",
			}
			meta := api.start(t)
			ctx := context.Background()
			r := ResourceDigitalOceanActionGatewayMCPServer()

			state := &terraform.InstanceState{}
			if tt.prior != nil {
				prior := schema.TestResourceDataRaw(t, r.Schema, tt.prior)
				prior.SetId("docs")
				requireNoErrors(t, resourceDigitalOceanActionGatewayMCPServerRead(ctx, prior, meta))
				state = prior.State()
			}
			state.RawConfig = testRawConfig(r, rawConnection)

			_, err := r.Diff(ctx, state, terraform.NewResourceConfigRaw(withTools), meta)
			if err == nil || !strings.Contains(err.Error(), "requires sync_user_id") {
				t.Fatalf("expected the plan to fail requiring sync_user_id, got %v", err)
			}
		})
	}
}

func TestValidateMCPServerCredentialConfig(t *testing.T) {
	config := func(attrs map[string]cty.Value) cty.Value {
		all := map[string]cty.Value{
			"credential_ref_source": cty.NullVal(cty.String),
			"credential_ref":        cty.NullVal(cty.String),
			"api_key":               cty.NullVal(cty.String),
			"oauth_client_id":       cty.NullVal(cty.String),
			"oauth_client_secret":   cty.NullVal(cty.String),
			"oauth_authorize_url":   cty.NullVal(cty.String),
			"oauth_token_url":       cty.NullVal(cty.String),
			"oauth_scopes":          cty.NullVal(cty.Set(cty.String)),
			"sync_user_id":          cty.NullVal(cty.String),
			"enabled_tools":         cty.NullVal(cty.Set(cty.String)),
		}
		for k, v := range attrs {
			all[k] = v
		}
		return cty.ObjectVal(all)
	}
	oauth := map[string]cty.Value{
		"credential_ref_source": cty.StringVal("connection"),
		"oauth_client_id":       cty.StringVal("client"),
		"oauth_client_secret":   cty.StringVal("secret"),
		"oauth_authorize_url":   cty.StringVal("https://mcp.example.com/authorize"),
		"oauth_token_url":       cty.StringVal("https://mcp.example.com/token"),
		"oauth_scopes":          cty.SetVal([]cty.Value{cty.StringVal("read")}),
	}
	without := func(attrs map[string]cty.Value, key string) map[string]cty.Value {
		out := map[string]cty.Value{}
		for k, v := range attrs {
			if k != key {
				out[k] = v
			}
		}
		return out
	}
	with := func(attrs map[string]cty.Value, key string, v cty.Value) map[string]cty.Value {
		out := without(attrs, key)
		out[key] = v
		return out
	}

	tests := []struct {
		name    string
		raw     cty.Value
		wantErr string
	}{
		{name: "no credentials", raw: config(nil)},
		{
			name:    "secret without a secret",
			raw:     config(map[string]cty.Value{"credential_ref_source": cty.StringVal("secret")}),
			wantErr: "requires one of credential_ref or api_key",
		},
		{
			name: "secret with an api key",
			raw:  config(map[string]cty.Value{"credential_ref_source": cty.StringVal("secret"), "api_key": cty.StringVal("sk")}),
		},
		{
			name: "secret with an unknown reference",
			raw:  config(map[string]cty.Value{"credential_ref_source": cty.StringVal("secret"), "credential_ref": cty.UnknownVal(cty.String)}),
		},
		{name: "connection with every OAuth argument", raw: config(oauth)},
		{
			name:    "connection without a token URL",
			raw:     config(without(oauth, "oauth_token_url")),
			wantErr: "requires oauth_token_url",
		},
		{
			name: "connection with partly unknown scopes",
			raw:  config(with(oauth, "oauth_scopes", cty.SetVal([]cty.Value{cty.UnknownVal(cty.String)}))),
		},
		{
			name: "unknown credential source",
			raw:  config(map[string]cty.Value{"credential_ref_source": cty.UnknownVal(cty.String)}),
		},
		{
			name:    "connection enabling tools without a sync user",
			raw:     config(with(oauth, "enabled_tools", cty.SetVal([]cty.Value{cty.StringVal("docs_search")}))),
			wantErr: "enabled_tools requires sync_user_id",
		},
		{
			name:    "connection enabling partly unknown tools without a sync user",
			raw:     config(with(oauth, "enabled_tools", cty.SetVal([]cty.Value{cty.UnknownVal(cty.String)}))),
			wantErr: "enabled_tools requires sync_user_id",
		},
		{
			name: "connection enabling tools with a sync user",
			raw:  config(with(with(oauth, "enabled_tools", cty.SetVal([]cty.Value{cty.StringVal("docs_search")})), "sync_user_id", cty.StringVal("alice"))),
		},
		{
			name: "connection enabling tools with an unknown sync user",
			raw:  config(with(with(oauth, "enabled_tools", cty.SetVal([]cty.Value{cty.StringVal("docs_search")})), "sync_user_id", cty.UnknownVal(cty.String))),
		},
		{
			name: "connection with empty enabled tools and no sync user",
			raw:  config(with(oauth, "enabled_tools", cty.SetValEmpty(cty.String))),
		},
		{
			name: "connection with unknown enabled tools and no sync user",
			raw:  config(with(oauth, "enabled_tools", cty.UnknownVal(cty.Set(cty.String)))),
		},
		{
			name: "secret enabling tools without a sync user",
			raw: config(map[string]cty.Value{
				"credential_ref_source": cty.StringVal("secret"),
				"api_key":               cty.StringVal("sk"),
				"enabled_tools":         cty.SetVal([]cty.Value{cty.StringVal("docs_search")}),
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMCPServerCredentialConfig(tt.raw)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %s", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestProvidersDataSourceFiltersListAttributes(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v2/action-gateway/tools/providers", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, http.StatusOK, map[string]interface{}{"providers": []interface{}{
			map[string]interface{}{"name": "github", "auth_types": []string{"digitalocean_oauth"}, "scopes": []string{"repo"}},
			map[string]interface{}{"name": "jira", "auth_types": []string{"digitalocean_oauth", "team_api_key"}, "scopes": []string{"read:jira-work", "write:jira-work"}},
			map[string]interface{}{"name": "linear", "auth_type": "team_api_key"},
		}})
	})
	meta := newTestMeta(t, mux)

	tests := []struct {
		name   string
		filter map[string]interface{}
		want   []string
	}{
		{name: "auth_types", filter: map[string]interface{}{"key": "auth_types", "values": []interface{}{"team_api_key"}}, want: []string{"jira", "linear"}},
		{name: "scopes", filter: map[string]interface{}{"key": "scopes", "values": []interface{}{"write:jira-work"}}, want: []string{"jira"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := DataSourceDigitalOceanActionGatewayProviders()
			d := schema.TestResourceDataRaw(t, r.Schema, map[string]interface{}{
				"filter": []interface{}{tt.filter},
				"sort":   []interface{}{map[string]interface{}{"key": "name"}},
			})
			requireNoErrors(t, r.ReadContext(context.Background(), d, meta))

			var got []string
			for _, p := range d.Get("providers").([]interface{}) {
				got = append(got, p.(map[string]interface{})["name"].(string))
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("providers = %v, expected %v", got, tt.want)
			}
		})
	}
}
