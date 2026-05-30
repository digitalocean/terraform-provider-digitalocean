package gradientai_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccDataSourceDigitalOceanCustomModels_List verifies that the list datasource
// can call the API without requiring a resource to be created first.
func TestAccDataSourceDigitalOceanCustomModels_List(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "digitalocean_gradientai_custom_models" "all" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_gradientai_custom_models.all", "custom_models.#"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanCustomModels_FilterByName(t *testing.T) {
	name := acceptance.RandomTestName() + "-cm"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceCustomModelsFilterConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_gradientai_custom_models.by_name", "custom_models.#", "1"),
					resource.TestCheckResourceAttr("data.digitalocean_gradientai_custom_models.by_name", "custom_models.0.name", name),
					resource.TestCheckResourceAttrSet("data.digitalocean_gradientai_custom_models.by_name", "custom_models.0.uuid"),
					resource.TestCheckResourceAttrSet("data.digitalocean_gradientai_custom_models.by_name", "custom_models.0.status"),
				),
			},
		},
	})
}

func testAccDataSourceCustomModelsFilterConfig(name string) string {
	return fmt.Sprintf(`
%s

data "digitalocean_gradientai_custom_models" "by_name" {
  filter {
    key    = "name"
    values = [%q]
  }

  depends_on = [digitalocean_gradientai_custom_model.test]
}
`, testAccCustomModelConfig_basic(name), name)
}
