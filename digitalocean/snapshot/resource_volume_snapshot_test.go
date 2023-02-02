package snapshot_test

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("digitalocean_volume_snapshot", &resource.Sweeper{
		Name:         "digitalocean_volume_snapshot",
		F:            testSweepVolumeSnapshots,
		Dependencies: []string{"digitalocean_volume"},
	})
}

func testSweepVolumeSnapshots(region string) error {
	meta, err := acceptance.SharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*config.CombinedConfig).GodoClient()

	opt := &godo.ListOptions{PerPage: 200}
	snapshots, _, err := client.Snapshots.ListVolume(context.Background(), opt)
	if err != nil {
		return err
	}

	for _, s := range snapshots {
		if strings.HasPrefix(s.Name, "snapshot-") {
			log.Printf("Destroying Volume Snapshot %s", s.Name)

			if _, err := client.Snapshots.Delete(context.Background(), s.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccDigitalOceanVolumeSnapshot_Basic(t *testing.T) {
	var snapshot godo.Snapshot
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVolumeSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanVolumeSnapshotConfig_basic, rInt, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVolumeSnapshotExists("digitalocean_volume_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr(
						"digitalocean_volume_snapshot.foobar", "name", fmt.Sprintf("snapshot-%d", rInt)),
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
  name        = "volume-%d"
  size        = 100
  description = "peace makes plenty"
}

resource "digitalocean_volume_snapshot" "foobar" {
  name = "snapshot-%d"
  volume_id = "${digitalocean_volume.foo.id}"
  tags = ["foo","bar"]
}`

func TestAccDigitalOceanVolumeSnapshot_UpdateTags(t *testing.T) {
	var snapshot godo.Snapshot
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanVolumeSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanVolumeSnapshotConfig_basic, rInt, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanVolumeSnapshotExists("digitalocean_volume_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr("digitalocean_volume_snapshot.foobar", "tags.#", "2"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanVolumeSnapshotConfig_basic_tag_update, rInt, rInt),
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
  name        = "volume-%d"
  size        = 100
  description = "peace makes plenty"
}

resource "digitalocean_volume_snapshot" "foobar" {
  name = "snapshot-%d"
  volume_id = "${digitalocean_volume.foo.id}"
  tags = ["foo","bar","baz"]
}`
