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

// requireEndpoints reads one endpoint and skips the test when the endpoints
// route answers 404, which is how a team without the endpoint feature flag
// sees it. TF_ACC_REQUIRE_ENDPOINTS=1 turns that skip into a failure, so the
// Action Infra run of the suite exercises the resource and a flag regression
// is visible there.
func requireEndpoints(t *testing.T, name string) {
	t.Helper()

	if os.Getenv(resource.EnvTfAcc) == "" {
		t.Skipf("Acceptance tests skipped unless env '%s' set", resource.EnvTfAcc)
	}
	if os.Getenv("DIGITALOCEAN_TOKEN") == "" {
		t.Fatal("DIGITALOCEAN_TOKEN must be set for acceptance tests")
	}
	acceptance.TestAccPreCheck(t)
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	_, resp, err := client.ActionGateway.Endpoints.List(context.Background(), nil)
	if err == nil {
		return
	}
	if resp == nil || resp.StatusCode != http.StatusNotFound {
		t.Fatalf("error reading Action Gateway endpoints: %s", err)
	}
	if os.Getenv("TF_ACC_REQUIRE_ENDPOINTS") == "1" {
		t.Fatalf("TF_ACC_REQUIRE_ENDPOINTS=1 but the endpoints route answered 404 for %s: the Action Gateway endpoint feature flag is off for this team", name)
	}
	t.Skipf("skipping %s: the Action Gateway endpoints route answered 404, so the feature flag is off for this team", name)
}

func TestAccDigitalOceanActionGatewayEndpoint_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceName := "digitalocean_action_gateway_endpoint.foobar"

	requireEndpoints(t, resourceName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckActionGatewayEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccActionGatewayEndpointConfig, name, `"exa_web_search"`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "version.0.version_number", "1"),
					resource.TestCheckResourceAttr(resourceName, "version.0.tools.#", "1"),
					resource.TestCheckTypeSetElemAttr(resourceName, "version.0.tools.*", "exa_web_search"),
					resource.TestCheckResourceAttr(resourceName, "version.0.pinned_tools.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "version.0.policy.0.default_action", "ask"),
					resource.TestCheckResourceAttr(resourceName, "version.0.policy.0.rules.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "mcp_url"),
					resource.TestCheckResourceAttrSet(resourceName, "version_mcp_url"),
				),
			},
			{
				// A version change publishes a new head in place.
				Config: fmt.Sprintf(testAccActionGatewayEndpointConfig, name, `"exa_web_fetch"`),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "version.0.version_number", "2"),
					resource.TestCheckTypeSetElemAttr(resourceName, "version.0.tools.*", "exa_web_fetch"),
					testAccCheckActionGatewayEndpointKeptEndpoint(resourceName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// An import reports the head pinned, as <tool_slug>@<version>.
				ImportStateVerifyIgnore: []string{"version.0.tools", "version.0.pinned_tools", "insights"},
			},
		},
	})
}

// testAccCheckActionGatewayEndpointKeptEndpoint requires a version publish to
// keep the endpoint and leave the previous version live.
func testAccCheckActionGatewayEndpointKeptEndpoint(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("not found: %s", name)
		}
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		found, _, err := client.ActionGateway.Endpoints.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}
		if found.Endpoint == nil || found.Endpoint.HeadVersion == nil {
			return fmt.Errorf("endpoint %s has no head version", rs.Primary.ID)
		}
		if number := found.Endpoint.HeadVersion.VersionNumber; number < 2 {
			return fmt.Errorf("head version = %d, expected the publish to have advanced it", number)
		}

		versions, _, err := client.ActionGateway.Endpoints.ListVersions(context.Background(), rs.Primary.ID, nil)
		if err != nil {
			return err
		}
		if len(versions.Versions) < 2 {
			return fmt.Errorf("live versions = %d, expected the previous version to stay live", len(versions.Versions))
		}
		return nil
	}
}

func TestAccDataSourceDigitalOceanActionGatewayEndpoint_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccActionGatewayEndpointConfig, name, `"exa_web_search"`)
	dataName := "data.digitalocean_action_gateway_endpoint.foobar"

	requireEndpoints(t, dataName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckActionGatewayEndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig + `
data "digitalocean_action_gateway_endpoint" "foobar" {
  name = digitalocean_action_gateway_endpoint.foobar.name
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(dataName, "id", "digitalocean_action_gateway_endpoint.foobar", "id"),
					resource.TestCheckResourceAttrPair(dataName, "mcp_url", "digitalocean_action_gateway_endpoint.foobar", "mcp_url"),
					resource.TestCheckResourceAttrPair(dataName, "version.0.version_number", "digitalocean_action_gateway_endpoint.foobar", "version.0.version_number"),
					resource.TestCheckResourceAttrPair(dataName, "version.0.pinned_tools.#", "digitalocean_action_gateway_endpoint.foobar", "version.0.pinned_tools.#"),
				),
			},
		},
	})
}

func testAccCheckActionGatewayEndpointDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_action_gateway_endpoint" {
			continue
		}
		_, resp, err := client.ActionGateway.Endpoints.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Action Gateway endpoint %s still exists", rs.Primary.ID)
		}
		if resp == nil || resp.StatusCode != http.StatusNotFound {
			return err
		}
	}
	return nil
}

// testAccActionGatewayEndpointConfig takes the endpoint name and its tool
// selection, so a step can change the version block without renaming the
// endpoint.
const testAccActionGatewayEndpointConfig = `
resource "digitalocean_action_gateway_endpoint" "foobar" {
  name = "%s"

  insights {
    logs   = true
    traces = false
  }

  version {
    tools = [%s]

    policy {
      default_action = "ask"

      rules {
        tool   = %[2]s
        action = "allow"
      }
    }
  }
}
`
