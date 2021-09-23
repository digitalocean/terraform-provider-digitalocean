package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	slackChannels = `
slack {
	channel = "production-alerts"
	url		= "https://hooks.slack.com/services/T1234567/AAAAAAAA/ZZZZZZ"
}
	`

	multipleSlackChannel = `
slack {
	channel = "production-alerts"
	url		= "https://hooks.slack.com/services/T1234567/AAAAAAAA/ZZZZZZ"
}

slack {
	channel = "digitalocean-cloud-alerts"
	url		= "https://hooks.slack.com/services/T2345678/BBBBBBBB/XXXXXX"
}
`

	testAccAlertPolicy = `
resource "digitalocean_droplet" "web" {
	image  = "ubuntu-20-04-x64"
	name   = "web-1"
	region = "fra1"
	size   = "s-1vcpu-1gb"
  }

  resource "digitalocean_monitor_alert" "%s" {
	alerts  {
	  email 	= ["benny@digitalocean.com"]
      %s
	}
	window      = "%s"
	type        = "%s"
	compare     = "GreaterThan"
	value       = 95
	entities    = [digitalocean_droplet.web.id]
	description = "%s"
  }
`

	testAccAlertPolicySlackEmailAlerts = `
resource "digitalocean_droplet" "web" {
	image  = "ubuntu-20-04-x64"
	name   = "web-1"
	region = "fra1"
	size   = "s-1vcpu-1gb"
  }

  resource "digitalocean_monitor_alert" "%s" {
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

	testAccAlertPolicyWithTag = `
resource "digitalocean_tag" "test" {
    name = "%s"
}

resource "digitalocean_droplet" "web" {
	image  = "ubuntu-20-04-x64"
	name   = "web-1"
	region = "fra1"
	size   = "s-1vcpu-1gb"
	tags    = [digitalocean_tag.test.name]
  }

  resource "digitalocean_monitor_alert" "%s" {
	alerts  {
	  email 	= ["benny@digitalocean.com"]
	}
	window      = "%s"
	type        = "%s"
	compare     = "GreaterThan"
	value       = 95
	tags        = [digitalocean_tag.test.name]
	description = "%s"
}
`

	testAccAlertPolicyAddDroplet = `
resource "digitalocean_droplet" "web" {
	image  = "ubuntu-20-04-x64"
	name   = "web-1"
	region = "fra1"
	size   = "s-1vcpu-1gb"
  }

resource "digitalocean_droplet" "web2" {
	image  = "ubuntu-20-04-x64"
	name   = "web-2"
	region = "fra1"
	size   = "s-1vcpu-1gb"
}


resource "digitalocean_monitor_alert" "%s" {
	alerts  {
	  email 	= ["benny@digitalocean.com"]
      %s
	}
	window      = "%s"
	type        = "%s"
	compare     = "GreaterThan"
	value       = 95
	entities    = [digitalocean_droplet.web.id, digitalocean_droplet.web2.id]
	description = "%s"
  }
`
)

func TestAccDigitalOceanMonitorAlert(t *testing.T) {
	var randName = randomTestName()
	resourceName := fmt.Sprintf("digitalocean_monitor_alert.%s", randName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		ProviderFactories:         testAccProviderFactories,
		CheckDestroy:              testAccCheckDigitalOceanMonitorAlertDestroy,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccAlertPolicy, randName, "", "5m", "v1/insights/droplet/cpu", "Alert about CPU usage"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "v1/insights/droplet/cpu"),
					resource.TestCheckResourceAttr(resourceName, "compare", "GreaterThan"),
					resource.TestCheckResourceAttr(resourceName, "alerts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "alerts.0.email.0", "benny@digitalocean.com"),
				),
			},
		},
	})
}

func TestAccDigitalOceanMonitorAlertSlackEmailAlerts(t *testing.T) {
	var randName = randomTestName()
	resourceName := fmt.Sprintf("digitalocean_monitor_alert.%s", randName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		ProviderFactories:         testAccProviderFactories,
		CheckDestroy:              testAccCheckDigitalOceanMonitorAlertDestroy,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccAlertPolicySlackEmailAlerts, randName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "v1/insights/droplet/cpu"),
					resource.TestCheckResourceAttr(resourceName, "compare", "GreaterThan"),
					resource.TestCheckResourceAttr(resourceName, "alerts.#", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "alerts.0.email.0"),
					resource.TestCheckResourceAttr(resourceName, "alerts.0.slack.0.channel", "production-alerts"),
					resource.TestCheckResourceAttr(resourceName, "alerts.0.slack.0.url", "https://hooks.slack.com/services/T1234567/AAAAAAAA/ZZZZZZ"),
				),
			},
		},
	})
}

func TestAccDigitalOceanMonitorAlertUpdate(t *testing.T) {
	var randName = randomTestName()
	resourceName := fmt.Sprintf("digitalocean_monitor_alert.%s", randName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		ProviderFactories:         testAccProviderFactories,
		CheckDestroy:              testAccCheckDigitalOceanMonitorAlertDestroy,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccAlertPolicy, randName, "", "10m", "v1/insights/droplet/cpu", "Alert about CPU usage"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Alert about CPU usage"),
					resource.TestCheckResourceAttr(resourceName, "type", "v1/insights/droplet/cpu"),
					resource.TestCheckResourceAttr(resourceName, "compare", "GreaterThan"),
					resource.TestCheckResourceAttr(resourceName, "alerts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "entities.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "alerts.0.email.0", "benny@digitalocean.com"),
					resource.TestCheckResourceAttr(resourceName, "alerts.0.slack.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "window", "10m"),
				),
			},
			{
				Config: fmt.Sprintf(testAccAlertPolicyAddDroplet, randName, multipleSlackChannel, "5m", "v1/insights/droplet/memory_utilization_percent", "Alert about memory usage"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "description", "Alert about memory usage"),
					resource.TestCheckResourceAttr(resourceName, "type", "v1/insights/droplet/memory_utilization_percent"),
					resource.TestCheckResourceAttr(resourceName, "compare", "GreaterThan"),
					resource.TestCheckResourceAttr(resourceName, "alerts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "entities.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "alerts.0.email.0", "benny@digitalocean.com"),
					resource.TestCheckResourceAttr(resourceName, "alerts.0.slack.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "alerts.0.slack.0.channel", "production-alerts"),
					resource.TestCheckResourceAttr(resourceName, "alerts.0.slack.0.url", "https://hooks.slack.com/services/T1234567/AAAAAAAA/ZZZZZZ"),
					resource.TestCheckResourceAttr(resourceName, "alerts.0.slack.1.channel", "digitalocean-cloud-alerts"),
					resource.TestCheckResourceAttr(resourceName, "alerts.0.slack.1.url", "https://hooks.slack.com/services/T2345678/BBBBBBBB/XXXXXX"),
					resource.TestCheckResourceAttr(resourceName, "window", "5m"),
				),
			},
		},
	})
}

func TestAccDigitalOceanMonitorAlertWithTag(t *testing.T) {
	var (
		randName = randomTestName()
		tagName  = randomTestName()
	)
	resourceName := fmt.Sprintf("digitalocean_monitor_alert.%s", randName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		ProviderFactories:         testAccProviderFactories,
		CheckDestroy:              testAccCheckDigitalOceanMonitorAlertDestroy,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccAlertPolicyWithTag, tagName, randName, "5m", "v1/insights/droplet/cpu", "Alert about CPU usage"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "type", "v1/insights/droplet/cpu"),
					resource.TestCheckResourceAttr(resourceName, "compare", "GreaterThan"),
					resource.TestCheckResourceAttr(resourceName, "alerts.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "alerts.0.email.0", "benny@digitalocean.com"),
					resource.TestCheckResourceAttr(resourceName, "tags.0", tagName),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanMonitorAlertDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).godoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_monitor_alert" {
			continue
		}
		uuid := rs.Primary.Attributes["uuid"]

		// Try to find the monitor alert
		_, _, err := client.Monitoring.GetAlertPolicy(context.Background(), uuid)

		if err == nil {
			return fmt.Errorf("Monitor alert still exists")
		}
	}

	return nil
}
