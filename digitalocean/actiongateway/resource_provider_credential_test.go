package actiongateway_test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

// apiKeyProviderForTest returns a tool provider that accepts a team API key
// and a key for it. The catalog does not say which provider takes a key
// without further parameters, and the suite has no key of its own, so both
// come from the environment.
func apiKeyProviderForTest(t *testing.T) (provider, apiKey string) {
	t.Helper()

	provider = os.Getenv("DO_ACTION_GATEWAY_API_KEY_PROVIDER")
	apiKey = os.Getenv("DO_ACTION_GATEWAY_API_KEY")
	if provider == "" || apiKey == "" {
		t.Skip("DO_ACTION_GATEWAY_API_KEY_PROVIDER and DO_ACTION_GATEWAY_API_KEY must be set to a tool provider that accepts an API key with no credential parameters, and a key for it")
	}
	return provider, apiKey
}

func TestAccDigitalOceanActionGatewayProviderCredential_APIKey(t *testing.T) {
	provider, apiKey := apiKeyProviderForTest(t)
	name := acceptance.RandomTestName()
	resourceName := "digitalocean_action_gateway_provider_credential.foobar"
	var id string

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckActionGatewayProviderCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccActionGatewayProviderCredentialConfig, provider, name, apiKey),
				Check: resource.ComposeTestCheckFunc(
					testAccCaptureResourceID(resourceName, &id),
					resource.TestCheckResourceAttr(resourceName, "tool_provider", provider),
					resource.TestCheckResourceAttr(resourceName, "kind", "api_key"),
					resource.TestCheckResourceAttr(resourceName, "display_name", name),
					resource.TestCheckResourceAttr(resourceName, "status", "active"),
					resource.TestCheckResourceAttr(resourceName, "reusable_by_team", "false"),
					resource.TestCheckResourceAttr(resourceName, "scope", "user"),
					resource.TestCheckResourceAttrSet(resourceName, "owning_user_id"),
					resource.TestCheckResourceAttrSet(resourceName, "created_at"),
				),
			},
			{
				// A new display name changes the credential in place.
				Config: fmt.Sprintf(testAccActionGatewayProviderCredentialConfig, provider, name+"-renamed", apiKey),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr(resourceName, "id", &id),
					resource.TestCheckResourceAttr(resourceName, "display_name", name+"-renamed"),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
			},
		},
	})
}

func TestAccDataSourceDigitalOceanActionGatewayProviderCredential_ByDisplayName(t *testing.T) {
	provider, apiKey := apiKeyProviderForTest(t)
	name := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccActionGatewayProviderCredentialConfig, provider, name, apiKey)
	resourceName := "digitalocean_action_gateway_provider_credential.foobar"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckActionGatewayProviderCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + `
data "digitalocean_action_gateway_provider_credential" "by_name" {
  tool_provider = digitalocean_action_gateway_provider_credential.foobar.tool_provider
  display_name  = digitalocean_action_gateway_provider_credential.foobar.display_name
}

data "digitalocean_action_gateway_provider_credential" "by_id" {
  id = digitalocean_action_gateway_provider_credential.foobar.id
}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("data.digitalocean_action_gateway_provider_credential.by_name", "id", resourceName, "id"),
					resource.TestCheckResourceAttrPair("data.digitalocean_action_gateway_provider_credential.by_name", "kind", resourceName, "kind"),
					resource.TestCheckResourceAttrPair("data.digitalocean_action_gateway_provider_credential.by_id", "display_name", resourceName, "display_name"),
					resource.TestCheckResourceAttrPair("data.digitalocean_action_gateway_provider_credential.by_id", "status", resourceName, "status"),
				),
			},
		},
	})
}

func testAccCaptureResourceID(resourceName string, id *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource not found: %s", resourceName)
		}
		*id = rs.Primary.ID
		return nil
	}
}

func testAccCheckActionGatewayProviderCredentialDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_action_gateway_provider_credential" {
			continue
		}
		_, resp, err := client.ActionGateway.ProviderCredentials.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("Action Gateway provider credential %s still exists", rs.Primary.ID)
		}
		if resp == nil || resp.StatusCode != http.StatusNotFound {
			return err
		}
	}
	return nil
}

const testAccActionGatewayProviderCredentialConfig = `
resource "digitalocean_action_gateway_provider_credential" "foobar" {
  tool_provider = "%s"
  kind          = "api_key"
  display_name  = "%s"
  api_key       = "%s"
}
`
