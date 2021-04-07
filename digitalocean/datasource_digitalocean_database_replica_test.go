package digitalocean

import (
	"fmt"
	"testing"
	"time"

	"github.com/digitalocean/godo"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceDigitalOceanDatabaseReplica_Basic(t *testing.T) {
	var databaseReplica godo.DatabaseReplica
	var database godo.Database

	databaseName := randomTestName()
	databaseReplicaName := randomTestName()

	databaseConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterConfigBasic, databaseName)
	replicaConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseReplicaConfigBasic, databaseReplicaName)
	datasourceReplicaConfig := fmt.Sprintf(testAccCheckDigitalOceanDatasourceDatabaseReplicaConfigBasic, databaseReplicaName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseReplicaDestroy,
		Steps: []resource.TestStep{
			{
				Config: databaseConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseClusterExists("digitalocean_database_cluster.foobar", &database),
					resource.TestCheckFunc(
						func(s *terraform.State) error {
							time.Sleep(30 * time.Second)
							return nil
						},
					),
				),
			},
			{
				Config: databaseConfig + replicaConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseReplicaExists("digitalocean_database_replica.read-01", &databaseReplica),
					testAccCheckDigitalOceanDatabaseReplicaAttributes(&databaseReplica, databaseReplicaName),
				),
			},
			{
				Config: databaseConfig + replicaConfig + datasourceReplicaConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("digitalocean_database_replica.read-01", "cluster_id",
						"data.digitalocean_database_replica.my_db_replica", "cluster_id"),
					resource.TestCheckResourceAttrPair("digitalocean_database_replica.read-01", "name",
						"data.digitalocean_database_replica.my_db_replica", "name"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_database_replica.my_db_replica", "region", "nyc3"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_database_replica.my_db_replica", "name", databaseReplicaName),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_database_replica.my_db_replica", "host"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_database_replica.my_db_replica", "private_host"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_database_replica.my_db_replica", "port"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_database_replica.my_db_replica", "user"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_database_replica.my_db_replica", "uri"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_database_replica.my_db_replica", "private_uri"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_database_replica.my_db_replica", "password"),
					resource.TestCheckResourceAttr(
						"data.digitalocean_database_replica.my_db_replica", "tags.#", "1"),
					resource.TestCheckResourceAttrSet(
						"data.digitalocean_database_replica.my_db_replica", "private_network_uuid"),
				),
			},
		},
	})
}

const (
	testAccCheckDigitalOceanDatasourceDatabaseReplicaConfigBasic = `
data "digitalocean_database_replica" "my_db_replica" {
	cluster_id = digitalocean_database_cluster.foobar.id
	name       = "%s"
	}`
)
