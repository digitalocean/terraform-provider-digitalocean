package dedicatedinference_test

import (
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanDedicatedInferenceGPUModelConfig_Basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDedicatedInferenceGPUModelConfigConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_dedicated_inference_gpu_model_config.test", "gpu_model_configs.#"),
					resource.TestCheckResourceAttrSet("data.digitalocean_dedicated_inference_gpu_model_config.test", "gpu_model_configs.0.model_slug"),
					resource.TestCheckResourceAttrSet("data.digitalocean_dedicated_inference_gpu_model_config.test", "gpu_model_configs.0.model_name"),
					resource.TestCheckResourceAttrSet("data.digitalocean_dedicated_inference_gpu_model_config.test", "gpu_model_configs.0.gpu_slugs.#"),
				),
			},
		},
	})
}

func testAccDataSourceDedicatedInferenceGPUModelConfigConfig() string {
	return `
data "digitalocean_dedicated_inference_gpu_model_config" "test" {}
`
}
