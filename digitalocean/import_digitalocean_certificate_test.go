package digitalocean

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDigitalOceanCertificate_importBasic(t *testing.T) {
	resourceName := "digitalocean_certificate.foobar"
	rInt := acctest.RandInt()

	wd, _ := os.Getwd()
	leafCert := wd + "/test-fixtures/terraform.cert"
	privateKeyMaterial := wd + "/test-fixtures/terraform.key"
	certChain := wd + "/test-fixtures/terraform-chain.cert"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanCertificateConfig_basic(rInt, privateKeyMaterial, leafCert, certChain),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
