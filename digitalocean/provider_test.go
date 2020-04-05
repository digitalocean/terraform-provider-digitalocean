package digitalocean

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
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

func TestSpaceAPIDefaultEndpoint(t *testing.T) {
	rawProvider := Provider()
	raw := map[string]interface{}{
		"token":             "12345",
		"spaces_access_id":  "abcdef",
		"spaces_secret_key": "xyzzy",
	}

	err := rawProvider.Configure(terraform.NewResourceConfigRaw(raw))
	meta := rawProvider.(*schema.Provider).Meta()
	if meta == nil {
		t.Fatalf("Expected metadata, got nil: err: %s", err)
	}

	client, err := meta.(*CombinedConfig).spacesClient("sfo2")
	if err != nil {
		t.Fatalf("Failed to create Spaces client: %s", err)
	}

	expectedEndpoint := "https://sfo2.digitaloceanspaces.com"
	if *client.Config.Endpoint != expectedEndpoint {
		t.Fatalf("Expected %s, got %s", expectedEndpoint, *client.Config.Endpoint)
	}
}

func TestSpaceAPIEndpointOverride(t *testing.T) {
	customSpacesEndpoint := "https://{{.Region}}.not-digitalocean-domain.com"

	rawProvider := Provider()
	raw := map[string]interface{}{
		"token":             "12345",
		"spaces_endpoint":   customSpacesEndpoint,
		"spaces_access_id":  "abcdef",
		"spaces_secret_key": "xyzzy",
	}

	err := rawProvider.Configure(terraform.NewResourceConfigRaw(raw))
	meta := rawProvider.(*schema.Provider).Meta()
	if meta == nil {
		t.Fatalf("Expected metadata, got nil: err: %s", err)
	}

	client, err := meta.(*CombinedConfig).spacesClient("sfo2")
	if err != nil {
		t.Fatalf("Failed to create Spaces client: %s", err)
	}

	expectedEndpoint := "https://sfo2.not-digitalocean-domain.com"
	if *client.Config.Endpoint != expectedEndpoint {
		t.Fatalf("Expected %s, got %s", expectedEndpoint, *client.Config.Endpoint)
	}
}

func randomTestName() string {
	return randomName(testNamePrefix, 10)
}

func randomName(prefix string, length int) string {
	return fmt.Sprintf("%s%s", prefix, acctest.RandString(length))
}
