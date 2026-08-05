package gradientai_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// TestAccDigitalOceanAgentGuardrailAttachment_Basic attaches an existing guardrail to a
// freshly created agent. Guardrails are account-specific, so the guardrail UUID is read
// from the GENAI_GUARDRAIL_UUID environment variable and the test is skipped when unset.
func TestAccDigitalOceanAgentGuardrailAttachment_Basic(t *testing.T) {
	guardrailUUID := os.Getenv("GENAI_GUARDRAIL_UUID")
	agentName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			if guardrailUUID == "" {
				t.Skip("GENAI_GUARDRAIL_UUID must be set for the guardrail attachment acceptance test")
			}
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanAgentGuardrailAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanAgentGuardrailAttachmentConfig(agentName, guardrailUUID, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAgentGuardrailAttachmentExists("digitalocean_gradientai_agent_guardrail_attachment.test"),
					resource.TestCheckResourceAttrSet("digitalocean_gradientai_agent_guardrail_attachment.test", "agent_uuid"),
					resource.TestCheckResourceAttr("digitalocean_gradientai_agent_guardrail_attachment.test", "guardrail_uuid", guardrailUUID),
					resource.TestCheckResourceAttr("digitalocean_gradientai_agent_guardrail_attachment.test", "priority", "1"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanAgentGuardrailAttachmentExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No attachment ID is set")
		}

		agentUUID := rs.Primary.Attributes["agent_uuid"]
		guardrailUUID := rs.Primary.Attributes["guardrail_uuid"]

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		agent, _, err := client.GradientAI.GetAgent(context.Background(), agentUUID)
		if err != nil {
			return err
		}
		for _, g := range agent.Guardrails {
			if g != nil && (g.GuardrailUuid == guardrailUUID || g.Uuid == guardrailUUID) {
				return nil
			}
		}
		return fmt.Errorf("guardrail %s not attached to agent %s", guardrailUUID, agentUUID)
	}
}

func testAccCheckDigitalOceanAgentGuardrailAttachmentDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_gradientai_agent_guardrail_attachment" {
			continue
		}

		agentUUID := rs.Primary.Attributes["agent_uuid"]
		guardrailUUID := rs.Primary.Attributes["guardrail_uuid"]

		agent, _, err := client.GradientAI.GetAgent(context.Background(), agentUUID)
		if err != nil {
			// The agent is likely gone too, which means the attachment is gone.
			continue
		}
		for _, g := range agent.Guardrails {
			if g != nil && (g.GuardrailUuid == guardrailUUID || g.Uuid == guardrailUUID) {
				return fmt.Errorf("guardrail attachment still exists: %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckDigitalOceanAgentGuardrailAttachmentConfig(name, guardrailUUID string, priority int) string {
	return fmt.Sprintf(`
resource "digitalocean_gradientai_agent" "test" {
  name        = "%[1]s"
  instruction = "You are a helpful AI assistant with a guardrail."
  model_uuid  = "%[2]s"
  project_id  = "%[3]s"
  region      = "tor1"
}

resource "digitalocean_gradientai_agent_guardrail_attachment" "test" {
  agent_uuid     = digitalocean_gradientai_agent.test.id
  guardrail_uuid = "%[4]s"
  priority       = %[5]d
}
`, name, defaultModelUUID, defaultProjecID, guardrailUUID, priority)
}
