package database_test

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceDigitalOceanDatabaseCluster_Basic(t *testing.T) {
	var database godo.Database
	databaseName := acceptance.RandomTestName()

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
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_database_cluster.foobar", "project_id"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_database_cluster.foobar", "storage_size_mib"),
					testAccCheckDataSourceDigitalOceanDatabaseClusterMetricsEndpoints("data.digitalocean_database_cluster.foobar"),
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

func testAccCheckDataSourceDigitalOceanDatabaseClusterMetricsEndpoints(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Not found: %s", resourceName)
		}

		// Check that metrics_endpoints is set and has at least one element
		count, err := strconv.Atoi(rs.Primary.Attributes["metrics_endpoints.#"])
		if err != nil {
			return fmt.Errorf("Error parsing metrics_endpoints count: %s", err)
		}
		if count == 0 {
			return fmt.Errorf("metrics_endpoints is empty")
		}

		// Check that the first endpoint is a valid URL
		firstEndpoint := rs.Primary.Attributes["metrics_endpoints.0"]
		if firstEndpoint == "" {
			return fmt.Errorf("First endpoint in metrics_endpoints is empty")
		}

		return nil
	}
}

const testAccCheckDataSourceDigitalOceanDatabaseClusterConfigBasic = `
resource "digitalocean_database_cluster" "foobar" {
  name             = "%s"
  engine           = "pg"
  version          = "15"
  size             = "db-s-1vcpu-1gb"
  region           = "nyc1"
  node_count       = 1
  tags             = ["production"]
  storage_size_mib = 10240
}
`

const testAccCheckDataSourceDigitalOceanDatabaseClusterConfigWithDatasource = `
resource "digitalocean_database_cluster" "foobar" {
  name             = "%s"
  engine           = "pg"
  version          = "15"
  size             = "db-s-1vcpu-1gb"
  region           = "nyc1"
  node_count       = 1
  tags             = ["production"]
  storage_size_mib = 10240
}

data "digitalocean_database_cluster" "foobar" {
  name = digitalocean_database_cluster.foobar.name
}
`
