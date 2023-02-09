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

func TestAccDigitalOceanVolumeSnapshot_Basic(t *testing.T) {
	var snapshot godo.Snapshot
	volName := acceptance.RandomTestName("volume")
	snapName := acceptance.RandomTestName("snapshot")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVolumeSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanVolumeSnapshotConfig_basic, volName, snapName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVolumeSnapshotExists("digitalocean_volume_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr(
						"digitalocean_volume_snapshot.foobar", "name", snapName),
					resource.TestCheckResourceAttr(
						"digitalocean_volume_snapshot.foobar", "size", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_volume_snapshot.foobar", "regions.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_volume_snapshot.foobar", "min_disk_size", "100"),
					resource.TestCheckResourceAttr(
						"digitalocean_volume_snapshot.foobar", "tags.#", "2"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_volume_snapshot.foobar", "volume_id"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanVolumeSnapshotExists(n string, snapshot *godo.Snapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Volume Snapshot ID is set")
		}

		foundSnapshot, _, err := client.Snapshots.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundSnapshot.ID != rs.Primary.ID {
			return fmt.Errorf("Volume Snapshot not found")
		}

		*snapshot = *foundSnapshot

		return nil
	}
}

func testAccCheckDigitalOceanVolumeSnapshotDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_volume_snapshot" {
			continue
		}

		// Try to find the snapshot
		_, _, err := client.Snapshots.Get(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Volume Snapshot still exists")
		}
	}

	return nil
}

const testAccCheckDigitalOceanVolumeSnapshotConfig_basic = `
resource "digitalocean_volume" "foo" {
  region      = "nyc1"
  name        = "%s"
  size        = 100
  description = "peace makes plenty"
}

resource "digitalocean_volume_snapshot" "foobar" {
  name      = "%s"
  volume_id = digitalocean_volume.foo.id
  tags      = ["foo", "bar"]
}`

func TestAccDigitalOceanVolumeSnapshot_UpdateTags(t *testing.T) {
	var snapshot godo.Snapshot
	volName := acceptance.RandomTestName("volume")
	snapName := acceptance.RandomTestName("snapshot")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVolumeSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanVolumeSnapshotConfig_basic, volName, snapName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVolumeSnapshotExists("digitalocean_volume_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr("digitalocean_volume_snapshot.foobar", "tags.#", "2"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanVolumeSnapshotConfig_basic_tag_update, volName, snapName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVolumeSnapshotExists("digitalocean_volume_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr("digitalocean_volume_snapshot.foobar", "tags.#", "3"),
				),
			},
		},
	})
}

const testAccCheckDigitalOceanVolumeSnapshotConfig_basic_tag_update = `
resource "digitalocean_volume" "foo" {
  region      = "nyc1"
  name        = "%s"
  size        = 100
  description = "peace makes plenty"
}

resource "digitalocean_volume_snapshot" "foobar" {
  name      = "%s"
  volume_id = digitalocean_volume.foo.id
  tags      = ["foo", "bar", "baz"]
}`
