package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccCheckDigitalOceanUptimeAlertConfig_basic = `
resource "digitalocean_uptime_check" "test" {
	name  = "terraform-test"
	target = "https://www.landingpage.com"
	regions = ["us_east", "eu_west"]
  }
  resource "digitalocean_uptime_alert" "foobar"  {
	check_id = "${digitalocean_uptime_check.test.id}"
	name = "%s"
  	type = "latency"
	threshold = %s
	comparison = "greater_than"
	notifications {
		email = [data.digitalocean_account.test.email]
	}
	period = "2m"
  }
  `

func TestAccDigitalOceanUptimeAlert_Basic(t *testing.T) {
	originalAlertName := randomTestName()
	originalThreshold := "300"
	updateThreshold := "250"

	checkCreateConfig := fmt.Sprintf(testAccCheckDigitalOceanUptimeAlertConfig_basic, originalAlertName, originalThreshold)
	checkUpdateConfig := fmt.Sprintf(testAccCheckDigitalOceanUptimeAlertConfig_basic, originalAlertName, updateThreshold)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanUptimeAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: checkCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanUptimeAlertExists("digitalocean_uptime_alert.foobar"),
					resource.TestCheckResourceAttr(
						"digitalocean_uptime_alert.foobar", "name", originalAlertName),
					resource.TestCheckResourceAttr(
						"digitalocean_uptime_alert.foobar", "threshold", originalThreshold),
					resource.TestCheckResourceAttr(
						"digitalocean_uptime_alert.foobar", "type", "latency"),
					resource.TestCheckResourceAttr(
						"digitalocean_uptime_alert.foobar", "comparison", "greater_than"),
					resource.TestCheckResourceAttr(
						"digitalocean_uptime_alert.foobar", "period", "2m"),
				),
			},
			{
				Config: checkUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanUptimeAlertExists("digitalocean_uptime_alert.foobar"),
					resource.TestCheckResourceAttr(
						"digitalocean_uptime_alert.foobar", "name", originalAlertName),
					resource.TestCheckResourceAttr(
						"digitalocean_uptime_alert.foobar", "threshold", updateThreshold),
					resource.TestCheckResourceAttr(
						"digitalocean_uptime_alert.foobar", "type", "latency"),
					resource.TestCheckResourceAttr(
						"digitalocean_uptime_alert.foobar", "comparison", "greater_than"),
					resource.TestCheckResourceAttr(
						"digitalocean_uptime_alert.foobar", "period", "2m"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanUptimeAlertDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).godoClient()

	var checkID string

	// get check ID
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_uptime_check" {
			continue
		}
		checkID = rs.Primary.ID
	}

	// check if alert exists, error if it does
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_uptime_alert" {
			continue
		}
		_, _, err := client.UptimeChecks.GetAlert(context.Background(), checkID, rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Uptime Alert resource still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanUptimeAlertExists(resource string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		rs, ok := s.RootModule().Resources[resource]

		if !ok {
			return fmt.Errorf("Not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID set for resource: %s", resource)
		}

		err, checkID := getCheckID("digitalocean_uptime_check.test", s)
		if err != nil {
			return fmt.Errorf("Error retrieve check ID for alert: %s", resource)
		}

		foundUptimeAlert, _, err := client.UptimeChecks.GetAlert(context.Background(), checkID, rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundUptimeAlert.ID != rs.Primary.ID {
			return fmt.Errorf("Resource not found: %s : %s", resource, rs.Primary.ID)
		}

		return nil
	}
}

func getCheckID(resource string, s *terraform.State) (error, string) {
	client := testAccProvider.Meta().(*CombinedConfig).godoClient()

	rs, ok := s.RootModule().Resources[resource]

	if !ok {
		return fmt.Errorf("Not found: %s", resource), ""
	}

	if rs.Primary.ID == "" {
		return fmt.Errorf("No ID set for resource: %s", resource), ""
	}

	foundUptimeCheck, _, err := client.UptimeChecks.Get(context.Background(), rs.Primary.ID)

	if err != nil {
		return err, ""
	}

	if foundUptimeCheck.ID != rs.Primary.ID {
		return fmt.Errorf("Resource not found: %s : %s", resource, rs.Primary.ID), ""
	}

	return nil, foundUptimeCheck.ID
}
