package nfs_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanNfsSnapshot_Basic(t *testing.T) {
	var snapshot godo.NfsSnapshot
	name := "tf-snap"
	nfsName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanNfsSnapshotDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckDigitalOceanNfsSnapshotConfig_basic(nfsName, name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanNfsIsActive("digitalocean_nfs.foo"),
					testAccCheckDigitalOceanNfsSnapshotExists("digitalocean_nfs_snapshot.foobar", &snapshot),
					testAccCheckDigitalOceanNfsSnapshotIsActive("digitalocean_nfs_snapshot.foobar"),
					resource.TestCheckResourceAttr("digitalocean_nfs_snapshot.foobar", "name", name),
					resource.TestCheckResourceAttr("digitalocean_nfs_snapshot.foobar", "region", "atl1"),
					resource.TestCheckResourceAttrSet("digitalocean_nfs_snapshot.foobar", "share_id"),
					resource.TestCheckResourceAttrSet("digitalocean_nfs_snapshot.foobar", "created_at"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanNfsSnapshotConfig_basic(nfsName, snapshotName string) string {
	return fmt.Sprintf(`
resource "digitalocean_vpc" "foobar" {
  name   = "%s"
  region = "atl1"
}

resource "digitalocean_nfs" "foo" {
  region = "atl1"
  name   = "%s"
  size   = 50
  vpc_id = digitalocean_vpc.foobar.id
}

resource "digitalocean_nfs_snapshot" "foobar" {
  name     = "%s"
  share_id = digitalocean_nfs.foo.id
  region   = "atl1"
}
`, nfsName, nfsName, snapshotName)
}

func testAccCheckDigitalOceanNfsSnapshotExists(n string, snapshot *godo.NfsSnapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No NFS Snapshot ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		foundSnapshot, _, err := client.Nfs.GetSnapshot(context.Background(), rs.Primary.ID, rs.Primary.Attributes["region"])
		if err != nil {
			return err
		}
		if foundSnapshot.ID != rs.Primary.ID {
			return fmt.Errorf("NFS Snapshot not found")
		}

		*snapshot = *foundSnapshot
		return nil
	}
}

func testAccCheckDigitalOceanNfsSnapshotIsActive(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No NFS Snapshot ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		for i := 0; i < 20; i++ {
			snapshot, _, err := client.Nfs.GetSnapshot(context.Background(), rs.Primary.ID, rs.Primary.Attributes["region"])
			if err != nil {
				return err
			}
			if snapshot.Status == "SNAPSHOT_AVAILABLE" {
				time.Sleep(5 * time.Second)
				return nil
			}
			time.Sleep(10 * time.Second)
		}
		return fmt.Errorf("NFS Snapshot did not become AVAILABLE in time")
	}
}

func testAccCheckDigitalOceanNfsSnapshotDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_nfs_snapshot" {
			continue
		}

		// Wait for snapshot to be fully deleted (up to 60 seconds)
		for i := 0; i < 30; i++ {
			_, resp, err := client.Nfs.GetSnapshot(context.Background(), rs.Primary.ID, rs.Primary.Attributes["region"])
			if err != nil && resp != nil && resp.StatusCode == 404 {
				return nil
			}
			time.Sleep(2 * time.Second)
		}
		return fmt.Errorf("NFS Snapshot still exists")
	}
	return nil
}
