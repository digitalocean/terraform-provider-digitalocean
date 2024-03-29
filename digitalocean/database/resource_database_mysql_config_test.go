package database_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanDatabaseMySQLConfig_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	dbConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterMySQL, name, "8")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseMySQLConfigConfigBasic, dbConfig, 10, "UTC", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_database_mysql_config.foobar", "connect_timeout", "10"),
					resource.TestCheckResourceAttr("digitalocean_database_mysql_config.foobar", "default_time_zone", "UTC"),
					resource.TestCheckResourceAttr("digitalocean_database_mysql_config.foobar", "sql_require_primary_key", "false"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseMySQLConfigConfigBasic, dbConfig, 15, "SYSTEM", false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_database_mysql_config.foobar", "connect_timeout", "15"),
					resource.TestCheckResourceAttr("digitalocean_database_mysql_config.foobar", "default_time_zone", "SYSTEM"),
					resource.TestCheckResourceAttr("digitalocean_database_mysql_config.foobar", "sql_require_primary_key", "false"),
				),
			},
		},
	})
}

const testAccCheckDigitalOceanDatabaseMySQLConfigConfigBasic = `
%s

resource "digitalocean_database_mysql_config" "foobar" {
  cluster_id              = digitalocean_database_cluster.foobar.id
  connect_timeout         = %d
  default_time_zone       = "%s"
  sql_require_primary_key = "%t"
}`
