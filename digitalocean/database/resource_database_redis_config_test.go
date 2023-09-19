package database_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanDatabaseRedisConfig_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	dbConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterRedis, name, "7")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseRedisConfigConfigBasic, dbConfig, "noeviction", 3600, "KA"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"digitalocean_database_redis_config.foobar", "maxmemory_policy", "noeviction"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_redis_config.foobar", "timeout", "3600"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_redis_config.foobar", "notify_keyspace_events", "KA"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_redis_config.foobar", "ssl", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_redis_config.foobar", "persistence", "rdb"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseRedisConfigConfigBasic, dbConfig, "allkeys-lru", 600, "KEA"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"digitalocean_database_redis_config.foobar", "maxmemory_policy", "allkeys-lru"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_redis_config.foobar", "timeout", "600"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_redis_config.foobar", "notify_keyspace_events", "KEA"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_redis_config.foobar", "ssl", "true"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_redis_config.foobar", "persistence", "rdb"),
				),
			},
		},
	})
}

const testAccCheckDigitalOceanDatabaseRedisConfigConfigBasic = `
%s

resource "digitalocean_database_redis_config" "foobar" {
  cluster_id             = digitalocean_database_cluster.foobar.id
  maxmemory_policy       = "%s"
  timeout                = %d
  notify_keyspace_events = "%s"
}`
