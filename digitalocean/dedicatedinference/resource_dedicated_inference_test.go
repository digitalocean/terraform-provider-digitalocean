package dedicatedinference_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const (
	testDIRegion          = "atl1"
	testDIModelSlug       = "Qwen/Qwen2.5-14B-Instruct"
	testDIModelProvider   = "hugging_face"
	testDIAcceleratorSlug = "gpu-mi300x1-192gb"
	testDIAcceleratorType = "prefill_decode"
)

func testDIVPCUUID(t *testing.T) string {
	t.Helper()
	v := os.Getenv("DO_DEDICATED_INFERENCE_VPC_UUID")
	if v == "" {
		t.Skip("DO_DEDICATED_INFERENCE_VPC_UUID must be set for dedicated inference acceptance tests")
	}
	return v
}

func TestAccDigitalOceanDedicatedInference_Basic(t *testing.T) {
	var di godo.DedicatedInference
	name := acceptance.RandomTestName() + "-di"
	vpcUUID := testDIVPCUUID(t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDedicatedInferenceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedInferenceConfig_basic(name, vpcUUID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDedicatedInferenceExists("digitalocean_dedicated_inference.test", &di),
					resource.TestCheckResourceAttr("digitalocean_dedicated_inference.test", "name", name),
					resource.TestCheckResourceAttr("digitalocean_dedicated_inference.test", "region", testDIRegion),
					resource.TestCheckResourceAttr("digitalocean_dedicated_inference.test", "enable_public_endpoint", "true"),
					resource.TestCheckResourceAttr("digitalocean_dedicated_inference.test", "model_deployments.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_dedicated_inference.test", "model_deployments.0.model_slug", testDIModelSlug),
					resource.TestCheckResourceAttr("digitalocean_dedicated_inference.test", "model_deployments.0.model_provider", testDIModelProvider),
					resource.TestCheckResourceAttr("digitalocean_dedicated_inference.test", "model_deployments.0.accelerators.#", "1"),
					resource.TestCheckResourceAttr("digitalocean_dedicated_inference.test", "model_deployments.0.accelerators.0.accelerator_slug", testDIAcceleratorSlug),
					resource.TestCheckResourceAttr("digitalocean_dedicated_inference.test", "model_deployments.0.accelerators.0.scale", "1"),
					resource.TestCheckResourceAttr("digitalocean_dedicated_inference.test", "model_deployments.0.accelerators.0.type", testDIAcceleratorType),
					resource.TestCheckResourceAttrSet("digitalocean_dedicated_inference.test", "status"),
					resource.TestCheckResourceAttrSet("digitalocean_dedicated_inference.test", "created_at"),
					resource.TestCheckResourceAttrSet("digitalocean_dedicated_inference.test", "updated_at"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDedicatedInference_Update(t *testing.T) {
	var di godo.DedicatedInference
	name := acceptance.RandomTestName() + "-di"
	updatedName := name + "-updated"
	vpcUUID := testDIVPCUUID(t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDedicatedInferenceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedInferenceConfig_basic(name, vpcUUID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDedicatedInferenceExists("digitalocean_dedicated_inference.test", &di),
					resource.TestCheckResourceAttr("digitalocean_dedicated_inference.test", "name", name),
					resource.TestCheckResourceAttr("digitalocean_dedicated_inference.test", "enable_public_endpoint", "true"),
				),
			},
			{
				Config: testAccDedicatedInferenceConfig_updated(updatedName, vpcUUID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDedicatedInferenceExists("digitalocean_dedicated_inference.test", &di),
					resource.TestCheckResourceAttr("digitalocean_dedicated_inference.test", "name", updatedName),
					resource.TestCheckResourceAttr("digitalocean_dedicated_inference.test", "enable_public_endpoint", "true"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDedicatedInference_UpdateModelDeployments(t *testing.T) {
	var di godo.DedicatedInference
	name := acceptance.RandomTestName() + "-di"
	vpcUUID := testDIVPCUUID(t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDedicatedInferenceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedInferenceConfig_basic(name, vpcUUID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDedicatedInferenceExists("digitalocean_dedicated_inference.test", &di),
					resource.TestCheckResourceAttr("digitalocean_dedicated_inference.test", "model_deployments.0.accelerators.0.scale", "1"),
				),
			},
			{
				Config: testAccDedicatedInferenceConfig_scaledAccelerator(name, vpcUUID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDedicatedInferenceExists("digitalocean_dedicated_inference.test", &di),
					resource.TestCheckResourceAttr("digitalocean_dedicated_inference.test", "model_deployments.0.accelerators.0.scale", "2"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDedicatedInference_Delete(t *testing.T) {
	name := acceptance.RandomTestName() + "-di"
	vpcUUID := testDIVPCUUID(t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDedicatedInferenceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedInferenceConfig_basic(name, vpcUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_dedicated_inference.test", "name", name),
				),
			},
		},
	})
}

func TestAccDigitalOceanDedicatedInference_Import(t *testing.T) {
	name := acceptance.RandomTestName() + "-di"
	vpcUUID := testDIVPCUUID(t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDedicatedInferenceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedInferenceConfig_basic(name, vpcUUID),
			},
			{
				ResourceName:            "digitalocean_dedicated_inference.test",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"hugging_face_token"},
			},
		},
	})
}

func testAccCheckDigitalOceanDedicatedInferenceDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_dedicated_inference" {
			continue
		}

		_, _, err := client.DedicatedInference.Get(context.Background(), rs.Primary.ID)
		if err == nil {
			return fmt.Errorf("dedicated inference endpoint still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckDigitalOceanDedicatedInferenceExists(resourceName string, di *godo.DedicatedInference) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("not found: %s", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no dedicated inference ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		found, _, err := client.DedicatedInference.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("dedicated inference endpoint not found")
		}

		*di = *found
		return nil
	}
}

func testAccDedicatedInferenceConfig_basic(name, vpcUUID string) string {
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
`, name, testDIRegion, vpcUUID, testDIModelSlug, testDIModelProvider, testDIAcceleratorSlug, testDIAcceleratorType)
}

func testAccDedicatedInferenceConfig_updated(name, vpcUUID string) string {
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
`, name, testDIRegion, vpcUUID, testDIModelSlug, testDIModelProvider, testDIAcceleratorSlug, testDIAcceleratorType)
}

func testAccDedicatedInferenceConfig_scaledAccelerator(name, vpcUUID string) string {
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
      scale            = 2
      type             = "%s"
    }
  }
}
`, name, testDIRegion, vpcUUID, testDIModelSlug, testDIModelProvider, testDIAcceleratorSlug, testDIAcceleratorType)
}
