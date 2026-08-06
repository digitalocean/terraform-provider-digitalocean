package database

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Kafka topic creation is asynchronous: when the topic is not provisioned
// in time, the API responds 202 with an acknowledgment body instead of the
// documented 201 with the topic, and godo yields a nil topic with no
// error. Create must not depend on the response payload, and must wait
// until the topic is readable.
func TestDatabaseKafkaTopicCreateNoTopicInResponse(t *testing.T) {
	clusterID := "3f2549b8-9257-4e2f-a04c-e23547b2d685"
	topicName := "topic-foobar"

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	defer server.Close()

	mux.HandleFunc(fmt.Sprintf("/v2/databases/%s/topics", clusterID), func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %v, expected %v", r.Method, http.MethodPost)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprint(w, `{"message":"topic cannot be performed. Reason: topic is still provisioning","id":"accepted"}`)
	})

	var gets int
	mux.HandleFunc(fmt.Sprintf("/v2/databases/%s/topics/%s", clusterID, topicName), func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		gets++
		if gets == 1 {
			// Still provisioning: a payload without the topic key.
			fmt.Fprint(w, `{"id":"accepted"}`)
			return
		}
		fmt.Fprintf(w, `{"topic": {"name": %q, "state": "active", "replication_factor": 2}}`, topicName)
	})

	cfg := &config.Config{
		Token:       "test-token",
		APIEndpoint: server.URL,
	}
	client, err := cfg.Client()
	if err != nil {
		t.Fatalf("error building client: %s", err)
	}

	d := schema.TestResourceDataRaw(t, ResourceDigitalOceanDatabaseKafkaTopic().Schema, map[string]interface{}{
		"cluster_id":         clusterID,
		"name":               topicName,
		"partition_count":    3,
		"replication_factor": 2,
	})

	if diags := resourceDigitalOceanDatabaseKafkaTopicCreate(context.Background(), d, client); diags.HasError() {
		t.Fatalf("create returned error: %v", diags)
	}

	if want := makeKafkaTopicID(clusterID, topicName); d.Id() != want {
		t.Errorf("id = %q, expected %q", d.Id(), want)
	}
	if state := d.Get("state").(string); state != "active" {
		t.Errorf("state = %q, expected %q", state, "active")
	}
}
