package partnernetworkconnect_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanPartnerNetworkConnect_ByID(t *testing.T) {
	var partnerNetworkConnect godo.PartnerNetworkConnect
	partnerNetworkConnectName := acceptance.RandomTestName()
	vpcName1 := acceptance.RandomTestName()
	vpcName2 := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanPartnerNetworkConnectConfig_Basic, vpcName1, vpcName2, partnerNetworkConnectName)
	dataSourceConfig := `
data "digitalocean_partner_network_connect" "foobar" {
  id = digitalocean_partner_network_connect.foobar.id
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanPartnerNetworkConnectExists("data.digitalocean_partner_network_connect.foobar", &partnerNetworkConnect),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_network_connect.foobar", "name", partnerNetworkConnectName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_network_connect.foobar", "connection_bandwidth_in_mbps", "100"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_network_connect.foobar", "region", "nyc"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_network_connect.foobar", "naas_provider", "MEGAPORT"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_network_connect.foobar", "vpc_ids.#", "2"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_network_connect.foobar", "bgp.0.local_router_ip", "169.254.0.1/29"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_network_connect.foobar", "bgp.0.peer_router_asn", "133937"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_network_connect.foobar", "bgp.0.peer_router_ip", "169.254.0.6/29"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_partner_network_connect.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_partner_network_connect.foobar", "state"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanPartnerNetworkConnect_ByName(t *testing.T) {
	var partnerNetworkConnect godo.PartnerNetworkConnect
	partnerNetworkConnectName := acceptance.RandomTestName()
	vpcName1 := acceptance.RandomTestName()
	vpcName2 := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanPartnerNetworkConnectConfig_Basic, vpcName1, vpcName2, partnerNetworkConnectName)
	dataSourceConfig := `
data "digitalocean_partner_network_connect" "foobar" {
  name = digitalocean_partner_network_connect.foobar.name
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanPartnerNetworkConnectExists("data.digitalocean_partner_network_connect.foobar", &partnerNetworkConnect),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_network_connect.foobar", "name", partnerNetworkConnectName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_network_connect.foobar", "connection_bandwidth_in_mbps", "100"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_network_connect.foobar", "region", "nyc"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_network_connect.foobar", "naas_provider", "MEGAPORT"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_network_connect.foobar", "vpc_ids.#", "2"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_network_connect.foobar", "bgp.0.local_router_ip", "169.254.0.1/29"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_network_connect.foobar", "bgp.0.peer_router_asn", "133937"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_network_connect.foobar", "bgp.0.peer_router_ip", "169.254.0.6/29"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_partner_network_connect.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_partner_network_connect.foobar", "state"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanPartnerNetworkConnect_ExpectErrors(t *testing.T) {
	partnenrNetworkConnectName := acceptance.RandomTestName()
	partnerNetworkConnectNotExists := fmt.Sprintf(testAccCheckDataSourceDigitalOceanPartnerNetworkConnectConfig_DoesNotExist, partnenrNetworkConnectName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      partnerNetworkConnectNotExists,
				ExpectError: regexp.MustCompile(`no Partner Network Connect found with name`),
			},
		},
	})
}

const testAccCheckDataSourceDigitalOceanPartnerNetworkConnectConfig_Basic = `
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
  depends_on = [
    digitalocean_vpc.vpc1,
    digitalocean_vpc.vpc2
  ]
  bgp {
    local_router_ip = "169.254.0.1/29"
    peer_router_asn = 133937
    peer_router_ip  = "169.254.0.6/29"
    auth_key        = "BGPAu7hK3y!"
  }
}
`

const testAccCheckDataSourceDigitalOceanPartnerNetworkConnectConfig_DoesNotExist = `
data "digitalocean_partner_network_connect" "foobar" {
  name = "%s"
}
`
