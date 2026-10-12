package reservedip

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	testReservedIP        = "192.0.2.10"
	testReservedDropletID = 606930048
)

// fakeReservedIPAPI serves the reserved IP endpoints used by the assignment
// resource. Assign/unassign requests are handled by onAction, which may
// mutate the assigned droplet and choose the response.
type fakeReservedIPAPI struct {
	mu sync.Mutex
	// dropletID is the droplet the reserved IP is assigned to; 0 when unassigned.
	dropletID int
	// staleGets is the number of upcoming Gets that report no droplet,
	// regardless of dropletID.
	staleGets   int
	actionPosts int
	onAction    func(api *fakeReservedIPAPI, w http.ResponseWriter, actionType string, dropletID int)
}

func (api *fakeReservedIPAPI) start(t *testing.T) *config.CombinedConfig {
	t.Helper()

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	mux.HandleFunc(fmt.Sprintf("/v2/droplets/%d/actions", testReservedDropletID), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"actions": []}`)
	})

	mux.HandleFunc(fmt.Sprintf("/v2/reserved_ips/%s", testReservedIP), func(w http.ResponseWriter, r *http.Request) {
		api.mu.Lock()
		defer api.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		if api.dropletID == 0 || api.staleGets > 0 {
			if api.staleGets > 0 {
				api.staleGets--
			}
			fmt.Fprintf(w, `{"reserved_ip": {"ip": %q, "droplet": null}}`, testReservedIP)
			return
		}
		fmt.Fprintf(w, `{"reserved_ip": {"ip": %q, "droplet": {"id": %d}}}`, testReservedIP, api.dropletID)
	})

	mux.HandleFunc(fmt.Sprintf("/v2/reserved_ips/%s/actions", testReservedIP), func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %v, expected %v", r.Method, http.MethodPost)
		}

		var body struct {
			Type      string `json:"type"`
			DropletID int    `json:"droplet_id"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("error decoding action request: %s", err)
		}

		api.mu.Lock()
		defer api.mu.Unlock()

		api.actionPosts++
		w.Header().Set("Content-Type", "application/json")
		api.onAction(api, w, body.Type, body.DropletID)
	})

	mux.HandleFunc(fmt.Sprintf("/v2/reserved_ips/%s/actions/1", testReservedIP), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"action": {"id": 1, "status": "completed"}}`)
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

func newTestReservedIPAssignmentData(t *testing.T) *schema.ResourceData {
	t.Helper()

	return schema.TestResourceDataRaw(t, ResourceDigitalOceanReservedIPAssignment().Schema, map[string]interface{}{
		"ip_address": testReservedIP,
		"droplet_id": testReservedDropletID,
	})
}

func writeActionInProgress(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, `{"action": {"id": 1, "status": "in-progress"}}`)
}

func writeInternalServerError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, `{"id": "server_error", "message": "Server Error"}`)
}

func writeUnprocessable(w http.ResponseWriter, message string) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	fmt.Fprintf(w, `{"id": "unprocessable_entity", "message": %q}`, message)
}

// An assign request that fails with a 5xx after being applied is retried by
// the HTTP client, and the API rejects the retry. Create must succeed because
// the reserved IP ended up on the target droplet.
func TestReservedIPAssignmentCreateRetriedAssignAlreadyApplied(t *testing.T) {
	t.Parallel()

	api := &fakeReservedIPAPI{
		onAction: func(api *fakeReservedIPAPI, w http.ResponseWriter, actionType string, dropletID int) {
			if api.actionPosts == 1 {
				api.dropletID = dropletID
				writeInternalServerError(w)
				return
			}
			writeUnprocessable(w, "Droplet is already assigned to another reserved IP.")
		},
	}
	client := api.start(t)
	d := newTestReservedIPAssignmentData(t)

	if diags := resourceDigitalOceanReservedIPAssignmentCreate(context.Background(), d, client); diags.HasError() {
		t.Fatalf("create returned error: %v", diags)
	}

	if d.Id() == "" {
		t.Error("expected resource ID to be set")
	}
	if api.actionPosts != 2 {
		t.Errorf("action posts = %d, expected 2", api.actionPosts)
	}
}

// Reads can briefly show the previous assignment after the assign action
// completes. Create must not drop the resource from state.
func TestReservedIPAssignmentCreateAssignmentNotYetVisible(t *testing.T) {
	t.Parallel()

	api := &fakeReservedIPAPI{
		staleGets: 1,
		onAction: func(api *fakeReservedIPAPI, w http.ResponseWriter, actionType string, dropletID int) {
			api.dropletID = dropletID
			writeActionInProgress(w)
		},
	}
	client := api.start(t)
	d := newTestReservedIPAssignmentData(t)

	if diags := resourceDigitalOceanReservedIPAssignmentCreate(context.Background(), d, client); diags.HasError() {
		t.Fatalf("create returned error: %v", diags)
	}

	if d.Id() == "" {
		t.Error("expected resource ID to be set")
	}
}

func TestReservedIPAssignmentDeleteRetriedUnassignAlreadyApplied(t *testing.T) {
	t.Parallel()

	api := &fakeReservedIPAPI{
		dropletID: testReservedDropletID,
		onAction: func(api *fakeReservedIPAPI, w http.ResponseWriter, actionType string, dropletID int) {
			if api.actionPosts == 1 {
				api.dropletID = 0
				writeInternalServerError(w)
				return
			}
			writeUnprocessable(w, "The reserved IP is not assigned to a Droplet.")
		},
	}
	client := api.start(t)
	d := newTestReservedIPAssignmentData(t)
	d.SetId("assignment")

	if diags := resourceDigitalOceanReservedIPAssignmentDelete(context.Background(), d, client); diags.HasError() {
		t.Fatalf("delete returned error: %v", diags)
	}

	if d.Id() != "" {
		t.Errorf("id = %q, expected it to be cleared", d.Id())
	}
	if api.actionPosts != 2 {
		t.Errorf("action posts = %d, expected 2", api.actionPosts)
	}
}

func TestReservedIPAssignmentDeleteAlreadyUnassigned(t *testing.T) {
	t.Parallel()

	api := &fakeReservedIPAPI{
		onAction: func(api *fakeReservedIPAPI, w http.ResponseWriter, actionType string, dropletID int) {
			t.Errorf("unexpected %s action", actionType)
		},
	}
	client := api.start(t)
	d := newTestReservedIPAssignmentData(t)
	d.SetId("assignment")

	if diags := resourceDigitalOceanReservedIPAssignmentDelete(context.Background(), d, client); diags.HasError() {
		t.Fatalf("delete returned error: %v", diags)
	}

	if d.Id() != "" {
		t.Errorf("id = %q, expected it to be cleared", d.Id())
	}
}
