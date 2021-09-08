package digitalocean

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanDatabaseMonitorAlert_importBasic(t *testing.T) {
	resourceName := "digitalocean_monitor_alert.cpu_alert"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseDBDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccAlertPolicy("10m")),
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
				ImportStateId:     fmt.Sprintf("%s,%s", "this-cluster-id-does-not-exist", monitor_alert_test_name),
				ExpectError:       regexp.MustCompile(`(Please verify the ID is correct|Cannot import non-existent remote object)`),
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
