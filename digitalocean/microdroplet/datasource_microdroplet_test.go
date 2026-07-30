package microdroplet_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanMicroDroplet_ByID(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccMicroDropletConfig_Basic, name, testMicroDropletImage)
	dataSourceConfig := `
data "digitalocean_microdroplet" "byid" {
  id = digitalocean_microdroplet.foobar.id
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{Config: resourceConfig},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_microdroplet.byid", "name", name),
					resource.TestCheckResourceAttrSet("data.digitalocean_microdroplet.byid", "urn"),
					resource.TestCheckResourceAttrSet("data.digitalocean_microdroplet.byid", "endpoint"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanMicroDroplet_ByName(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccMicroDropletConfig_Basic, name, testMicroDropletImage)
	dataSourceConfig := `
data "digitalocean_microdroplet" "byname" {
  name = digitalocean_microdroplet.foobar.name
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{Config: resourceConfig},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_microdroplet.byname", "name", name),
					resource.TestCheckResourceAttrSet("data.digitalocean_microdroplet.byname", "id"),
				),
			},
		},
	})
}
