package genai_test

import (
	"context"
	"testing"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/genai" // Updated import path
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestDataSourceDigitalOceanAgent(t *testing.T) {
	resourceName := "data.digitalocean_agent.test"

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders(), // Use the testAccProviders function
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
					resource.TestCheckResourceAttr(resourceName, "description", "Test Description"), // Add description check
				),
			},
		},
	})
}

func TestDataSourceDigitalOceanAgents(t *testing.T) {
	resourceName := "data.digitalocean_agent_list.test"

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders(), // Use the testAccProviders function
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

	result, err := genai.FlattenDigitalOceanAgent(agent) // Use the method on the genai package
	assert.NoError(t, err)
	assert.Equal(t, "test-agent-id", result["agent_id"])
	assert.Equal(t, "Test Agent", result["name"])
	assert.Equal(t, "nyc3", result["region"])
	assert.Equal(t, "Test Description", result["description"])
}

// Define the testAccProviders function
func testAccProviders() map[string]*schema.Provider {
	return map[string]*schema.Provider{
		"digitalocean": testAccProvider(),
	}
}

// Define the testAccProvider function
func testAccProvider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{},
		ResourcesMap: map[string]*schema.Resource{
			"digitalocean_agent":      genai.ResourceDigitalOceanAgent(),
			"digitalocean_agent_list": genai.DataSourceDigitalOceanAgents(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"digitalocean_agent":      genai.DataSourceDigitalOceanAgent(),
			"digitalocean_agent_list": genai.DataSourceDigitalOceanAgents(),
		},
		ConfigureFunc: func(d *schema.ResourceData) (interface{}, error) {
			return &mockGodoClient{}, nil
		},
	}
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
