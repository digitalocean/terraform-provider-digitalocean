package reservedip_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanFloatingIP_Region(t *testing.T) {
	var floatingIP godo.FloatingIP

	expectedURNRegEx, _ := regexp.Compile(`do:floatingip:(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanFloatingIPConfig_region,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanFloatingIPExists("digitalocean_floating_ip.foobar", &floatingIP),
					resource.TestCheckResourceAttr(
						"digitalocean_floating_ip.foobar", "region", "nyc3"),
					resource.TestMatchResourceAttr("digitalocean_floating_ip.foobar", "urn", expectedURNRegEx),
				),
			},
		},
	})
}

func TestAccDigitalOceanFloatingIP_Droplet(t *testing.T) {
	var floatingIP godo.FloatingIP
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanFloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanFloatingIPConfig_droplet(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanFloatingIPExists("digitalocean_floating_ip.foobar", &floatingIP),
					resource.TestCheckResourceAttr(
						"digitalocean_floating_ip.foobar", "region", "nyc3"),
				),
			},
			{
				Config: testAccCheckDigitalOceanFloatingIPConfig_Reassign(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanFloatingIPExists("digitalocean_floating_ip.foobar", &floatingIP),
					resource.TestCheckResourceAttr(
						"digitalocean_floating_ip.foobar", "region", "nyc3"),
				),
			},
			{
				Config: testAccCheckDigitalOceanFloatingIPConfig_Unassign(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanFloatingIPExists("digitalocean_floating_ip.foobar", &floatingIP),
					resource.TestCheckResourceAttr(
						"digitalocean_floating_ip.foobar", "region", "nyc3"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanFloatingIPDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_floating_ip" {
			continue
		}

		// Try to find the key
		_, _, err := client.FloatingIPs.Get(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Floating IP still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanFloatingIPExists(n string, floatingIP *godo.FloatingIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		// Try to find the FloatingIP
		foundFloatingIP, _, err := client.FloatingIPs.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundFloatingIP.IP != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*floatingIP = *foundFloatingIP

		return nil
	}
}

var testAccCheckDigitalOceanFloatingIPConfig_region = `
resource "digitalocean_floating_ip" "foobar" {
  region = "nyc3"
}`

func testAccCheckDigitalOceanFloatingIPConfig_droplet(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name               = "%s"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}

resource "digitalocean_floating_ip" "foobar" {
  droplet_id = digitalocean_droplet.foobar.id
  region     = digitalocean_droplet.foobar.region
}`, name)
}

func testAccCheckDigitalOceanFloatingIPConfig_Reassign(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "baz" {
  name               = "%s"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}

resource "digitalocean_floating_ip" "foobar" {
  droplet_id = digitalocean_droplet.baz.id
  region     = digitalocean_droplet.baz.region
}`, name)
}

func testAccCheckDigitalOceanFloatingIPConfig_Unassign(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "baz" {
  name               = "%s"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}

resource "digitalocean_floating_ip" "foobar" {
  region = "nyc3"
}`, name)
}
