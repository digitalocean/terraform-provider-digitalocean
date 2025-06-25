package genai_test

import (
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanFunctionRoute_Basic(t *testing.T) {
	var agent godo.Agent
	function_name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanAgentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanAgentConfig_basic(function_name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanAgentExists("digitalocean_genai_function.test", &agent),
					resource.TestCheckResourceAttr("digitalocean_genai_function.test", "name", function_name),
					resource.TestCheckResourceAttr("digitalocean_genai_function.test", "description", "Adding a function route and this will also tell temperature"),
					resource.TestCheckResourceAttrSet("digitalocean_genai_function.test", "faas_name"),
					resource.TestCheckResourceAttrSet("digitalocean_genai_function.test", "faas_namespace"),
					resource.TestCheckResourceAttrSet("digitalocean_genai_function.test", "input_schema"),
					resource.TestCheckResourceAttrSet("digitalocean_genai_function.test", "output_schema"),
				),
			},
		},
	})
}
