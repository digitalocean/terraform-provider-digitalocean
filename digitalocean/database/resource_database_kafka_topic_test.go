package database_test

import (
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDigitalOceanDatabaseKafkaTopic(t *testing.T) {
	name := acceptance.RandomTestName()
	dbConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterKafka, name, "3.5")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseClusterDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseKafkaTopicBasic, dbConfig, "topic-foobar"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"digitalocean_database_kafka_topic.foobar", "name", "topic-foobar"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_kafka_topic.foobar", "state", "active"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_kafka_topic.foobar", "replication_factor", "2"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_kafka_topic.foobar", "partition_count", "3"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseKafkaTopicComplete, dbConfig, "topic-foobar", 5, 3, "compact", "snappy", 80000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"digitalocean_database_kafka_topic.foobaz", "name", "topic-foobar"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_kafka_topic.foobaz", "state", "active"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_kafka_topic.foobaz", "replication_factor", "3"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_kafka_topic.foobaz", "partition_count", "5"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_kafka_topic.foobaz.config", "cleanup_policy", "compact"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_kafka_topic.foobaz.config", "compression_type", "snappy"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_kafka_topic.foobaz.config", "delete_retention_ms", "80000"),
				),
			},
		},
	})
}

const testAccCheckDigitalOceanDatabaseKafkaTopicBasic = `
%s

resource "digitalocean_database_kafka_topic" "foobar" {
  cluster_id         = digitalocean_database_cluster.foobar.id
  name       		 = "%s"
}`

const testAccCheckDigitalOceanDatabaseKafkaTopicComplete = `
%s

resource "digitalocean_database_kafka_topic" "foobaz" {
  cluster_id         = digitalocean_database_cluster.foobar.id
  name       		 = "%s"
  partition_count    = %d
  replication_factor = %d
  config = {
	cleanup_policy = "%s"
	compression_type = "%s"
	delete_retention_ms = %d
  }
}`
