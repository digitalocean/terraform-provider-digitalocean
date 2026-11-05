package dedicatedinference_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanDedicatedInferenceTokens_Basic(t *testing.T) {
	diName := acceptance.RandomTestName() + "-di"
	tokenName := acceptance.RandomTestName() + "-token"
	vpcUUID := testDIVPCUUID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDedicatedInferenceTokensConfig(diName, tokenName, vpcUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_dedicated_inference_tokens.test", "tokens.#"),
					resource.TestCheckResourceAttrSet("data.digitalocean_dedicated_inference_tokens.test", "tokens.0.id"),
					resource.TestCheckResourceAttrSet("data.digitalocean_dedicated_inference_tokens.test", "tokens.0.name"),
					resource.TestCheckResourceAttrSet("data.digitalocean_dedicated_inference_tokens.test", "tokens.0.created_at"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanDedicatedInferenceTokens_FilterByName(t *testing.T) {
	diName := acceptance.RandomTestName() + "-di"
	tokenName := acceptance.RandomTestName() + "-token"
	vpcUUID := testDIVPCUUID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDedicatedInferenceTokensFilterConfig(diName, tokenName, vpcUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_dedicated_inference_tokens.by_name", "tokens.#", "1"),
					resource.TestCheckResourceAttr("data.digitalocean_dedicated_inference_tokens.by_name", "tokens.0.name", tokenName),
				),
			},
		},
	})
}

func testAccDataSourceDedicatedInferenceTokensConfig(diName, tokenName, vpcUUID string) string {
	return fmt.Sprintf(`
resource "digitalocean_dedicated_inference" "test" {
  name                   = "%s"
  region                 = "%s"
  vpc_uuid               = "%s"
  enable_public_endpoint = true

  model_deployments {
    model_slug        = "%s"
    model_provider    = "%s"
    provider_model_id = "%s"

    accelerators {
      accelerator_slug = "%s"
      scale            = 1
      type             = "%s"
    }
  }
}

resource "digitalocean_dedicated_inference_token" "test" {
  dedicated_inference_id = digitalocean_dedicated_inference.test.id
  name                   = "%s"
}

data "digitalocean_dedicated_inference_tokens" "test" {
  dedicated_inference_id = digitalocean_dedicated_inference.test.id
  depends_on             = [digitalocean_dedicated_inference_token.test]
}
`, diName, testDIRegion, vpcUUID, testDIModelSlug, testDIModelProvider, testDIProviderModelID, testDIAcceleratorSlug, testDIAcceleratorType, tokenName)
}

func testAccDataSourceDedicatedInferenceTokensFilterConfig(diName, tokenName, vpcUUID string) string {
	return fmt.Sprintf(`
resource "digitalocean_dedicated_inference" "test" {
  name                   = "%s"
  region                 = "%s"
  vpc_uuid               = "%s"
  enable_public_endpoint = true

  model_deployments {
    model_slug        = "%s"
    model_provider    = "%s"
    provider_model_id = "%s"

    accelerators {
      accelerator_slug = "%s"
      scale            = 1
      type             = "%s"
    }
  }
}

resource "digitalocean_dedicated_inference_token" "test" {
  dedicated_inference_id = digitalocean_dedicated_inference.test.id
  name                   = "%s"
}

data "digitalocean_dedicated_inference_tokens" "by_name" {
  dedicated_inference_id = digitalocean_dedicated_inference.test.id

  filter {
    key    = "name"
    values = ["%s"]
  }

  depends_on = [digitalocean_dedicated_inference_token.test]
}
`, diName, testDIRegion, vpcUUID, testDIModelSlug, testDIModelProvider, testDIProviderModelID, testDIAcceleratorSlug, testDIAcceleratorType, tokenName, tokenName)
}
