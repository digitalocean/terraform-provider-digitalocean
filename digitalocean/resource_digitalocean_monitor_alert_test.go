package digitalocean

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// update and delete tests missing

const (
	monitor_alert_test_name = "cpu_alert"
)

func testAccAlertPolicy(window string) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "web" {
	image  = "ubuntu-20-04-x64"
	name   = "web-1"
	region = "fra1"
	size   = "s-1vcpu-1gb"
  }
  
  resource "digitalocean_monitor_alert" "%s" {
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
`, monitor_alert_test_name, window)
}

func testAccAlertPolicyEmailAlerts() string {
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

func testAccAlertPolicySlackEmailAlerts() string {
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
	  slack {
		channel = "production-alerts"
		url		= "https://hooks.slack.com/services/T1234567/AAAAAAAA/ZZZZZZ"
	  }
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
		CheckDestroy:              testAccCheckDigitalOceanDropletDestroy,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertPolicy("5m"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_spaces_monitor_alert.cpu_alert", "type", "v1/insights/droplet/cpu"),
					resource.TestCheckResourceAttr("digitalocean_spaces_monitor_alert.cpu_alert", "compare", "GreaterThan"),
					resource.TestCheckResourceAttr("digitalocean_spaces_monitor_alert.cpu_alert", "alerts.email.0", "benny@digitalocean.com"),
					// resource.TestCheckResourceAttr("digitalocean_spaces_monitor_alert.cpu_alert", "alerts.slack.0.channel", "Production Alerts"),
					// resource.TestCheckResourceAttr("digitalocean_spaces_monitor_alert.cpu_alert", "alerts.slack.0.url", "https://hooks.slack.com/services/T1234567/AAAAAAAA/ZZZZZZ"),
				),
			},
		},
	})
}

func TestAccDigitalOceanMonitorAlertEmailAlerts(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		ProviderFactories:         testAccProviderFactories,
		CheckDestroy:              testAccCheckDigitalOceanDropletDestroy,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertPolicyEmailAlerts(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_spaces_monitor_alert.cpu_alert", "type", "v1/insights/droplet/cpu"),
					resource.TestCheckResourceAttr("digitalocean_spaces_monitor_alert.cpu_alert", "compare", "GreaterThan"),
					resource.TestCheckNoResourceAttr("digitalocean_spaces_monitor_alert.cpu_alert", "alerts"),
					resource.TestCheckNoResourceAttr("digitalocean_spaces_monitor_alert.alerts.email.0", "benny@digitalocean.com"),
				),
			},
		},
	})
}

func TestAccDigitalOceanMonitorAlertSlackEmailAlerts(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		ProviderFactories:         testAccProviderFactories,
		CheckDestroy:              testAccCheckDigitalOceanDropletDestroy,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertPolicySlackEmailAlerts(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_monitor_alert.cpu_alert", "type", "v1/insights/droplet/cpu"),
					resource.TestCheckResourceAttr("digitalocean_monitor_alert.cpu_alert", "compare", "GreaterThan"),
					resource.TestCheckResourceAttr("digitalocean_monitor_alert.cpu_alert", "alerts.#", "1"),
					resource.TestCheckResourceAttrSet("digitalocean_monitor_alert.cpu_alert", "alerts.0.email.0"),
					resource.TestCheckResourceAttr("digitalocean_monitor_alert.cpu_alert", "alerts.0.slack.0.channel", "production-alerts"),
					resource.TestCheckResourceAttr("digitalocean_monitor_alert.cpu_alert", "alerts.0.slack.0.url", "https://hooks.slack.com/services/T1234567/AAAAAAAA/ZZZZZZ"),
				),
			},
		},
	})
}

// change the type, and see that it'll change
func TestAccDigitalOceanMonitorAlertUpdate(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		ProviderFactories:         testAccProviderFactories,
		CheckDestroy:              testAccCheckDigitalOceanDropletDestroy,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccAlertPolicySlackEmailAlerts(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_spaces_monitor_alert.cpu_alert", "type", "v1/insights/droplet/cpu"),
					resource.TestCheckResourceAttr("digitalocean_spaces_monitor_alert.cpu_alert", "compare", "GreaterThan"),
					resource.TestCheckNoResourceAttr("digitalocean_spaces_monitor_alert.cpu_alert", "alerts"),
					resource.TestCheckNoResourceAttr("digitalocean_spaces_monitor_alert.alerts.email.0", "benny@digitalocean.com"),
					resource.TestCheckNoResourceAttr("digitalocean_spaces_monitor_alert.alerts.slack.0.channel", "production-alerts"),
					resource.TestCheckNoResourceAttr("digitalocean_spaces_monitor_alert.alerts.slack.0.url", "https://hooks.slack.com/services/T1234567/AAAAAAAA/ZZZZZZ"),
				),
			},
			{
				Config: testAccAlertPolicySlackEmailAlerts(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_monitor_alert.cpu_alert", "type", "v1/insights/droplet/cpu"),
					resource.TestCheckResourceAttr("digitalocean_monitor_alert.cpu_alert", "compare", "GreaterThan"),
					resource.TestCheckResourceAttr("digitalocean_monitor_alert.cpu_alert", "alerts.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_monitor_alert.cpu_alert", "alerts.0.email.0", "benny@digitalocean.com"),
					resource.TestCheckResourceAttr("digitalocean_monitor_alert.cpu_alert", "alerts.0.slack.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_monitor_alert.cpu_alert", "alerts.0.slack.0.channel", "production-alerts"),
					resource.TestCheckResourceAttr("digitalocean_monitor_alert.cpu_alert", "alerts.0.slack.0.url", "https://hooks.slack.com/services/T1234567/AAAAAAAA/ZZZZZZ"),
				),
			},
		},
	})
}
