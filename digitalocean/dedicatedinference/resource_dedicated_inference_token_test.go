package dedicatedinference_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanDedicatedInferenceToken_Basic(t *testing.T) {
	diName := acceptance.RandomTestName() + "-di"
	tokenName := acceptance.RandomTestName() + "-token"
	vpcUUID := testDIVPCUUID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDedicatedInferenceTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedInferenceTokenConfig(diName, tokenName, vpcUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_dedicated_inference_token.test", "name", tokenName),
					resource.TestCheckResourceAttrSet("digitalocean_dedicated_inference_token.test", "dedicated_inference_id"),
					resource.TestCheckResourceAttrSet("digitalocean_dedicated_inference_token.test", "token"),
					resource.TestCheckResourceAttrSet("digitalocean_dedicated_inference_token.test", "created_at"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDedicatedInferenceToken_Import(t *testing.T) {
	diName := acceptance.RandomTestName() + "-di"
	tokenName := acceptance.RandomTestName() + "-token"
	vpcUUID := testDIVPCUUID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDedicatedInferenceTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedInferenceTokenConfig(diName, tokenName, vpcUUID),
			},
			{
				ResourceName:            "digitalocean_dedicated_inference_token.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token"},
			},
		},
	})
}

func testAccCheckDigitalOceanDedicatedInferenceTokenDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_dedicated_inference_token" {
			continue
		}

		diID := rs.Primary.Attributes["dedicated_inference_id"]
		tokenID := rs.Primary.ID

		tokens, _, err := client.DedicatedInference.ListTokens(context.Background(), diID, nil)
		if err != nil {
			continue
		}
		for _, t := range tokens {
			if t.ID == tokenID {
				return fmt.Errorf("dedicated inference token still exists: %s", tokenID)
			}
		}
	}

	return nil
}

func testAccDedicatedInferenceTokenConfig(diName, tokenName, vpcUUID string) string {
	return fmt.Sprintf(`
resource "digitalocean_dedicated_inference" "test" {
  name                   = "%s"
  region                 = "%s"
  vpc_uuid               = "%s"
  enable_public_endpoint = true

  model_deployments {
    model_slug     = "%s"
    model_provider = "%s"

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
`, diName, testDIRegion, vpcUUID, testDIModelSlug, testDIModelProvider, testDIAcceleratorSlug, testDIAcceleratorType, tokenName)
}
