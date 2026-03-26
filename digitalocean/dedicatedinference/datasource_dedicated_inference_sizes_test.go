package dedicatedinference_test

import (
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanDedicatedInferenceSizes_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDedicatedInferenceSizesConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_dedicated_inference_sizes.test", "sizes.#"),
					resource.TestCheckResourceAttrSet("data.digitalocean_dedicated_inference_sizes.test", "enabled_regions.#"),
					resource.TestCheckResourceAttrSet("data.digitalocean_dedicated_inference_sizes.test", "sizes.0.gpu_slug"),
					resource.TestCheckResourceAttrSet("data.digitalocean_dedicated_inference_sizes.test", "sizes.0.price_per_hour"),
					resource.TestCheckResourceAttrSet("data.digitalocean_dedicated_inference_sizes.test", "sizes.0.currency"),
				),
			},
		},
	})
}

func testAccDataSourceDedicatedInferenceSizesConfig() string {
	return `
data "digitalocean_dedicated_inference_sizes" "test" {}
`
}
