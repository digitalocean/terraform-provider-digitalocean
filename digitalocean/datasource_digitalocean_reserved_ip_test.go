package digitalocean

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceDigitalOceanReservedIP_Basic(t *testing.T) {
	var reservedIP godo.ReservedIP

	expectedURNRegEx, _ := regexp.Compile(`do:reservedip:(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDigitalOceanReservedIPConfig_Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanReservedIPExists("data.digitalocean_reserved_ip.foobar", &reservedIP),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_reserved_ip.foobar", "ip_address"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_reserved_ip.foobar", "region", "nyc3"),
					resource.TestMatchResourceAttr("data.digitalocean_reserved_ip.foobar", "urn", expectedURNRegEx),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanReservedIP_FindsFloatingIP(t *testing.T) {
	var reservedIP godo.ReservedIP

	expectedURNRegEx, _ := regexp.Compile(`do:reservedip:(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDigitalOceanReservedIPConfig_FindsFloatingIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanReservedIPExists("data.digitalocean_reserved_ip.foobar", &reservedIP),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_reserved_ip.foobar", "ip_address"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_reserved_ip.foobar", "region", "nyc3"),
					resource.TestMatchResourceAttr("data.digitalocean_reserved_ip.foobar", "urn", expectedURNRegEx),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanReservedIPExists(n string, reservedIP *godo.ReservedIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No reserved IP ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		foundReservedIP, _, err := client.ReservedIPs.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundReservedIP.IP != rs.Primary.ID {
			return fmt.Errorf("reserved IP not found")
		}

		*reservedIP = *foundReservedIP

		return nil
	}
}

const testAccCheckDataSourceDigitalOceanReservedIPConfig_FindsFloatingIP = `
resource "digitalocean_floating_ip" "foo" {
  region = "nyc3"
}

data "digitalocean_reserved_ip" "foobar" {
  ip_address = digitalocean_floating_ip.foo.ip_address
}`

const testAccCheckDataSourceDigitalOceanReservedIPConfig_Basic = `
resource "digitalocean_reserved_ip" "foo" {
  region = "nyc3"
}

data "digitalocean_reserved_ip" "foobar" {
  ip_address = digitalocean_reserved_ip.foo.ip_address
}`
