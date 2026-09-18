package microvm_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanMicroVMs_All(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccMicroVMConfig_Basic, name, testMicroVMOCIRef)
	dataSourceConfig := `
data "digitalocean_microvms" "all" {
  depends_on = [digitalocean_microvm.foobar]
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{Config: resourceConfig},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_microvms.all", "micro_vms.#"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanMicroVMs_ByRegion(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccMicroVMConfig_Basic, name, testMicroVMOCIRef)
	dataSourceConfig := `
data "digitalocean_microvms" "region" {
  region     = "nyc3"
  depends_on = [digitalocean_microvm.foobar]
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{Config: resourceConfig},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_microvms.region", "micro_vms.#"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanMicroVMs_ByName(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccMicroVMConfig_Basic, name, testMicroVMOCIRef)
	dataSourceConfig := fmt.Sprintf(`
data "digitalocean_microvms" "byname" {
  name       = "%s"
  depends_on = [digitalocean_microvm.foobar]
}`, name)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{Config: resourceConfig},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_microvms.byname", "micro_vms.#", "1"),
					resource.TestCheckResourceAttr("data.digitalocean_microvms.byname", "micro_vms.0.name", name),
				),
			},
		},
	})
}
