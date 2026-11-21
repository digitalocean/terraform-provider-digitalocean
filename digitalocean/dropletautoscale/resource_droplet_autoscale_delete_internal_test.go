package dropletautoscale

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// After the delete request is accepted the API keeps returning the pool with
// status "deleting" for a while before it answers 404. The delete waiter must
// keep polling through that state instead of failing.
func TestDropletAutoscaleDeleteWaitsWhileDeleting(t *testing.T) {
	defer func(delay, minTimeout time.Duration) {
		dropletAutoscaleDeleteDelay = delay
		dropletAutoscaleDeleteMinTimeout = minTimeout
	}(dropletAutoscaleDeleteDelay, dropletAutoscaleDeleteMinTimeout)
	dropletAutoscaleDeleteDelay = 10 * time.Millisecond
	dropletAutoscaleDeleteMinTimeout = 10 * time.Millisecond

	poolID := "614fb0ef-b806-473c-bb85-a960a6d3d494"

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	var deleted bool
	mux.HandleFunc(fmt.Sprintf("/v2/droplets/autoscale/%s/dangerous", poolID), func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %v, expected %v", r.Method, http.MethodDelete)
		}
		deleted = true
		w.WriteHeader(http.StatusAccepted)
	})

	var gets int
	mux.HandleFunc(fmt.Sprintf("/v2/droplets/autoscale/%s", poolID), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		gets++
		if gets <= 2 {
			fmt.Fprintf(w, `{"autoscale_pool": {"id": %q, "name": "pool", "status": "deleting"}}`, poolID)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, `{"id": "not_found", "message": "autoscale group with id %s not found"}`, poolID)
	})

	cfg := &config.Config{
		Token:       "test-token",
		APIEndpoint: server.URL,
	}
	client, err := cfg.Client()
	if err != nil {
		t.Fatalf("error building client: %s", err)
	}

	d := schema.TestResourceDataRaw(t, ResourceDigitalOceanDropletAutoscale().Schema, map[string]interface{}{
		"name": "pool",
	})
	d.SetId(poolID)

	if diags := resourceDigitalOceanDropletAutoscaleDelete(context.Background(), d, client); diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if !deleted {
		t.Fatal("expected the delete request to be sent")
	}
	if gets < 3 {
		t.Fatalf("expected the waiter to poll until the pool was gone, got %d GETs", gets)
	}
	if d.Id() != "" {
		t.Fatalf("expected the ID to be cleared, got %q", d.Id())
	}
}
