package genai_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanOpenAIApiKey_Basic(t *testing.T) {
	keyName := acceptance.RandomTestName() + "-openai-key"
	resourceConfig := fmt.Sprintf(`
resource "digitalocean_genai_openai_api_key" "test" {
  api_key = "sk-proj-testkey"
  name    = "%s"
}
`, keyName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_genai_openai_api_key.test", "name", keyName),
					resource.TestCheckResourceAttr("digitalocean_genai_openai_api_key.test", "api_key", "sk-proj-testkey"),
				),
			},
		},
	})
}

func TestAccDigitalOceanOpenAIApiKey_Update(t *testing.T) {
	keyName := acceptance.RandomTestName() + "-openai-key"
	updatedKeyName := keyName + "-updated"

	resourceConfig := fmt.Sprintf(`
resource "digitalocean_genai_openai_api_key" "test" {
  api_key = "sk-proj-testkey"
  name    = "%s"
}
`, keyName)

	updatedResourceConfig := fmt.Sprintf(`
resource "digitalocean_genai_openai_api_key" "test" {
  api_key = "sk-proj-testkey"
  name    = "%s"
}
`, updatedKeyName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_genai_openai_api_key.test", "name", keyName),
				),
			},
			{
				Config: updatedResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_genai_openai_api_key.test", "name", updatedKeyName),
				),
			},
		},
	})
}

func TestAccDigitalOceanOpenAIApiKey_Delete(t *testing.T) {
	keyName := acceptance.RandomTestName() + "-openai-key"
	resourceConfig := fmt.Sprintf(`
resource "digitalocean_genai_openai_api_key" "test" {
  api_key = "sk-proj-testkey"
  name    = "%s"
}
`, keyName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_genai_openai_api_key.test", "name", keyName),
				),
			},
			{
				ResourceName:      "digitalocean_genai_openai_api_key.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
