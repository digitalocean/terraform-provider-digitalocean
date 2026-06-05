package gradientai_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanCustomModel_Basic(t *testing.T) {
	name := acceptance.RandomTestName() + "-cm"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCustomModelConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_gradientai_custom_model.test", "name", name),
					resource.TestCheckResourceAttrPair(
						"data.digitalocean_gradientai_custom_model.test", "uuid",
						"digitalocean_gradientai_custom_model.test", "uuid",
					),
					resource.TestCheckResourceAttrPair(
						"data.digitalocean_gradientai_custom_model.test", "name",
						"digitalocean_gradientai_custom_model.test", "name",
					),
					resource.TestCheckResourceAttrPair(
						"data.digitalocean_gradientai_custom_model.test", "status",
						"digitalocean_gradientai_custom_model.test", "status",
					),
					resource.TestCheckResourceAttr("data.digitalocean_gradientai_custom_model.test", "source_ref.0.repo_id", testCustomModelRepo()),
					resource.TestCheckResourceAttrSet("data.digitalocean_gradientai_custom_model.test", "created_at"),
				),
			},
		},
	})
}

func testAccDataSourceCustomModelConfig(name string) string {
	return fmt.Sprintf(`
%s

data "digitalocean_gradientai_custom_model" "test" {
  uuid = digitalocean_gradientai_custom_model.test.uuid
}
`, testAccCustomModelConfig_basic(name))
}
