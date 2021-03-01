package digitalocean

import (
	"context"
	"fmt"
	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanFirewall_Basic(t *testing.T) {
	fwDataConfig := `
data "digitalocean_firewall" "foobar" {
	firewall_id = digitalocean_firewall.foobar.firewall_id
	name = digitalocean_firewall.foobar.name
}`

	var firewall godo.Firewall
	fwName := randomTestName()

	fwCreateConfig := fmt.Sprintf(testAccDigitalOceanFirewallConfig_OnlyInbound(fwName))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fwCreateConfig,
			},
			{
				Config: fwCreateConfig + fwDataConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatasourceDigitalOceanFirewallExists("data.digitalocean_firewall.foobar", &firewall),
					resource.TestCheckResourceAttr("data.digitalocean_firewall.foobar", "name", "foobar-"+fwName),
				),
			},
		},
	})
}

func testAccCheckDatasourceDigitalOceanFirewallExists(resource string, firewall *godo.Firewall) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("not found: %s", resource)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no firewall ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		foundFirewall, _, err := client.Firewalls.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundFirewall.ID != rs.Primary.ID {
			return fmt.Errorf("firewall not found")
		}

		*firewall = *foundFirewall

		return nil
	}
}
