package nfs_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/nfs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceDigitalOceanNfsAccessPoint_byID(t *testing.T) {
	vpcName := acceptance.RandomTestName()
	nfsName := acceptance.RandomTestName()
	apName := acceptance.RandomTestName("ap")

	nfsConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanNfsAccessPoint_nfs, vpcName, nfsName)
	accessPointConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanNfsAccessPoint_accessPoint, apName)
	dataSourceConfig := `
data "digitalocean_nfs_access_point" "foobar" {
  id = digitalocean_nfs_access_point.foo.id
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
				Config: nfsConfig + accessPointConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanNfsAccessPointIsActive("digitalocean_nfs_access_point.foo"),
				),
			},
			{
				Config: nfsConfig + accessPointConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanNfsAccessPointExists("data.digitalocean_nfs_access_point.foobar"),
					resource.TestCheckResourceAttr("data.digitalocean_nfs_access_point.foobar", "name", apName),
					resource.TestCheckResourceAttr("data.digitalocean_nfs_access_point.foobar", "path", "/exports/data"),
					resource.TestCheckResourceAttrSet("data.digitalocean_nfs_access_point.foobar", "share_id"),
					resource.TestCheckResourceAttrSet("data.digitalocean_nfs_access_point.foobar", "status"),
				),
			},
		},
	})
}

func TestAccDataSourceDigitalOceanNfsAccessPoint_byName(t *testing.T) {
	vpcName := acceptance.RandomTestName()
	nfsName := acceptance.RandomTestName()
	apName := acceptance.RandomTestName("ap")

	nfsConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanNfsAccessPoint_nfs, vpcName, nfsName)
	accessPointConfig := fmt.Sprintf(testAccCheckDataSourceDigitalOceanNfsAccessPoint_accessPoint, apName)
	dataSourceConfig := `
data "digitalocean_nfs_access_point" "foobar" {
  name     = digitalocean_nfs_access_point.foo.name
  share_id = digitalocean_nfs.foo.id
  vpc_id   = digitalocean_vpc.foobar.id
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
				Config: nfsConfig + accessPointConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanNfsAccessPointIsActive("digitalocean_nfs_access_point.foo"),
				),
			},
			{
				Config: nfsConfig + accessPointConfig + dataSourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanNfsAccessPointExists("data.digitalocean_nfs_access_point.foobar"),
					resource.TestCheckResourceAttr("data.digitalocean_nfs_access_point.foobar", "name", apName),
					resource.TestCheckResourceAttrSet("data.digitalocean_nfs_access_point.foobar", "share_id"),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanNfsAccessPointExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("No NFS access point ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		accessPoint, _, err := nfs.GetNfsAccessPoint(context.Background(), client, rs.Primary.ID)
		if err != nil {
			return err
		}
		if accessPoint.ID != rs.Primary.ID {
			return fmt.Errorf("NFS access point not found")
		}

		return nil
	}
}

const testAccCheckDataSourceDigitalOceanNfsAccessPoint_nfs = `
resource "digitalocean_vpc" "foobar" {
  name   = "%s"
  region = "atl1"
}

resource "digitalocean_nfs" "foo" {
  region           = "atl1"
  name             = "%s"
  size             = 50
  vpc_id           = digitalocean_vpc.foobar.id
  performance_tier = "standard"
}
`

const testAccCheckDataSourceDigitalOceanNfsAccessPoint_accessPoint = `
resource "digitalocean_nfs_access_point" "foo" {
  name     = "%s"
  share_id = digitalocean_nfs.foo.id
  path     = "/exports/data"
  vpc_id   = digitalocean_vpc.foobar.id

  access_policy {
    protocols     = ["NFS4"]
    squash_config = "ROOT_SQUASH"
  }
}
`
