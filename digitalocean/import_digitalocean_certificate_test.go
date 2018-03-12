package digitalocean

import (
	"testing"

	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDigitalOceanCertificate_importBasic(t *testing.T) {
	resourceName := "digitalocean_certificate.foobar"
	rInt := acctest.RandInt()
	leafCertMaterial, privateKeyMaterial, err := acctest.RandTLSCert("Acme Co")
	if err != nil {
		t.Fatalf("Cannot generate test TLS certificate: %s", err)
	}
	rootCertMaterial, _, err := acctest.RandTLSCert("Acme Go")
	if err != nil {
		t.Fatalf("Cannot generate test TLS certificate: %s", err)
	}
	certChainMaterial := fmt.Sprintf("%s\n%s", strings.TrimSpace(rootCertMaterial), leafCertMaterial)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanCertificateConfig_basic(rInt, privateKeyMaterial, leafCertMaterial, certChainMaterial),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"certificate_chain"}, //we ignore the IP Address as we do not set to state
			},
		},
	})
}
