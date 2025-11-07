package nfs_test

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

func TestAccDataSourceDigitalOceanNfsSnapshot_basic(t *testing.T) {
	var snapshot godo.NfsSnapshot
	vpcName := acceptance.RandomTestName()
	nfsName := acceptance.RandomTestName()
	snapName := "tf-snap"

	nfsConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanNfsSnapshot_nfs, vpcName, nfsName)
	snapshotConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanNfsSnapshot_snapshot, snapName)
	dataSourceConfig := `
data "digitalocean_nfs_snapshot" "foobar" {
  name      = digitalocean_nfs_snapshot.foo.name
  share_id  = digitalocean_nfs.foo.id
  region    = "atl1"
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: nfsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanNfsIsActive("digitalocean_nfs.foo"),
				),
			},
			{
				Config: nfsConfig + snapshotConfig,
			},
			{
				Config: nfsConfig + snapshotConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanNfsSnapshotExists("data.digitalocean_nfs_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr("data.digitalocean_nfs_snapshot.foobar", "name", snapName),
					resource.TestCheckResourceAttr("data.digitalocean_nfs_snapshot.foobar", "region", "atl1"),
					resource.TestCheckResourceAttrSet("data.digitalocean_nfs_snapshot.foobar", "share_id"),
					resource.TestCheckResourceAttrSet("data.digitalocean_nfs_snapshot.foobar", "size"),
					resource.TestCheckResourceAttrSet("data.digitalocean_nfs_snapshot.foobar", "created_at"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanNfsSnapshot_regex(t *testing.T) {
	var snapshot godo.NfsSnapshot
	vpcName := acceptance.RandomTestName()
	nfsName := acceptance.RandomTestName()
	snapName := "tf-snap"

	nfsConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanNfsSnapshot_nfs, vpcName, nfsName)
	snapshotConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanNfsSnapshot_snapshot, snapName)
	dataSourceConfig := `
data "digitalocean_nfs_snapshot" "foobar" {
 name_regex = "^${digitalocean_nfs_snapshot.foo.name}"
 share_id   = digitalocean_nfs.foo.id
 region     = "atl1"
}`

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: nfsConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanNfsIsActive("digitalocean_nfs.foo"),
				),
			},
			{
				Config: nfsConfig + snapshotConfig,
			},
			{
				Config: nfsConfig + snapshotConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanNfsSnapshotExists("data.digitalocean_nfs_snapshot.foobar", &snapshot),
					resource.TestCheckResourceAttr("data.digitalocean_nfs_snapshot.foobar", "name", snapName),
					resource.TestCheckResourceAttr("data.digitalocean_nfs_snapshot.foobar", "region", "atl1"),
					resource.TestCheckResourceAttrSet("data.digitalocean_nfs_snapshot.foobar", "share_id"),
					resource.TestCheckResourceAttrSet("data.digitalocean_nfs_snapshot.foobar", "created_at"),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanNfsSnapshotExists(n string, snapshot *godo.NfsSnapshot) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No snapshot ID is set")
		}

		region := rs.Primary.Attributes["region"]
		foundSnapshot, _, err := client.Nfs.GetSnapshot(context.Background(), rs.Primary.ID, region)
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

const testAccCheckDataSourceDigitalOceanNfsSnapshot_nfs = `
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
`

const testAccCheckDataSourceDigitalOceanNfsSnapshot_snapshot = `
resource "digitalocean_nfs_snapshot" "foo" {
  name     = "%s"
  share_id = digitalocean_nfs.foo.id
  region   = "atl1"
}
`
