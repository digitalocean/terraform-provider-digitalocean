package actiongateway_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanActionGatewayToolbelt_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceName := "digitalocean_action_gateway_toolbelt.foobar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckActionGatewayToolbeltDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccActionGatewayToolbeltConfig, name, `"exa_web_search"`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "display_name", "Search Tools"),
					resource.TestCheckResourceAttr(resourceName, "tools.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tools.*", "exa_web_search"),
					resource.TestCheckResourceAttr(resourceName, "pinned_tools.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "version", "1"),
					resource.TestCheckResourceAttr(resourceName, "reference", name+"@1"),
					resource.TestCheckResourceAttr(resourceName, "tool_count", "1"),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
				),
			},
			{
				Config: fmt.Sprintf(testAccActionGatewayToolbeltConfig, name, `"exa_web_search", "exa_web_fetch"`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "tools.#", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "tools.*", "exa_web_fetch"),
					resource.TestCheckResourceAttr(resourceName, "version", "2"),
					resource.TestCheckResourceAttr(resourceName, "reference", name+"@2"),
					resource.TestCheckResourceAttr(resourceName, "tool_count", "2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// An import reports members pinned, as <tool_slug>@<version>.
				ImportStateVerifyIgnore: []string{"tools"},
			},
		},
	})
}

func TestAccDataSourceDigitalOceanActionGatewayToolbelt_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccActionGatewayToolbeltConfig, name, `"exa_web_search"`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckActionGatewayToolbeltDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + `
data "digitalocean_action_gateway_toolbelt" "foobar" {
  name = digitalocean_action_gateway_toolbelt.foobar.name
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.digitalocean_action_gateway_toolbelt.foobar", "reference", "digitalocean_action_gateway_toolbelt.foobar", "reference"),
					resource.TestCheckResourceAttr("data.digitalocean_action_gateway_toolbelt.foobar", "display_name", "Search Tools"),
					resource.TestCheckResourceAttr("data.digitalocean_action_gateway_toolbelt.foobar", "tools.#", "1"),
					resource.TestCheckTypeSetElemAttr("data.digitalocean_action_gateway_toolbelt.foobar", "tools.*", "exa_web_search@1"),
				),
			},
		},
	})
}

func testAccCheckActionGatewayToolbeltDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_action_gateway_toolbelt" {
			continue
		}
		_, resp, err := client.ActionGateway.Toolbelts.Get(context.Background(), rs.Primary.ID, nil)
		if err == nil {
			return fmt.Errorf("Action Gateway toolbelt %s still exists", rs.Primary.ID)
		}
		if resp == nil || resp.StatusCode != http.StatusNotFound {
			return err
		}
	}
	return nil
}

const testAccActionGatewayToolbeltConfig = `
resource "digitalocean_action_gateway_toolbelt" "foobar" {
  name         = "%s"
  display_name = "Search Tools"
  description  = "Tools for searching and fetching public web pages."
  tools        = [%s]
}
`
