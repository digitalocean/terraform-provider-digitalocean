package firewall_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanFirewall_Basic(t *testing.T) {
	fwDataConfig := `
data "digitalocean_firewall" "foobar" {
  firewall_id = digitalocean_firewall.foobar.id
}`

	var firewall godo.Firewall
	fwName := acceptance.RandomTestName()

	fwCreateConfig := fmt.Sprintf(testAccDigitalOceanFirewallConfig_OnlyInbound(fwName))
	updatedFWCreateConfig := testAccDigitalOceanFirewallConfig_OnlyMultipleInbound(fwName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fwCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanFirewallExists("digitalocean_firewall.foobar", &firewall),
				),
			},
			{
				Config: fwCreateConfig + fwDataConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_firewall.foobar", "name", "foobar-"+fwName),
					resource.TestCheckResourceAttrPair("digitalocean_firewall.foobar", "id",
						"data.digitalocean_firewall.foobar", "firewall_id"),
					resource.TestCheckResourceAttrPair("digitalocean_firewall.foobar", "droplet_ids",
						"data.digitalocean_firewall.foobar", "droplet_ids"),
					resource.TestCheckResourceAttrPair("digitalocean_firewall.foobar", "inbound_rule",
						"data.digitalocean_firewall.foobar", "inbound_rule"),
					resource.TestCheckResourceAttrPair("digitalocean_firewall.foobar", "outbound_rule",
						"data.digitalocean_firewall.foobar", "outbound_rule"),
					resource.TestCheckResourceAttrPair("digitalocean_firewall.foobar", "status",
						"data.digitalocean_firewall.foobar", "status"),
					resource.TestCheckResourceAttrPair("digitalocean_firewall.foobar", "created_at",
						"data.digitalocean_firewall.foobar", "created_at"),
					resource.TestCheckResourceAttrPair("digitalocean_firewall.foobar", "pending_changes",
						"data.digitalocean_firewall.foobar", "pending_changes"),
					resource.TestCheckResourceAttrPair("digitalocean_firewall.foobar", "tags",
						"data.digitalocean_firewall.foobar", "tags"),
					resource.TestCheckResourceAttr("data.digitalocean_firewall.foobar", "inbound_rule.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("data.digitalocean_firewall.foobar", "inbound_rule.0.port_range", "22"),
					resource.TestCheckResourceAttr("data.digitalocean_firewall.foobar", "inbound_rule.0.source_addresses.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("data.digitalocean_firewall.foobar", "inbound_rule.0.source_addresses.1", "::/0"),
				),
			},
			{
				Config: updatedFWCreateConfig,
			},
			{
				Config: updatedFWCreateConfig + fwDataConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_firewall.foobar", "inbound_rule.0.protocol", "tcp"),
					resource.TestCheckResourceAttr("data.digitalocean_firewall.foobar", "inbound_rule.0.port_range", "22"),
					resource.TestCheckResourceAttr("data.digitalocean_firewall.foobar", "inbound_rule.0.source_addresses.0", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("data.digitalocean_firewall.foobar", "inbound_rule.0.source_addresses.1", "::/0"),
					resource.TestCheckResourceAttr("data.digitalocean_firewall.foobar", "inbound_rule.1.protocol", "tcp"),
					resource.TestCheckResourceAttr("data.digitalocean_firewall.foobar", "inbound_rule.1.port_range", "80"),
					resource.TestCheckResourceAttr("data.digitalocean_firewall.foobar", "inbound_rule.1.source_addresses.0", "1.2.3.0/24"),
					resource.TestCheckResourceAttr("data.digitalocean_firewall.foobar", "inbound_rule.1.source_addresses.1", "2002::/16"),
				),
			},
		},
	})
}
