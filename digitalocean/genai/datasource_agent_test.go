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
data "digitalocean_genai_agent" "foobar" {
  agent_id = digitalocean_genai_agent.foo.agent_id
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
					testAccCheckDataSourceDigitalOceanAgentExists("data.digitalocean_genai_agent.foobar", &agent),
					resource.TestCheckResourceAttr(
						"data.digitalocean_genai_agent.foobar", "name", agentName),
					resource.TestCheckResourceAttrSet("data.digitalocean_genai_agent.foobar", "agent_id"),
					resource.TestCheckResourceAttrSet("data.digitalocean_genai_agent.foobar", "uuid"),
					resource.TestCheckResourceAttrSet("data.digitalocean_genai_agent.foobar", "created_at"),
					resource.TestCheckResourceAttrSet("data.digitalocean_genai_agent.foobar", "updated_at"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanAgent_CompleteConfiguration(t *testing.T) {
	var agent godo.Agent
	name := acceptance.RandomTestName()
	description := "Test GenAI agent for acceptance testing"
	instruction := "You are a helpful AI assistant for testing purposes."
	model_uuid := defaultModelUUID
	project_id := defaultProjecID
	region := "tor1"

	resourceConfig := testAccCheckDataSourceDigitalOceanAgentConfig_complete(name, description, instruction, model_uuid, project_id, region)
	dataSourceConfig := `
data "digitalocean_genai_agent" "foobar" {
  agent_id = digitalocean_genai_agent.foo.agent_id
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
					testAccCheckDataSourceDigitalOceanAgentExists("data.digitalocean_genai_agent.foobar", &agent),
					resource.TestCheckResourceAttr(
						"data.digitalocean_genai_agent.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"data.digitalocean_genai_agent.foobar", "description", description),
					resource.TestCheckResourceAttr(
						"data.digitalocean_genai_agent.foobar", "instruction", instruction),
					resource.TestCheckResourceAttr(
						"data.digitalocean_genai_agent.foobar", "model_uuid", model_uuid),
					resource.TestCheckResourceAttrSet("data.digitalocean_genai_agent.foobar", "agent_id"),
					resource.TestCheckResourceAttrSet("data.digitalocean_genai_agent.foobar", "uuid"),
					resource.TestCheckResourceAttrSet("data.digitalocean_genai_agent.foobar", "created_at"),
					resource.TestCheckResourceAttrSet("data.digitalocean_genai_agent.foobar", "updated_at"),
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
data "digitalocean_genai_agent" "foobar" {
  agent_id = digitalocean_genai_agent.foo.agent_id
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
					testAccCheckDataSourceDigitalOceanAgentExists("data.digitalocean_genai_agent.foobar", &agent),
					resource.TestCheckResourceAttr(
						"data.digitalocean_genai_agent.foobar", "name", agentName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_genai_agent.foobar", "tags.#", "1"),
					resource.TestCheckResourceAttrSet("data.digitalocean_genai_agent.foobar", "agent_id"),
					resource.TestCheckResourceAttrSet("data.digitalocean_genai_agent.foobar", "uuid"),
					resource.TestCheckResourceAttrSet("data.digitalocean_genai_agent.foobar", "created_at"),
					resource.TestCheckResourceAttrSet("data.digitalocean_genai_agent.foobar", "updated_at"),
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

func testAccCheckDataSourceDigitalOceanAgentConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_genai_agent" "foo" {
  name        = "%s"
  description = "Basic test agent"
  instruction = "You are a test assistant."
  model_uuid  = "%s"
}`, name, defaultModelUUID)
}

func testAccCheckDataSourceDigitalOceanAgentConfig_complete(name, description, instruction, model_uuid, project_id, region string) string {
	return fmt.Sprintf(`
resource "digitalocean_genai_agent" "foo" {
  name        = "%s"
  instruction = "%s"
  description = "%s"
  model_uuid  = "%s"
  project_id  = "%s"
  region      = "tor1"

}`, name, instruction, description, model_uuid, project_id)
}

func testAccCheckDataSourceDigitalOceanAgentConfig_withTags(name, tagName string) string {
	return fmt.Sprintf(`
resource "digitalocean_genai_agent" "foo" {
  name        = "%s"
  instruction = "You are a tagged test assistant."
  description = "Test agent with tags"
  model_uuid  = "%s"
  project_id  = "%s"

  tags = [digitalocean_tag.foo.id]
}`, name, defaultProjecID, defaultModelUUID)
}

func testAccCheckDataSourceDigitalOceanAgentConfig_nonExistent() string {
	return `
data "digitalocean_genai_agent" "foobar" {
  agent_id = "non-existent-agent-id-12345"
}`
}

func TestAccDataSourceDigitalOceanAgent_WithKnowledgeBase(t *testing.T) {
	var agent godo.Agent
	agentName := acceptance.RandomTestName()
	resourceConfig := testAccCheckDataSourceDigitalOceanAgentConfig_withKnowledgeBase(agentName)
	dataSourceConfig := `
data "digitalocean_genai_agent" "foobar" {
  agent_id = digitalocean_genai_agent.foo.agent_id
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
					testAccCheckDataSourceDigitalOceanAgentExists("data.digitalocean_genai_agent.foobar", &agent),
					resource.TestCheckResourceAttr(
						"data.digitalocean_genai_agent.foobar", "name", agentName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_genai_agent.foobar", "knowledge_base.#", "1"),
					resource.TestCheckResourceAttrSet("data.digitalocean_genai_agent.foobar", "agent_id"),
					resource.TestCheckResourceAttrSet("data.digitalocean_genai_agent.foobar", "uuid"),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanAgentConfig_withKnowledgeBase(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_genai_agent" "foo" {
  name        = "%s"
  description = "Test agent with knowledge base"
  instruction = "You are an assistant with access to a knowledge base."
  model_uuid  = "%s"

  knowledge_base {
    name        = "test-kb"
    description = "Test knowledge base"
  }
}`, name, defaultModelUUID)
}
