package vpcnatgateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	testNATGatewayID = "70e1b58d-cdec-4e95-b3ee-2d4d95feff51"
	testNATVPCUUID   = "0eb1752f-807b-4562-a077-8018e13ab1fb"
)

// An in-place update must send the configured VPC attachments, including
// default_gateway, instead of "vpcs": null (which clears the default route).
func TestVPCNATGatewayUpdateSendsVPCs(t *testing.T) {
	var (
		mu         sync.Mutex
		updateBody []byte
	)

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	mux.HandleFunc(fmt.Sprintf("/v2/vpc_nat_gateways/%s", testNATGatewayID), func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPut {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Errorf("error reading update request: %s", err)
			}
			mu.Lock()
			updateBody = body
			mu.Unlock()
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"vpc_nat_gateway": {
			"id": %q,
			"name": "example-egress",
			"type": "PUBLIC",
			"state": "ACTIVE",
			"region": "lon1",
			"size": 1,
			"vpcs": [{"vpc_uuid": %q, "gateway_ip": "10.20.0.5", "default_gateway": true}],
			"tcp_timeout_seconds": 300,
			"udp_timeout_seconds": 30,
			"icmp_timeout_seconds": 30,
			"created_at": "2026-10-23T10:00:00Z",
			"updated_at": "2026-10-23T10:05:00Z"
		}}`, testNATGatewayID, testNATVPCUUID)
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

	d := schema.TestResourceDataRaw(t, ResourceDigitalOceanVPCNATGateway().Schema, map[string]interface{}{
		"name":   "example-egress",
		"type":   "PUBLIC",
		"region": "lon1",
		"size":   1,
		"vpcs": []interface{}{
			map[string]interface{}{
				"vpc_uuid":        testNATVPCUUID,
				"default_gateway": true,
			},
		},
		"udp_timeout_seconds":  30,
		"icmp_timeout_seconds": 30,
		"tcp_timeout_seconds":  300,
	})
	d.SetId(testNATGatewayID)

	if diags := resourceDigitalOceanVPCNATGatewayUpdate(context.Background(), d, client); diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}

	mu.Lock()
	defer mu.Unlock()

	if updateBody == nil {
		t.Fatal("expected an update request")
	}

	var body struct {
		TCPTimeoutSeconds int               `json:"tcp_timeout_seconds"`
		VPCs              []json.RawMessage `json:"vpcs"`
	}
	if err := json.Unmarshal(updateBody, &body); err != nil {
		t.Fatalf("error decoding update request %s: %s", updateBody, err)
	}
	if body.TCPTimeoutSeconds != 300 {
		t.Errorf("tcp_timeout_seconds = %d, expected 300", body.TCPTimeoutSeconds)
	}
	if len(body.VPCs) != 1 {
		t.Fatalf("expected 1 VPC in update request, got %s", updateBody)
	}

	var vpc map[string]interface{}
	if err := json.Unmarshal(body.VPCs[0], &vpc); err != nil {
		t.Fatalf("error decoding VPC %s: %s", body.VPCs[0], err)
	}
	if vpc["vpc_uuid"] != testNATVPCUUID {
		t.Errorf("vpc_uuid = %v, expected %v", vpc["vpc_uuid"], testNATVPCUUID)
	}
	if vpc["default_gateway"] != true {
		t.Errorf("default_gateway = %v, expected true", vpc["default_gateway"])
	}
	if _, ok := vpc["gateway_ip"]; ok {
		t.Errorf("gateway_ip should not be sent on update, got %s", body.VPCs[0])
	}
}
