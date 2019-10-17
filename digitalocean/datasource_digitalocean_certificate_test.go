package digitalocean

import (
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceDigitalOceanCertificate_Basic(t *testing.T) {
	var certificate godo.Certificate
	name := fmt.Sprintf("certificate-%s", acctest.RandString(10))

	privateKeyMaterial, leafCertMaterial, certChainMaterial := generateTestCertMaterial(t)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDigitalOceanCertificateConfig_basic(name, privateKeyMaterial, leafCertMaterial, certChainMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanCertificateExists("data.digitalocean_certificate.foobar", &certificate),
					resource.TestCheckResourceAttr(
						"data.digitalocean_certificate.foobar", "id", name),
					resource.TestCheckResourceAttr(
						"data.digitalocean_certificate.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"data.digitalocean_certificate.foobar", "type", "custom"),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanCertificateExists(n string, certificate *godo.Certificate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No certificate ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		foundCertificate, err := findCertificateByName(client, rs.Primary.ID)
		if err != nil {
			return err
		}

		*certificate = *foundCertificate

		return nil
	}
}

func testAccCheckDataSourceDigitalOceanCertificateConfig_basic(name, privateKeyMaterial, leafCert, certChain string) string {
	return fmt.Sprintf(`
resource "digitalocean_certificate" "foo" {
  name = "%s"
  private_key = <<EOF
%s
EOF
  leaf_certificate = <<EOF
%s
EOF
  certificate_chain = <<EOF
%s
EOF
}

data "digitalocean_certificate" "foobar" {
  name = "${digitalocean_certificate.foo.name}"
}
`, name, privateKeyMaterial, leafCert, certChain)
}
