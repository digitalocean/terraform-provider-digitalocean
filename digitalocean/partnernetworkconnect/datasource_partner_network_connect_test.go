package partnernetworkconnect_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
)

func TestAccDataSourceDigitalOceanPartnerAttachment_ByID(t *testing.T) {
	var partnerAttachment godo.PartnerAttachment
	partnerAttachmentName := acceptance.RandomTestName()
	vpcName1 := acceptance.RandomTestName()
	vpcName2 := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanPartnerAttachmentConfig_Basic, vpcName1, vpcName2, partnerAttachmentName)
	dataSourceConfig := `
data "digitalocean_partner_attachment" "foobar" {
  id = digitalocean_partner_attachment.foobar.id
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
					testAccCheckDigitalOceanPartnerAttachmentExists("data.digitalocean_partner_attachment.foobar", &partnerAttachment),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_attachment.foobar", "name", partnerAttachmentName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_attachment.foobar", "connection_bandwidth_in_mbps", "1000"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_attachment.foobar", "region", "nyc"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_attachment.foobar", "naas_provider", "MEGAPORT"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_attachment.foobar", "redundancy_zone", "MEGAPORT_RED"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_attachment.foobar", "vpc_ids.#", "2"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_attachment.foobar", "bgp.0.local_router_ip", "169.254.100.1/29"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_attachment.foobar", "bgp.0.peer_router_asn", "133937"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_attachment.foobar", "bgp.0.peer_router_ip", "169.254.100.6/29"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_partner_attachment.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_partner_attachment.foobar", "state"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_partner_attachment.foobar", "parent_uuid"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_partner_attachment.foobar", "children"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanPartnerAttachment_ByName(t *testing.T) {
	var partnerAttachment godo.PartnerAttachment
	partnerAttachmentName := acceptance.RandomTestName()
	vpcName1 := acceptance.RandomTestName()
	vpcName2 := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanPartnerAttachmentConfig_Basic, vpcName1, vpcName2, partnerAttachmentName)
	dataSourceConfig := `
data "digitalocean_partner_attachment" "foobar" {
  name = digitalocean_partner_attachment.foobar.name
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
					testAccCheckDigitalOceanPartnerAttachmentExists("data.digitalocean_partner_attachment.foobar", &partnerAttachment),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_attachment.foobar", "name", partnerAttachmentName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_attachment.foobar", "connection_bandwidth_in_mbps", "1000"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_attachment.foobar", "region", "nyc"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_attachment.foobar", "naas_provider", "MEGAPORT"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_attachment.foobar", "redundancy_zone", "MEGAPORT_RED"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_attachment.foobar", "vpc_ids.#", "2"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_attachment.foobar", "bgp.0.local_router_ip", "169.254.100.1/29"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_attachment.foobar", "bgp.0.peer_router_asn", "133937"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_partner_attachment.foobar", "bgp.0.peer_router_ip", "169.254.100.6/29"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_partner_attachment.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_partner_attachment.foobar", "state"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_partner_attachment.foobar", "parent_uuid"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_partner_attachment.foobar", "children"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanPartnerAttachment_ExpectErrors(t *testing.T) {
	partnenrAttachmentName := acceptance.RandomTestName()
	partnerAttachmentNotExists := fmt.Sprintf(testAccCheckDataSourceDigitalOceanPartnerAttachmentConfig_DoesNotExist, partnenrAttachmentName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      partnerAttachmentNotExists,
				ExpectError: regexp.MustCompile(`no Partner Attachment found with name`),
			},
		},
	})
}

const testAccCheckDataSourceDigitalOceanPartnerAttachmentConfig_Basic = `
resource "digitalocean_vpc" "vpc1" {
  name   = "%s"
  region = "nyc3"
}
resource "digitalocean_vpc" "vpc2" {
  name   = "%s"
  region = "nyc3"
}
resource "digitalocean_partner_attachment" "foobar" {
  name                         = "%s"
  connection_bandwidth_in_mbps = 1000
  region                       = "nyc"
  naas_provider                = "MEGAPORT"
  redundancy_zone              = "MEGAPORT_RED"
  vpc_ids = [
    digitalocean_vpc.vpc1.id,
    digitalocean_vpc.vpc2.id
  ]
  depends_on = [
    digitalocean_vpc.vpc1,
    digitalocean_vpc.vpc2
  ]
  bgp {
    local_router_ip = "169.254.100.1/29"
    peer_router_asn = 133937
    peer_router_ip  = "169.254.100.6/29"
    auth_key        = "BGPAu7hK3y!"
  }
  parent_uuid = "00000000-0000-0000-0000-000000000000"
  children = "11111111-1111-1111-1111-111111111111"
}
`

const testAccCheckDataSourceDigitalOceanPartnerAttachmentConfig_DoesNotExist = `
data "digitalocean_partner_attachment" "foobar" {
  name = "%s"
}
`
