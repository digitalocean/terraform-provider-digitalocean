package database_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanDatabaseKafkaConfig_Basic(t *testing.T) {
	name := acceptance.RandomTestName()
	dbConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterKafka, name, "3.7")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseKafkaConfigConfigBasic, dbConfig, 3000, true, "9223372036854776000"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_database_kafka_config.foobar", "group_initial_rebalance_delay_ms", "3000"),
					resource.TestCheckResourceAttr("digitalocean_database_kafka_config.foobar", "log_message_downconversion_enable", "true"),
					resource.TestCheckResourceAttr("digitalocean_database_kafka_config.foobar", "log_message_timestamp_difference_max_ms", "9223372036854776000"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseKafkaConfigConfigBasic, dbConfig, 300000, false, "0"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("digitalocean_database_kafka_config.foobar", "group_initial_rebalance_delay_ms", "300000"),
					resource.TestCheckResourceAttr("digitalocean_database_kafka_config.foobar", "log_message_downconversion_enable", "false"),
					resource.TestCheckResourceAttr("digitalocean_database_kafka_config.foobar", "log_message_timestamp_difference_max_ms", "0"),
				),
			},
		},
	})
}

const testAccCheckDigitalOceanDatabaseKafkaConfigConfigBasic = `
%s

resource "digitalocean_database_kafka_config" "foobar" {
  cluster_id                         = digitalocean_database_cluster.foobar.id
  group_initial_rebalance_delay_ms = %d
  log_message_downconversion_enable" = %t
  log_message_timestamp_difference_max_ms = %s
}`
