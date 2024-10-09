package database_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanDatabaseMongoDBConfig_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	dbConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterMongoDB, name, "7")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseMongoDBConfigBasic, dbConfig, "available", 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_database_mongodb_config.foobar", "default_read_concern", "available"),
					resource.TestCheckResourceAttr("digitalocean_database_mongodb_config.foobar", "transaction_lifetime_limit_seconds", "1"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseMongoDBConfigBasic, dbConfig, "majority", 100),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_database_mongodb_config.foobar", "default_read_concern", "majority"),
					resource.TestCheckResourceAttr("digitalocean_database_mongodb_config.foobar", "transaction_lifetime_limit_seconds", "100"),
				),
			},
		},
	})
}

const testAccCheckDigitalOceanDatabaseMongoDBConfigBasic = `
%s

resource "digitalocean_database_mongodb_config" "foobar" {
  cluster_id                         = digitalocean_database_cluster.foobar.id
  default_read_concern               = "%s"
  transaction_lifetime_limit_seconds = %d
}`
