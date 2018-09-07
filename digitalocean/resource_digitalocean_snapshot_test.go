package digitalocean

import (
	"context"
	"fmt"
	"log"
	"testing"

	"strings"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func init() {
	resource.AddTestSweepers("digitalocean_snapshot", &resource.Sweeper{
		Name:         "digitalocean_snapshot",
		F:            testSweepSnapshots,
		Dependencies: []string{"digitalocean_volume"},
	})
}

func testSweepSnapshots(region string) error {
	meta, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*godo.Client)

	snapshots, _, err := client.Snapshots.ListVolume(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, s := range snapshots {
		if strings.HasPrefix(s.Name, "snapshot-") {
			log.Printf("Destroying Snapshot %s", s.Name)

			if _, err := client.Snapshots.Delete(context.Background(), s.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccDigitalOceanSnapshot_Basic(t *testing.T) {
	var snapshot godo.Snapshot
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanSnapshotConfig_basic, rInt, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanSnapshotExists("digitalocean_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr(
						"digitalocean_snapshot.foobar", "name", fmt.Sprintf("snapshot-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_snapshot.foobar", "size", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_snapshot.foobar", "regions.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_snapshot.foobar", "min_disk_size", "100"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_snapshot.foobar", "volume_id"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanSnapshotExists(n string, snapshot *godo.Snapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*godo.Client)

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Snapshot ID is set")
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

func testAccCheckDigitalOceanSnapshotDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*godo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_snapshot" {
			continue
		}

		// Try to find the snapshot
		_, _, err := client.Snapshots.Get(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Snapshot still exists")
		}
	}

	return nil
}

const testAccCheckDigitalOceanSnapshotConfig_basic = `
resource "digitalocean_volume" "foo" {
  region      = "nyc1"
  name        = "volume-%d"
  size        = 100
  description = "peace makes plenty"
}

resource "digitalocean_snapshot" "foobar" {
  name = "snapshot-%d"
  volume_id = "${digitalocean_volume.foo.id}"
}`
