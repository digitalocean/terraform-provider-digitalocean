package agentplatform_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// openAITestKey returns a real OpenAI API key for acceptance tests. The
// DigitalOcean API validates the key, so a placeholder does not work; tests
// are skipped instead when the environment variable is unset.
func openAITestKey(t *testing.T, envVar string) string {
	t.Helper()

	key := os.Getenv(envVar)
	if key == "" {
		t.Skipf("%s must be set to a real OpenAI API key", envVar)
	}

	return key
}

func TestAccDigitalOceanOpenAIApiKey_Basic(t *testing.T) {
	apiKey := openAITestKey(t, "OPENAI_TEST_KEY")
	keyName := acceptance.RandomTestName() + "-openai-key"
	resourceConfig := fmt.Sprintf(`
resource "digitalocean_agent_platform_openai_api_key" "test" {
  api_key = "%s"
  name    = "%s"
}
`, apiKey, keyName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_agent_platform_openai_api_key.test", "name", keyName),
					resource.TestCheckResourceAttr("digitalocean_agent_platform_openai_api_key.test", "api_key", apiKey),
				),
			},
		},
	})
}

func TestAccDigitalOceanOpenAIApiKey_Update(t *testing.T) {
	apiKey := openAITestKey(t, "OPENAI_TEST_KEY")
	rotatedKey := os.Getenv("OPENAI_TEST_KEY_ROTATED")
	if rotatedKey == "" {
		rotatedKey = apiKey
	}

	keyName := acceptance.RandomTestName() + "-openai-key"
	updatedKeyName := keyName + "-updated"

	resourceConfig := fmt.Sprintf(`
resource "digitalocean_agent_platform_openai_api_key" "test" {
  api_key = "%s"
  name    = "%s"
}
`, apiKey, keyName)

	updatedResourceConfig := fmt.Sprintf(`
resource "digitalocean_agent_platform_openai_api_key" "test" {
  api_key = "%s"
  name    = "%s"
}
`, rotatedKey, updatedKeyName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_agent_platform_openai_api_key.test", "name", keyName),
					resource.TestCheckResourceAttr("digitalocean_agent_platform_openai_api_key.test", "api_key", apiKey),
				),
			},
			{
				Config: updatedResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_agent_platform_openai_api_key.test", "name", updatedKeyName),
					resource.TestCheckResourceAttr("digitalocean_agent_platform_openai_api_key.test", "api_key", rotatedKey),
				),
			},
		},
	})
}

// TestAccDigitalOceanOpenAIApiKey_Validation covers the argument combinations
// that must be rejected. It needs a DigitalOcean token but no OpenAI key,
// because validation fails before any API call.
func TestAccDigitalOceanOpenAIApiKey_Validation(t *testing.T) {
	keyName := acceptance.RandomTestName() + "-openai-key"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
resource "digitalocean_agent_platform_openai_api_key" "test" {
  api_key            = "first"
  api_key_wo         = "second"
  api_key_wo_version = 1
  name               = "%s"
}
`, keyName),
				ExpectError: regexp.MustCompile(`(?s)only\s+one\s+of\s+.*can\s+be\s+specified`),
			},
			{
				Config: fmt.Sprintf(`
resource "digitalocean_agent_platform_openai_api_key" "test" {
  name = "%s"
}
`, keyName),
				ExpectError: regexp.MustCompile(`(?s)one\s+of\s+.*must\s+be\s+specified`),
			},
			{
				Config: fmt.Sprintf(`
resource "digitalocean_agent_platform_openai_api_key" "test" {
  api_key_wo_version = 1
  name               = "%s"
}
`, keyName),
				ExpectError: regexp.MustCompile(`(?s)(one|all)\s+of\s+.*must\s+be\s+specified`),
			},
			{
				Config: fmt.Sprintf(`
resource "digitalocean_agent_platform_openai_api_key" "test" {
  api_key_wo = "only-the-key"
  name       = "%s"
}
`, keyName),
				ExpectError: regexp.MustCompile(`(?s)all\s+of\s+.*must\s+be\s+specified`),
			},
			{
				Config: fmt.Sprintf(`
resource "digitalocean_agent_platform_openai_api_key" "test" {
  api_key            = "legacy-key"
  api_key_wo_version = 1
  name               = "%s"
}
`, keyName),
				ExpectError: regexp.MustCompile(`(?s)all\s+of\s+.*must\s+be\s+specified`),
			},
		},
	})
}

func TestAccDigitalOceanOpenAIApiKey_WriteOnly(t *testing.T) {
	apiKey := openAITestKey(t, "OPENAI_TEST_KEY")

	// A second, different real key proves the update sends the new value
	// instead of re-sending the one already used at create time.
	rotatedKey := os.Getenv("OPENAI_TEST_KEY_ROTATED")
	if rotatedKey == "" {
		t.Log("OPENAI_TEST_KEY_ROTATED is unset; the version bump re-sends the same key and only proves that an update is triggered")
		rotatedKey = apiKey
	}

	keyName := acceptance.RandomTestName() + "-openai-key"
	resourceConfig := `
resource "digitalocean_agent_platform_openai_api_key" "test" {
  api_key_wo         = "%s"
  api_key_wo_version = %d
  name               = "%s"
}
`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(resourceConfig, apiKey, 1, keyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_agent_platform_openai_api_key.test", "name", keyName),
					resource.TestCheckResourceAttr("digitalocean_agent_platform_openai_api_key.test", "api_key_wo_version", "1"),
					resource.TestCheckNoResourceAttr("digitalocean_agent_platform_openai_api_key.test", "api_key_wo"),
					resource.TestCheckNoResourceAttr("digitalocean_agent_platform_openai_api_key.test", "api_key"),
				),
			},
			{
				Config: fmt.Sprintf(resourceConfig, rotatedKey, 2, keyName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_agent_platform_openai_api_key.test", "api_key_wo_version", "2"),
					resource.TestCheckNoResourceAttr("digitalocean_agent_platform_openai_api_key.test", "api_key_wo"),
					resource.TestCheckNoResourceAttr("digitalocean_agent_platform_openai_api_key.test", "api_key"),
				),
			},
			{
				ResourceName:            "digitalocean_agent_platform_openai_api_key.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key_wo", "api_key_wo_version"},
			},
		},
	})
}

func TestAccDigitalOceanOpenAIApiKey_Delete(t *testing.T) {
	apiKey := openAITestKey(t, "OPENAI_TEST_KEY")
	keyName := acceptance.RandomTestName() + "-openai-key"
	resourceConfig := fmt.Sprintf(`
resource "digitalocean_agent_platform_openai_api_key" "test" {
  api_key = "%s"
  name    = "%s"
}
`, apiKey, keyName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_agent_platform_openai_api_key.test", "name", keyName),
				),
			},
			{
				ResourceName:            "digitalocean_agent_platform_openai_api_key.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"api_key"},
			},
		},
	})
}
