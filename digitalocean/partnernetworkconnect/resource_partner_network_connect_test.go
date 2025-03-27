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

func TestAccDigitalOceanPartnerAttachment_Basic(t *testing.T) {
	var partnerAttachment godo.PartnerAttachment

	vpc1Name := acceptance.RandomTestName()
	vpc2Name := acceptance.RandomTestName()
	vpc3Name := acceptance.RandomTestName()
	partnerAttachmentName := acceptance.RandomTestName()
	partnerAttachmentCreateConfig := fmt.Sprintf(testAccCheckDigitalOceanPartnerAttachmentConfig_Basic, vpc1Name, vpc2Name, partnerAttachmentName)

	updatePartnerAttachmentName := acceptance.RandomTestName()
	partnerAttachmentUpdateConfig := fmt.Sprintf(testAccCheckDigitalOceanPartnerAttachmentConfig_Basic, vpc1Name, vpc2Name, updatePartnerAttachmentName)
	partnerAttachmentVPCUpdateConfig := fmt.Sprintf(testAccCheckDigitalOceanPartnerAttachmentConfig_VPCUpdate, vpc1Name, vpc2Name, vpc3Name, updatePartnerAttachmentName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanPartnerAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: partnerAttachmentCreateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanPartnerAttachmentExists("digitalocean_partner_attachment.foobar", &partnerAttachment),
					resource.TestCheckResourceAttr(
						"digitalocean_partner_attachment.foobar", "name", partnerAttachmentName),
					resource.TestCheckResourceAttr(
						"digitalocean_partner_attachment.foobar", "vpc_ids.#", "2"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_partner_attachment.foobar", "created_at"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_partner_attachment.foobar", "state"),
				),
			},
			{
				Config: partnerAttachmentUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanPartnerAttachmentExists("digitalocean_partner_attachment.foobar", &partnerAttachment),
					resource.TestCheckResourceAttr(
						"digitalocean_partner_attachment.foobar", "name", updatePartnerAttachmentName),
				),
			},
			{
				Config: partnerAttachmentVPCUpdateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanPartnerAttachmentExists("digitalocean_partner_attachment.foobar", &partnerAttachment),
					resource.TestCheckResourceAttr(
						"digitalocean_partner_attachment.foobar", "vpc_ids.#", "3"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanPartnerAttachmentDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_partner_attachment" {
			continue
		}

		_, _, err := client.PartnerAttachment.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Partner Attachment still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanPartnerAttachmentExists(n string, partnerAttachment *godo.PartnerAttachment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Partner Attachment not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Partner Attachment ID is not set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		foundPartnerAttachment, _, err := client.PartnerAttachment.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("error fetching Partner Attachment (%s): %s", rs.Primary.ID, err)
		}

		*partnerAttachment = *foundPartnerAttachment

		return nil
	}
}

const testAccCheckDigitalOceanPartnerAttachmentConfig_Basic = `
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

const testAccCheckDigitalOceanPartnerAttachmentConfig_VPCUpdate = `
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
resource "digitalocean_partner_attachment" "foobar" {
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
