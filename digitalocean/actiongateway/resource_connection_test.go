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

func TestAccDigitalOceanActionGatewayConnection_OAuth(t *testing.T) {
	userID := acceptance.RandomTestName()
	resourceName := "digitalocean_action_gateway_connection.foobar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckActionGatewayConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccActionGatewayConnectionConfig, userID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "tool_provider", "github"),
					resource.TestCheckResourceAttr(resourceName, "user_id", userID),
					resource.TestCheckResourceAttr(resourceName, "status", "pending"),
					resource.TestCheckResourceAttr(resourceName, "credential_kind", "digitalocean_oauth"),
					resource.TestCheckResourceAttrSet(resourceName, "tool_provider_display_name"),
					resource.TestCheckResourceAttrSet(resourceName, "owning_user_id"),
					resource.TestCheckResourceAttr(resourceName, "authorization.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "authorization.0.connect_url"),
					resource.TestCheckResourceAttrSet(resourceName, "authorization.0.expires_at"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"authorization"},
			},
		},
	})
}

func TestAccDataSourceDigitalOceanActionGatewayConnection_ByUserID(t *testing.T) {
	userID := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccActionGatewayConnectionConfig, userID)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckActionGatewayConnectionDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + `
data "digitalocean_action_gateway_connection" "by_user" {
  tool_provider = digitalocean_action_gateway_connection.foobar.tool_provider
  user_id       = digitalocean_action_gateway_connection.foobar.user_id
}

data "digitalocean_action_gateway_connection" "by_id" {
  id = digitalocean_action_gateway_connection.foobar.id
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.digitalocean_action_gateway_connection.by_user", "id", "digitalocean_action_gateway_connection.foobar", "id"),
					resource.TestCheckResourceAttrPair("data.digitalocean_action_gateway_connection.by_user", "status", "digitalocean_action_gateway_connection.foobar", "status"),
					resource.TestCheckResourceAttrPair("data.digitalocean_action_gateway_connection.by_id", "user_id", "digitalocean_action_gateway_connection.foobar", "user_id"),
					resource.TestCheckResourceAttrPair("data.digitalocean_action_gateway_connection.by_id", "credential_kind", "digitalocean_action_gateway_connection.foobar", "credential_kind"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanActionGatewayProviders_Filter(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
data "digitalocean_action_gateway_providers" "github" {
  filter {
    key    = "name"
    values = ["github"]
  }
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_action_gateway_providers.github", "providers.#", "1"),
					resource.TestCheckResourceAttr("data.digitalocean_action_gateway_providers.github", "providers.0.name", "github"),
					resource.TestCheckResourceAttrSet("data.digitalocean_action_gateway_providers.github", "providers.0.display_name"),
				),
			},
		},
	})
}

func testAccCheckActionGatewayConnectionDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_action_gateway_connection" {
			continue
		}
		found, resp, err := client.ActionGateway.Connections.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				continue
			}
			return err
		}
		if found.Connection != nil && found.Connection.Status != "revoked" {
			return fmt.Errorf("Action Gateway connection %s still exists (status %s)", rs.Primary.ID, found.Connection.Status)
		}
	}
	return nil
}

const testAccActionGatewayConnectionConfig = `
resource "digitalocean_action_gateway_connection" "foobar" {
  tool_provider = "github"
  user_id       = "%s"
}
`
