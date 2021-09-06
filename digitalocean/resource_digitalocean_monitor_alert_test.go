package digitalocean

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// update and delete tests missing

func testAccAlertPolicy(window string) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "web" {
	image  = "ubuntu-20-04-x64"
	name   = "web-1"
	region = "fra1"
	size   = "s-1vcpu-1gb"
  }
  
  resource "digitalocean_monitor_alert" "cpu_alert" {
	alerts  {
	  email 	= ["benny@digitalocean.com"]
	}
	window      = "%s"
	type        = "v1/insights/droplet/cpu"
	compare     = "GreaterThan"
	value       = 95
	entities    = [digitalocean_droplet.web.id]
	description = "Alert about CPU usage"
  }
`, window)
}

func testAccAlertPolicyNoAlerts() string {
	return `
resource "digitalocean_droplet" "web" {
	image  = "ubuntu-20-04-x64"
	name   = "web-1"
	region = "fra1"
	size   = "s-1vcpu-1gb"
  }
  
  resource "digitalocean_monitor_alert" "cpu_alert" {
	alerts {
	  email 	= ["benny@digitalocean.com"]
	}
	window      = "5m"
	type        = "v1/insights/droplet/cpu"
	compare     = "GreaterThan"
	value       = 95
	entities    = [digitalocean_droplet.web.id]
	description = "Alert about CPU usage"
  }
`
}

func TestAccDigitalOceanMonitorAlert(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		ProviderFactories:         testAccProviderFactories,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertPolicy("5m"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_spaces_monitor_alert.cpu_alert", "type", "v1/insights/droplet/cpu"),
					resource.TestCheckResourceAttr("digitalocean_spaces_monitor_alert.cpu_alert", "compare", "GreaterThan"),
					resource.TestCheckResourceAttr("digitalocean_spaces_monitor_alert.cpu_alert", "alerts.0.email", "benny@digitalocean.com"),
					// resource.TestCheckResourceAttr("digitalocean_spaces_monitor_alert.cpu_alert", "alerts.slack.0.channel", "Production Alerts"),
					// resource.TestCheckResourceAttr("digitalocean_spaces_monitor_alert.cpu_alert", "alerts.slack.0.url", "https://hooks.slack.com/services/T1234567/AAAAAAAA/ZZZZZZ"),
				),
			},
		},
	})
}

func TestAccDigitalOceanMonitorAlertNoAlerts(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		ProviderFactories:         testAccProviderFactories,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertPolicyNoAlerts(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_spaces_monitor_alert.cpu_alert", "type", "v1/insights/droplet/cpu"),
					resource.TestCheckResourceAttr("digitalocean_spaces_monitor_alert.cpu_alert", "compare", "GreaterThan"),
					resource.TestCheckNoResourceAttr("digitalocean_spaces_monitor_alert.cpu_alert", "alerts"),
					resource.TestCheckNoResourceAttr("digitalocean_spaces_monitor_alert.alerts.0.email", "benny@digitalocean.com"),
				),
				// ExpectError: "",
			},
		},
	})
}

// ideas for tests:
// email/slack required
//
