package functions_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// Functions themselves are deployed with doctl serverless rather than the
// public API, so trigger tests run against a namespace that already has a
// deployed function, e.g. the sample/hello function from
// `doctl serverless init --language js`.
func testAccFunctionsTriggerPreCheck(t *testing.T) (string, string) {
	acceptance.TestAccPreCheck(t)

	namespaceID := os.Getenv("DIGITALOCEAN_FUNCTIONS_NAMESPACE_ID")
	function := os.Getenv("DIGITALOCEAN_FUNCTIONS_FUNCTION")
	if namespaceID == "" || function == "" {
		t.Skip("DIGITALOCEAN_FUNCTIONS_NAMESPACE_ID and DIGITALOCEAN_FUNCTIONS_FUNCTION must be set to run Functions trigger acceptance tests")
	}
	return namespaceID, function
}

func TestAccDigitalOceanFunctionsTrigger_Basic(t *testing.T) {
	if os.Getenv(resource.EnvTfAcc) == "" {
		t.Skipf("Acceptance tests skipped unless env '%s' set", resource.EnvTfAcc)
	}
	namespaceID, function := testAccFunctionsTriggerPreCheck(t)
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanFunctionsTriggerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDigitalOceanFunctionsTriggerConfig(namespaceID, name, function, true, "*/30 * * * *", `{"name":"terraform"}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanFunctionsTriggerExists("digitalocean_functions_trigger.test"),
					resource.TestCheckResourceAttr("digitalocean_functions_trigger.test", "type", "SCHEDULED"),
					resource.TestCheckResourceAttr("digitalocean_functions_trigger.test", "enabled", "true"),
					resource.TestCheckResourceAttr("digitalocean_functions_trigger.test", "scheduled_details.0.cron", "*/30 * * * *"),
					resource.TestCheckResourceAttr("digitalocean_functions_trigger.test", "scheduled_details.0.body", `{"name":"terraform"}`),
					resource.TestCheckResourceAttrSet("digitalocean_functions_trigger.test", "next_run_at"),
					resource.TestCheckResourceAttrPair("data.digitalocean_functions_trigger.test", "function", "digitalocean_functions_trigger.test", "function"),
				),
			},
			{
				Config: testAccDigitalOceanFunctionsTriggerConfig(namespaceID, name, function, false, "0 6 * * *", `{"name":"updated"}`),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanFunctionsTriggerExists("digitalocean_functions_trigger.test"),
					resource.TestCheckResourceAttr("digitalocean_functions_trigger.test", "enabled", "false"),
					resource.TestCheckResourceAttr("digitalocean_functions_trigger.test", "scheduled_details.0.cron", "0 6 * * *"),
					resource.TestCheckResourceAttr("digitalocean_functions_trigger.test", "scheduled_details.0.body", `{"name":"updated"}`),
				),
			},
			{
				ResourceName:      "digitalocean_functions_trigger.test",
				ImportState:       true,
				ImportStateId:     fmt.Sprintf("%s,%s", namespaceID, name),
				ImportStateVerify: true,
				// Run times move on their own between steps.
				ImportStateVerifyIgnore: []string{"last_run_at", "next_run_at"},
			},
		},
	})
}

func testAccCheckDigitalOceanFunctionsTriggerDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_functions_trigger" {
			continue
		}

		_, _, err := client.Functions.GetTrigger(context.Background(), rs.Primary.Attributes["namespace_id"], rs.Primary.Attributes["name"])
		if err == nil {
			return fmt.Errorf("Functions trigger %s still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckDigitalOceanFunctionsTriggerExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		_, _, err := client.Functions.GetTrigger(context.Background(), rs.Primary.Attributes["namespace_id"], rs.Primary.Attributes["name"])
		return err
	}
}

func testAccDigitalOceanFunctionsTriggerConfig(namespaceID, name, function string, enabled bool, cron, body string) string {
	return fmt.Sprintf(`
resource "digitalocean_functions_trigger" "test" {
  namespace_id = "%s"
  name         = "%s"
  function     = "%s"
  enabled      = %t

  scheduled_details {
    cron = "%s"
    body = jsonencode(%s)
  }
}

data "digitalocean_functions_trigger" "test" {
  namespace_id = digitalocean_functions_trigger.test.namespace_id
  name         = digitalocean_functions_trigger.test.name
}
`, namespaceID, name, function, enabled, cron, body)
}
