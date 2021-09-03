package digitalocean

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func testAccAlertPolicy(window string) string {
	return fmt.Sprint(`
resource "digitalocean_droplet" "web" {
	image  = "ubuntu-20-04-x64"
	name   = "web-1"
	region = "nyc2"
	size   = "s-1vcpu-1gb"
  }
  
  resource "digitalocean_monitoring" "cpu_alert" {
	alerts      = {
	  email = ["benny@digitalocean.com"]
	  slack = {
		"channel"   = "Production Alerts",
		"url"       = "https://hooks.slack.com/services/T1234567/AAAAAAAA/ZZZZZZ"
	}
	window      = %s
	type        = "v1/insights/droplet/cpu"
	compare     = "GreaterThan"
	value       = 95
	enabled     = true
	entities    = [digitalocean_droplet.web.id]
	description = "Alert about CPU usage"
  }
`, window)
}

func TestAccDigitalOceanMonitorAlert(t *testing.T) {
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                  func() { testAccPreCheck(t) },
		ProviderFactories:         testAccProviderFactories,
		PreventPostDestroyRefresh: true,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDigitalOceanSpacesObjectsConfigResources(rInt), // NOTE: contains no data source
				// Does not need Check
			},
			{
				Config: testAccDataSourceDigitalOceanSpacesObjectsConfigMaxKeys(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSpacesObjectsDataSourceExists("data.digitalocean_spaces_bucket_objects.yesh"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "keys.#", "2"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "keys.0", "arch/courthouse_towers/landscape"),
					resource.TestCheckResourceAttr("data.digitalocean_spaces_bucket_objects.yesh", "keys.1", "arch/navajo/north_window"),
				),
			},
		},
	})
}
