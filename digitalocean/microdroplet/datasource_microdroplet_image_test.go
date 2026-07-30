package microdroplet_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanMicroDropletImage_ByID(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccMicroDropletImageConfig_Basic, name, testMicroDropletImage)
	dataSourceConfig := `
data "digitalocean_microdroplet_image" "byid" {
  id = digitalocean_microdroplet_image.foobar.id
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{Config: resourceConfig},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_microdroplet_image.byid", "name", name),
					resource.TestCheckResourceAttrSet("data.digitalocean_microdroplet_image.byid", "urn"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanMicroDropletImage_ByName(t *testing.T) {
	name := acceptance.RandomTestName()
	resourceConfig := fmt.Sprintf(testAccMicroDropletImageConfig_Basic, name, testMicroDropletImage)
	dataSourceConfig := `
data "digitalocean_microdroplet_image" "byname" {
  name = digitalocean_microdroplet_image.foobar.name
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{Config: resourceConfig},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_microdroplet_image.byname", "name", name),
					resource.TestCheckResourceAttrSet("data.digitalocean_microdroplet_image.byname", "id"),
				),
			},
		},
	})
}
