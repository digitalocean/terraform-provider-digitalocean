package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceDigitalOceanContainerRegistry_Basic(t *testing.T) {
	var reg godo.Registry
	regName := randomTestName()

	resourceConfig := fmt.Sprintf(`
resource "digitalocean_container_registry" "foo" {
  name                   = "%s"
  subscription_tier_slug = "basic"
}
`, regName)

	dataSourceConfig := `
data "digitalocean_container_registry" "foobar" {
  name = digitalocean_container_registry.foo.name
}
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanContainerRegistryExists("data.digitalocean_container_registry.foobar", &reg),
					resource.TestCheckResourceAttr(
						"data.digitalocean_container_registry.foobar", "name", regName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_container_registry.foobar", "subscription_tier_slug", "basic"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_container_registry.foobar", "region"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_container_registry.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_container_registry.foobar", "storage_usage_bytes"),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanContainerRegistryExists(n string, reg *godo.Registry) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No registry ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		foundReg, _, err := client.Registry.Get(context.Background())

		if err != nil {
			return err
		}

		if foundReg.Name != rs.Primary.ID {
			return fmt.Errorf("Registry not found")
		}

		*reg = *foundReg

		return nil
	}
}
