package microdroplet_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanMicroDroplets_All(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccMicroDropletConfig_Basic, name, testMicroDropletImage)
	dataSourceConfig := `
data "digitalocean_microdroplets" "all" {
  depends_on = [digitalocean_microdroplet.foobar]
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{Config: resourceConfig},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_microdroplets.all", "micro_droplets.#"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanMicroDroplets_ByRegion(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccMicroDropletConfig_Basic, name, testMicroDropletImage)
	dataSourceConfig := `
data "digitalocean_microdroplets" "region" {
  region     = "nyc3"
  depends_on = [digitalocean_microdroplet.foobar]
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{Config: resourceConfig},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_microdroplets.region", "micro_droplets.#"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanMicroDroplets_ByName(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccMicroDropletConfig_Basic, name, testMicroDropletImage)
	dataSourceConfig := fmt.Sprintf(`
data "digitalocean_microdroplets" "byname" {
  name       = "%s"
  depends_on = [digitalocean_microdroplet.foobar]
}`, name)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{Config: resourceConfig},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_microdroplets.byname", "micro_droplets.#", "1"),
					resource.TestCheckResourceAttr("data.digitalocean_microdroplets.byname", "micro_droplets.0.name", name),
				),
			},
		},
	})
}
