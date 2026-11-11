package actiongateway

import (
	"context"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	testEndpointID      = "1b4d4e0f-9a7c-4d5b-8f3c-2a6e6e8f0a13"
	testEndpointVPCUUID = "5a4981aa-9653-4bd1-bef5-d6bff52042e4"
)

// fakeEndpointAPI is an in-memory endpoint lineage: create makes version 1,
// each publish of a different spec adds a version, and read reports the head.
type fakeEndpointAPI struct {
	mu        sync.Mutex
	name      string
	actorID   string
	insights  map[string]interface{}
	versions  []map[string]interface{}
	next      int
	publishes int
	patches   []map[string]interface{}
}

func newFakeEndpointAPI() *fakeEndpointAPI {
	return &fakeEndpointAPI{name: "support-agent"}
}

func (api *fakeEndpointAPI) version() map[string]interface{} {
	number := len(api.versions)
	return api.versions[number-1]
}

func (api *fakeEndpointAPI) endpoint() map[string]interface{} {
	head := map[string]interface{}{}
	if len(api.versions) > 0 {
		head = api.version()
	}
	out := map[string]interface{}{
		"id":           testEndpointID,
		"name":         api.name,
		"actor_id":     api.actorID,
		"is_default":   false,
		"head_version": head,
		"mcp_url":      "https://actions.do-ai.run/mcp/endpoint/" + testEndpointID,
	}
	if api.insights != nil {
		out["insights"] = api.insights
	}
	return out
}

func (api *fakeEndpointAPI) newVersion(body map[string]interface{}) map[string]interface{} {
	api.next++
	spec := body["version"]
	if spec == nil {
		spec = body
	}
	specMap, _ := spec.(map[string]interface{})
	if specMap == nil {
		specMap = map[string]interface{}{}
	}
	if actor, ok := specMap["actor_id"].(string); ok && api.actorID == "" {
		api.actorID = actor
	}
	number := api.next
	version := map[string]interface{}{
		"id":                           testEndpointID + ".v" + strconv.Itoa(number),
		"endpoint_id":                  testEndpointID,
		"version_number":               number,
		"name":                         api.name,
		"actor_id":                     api.actorID,
		"mcp_url":                      "https://actions.do-ai.run/mcp/endpoint/" + testEndpointID + "/version/v" + strconv.Itoa(number),
		"created_at":                   "2026-11-10T12:00:00Z",
		"updated_at":                   "2026-11-10T12:00:00Z",
		"promoted_from_version_number": 0,
	}
	if policy, ok := specMap["policy"]; ok {
		version["policy"] = policy
	}
	if tools, ok := specMap["tools"]; ok {
		// The API canonicalizes the selection into pinned references.
		refs := []interface{}{}
		for _, raw := range tools.([]interface{}) {
			slug, pin := splitToolPin(raw.(string))
			if pin == "" {
				pin = "1"
			}
			kind := "SESSION_TOOL_REFERENCE_KIND_TOOL"
			if name, ok := trimToolbeltPrefix(slug); ok {
				kind, slug = "SESSION_TOOL_REFERENCE_KIND_TOOLBELT", name
			}
			refs = append(refs, map[string]interface{}{"kind": kind, "name": slug, "version": pin})
		}
		version["tools"] = map[string]interface{}{"references": refs}
	}
	if config, ok := specMap["config"]; ok {
		version["config"] = config
	}
	if network, ok := specMap["network"]; ok {
		binding := network.(map[string]interface{})
		// The wire spelling of the version's VPC binding is vpcUuid.
		version["network"] = map[string]interface{}{"vpcUuid": binding["vpcUuid"], "region": "nyc3"}
	}
	api.versions = append(api.versions, version)
	return version
}

func trimToolbeltPrefix(ref string) (string, bool) {
	const prefix = "toolbelt:"
	if len(ref) > len(prefix) && ref[:len(prefix)] == prefix {
		return ref[len(prefix):], true
	}
	return ref, false
}

func (api *fakeEndpointAPI) start(t *testing.T) *config.CombinedConfig {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /v2/action-gateway/endpoints", func(w http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		defer api.mu.Unlock()
		body := decodeJSON(t, r)
		api.name, _ = body["name"].(string)
		if insights, ok := body["insights"].(map[string]interface{}); ok {
			api.insights = insights
		}
		version := api.newVersion(body)
		writeJSON(t, w, http.StatusCreated, map[string]interface{}{
			"endpoint":        api.endpoint(),
			"version":         version,
			"mcp_url":         api.endpoint()["mcp_url"],
			"version_mcp_url": version["mcp_url"],
		})
	})
	mux.HandleFunc("GET /v2/action-gateway/endpoints/"+testEndpointID, func(w http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		defer api.mu.Unlock()
		writeJSON(t, w, http.StatusOK, map[string]interface{}{"endpoint": api.endpoint()})
	})
	mux.HandleFunc("PATCH /v2/action-gateway/endpoints/"+testEndpointID, func(w http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		defer api.mu.Unlock()
		body := decodeJSON(t, r)
		api.patches = append(api.patches, body)
		if name, ok := body["name"].(string); ok && name != "" {
			api.name = name
		}
		if insights, ok := body["insights"].(map[string]interface{}); ok {
			if api.insights == nil {
				api.insights = map[string]interface{}{}
			}
			for k, v := range insights {
				api.insights[k] = v
			}
		}
		writeJSON(t, w, http.StatusOK, map[string]interface{}{"endpoint": api.endpoint()})
	})
	mux.HandleFunc("POST /v2/action-gateway/endpoints/"+testEndpointID+"/versions", func(w http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		defer api.mu.Unlock()
		api.publishes++
		version := api.newVersion(decodeJSON(t, r))
		writeJSON(t, w, http.StatusCreated, map[string]interface{}{
			"endpoint":        api.endpoint(),
			"version":         version,
			"mcp_url":         api.endpoint()["mcp_url"],
			"version_mcp_url": version["mcp_url"],
		})
	})
	mux.HandleFunc("DELETE /v2/action-gateway/endpoints/"+testEndpointID, func(w http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		defer api.mu.Unlock()
		api.versions = nil
		w.WriteHeader(http.StatusNoContent)
	})
	return newTestMeta(t, mux)
}

// endpointRawConfig builds the configuration value Terraform sends the
// provider, so a test sees the same omitted-versus-empty distinction the
// provider does. version holds only the attributes the test sets; the rest are
// null.
func endpointRawConfig(r *schema.Resource, name string, version map[string]cty.Value, insights map[string]cty.Value) cty.Value {
	versionAttrs := map[string]cty.Value{}
	for attr, ty := range endpointVersionType().AttributeTypes() {
		versionAttrs[attr] = cty.NullVal(ty)
	}
	for attr, v := range version {
		versionAttrs[attr] = v
	}
	attrs := map[string]cty.Value{
		"name":    cty.StringVal(name),
		"version": cty.ListVal([]cty.Value{cty.ObjectVal(versionAttrs)}),
	}
	if insights != nil {
		insightAttrs := map[string]cty.Value{}
		for attr, ty := range insightsType().AttributeTypes() {
			insightAttrs[attr] = cty.NullVal(ty)
		}
		for attr, v := range insights {
			insightAttrs[attr] = v
		}
		attrs["insights"] = cty.ListVal([]cty.Value{cty.ObjectVal(insightAttrs)})
	}
	return testRawConfig(r, attrs)
}

func endpointVersionType() cty.Type {
	return cty.Object(map[string]cty.Type{
		"actor_id":                     cty.String,
		"tools":                        cty.Set(cty.String),
		"pinned_tools":                 cty.Set(cty.String),
		"policy":                       cty.List(cty.Object(map[string]cty.Type{"default_action": cty.String, "rules": cty.List(cty.Object(map[string]cty.Type{"tool": cty.String, "action": cty.String, "match": cty.Map(cty.String)}))})),
		"config":                       cty.String,
		"network":                      cty.List(cty.Object(map[string]cty.Type{"vpc_uuid": cty.String, "region": cty.String})),
		"id":                           cty.String,
		"version_number":               cty.Number,
		"promoted_from_version_number": cty.Number,
	})
}

func insightsType() cty.Type {
	return cty.Object(map[string]cty.Type{"metrics": cty.Bool, "logs": cty.Bool, "traces": cty.Bool})
}

// applyEndpointConfig plans and applies one configuration, the way Terraform
// calls the provider.
func applyEndpointConfig(t *testing.T, r *schema.Resource, meta *config.CombinedConfig, state *terraform.InstanceState, config map[string]interface{}, raw cty.Value) *terraform.InstanceState {
	t.Helper()
	ctx := context.Background()
	diff, err := r.Diff(ctx, state, terraform.NewResourceConfigRaw(config), meta)
	if err != nil {
		t.Fatalf("error planning: %s", err)
	}
	if diff == nil || diff.Empty() {
		// Nothing to apply, which is what Terraform does with an empty plan.
		return state
	}
	diff.RawConfig = raw
	next, diags := r.Apply(ctx, state, diff, meta)
	requireNoErrors(t, diags)
	return next
}

func planEndpointConfig(t *testing.T, r *schema.Resource, meta *config.CombinedConfig, state *terraform.InstanceState, config map[string]interface{}) *terraform.InstanceDiff {
	t.Helper()
	diff, err := r.Diff(context.Background(), state, terraform.NewResourceConfigRaw(config), meta)
	if err != nil {
		t.Fatalf("error planning: %s", err)
	}
	return diff
}

func TestEndpointCreatePublishesVersionOne(t *testing.T) {
	t.Parallel()

	api := newFakeEndpointAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayEndpoint()

	name := "support-agent"
	version := map[string]interface{}{
		"actor_id": "user-123",
		"tools":    []interface{}{"exa_web_search"},
	}
	config := map[string]interface{}{
		"name":     name,
		"version":  []interface{}{version},
		"insights": []interface{}{map[string]interface{}{"logs": true}},
	}
	raw := endpointRawConfig(r, name, map[string]cty.Value{
		"actor_id": cty.StringVal("user-123"),
		"tools":    cty.SetVal([]cty.Value{cty.StringVal("exa_web_search")}),
	}, map[string]cty.Value{"logs": cty.True})

	state := applyEndpointConfig(t, r, meta, nil, config, raw)
	if got := state.ID; got != testEndpointID {
		t.Fatalf("id = %q, expected the endpoint UUID", got)
	}
	d := r.Data(state)

	if got := d.Get("version.0.version_number"); got != 1 {
		t.Errorf("version_number = %v, expected 1", got)
	}
	if got := d.Get("version.0.id"); got != testEndpointID+".v1" {
		t.Errorf("version id = %v", got)
	}
	if got := expandStringSet(d.Get("version.0.tools")); !reflect.DeepEqual(got, []string{"exa_web_search"}) {
		t.Errorf("tools = %v, expected the configured unpinned slug", got)
	}
	if got := expandStringSet(d.Get("version.0.pinned_tools")); !reflect.DeepEqual(got, []string{"exa_web_search@1"}) {
		t.Errorf("pinned_tools = %v", got)
	}
	if got := d.Get("version.0.actor_id"); got != "user-123" {
		t.Errorf("actor_id = %v", got)
	}
	if got := d.Get("version_mcp_url"); got != "https://actions.do-ai.run/mcp/endpoint/"+testEndpointID+"/version/v1" {
		t.Errorf("version_mcp_url = %v", got)
	}
	api.mu.Lock()
	defer api.mu.Unlock()
	if enabled, _ := api.insights["logs"].(bool); !enabled {
		t.Errorf("insights = %v, expected the create request to carry logs=true", api.insights)
	}
}

func TestEndpointCreateSendsOmittedToolsAsAll(t *testing.T) {
	t.Parallel()

	api := newFakeEndpointAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayEndpoint()

	name := "support-agent"
	config := map[string]interface{}{
		"name":    name,
		"version": []interface{}{map[string]interface{}{}},
	}
	raw := endpointRawConfig(r, name, nil, nil)

	applyEndpointConfig(t, r, meta, nil, config, raw)

	if got, sent := api.version()["tools"]; sent {
		t.Errorf("tools = %v, expected the version to carry no tools selection", got)
	}
	if got := api.version()["tools"]; got != nil {
		t.Errorf("tools = %v, expected the request to omit tools", got)
	}
}

func TestEndpointExplicitlyEmptyToolsSendsNone(t *testing.T) {
	t.Parallel()

	api := newFakeEndpointAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayEndpoint()

	name := "support-agent"
	config := map[string]interface{}{
		"name": name,
		"version": []interface{}{map[string]interface{}{
			"tools": []interface{}{},
		}},
	}
	raw := endpointRawConfig(r, name, map[string]cty.Value{
		"tools": cty.SetValEmpty(cty.String),
	}, nil)

	applyEndpointConfig(t, r, meta, nil, config, raw)

	refs, ok := api.version()["tools"].(map[string]interface{})
	if !ok {
		t.Fatalf("tools = %v, expected an explicit empty selection", api.version()["tools"])
	}
	if got := refs["references"]; !reflect.DeepEqual(got, []interface{}{}) {
		t.Errorf("references = %v, expected none", got)
	}
	if count := api.version()["tools"]; count == nil {
		t.Errorf("expected the request to carry an explicit empty tools list")
	}
}

func TestEndpointOmittedOptionalBlocksStayEmpty(t *testing.T) {
	t.Parallel()

	api := newFakeEndpointAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayEndpoint()

	name := "support-agent"
	config := map[string]interface{}{
		"name":    name,
		"version": []interface{}{map[string]interface{}{"tools": []interface{}{"exa_web_search"}}},
	}
	raw := endpointRawConfig(r, name, map[string]cty.Value{
		"tools": cty.SetVal([]cty.Value{cty.StringVal("exa_web_search")}),
	}, nil)
	state := applyEndpointConfig(t, r, meta, nil, config, raw)

	// The API reports an omitted policy as an empty object and no insights at
	// all, which must not show up as a persistent diff.
	if plan := planEndpointConfig(t, r, meta, state, config); plan != nil && !plan.Empty() {
		t.Errorf("expected an empty plan with the optional blocks omitted, got:\n%s", describeDiff(plan))
	}
	if got := r.Data(state).Get("version.0.policy.#"); got != 0 {
		t.Errorf("policy.# = %v, expected an omitted policy to stay empty", got)
	}
	if got := r.Data(state).Get("insights.#"); got != 0 {
		t.Errorf("insights.# = %v, expected omitted consent to stay empty", got)
	}
}

func TestEndpointPolicyWithoutDefaultActionKeepsEmpty(t *testing.T) {
	t.Parallel()

	api := newFakeEndpointAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayEndpoint()

	name := "support-agent"
	config := map[string]interface{}{
		"name": name,
		"version": []interface{}{map[string]interface{}{
			"tools": []interface{}{"exa_web_search"},
			"policy": []interface{}{map[string]interface{}{
				"rules": []interface{}{map[string]interface{}{"tool": "exa_web_search", "action": "allow"}},
			}},
		}},
	}
	raw := endpointRawConfig(r, name, map[string]cty.Value{
		"tools": cty.SetVal([]cty.Value{cty.StringVal("exa_web_search")}),
		"policy": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
			"default_action": cty.NullVal(cty.String),
			"rules": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
				"tool":   cty.StringVal("exa_web_search"),
				"action": cty.StringVal("allow"),
				"match":  cty.NullVal(cty.Map(cty.String)),
			})}),
		})}),
	}, nil)
	state := applyEndpointConfig(t, r, meta, nil, config, raw)

	if plan := planEndpointConfig(t, r, meta, state, config); plan != nil && !plan.Empty() {
		t.Errorf("expected an empty plan with default_action omitted, got:\n%s", describeDiff(plan))
	}
	d := r.Data(state)
	if got := d.Get("version.0.policy.0.rules.#"); got != 1 {
		t.Errorf("rules.# = %v, expected the rule to round trip", got)
	}
	if got := d.Get("version.0.policy.0.rules.0.action"); got != "allow" {
		t.Errorf("action = %v, expected the lowercase vocabulary", got)
	}
}

func TestEndpointVersionToolsSelection(t *testing.T) {
	t.Parallel()

	r := ResourceDigitalOceanActionGatewayEndpoint()
	tests := []struct {
		name  string
		tools cty.Value
		known bool
		want  []string
	}{
		{
			name:  "omitted",
			tools: cty.NullVal(cty.Set(cty.String)),
			known: true,
		},
		{
			name:  "explicitly empty",
			tools: cty.SetValEmpty(cty.String),
			known: true,
			want:  []string{},
		},
		{
			name:  "one unpinned",
			tools: cty.SetVal([]cty.Value{cty.StringVal("exa_web_search")}),
			known: true,
			want:  []string{"exa_web_search"},
		},
		{
			name:  "unknown at plan time",
			tools: cty.UnknownVal(cty.Set(cty.String)),
			known: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			raw := testRawConfig(r, map[string]cty.Value{
				"name": cty.StringVal("support-agent"),
				"version": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
					"tools": tt.tools,
				})}),
			})
			got, known := versionToolsSelection(raw)
			if known != tt.known {
				t.Fatalf("known = %v, expected %v", known, tt.known)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tools = %v, expected %v", got, tt.want)
			}
		})
	}
}

func TestEndpointPlanAfterApplyIsEmpty(t *testing.T) {
	t.Parallel()

	api := newFakeEndpointAPI()
	meta := api.start(t)

	r := ResourceDigitalOceanActionGatewayEndpoint()
	cfg := terraform.NewResourceConfigRaw(map[string]interface{}{
		"name": "support-agent",
		"version": []interface{}{map[string]interface{}{
			"actor_id": "user-123",
			"tools":    []interface{}{"exa_web_search"},
			"policy": []interface{}{map[string]interface{}{
				"default_action": "ask",
				"rules": []interface{}{map[string]interface{}{
					"tool":   "exa_web_search",
					"action": "allow",
				}},
			}},
			"config": `{"customInstructions":"Be brief."}`,
		}},
	})
	raw := testRawConfig(r, map[string]cty.Value{
		"name": cty.StringVal("support-agent"),
		"version": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
			"actor_id": cty.StringVal("user-123"),
			"tools":    cty.SetVal([]cty.Value{cty.StringVal("exa_web_search")}),
		})}),
	})

	ctx := context.Background()
	diff, err := r.Diff(ctx, nil, cfg, meta)
	if err != nil {
		t.Fatalf("error planning create: %s", err)
	}
	diff.RawConfig = raw
	state, diags := r.Apply(ctx, nil, diff, meta)
	requireNoErrors(t, diags)

	state, diags = r.RefreshWithoutUpgrade(ctx, state, meta)
	requireNoErrors(t, diags)

	next, err := r.Diff(ctx, state, cfg, meta)
	if err != nil {
		t.Fatalf("error planning after apply: %s", err)
	}
	if next != nil && !next.Empty() {
		t.Errorf("expected an empty plan after apply, got:\n%s", describeDiff(next))
	}
}

func TestEndpointUpdatePatchesWithoutReplacement(t *testing.T) {
	t.Parallel()

	api := newFakeEndpointAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayEndpoint()

	name := "support-agent"
	version := map[string]interface{}{"actor_id": "user-123", "tools": []interface{}{"exa_web_search"}}
	insights := map[string]interface{}{"logs": true, "traces": false}
	base := map[string]interface{}{
		"name":     name,
		"version":  []interface{}{version},
		"insights": []interface{}{insights},
	}
	rawVersion := map[string]cty.Value{
		"actor_id": cty.StringVal("user-123"),
		"tools":    cty.SetVal([]cty.Value{cty.StringVal("exa_web_search")}),
	}
	rawInsights := map[string]cty.Value{"logs": cty.True, "traces": cty.False}
	state := applyEndpointConfig(t, r, meta, nil, base, endpointRawConfig(r, name, rawVersion, rawInsights))

	if plan := planEndpointConfig(t, r, meta, state, base); plan != nil && !plan.Empty() {
		t.Fatalf("expected an empty plan after apply, got:\n%s", describeDiff(plan))
	}

	// A rename alone patches the endpoint and publishes nothing.
	renamed := map[string]interface{}{"name": "renamed-agent", "version": []interface{}{version}, "insights": base["insights"]}
	state = applyEndpointConfig(t, r, meta, state, renamed, endpointRawConfig(r, "renamed-agent", rawVersion, rawInsights))
	if api.publishes != 0 {
		t.Errorf("publishes = %d, expected a rename to publish no version", api.publishes)
	}
	if len(api.patches) != 1 || api.patches[0]["name"] != "renamed-agent" {
		t.Errorf("patches = %v, expected one rename", api.patches)
	}
	if _, sent := api.patches[0]["insights"]; sent {
		t.Errorf("patch = %v, expected a rename to send no insights", api.patches[0])
	}

	// A tools change publishes a new head in place and leaves the old one live.
	published := map[string]interface{}{
		"name":     "renamed-agent",
		"version":  []interface{}{map[string]interface{}{"actor_id": "user-123", "tools": []interface{}{"exa_web_fetch"}}},
		"insights": base["insights"],
	}
	rawPublished := map[string]cty.Value{"actor_id": cty.StringVal("user-123"), "tools": cty.SetVal([]cty.Value{cty.StringVal("exa_web_fetch")})}
	state = applyEndpointConfig(t, r, meta, state, published, endpointRawConfig(r, "renamed-agent", rawPublished, rawInsights))

	if api.publishes != 1 {
		t.Fatalf("publishes = %d, expected exactly one version", api.publishes)
	}
	if len(api.versions) != 2 {
		t.Errorf("versions = %d, expected the previous version to stay live", len(api.versions))
	}
	d := r.Data(state)
	if got := d.Get("version.0.version_number"); got != 2 {
		t.Errorf("version_number = %v, expected the new head", got)
	}
	if got := d.Id(); got != testEndpointID {
		t.Errorf("id = %q, expected a publish to keep the endpoint", got)
	}
	if got := expandStringSet(d.Get("version.0.tools")); !reflect.DeepEqual(got, []string{"exa_web_fetch"}) {
		t.Errorf("tools = %v", got)
	}
	if plan := planEndpointConfig(t, r, meta, state, published); plan != nil && !plan.Empty() {
		t.Errorf("expected an empty plan after the publish, got:\n%s", describeDiff(plan))
	}
}

// TestEndpointRenameWithOptionalBlocksDoesNotPublish covers a rename of an
// endpoint whose version sets the optional blocks: those blocks hold computed
// members, and a rename must not publish a version because of them.
func TestEndpointRenameWithOptionalBlocksDoesNotPublish(t *testing.T) {
	t.Parallel()

	api := newFakeEndpointAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayEndpoint()

	name := "support-agent"
	version := map[string]interface{}{
		"tools":   []interface{}{"exa_web_search"},
		"network": []interface{}{map[string]interface{}{"vpc_uuid": testEndpointVPCUUID}},
		"policy": []interface{}{map[string]interface{}{
			"rules": []interface{}{map[string]interface{}{"tool": "exa_web_search", "action": "allow"}},
		}},
	}
	rawVersion := map[string]cty.Value{
		"tools": cty.SetVal([]cty.Value{cty.StringVal("exa_web_search")}),
		"network": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
			"vpc_uuid": cty.StringVal(testEndpointVPCUUID), "region": cty.NullVal(cty.String),
		})}),
		"policy": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
			"default_action": cty.NullVal(cty.String),
			"rules": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
				"tool": cty.StringVal("exa_web_search"), "action": cty.StringVal("allow"), "match": cty.NullVal(cty.Map(cty.String)),
			})}),
		})}),
	}
	base := map[string]interface{}{"name": name, "version": []interface{}{version}}
	state := applyEndpointConfig(t, r, meta, nil, base, endpointRawConfig(r, name, rawVersion, nil))

	renamed := map[string]interface{}{"name": "renamed-agent", "version": []interface{}{version}}
	state = applyEndpointConfig(t, r, meta, state, renamed, endpointRawConfig(r, "renamed-agent", rawVersion, nil))

	if api.publishes != 0 {
		t.Errorf("publishes = %d, expected a rename to publish no version", api.publishes)
	}
	if got := r.Data(state).Get("name"); got != "renamed-agent" {
		t.Errorf("name = %v, expected the rename to land", got)
	}
	if plan := planEndpointConfig(t, r, meta, state, renamed); plan != nil && !plan.Empty() {
		t.Errorf("expected an empty plan after the rename, got:\n%s", describeDiff(plan))
	}
}

func TestEndpointInsightsChangePatchesWithoutPublish(t *testing.T) {
	t.Parallel()

	api := newFakeEndpointAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayEndpoint()

	name := "support-agent"
	version := map[string]interface{}{"actor_id": "user-123"}
	base := map[string]interface{}{"name": name, "version": []interface{}{version}}
	rawVersion := map[string]cty.Value{"actor_id": cty.StringVal("user-123")}
	state := applyEndpointConfig(t, r, meta, nil, base, endpointRawConfig(r, name, rawVersion, nil))

	changed := map[string]interface{}{
		"name":     name,
		"version":  []interface{}{version},
		"insights": []interface{}{map[string]interface{}{"traces": true}},
	}
	applyEndpointConfig(t, r, meta, state, changed, endpointRawConfig(r, name, rawVersion, map[string]cty.Value{"traces": cty.True}))

	if api.publishes != 0 {
		t.Errorf("publishes = %d, expected a consent change to publish no version", api.publishes)
	}
	if len(api.patches) != 1 {
		t.Fatalf("patches = %v, expected one consent update", api.patches)
	}
	insights, ok := api.patches[0]["insights"].(map[string]interface{})
	if !ok {
		t.Fatalf("patch = %v, expected an insights object", api.patches[0])
	}
	if _, sent := insights["logs"]; sent {
		t.Errorf("insights = %v, expected an omitted signal to be left out", insights)
	}
	if enabled, _ := insights["traces"].(bool); !enabled {
		t.Errorf("insights = %v, expected traces=true", insights)
	}
}

func TestEndpointReadReportsOutOfBandVersionAsDiff(t *testing.T) {
	t.Parallel()

	api := newFakeEndpointAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayEndpoint()

	name := "support-agent"
	config := map[string]interface{}{
		"name":    name,
		"version": []interface{}{map[string]interface{}{"actor_id": "user-123", "tools": []interface{}{"exa_web_search"}}},
	}
	raw := endpointRawConfig(r, name, map[string]cty.Value{
		"actor_id": cty.StringVal("user-123"),
		"tools":    cty.SetVal([]cty.Value{cty.StringVal("exa_web_search")}),
	}, nil)
	state := applyEndpointConfig(t, r, meta, nil, config, raw)

	// Another client publishes a version with different tools and policy.
	api.mu.Lock()
	api.newVersion(map[string]interface{}{"version": map[string]interface{}{
		"actor_id": "user-123",
		"tools":    []interface{}{"exa_web_fetch"},
		"policy":   map[string]interface{}{"defaultAction": "deny"},
	}})
	api.mu.Unlock()

	refreshed, diags := r.RefreshWithoutUpgrade(context.Background(), state, meta)
	requireNoErrors(t, diags)

	next := planEndpointConfig(t, r, meta, refreshed, config)
	if next == nil || next.Empty() {
		t.Fatal("expected an out-of-band publish to appear as a diff")
	}
	diff := describeDiff(next)
	if !strings.Contains(diff, "version.0.tools.") {
		t.Errorf("diff = \n%s\nexpected the head's tools to differ", diff)
	}
	if !strings.Contains(diff, "version.0.policy.") {
		t.Errorf("diff = \n%s\nexpected the head's policy to differ", diff)
	}
	if next.RequiresNew() {
		t.Errorf("diff = \n%s\nexpected a converge by publishing, not a replacement", diff)
	}
}

// TestEndpointWithoutHeadConvergesByPublishing covers an endpoint whose live
// versions were all deleted out of band: the next apply publishes the
// configured version rather than adopting an empty head.
func TestEndpointWithoutHeadConvergesByPublishing(t *testing.T) {
	t.Parallel()

	api := newFakeEndpointAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayEndpoint()

	name := "support-agent"
	config := map[string]interface{}{
		"name":    name,
		"version": []interface{}{map[string]interface{}{"actor_id": "user-123", "tools": []interface{}{"exa_web_search"}}},
	}
	raw := endpointRawConfig(r, name, map[string]cty.Value{
		"actor_id": cty.StringVal("user-123"),
		"tools":    cty.SetVal([]cty.Value{cty.StringVal("exa_web_search")}),
	}, nil)
	state := applyEndpointConfig(t, r, meta, nil, config, raw)

	// Another client deletes every live version; the endpoint itself remains.
	api.mu.Lock()
	api.versions = nil
	api.next = 1
	api.mu.Unlock()

	refreshed, diags := r.RefreshWithoutUpgrade(context.Background(), state, meta)
	requireNoErrors(t, diags)

	plan := planEndpointConfig(t, r, meta, refreshed, config)
	if plan == nil || plan.Empty() {
		t.Fatal("expected an endpoint with no live version to plan a publish")
	}
	if plan.RequiresNew() {
		t.Errorf("expected a publish in place, got:\n%s", describeDiff(plan))
	}

	state = applyEndpointConfig(t, r, meta, refreshed, config, raw)
	if api.publishes != 1 {
		t.Errorf("publishes = %d, expected one version", api.publishes)
	}
	if got := r.Data(state).Get("version.0.version_number"); got == 0 || got == nil {
		t.Errorf("version_number = %v, expected the published head", got)
	}
}

func TestEndpointImportKeepsHead(t *testing.T) {
	t.Parallel()

	api := newFakeEndpointAPI()
	api.versions = []map[string]interface{}{{
		"id":             testEndpointID + ".v4",
		"endpoint_id":    testEndpointID,
		"version_number": 4,
		"actor_id":       "user-123",
		"tools":          map[string]interface{}{"references": []interface{}{map[string]interface{}{"kind": "SESSION_TOOL_REFERENCE_KIND_TOOL", "name": "exa_web_search", "version": "3"}}},
		"network":        map[string]interface{}{"vpcUuid": testEndpointVPCUUID, "region": "nyc3"},
		"mcp_url":        "https://actions.do-ai.run/mcp/endpoint/" + testEndpointID + "/version/v4",
	}}
	api.next = 4
	meta := api.start(t)

	d := schema.TestResourceDataRaw(t, ResourceDigitalOceanActionGatewayEndpoint().Schema, map[string]interface{}{})
	d.SetId(testEndpointID)
	requireNoErrors(t, resourceDigitalOceanActionGatewayEndpointRead(context.Background(), d, meta))

	if got := d.Get("version.0.version_number"); got != 4 {
		t.Errorf("version_number = %v, expected the imported head", got)
	}
	// Without a configuration the API's own spelling is reported, pins included.
	if got := expandStringSet(d.Get("version.0.tools")); !reflect.DeepEqual(got, []string{"exa_web_search@3"}) {
		t.Errorf("tools = %v, expected the pinned spelling on import", got)
	}
	if got := d.Get("version.0.network.0.vpc_uuid"); got != testEndpointVPCUUID {
		t.Errorf("vpc_uuid = %v", got)
	}
	if got := d.Get("version.0.network.0.region"); got != "nyc3" {
		t.Errorf("region = %v, expected it to be read from the response", got)
	}
}

// TestEndpointImportThenPlanIsEmpty covers the documented import promise: the
// first plan after importing an endpoint whose head matches the configuration
// offers no change.
func TestEndpointImportThenPlanIsEmpty(t *testing.T) {
	t.Parallel()

	api := newFakeEndpointAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayEndpoint()

	// Create the endpoint the way another client would, with an unpinned tool
	// selection that the API canonicalizes.
	api.newVersion(map[string]interface{}{"version": map[string]interface{}{
		"actor_id": "user-123",
		"tools":    []interface{}{"exa_web_search"},
		"policy":   map[string]interface{}{"defaultAction": "ask"},
	}})

	imported := r.TestResourceData()
	imported.SetId(testEndpointID)
	requireNoErrors(t, resourceDigitalOceanActionGatewayEndpointRead(context.Background(), imported, meta))

	// An import reports the head's members pinned, so a configuration that
	// matches uses the pinned spelling.
	config := map[string]interface{}{
		"name": "support-agent",
		"version": []interface{}{map[string]interface{}{
			"actor_id": "user-123",
			"tools":    []interface{}{"exa_web_search@1"},
			"policy":   []interface{}{map[string]interface{}{"default_action": "ask"}},
		}},
	}
	plan := planEndpointConfig(t, r, meta, imported.State(), config)
	if plan != nil && !plan.Empty() {
		t.Errorf("expected no plan after importing a matching endpoint, got:\n%s", describeDiff(plan))
	}
}

func TestEndpointDeleteAlreadyGone(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /v2/action-gateway/endpoints/"+testEndpointID, func(w http.ResponseWriter, r *http.Request) {
		writeNotFound(t, w)
	})
	meta := newTestMeta(t, mux)

	d := schema.TestResourceDataRaw(t, ResourceDigitalOceanActionGatewayEndpoint().Schema, map[string]interface{}{})
	d.SetId(testEndpointID)
	requireNoErrors(t, resourceDigitalOceanActionGatewayEndpointDelete(context.Background(), d, meta))

	if d.Id() != "" {
		t.Errorf("id = %q, expected it to be cleared", d.Id())
	}
}

func TestEndpointDataSourceByName(t *testing.T) {
	t.Parallel()

	api := newFakeEndpointAPI()
	api.newVersion(map[string]interface{}{"version": map[string]interface{}{
		"actor_id": "user-123",
		"tools":    []interface{}{"exa_web_search"},
	}})

	mux := http.NewServeMux()
	mux.HandleFunc("GET /v2/action-gateway/endpoints", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(t, w, http.StatusOK, map[string]interface{}{
			"endpoints": []interface{}{
				map[string]interface{}{"id": testEndpointID, "name": "support-agent", "actor_id": "user-123"},
				map[string]interface{}{"id": "33333333-3333-4333-8333-333333333333", "name": "other-agent", "actor_id": "user-123"},
			},
		})
	})
	mux.HandleFunc("GET /v2/action-gateway/endpoints/"+testEndpointID, func(w http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		defer api.mu.Unlock()
		writeJSON(t, w, http.StatusOK, map[string]interface{}{"endpoint": api.endpoint()})
	})
	dataMeta := newTestMeta(t, mux)

	d := schema.TestResourceDataRaw(t, DataSourceDigitalOceanActionGatewayEndpoint().Schema, map[string]interface{}{
		"name": "support-agent",
	})
	requireNoErrors(t, dataSourceDigitalOceanActionGatewayEndpointRead(context.Background(), d, dataMeta))

	if got := d.Id(); got != testEndpointID {
		t.Errorf("id = %q, expected the resolved endpoint", got)
	}
	if got := d.Get("mcp_url"); got != api.endpoint()["mcp_url"] {
		t.Errorf("mcp_url = %v", got)
	}
	if got := d.Get("version.0.version_number"); got != 1 {
		t.Errorf("version_number = %v, expected the head", got)
	}
	if got := expandStringSet(d.Get("version.0.pinned_tools")); !reflect.DeepEqual(got, []string{"exa_web_search@1"}) {
		t.Errorf("pinned_tools = %v", got)
	}
}

func TestFlattenToolReferences(t *testing.T) {
	t.Parallel()

	got := flattenToolReferences(&godo.ActionGatewayEndpointVersionToolSet{
		References: []godo.ActionGatewaySessionToolReference{
			{Kind: "SESSION_TOOL_REFERENCE_KIND_TOOL", Name: "exa_web_search", Version: "1"},
			{Kind: "SESSION_TOOL_REFERENCE_KIND_TOOLBELT", Name: "read-only", Version: "2"},
		},
	})
	want := []string{"exa_web_search@1", "toolbelt:read-only@2"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("tools = %v, expected %v", got, want)
	}
	if got := flattenToolReferences(nil); got != nil {
		t.Errorf("tools = %v, expected nil for a version that selects every tool", got)
	}
}

func TestNormalizePolicyAction(t *testing.T) {
	t.Parallel()

	for input, want := range map[string]string{
		"allow":                             "allow",
		"SESSION_POLICY_ACTION_ASK":         "ask",
		"SESSION_POLICY_ACTION_DENY":        "deny",
		"SESSION_POLICY_ACTION_UNSPECIFIED": "",
		"":                                  "",
		"unspecified":                       "",
	} {
		if got := normalizePolicyAction(input); got != want {
			t.Errorf("normalizePolicyAction(%q) = %q, expected %q", input, got, want)
		}
	}
}

func TestFindEndpointIDRequiresOneMatch(t *testing.T) {
	t.Parallel()

	endpoints := []map[string]interface{}{
		{"id": "22222222-2222-4222-8222-222222222222", "name": "support-agent", "actor_id": "user-123"},
		{"id": "33333333-3333-4333-8333-333333333333", "name": "support-agent", "actor_id": "user-456"},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v2/action-gateway/endpoints", func(w http.ResponseWriter, r *http.Request) {
		body := map[string]interface{}{"endpoints": []interface{}{}}
		list := make([]interface{}, 0, len(endpoints))
		for _, endpoint := range endpoints {
			if name := r.URL.Query().Get("name"); name != "" && endpoint["name"] != name {
				continue
			}
			if actor := r.URL.Query().Get("actor_id"); actor != "" && endpoint["actor_id"] != actor {
				continue
			}
			list = append(list, endpoint)
		}
		body["endpoints"] = list
		writeJSON(t, w, http.StatusOK, body)
	})
	meta := newTestMeta(t, mux)
	client := meta.GodoClient()

	ctx := context.Background()

	id, err := findEndpointID(ctx, client, "support-agent", "user-123")
	if err != nil || id != "22222222-2222-4222-8222-222222222222" {
		t.Errorf("id = %q, err = %v, expected the actor's endpoint", id, err)
	}

	if _, err := findEndpointID(ctx, client, "support-agent", ""); err == nil {
		t.Error("expected an ambiguous name lookup to fail")
	} else {
		for _, want := range []string{"22222222-2222-4222-8222-222222222222", "33333333-3333-4333-8333-333333333333"} {
			if !strings.Contains(err.Error(), want) {
				t.Errorf("error %q does not name candidate %s", err, want)
			}
		}
	}

	if _, err := findEndpointID(ctx, client, "missing-agent", ""); err == nil {
		t.Error("expected a missing endpoint to fail")
	}
}

func TestEndpointNetworkRoundTripsWithoutDiff(t *testing.T) {
	t.Parallel()

	api := newFakeEndpointAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayEndpoint()

	name := "support-agent"
	config := map[string]interface{}{
		"name": name,
		"version": []interface{}{map[string]interface{}{
			"tools":   []interface{}{"exa_web_search"},
			"network": []interface{}{map[string]interface{}{"vpc_uuid": testEndpointVPCUUID}},
		}},
	}
	raw := endpointRawConfig(r, name, map[string]cty.Value{
		"tools": cty.SetVal([]cty.Value{cty.StringVal("exa_web_search")}),
		"network": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
			"vpc_uuid": cty.StringVal(testEndpointVPCUUID),
			"region":   cty.NullVal(cty.String),
		})}),
	}, nil)
	state := applyEndpointConfig(t, r, meta, nil, config, raw)

	if got := r.Data(state).Get("version.0.network.0.region"); got != "nyc3" {
		t.Errorf("region = %v, expected it to be read from the response", got)
	}
	if plan := planEndpointConfig(t, r, meta, state, config); plan != nil && !plan.Empty() {
		t.Errorf("expected an empty plan with network set, got:\n%s", describeDiff(plan))
	}
	publishes := api.publishes
	applyEndpointConfig(t, r, meta, state, config, raw)
	if api.publishes != publishes {
		t.Errorf("publishes = %d, expected a repeat apply to publish no version", api.publishes)
	}
}

func TestEndpointActorChangeForcesReplacement(t *testing.T) {
	t.Parallel()

	api := newFakeEndpointAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayEndpoint()

	name := "support-agent"
	before := map[string]interface{}{
		"name":    name,
		"version": []interface{}{map[string]interface{}{"actor_id": "user-123", "tools": []interface{}{"exa_web_search"}}},
	}
	raw := endpointRawConfig(r, name, map[string]cty.Value{
		"actor_id": cty.StringVal("user-123"),
		"tools":    cty.SetVal([]cty.Value{cty.StringVal("exa_web_search")}),
	}, nil)
	state := applyEndpointConfig(t, r, meta, nil, before, raw)

	after := map[string]interface{}{
		"name":    name,
		"version": []interface{}{map[string]interface{}{"actor_id": "user-456", "tools": []interface{}{"exa_web_search"}}},
	}
	plan := planEndpointConfig(t, r, meta, state, after)
	if plan == nil || plan.Empty() {
		t.Fatal("expected an actor change to plan a replacement")
	}
	if !plan.RequiresNew() {
		t.Errorf("expected an actor change to force a new endpoint, got:\n%s", describeDiff(plan))
	}
}

func TestEndpointPolicyRuleChangePublishes(t *testing.T) {
	t.Parallel()

	api := newFakeEndpointAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayEndpoint()

	name := "support-agent"
	rule := map[string]interface{}{"tool": "exa_web_search", "action": "allow"}
	version := map[string]interface{}{
		"tools":  []interface{}{"exa_web_search", "exa_web_fetch"},
		"policy": []interface{}{map[string]interface{}{"default_action": "ask", "rules": []interface{}{rule}}},
	}
	rawVersion := map[string]cty.Value{
		"tools": cty.SetVal([]cty.Value{cty.StringVal("exa_web_search"), cty.StringVal("exa_web_fetch")}),
		"policy": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
			"default_action": cty.StringVal("ask"),
			"rules": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
				"tool": cty.StringVal("exa_web_search"), "action": cty.StringVal("allow"), "match": cty.NullVal(cty.Map(cty.String)),
			})}),
		})}),
	}
	base := map[string]interface{}{"name": name, "version": []interface{}{version}}
	state := applyEndpointConfig(t, r, meta, nil, base, endpointRawConfig(r, name, rawVersion, nil))

	// Removing the rule is a policy change and publishes a version.
	withoutRule := map[string]interface{}{
		"name": name,
		"version": []interface{}{map[string]interface{}{
			"tools":  []interface{}{"exa_web_search", "exa_web_fetch"},
			"policy": []interface{}{map[string]interface{}{"default_action": "deny"}},
		}},
	}
	rawWithoutRule := map[string]cty.Value{
		"tools": cty.SetVal([]cty.Value{cty.StringVal("exa_web_search"), cty.StringVal("exa_web_fetch")}),
		"policy": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
			"default_action": cty.StringVal("deny"),
			"rules":          cty.ListValEmpty(cty.Object(map[string]cty.Type{"tool": cty.String, "action": cty.String, "match": cty.Map(cty.String)})),
		})}),
	}
	applyEndpointConfig(t, r, meta, state, withoutRule, endpointRawConfig(r, name, rawWithoutRule, nil))

	if api.publishes != 1 {
		t.Errorf("publishes = %d, expected a policy change to publish one version", api.publishes)
	}
	if got := api.version()["policy"]; got == nil {
		t.Fatalf("policy = %v, expected it in the published version", got)
	}
}

// TestEndpointRemovingOptionalBlocksPublishes covers taking the policy and
// network blocks back out: both are stored with the version, so removing them
// publishes a version that carries neither.
func TestEndpointRemovingOptionalBlocksPublishes(t *testing.T) {
	t.Parallel()

	api := newFakeEndpointAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayEndpoint()

	name := "support-agent"
	full := map[string]interface{}{
		"name": name,
		"version": []interface{}{map[string]interface{}{
			"actor_id": "user-123",
			"tools":    []interface{}{"exa_web_search"},
			"policy": []interface{}{map[string]interface{}{"default_action": "deny",
				"rules": []interface{}{map[string]interface{}{"tool": "exa_web_search", "action": "allow"}}}},
			"network": []interface{}{map[string]interface{}{"vpc_uuid": testEndpointVPCUUID}},
		}},
	}
	rawFull := endpointRawConfig(r, name, map[string]cty.Value{
		"actor_id": cty.StringVal("user-123"),
		"tools":    cty.SetVal([]cty.Value{cty.StringVal("exa_web_search")}),
		"policy": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
			"default_action": cty.StringVal("deny"),
			"rules": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
				"tool": cty.StringVal("exa_web_search"), "action": cty.StringVal("allow"), "match": cty.NullVal(cty.Map(cty.String)),
			})}),
		})}),
		"network": cty.ListVal([]cty.Value{cty.ObjectVal(map[string]cty.Value{
			"vpc_uuid": cty.StringVal(testEndpointVPCUUID), "region": cty.NullVal(cty.String),
		})}),
	}, nil)
	state := applyEndpointConfig(t, r, meta, nil, full, rawFull)

	reduced := map[string]interface{}{
		"name":    name,
		"version": []interface{}{map[string]interface{}{"actor_id": "user-123", "tools": []interface{}{"exa_web_search"}}},
	}
	rawReduced := endpointRawConfig(r, name, map[string]cty.Value{
		"actor_id": cty.StringVal("user-123"),
		"tools":    cty.SetVal([]cty.Value{cty.StringVal("exa_web_search")}),
	}, nil)

	plan := planEndpointConfig(t, r, meta, state, reduced)
	if plan == nil || plan.Empty() {
		t.Fatal("expected removing the blocks to plan a publish")
	}
	if plan.RequiresNew() {
		t.Errorf("expected a publish in place, got:\n%s", describeDiff(plan))
	}

	state = applyEndpointConfig(t, r, meta, state, reduced, rawReduced)
	if api.publishes != 1 {
		t.Errorf("publishes = %d, expected one version", api.publishes)
	}

	// The published version carries no policy and no network.
	published, ok := api.version()["policy"]
	if ok && published != nil {
		t.Errorf("policy = %v, expected the published version to carry none", published)
	}
	if _, ok := api.version()["network"]; ok {
		t.Errorf("network = %v, expected the published version to carry none", api.version()["network"])
	}

	if plan := planEndpointConfig(t, r, meta, state, reduced); plan != nil && !plan.Empty() {
		t.Errorf("expected an empty plan after the publish, got:\n%s", describeDiff(plan))
	}
}

func TestEndpointConfigKeyOrderDoesNotDiff(t *testing.T) {
	t.Parallel()

	api := newFakeEndpointAPI()
	meta := api.start(t)
	r := ResourceDigitalOceanActionGatewayEndpoint()

	name := "support-agent"
	original := map[string]interface{}{
		"name": name,
		"version": []interface{}{map[string]interface{}{
			"config": `{"customInstructions":"Be brief.","preloadTools":["exa_web_search"]}`,
		}},
	}
	raw := endpointRawConfig(r, name, map[string]cty.Value{
		"config": cty.StringVal(`{"customInstructions":"Be brief.","preloadTools":["exa_web_search"]}`),
	}, nil)
	state := applyEndpointConfig(t, r, meta, nil, original, raw)

	reordered := map[string]interface{}{
		"name": name,
		"version": []interface{}{map[string]interface{}{
			"config": "{ \"preloadTools\": [\"exa_web_search\"], \"customInstructions\": \"Be brief.\" }",
		}},
	}
	if plan := planEndpointConfig(t, r, meta, state, reordered); plan != nil && !plan.Empty() {
		t.Errorf("expected key order and whitespace not to diff, got:\n%s", describeDiff(plan))
	}
}

func TestEndpointVersionSchemaPins(t *testing.T) {
	t.Parallel()

	r := ResourceDigitalOceanActionGatewayEndpoint()
	version := r.Schema["version"].Elem.(*schema.Resource)

	if !version.Schema["actor_id"].ForceNew {
		t.Error("actor_id must force a new endpoint, because no operation changes an endpoint's actor")
	}
	if version.Schema["tools"].Type != schema.TypeSet {
		t.Error("tools must be a set so reordering it is not a change")
	}
	if !version.Schema["pinned_tools"].Computed || version.Schema["pinned_tools"].Optional {
		t.Error("pinned_tools must be computed only")
	}
	if got := version.Schema["config"].StateFunc("{\"b\": 1, \"a\": 2}"); got != `{"a":2,"b":1}` {
		t.Errorf("config StateFunc = %v, expected normalized JSON", got)
	}
}
