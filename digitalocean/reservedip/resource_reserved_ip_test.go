package reservedip_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanReservedIP_Region(t *testing.T) {
	var reservedIP godo.ReservedIP

	expectedURNRegEx, _ := regexp.Compile(`do:reservedip:(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanReservedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanReservedIPConfig_region,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanReservedIPExists("digitalocean_reserved_ip.foobar", &reservedIP),
					resource.TestCheckResourceAttr(
						"digitalocean_reserved_ip.foobar", "region", "nyc3"),
					resource.TestMatchResourceAttr("digitalocean_reserved_ip.foobar", "urn", expectedURNRegEx),
				),
			},
		},
	})
}

func TestAccDigitalOceanReservedIP_Droplet(t *testing.T) {
	var reservedIP godo.ReservedIP
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanReservedIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanReservedIPConfig_droplet(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanReservedIPExists("digitalocean_reserved_ip.foobar", &reservedIP),
					resource.TestCheckResourceAttr(
						"digitalocean_reserved_ip.foobar", "region", "nyc3"),
				),
			},
			{
				Config: testAccCheckDigitalOceanReservedIPConfig_Reassign(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanReservedIPExists("digitalocean_reserved_ip.foobar", &reservedIP),
					resource.TestCheckResourceAttr(
						"digitalocean_reserved_ip.foobar", "region", "nyc3"),
				),
			},
			{
				Config: testAccCheckDigitalOceanReservedIPConfig_Unassign(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanReservedIPExists("digitalocean_reserved_ip.foobar", &reservedIP),
					resource.TestCheckResourceAttr(
						"digitalocean_reserved_ip.foobar", "region", "nyc3"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanReservedIPDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_reserved_ip" {
			continue
		}

		// Try to find the key
		_, _, err := client.ReservedIPs.Get(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Reserved IP still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanReservedIPExists(n string, reservedIP *godo.ReservedIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		// Try to find the ReservedIP
		foundReservedIP, _, err := client.ReservedIPs.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundReservedIP.IP != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*reservedIP = *foundReservedIP

		return nil
	}
}

var testAccCheckDigitalOceanReservedIPConfig_region = `
resource "digitalocean_reserved_ip" "foobar" {
  region = "nyc3"
}`

func testAccCheckDigitalOceanReservedIPConfig_droplet(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name               = "tf-acc-test-%d"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}

resource "digitalocean_reserved_ip" "foobar" {
  droplet_id = digitalocean_droplet.foobar.id
  region     = digitalocean_droplet.foobar.region
}`, rInt)
}

func testAccCheckDigitalOceanReservedIPConfig_Reassign(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "baz" {
  name               = "tf-acc-test-%d"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}

resource "digitalocean_reserved_ip" "foobar" {
  droplet_id = digitalocean_droplet.baz.id
  region     = digitalocean_droplet.baz.region
}`, rInt)
}

func testAccCheckDigitalOceanReservedIPConfig_Unassign(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "baz" {
  name               = "tf-acc-test-%d"
  size               = "s-1vcpu-1gb"
  image              = "ubuntu-22-04-x64"
  region             = "nyc3"
  ipv6               = true
  private_networking = true
}

resource "digitalocean_reserved_ip" "foobar" {
  region     = "nyc3"
}`, rInt)
}
