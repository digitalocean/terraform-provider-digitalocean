package database_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanDatabaseValkeyConfig_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	dbConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterValkey, name, "8")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseValkeyConfigConfigBasic, dbConfig, 3600, "KA"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"digitalocean_database_valkey_config.foobar", "timeout", "3600"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_valkey_config.foobar", "notify_keyspace_events", "KA"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_valkey_config.foobar", "ssl", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_valkey_config.foobar", "persistence", "rdb"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseValkeyConfigConfigBasic, dbConfig, 0, "KEA"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"digitalocean_database_valkey_config.foobar", "timeout", "0"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_valkey_config.foobar", "notify_keyspace_events", "KEA"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_valkey_config.foobar", "ssl", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_valkey_config.foobar", "persistence", "rdb"),
				),
			},
		},
	})
}

const testAccCheckDigitalOceanDatabaseValkeyConfigConfigBasic = `
%s

resource "digitalocean_database_valkey_config" "foobar" {
  cluster_id             = digitalocean_database_cluster.foobar.id
  timeout                = %d
  notify_keyspace_events = "%s"
}`
