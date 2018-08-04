package digitalocean

import (
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDigitalOceanFloatingIPAssignment(t *testing.T) {
	var floatingIP godo.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanFloatingIPAssignmentConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanFloatingIPExists("digitalocean_floating_ip_assignment.foobar", &floatingIP),
					resource.TestMatchResourceAttr(
						"digitalocean_floating_ip_assignment.foobar", "id", regexp.MustCompile("[0-9.]+")),
					resource.TestMatchResourceAttr(
						"digitalocean_floating_ip_assignment.foobar", "droplet_id", regexp.MustCompile("[0-9]+")),
				),
			},
			{
				Config: testAccCheckDigitalOceanFloatingIPAssignmentConfig_NoDropletId,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanFloatingIPExists("digitalocean_floating_ip_assignment.foobar", &floatingIP),
					resource.TestMatchResourceAttr(
						"digitalocean_floating_ip_assignment.foobar", "id", regexp.MustCompile("[0-9.]+")),
					resource.TestCheckResourceAttr(
						"digitalocean_floating_ip_assignment.foobar", "droplet_id", "0"),
				),
			},
		},
	})
}

var testAccCheckDigitalOceanFloatingIPAssignmentConfig = `
resource "digitalocean_floating_ip" "foobar" {
  region = "nyc3"
  ignore_droplet_id = true
}

resource "digitalocean_droplet" "foobar" {
  name               = "foobar"
  size               = "1gb"
  image              = "centos-7-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}

resource "digitalocean_floating_ip_assignment" "foobar" {
  ip_address = "${digitalocean_floating_ip.foobar.ip_address}"
  droplet_id = "${digitalocean_droplet.foobar.id}"
}
`

var testAccCheckDigitalOceanFloatingIPAssignmentConfig_NoDropletId = `
resource "digitalocean_floating_ip" "foobar" {
  region = "nyc3"
  ignore_droplet_id = true
}

resource "digitalocean_floating_ip_assignment" "foobar" {
  ip_address = "${digitalocean_floating_ip.foobar.ip_address}"
}
`
