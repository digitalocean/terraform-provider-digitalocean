package droplet_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanDroplets_Basic(t *testing.T) {
	name1 := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(10))
	name2 := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(10))

	resourcesConfig := fmt.Sprintf(`
resource "digitalocean_droplet" "foo" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
}

resource "digitalocean_droplet" "bar" {
  name   = "%s"
  size   = "s-1vcpu-1gb"
  image  = "ubuntu-22-04-x64"
  region = "nyc3"
}
`, name1, name2)

	datasourceConfig := fmt.Sprintf(`
data "digitalocean_droplets" "result" {
  filter {
    key    = "name"
    values = ["%s"]
  }
}
`, name1)
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourcesConfig,
			},
			{
				Config: resourcesConfig + datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.digitalocean_droplets.result", "droplets.#", "1"),
					resource.TestCheckResourceAttr("data.digitalocean_droplets.result", "droplets.0.name", name1),
					resource.TestCheckResourceAttrPair("data.digitalocean_droplets.result", "droplets.0.id", "digitalocean_droplet.foo", "id"),
				),
			},
			{
				Config: resourcesConfig,
			},
		},
	})
}
