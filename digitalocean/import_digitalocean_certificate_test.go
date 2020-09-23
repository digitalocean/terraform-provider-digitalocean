package digitalocean

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanCertificate_importBasic(t *testing.T) {
	resourceName := "digitalocean_certificate.foobar"
	rInt := acctest.RandInt()

	privateKeyMaterial, leafCertMaterial, certChainMaterial := generateTestCertMaterial(t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanCertificateConfig_basic(rInt, privateKeyMaterial, leafCertMaterial, certChainMaterial),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"certificate_chain", "leaf_certificate", "private_key"}, // We ignore these as they are not returned by the API

			},
		},
	})
}
