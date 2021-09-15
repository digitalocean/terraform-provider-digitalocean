package digitalocean

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanMonitorAlert_importBasic(t *testing.T) {
	// copied from the database import test, but 3/3 does not fail ...

	randName := randomTestName()
	resourceName := fmt.Sprintf("digitalocean_monitor_alert.%s", randName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanMonitorAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccAlertPolicy, randName, "", "10m", "v1/insights/droplet/memory_utilization_percent", "Alert about memory usage"),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccMonitorAlertImportID(resourceName),
			},
			// Test importing non-existent resource provides expected error.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "this-monitor-alert-id-does-not-exist",
				ExpectError:       regexp.MustCompile(`some fancy error`),
			},
		},
	})
}

func testAccMonitorAlertImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("Not found: %s", n)
		}

		monitorAlertID := rs.Primary.Attributes["monitor_alert_id"]
		name := rs.Primary.Attributes["name"]

		return fmt.Sprintf("%s,%s", monitorAlertID, name), nil
	}
}
