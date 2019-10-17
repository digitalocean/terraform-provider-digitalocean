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
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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
	name := fmt.Sprintf("certificate-%d", rInt)
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
						"digitalocean_certificate.foobar", "id", name),
					resource.TestCheckResourceAttr(
						"digitalocean_certificate.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"digitalocean_certificate.foobar", "private_key", HashString(fmt.Sprintf("%s\n", privateKeyMaterial))),
					resource.TestCheckResourceAttr(
						"digitalocean_certificate.foobar", "leaf_certificate", HashString(fmt.Sprintf("%s\n", leafCertMaterial))),
					resource.TestCheckResourceAttr(
						"digitalocean_certificate.foobar", "certificate_chain", HashString(fmt.Sprintf("%s\n", certChainMaterial))),
				),
			},
		},
	})
}

func TestAccDigitalOceanCertificate_ExpectedErrors(t *testing.T) {
	rInt := acctest.RandInt()
	privateKeyMaterial, leafCertMaterial, certChainMaterial := generateTestCertMaterial(t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanCertificateDestroy,
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckDigitalOceanCertificateConfig_customNoLeaf(rInt, privateKeyMaterial, certChainMaterial),
				ExpectError: regexp.MustCompile("`leaf_certificate` is required for when type is `custom` or empty"),
			},
			{
				Config:      testAccCheckDigitalOceanCertificateConfig_customNoKey(rInt, leafCertMaterial, certChainMaterial),
				ExpectError: regexp.MustCompile("`private_key` is required for when type is `custom` or empty"),
			},
			{
				Config:      testAccCheckDigitalOceanCertificateConfig_noDomains(rInt),
				ExpectError: regexp.MustCompile("`domains` is required for when type is `lets_encrypt`"),
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

		_, err := findCertificateByName(client, rs.Primary.ID)

		if err != nil && !strings.Contains(err.Error(), "not found") {
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

		c, err := findCertificateByName(client, rs.Primary.ID)
		if err != nil {
			return err
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

func testAccCheckDigitalOceanCertificateConfig_customNoLeaf(rInt int, privateKeyMaterial, certChain string) string {
	return fmt.Sprintf(`
resource "digitalocean_certificate" "foobar" {
  name = "certificate-%d"
  private_key = <<EOF
%s
EOF
  certificate_chain = <<EOF
%s
EOF
}`, rInt, privateKeyMaterial, certChain)
}

func testAccCheckDigitalOceanCertificateConfig_customNoKey(rInt int, leafCert, certChain string) string {
	return fmt.Sprintf(`
resource "digitalocean_certificate" "foobar" {
  name = "certificate-%d"
  leaf_certificate = <<EOF
%s
EOF
  certificate_chain = <<EOF
%s
EOF
}`, rInt, leafCert, certChain)
}

func testAccCheckDigitalOceanCertificateConfig_noDomains(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_certificate" "foobar" {
  name = "certificate-%d"
  type = "lets_encrypt"
}`, rInt)
}
