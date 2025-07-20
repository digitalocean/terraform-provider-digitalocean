package genai_test

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

func TestAccDigitalOceanAgentRoute_Create(t *testing.T) {
	var agentRoute godo.Agent

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanAgentRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanAgentRouteConfig_create(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAgentRouteExists("digitalocean_genai_agent_route.test", &agentRoute),
					resource.TestCheckResourceAttr("digitalocean_genai_agent_route.test", "parent_agent_uuid", "12345678-1234-1234-1234-123456789012"),
					resource.TestCheckResourceAttr("digitalocean_genai_agent_route.test", "child_agent_uuid", "12345678-1234-1234-1234-123456789012"),
					resource.TestCheckResourceAttr("digitalocean_genai_agent_route.test", "route_name", "weather_route"),
					resource.TestCheckResourceAttr("digitalocean_genai_agent_route.test", "if_case", "use this to get weather information"),
				),
			},
		},
	})
}

func TestAccDigitalOceanAgentRoute_Update(t *testing.T) {
	var agentRoute godo.Agent

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanAgentRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanAgentRouteConfig_create(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAgentRouteExists("digitalocean_genai_agent_route.test", &agentRoute),
					resource.TestCheckResourceAttr("digitalocean_genai_agent_route.test", "parent_agent_uuid", "12345678-1234-1234-1234-123456789012"),
					resource.TestCheckResourceAttr("digitalocean_genai_agent_route.test", "child_agent_uuid", "12345678-1234-1234-1234-123456789012"),
				),
			},
			{
				Config: testAccCheckDigitalOceanAgentRouteConfig_update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAgentRouteExists("digitalocean_genai_agent_route.test", &agentRoute),
					resource.TestCheckResourceAttr("digitalocean_genai_agent_route.test", "parent_agent_uuid", "12345678-1234-1234-1234-123456789012"),
					resource.TestCheckResourceAttr("digitalocean_genai_agent_route.test", "child_agent_uuid", "12345678-1234-1234-1234-123456789012"),
					resource.TestCheckResourceAttr("digitalocean_genai_agent_route.test", "route_name", "updated_weather_route"),
					resource.TestCheckResourceAttr("digitalocean_genai_agent_route.test", "if_case", "updated: use this to get weather information"),
				),
			},
		},
	})
}

func TestAccDigitalOceanAgentRoute_RequiredFieldsOnly(t *testing.T) {
	var agentRoute godo.Agent

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanAgentRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanAgentRouteConfig_requiredOnly(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAgentRouteExists("digitalocean_genai_agent_route.test", &agentRoute),
					resource.TestCheckResourceAttr("digitalocean_genai_agent_route.test", "parent_agent_uuid", "12345678-1234-1234-1234-123456789012"),
					resource.TestCheckResourceAttr("digitalocean_genai_agent_route.test", "child_agent_uuid", "12345678-1234-1234-1234-123456789012"),
					resource.TestCheckResourceAttr("digitalocean_genai_agent_route.test", "route_name", ""),
					resource.TestCheckResourceAttr("digitalocean_genai_agent_route.test", "if_case", ""),
				),
			},
		},
	})
}

func TestAccDigitalOceanAgentRoute_Delete(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanAgentRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanAgentRouteConfig_create(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAgentRouteExists("digitalocean_genai_agent_route.test", nil),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanAgentRouteConfig_create() string {
	return `
resource "digitalocean_genai_agent_route" "test" {
  parent_agent_uuid = "12345678-1234-1234-1234-123456789012"
  child_agent_uuid  = "12345678-1234-1234-1234-123456789012"
  route_name        = "weather_route"
  if_case          = "use this to get weather information"
}
`
}

func testAccCheckDigitalOceanAgentRouteConfig_update() string {
	return `
resource "digitalocean_genai_agent_route" "test" {
  parent_agent_uuid = "12345678-1234-1234-1234-123456789012"
  child_agent_uuid  = "12345678-1234-1234-1234-123456789012"
  route_name        = "updated_weather_route"
  if_case          = "updated: use this to get weather information"
}
`
}

func testAccCheckDigitalOceanAgentRouteConfig_requiredOnly() string {
	return `
resource "digitalocean_genai_agent_route" "test" {
  parent_agent_uuid = "12345678-1234-1234-1234-123456789012"
  child_agent_uuid  = "12345678-1234-1234-1234-123456789012"
}
`
}

func testAccCheckDigitalOceanAgentRouteExists(n string, agentRoute *godo.Agent) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("resource %s not found in state", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource ID is not set")
		}

		if rs.Primary.Attributes["parent_agent_uuid"] == "" {
			return fmt.Errorf("parent_agent_uuid is required but not set")
		}

		if rs.Primary.Attributes["child_agent_uuid"] == "" {
			return fmt.Errorf("child_agent_uuid is required but not set")
		}

		return nil
	}
}

func testAccCheckDigitalOceanAgentRouteDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_genai_agent_route" {
			continue
		}

		routeUUID := rs.Primary.ID
		if routeUUID == "" {
			continue
		}

		_, _, err := client.GenAI.DeleteAgentRoute(ctx, routeUUID, rs.Primary.Attributes["parent_agent_uuid"])
		if err != nil {
			if !isNotFoundError(err) {
				return fmt.Errorf("unexpected error when checking if agent route %s was destroyed: %v", routeUUID, err)
			}
		}
	}
	return nil
}

func isNotFoundError(err error) bool {
	return err != nil && (err.Error() == "404 Not Found" ||
		err.Error() == "not found" ||
		fmt.Sprintf("%v", err) == "404")
}
