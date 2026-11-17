package functions_test

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

func TestAccDigitalOceanFunctionsNamespace_Basic(t *testing.T) {
	var ns godo.FunctionsNamespace
	label := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanFunctionsNamespaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanFunctionsNamespaceConfig(label, "nyc1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanFunctionsNamespaceExists("digitalocean_functions_namespace.test", &ns),
					resource.TestCheckResourceAttr("digitalocean_functions_namespace.test", "label", label),
					resource.TestCheckResourceAttr("digitalocean_functions_namespace.test", "region", "nyc1"),
					resource.TestCheckResourceAttrSet("digitalocean_functions_namespace.test", "uuid"),
					resource.TestCheckResourceAttrSet("digitalocean_functions_namespace.test", "api_host"),
					resource.TestCheckResourceAttrSet("digitalocean_functions_namespace.test", "created_at"),
				),
			},
			{
				ResourceName:      "digitalocean_functions_namespace.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckDigitalOceanFunctionsNamespaceDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_functions_namespace" {
			continue
		}

		_, _, err := client.Functions.GetNamespace(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Functions namespace %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckDigitalOceanFunctionsNamespaceExists(n string, ns *godo.FunctionsNamespace) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No Functions namespace ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		found, _, err := client.Functions.GetNamespace(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}
		if found.Namespace != rs.Primary.ID {
			return fmt.Errorf("Functions namespace not found")
		}

		*ns = *found
		return nil
	}
}

func testAccCheckDigitalOceanFunctionsNamespaceConfig(label, region string) string {
	return fmt.Sprintf(`
resource "digitalocean_functions_namespace" "test" {
  label  = "%s"
  region = "%s"
}
`, label, region)
}
