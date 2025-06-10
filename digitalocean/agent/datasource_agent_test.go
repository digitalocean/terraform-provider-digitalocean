package agent

import (
	"context"
	"testing"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestDataSourceDigitalOceanAgent(t *testing.T) {
	resourceName := "data.digitalocean_agent.test"

	// Removed unused client variable

	resource.Test(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"digitalocean": {
				ResourcesMap: map[string]*schema.Resource{
					"digitalocean_agent": DataSourceDigitalOceanAgent(),
				},
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
					data "digitalocean_agent" "test" {
						agent_id = "test-agent-id"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "agent_id", "test-agent-id"),
					resource.TestCheckResourceAttr(resourceName, "name", "Test Agent"),
					resource.TestCheckResourceAttr(resourceName, "region", "nyc3"),
				),
			},
		},
	})
}

func TestDataSourceDigitalOceanAgentList(t *testing.T) {
	resourceName := "data.digitalocean_agent_list.test"

	// Removed unused meta variable

	resource.Test(t, resource.TestCase{
		Providers: map[string]*schema.Provider{
			"digitalocean": {
				ResourcesMap: map[string]*schema.Resource{
					"digitalocean_agent_list": DataSourceDigitalOceanAgentList(),
				},
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
					data "digitalocean_agent_list" "test" {}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "agents.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "agents.0.name", "Agent 1"),
					resource.TestCheckResourceAttr(resourceName, "agents.1.name", "Agent 2"),
				),
			},
		},
	})
}

func TestFlattenDigitalOceanAgent(t *testing.T) {
	agent := &godo.Agent{
		Uuid:        "test-agent-id",
		Name:        "Test Agent",
		Instruction: "Test Instruction",
		ProjectId:   "test-project-id",
		Region:      "nyc3",
		Description: "Test Description",
		CreatedAt:   &godo.Timestamp{Time: time.Now()},
		UpdatedAt:   &godo.Timestamp{Time: time.Now()},
		Tags:        []string{"tag1", "tag2"},
	}

	result, err := flattenDigitalOceanAgent(agent)
	assert.NoError(t, err)
	assert.Equal(t, "test-agent-id", result["agent_id"])
	assert.Equal(t, "Test Agent", result["name"])
	assert.Equal(t, "nyc3", result["region"])
	assert.Equal(t, "Test Description", result["description"])
}

type mockGodoClient struct{}

func (m *mockGodoClient) GenAI() *mockGenAIService {
	return &mockGenAIService{}
}

type mockGenAIService struct{}

func (m *mockGenAIService) GetAgent(ctx context.Context, agentID string) (*godo.Agent, *godo.Response, error) {
	return &godo.Agent{
		Uuid:        "test-agent-id",
		Name:        "Test Agent",
		Region:      "nyc3",
		Description: "Test Description",
	}, nil, nil
}

func (m *mockGenAIService) ListAgents(ctx context.Context, opts *godo.ListOptions) ([]godo.Agent, *godo.Response, error) {
	return []godo.Agent{
		{Uuid: "agent-1", Name: "Agent 1"},
		{Uuid: "agent-2", Name: "Agent 2"},
	}, nil, nil
}
