package microdroplet_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanMicroDropletImages_All(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccMicroDropletImageConfig_Basic, name, testMicroDropletImage)
	dataSourceConfig := `
data "digitalocean_microdroplet_images" "all" {
  depends_on = [digitalocean_microdroplet_image.foobar]
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{Config: resourceConfig},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.digitalocean_microdroplet_images.all", "micro_droplet_images.#"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanMicroDropletImages_Filter(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccMicroDropletImageConfig_Basic, name, testMicroDropletImage)
	dataSourceConfig := fmt.Sprintf(`
data "digitalocean_microdroplet_images" "byname" {
  depends_on = [digitalocean_microdroplet_image.foobar]

  filter {
    key    = "name"
    values = ["%s"]
  }
}`, name)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{Config: resourceConfig},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_microdroplet_images.byname", "micro_droplet_images.#", "1"),
					resource.TestCheckResourceAttr("data.digitalocean_microdroplet_images.byname", "micro_droplet_images.0.name", name),
				),
			},
		},
	})
}
