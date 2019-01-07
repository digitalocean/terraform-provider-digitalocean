package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceDigitalOceanVolumeSnapshot_basic(t *testing.T) {
	var snapshot godo.Snapshot
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanVolumeSnapshot_basic, rInt, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanVolumeSnapshotExists("data.digitalocean_volume_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "name", fmt.Sprintf("snapshot-%d", rInt)),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "size", "0"),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "min_disk_size", "100"),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "regions.#", "1"),
					resource.TestCheckResourceAttrSet("data.digitalocean_volume_snapshot.foobar", "volume_id"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanVolumeSnapshot_regex(t *testing.T) {
	var snapshot godo.Snapshot
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanVolumeSnapshot_regex, rInt, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanVolumeSnapshotExists("data.digitalocean_volume_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "name", fmt.Sprintf("snapshot-%d", rInt)),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "size", "0"),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "min_disk_size", "100"),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "regions.#", "1"),
					resource.TestCheckResourceAttrSet("data.digitalocean_volume_snapshot.foobar", "volume_id"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanVolumeSnapshot_region(t *testing.T) {
	var snapshot godo.Snapshot
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanVolumeSnapshot_region, rInt, rInt, rInt, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanVolumeSnapshotExists("data.digitalocean_volume_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "name", fmt.Sprintf("snapshot-%d", rInt)),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "size", "0"),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "min_disk_size", "100"),
					resource.TestCheckResourceAttr("data.digitalocean_volume_snapshot.foobar", "regions.#", "1"),
					resource.TestCheckResourceAttrSet("data.digitalocean_volume_snapshot.foobar", "volume_id"),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanVolumeSnapshotExists(n string, snapshot *godo.Snapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

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
  name        = "volume-%d"
  size        = 100
  description = "peace makes plenty"
}

resource "digitalocean_volume_snapshot" "foo" {
  name = "snapshot-%d"
  volume_id = "${digitalocean_volume.foo.id}"
}

data "digitalocean_volume_snapshot" "foobar" {
  most_recent = true
  name = "${digitalocean_volume_snapshot.foo.name}"
}`

const testAccCheckDataSourceDigitalOceanVolumeSnapshot_regex = `
resource "digitalocean_volume" "foo" {
  region      = "nyc1"
  name        = "volume-%d"
  size        = 100
  description = "peace makes plenty"
}

resource "digitalocean_volume_snapshot" "foo" {
  name = "snapshot-%d"
  volume_id = "${digitalocean_volume.foo.id}"
}

data "digitalocean_volume_snapshot" "foobar" {
  most_recent = true
  name_regex = "^${digitalocean_volume_snapshot.foo.name}"
}`

const testAccCheckDataSourceDigitalOceanVolumeSnapshot_region = `
resource "digitalocean_volume" "foo" {
  region      = "nyc1"
  name        = "volume-nyc-%d"
  size        = 100
  description = "peace makes plenty"
}

resource "digitalocean_volume" "bar" {
  region      = "lon1"
  name        = "volume-lon-%d"
  size        = 100
  description = "peace makes plenty"
}

resource "digitalocean_volume_snapshot" "foo" {
  name = "snapshot-%d"
  volume_id = "${digitalocean_volume.foo.id}"
}

resource "digitalocean_volume_snapshot" "bar" {
  name = "snapshot-%d"
  volume_id = "${digitalocean_volume.bar.id}"
}

data "digitalocean_volume_snapshot" "foobar" {
  name = "${digitalocean_volume_snapshot.bar.name}"
  region = "lon1"
}`
