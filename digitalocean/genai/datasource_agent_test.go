package genai_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceDigitalOceanAgent_BasicByAgentID(t *testing.T) {
	var agent godo.Agent
	agentName := acceptance.RandomTestName()
	resourceConfig := testAccCheckDataSourceDigitalOceanAgentConfig_basic(agentName)
	dataSourceConfig := `
data "digitalocean_agent" "foobar" {
  agent_id = digitalocean_agent.foo.agent_id
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanAgentExists("data.digitalocean_agent.foobar", &agent),
					resource.TestCheckResourceAttr(
						"data.digitalocean_agent.foobar", "name", agentName),
					resource.TestCheckResourceAttrSet("data.digitalocean_agent.foobar", "agent_id"),
					resource.TestCheckResourceAttrSet("data.digitalocean_agent.foobar", "uuid"),
					resource.TestCheckResourceAttrSet("data.digitalocean_agent.foobar", "created_at"),
					resource.TestCheckResourceAttrSet("data.digitalocean_agent.foobar", "updated_at"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanAgent_CompleteConfiguration(t *testing.T) {
	var agent godo.Agent
	agentName := acceptance.RandomTestName()
	description := "Test GenAI agent for acceptance testing"
	prompt := "You are a helpful AI assistant for testing purposes."
	model := "gpt-4"

	resourceConfig := testAccCheckDataSourceDigitalOceanAgentConfig_complete(agentName, description, prompt, model)
	dataSourceConfig := `
data "digitalocean_agent" "foobar" {
  agent_id = digitalocean_agent.foo.agent_id
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanAgentExists("data.digitalocean_agent.foobar", &agent),
					resource.TestCheckResourceAttr(
						"data.digitalocean_agent.foobar", "name", agentName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_agent.foobar", "description", description),
					resource.TestCheckResourceAttr(
						"data.digitalocean_agent.foobar", "prompt", prompt),
					resource.TestCheckResourceAttr(
						"data.digitalocean_agent.foobar", "model", model),
					resource.TestCheckResourceAttrSet("data.digitalocean_agent.foobar", "agent_id"),
					resource.TestCheckResourceAttrSet("data.digitalocean_agent.foobar", "uuid"),
					resource.TestCheckResourceAttrSet("data.digitalocean_agent.foobar", "created_at"),
					resource.TestCheckResourceAttrSet("data.digitalocean_agent.foobar", "updated_at"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanAgent_WithTags(t *testing.T) {
	var agent godo.Agent
	agentName := acceptance.RandomTestName()
	tagName := acceptance.RandomTestName("tag")

	resourceConfig := testAccCheckDataSourceDigitalOceanAgentConfig_withTags(agentName, tagName)
	dataSourceConfig := `
data "digitalocean_agent" "foobar" {
  agent_id = digitalocean_agent.foo.agent_id
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanAgentExists("data.digitalocean_agent.foobar", &agent),
					resource.TestCheckResourceAttr(
						"data.digitalocean_agent.foobar", "name", agentName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_agent.foobar", "tags.#", "1"),
					resource.TestCheckResourceAttrSet("data.digitalocean_agent.foobar", "agent_id"),
					resource.TestCheckResourceAttrSet("data.digitalocean_agent.foobar", "uuid"),
					resource.TestCheckResourceAttrSet("data.digitalocean_agent.foobar", "created_at"),
					resource.TestCheckResourceAttrSet("data.digitalocean_agent.foobar", "updated_at"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanAgent_NonExistentAgent(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckDataSourceDigitalOceanAgentConfig_nonExistent(),
				ExpectError: regexp.MustCompile("not found"),
			},
		},
	})
}

// Helper function to check if the agent exists in the data source state
func testAccCheckDataSourceDigitalOceanAgentExists(n string, agent *godo.Agent) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No agent ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		agentID := rs.Primary.Attributes["agent_id"]
		if agentID == "" {
			return fmt.Errorf("No agent_id is set")
		}

		foundAgent, _, err := client.GenAI.GetAgent(context.Background(), agentID)
		if err != nil {
			return err
		}

		if foundAgent.Uuid != rs.Primary.ID {
			return fmt.Errorf("Agent not found: expected UUID %s, got %s", rs.Primary.ID, foundAgent.Uuid)
		}

		*agent = *foundAgent

		return nil
	}
}

// Configuration templates for different test scenarios

func testAccCheckDataSourceDigitalOceanAgentConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_agent" "foo" {
  name        = "%s"
  description = "Basic test agent"
  prompt      = "You are a test assistant."
  model       = "gpt-3.5-turbo"
}`, name)
}

func testAccCheckDataSourceDigitalOceanAgentConfig_complete(name, description, prompt, model string) string {
	return fmt.Sprintf(`
resource "digitalocean_agent" "foo" {
  name        = "%s"
  description = "%s"
  prompt      = "%s"
  model       = "%s"
  
  # Additional configuration options if supported
  temperature = 0.7
  max_tokens  = 1000
  
  # Enable/disable features if supported
  web_search_enabled = true
  code_interpreter_enabled = false
}`, name, description, prompt, model)
}

func testAccCheckDataSourceDigitalOceanAgentConfig_withTags(name, tagName string) string {
	return fmt.Sprintf(`
resource "digitalocean_tag" "foo" {
  name = "%s"
}

resource "digitalocean_agent" "foo" {
  name        = "%s"
  description = "Test agent with tags"
  prompt      = "You are a tagged test assistant."
  model       = "gpt-3.5-turbo"
  tags        = [digitalocean_tag.foo.id]
}`, tagName, name)
}

func testAccCheckDataSourceDigitalOceanAgentConfig_nonExistent() string {
	return `
data "digitalocean_agent" "foobar" {
  agent_id = "non-existent-agent-id-12345"
}`
}

// Additional test for testing agent configurations if your schema supports them
func TestAccDataSourceDigitalOceanAgent_WithKnowledgeBase(t *testing.T) {
	var agent godo.Agent
	agentName := acceptance.RandomTestName()

	// This test assumes your agent supports knowledge base configuration
	resourceConfig := testAccCheckDataSourceDigitalOceanAgentConfig_withKnowledgeBase(agentName)
	dataSourceConfig := `
data "digitalocean_agent" "foobar" {
  agent_id = digitalocean_agent.foo.agent_id
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanAgentExists("data.digitalocean_agent.foobar", &agent),
					resource.TestCheckResourceAttr(
						"data.digitalocean_agent.foobar", "name", agentName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_agent.foobar", "knowledge_base.#", "1"),
					resource.TestCheckResourceAttrSet("data.digitalocean_agent.foobar", "agent_id"),
					resource.TestCheckResourceAttrSet("data.digitalocean_agent.foobar", "uuid"),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanAgentConfig_withKnowledgeBase(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_agent" "foo" {
  name        = "%s"
  description = "Test agent with knowledge base"
  prompt      = "You are an assistant with access to a knowledge base."
  model       = "gpt-4"
  
  knowledge_base {
    name        = "test-kb"
    description = "Test knowledge base"
    # Add other knowledge base configuration as needed
  }
}`, name)
}
