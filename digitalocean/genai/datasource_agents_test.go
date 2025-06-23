package genai_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	defaultDescription = "Agent will be used to help users find the perfect solution for their technical needs."
	defaultInstruction = "You are DigitalOcean's Solutions Architect Assistant, designed to help users find the perfect solution for their technical needs. Your primary role is to understand user requirements and recommend appropriate solutions from the DigitalOcean Marketplace and product portfolio."
	defaultModelUUID   = "d111f1d1-d1f1-11ef-bf1f-1e111e1ggge1"
	defaultProjecID    = "11e1e111-ee11-11ac-11ff-1111cf1111e1"
)

func TestAccDataSourceDigitalOceanAgents_Basic(t *testing.T) {
	name1 := acceptance.RandomTestName() + "-agent"
	name2 := acceptance.RandomTestName() + "-agent"

	resourcesConfig := fmt.Sprintf(`
resource "digitalocean_genai_agent" "foo" {
  name        = "%s"
  description = "%s"
  instruction = "%s"
  model_uuid  = "%s"
  project_id  = "%s"
  region      = "tor1"

}

resource "digitalocean_genai_agent" "bar" {
  name        = "%s"
  description = "%s"
  instruction = "%s"
  model_uuid  = "%s"
  project_id  = "%s"
  region      = "tor1"
}
`, name1, defaultDescription, defaultInstruction, defaultModelUUID, defaultProjecID, name2, defaultDescription, defaultInstruction, defaultModelUUID, defaultProjecID)

	datasourceConfig := fmt.Sprintf(`
data "digitalocean_genai_agents" "result" {
  filter {
    key    = "name"
    values = ["%s"]
  }
}
`, name1)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_genai_agents.result", "agents.#", "1"),
					resource.TestCheckResourceAttr("data.digitalocean_genai_agents.result", "agents.0.name", name1),
					resource.TestCheckResourceAttrPair("data.digitalocean_genai_agents.result", "agents.0.id", "digitalocean_domain.foo", "id")),
			},
			{
				Config: resourcesConfig,
			},
		},
	})
}
