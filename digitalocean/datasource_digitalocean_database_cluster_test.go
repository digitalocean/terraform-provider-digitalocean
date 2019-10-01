package digitalocean

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/terraform"

	"github.com/digitalocean/godo"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceDigitalOceanDatabaseCluster_Basic(t *testing.T) {
	var database godo.Database
	databaseName := fmt.Sprintf("foobar-test-terraform-%s", acctest.RandString(10))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanDatabaseClusterConfigBasic, databaseName),
			},
			{
				Config: fmt.Sprintf(testAccCheckDataSourceDigitalOceanDatabaseClusterConfigWithDatasource, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceDigitalOceanDatabaseClusterExists("data.digitalocean_database_cluster.foobar", &database),
					resource.TestCheckResourceAttr(
						"data.digitalocean_database_cluster.foobar", "name", databaseName),
					resource.TestCheckResourceAttr(
						"data.digitalocean_database_cluster.foobar", "engine", "pg"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_database_cluster.foobar", "host"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_database_cluster.foobar", "private_host"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_database_cluster.foobar", "port"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_database_cluster.foobar", "user"),
				),
			},
		},
	})
}

func testAccCheckDataSourceDigitalOceanDatabaseClusterExists(n string, databaseCluster *godo.Database) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No Record ID is set")
		}

		client := testAccProvider.Meta().(*CombinedConfig).godoClient()

		foundCluster, _, err := client.Databases.Get(context.Background(), rs.Primary.ID)
		if err != nil {
			return err
		}

		if foundCluster.ID != rs.Primary.ID {
			return fmt.Errorf("DatabaseCluster not found")
		}

		*databaseCluster = *foundCluster

		return nil
	}
}

const testAccCheckDataSourceDigitalOceanDatabaseClusterConfigBasic = `
resource "digitalocean_database_cluster" "foobar" {
	name       = "%s"
	engine     = "pg"
	version    = "11"
	size       = "db-s-1vcpu-1gb"
	region     = "nyc1"
    node_count = 1
	tags       = ["production"]
}
`

const testAccCheckDataSourceDigitalOceanDatabaseClusterConfigWithDatasource = `
resource "digitalocean_database_cluster" "foobar" {
	name       = "%s"
	engine     = "pg"
	version    = "11"
	size       = "db-s-1vcpu-1gb"
	region     = "nyc1"
    node_count = 1
	tags       = ["production"]
}

data "digitalocean_database_cluster" "foobar" {
   name = "${digitalocean_database_cluster.foobar.name}"
}
`
