package database_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDigitalOceanDatabaseConnectionPool_Basic(t *testing.T) {
	var pool godo.DatabasePool

	databaseName := acceptance.RandomTestName()
	poolName := acceptance.RandomTestName()

	resourceConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseConnectionPoolConfigBasic, databaseName, poolName)
	datasourceConfig := fmt.Sprintf(testAccCheckDigitalOceanDatasourceDatabaseConnectionPoolConfigBasic, poolName)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseConnectionPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseConnectionPoolExists("digitalocean_database_connection_pool.pool-01", &pool),
					testAccCheckDigitalOceanDatabaseConnectionPoolAttributes(&pool, poolName),
					resource.TestCheckResourceAttr(
						"digitalocean_database_connection_pool.pool-01", "name", poolName),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_connection_pool.pool-01", "cluster_id"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_connection_pool.pool-01", "size", "10"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_connection_pool.pool-01", "mode", "transaction"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_connection_pool.pool-01", "db_name", "defaultdb"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_connection_pool.pool-01", "user", "doadmin"),
				),
			},
			{
				Config: resourceConfig + datasourceConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair("digitalocean_database_connection_pool.pool-01", "name",
						"data.digitalocean_database_connection_pool.pool-01", "name"),
					resource.TestCheckResourceAttrPair("digitalocean_database_connection_pool.pool-01", "mode",
						"data.digitalocean_database_connection_pool.pool-01", "mode"),
					resource.TestCheckResourceAttrPair("digitalocean_database_connection_pool.pool-01", "size",
						"data.digitalocean_database_connection_pool.pool-01", "size"),
					resource.TestCheckResourceAttrPair("digitalocean_database_connection_pool.pool-01", "db_name",
						"data.digitalocean_database_connection_pool.pool-01", "db_name"),
					resource.TestCheckResourceAttrPair("digitalocean_database_connection_pool.pool-01", "user",
						"data.digitalocean_database_connection_pool.pool-01", "user"),
				),
			},
		},
	})
}

const testAccCheckDigitalOceanDatasourceDatabaseConnectionPoolConfigBasic = `
data "digitalocean_database_connection_pool" "pool-01" {
  cluster_id = digitalocean_database_cluster.foobar.id
  name       = "%s"
}`
