package partnerinterconnectattachment_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanPartnerInterconnectAttachment_importBasic(t *testing.T) {
	resourceName := "digitalocean_partner_interconnect_attachment.foobar"
	vpc1Name := acceptance.RandomTestName()
	vpc2Name := acceptance.RandomTestName()
	partnerInterconnectAttachmentName := acceptance.RandomTestName()
	partnerInterconnectAttachmentCreateConfig := fmt.Sprintf(testAccCheckDigitalOceanPartnerInterconnectAttachmentConfig_Basic, vpc1Name, vpc2Name, partnerInterconnectAttachmentName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanPartnerInterconnectAttachmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: partnerInterconnectAttachmentCreateConfig,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"bgp.0.auth_key"},
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       false,
				ImportStateVerifyIgnore: []string{"bgp.0.auth_key"},
				ImportStateId:           "123abc",
				ExpectError:             regexp.MustCompile(`(Please verify the ID is correct|Cannot import non-existent remote object)`),
			},
		},
	})
}
