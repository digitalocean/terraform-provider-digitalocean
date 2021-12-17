package digitalocean

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testNamePrefix = "tf-acc-test-"

var testAccProvider *schema.Provider
var testAccProviders map[string]*schema.Provider
var testAccProviderFactories map[string]func() (*schema.Provider, error)

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"digitalocean": testAccProvider,
	}
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"digitalocean": func() (*schema.Provider, error) {
			return testAccProvider, nil
		},
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("DIGITALOCEAN_TOKEN"); v == "" {
		t.Fatal("DIGITALOCEAN_TOKEN must be set for acceptance tests")
	}

	err := testAccProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(nil))
	if err != nil {
		t.Fatal(err)
	}
}

func TestURLOverride(t *testing.T) {
	customEndpoint := "https://mock-api.internal.example.com/"

	rawProvider := Provider()
	raw := map[string]interface{}{
		"token":        "12345",
		"api_endpoint": customEndpoint,
	}

	diags := rawProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("provider configure failed: %s", diagnosticsToString(diags))
	}

	meta := rawProvider.Meta()
	if meta == nil {
		t.Fatalf("Expected metadata, got nil")
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

	diags := rawProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("provider configure failed: %s", diagnosticsToString(diags))
	}

	meta := rawProvider.Meta()
	if meta == nil {
		t.Fatal("Expected metadata, got nil")
	}

	client := meta.(*CombinedConfig).godoClient()
	if client.BaseURL.String() != "https://api.digitalocean.com" {
		t.Fatalf("Expected %s, got %s", "https://api.digitalocean.com", client.BaseURL.String())
	}
}

func diagnosticsToString(diags diag.Diagnostics) string {
	diagsAsStrings := make([]string, len(diags))
	for i, diag := range diags {
		diagsAsStrings[i] = diag.Summary
	}

	return strings.Join(diagsAsStrings, "; ")
}

func TestSpaceAPIDefaultEndpoint(t *testing.T) {
	rawProvider := Provider()
	raw := map[string]interface{}{
		"token":             "12345",
		"spaces_access_id":  "abcdef",
		"spaces_secret_key": "xyzzy",
	}

	diags := rawProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("provider configure failed: %s", diagnosticsToString(diags))
	}

	meta := rawProvider.Meta()
	if meta == nil {
		t.Fatalf("Expected metadata, got nil")
	}

	var client *session.Session
	var err error
	client, err = meta.(*CombinedConfig).spacesClient("sfo2")
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

	diags := rawProvider.Configure(context.Background(), terraform.NewResourceConfigRaw(raw))
	if diags.HasError() {
		t.Fatalf("provider configure failed: %s", diagnosticsToString(diags))
	}

	meta := rawProvider.Meta()
	if meta == nil {
		t.Fatal("Expected metadata, got nil")
	}

	var client *session.Session
	var err error
	client, err = meta.(*CombinedConfig).spacesClient("sfo2")
	if err != nil {
		t.Fatalf("Failed to create Spaces client: %s", err)
	}

	expectedEndpoint := "https://sfo2.not-digitalocean-domain.com"
	if *client.Config.Endpoint != expectedEndpoint {
		t.Fatalf("Expected %s, got %s", expectedEndpoint, *client.Config.Endpoint)
	}
}

func randomTestName(additionalNames ...string) string {
	prefix := testNamePrefix
	for _, n := range additionalNames {
		prefix += "-" + strings.Replace(n, " ", "_", -1)
	}
	return randomName(prefix, 10)
}

func randomName(prefix string, length int) string {
	return fmt.Sprintf("%s%s", prefix, acctest.RandString(length))
}
