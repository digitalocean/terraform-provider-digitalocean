package database_test

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

func TestAccDataSourceDigitalOceanDatabaseCluster_Basic(t *testing.T) {
	var database godo.Database
	databaseName := fmt.Sprintf("foobar-test-terraform-%s", acctest.RandString(10))

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
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
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_database_cluster.foobar", "password"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_database_cluster.foobar", "private_network_uuid"),
					testAccCheckDigitalOceanDatabaseClusterURIPassword(
						"digitalocean_database_cluster.foobar", "uri"),
					testAccCheckDigitalOceanDatabaseClusterURIPassword(
						"digitalocean_database_cluster.foobar", "private_uri"),
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

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

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
  name = digitalocean_database_cluster.foobar.name
}
`
