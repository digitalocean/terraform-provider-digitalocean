package digitalocean

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceDigitalOceanDroplet_BasicByName(t *testing.T) {
	var droplet godo.Droplet
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDigitalOceanDropletConfig_basicByName(name),
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
					resource.TestCheckResourceAttrSet("data.digitalocean_droplet.foobar", "created_at"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanDroplet_BasicByTag(t *testing.T) {
	var droplet godo.Droplet
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandString(10))
	tagName := fmt.Sprintf("tf-acc-test-tag-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDigitalOceanDropletConfig_basicWithTag(tagName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletExists("digitalocean_droplet.foo", &droplet),
				),
			},
			{
				Config: testAccCheckDataSourceDigitalOceanDropletConfig_basicByTag(tagName, name),
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
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet.foobar", "tags.#", "1"),
					resource.TestCheckResourceAttrSet("data.digitalocean_droplet.foobar", "urn"),
					resource.TestCheckResourceAttrSet("data.digitalocean_droplet.foobar", "created_at"),
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

func testAccCheckDataSourceDigitalOceanDropletConfig_basicByName(name string) string {
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

func testAccCheckDataSourceDigitalOceanDropletConfig_basicWithTag(tagName string, name string) string {
	return fmt.Sprintf(`
resource "digitalocean_tag" "foo" {
  name = "%s"
}

resource "digitalocean_droplet" "foo" {
  name   = "%s"
  size   = "512mb"
  image  = "centos-7-x64"
  region = "nyc3"
  ipv6   = true
  tags   = ["${digitalocean_tag.foo.id}"]
}
`, tagName, name)
}

func testAccCheckDataSourceDigitalOceanDropletConfig_basicByTag(tagName string, name string) string {
	return fmt.Sprintf(`
resource "digitalocean_tag" "foo" {
  name = "%s"
}

resource "digitalocean_droplet" "foo" {
  name   = "%s"
  size   = "512mb"
  image  = "centos-7-x64"
  region = "nyc3"
  ipv6   = true
  tags   = ["${digitalocean_tag.foo.id}"]
}

data "digitalocean_droplet" "foobar" {
  tag = "${digitalocean_tag.foo.id}"
}
`, tagName, name)
}
