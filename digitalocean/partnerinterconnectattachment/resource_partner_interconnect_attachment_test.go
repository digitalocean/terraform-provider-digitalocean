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
	partnerInterconnectAttachmentName := acceptance.RandomTestName()
	partnerInterconnectAttachmentCreateConfig := fmt.Sprintf(testAccCheckDigitalOceanPartnerInterconnectAttachmentConfig_Basic, partnerInterconnectAttachmentName)

	updatePartnerInterconnectAttachmentName := acceptance.RandomTestName()
	partnerInterconnectAttachmentUpdateConfig := fmt.Sprintf(testAccCheckDigitalOceanPartnerInterconnectAttachmentConfig_Basic, updatePartnerInterconnectAttachmentName)
	partnerInterconnectAttachmentVPCUpdateConfig := fmt.Sprintf(testAccCheckDigitalOceanPartnerInterconnectAttachmentConfig_VPCUpdate, updatePartnerInterconnectAttachmentName)

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
					resource.TestCheckResourceAttrPair(
						"digitalocean_partner_interconnect_attachment.foobar", "vpc_ids.0", "digitalocean_vpc.vpc1", "id"),
					resource.TestCheckResourceAttrPair(
						"digitalocean_partner_interconnect_attachment.foobar", "vpc_ids.1", "digitalocean_vpc.vpc2", "id"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_partner_interconnect_attachment.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_partner_interconnect_attachment.foobar", "status"),
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
                        "digitalocean_partner_interconnect_attachment.foobar", "vpc_ids.#", "2"),
                    resource.TestCheckResourceAttrPair(
                        "digitalocean_partner_interconnect_attachment.foobar", "vpc_ids.0", "digitalocean_vpc.vpc1", "id"),
                    resource.TestCheckResourceAttrPair(
                        "digitalocean_partner_interconnect_attachment.foobar", "vpc_ids.1", "digitalocean_vpc.vpc3", "id"),
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
  name   = "vpc1"
  region = "nyc3"
}

resource "digitalocean_vpc" "vpc2" {
  name   = "vpc2"
  region = "nyc3"
}

resource "digitalocean_partner_interconnect_attachment" "foobar" {
  name = "%s"
  connection_bandwidth_in_mbps = 100
  region = "nyc"
  naas_provider = "megaport"
  vpc_ids = [
    digitalocean_vpc.vpc1.id,
	digitalocean_vpc.vpc2.id
  ]
  bgp {
	local_asn = 64532
	local_router_ip = "169.254.0.1/29"
	peer_asn = 133937
	peer_router_ip = "169.254.0.6/29"
  }
}
`

const testAccCheckDigitalOceanPartnerInterconnectAttachmentConfig_VPCUpdate = `
resource "digitalocean_vpc" "vpc1" {
  name   = "vpc1"
  region = "nyc3"
}

resource "digitalocean_vpc" "vpc3" {
  name   = "vpc3"
  region = "nyc3"
}

resource "digitalocean_partner_interconnect_attachment" "foobar" {
  name = "%s"
  connection_bandwidth_in_mbps = 100
  region = "nyc"
  naas_provider = "megaport"
  vpc_ids = [
    digitalocean_vpc.vpc1.id,
    digitalocean_vpc.vpc3.id
  ]
  bgp {
    local_asn = 64532
    local_router_ip = "169.254.0.1/29"
    peer_asn = 133937
    peer_router_ip = "169.254.0.6/29"
  }
}
`
