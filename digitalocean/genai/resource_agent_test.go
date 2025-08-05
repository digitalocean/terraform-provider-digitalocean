package genai_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanAgent_Basic(t *testing.T) {
	var agent godo.Agent
	agentName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanAgentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanAgentConfig_basic(agentName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAgentExists("digitalocean_genai_agent.test", &agent),
					resource.TestCheckResourceAttr("digitalocean_genai_agent.test", "name", agentName),
					resource.TestCheckResourceAttr("digitalocean_genai_agent.test", "instruction", "You are a helpful AI assistant."),
					resource.TestCheckResourceAttrSet("digitalocean_genai_agent.test", "model_uuid"),
					resource.TestCheckResourceAttrSet("digitalocean_genai_agent.test", "project_id"),
					resource.TestCheckResourceAttrSet("digitalocean_genai_agent.test", "region"),
					resource.TestCheckResourceAttrSet("digitalocean_genai_agent.test", "created_at"),
					resource.TestCheckResourceAttrSet("digitalocean_genai_agent.test", "updated_at"),
				),
			},
		},
	})
}

func TestAccDigitalOceanAgent_WithOptionalFields(t *testing.T) {
	var agent godo.Agent
	agentName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanAgentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanAgentConfig_withOptionalFields(agentName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAgentExists("digitalocean_genai_agent.test", &agent),
					resource.TestCheckResourceAttr("digitalocean_genai_agent.test", "name", agentName),
					resource.TestCheckResourceAttr("digitalocean_genai_agent.test", "description", "Test agent with optional fields"),
					resource.TestCheckResourceAttr("digitalocean_genai_agent.test", "tags.#", "2"),
					resource.TestCheckResourceAttr("digitalocean_genai_agent.test", "tags.0", "test"),
					resource.TestCheckResourceAttr("digitalocean_genai_agent.test", "tags.1", "ai"),
				),
			},
		},
	})
}

func TestAccDigitalOceanAgent_Update(t *testing.T) {
	var agent godo.Agent
	agentName := acceptance.RandomTestName()
	updatedName := agentName + "-updated"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanAgentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanAgentConfig_basic(agentName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAgentExists("digitalocean_genai_agent.test", &agent),
					resource.TestCheckResourceAttr("digitalocean_genai_agent.test", "name", agentName),
					resource.TestCheckResourceAttr("digitalocean_genai_agent.test", "instruction", "You are a helpful AI assistant."),
				),
			},
			{
				Config: testAccCheckDigitalOceanAgentConfig_updated(updatedName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAgentExists("digitalocean_genai_agent.test", &agent),
					resource.TestCheckResourceAttr("digitalocean_genai_agent.test", "name", updatedName),
					resource.TestCheckResourceAttr("digitalocean_genai_agent.test", "instruction", "You are an updated AI assistant with new capabilities."),
					resource.TestCheckResourceAttr("digitalocean_genai_agent.test", "description", "Updated test agent"),
				),
			},
		},
	})
}

func TestAccDigitalOceanAgent_WithKnowledgeBase(t *testing.T) {
	var agent godo.Agent
	agentName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanAgentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanAgentConfig_withKnowledgeBase(agentName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAgentExists("digitalocean_genai_agent.test", &agent),
					resource.TestCheckResourceAttr("digitalocean_genai_agent.test", "name", agentName),
					resource.TestCheckResourceAttr("digitalocean_genai_agent.test", "knowledge_base_uuid.#", "1"),
					resource.TestCheckResourceAttrSet("digitalocean_genai_agent.test", "knowledge_base_uuid.0"),
				),
			},
		},
	})
}

func TestAccDigitalOceanAgent_WithDeployment(t *testing.T) {
	var agent godo.Agent
	agentName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanAgentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanAgentConfig_withDeployment(agentName, "private"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAgentExists("digitalocean_genai_agent.test", &agent),
					resource.TestCheckResourceAttr("digitalocean_genai_agent.test", "name", agentName),
					resource.TestCheckResourceAttr("digitalocean_genai_agent.test", "deployment.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_genai_agent.test", "deployment.0.visibility", "private"),
				),
			},
			{
				Config: testAccCheckDigitalOceanAgentConfig_withDeployment(agentName, "public"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAgentExists("digitalocean_genai_agent.test", &agent),
					resource.TestCheckResourceAttr("digitalocean_genai_agent.test", "deployment.0.visibility", "public"),
				),
			},
		},
	})
}

func TestAccDigitalOceanAgent_ImportBasic(t *testing.T) {
	agentName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanAgentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanAgentConfig_basic(agentName),
			},
			{
				ResourceName:      "digitalocean_genai_agent.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDigitalOceanAgentDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_genai_agent" {
			continue
		}

		_, _, err := client.GenAI.GetAgent(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Agent still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckDigitalOceanAgentExists(resource string, agent *godo.Agent) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Agent ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		foundAgent, _, err := client.GenAI.GetAgent(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundAgent.Uuid != rs.Primary.ID {
			return fmt.Errorf("Agent not found")
		}

		*agent = *foundAgent
		return nil
	}
}

func testAccCheckDigitalOceanAgentConfig_basic(name string) string {
	return fmt.Sprintf(`

resource "digitalocean_genai_agent" "test" {
  name        = "%s"
  instruction = "You are a helpful AI assistant."
  model_uuid  = "%s"
  project_id  = "%s"
  region      = "tor1"
}`, name, defaultModelUUID, defaultProjecID)
}

func testAccCheckDigitalOceanAgentConfig_withOptionalFields(name string) string {
	return fmt.Sprintf(`

resource "digitalocean_genai_agent" "test" {
  name        = "%s"
  instruction = "You are a helpful AI assistant."
  description = "Test agent with optional fields"
  model_uuid  = "%s"
  project_id  = "%s"
  region      = "tor1"
}`, name, defaultModelUUID, defaultProjecID)
}

func testAccCheckDigitalOceanAgentConfig_updated(name string) string {
	return fmt.Sprintf(`

resource "digitalocean_genai_agent" "test" {
  name        = "%s"
  instruction = "You are an updated AI assistant with new capabilities."
  description = "Updated test agent"
  model_uuid  = "%s"
  project_id  = "%s"
  region      = "tor1"

}`, name, defaultModelUUID, defaultProjecID)
}

func testAccCheckDigitalOceanAgentConfig_withKnowledgeBase(name string) string {
	return fmt.Sprintf(`

resource "digitalocean_knowledge_base" "test" {
  name        = "%s-kb"
  description = "Test knowledge base"
  project_id  = digitalocean_project.test.id
  region      = "tor1"
}

resource "digitalocean_genai_agent" "test" {
  name                = "%s"
  instruction         = "You are a helpful AI assistant with knowledge base access."
  model_uuid          = "%s"
  project_id          = "%s"
  region              = "tor1"
  knowledge_base_uuid = [digitalocean_knowledge_base.test.id]
}`, name, name, defaultModelUUID, defaultProjecID)
}

func testAccCheckDigitalOceanAgentConfig_withDeployment(name, visibility string) string {
	return fmt.Sprintf(`

resource "digitalocean_genai_agent" "test" {
  name        = "%s"
  instruction = "You are a helpful AI assistant."
  model_uuid  = "%s"
  project_id  = "%s"
  region      = "tor1"

  deployment {
    visibility = "%s"
  }
}`, name, defaultModelUUID, defaultProjecID, visibility)
}
