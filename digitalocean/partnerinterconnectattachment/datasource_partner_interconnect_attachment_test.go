package partnerinterconnectattachment_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanPartnerInterconnectAttachment_ByID(t *testing.T) {
	var partnerInterconnectAttachment godo.PartnerInterconnectAttachment
	partnerInterconnectAttachmentName := acceptance.RandomTestName()
	vpcName1 := acceptance.RandomTestName()
	vpcName2 := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanPartnerInterconnectAttachmentConfig_Basic, vpcName1, vpcName2, partnerInterconnectAttachmentName)
	dataSourceConfig := `
data "digitalocean_partner_interconnect_attachment" "foobar" {
	id = digitalocean_partner_interconnect_attachment.foobar.id
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
					testAccCheckDigitalOceanPartnerInterconnectAttachmentExists("data.digitalocean_partner_interconnect_attachment.foobar", &partnerInterconnectAttachment),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_interconnect_attachment.foobar", "name", partnerInterconnectAttachmentName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_interconnect_attachment.foobar", "connection_bandwidth_in_mbps", "100"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_interconnect_attachment.foobar", "region", "nyc"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_interconnect_attachment.foobar", "naas_provider", "MEGAPORT"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_interconnect_attachment.foobar", "vpc_ids.#", "2"),
					resource.TestCheckResourceAttrPair(
						"data.digitalocean_partner_interconnect_attachment.foobar", "vpc_ids.0", "digitalocean_vpc.vpc1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.digitalocean_partner_interconnect_attachment.foobar", "vpc_ids.1", "digitalocean_vpc.vpc2", "id"),
					resource.TestCheckResourceAttrPair(
						"data.digitalocean_partner_interconnect_attachment.foobar", "bgp.0.local_router_ip", "169.254.0.1/29", "local_router_ip"),
					resource.TestCheckResourceAttrPair(
						"data.digitalocean_partner_interconnect_attachment.foobar", "bgp.0.peer_router_asn", "133937", "peer_router_asn"),
					resource.TestCheckResourceAttrPair(
						"data.digitalocean_partner_interconnect_attachment.foobar", "bgp.0.peer_router_ip", "169.254.0.6/29", "peer_router_ip"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_partner_interconnect_attachment.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_partner_interconnect_attachment.foobar", "state"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanPartnerInterconnectAttachment_ByName(t *testing.T) {
	var partnerInterconnectAttachment godo.PartnerInterconnectAttachment
	partnerInterconnectAttachmentName := acceptance.RandomTestName()
	vpcName1 := acceptance.RandomTestName()
	vpcName2 := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanPartnerInterconnectAttachmentConfig_Basic, vpcName1, vpcName2, partnerInterconnectAttachmentName)
	dataSourceConfig := `
data "digitalocean_partner_interconnect_attachment" "foobar" {
	name = digitalocean_partner_interconnect_attachment.foobar.name
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
					testAccCheckDigitalOceanPartnerInterconnectAttachmentExists("data.digitalocean_partner_interconnect_attachment.foobar", &partnerInterconnectAttachment),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_interconnect_attachment.foobar", "name", partnerInterconnectAttachmentName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_interconnect_attachment.foobar", "connection_bandwidth_in_mbps", "100"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_interconnect_attachment.foobar", "region", "nyc"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_interconnect_attachment.foobar", "naas_provider", "MEGAPORT"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_interconnect_attachment.foobar", "vpc_ids.#", "2"),
					resource.TestCheckResourceAttrPair(
						"data.digitalocean_partner_interconnect_attachment.foobar", "vpc_ids.0", "digitalocean_vpc.vpc1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.digitalocean_partner_interconnect_attachment.foobar", "vpc_ids.1", "digitalocean_vpc.vpc2", "id"),
					resource.TestCheckResourceAttrPair(
						"data.digitalocean_partner_interconnect_attachment.foobar", "bgp.0.local_router_ip", "169.254.0.1/29", "local_router_ip"),
					resource.TestCheckResourceAttrPair(
						"data.digitalocean_partner_interconnect_attachment.foobar", "bgp.0.peer_router_asn", "133937", "peer_router_asn"),
					resource.TestCheckResourceAttrPair(
						"data.digitalocean_partner_interconnect_attachment.foobar", "bgp.0.peer_router_ip", "169.254.0.6/29", "peer_router_ip"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_partner_interconnect_attachment.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_partner_interconnect_attachment.foobar", "state"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanPartnerInterconnectAttachment_ExpectErrors(t *testing.T) {
	partnenrInterconnectAttachmentName := acceptance.RandomTestName()
	partnerInterconnectAttachmentNotExists := fmt.Sprintf(testAccCheckDataSourceDigitalOceanPartnerInterconnectAttachmentConfig_DoesNotExist, partnenrInterconnectAttachmentName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      partnerInterconnectAttachmentNotExists,
				ExpectError: regexp.MustCompile(`Error retrieving Partner Interconnect Attachment`),
			},
		},
	})
}

const testAccCheckDataSourceDigitalOceanPartnerInterconnectAttachmentConfig_Basic = `
resource "digitalocean_vpc" "vpc1" {
  name   = "%s"
  region = "nyc3"
}

resource "digitalocean_vpc" "vpc2" {
  name   = "%s"
  region = "nyc3"
}

resource "digitalocean_partner_interconnect_attachment" "foobar" {
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
  }
}
`

const testAccCheckDataSourceDigitalOceanPartnerInterconnectAttachmentConfig_DoesNotExist = `
data "digitalocean_partner_interconnect_attachment" "foobar" {
	id = "%s"
}
`
