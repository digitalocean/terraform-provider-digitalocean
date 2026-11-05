package database_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanDatabaseAdvancedPostgreSQLConfig_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	dbConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterAdvancedPostgreSQL, name, "16")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseAdvancedPostgreSQLConfigBasic, dbConfig, "UTC", "32"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_database_advanced_postgresql_config.foobar", "pg_parameters.timezone", "UTC"),
					resource.TestCheckResourceAttr("digitalocean_database_advanced_postgresql_config.foobar", "pg_parameters.work_mem", "32"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseAdvancedPostgreSQLConfigBasic, dbConfig, "UTC", "16"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_database_advanced_postgresql_config.foobar", "pg_parameters.timezone", "UTC"),
					resource.TestCheckResourceAttr("digitalocean_database_advanced_postgresql_config.foobar", "pg_parameters.work_mem", "16"),
				),
			},
		},
	})
}

const testAccCheckDigitalOceanDatabaseClusterAdvancedPostgreSQL = `
resource "digitalocean_database_cluster" "foobar" {
  name       = "%s"
  engine     = "advanced_pg"
  version    = "%s"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
  tags       = ["production"]
}`

const testAccCheckDigitalOceanDatabaseAdvancedPostgreSQLConfigBasic = `
%s

resource "digitalocean_database_advanced_postgresql_config" "foobar" {
  cluster_id = digitalocean_database_cluster.foobar.id

  pg_parameters = {
    timezone = "%s"
    work_mem = "%s"
  }
}`
