package volume_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceDigitalOceanVolume_Basic(t *testing.T) {
	var volume godo.Volume
	testName := acceptance.RandomTestName()
	resourceConfig := testAccCheckDataSourceDigitalOceanVolumeConfig_basic(testName)
	dataSourceConfig := `
data "digitalocean_volume" "foobar" {
  name = digitalocean_volume.foo.name
}`

	expectedURNRegEx, _ := regexp.Compile(`do:volume:[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`)

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
					testAccCheckDataSourceDigitalOceanVolumeExists("data.digitalocean_volume.foobar", &volume),
					resource.TestCheckResourceAttr(
						"data.digitalocean_volume.foobar", "name", fmt.Sprintf("%s-volume", testName)),
					resource.TestCheckResourceAttr(
						"data.digitalocean_volume.foobar", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_volume.foobar", "size", "10"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_volume.foobar", "droplet_ids.#", "0"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_volume.foobar", "tags.#", "2"),
					resource.TestMatchResourceAttr("data.digitalocean_volume.foobar", "urn", expectedURNRegEx),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanVolume_RegionScoped(t *testing.T) {
	var volume godo.Volume
	testName := acceptance.RandomTestName()
	resourceConfig := testAccCheckDataSourceDigitalOceanVolumeConfig_region_scoped(testName)
	dataSourceConfig := `
data "digitalocean_volume" "foobar" {
  name   = digitalocean_volume.foo.name
  region = "lon1"
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
					testAccCheckDataSourceDigitalOceanVolumeExists("data.digitalocean_volume.foobar", &volume),
					resource.TestCheckResourceAttr(
						"data.digitalocean_volume.foobar", "name", fmt.Sprintf("%s-volume", testName)),
					resource.TestCheckResourceAttr(
						"data.digitalocean_volume.foobar", "region", "lon1"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_volume.foobar", "size", "20"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_volume.foobar", "droplet_ids.#", "0"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_volume.foobar", "tags.#", "0"),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanVolumeExists(n string, volume *godo.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Volume ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		foundVolume, _, err := client.Storage.GetVolume(context.Background(), rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundVolume.ID != rs.Primary.ID {
			return fmt.Errorf("Volume not found")
		}

		*volume = *foundVolume

		return nil
	}
}

func testAccCheckDataSourceDigitalOceanVolumeConfig_basic(testName string) string {
	return fmt.Sprintf(`
resource "digitalocean_volume" "foo" {
  region = "nyc3"
  name   = "%s-volume"
  size   = 10
  tags   = ["foo","bar"]
}`, testName)
}

func testAccCheckDataSourceDigitalOceanVolumeConfig_region_scoped(testName string) string {
	return fmt.Sprintf(`
resource "digitalocean_volume" "foo" {
  region = "nyc3"
  name   = "%s-volume"
  size   = 10
  tags   = ["foo","bar"]
}

resource "digitalocean_volume" "bar" {
  region = "lon1"
  name   = "%s-volume"
  size   = 20
}`, testName, testName)
}
