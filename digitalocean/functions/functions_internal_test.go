package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	testNamespaceID = "fn-b90faf52-2b42-49c2-9792-75edfbb6f397"
	testTimestamp   = "2022-09-14T04:16:45Z"
)

// fakeFunctionsAPI is an in-memory implementation of the Functions namespace
// and trigger endpoints, shaped like the responses documented in the public
// API reference.
type fakeFunctionsAPI struct {
	mu         sync.Mutex
	namespaces map[string]map[string]interface{}
	triggers   map[string]map[string]map[string]interface{}
	// requests records the decoded JSON body of every write, keyed by
	// "<METHOD> <path>".
	requests map[string][]map[string]interface{}
	created  int
}

func newFakeFunctionsAPI() *fakeFunctionsAPI {
	return &fakeFunctionsAPI{
		namespaces: map[string]map[string]interface{}{},
		triggers:   map[string]map[string]map[string]interface{}{},
		requests:   map[string][]map[string]interface{}{},
	}
}

func (api *fakeFunctionsAPI) addNamespace(id, label string) {
	api.namespaces[id] = map[string]interface{}{
		"namespace":  id,
		"label":      label,
		"region":     "nyc1",
		"uuid":       "b90faf52-2b42-49c2-9792-75edfbb6f397",
		"api_host":   "https://faas-nyc1-2ef2e6cc.doserverless.co",
		"key":        "secret-namespace-key",
		"created_at": testTimestamp,
		"updated_at": testTimestamp,
	}
	api.triggers[id] = map[string]map[string]interface{}{}
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeNotFound(w http.ResponseWriter) {
	writeJSON(w, http.StatusNotFound, map[string]string{"id": "not_found", "message": "The resource you requested could not be found."})
}

func (api *fakeFunctionsAPI) record(r *http.Request) map[string]interface{} {
	var body map[string]interface{}
	_ = json.NewDecoder(r.Body).Decode(&body)
	key := r.Method + " " + r.URL.Path
	api.requests[key] = append(api.requests[key], body)
	return body
}

func (api *fakeFunctionsAPI) handler(w http.ResponseWriter, r *http.Request) {
	api.mu.Lock()
	defer api.mu.Unlock()

	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/v2/functions/namespaces"), "/")
	// parts[0] is always "" because of the leading slash.
	switch {
	case len(parts) == 1 && r.Method == http.MethodGet:
		list := []map[string]interface{}{}
		for _, ns := range api.namespaces {
			// The list endpoint does not return keys.
			cp := map[string]interface{}{}
			for k, v := range ns {
				if k != "key" {
					cp[k] = v
				}
			}
			list = append(list, cp)
		}
		writeJSON(w, http.StatusOK, map[string]interface{}{"namespaces": list})
	case len(parts) == 1 && r.Method == http.MethodPost:
		body := api.record(r)
		api.created++
		id := fmt.Sprintf("fn-%d", api.created)
		api.addNamespace(id, body["label"].(string))
		api.namespaces[id]["region"] = body["region"]
		writeJSON(w, http.StatusOK, map[string]interface{}{"namespace": api.namespaces[id]})
	case len(parts) == 2:
		ns, ok := api.namespaces[parts[1]]
		if !ok {
			writeNotFound(w)
			return
		}
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, http.StatusOK, map[string]interface{}{"namespace": ns})
		case http.MethodDelete:
			delete(api.namespaces, parts[1])
			delete(api.triggers, parts[1])
			w.WriteHeader(http.StatusNoContent)
		}
	case len(parts) == 3 && parts[2] == "triggers" && r.Method == http.MethodPost:
		triggers, ok := api.triggers[parts[1]]
		if !ok {
			writeNotFound(w)
			return
		}
		body := api.record(r)
		t := map[string]interface{}{
			"namespace":         parts[1],
			"name":              body["name"],
			"function":          body["function"],
			"type":              body["type"],
			"is_enabled":        body["is_enabled"],
			"scheduled_details": body["scheduled_details"],
			"created_at":        testTimestamp,
			"updated_at":        testTimestamp,
			"scheduled_runs":    map[string]interface{}{"last_run_at": nil, "next_run_at": "2022-11-03T17:03:02Z"},
		}
		triggers[body["name"].(string)] = t
		writeJSON(w, http.StatusOK, map[string]interface{}{"trigger": t})
	case len(parts) == 4 && parts[2] == "triggers":
		t, ok := api.triggers[parts[1]][parts[3]]
		if !ok {
			writeNotFound(w)
			return
		}
		switch r.Method {
		case http.MethodGet:
			writeJSON(w, http.StatusOK, map[string]interface{}{"trigger": t})
		case http.MethodPut:
			body := api.record(r)
			if v, ok := body["is_enabled"]; ok {
				t["is_enabled"] = v
			}
			if v, ok := body["scheduled_details"]; ok {
				t["scheduled_details"] = v
			}
			writeJSON(w, http.StatusOK, map[string]interface{}{"trigger": t})
		case http.MethodDelete:
			delete(api.triggers[parts[1]], parts[3])
			w.WriteHeader(http.StatusNoContent)
		}
	default:
		t := fmt.Sprintf("unexpected request %s %s", r.Method, r.URL.Path)
		writeJSON(w, http.StatusTeapot, map[string]string{"id": "unexpected", "message": t})
	}
}

func (api *fakeFunctionsAPI) start(t *testing.T) *config.CombinedConfig {
	t.Helper()

	mux := http.NewServeMux()
	mux.HandleFunc("/v2/functions/namespaces", api.handler)
	mux.HandleFunc("/v2/functions/namespaces/", api.handler)
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

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

func TestFunctionsNamespaceLifecycle(t *testing.T) {
	t.Parallel()

	api := newFakeFunctionsAPI()
	meta := api.start(t)
	ctx := context.Background()

	d := schema.TestResourceDataRaw(t, ResourceDigitalOceanFunctionsNamespace().Schema, map[string]interface{}{
		"label":  "my-namespace",
		"region": "ams3",
	})

	if diags := resourceDigitalOceanFunctionsNamespaceCreate(ctx, d, meta); diags.HasError() {
		t.Fatalf("create: %v", diags)
	}
	if d.Id() != "fn-1" {
		t.Fatalf("id = %q, want fn-1", d.Id())
	}
	want := map[string]string{
		"label":      "my-namespace",
		"region":     "ams3",
		"api_host":   "https://faas-nyc1-2ef2e6cc.doserverless.co",
		"key":        "secret-namespace-key",
		"uuid":       "b90faf52-2b42-49c2-9792-75edfbb6f397",
		"created_at": testTimestamp,
	}
	for k, v := range want {
		if got := d.Get(k).(string); got != v {
			t.Errorf("%s = %q, want %q", k, got, v)
		}
	}
	reqs := api.requests["POST /v2/functions/namespaces"]
	if len(reqs) != 1 || reqs[0]["label"] != "my-namespace" || reqs[0]["region"] != "ams3" {
		t.Errorf("unexpected create requests: %#v", reqs)
	}

	if diags := resourceDigitalOceanFunctionsNamespaceDelete(ctx, d, meta); diags.HasError() {
		t.Fatalf("delete: %v", diags)
	}
	if _, ok := api.namespaces["fn-1"]; ok {
		t.Fatal("namespace was not deleted")
	}

	// Deleting again (e.g. removed out of band) is not an error.
	d.SetId("fn-1")
	if diags := resourceDigitalOceanFunctionsNamespaceDelete(ctx, d, meta); diags.HasError() {
		t.Fatalf("second delete: %v", diags)
	}
}

func TestFunctionsNamespaceReadRemovesMissing(t *testing.T) {
	t.Parallel()

	api := newFakeFunctionsAPI()
	meta := api.start(t)

	d := schema.TestResourceDataRaw(t, ResourceDigitalOceanFunctionsNamespace().Schema, map[string]interface{}{})
	d.SetId("fn-missing")
	if diags := resourceDigitalOceanFunctionsNamespaceRead(context.Background(), d, meta); diags.HasError() {
		t.Fatalf("read: %v", diags)
	}
	if d.Id() != "" {
		t.Fatalf("expected namespace to be removed from state, id = %q", d.Id())
	}
}

func TestFunctionsNamespaceDataSource(t *testing.T) {
	t.Parallel()

	api := newFakeFunctionsAPI()
	api.addNamespace(testNamespaceID, "prod")
	api.addNamespace("fn-dup-1", "dup")
	api.addNamespace("fn-dup-2", "dup")
	meta := api.start(t)
	ctx := context.Background()
	s := DataSourceDigitalOceanFunctionsNamespace().Schema

	t.Run("by label fetches key", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, s, map[string]interface{}{"label": "prod"})
		if diags := dataSourceDigitalOceanFunctionsNamespaceRead(ctx, d, meta); diags.HasError() {
			t.Fatalf("read: %v", diags)
		}
		if d.Id() != testNamespaceID || d.Get("namespace_id") != testNamespaceID {
			t.Errorf("id = %q, namespace_id = %q", d.Id(), d.Get("namespace_id"))
		}
		if d.Get("key") != "secret-namespace-key" {
			t.Errorf("key = %q, want key from the get endpoint", d.Get("key"))
		}
	})

	t.Run("by id", func(t *testing.T) {
		d := schema.TestResourceDataRaw(t, s, map[string]interface{}{"namespace_id": testNamespaceID})
		if diags := dataSourceDigitalOceanFunctionsNamespaceRead(ctx, d, meta); diags.HasError() {
			t.Fatalf("read: %v", diags)
		}
		if d.Get("label") != "prod" {
			t.Errorf("label = %q, want prod", d.Get("label"))
		}
	})

	for name, tc := range map[string]struct{ label, wantErr string }{
		"missing label":   {"nope", `no Functions namespace found with label "nope"`},
		"ambiguous label": {"dup", `found 2 Functions namespaces with label "dup"`},
	} {
		tc := tc
		t.Run(name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, s, map[string]interface{}{"label": tc.label})
			diags := dataSourceDigitalOceanFunctionsNamespaceRead(ctx, d, meta)
			if !diags.HasError() || !strings.Contains(diags[0].Summary, tc.wantErr) {
				t.Fatalf("diags = %v, want error containing %q", diags, tc.wantErr)
			}
		})
	}
}

func TestFunctionsTriggerLifecycle(t *testing.T) {
	t.Parallel()

	api := newFakeFunctionsAPI()
	api.addNamespace(testNamespaceID, "prod")
	meta := api.start(t)
	ctx := context.Background()

	d := schema.TestResourceDataRaw(t, ResourceDigitalOceanFunctionsTrigger().Schema, map[string]interface{}{
		"namespace_id": testNamespaceID,
		"name":         "nightly",
		"function":     "jobs/report",
		"scheduled_details": []interface{}{map[string]interface{}{
			"cron": "0 3 * * *",
			"body": `{"name": "world", "count": 2}`,
		}},
	})

	if diags := resourceDigitalOceanFunctionsTriggerCreate(ctx, d, meta); diags.HasError() {
		t.Fatalf("create: %v", diags)
	}
	if want := testNamespaceID + "/nightly"; d.Id() != want {
		t.Fatalf("id = %q, want %q", d.Id(), want)
	}

	createPath := "POST /v2/functions/namespaces/" + testNamespaceID + "/triggers"
	reqs := api.requests[createPath]
	if len(reqs) != 1 {
		t.Fatalf("create requests = %d, want 1", len(reqs))
	}
	req := reqs[0]
	if req["type"] != "SCHEDULED" || req["is_enabled"] != true || req["function"] != "jobs/report" {
		t.Errorf("unexpected create request: %#v", req)
	}
	details := req["scheduled_details"].(map[string]interface{})
	if details["cron"] != "0 3 * * *" {
		t.Errorf("cron = %v", details["cron"])
	}
	if body := details["body"].(map[string]interface{}); body["name"] != "world" || body["count"] != float64(2) {
		t.Errorf("body sent as %#v, want a JSON object", body)
	}

	if got := d.Get("next_run_at"); got != "2022-11-03T17:03:02Z" {
		t.Errorf("next_run_at = %q", got)
	}
	if got := d.Get("last_run_at"); got != "" {
		t.Errorf("last_run_at = %q, want empty for a trigger that never ran", got)
	}
	if got := d.Get("scheduled_details.0.body").(string); !suppressEquivalentJSON("", `{"name": "world", "count": 2}`, got, nil) {
		t.Errorf("body in state = %q is not equivalent to the configured body", got)
	}

	// Disabling must send an explicit false.
	if err := d.Set("enabled", false); err != nil {
		t.Fatal(err)
	}
	if diags := resourceDigitalOceanFunctionsTriggerUpdate(ctx, d, meta); diags.HasError() {
		t.Fatalf("update: %v", diags)
	}
	updatePath := "PUT /v2/functions/namespaces/" + testNamespaceID + "/triggers/nightly"
	ups := api.requests[updatePath]
	if len(ups) != 1 {
		t.Fatalf("update requests = %d, want 1", len(ups))
	}
	if v, ok := ups[0]["is_enabled"]; !ok || v != false {
		t.Errorf("update request = %#v, want is_enabled: false", ups[0])
	}
	if d.Get("enabled").(bool) {
		t.Error("enabled should be false after update")
	}

	if diags := resourceDigitalOceanFunctionsTriggerDelete(ctx, d, meta); diags.HasError() {
		t.Fatalf("delete: %v", diags)
	}
	if _, ok := api.triggers[testNamespaceID]["nightly"]; ok {
		t.Fatal("trigger was not deleted")
	}

	// A trigger deleted out of band is removed from state on refresh.
	d.SetId(makeTriggerID(testNamespaceID, "nightly"))
	if diags := resourceDigitalOceanFunctionsTriggerRead(ctx, d, meta); diags.HasError() {
		t.Fatalf("read: %v", diags)
	}
	if d.Id() != "" {
		t.Fatalf("expected trigger to be removed from state, id = %q", d.Id())
	}
}

func TestFunctionsTriggerImportAndDataSource(t *testing.T) {
	t.Parallel()

	api := newFakeFunctionsAPI()
	api.addNamespace(testNamespaceID, "prod")
	api.triggers[testNamespaceID]["hourly"] = map[string]interface{}{
		"namespace":         testNamespaceID,
		"name":              "hourly",
		"function":          "ping",
		"type":              "SCHEDULED",
		"is_enabled":        false,
		"scheduled_details": map[string]interface{}{"cron": "0 * * * *", "body": map[string]interface{}{"b": 1, "a": "x"}},
		"created_at":        testTimestamp,
		"updated_at":        testTimestamp,
		"scheduled_runs":    map[string]interface{}{"last_run_at": "2022-11-03T16:00:00Z", "next_run_at": "2022-11-03T17:00:00Z"},
	}
	meta := api.start(t)
	ctx := context.Background()

	for _, importID := range []string{testNamespaceID + ",hourly", testNamespaceID + "/hourly"} {
		d := schema.TestResourceDataRaw(t, ResourceDigitalOceanFunctionsTrigger().Schema, map[string]interface{}{})
		d.SetId(importID)
		out, err := resourceDigitalOceanFunctionsTriggerImport(ctx, d, meta)
		if err != nil {
			t.Fatalf("import %q: %s", importID, err)
		}
		if diags := resourceDigitalOceanFunctionsTriggerRead(ctx, out[0], meta); diags.HasError() {
			t.Fatalf("read after import %q: %v", importID, diags)
		}
		got := out[0]
		if got.Id() != testNamespaceID+"/hourly" || got.Get("namespace_id") != testNamespaceID || got.Get("function") != "ping" {
			t.Errorf("import %q: id=%q namespace_id=%q function=%q", importID, got.Id(), got.Get("namespace_id"), got.Get("function"))
		}
		if got.Get("scheduled_details.0.body") != `{"a":"x","b":1}` {
			t.Errorf("import %q: body = %q", importID, got.Get("scheduled_details.0.body"))
		}
	}

	d := schema.TestResourceDataRaw(t, DataSourceDigitalOceanFunctionsTrigger().Schema, map[string]interface{}{
		"namespace_id": testNamespaceID,
		"name":         "hourly",
	})
	if diags := dataSourceDigitalOceanFunctionsTriggerRead(ctx, d, meta); diags.HasError() {
		t.Fatalf("data source read: %v", diags)
	}
	if d.Get("enabled").(bool) || d.Get("scheduled_details.0.cron") != "0 * * * *" || d.Get("last_run_at") != "2022-11-03T16:00:00Z" {
		t.Errorf("unexpected data source state: enabled=%v cron=%q last_run_at=%q", d.Get("enabled"), d.Get("scheduled_details.0.cron"), d.Get("last_run_at"))
	}
}

func TestParseTriggerID(t *testing.T) {
	t.Parallel()

	for _, bad := range []string{"", "fn-1", "fn-1,", ",name", "/name"} {
		if _, _, err := parseTriggerID(bad); err == nil {
			t.Errorf("parseTriggerID(%q) succeeded, want error", bad)
		}
	}
}

func TestExpandTriggerBody(t *testing.T) {
	t.Parallel()

	if body, err := expandTriggerBody("  "); err != nil || body != nil {
		t.Errorf("blank body = %#v, %v; want nil, nil", body, err)
	}
	if body, err := expandTriggerBody(`{"nested": {"k": [1, 2]}}`); err != nil || body["nested"] == nil {
		t.Errorf("object body = %#v, %v", body, err)
	}
	for _, bad := range []string{`[1, 2]`, `"text"`, `null`, `{"a":`} {
		if _, errs := validateTriggerBody(bad, "body"); len(errs) == 0 {
			t.Errorf("validateTriggerBody(%q) passed, want error", bad)
		}
	}
}

func TestSuppressEquivalentJSON(t *testing.T) {
	t.Parallel()

	cases := []struct {
		old, new string
		want     bool
	}{
		{`{"a":1,"b":2}`, "{\n  \"b\": 2,\n  \"a\": 1\n}", true},
		{`{"a":1}`, `{"a":2}`, false},
		{`{"a":1}`, ``, false},
		{``, ``, true},
	}
	for _, c := range cases {
		if got := suppressEquivalentJSON("body", c.old, c.new, nil); got != c.want {
			t.Errorf("suppressEquivalentJSON(%q, %q) = %v, want %v", c.old, c.new, got, c.want)
		}
	}
}

// Removing the body cannot be expressed through the update endpoint, so the
// plan must replace the trigger. Changing the body or the cron updates in place.
func TestFunctionsTriggerBodyRemovalForcesReplacement(t *testing.T) {
	t.Parallel()

	r := ResourceDigitalOceanFunctionsTrigger()
	state := &terraform.InstanceState{
		ID: testNamespaceID + "/nightly",
		Attributes: map[string]string{
			"id":                       testNamespaceID + "/nightly",
			"namespace_id":             testNamespaceID,
			"name":                     "nightly",
			"function":                 "jobs/report",
			"type":                     "SCHEDULED",
			"enabled":                  "true",
			"scheduled_details.#":      "1",
			"scheduled_details.0.cron": "0 3 * * *",
			"scheduled_details.0.body": `{"a":1}`,
		},
	}

	plan := func(details map[string]interface{}) *terraform.InstanceDiff {
		t.Helper()
		cfg := terraform.NewResourceConfigRaw(map[string]interface{}{
			"namespace_id":      testNamespaceID,
			"name":              "nightly",
			"function":          "jobs/report",
			"scheduled_details": []interface{}{details},
		})
		diff, err := r.Diff(context.Background(), state, cfg, nil)
		if err != nil {
			t.Fatalf("diff: %s", err)
		}
		return diff
	}

	if diff := plan(map[string]interface{}{"cron": "0 3 * * *"}); diff == nil || !diff.RequiresNew() {
		t.Errorf("removing body: RequiresNew = false, want true (diff: %#v)", diff)
	}
	if diff := plan(map[string]interface{}{"cron": "0 3 * * *", "body": `{"a":2}`}); diff == nil || diff.RequiresNew() {
		t.Errorf("changing body: want in-place update, got %#v", diff)
	}
	if diff := plan(map[string]interface{}{"cron": "0 4 * * *", "body": `{ "a": 1 }`}); diff == nil || diff.RequiresNew() || diff.Attributes["scheduled_details.0.body"] != nil {
		t.Errorf("changing cron with an equivalent body: want in-place update without a body diff, got %#v", diff)
	}
}
