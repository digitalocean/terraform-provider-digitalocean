package reservedipv6_test

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

func TestAccDigitalOceanReservedIPV6_RegionSlug(t *testing.T) {
	var reservedIPv6 godo.ReservedIPV6

	expectedURNRegex, _ := regexp.Compile(`do:reservedipv6:/^[0-9A-F]{8}-[0-9A-F]{4}-[4][0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$/i`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanReservedIPV6Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanReservedIPV6Config_regionSlug,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanReservedIPV6Exists("digitalocean_reserved_ipv6.foobar", &reservedIPv6),
					resource.TestCheckResourceAttr(
						"digitalocean_reserved_ipv6.foobar", "region_slug", "nyc3"),
					resource.TestMatchResourceAttr("digitalocean_reserved_ipv6.foobar", "urn", expectedURNRegex),
				),
			},
		},
	})
}

func TestAccDigitalOceanReservedIPV6_Droplet(t *testing.T) {
	var reservedIPv6 godo.ReservedIPV6
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanReservedIPV6Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanReservedIPV6Config_droplet(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanReservedIPV6Exists("digitalocean_reserved_ipv6.foobar", &reservedIPv6),
					resource.TestCheckResourceAttr(
						"digitalocean_reserved_ipv6.foobar", "region_slug", "nyc3"),
				),
			},
			{
				Config: testAccCheckDigitalOceanReservedIPv6Config_Reassign(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanReservedIPV6Exists("digitalocean_reserved_ipv6.foobar", &reservedIPv6),
					resource.TestCheckResourceAttr(
						"digitalocean_reserved_ipv6.foobar", "region_slug", "nyc3"),
				),
			},
			{
				Config: testAccCheckDigitalOceanReservedIPV6Config_Unassign(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanReservedIPV6Exists("digitalocean_reserved_ipv6.foobar", &reservedIPv6),
					resource.TestCheckResourceAttr(
						"digitalocean_reserved_ipv6.foobar", "region_slug", "nyc3"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanReservedIPV6Destroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_reserved_ipv6" {
			continue
		}

		// Try to find the key
		_, _, err := client.ReservedIPV6s.Get(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Reserved IPv6 still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanReservedIPV6Exists(n string, reservedIPv6 *godo.ReservedIPV6) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		// Try to find the ReservedIPv6
		foundReservedIP, _, err := client.ReservedIPV6s.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundReservedIP.IP != rs.Primary.ID {
			return fmt.Errorf("Record not found")
		}

		*reservedIPv6 = *foundReservedIP

		return nil
	}
}

var testAccCheckDigitalOceanReservedIPV6Config_regionSlug = `
resource "digitalocean_reserved_ipv6" "foobar" {
  region_slug = "nyc3"
}`

func testAccCheckDigitalOceanReservedIPV6Config_droplet(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foobar" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
  ipv6   = true
}

resource "digitalocean_reserved_ipv6" "foobar" {
  droplet_id  = digitalocean_droplet.foobar.id
  region_slug = digitalocean_droplet.foobar.region
}`, name)
}

func testAccCheckDigitalOceanReservedIPv6Config_Reassign(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "baz" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
  ipv6   = true
}

resource "digitalocean_reserved_ipv6" "foobar" {
  droplet_id  = digitalocean_droplet.baz.id
  region_slug = digitalocean_droplet.baz.region
}`, name)
}

func testAccCheckDigitalOceanReservedIPV6Config_Unassign(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "baz" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
  ipv6   = true
}

resource "digitalocean_reserved_ipv6" "foobar" {
  region_slug = "nyc3"
}`, name)
}
