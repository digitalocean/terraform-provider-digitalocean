package dedicatedinference_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanDedicatedInferenceAccelerators_Basic(t *testing.T) {
	name := acceptance.RandomTestName() + "-di"
	vpcUUID := testDIVPCUUID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDedicatedInferenceAcceleratorsConfig(name, vpcUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_dedicated_inference_accelerators.test", "accelerators.#"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanDedicatedInferenceAccelerators_FilterBySlug(t *testing.T) {
	name := acceptance.RandomTestName() + "-di"
	vpcUUID := testDIVPCUUID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDedicatedInferenceAcceleratorsFilterConfig(name, vpcUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_dedicated_inference_accelerators.by_slug", "accelerators.#"),
					resource.TestCheckResourceAttr("data.digitalocean_dedicated_inference_accelerators.by_slug", "accelerators.0.slug", testDIAcceleratorSlug),
					resource.TestCheckResourceAttrSet("data.digitalocean_dedicated_inference_accelerators.by_slug", "accelerators.0.id"),
					resource.TestCheckResourceAttrSet("data.digitalocean_dedicated_inference_accelerators.by_slug", "accelerators.0.status"),
				),
			},
		},
	})
}

func testAccDataSourceDedicatedInferenceAcceleratorsConfig(name, vpcUUID string) string {
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

data "digitalocean_dedicated_inference_accelerators" "test" {
  dedicated_inference_id = digitalocean_dedicated_inference.test.id
}
`, name, testDIRegion, vpcUUID, testDIModelSlug, testDIModelProvider, testDIProviderModelID, testDIAcceleratorSlug, testDIAcceleratorType)
}

func testAccDataSourceDedicatedInferenceAcceleratorsFilterConfig(name, vpcUUID string) string {
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

data "digitalocean_dedicated_inference_accelerators" "by_slug" {
  dedicated_inference_id = digitalocean_dedicated_inference.test.id

  filter {
    key    = "slug"
    values = ["%s"]
  }
}
`, name, testDIRegion, vpcUUID, testDIModelSlug, testDIModelProvider, testDIProviderModelID, testDIAcceleratorSlug, testDIAcceleratorType, testDIAcceleratorSlug)
}
