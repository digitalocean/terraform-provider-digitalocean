package acceptance

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

// GenerateTestCertMaterial generates a TLS certificate for use in acceptance test fixtures.
func GenerateTestCertMaterial(t *testing.T) (string, string, string) {
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

	cert, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
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
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
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
