package partnernetworkconnect_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanPartnerNetworkConnect_Basic(t *testing.T) {
	var partnerNetworkConnect godo.PartnerNetworkConnect

	vpc1Name := acceptance.RandomTestName()
	vpc2Name := acceptance.RandomTestName()
	vpc3Name := acceptance.RandomTestName()
	partnerNetworkConnectName := acceptance.RandomTestName()
	partnerNetworkConnectCreateConfig := fmt.Sprintf(testAccCheckDigitalOceanPartnerNetworkConnectConfig_Basic, vpc1Name, vpc2Name, partnerNetworkConnectName)

	updatePartnerNetworkConnectName := acceptance.RandomTestName()
	partnerNetworkConnectUpdateConfig := fmt.Sprintf(testAccCheckDigitalOceanPartnerNetworkConnectConfig_Basic, vpc1Name, vpc2Name, updatePartnerNetworkConnectName)
	partnerNetworkConnectVPCUpdateConfig := fmt.Sprintf(testAccCheckDigitalOceanPartnerNetworkConnectConfig_VPCUpdate, vpc1Name, vpc2Name, vpc3Name, updatePartnerNetworkConnectName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanPartnerNetworkConnectDestroy,
		Steps: []resource.TestStep{
			{
				Config: partnerNetworkConnectCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanPartnerNetworkConnectExists("digitalocean_partner_network_connect.foobar", &partnerNetworkConnect),
					resource.TestCheckResourceAttr(
						"digitalocean_partner_network_connect.foobar", "name", partnerNetworkConnectName),
					resource.TestCheckResourceAttr(
						"digitalocean_partner_network_connect.foobar", "vpc_ids.#", "2"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_partner_network_connect.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_partner_network_connect.foobar", "state"),
				),
			},
			{
				Config: partnerNetworkConnectUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanPartnerNetworkConnectExists("digitalocean_partner_network_connect.foobar", &partnerNetworkConnect),
					resource.TestCheckResourceAttr(
						"digitalocean_partner_network_connect.foobar", "name", updatePartnerNetworkConnectName),
				),
			},
			{
				Config: partnerNetworkConnectVPCUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanPartnerNetworkConnectExists("digitalocean_partner_network_connect.foobar", &partnerNetworkConnect),
					resource.TestCheckResourceAttr(
						"digitalocean_partner_network_connect.foobar", "vpc_ids.#", "3"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanPartnerNetworkConnectDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_partner_network_connect" {
			continue
		}

		_, _, err := client.PartnerNetworkConnect.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Partner Network Connect still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanPartnerNetworkConnectExists(n string, partnerNetworkConnect *godo.PartnerNetworkConnect) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Partner Network Connect not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Partner Network Connect ID is not set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		foundPartnerNetworkConnect, _, err := client.PartnerNetworkConnect.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching Partner Network Connect (%s): %s", rs.Primary.ID, err)
		}

		*partnerNetworkConnect = *foundPartnerNetworkConnect

		return nil
	}
}

const testAccCheckDigitalOceanPartnerNetworkConnectConfig_Basic = `
resource "digitalocean_vpc" "vpc1" {
  name   = "%s"
  region = "nyc3"
}
resource "digitalocean_vpc" "vpc2" {
  name   = "%s"
  region = "nyc3"
}
resource "digitalocean_partner_network_connect" "foobar" {
  name                         = "%s"
  connection_bandwidth_in_mbps = 100
  region                       = "nyc"
  naas_provider                = "MEGAPORT"
  vpc_ids = [
    digitalocean_vpc.vpc1.id,
    digitalocean_vpc.vpc2.id
  ]
  bgp {
    local_router_ip = "169.254.0.1/29"
    peer_router_asn = 133937
    peer_router_ip  = "169.254.0.6/29"
    auth_key        = "BGPAu7hK3y!"
  }
}
`

const testAccCheckDigitalOceanPartnerNetworkConnectConfig_VPCUpdate = `
resource "digitalocean_vpc" "vpc1" {
  name   = "%s"
  region = "nyc3"
}
resource "digitalocean_vpc" "vpc2" {
  name   = "%s"
  region = "nyc3"
}
resource "digitalocean_vpc" "vpc3" {
  name   = "%s"
  region = "nyc3"
}
resource "digitalocean_partner_network_connect" "foobar" {
  name                         = "%s"
  connection_bandwidth_in_mbps = 100
  region                       = "nyc"
  naas_provider                = "MEGAPORT"
  vpc_ids = [
    digitalocean_vpc.vpc1.id,
    digitalocean_vpc.vpc2.id,
    digitalocean_vpc.vpc3.id
  ]
  bgp {
    local_router_ip = "169.254.0.1/29"
    peer_router_asn = 133937
    peer_router_ip  = "169.254.0.6/29"
    auth_key        = "BGPAu7hK3y!"
  }
}
`
