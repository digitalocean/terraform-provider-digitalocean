package snapshot_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanDropletSnapshot_Basic(t *testing.T) {
	var snapshot godo.Snapshot
	rInt1 := acctest.RandInt()
	rInt2 := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDropletSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDropletSnapshotConfig_basic, rInt1, rInt2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDropletSnapshotExists("digitalocean_droplet_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr(
						"digitalocean_droplet_snapshot.foobar", "name", fmt.Sprintf("snapshot-one-%d", rInt2)),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanDropletSnapshotExists(n string, snapshot *godo.Snapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

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
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

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
resource "digitalocean_droplet" "foo" {
  name      = "foo-%d"
  size      = "s-1vcpu-1gb"
  image     = "ubuntu-22-04-x64"
  region    = "nyc3"
  user_data = "foobar"
}

resource "digitalocean_droplet_snapshot" "foobar" {
  droplet_id = digitalocean_droplet.foo.id
  name       = "snapshot-one-%d"
}
  `
