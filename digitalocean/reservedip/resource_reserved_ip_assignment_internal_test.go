package reservedip

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/digitalocean/godo"
)

func TestReservedIPActionNotFound(t *testing.T) {
	notFoundResp := &godo.Response{Response: &http.Response{StatusCode: http.StatusNotFound}}
	okResp := &godo.Response{Response: &http.Response{StatusCode: http.StatusOK}}
	godoNotFound := &godo.ErrorResponse{Response: &http.Response{StatusCode: http.StatusNotFound}}

	tests := []struct {
		name string
		resp *godo.Response
		err  error
		want bool
	}{
		{name: "no error", resp: okResp, err: nil, want: false},
		{name: "404 via response", resp: notFoundResp, err: fmt.Errorf("not found"), want: true},
		{name: "404 via godo.ErrorResponse", resp: nil, err: godoNotFound, want: true},
		{name: "non-404 error", resp: okResp, err: fmt.Errorf("boom"), want: false},
		{name: "error with nil response", resp: nil, err: fmt.Errorf("boom"), want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := reservedIPActionNotFound(tt.resp, tt.err); got != tt.want {
				t.Errorf("reservedIPActionNotFound() = %v, want %v", got, tt.want)
			}
		})
	}
}

// newTestGodoClient returns a godo client whose requests are routed to the given
// test server.
func newTestGodoClient(t *testing.T, serverURL string) *godo.Client {
	t.Helper()
	client := godo.NewClient(nil)
	u, err := url.Parse(serverURL + "/")
	if err != nil {
		t.Fatalf("parsing test server URL: %s", err)
	}
	client.BaseURL = u
	return client
}

// TestReservedIPAssignmentRefreshFunc_ActionNotFound verifies that when the
// action-status endpoint 404s (as it does permanently for BYOIP-prefix reserved
// IPs), the refresh func falls back to the reserved IP itself and reports
// "completed"/"in-progress" based on the actual attachment, instead of erroring
// or polling the missing action forever.
func TestReservedIPAssignmentRefreshFunc_ActionNotFound(t *testing.T) {
	const (
		ip        = "192.0.2.10"
		actionID  = 123456
		dropletID = 42
	)

	tests := []struct {
		name          string
		check         reservedIPAssignmentCheck
		reservedIPDOC string // body for GET /v2/reserved_ips/{ip}
		wantState     string
		wantNilResult bool
	}{
		{
			name:          "assign completed - attached to target droplet",
			check:         reservedIPAssignmentCheck{dropletID: dropletID, assign: true},
			reservedIPDOC: fmt.Sprintf(`{"reserved_ip":{"ip":%q,"droplet":{"id":%d},"region":{"slug":"ams3"}}}`, ip, dropletID),
			wantState:     "completed",
		},
		{
			name:          "assign pending - not yet attached",
			check:         reservedIPAssignmentCheck{dropletID: dropletID, assign: true},
			reservedIPDOC: fmt.Sprintf(`{"reserved_ip":{"ip":%q,"droplet":null,"region":{"slug":"ams3"}}}`, ip),
			wantState:     "in-progress",
		},
		{
			name:          "unassign completed - detached",
			check:         reservedIPAssignmentCheck{dropletID: dropletID, assign: false},
			reservedIPDOC: fmt.Sprintf(`{"reserved_ip":{"ip":%q,"droplet":null,"region":{"slug":"ams3"}}}`, ip),
			wantState:     "completed",
		},
		{
			name:          "unassign pending - still attached to target",
			check:         reservedIPAssignmentCheck{dropletID: dropletID, assign: false},
			reservedIPDOC: fmt.Sprintf(`{"reserved_ip":{"ip":%q,"droplet":{"id":%d},"region":{"slug":"ams3"}}}`, ip, dropletID),
			wantState:     "in-progress",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := http.NewServeMux()
			// The action ID never becomes queryable: always 404.
			mux.HandleFunc(fmt.Sprintf("/v2/reserved_ips/%s/actions/%d", ip, actionID),
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprint(w, `{"id":"not_found","message":"The resource you requested could not be found."}`)
				})
			// The reserved IP itself is the source of truth.
			mux.HandleFunc(fmt.Sprintf("/v2/reserved_ips/%s", ip),
				func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					fmt.Fprint(w, tt.reservedIPDOC)
				})

			server := httptest.NewServer(mux)
			defer server.Close()

			client := newTestGodoClient(t, server.URL)
			result, state, err := reservedIPAssignmentRefreshFunc(client, ip, actionID, tt.check)()
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}
			if state != tt.wantState {
				t.Errorf("state = %q, want %q", state, tt.wantState)
			}
			// A non-nil result is required so StateChangeConf does not count the
			// poll against NotFoundChecks while we wait in "in-progress".
			if result == nil {
				t.Errorf("result = nil, want non-nil reserved IP so the poll is not treated as not-found")
			}
		})
	}
}
