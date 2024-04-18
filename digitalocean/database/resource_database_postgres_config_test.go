package database_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanDatabasePostgreSQLConfig_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	dbConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterPostgreSQL, name, "15")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabasePostgreSQLConfigConfigBasic, dbConfig, "UTC", 30.5, 32, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_database_postgresql_config.foobar", "jit", "false"),
					resource.TestCheckResourceAttr("digitalocean_database_postgresql_config.foobar", "timezone", "UTC"),
					resource.TestCheckResourceAttr("digitalocean_database_postgresql_config.foobar", "shared_buffers_percentage", "30.5"),
					resource.TestCheckResourceAttr("digitalocean_database_postgresql_config.foobar", "work_mem", "32"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabasePostgreSQLConfigConfigBasic, dbConfig, "UTC", 20.0, 16, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_database_postgresql_config.foobar", "jit", "true"),
					resource.TestCheckResourceAttr("digitalocean_database_postgresql_config.foobar", "timezone", "UTC"),
					resource.TestCheckResourceAttr("digitalocean_database_postgresql_config.foobar", "shared_buffers_percentage", "20"),
					resource.TestCheckResourceAttr("digitalocean_database_postgresql_config.foobar", "work_mem", "16"),
				),
			},
		},
	})
}

const testAccCheckDigitalOceanDatabasePostgreSQLConfigConfigBasic = `
%s

resource "digitalocean_database_postgresql_config" "foobar" {
  cluster_id                = digitalocean_database_cluster.foobar.id
  timezone                  = "%s"
  shared_buffers_percentage = %f
  work_mem                  = %d
  jit                       = %t
}`
