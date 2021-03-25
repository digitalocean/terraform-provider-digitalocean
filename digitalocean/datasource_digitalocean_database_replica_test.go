package digitalocean

import (
	"fmt"
	"testing"

	"github.com/digitalocean/godo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanDatabaseReplica_Basic(t *testing.T) {
	var databaseReplica godo.DatabaseReplica
	var database godo.Database

	databaseName := randomTestName()
	databaseReplicaName := randomTestName()

	databaseConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigBasic, databaseName)
	replicaConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseReplicaConfigBasic, databaseReplicaName)
	datasourceReplicaConfig := fmt.Sprintf(testAccCheckDigitalOceanDatasourceDatabaseReplicaConfigBasic, "digitalocean_database_replica.read-01.cluster_id", databaseReplicaName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
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
				),
			},
			{
				Config: databaseConfig + replicaConfig + datasourceReplicaConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("digitalocean_database_replica.read-01", "cluster_id",
						"data.digitalocean_database_replica.my_db_replica", "cluster_id"),
					resource.TestCheckResourceAttrPair("digitalocean_database_replica.read-01", "name",
						"data.digitalocean_database_replica.my_db_replica", "name"),
				),
			},
		},
	})
}

const testAccCheckDigitalOceanDatasourceDatabaseReplicaConfigBasic = `
data "digitalocean_database_replica" "my_db_replica" {
	cluster_id = "%s"
	name       = "%s"
  }`
