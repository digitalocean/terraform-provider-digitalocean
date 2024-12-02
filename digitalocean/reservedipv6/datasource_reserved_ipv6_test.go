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

func TestAccDataSourceDigitalOceanReservedIPV6_Basic(t *testing.T) {
	var reservedIPv6 godo.ReservedIPV6

	expectedURNRegex, _ := regexp.Compile(`do:reservedipv6:/^[0-9A-F]{8}-[0-9A-F]{4}-[4][0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$/i`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDigitalOceanReservedIPConfig_Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanReservedIPV6Exists("data.digitalocean_reserved_ipv6.foobar", &reservedIPv6),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_reserved_ipv6.foobar", "ip"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_reserved_ipv6.foobar", "region_slug", "nyc3"),
					resource.TestMatchResourceAttr("data.digitalocean_reserved_ipv6.foobar", "urn", expectedURNRegex),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanReservedIPV6_FindsReservedIP(t *testing.T) {
	var reservedIPv6 godo.ReservedIPV6

	expectedURNRegex, _ := regexp.Compile(`do:reservedipv6:/^[0-9A-F]{8}-[0-9A-F]{4}-[4][0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$/i`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDigitalOceanReservedIPConfig_FindsFloatingIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanReservedIPV6Exists("data.digitalocean_reserved_ipv6.foobar", &reservedIPv6),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_reserved_ipv6.foobar", "ip"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_reserved_ipv6.foobar", "region_slug", "nyc3"),
					resource.TestMatchResourceAttr("data.digitalocean_reserved_ipv6.foobar", "urn", expectedURNRegex),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanReservedIPV6Exists(n string, reservedIPv6 *godo.ReservedIPV6) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No reserved IPv6 ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		foundReservedIP, _, err := client.ReservedIPV6s.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundReservedIP.IP != rs.Primary.ID {
			return fmt.Errorf("reserved IPv6 not found")
		}

		*reservedIPv6 = *foundReservedIP

		return nil
	}
}

const testAccCheckDataSourceDigitalOceanReservedIPConfig_FindsFloatingIP = `
resource "digitalocean_reserved_ipv6" "foo" {
  region_slug = "nyc3"
}

data "digitalocean_reserved_ipv6" "foobar" {
  ip = digitalocean_reserved_ipv6.foo.ip
}`

const testAccCheckDataSourceDigitalOceanReservedIPConfig_Basic = `
resource "digitalocean_reserved_ipv6" "foo" {
  region_slug = "nyc3"
}

data "digitalocean_reserved_ipv6" "foobar" {
  ip = digitalocean_reserved_ipv6.foo.ip
}`
