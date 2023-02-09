package monitoring_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanMonitorAlert_importBasic(t *testing.T) {
	randName := acceptance.RandomTestName()
	resourceName := fmt.Sprintf("digitalocean_monitor_alert.%s", randName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanMonitorAlertDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccAlertPolicy, randName, randName, "", "10m", "v1/insights/droplet/memory_utilization_percent", "Alert about memory usage"),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Test importing non-existent resource provides expected error.
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: false,
				ImportStateId:     "this-monitor-alert-id-does-not-exist",
				ExpectError:       regexp.MustCompile(`Please verify the ID is correct|Cannot import non-existent remote object`),
			},
		},
	})
}
