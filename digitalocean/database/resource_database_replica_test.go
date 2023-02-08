package database_test

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

func TestAccDigitalOceanDatabaseReplica_Basic(t *testing.T) {
	var databaseReplica godo.DatabaseReplica
	var database godo.Database

	databaseName := acceptance.RandomTestName()
	databaseReplicaName := acceptance.RandomTestName()

	databaseConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigBasic, databaseName)
	replicaConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseReplicaConfigBasic, databaseReplicaName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseReplicaDestroy,
		Steps: []resource.TestStep{
			{
				Config: databaseConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
				),
			},
			{
				Config: databaseConfig + replicaConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseReplicaExists("digitalocean_database_replica.read-01", &databaseReplica),
					testAccCheckDigitalOceanDatabaseReplicaAttributes(&databaseReplica, databaseReplicaName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_replica.read-01", "size", "db-s-2vcpu-4gb"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_replica.read-01", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_replica.read-01", "name", databaseReplicaName),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_replica.read-01", "host"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_replica.read-01", "private_host"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_replica.read-01", "port"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_replica.read-01", "user"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_replica.read-01", "uri"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_replica.read-01", "private_uri"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_replica.read-01", "password"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_replica.read-01", "tags.#", "1"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_replica.read-01", "private_network_uuid"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_replica.read-01", "uuid"),
				),
			},
		},
	})
}

func TestAccDigitalOceanDatabaseReplica_WithVPC(t *testing.T) {
	var database godo.Database
	var databaseReplica godo.DatabaseReplica

	vpcName := acceptance.RandomTestName()
	databaseName := acceptance.RandomTestName()
	databaseReplicaName := acceptance.RandomTestName()

	databaseConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigWithVPC, vpcName, databaseName)
	replicaConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseReplicaConfigWithVPC, databaseReplicaName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: databaseConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
				),
			},
			{
				Config: databaseConfig + replicaConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseReplicaExists("digitalocean_database_replica.read-01", &databaseReplica),
					testAccCheckDigitalOceanDatabaseClusterAttributes(&database, databaseName),
					resource.TestCheckResourceAttrPair(
						"digitalocean_database_replica.read-01", "private_network_uuid", "digitalocean_vpc.foobar", "id"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanDatabaseReplicaDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_database_replica" {
			continue
		}
		clusterId := rs.Primary.Attributes["cluster_id"]
		name := rs.Primary.Attributes["name"]
		// Try to find the database replica
		_, _, err := client.Databases.GetReplica(context.Background(), clusterId, name)

		if err == nil {
			return fmt.Errorf("DatabaseReplica still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanDatabaseReplicaExists(n string, database *godo.DatabaseReplica) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DatabaseReplica cluster ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		clusterId := rs.Primary.Attributes["cluster_id"]
		name := rs.Primary.Attributes["name"]
		uuid := rs.Primary.Attributes["uuid"]

		foundDatabaseReplica, _, err := client.Databases.GetReplica(context.Background(), clusterId, name)

		if err != nil {
			return err
		}

		if foundDatabaseReplica.Name != name {
			return fmt.Errorf("DatabaseReplica not found")
		}

		if foundDatabaseReplica.ID != uuid {
			return fmt.Errorf("DatabaseReplica UUID not found")
		}

		*database = *foundDatabaseReplica

		return nil
	}
}

func testAccCheckDigitalOceanDatabaseReplicaAttributes(databaseReplica *godo.DatabaseReplica, name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if databaseReplica.Name != name {
			return fmt.Errorf("Bad name: %s", databaseReplica.Name)
		}

		return nil
	}
}

const testAccCheckDigitalOceanDatabaseReplicaConfigBasic = `
resource "digitalocean_database_replica" "read-01" {
  cluster_id = digitalocean_database_cluster.foobar.id
  name       = "%s"
  region     = "nyc3"
  size       = "db-s-2vcpu-4gb"
  tags       = ["staging"]
}`

const testAccCheckDigitalOceanDatabaseReplicaConfigWithVPC = `


resource "digitalocean_database_replica" "read-01" {
  cluster_id           = digitalocean_database_cluster.foobar.id
  name                 = "%s"
  region               = "nyc1"
  size                 = "db-s-2vcpu-4gb"
  tags                 = ["staging"]
  private_network_uuid = digitalocean_vpc.foobar.id
}`
