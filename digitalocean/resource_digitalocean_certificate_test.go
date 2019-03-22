package digitalocean

import (
	"bytes"
	"context"
	crand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func init() {
	resource.AddTestSweepers("digitalocean_certificate", &resource.Sweeper{
		Name: "digitalocean_certificate",
		F:    testSweepCertificate,
	})

}

func testSweepCertificate(region string) error {
	meta, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*CombinedConfig).godoClient()

	opt := &godo.ListOptions{PerPage: 200}
	certs, _, err := client.Certificates.List(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, c := range certs {
		if strings.HasPrefix(c.Name, "certificate-") {
			log.Printf("Destroying certificate %s", c.Name)

			if _, err := client.Certificates.Delete(context.Background(), c.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccDigitalOceanCertificate_Basic(t *testing.T) {
	var cert godo.Certificate
	rInt := acctest.RandInt()
	privateKeyMaterial, leafCertMaterial, certChainMaterial := generateTestCertMaterial(t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanCertificateConfig_basic(rInt, privateKeyMaterial, leafCertMaterial, certChainMaterial),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanCertificateExists("digitalocean_certificate.foobar", &cert),
					resource.TestCheckResourceAttr(
						"digitalocean_certificate.foobar", "name", fmt.Sprintf("certificate-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_certificate.foobar", "private_key", fmt.Sprintf("%s\n", privateKeyMaterial)),
					resource.TestCheckResourceAttr(
						"digitalocean_certificate.foobar", "leaf_certificate", fmt.Sprintf("%s\n", leafCertMaterial)),
					resource.TestCheckResourceAttr(
						"digitalocean_certificate.foobar", "certificate_chain", fmt.Sprintf("%s\n", certChainMaterial)),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanCertificateDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*CombinedConfig).godoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_certificate" {
			continue
		}

		_, _, err := client.Certificates.Get(context.Background(), rs.Primary.ID)

		if err != nil && !strings.Contains(err.Error(), "404") {
			return fmt.Errorf(
				"Error waiting for certificate (%s) to be destroyed: %s",
				rs.Primary.ID, err)
		}
	}

	return nil
}

func testAccCheckDigitalOceanCertificateExists(n string, cert *godo.Certificate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Certificate ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		c, _, err := client.Certificates.Get(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if c.ID != rs.Primary.ID {
			return fmt.Errorf("Certificate not found")
		}

		*cert = *c

		return nil
	}
}

func generateTestCertMaterial(t *testing.T) (string, string, string) {
	leafCertMaterial, privateKeyMaterial, err := randTLSCert("Acme Co", "example.com")
	if err != nil {
		t.Fatalf("Cannot generate test TLS certificate: %s", err)
	}
	rootCertMaterial, _, err := randTLSCert("Acme Go", "example.com")
	if err != nil {
		t.Fatalf("Cannot generate test TLS certificate: %s", err)
	}
	certChainMaterial := fmt.Sprintf("%s\n%s", strings.TrimSpace(rootCertMaterial), leafCertMaterial)

	return privateKeyMaterial, leafCertMaterial, certChainMaterial
}

// Based on Terraform's acctest.RandTLSCert, but allows for passing DNS name.
func randTLSCert(orgName string, dnsName string) (string, string, error) {
	template := &x509.Certificate{
		SerialNumber: big.NewInt(int64(acctest.RandInt())),
		Subject: pkix.Name{
			Organization: []string{orgName},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{dnsName},
	}

	privateKey, privateKeyPEM, err := genPrivateKey()
	if err != nil {
		return "", "", err
	}

	cert, err := x509.CreateCertificate(crand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", "", err
	}

	certPEM, err := pemEncode(cert, "CERTIFICATE")
	if err != nil {
		return "", "", err
	}

	return certPEM, privateKeyPEM, nil
}

func genPrivateKey() (*rsa.PrivateKey, string, error) {
	privateKey, err := rsa.GenerateKey(crand.Reader, 1024)
	if err != nil {
		return nil, "", err
	}

	privateKeyPEM, err := pemEncode(x509.MarshalPKCS1PrivateKey(privateKey), "RSA PRIVATE KEY")
	if err != nil {
		return nil, "", err
	}

	return privateKey, privateKeyPEM, nil
}

func pemEncode(b []byte, block string) (string, error) {
	var buf bytes.Buffer
	pb := &pem.Block{Type: block, Bytes: b}
	if err := pem.Encode(&buf, pb); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func testAccCheckDigitalOceanCertificateConfig_basic(rInt int, privateKeyMaterial, leafCert, certChain string) string {
	return fmt.Sprintf(`
resource "digitalocean_certificate" "foobar" {
  name = "certificate-%d"
  private_key = <<EOF
%s
EOF
  leaf_certificate = <<EOF
%s
EOF
  certificate_chain = <<EOF
%s
EOF
}`, rInt, privateKeyMaterial, leafCert, certChain)
}
