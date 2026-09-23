package microvm_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanMicroVM_ByID(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccMicroVMConfig_Basic, name, testMicroVMOCIRef)
	dataSourceConfig := `
data "digitalocean_microvm" "byid" {
  id = digitalocean_microvm.foobar.id
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{Config: resourceConfig},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_microvm.byid", "name", name),
					resource.TestCheckResourceAttrSet("data.digitalocean_microvm.byid", "urn"),
					resource.TestCheckResourceAttrSet("data.digitalocean_microvm.byid", "urls.0.hostname"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanMicroVM_ByName(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccMicroVMConfig_Basic, name, testMicroVMOCIRef)
	dataSourceConfig := `
data "digitalocean_microvm" "byname" {
  name = digitalocean_microvm.foobar.name
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{Config: resourceConfig},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_microvm.byname", "name", name),
					resource.TestCheckResourceAttrSet("data.digitalocean_microvm.byname", "id"),
				),
			},
		},
	})
}
