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

func TestAccDataSourceDigitalOceanFloatingIp_Basic(t *testing.T) {
	var floatingIp godo.FloatingIP

	expectedURNRegEx, _ := regexp.Compile(`do:floatingip:(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDigitalOceanFloatingIpConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanFloatingIpExists("data.digitalocean_floating_ip.foobar", &floatingIp),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_floating_ip.foobar", "ip_address"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_floating_ip.foobar", "region", "nyc3"),
					resource.TestMatchResourceAttr("data.digitalocean_floating_ip.foobar", "urn", expectedURNRegEx),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanFloatingIpExists(n string, floatingIp *godo.FloatingIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No floating ip ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		foundFloatingIp, _, err := client.FloatingIPs.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundFloatingIp.IP != rs.Primary.ID {
			return fmt.Errorf("Floating ip not found")
		}

		*floatingIp = *foundFloatingIp

		return nil
	}
}

const testAccCheckDataSourceDigitalOceanFloatingIpConfig_basic = `
resource "digitalocean_floating_ip" "foo" {
  region = "nyc3"
}

data "digitalocean_floating_ip" "foobar" {
  ip_address = digitalocean_floating_ip.foo.ip_address
}`
