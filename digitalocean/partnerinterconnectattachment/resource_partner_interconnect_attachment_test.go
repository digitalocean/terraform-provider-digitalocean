package partnerinterconnectattachment_test

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

func TestAccDigitalOceanPartnerInterconnectAttachment_Basic(t *testing.T) {
	var partnerInterconnectAttachment godo.PartnerInterconnectAttachment

	vpc1Name := acceptance.RandomTestName()
	vpc2Name := acceptance.RandomTestName()
	vpc3Name := acceptance.RandomTestName()
	partnerInterconnectAttachmentName := acceptance.RandomTestName()
	partnerInterconnectAttachmentCreateConfig := fmt.Sprintf(testAccCheckDigitalOceanPartnerInterconnectAttachmentConfig_Basic, vpc1Name, vpc2Name, partnerInterconnectAttachmentName)

	updatePartnerInterconnectAttachmentName := acceptance.RandomTestName()
	partnerInterconnectAttachmentUpdateConfig := fmt.Sprintf(testAccCheckDigitalOceanPartnerInterconnectAttachmentConfig_Basic, vpc1Name, vpc2Name, updatePartnerInterconnectAttachmentName)
	partnerInterconnectAttachmentVPCUpdateConfig := fmt.Sprintf(testAccCheckDigitalOceanPartnerInterconnectAttachmentConfig_VPCUpdate, vpc1Name, vpc2Name, vpc3Name, updatePartnerInterconnectAttachmentName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanPartnerInterconnectAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: partnerInterconnectAttachmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanPartnerInterconnectAttachmentExists("digitalocean_partner_interconnect_attachment.foobar", &partnerInterconnectAttachment),
					resource.TestCheckResourceAttr(
						"digitalocean_partner_interconnect_attachment.foobar", "name", partnerInterconnectAttachmentName),
					resource.TestCheckResourceAttr(
						"digitalocean_partner_interconnect_attachment.foobar", "vpc_ids.#", "2"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_partner_interconnect_attachment.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_partner_interconnect_attachment.foobar", "state"),
				),
			},
			{
				Config: partnerInterconnectAttachmentUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanPartnerInterconnectAttachmentExists("digitalocean_partner_interconnect_attachment.foobar", &partnerInterconnectAttachment),
					resource.TestCheckResourceAttr(
						"digitalocean_partner_interconnect_attachment.foobar", "name", updatePartnerInterconnectAttachmentName),
				),
			},
			{
				Config: partnerInterconnectAttachmentVPCUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanPartnerInterconnectAttachmentExists("digitalocean_partner_interconnect_attachment.foobar", &partnerInterconnectAttachment),
					resource.TestCheckResourceAttr(
						"digitalocean_partner_interconnect_attachment.foobar", "vpc_ids.#", "3"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanPartnerInterconnectAttachmentDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_partner_interconnect_attachment" {
			continue
		}

		_, _, err := client.PartnerInterconnectAttachments.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Partner Interconnect Attachment still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanPartnerInterconnectAttachmentExists(n string, partnerInterconnectAttachment *godo.PartnerInterconnectAttachment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Partner Interconnect Attachment not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Partner Interconnect Attachment ID is not set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		foundPartnerInterconnectAttachment, _, err := client.PartnerInterconnectAttachments.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching Partner Interconnect Attachment (%s): %s", rs.Primary.ID, err)
		}

		*partnerInterconnectAttachment = *foundPartnerInterconnectAttachment

		return nil
	}
}

const testAccCheckDigitalOceanPartnerInterconnectAttachmentConfig_Basic = `
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
  bgp {
    local_router_ip = "169.254.0.1/29"
    peer_router_asn = 133937
    peer_router_ip  = "169.254.0.6/29"
    auth_key        = "BGPAu7hK3y!"
  }
}
`

const testAccCheckDigitalOceanPartnerInterconnectAttachmentConfig_VPCUpdate = `
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

resource "digitalocean_partner_interconnect_attachment" "foobar" {
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
