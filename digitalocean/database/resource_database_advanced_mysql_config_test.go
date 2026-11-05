package database_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanDatabaseAdvancedMySQLConfig_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	dbConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterAdvancedMySQL, name, "8")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseAdvancedMySQLConfigBasic, dbConfig, "UTC", "10"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_database_advanced_mysql_config.foobar", "mysql_parameters.time_zone", "UTC"),
					resource.TestCheckResourceAttr("digitalocean_database_advanced_mysql_config.foobar", "mysql_parameters.connect_timeout", "10"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseAdvancedMySQLConfigBasic, dbConfig, "SYSTEM", "15"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_database_advanced_mysql_config.foobar", "mysql_parameters.time_zone", "SYSTEM"),
					resource.TestCheckResourceAttr("digitalocean_database_advanced_mysql_config.foobar", "mysql_parameters.connect_timeout", "15"),
				),
			},
		},
	})
}

const testAccCheckDigitalOceanDatabaseClusterAdvancedMySQL = `
resource "digitalocean_database_cluster" "foobar" {
  name       = "%s"
  engine     = "advanced_mysql"
  version    = "%s"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
  tags       = ["production"]
}`

const testAccCheckDigitalOceanDatabaseAdvancedMySQLConfigBasic = `
%s

resource "digitalocean_database_advanced_mysql_config" "foobar" {
  cluster_id = digitalocean_database_cluster.foobar.id

  mysql_parameters = {
    time_zone       = "%s"
    connect_timeout = "%s"
  }
}`
