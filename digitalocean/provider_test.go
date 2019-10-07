package digitalocean

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/terraform-providers/terraform-provider-kubernetes/kubernetes"
)

const testNamePrefix = "tf-acc-test-"

var (
	testAccProviders map[string]terraform.ResourceProvider
	testAccProvider  *schema.Provider
)

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"digitalocean": testAccProvider,
		"kubernetes":   kubernetes.Provider(),
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("DIGITALOCEAN_TOKEN"); v == "" {
		t.Fatal("DIGITALOCEAN_TOKEN must be set for acceptance tests")
	}
}

func TestURLOverride(t *testing.T) {
	customEndpoint := "https://mock-api.internal.example.com/"

	rawProvider := Provider()
	raw := map[string]interface{}{
		"token":        "12345",
		"api_endpoint": customEndpoint,
	}

	err := rawProvider.Configure(terraform.NewResourceConfigRaw(raw))
	meta := rawProvider.(*schema.Provider).Meta()
	if meta == nil {
		t.Fatalf("Expected metadata, got nil: err: %s", err)
	}
	client := meta.(*CombinedConfig).godoClient()
	if client.BaseURL.String() != customEndpoint {
		t.Fatalf("Expected %s, got %s", customEndpoint, client.BaseURL.String())
	}
}

func TestURLDefault(t *testing.T) {
	rawProvider := Provider()
	raw := map[string]interface{}{
		"token": "12345",
	}

	err := rawProvider.Configure(terraform.NewResourceConfigRaw(raw))
	meta := rawProvider.(*schema.Provider).Meta()
	if meta == nil {
		t.Fatalf("Expected metadata, got nil: err: %s", err)
	}
	client := meta.(*CombinedConfig).godoClient()
	if client.BaseURL.String() != "https://api.digitalocean.com" {
		t.Fatalf("Expected %s, got %s", "https://api.digitalocean.com", client.BaseURL.String())
	}
}

func randomTestName() string {
	return randomName(testNamePrefix, 10)
}

func randomName(prefix string, length int) string {
	return fmt.Sprintf("%s%s", prefix, acctest.RandString(length))
}
