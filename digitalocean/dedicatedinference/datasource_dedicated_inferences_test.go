package dedicatedinference_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// TestAccDataSourceDigitalOceanDedicatedInferences_List verifies the datasource can call
// the List API and return results without requiring a resource to be created first.
func TestAccDataSourceDigitalOceanDedicatedInferences_List(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "digitalocean_dedicated_inferences" "all" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_dedicated_inferences.all", "dedicated_inferences.#"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanDedicatedInferences_Basic(t *testing.T) {
	name := acceptance.RandomTestName() + "-di"
	vpcUUID := testDIVPCUUID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDedicatedInferencesConfig(name, vpcUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_dedicated_inferences.all", "dedicated_inferences.#"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanDedicatedInferences_FilterByName(t *testing.T) {
	name := acceptance.RandomTestName() + "-di"
	vpcUUID := testDIVPCUUID(t)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDedicatedInferencesFilterConfig(name, vpcUUID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_dedicated_inferences.by_name", "dedicated_inferences.#", "1"),
					resource.TestCheckResourceAttr("data.digitalocean_dedicated_inferences.by_name", "dedicated_inferences.0.name", name),
					resource.TestCheckResourceAttr("data.digitalocean_dedicated_inferences.by_name", "dedicated_inferences.0.region", testDIRegion),
					resource.TestCheckResourceAttrSet("data.digitalocean_dedicated_inferences.by_name", "dedicated_inferences.0.id"),
					resource.TestCheckResourceAttrSet("data.digitalocean_dedicated_inferences.by_name", "dedicated_inferences.0.status"),
				),
			},
		},
	})
}

func testAccDataSourceDedicatedInferencesConfig(name, vpcUUID string) string {
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

data "digitalocean_dedicated_inferences" "all" {
  depends_on = [digitalocean_dedicated_inference.test]
}
`, name, testDIRegion, vpcUUID, testDIModelSlug, testDIModelProvider, testDIAcceleratorSlug, testDIAcceleratorType)
}

func testAccDataSourceDedicatedInferencesFilterConfig(name, vpcUUID string) string {
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

data "digitalocean_dedicated_inferences" "by_name" {
  filter {
    key    = "name"
    values = ["%s"]
  }

  depends_on = [digitalocean_dedicated_inference.test]
}
`, name, testDIRegion, vpcUUID, testDIModelSlug, testDIModelProvider, testDIAcceleratorSlug, testDIAcceleratorType, name)
}
