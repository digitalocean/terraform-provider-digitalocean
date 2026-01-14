package gradientai_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanOpenAIApiKeys_ListAll(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "digitalocean_gradientai_openai_api_keys" "all" {}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_gradientai_openai_api_keys.all", "openai_api_keys.#"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanAgentsByOpenAIApiKey_ListAgents(t *testing.T) {
	keyName := acceptance.RandomTestName() + "-openai-key"
	resourceConfig := fmt.Sprintf(`
resource "digitalocean_gradientai_openai_api_key" "test" {
  api_key = "sk-proj-testkey"
  name    = "%s"
}

data "digitalocean_gradientai_agents_by_openai_api_key" "by_key" {
  uuid = digitalocean_gradientai_openai_api_key.test.uuid
}
`, keyName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_gradientai_agents_by_openai_api_key.by_key", "agents"),
				),
			},
		},
	})
}
