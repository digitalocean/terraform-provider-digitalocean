package digitalocean

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceDigitalOceanVolume_Basic(t *testing.T) {
	var volume godo.Volume
	rInt := acctest.RandInt()

	expectedURNRegEx, _ := regexp.Compile(`do:volume:[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}`)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDigitalOceanVolumeConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanVolumeExists("data.digitalocean_volume.foobar", &volume),
					resource.TestCheckResourceAttr(
						"data.digitalocean_volume.foobar", "name", fmt.Sprintf("volume-%d", rInt)),
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
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDataSourceDigitalOceanVolumeConfig_region_scoped(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanVolumeExists("data.digitalocean_volume.foobar", &volume),
					resource.TestCheckResourceAttr(
						"data.digitalocean_volume.foobar", "name", fmt.Sprintf("volume-%d", rInt)),
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

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

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

func testAccCheckDataSourceDigitalOceanVolumeConfig_basic(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_volume" "foo" {
  region = "nyc3"
  name   = "volume-%d"
  size   = 10
  tags   = ["foo","bar"]
}

data "digitalocean_volume" "foobar" {
  name = "${digitalocean_volume.foo.name}"
}`, rInt)
}

func testAccCheckDataSourceDigitalOceanVolumeConfig_region_scoped(rInt int) string {
	return fmt.Sprintf(`
resource "digitalocean_volume" "foo" {
  region = "nyc3"
  name   = "volume-%d"
  size   = 10
  tags   = ["foo","bar"]
}

resource "digitalocean_volume" "bar" {
  region = "lon1"
  name   = "volume-%d"
  size   = 20
}

data "digitalocean_volume" "foobar" {
  name   = "${digitalocean_volume.foo.name}"
  region = "lon1"
}`, rInt, rInt)
}
