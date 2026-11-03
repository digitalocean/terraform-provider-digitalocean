package actiongateway

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
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
			conns = []map[string]interface{}{{"id": "c", "provider": "github", "user_id": "malice"}}
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

// fakeMCPServerAPI serves one MCP server. Discovery completes after
// pendingPolls reads of the server once a resync has been requested.
type fakeMCPServerAPI struct {
	mu           sync.Mutex
	server       map[string]interface{}
	tools        []map[string]interface{}
	resyncs      int
	pendingPolls int
	authorize    bool
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
		if api.resyncs > 0 && api.server["lastSyncedAt"] == "" {
			if api.pendingPolls == 0 {
				api.server["lastSyncedAt"] = "2026-09-20T15:01:00Z"
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
		api.resyncs++
		if api.authorize {
			writeJSON(t, w, http.StatusOK, map[string]interface{}{
				"mcpServer":     api.server,
				"authorization": pendingAuthorization(),
			})
			return
		}
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
