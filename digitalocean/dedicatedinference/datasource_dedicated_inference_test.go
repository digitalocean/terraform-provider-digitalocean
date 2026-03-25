package dedicatedinference_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanDedicatedInference_Basic(t *testing.T) {
	name := acceptance.RandomTestName() + "-di"
	vpcUUID := testDIVPCUUID(t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDedicatedInferenceConfig(name, vpcUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_dedicated_inference.test", "name", name),
					resource.TestCheckResourceAttr("digitalocean_dedicated_inference.test", "region", testDIRegion),
					resource.TestCheckResourceAttrPair(
						"data.digitalocean_dedicated_inference.test", "id",
						"digitalocean_dedicated_inference.test", "id",
					),
					resource.TestCheckResourceAttrPair(
						"data.digitalocean_dedicated_inference.test", "name",
						"digitalocean_dedicated_inference.test", "name",
					),
					resource.TestCheckResourceAttrPair(
						"data.digitalocean_dedicated_inference.test", "region",
						"digitalocean_dedicated_inference.test", "region",
					),
					resource.TestCheckResourceAttrPair(
						"data.digitalocean_dedicated_inference.test", "status",
						"digitalocean_dedicated_inference.test", "status",
					),
					resource.TestCheckResourceAttr("data.digitalocean_dedicated_inference.test", "model_deployments.#", "1"),
					resource.TestCheckResourceAttr("data.digitalocean_dedicated_inference.test", "model_deployments.0.model_slug", testDIModelSlug),
					resource.TestCheckResourceAttr("data.digitalocean_dedicated_inference.test", "model_deployments.0.model_provider", testDIModelProvider),
					resource.TestCheckResourceAttr("data.digitalocean_dedicated_inference.test", "model_deployments.0.accelerators.#", "1"),
					resource.TestCheckResourceAttr("data.digitalocean_dedicated_inference.test", "model_deployments.0.accelerators.0.accelerator_slug", testDIAcceleratorSlug),
					resource.TestCheckResourceAttrSet("data.digitalocean_dedicated_inference.test", "created_at"),
					resource.TestCheckResourceAttrSet("data.digitalocean_dedicated_inference.test", "updated_at"),
				),
			},
		},
	})
}

func testAccDataSourceDedicatedInferenceConfig(name, vpcUUID string) string {
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

data "digitalocean_dedicated_inference" "test" {
  id = digitalocean_dedicated_inference.test.id
}
`, name, testDIRegion, vpcUUID, testDIModelSlug, testDIModelProvider, testDIAcceleratorSlug, testDIAcceleratorType)
}
