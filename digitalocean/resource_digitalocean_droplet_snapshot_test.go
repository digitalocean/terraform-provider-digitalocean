package digitalocean

import (
	"context"
	"fmt"
	"log"
	"strings"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func init() {
	resource.AddTestSweepers("digitalocean_droplet_snapshot", &resource.Sweeper{
		Name:         "digitalocean_droplet_snapshot",
		F:            testSweepDropletSnapshots,
		Dependencies: []string{"digitalocean_droplet"},
	})
}

func testSweepDropletSnapshots(region string) error {
	meta, err := sharedConfigForRegion(region)
	if err != nil {
		return err
	}

	client := meta.(*godo.Client)

	snapshots, _, err := client.Snapshots.ListDroplet(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, s := range snapshots {
		if strings.HasPrefix(s.Name, "snapshot-") {
			log.Printf("Destroying Droplet Snapshot %s", s.Name)

			if _, err := client.Snapshots.Delete(context.Background(), s.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestAccDigitalOceanDropletSnapshot_Basic(t *testing.T) {
	var snapshot godo.Snapshot
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDigitalOceanDropletSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDropletSnapshotConfig_basic, rInt, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletSnapshotExists("digitalocean_droplet_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_snapshot.foobar", "name", fmt.Sprintf("snapshot-%d", rInt)),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_snapshot.foobar", "regions.#", "1"),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_snapshot.foobar", "min_disk_size", "20"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_droplet_snapshot.foobar", "resource_id"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanDropletSnapshotExists(n string, snapshot *godo.Snapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*godo.Client)

		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Droplet Snapshot ID is set")
		}

		foundSnapshot, _, err := client.Snapshots.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundSnapshot.ID != rs.Primary.ID {
			return fmt.Errorf("Droplet Snapshot not found")
		}

		*snapshot = *foundSnapshot

		return nil
	}
}

func testAccCheckDigitalOceanDropletSnapshotDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*godo.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_droplet_snapshot" {
			continue
		}

		// Try to find the snapshot
		_, _, err := client.Snapshots.Get(context.Background(), rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("Droplet Snapshot still exists")
		}
	}

	return nil
}

const testAccCheckDigitalOceanDropletSnapshotConfig_basic = `
resource "digitalocean_droplet_snapshot" "foobar" {
  name = "snapshot-%d"
  resource_id = "${digitalocean_droplet.foo.id}"
}
resource "digitalocean_droplet" "foo" {
	name      = "foo-%d"
	size      = "512mb"
	image     = "centos-7-x64"
	region    = "nyc3"
	user_data = "foobar"
  }
`
