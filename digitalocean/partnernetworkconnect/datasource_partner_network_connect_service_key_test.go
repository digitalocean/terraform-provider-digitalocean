package partnernetworkconnect_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanPartnerAttachmentServiceKey_Found(t *testing.T) {
	//var serviceKey godo.ServiceKey

	partnerAttachmentName := acceptance.RandomTestName()
	vpcName1 := acceptance.RandomTestName()
	vpcName2 := acceptance.RandomTestName()

	resourceConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanPartnerAttachmentConfig_Basic, vpcName1, vpcName2, partnerAttachmentName)
	dataSourceConfig := `
data "digitalocean_partner_attachment_service_key" "foobar" {
  attachment_id = digitalocean_partner_attachment.foobar.id
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
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_partner_attachment_service_key.foobar", "value",
					),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_partner_attachment_service_key.foobar", "state",
					),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_partner_attachment_service_key.foobar", "created_at",
					),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanPartnerAttachmentServiceKey_ExpectError(t *testing.T) {
	invalidID := "non-existent-id"

	// just the data source pointing at a bogus ID
	errorConfig := fmt.Sprintf(`
data "digitalocean_partner_attachment_service_key" "test" {
  attachment_id = "%s"
}`, invalidID)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      errorConfig,
				ExpectError: regexp.MustCompile(`error retrieving service key for partner attachment`),
			},
		},
	})
}
