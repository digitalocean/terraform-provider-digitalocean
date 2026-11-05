package actiongateway_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanActionGatewayMCPServer_Basic(t *testing.T) {
	endpoint := os.Getenv("DO_ACTION_GATEWAY_MCP_ENDPOINT")
	if endpoint == "" {
		t.Skip("DO_ACTION_GATEWAY_MCP_ENDPOINT must be set to a public HTTPS MCP endpoint that needs no authentication")
	}
	serverRef := acceptance.RandomTestName()
	resourceName := "digitalocean_action_gateway_mcp_server.foobar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckActionGatewayMCPServerDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccActionGatewayMCPServerConfig, serverRef, endpoint, "Acceptance test server"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "server_ref", serverRef),
					resource.TestCheckResourceAttr(resourceName, "endpoint", endpoint),
					resource.TestCheckResourceAttr(resourceName, "transport", "streamable_http"),
					resource.TestCheckResourceAttr(resourceName, "credential_ref_source", "none"),
					resource.TestCheckResourceAttr(resourceName, "description", "Acceptance test server"),
					resource.TestCheckResourceAttrSet(resourceName, "last_synced_at"),
					resource.TestCheckResourceAttr(resourceName, "sync_error", ""),
				),
			},
			{
				Config: fmt.Sprintf(testAccActionGatewayMCPServerConfig, serverRef, endpoint, "Updated description"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Updated description"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckActionGatewayMCPServerDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_action_gateway_mcp_server" {
			continue
		}
		_, resp, err := client.ActionGateway.MCPServers.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Action Gateway MCP server %s still exists", rs.Primary.ID)
		}
		if resp == nil || resp.StatusCode != http.StatusNotFound {
			return err
		}
	}
	return nil
}

const testAccActionGatewayMCPServerConfig = `
resource "digitalocean_action_gateway_mcp_server" "foobar" {
  server_ref  = "%s"
  endpoint    = "%s"
  description = "%s"
}
`
