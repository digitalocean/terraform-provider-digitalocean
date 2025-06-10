package genai_test

import (
	"context"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/genai"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// fakeGenAI implements a fake version of the GenAI interface.
type fakeGenAI struct{}

func (f *fakeGenAI) CreateAgent(ctx context.Context, req *godo.AgentCreateRequest) (*godo.Agent, *godo.Response, error) {
	return &godo.Agent{
		Uuid:        "agent-123",
		Name:        req.Name,
		Instruction: req.Instruction,
		Model:       nil,
		Region:      req.Region,
		ProjectId:   req.ProjectId,
		Description: req.Description,
	}, &godo.Response{}, nil
}

func (f *fakeGenAI) GetAgent(ctx context.Context, uuid string) (*godo.Agent, *godo.Response, error) {
	return &godo.Agent{
		Uuid:        uuid,
		Name:        "test-agent",
		Instruction: "test-instruction",
		Model:       nil,
		Region:      "nyc1",
		ProjectId:   "proj-123",
		Description: "A test agent",
	}, &godo.Response{}, nil
}

func (f *fakeGenAI) UpdateAgent(ctx context.Context, id string, req *godo.AgentUpdateRequest) (*godo.Agent, *godo.Response, error) {
	return &godo.Agent{
		Uuid:        id,
		Name:        req.Name,
		Instruction: req.Instruction,
		Model:       nil,
		Region:      req.Region,
		ProjectId:   req.ProjectId,
	}, &godo.Response{}, nil
}

func (f *fakeGenAI) UpdateAgentVisibility(ctx context.Context, id string, req *godo.AgentVisibilityUpdateRequest) (*godo.Agent, *godo.Response, error) {
	return &godo.Agent{
		Uuid: id,
	}, &godo.Response{}, nil
}

func (f *fakeGenAI) DeleteAgent(ctx context.Context, id string) (*godo.Agent, *godo.Response, error) {
	return nil, &godo.Response{}, nil
}

// fakeGodoClient bundles the fake GenAI service.
type fakeGodoClient struct {
	GenAI *fakeGenAI
}

// fakeCombinedConfig satisfies the meta interface used in our resource.
type fakeCombinedConfig struct {
	client *fakeGodoClient
}

func (c *fakeCombinedConfig) GodoClient() *fakeGodoClient {
	return c.client
}

// getTestResourceData creates a ResourceData with default values for testing.
func getTestResourceData(t *testing.T) *schema.ResourceData {
	resourceSchema := genai.ResourceDigitalOceanAgent().Schema
	d := schema.TestResourceDataRaw(t, resourceSchema, map[string]interface{}{
		"name":        "test-agent",
		"instruction": "test instruction",
		"model_uuid":  "model-123",
		"region":      "nyc1",
		"project_id":  "proj-123",
		"description": "A test agent",
		"visibility":  "public",
	})
	return d
}

func TestResourceDigitalOceanAgentCRUD(t *testing.T) {
	ctx := context.Background()

	// Set up fake client and meta config.
	fakeClient := &fakeGodoClient{GenAI: &fakeGenAI{}}
	meta := &fakeCombinedConfig{client: fakeClient}

	d := getTestResourceData(t)

	// Test Create
	diags := genai.ResourceDigitalOceanAgent().CreateContext(ctx, d, meta)
	if diags.HasError() {
		t.Fatalf("Create failed: %v", diags)
	}
	if d.Id() != "agent-123" {
		t.Fatalf("Unexpected agent ID: got %s, want agent-123", d.Id())
	}

	// Test Read
	diags = genai.ResourceDigitalOceanAgent().ReadContext(ctx, d, meta)
	if diags.HasError() {
		t.Fatalf("Read failed: %v", diags)
	}

	// Test Update - simulate a change in the agent's name and visibility.
	if err := d.Set("name", "updated-agent"); err != nil {
		t.Fatalf("Failed to set name: %v", err)
	}
	if err := d.Set("visibility", "private"); err != nil {
		t.Fatalf("Failed to set visibility: %v", err)
	}
	diags = genai.ResourceDigitalOceanAgent().UpdateContext(ctx, d, meta)
	if diags.HasError() {
		t.Fatalf("Update failed: %v", diags)
	}

	// Test Delete
	diags = genai.ResourceDigitalOceanAgent().DeleteContext(ctx, d, meta)
	if diags.HasError() {
		t.Fatalf("Delete failed: %v", diags)
	}
	if d.Id() != "" {
		t.Fatalf("Agent not deleted properly, id still set: %s", d.Id())
	}
}
