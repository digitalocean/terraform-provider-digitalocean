package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceDigitalOceanContainerRegistry_Basic(t *testing.T) {
	var reg godo.Registry
	regName := fmt.Sprintf("foo-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanContainerRegistryConfig_basic, regName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanContainerRegistryExists("data.digitalocean_container_registry.foobar", &reg),
					resource.TestCheckResourceAttr(
						"data.digitalocean_container_registry.foobar", "name", regName),
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

const testAccCheckDataSourceDigitalOceanContainerRegistryConfig_basic = `
resource "digitalocean_container_registry" "foo" {
  name = "%s"
}

data "digitalocean_container_registry" "foobar" {
  name = "${digitalocean_container_registry.foo.name}"
}`
