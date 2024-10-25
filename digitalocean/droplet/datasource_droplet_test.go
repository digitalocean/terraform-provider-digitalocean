package droplet_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceDigitalOceanDroplet_BasicByName(t *testing.T) {
	var droplet godo.Droplet
	name := acceptance.RandomTestName()
	resourceConfig := testAccCheckDataSourceDigitalOceanDropletConfig_basicByName(name)
	dataSourceConfig := `
data "digitalocean_droplet" "foobar" {
  name = digitalocean_droplet.foo.name
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanDropletExists("data.digitalocean_droplet.foobar", &droplet),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet.foobar", "image", "ubuntu-22-04-x64"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet.foobar", "ipv6", "true"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet.foobar", "private_networking", "true"),
					resource.TestCheckResourceAttrSet("data.digitalocean_droplet.foobar", "urn"),
					resource.TestCheckResourceAttrSet("data.digitalocean_droplet.foobar", "created_at"),
					resource.TestCheckResourceAttrSet("data.digitalocean_droplet.foobar", "vpc_uuid"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanDroplet_GPUByName(t *testing.T) {
	runGPU := os.Getenv(runGPUEnvVar)
	if runGPU == "" {
		t.Skip("'DO_RUN_GPU_TESTS' env var not set; Skipping tests that requires a GPU Droplet")
	}

	keyName := acceptance.RandomTestName()
	publicKeyMaterial, _, err := acctest.RandSSHKeyPair("digitalocean@ssh-acceptance-test")
	if err != nil {
		t.Fatalf("Cannot generate test SSH key pair: %s", err)
	}

	var droplet godo.Droplet
	name := acceptance.RandomTestName()
	resourceConfig := testAccCheckDataSourceDigitalOceanDropletConfig_gpuByName(keyName, publicKeyMaterial, name)
	dataSourceConfig := `
data "digitalocean_droplet" "foobar" {
  name = digitalocean_droplet.foo.name
  gpu  = true
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
			},
			{
				Config: resourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanDropletExists("data.digitalocean_droplet.foobar", &droplet),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet.foobar", "image", gpuImage),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet.foobar", "region", "tor1"),
					resource.TestCheckResourceAttrSet("data.digitalocean_droplet.foobar", "urn"),
					resource.TestCheckResourceAttrSet("data.digitalocean_droplet.foobar", "created_at"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanDroplet_BasicById(t *testing.T) {
	var droplet godo.Droplet
	name := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDigitalOceanDropletConfig_basicById(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanDropletExists("data.digitalocean_droplet.foobar", &droplet),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet.foobar", "image", "ubuntu-22-04-x64"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet.foobar", "ipv6", "true"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet.foobar", "private_networking", "true"),
					resource.TestCheckResourceAttrSet("data.digitalocean_droplet.foobar", "urn"),
					resource.TestCheckResourceAttrSet("data.digitalocean_droplet.foobar", "created_at"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanDroplet_BasicByTag(t *testing.T) {
	var droplet godo.Droplet
	name := acceptance.RandomTestName()
	tagName := acceptance.RandomTestName("tag")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDigitalOceanDropletConfig_basicWithTag(tagName, name),
				Check: resource.ComposeTestCheckFunc(
					acceptance.TestAccCheckDigitalOceanDropletExists("digitalocean_droplet.foo", &droplet),
				),
			},
			{
				Config: testAccCheckDataSourceDigitalOceanDropletConfig_basicByTag(tagName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanDropletExists("data.digitalocean_droplet.foobar", &droplet),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet.foobar", "name", name),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet.foobar", "image", "ubuntu-22-04-x64"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet.foobar", "ipv6", "true"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_droplet.foobar", "private_networking", "true"),
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

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

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
resource "digitalocean_vpc" "foobar" {
  name   = "%s"
  region = "nyc3"
}

resource "digitalocean_droplet" "foo" {
  name     = "%s"
  size     = "%s"
  image    = "%s"
  region   = "nyc3"
  ipv6     = true
  vpc_uuid = digitalocean_vpc.foobar.id
}`, acceptance.RandomTestName(), name, defaultSize, defaultImage)
}

func testAccCheckDataSourceDigitalOceanDropletConfig_gpuByName(keyName, key, name string) string {
	return fmt.Sprintf(`
resource "digitalocean_ssh_key" "foobar" {
  name       = "%s"
  public_key = "%s"
}

resource "digitalocean_droplet" "foo" {
  name     = "%s"
  size     = "%s"
  image    = "%s"
  region   = "tor1"
  ssh_keys = [digitalocean_ssh_key.foobar.id]
}`, keyName, key, name, gpuSize, gpuImage)
}

func testAccCheckDataSourceDigitalOceanDropletConfig_basicById(name string) string {
	return fmt.Sprintf(`
resource "digitalocean_droplet" "foo" {
  name   = "%s"
  size   = "%s"
  image  = "%s"
  region = "nyc3"
  ipv6   = true
}

data "digitalocean_droplet" "foobar" {
  id = digitalocean_droplet.foo.id
}
`, name, defaultSize, defaultImage)
}

func testAccCheckDataSourceDigitalOceanDropletConfig_basicWithTag(tagName string, name string) string {
	return fmt.Sprintf(`
resource "digitalocean_tag" "foo" {
  name = "%s"
}

resource "digitalocean_droplet" "foo" {
  name   = "%s"
  size   = "%s"
  image  = "%s"
  region = "nyc3"
  ipv6   = true
  tags   = [digitalocean_tag.foo.id]
}
`, tagName, name, defaultSize, defaultImage)
}

func testAccCheckDataSourceDigitalOceanDropletConfig_basicByTag(tagName string, name string) string {
	return fmt.Sprintf(`
resource "digitalocean_tag" "foo" {
  name = "%s"
}

resource "digitalocean_droplet" "foo" {
  name   = "%s"
  size   = "%s"
  image  = "%s"
  region = "nyc3"
  ipv6   = true
  tags   = [digitalocean_tag.foo.id]
}

data "digitalocean_droplet" "foobar" {
  tag = digitalocean_tag.foo.id
}
`, tagName, name, defaultSize, defaultImage)
}
