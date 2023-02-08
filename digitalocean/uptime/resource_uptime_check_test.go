package uptime_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccCheckDigitalOceanUptimeCheckConfig_Basic = `
resource "digitalocean_uptime_check" "foobar" {
  name    = "%s"
  target  = "%s"
  regions = ["%s"]
}
`

func TestAccDigitalOceanUptimeCheck_Basic(t *testing.T) {
	checkName := acceptance.RandomTestName()
	checkTarget := "https://www.landingpage.com"
	checkRegions := "eu_west"
	checkCreateConfig := fmt.Sprintf(testAccCheckDigitalOceanUptimeCheckConfig_Basic, checkName, checkTarget, checkRegions)

	updatedCheckName := acceptance.RandomTestName()
	updatedCheckRegions := "us_east"
	checkUpdateConfig := fmt.Sprintf(testAccCheckDigitalOceanUptimeCheckConfig_Basic, updatedCheckName, checkTarget, updatedCheckRegions)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanUptimeCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: checkCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanUptimeCheckExists("digitalocean_uptime_check.foobar"),
					resource.TestCheckResourceAttr(
						"digitalocean_uptime_check.foobar", "name", checkName),
					resource.TestCheckResourceAttr(
						"digitalocean_uptime_check.foobar", "target", checkTarget),
					resource.TestCheckResourceAttr("digitalocean_uptime_check.foobar", "regions.#", "1"),
					resource.TestCheckTypeSetElemAttr("digitalocean_uptime_check.foobar", "regions.*", "eu_west"),
				),
			},
			{
				Config: checkUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanUptimeCheckExists("digitalocean_uptime_check.foobar"),
					resource.TestCheckResourceAttr(
						"digitalocean_uptime_check.foobar", "name", updatedCheckName),
					resource.TestCheckResourceAttr(
						"digitalocean_uptime_check.foobar", "target", checkTarget),
					resource.TestCheckResourceAttr("digitalocean_uptime_check.foobar", "regions.#", "1"),
					resource.TestCheckTypeSetElemAttr("digitalocean_uptime_check.foobar", "regions.*", "us_east"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanUptimeCheckDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {

		if rs.Type != "digitalocean_uptime_check" {
			continue
		}

		_, _, err := client.UptimeChecks.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Uptime Check resource still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanUptimeCheckExists(resource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		rs, ok := s.RootModule().Resources[resource]

		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID set for resource: %s", resource)
		}

		foundUptimeCheck, _, err := client.UptimeChecks.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundUptimeCheck.ID != rs.Primary.ID {
			return fmt.Errorf("Resource not found: %s : %s", resource, rs.Primary.ID)
		}

		return nil
	}
}
