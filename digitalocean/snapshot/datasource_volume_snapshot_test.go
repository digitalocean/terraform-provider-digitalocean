package snapshot_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceDigitalOceanVolumeSnapshot_basic(t *testing.T) {
	var snapshot godo.Snapshot
	volName := acceptance.RandomTestName("volume")
	snapName := acceptance.RandomTestName("snapshot")
	resourceConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanVolumeSnapshot_basic, volName, snapName)
	dataSourceConfig := `
data "digitalocean_volume_snapshot" "foobar" {
  most_recent = true
  name        = digitalocean_volume_snapshot.foo.name
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
					testAccCheckDataSourceDigitalOceanVolumeSnapshotExists("data.digitalocean_volume_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "name", snapName),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "size", "0"),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "min_disk_size", "100"),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "regions.#", "1"),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttrSet("data.digitalocean_volume_snapshot.foobar", "volume_id"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanVolumeSnapshot_regex(t *testing.T) {
	var snapshot godo.Snapshot
	volName := acceptance.RandomTestName("volume")
	snapName := acceptance.RandomTestName("snapshot")
	resourceConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanVolumeSnapshot_basic, volName, snapName)
	dataSourceConfig := `
data "digitalocean_volume_snapshot" "foobar" {
  most_recent = true
  name_regex  = "^${digitalocean_volume_snapshot.foo.name}"
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
					testAccCheckDataSourceDigitalOceanVolumeSnapshotExists("data.digitalocean_volume_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "name", snapName),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "size", "0"),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "min_disk_size", "100"),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "regions.#", "1"),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttrSet("data.digitalocean_volume_snapshot.foobar", "volume_id"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanVolumeSnapshot_region(t *testing.T) {
	var snapshot godo.Snapshot
	dropletName := acceptance.RandomTestName()
	snapName := acceptance.RandomTestName()
	nycResourceConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanVolumeSnapshot_basic, dropletName, snapName)
	lonResourceConfig := fmt.Sprintf(`
resource "digitalocean_volume" "bar" {
  region      = "lon1"
  name        = "%s"
  size        = 100
  description = "peace makes plenty"
}

resource "digitalocean_volume_snapshot" "bar" {
  name      = "%s"
  volume_id = digitalocean_volume.bar.id
}`, dropletName, snapName)
	dataSourceConfig := `
data "digitalocean_volume_snapshot" "foobar" {
  name   = digitalocean_volume_snapshot.bar.name
  region = "lon1"
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: nycResourceConfig + lonResourceConfig,
			},
			{
				Config: nycResourceConfig + lonResourceConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanVolumeSnapshotExists("data.digitalocean_volume_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "name", snapName),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "size", "0"),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "min_disk_size", "100"),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "regions.#", "1"),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "tags.#", "0"),
					resource.TestCheckResourceAttrSet("data.digitalocean_volume_snapshot.foobar", "volume_id"),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanVolumeSnapshotExists(n string, snapshot *godo.Snapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No snapshot ID is set")
		}

		foundSnapshot, _, err := client.Snapshots.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundSnapshot.ID != rs.Primary.ID {
			return fmt.Errorf("Snapshot not found")
		}

		*snapshot = *foundSnapshot

		return nil
	}
}

const testAccCheckDataSourceDigitalOceanVolumeSnapshot_basic = `
resource "digitalocean_volume" "foo" {
  region      = "nyc1"
  name        = "%s"
  size        = 100
  description = "peace makes plenty"
}

resource "digitalocean_volume_snapshot" "foo" {
  name      = "%s"
  volume_id = digitalocean_volume.foo.id
  tags      = ["foo", "bar"]
}
`
