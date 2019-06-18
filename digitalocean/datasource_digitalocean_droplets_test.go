package digitalocean

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceDigitalOceanDroplet_Basic(t *testing.T) {
	var droplet godo.Droplet
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDigitalOceanDropletConfig_basic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanDropletExists("data.digitalocean_droplet.foobar", &droplet),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet.foobar", "image", "centos-7-x64"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet.foobar", "ipv6", "true"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet.foobar", "private_networking", "false"),
					resource.TestCheckResourceAttrSet("data.digitalocean_droplet.foobar", "urn"),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanDropletExists(n string, droplet *godo.Droplet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No droplet ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		foundDroplet, _, err := client.Droplets.Get(context.Background(), id)

		if err != nil {
			return err
		}

		if foundDroplet.ID != id {
			return fmt.Errorf("Droplet not found")
		}

		*droplet = *foundDroplet

		return nil
	}
}

func testAccCheckDataSourceDigitalOceanDropletConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foo" {
  name   = "%s"
  size   = "512mb"
  image  = "centos-7-x64"
  region = "nyc3"
  ipv6   = true
}

data "digitalocean_droplet" "foobar" {
  name = "${digitalocean_droplet.foo.name}"
}
`, name)
}
