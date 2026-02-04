package gradientai_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanOpenAIApiKey_ByID(t *testing.T) {
	keyName := acceptance.RandomTestName() + "-openai-key"
	resourceConfig := fmt.Sprintf(`
resource "digitalocean_gradientai_openai_api_key" "test" {
  api_key = "sk-proj-testkey"
  name    = "%s"
}

data "digitalocean_gradientai_openai_api_key" "by_id" {
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
					resource.TestCheckResourceAttr("data.digitalocean_gradientai_openai_api_key.by_id", "name", keyName),
					resource.TestCheckResourceAttrSet("data.digitalocean_gradientai_openai_api_key.by_id", "uuid"),
				),
			},
		},
	})
}
