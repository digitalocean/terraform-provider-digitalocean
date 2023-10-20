package database_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanDatabaseKafkaTopic(t *testing.T) {
	name := acceptance.RandomTestName()
	dbConfig := fmt.Sprintf(testAccCheckDigitalOceanDatabaseClusterKafka, name, "3.5")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseKafkaTopicDestroy,
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
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.cleanup_policy"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.compression_type"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.delete_retention_ms"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.file_delete_delay_ms"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.flush_messages"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.flush_ms"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.index_interval_bytes"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.max_compaction_lag_ms"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.message_down_conversion_enable"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.message_format_version"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.message_timestamp_difference_max_ms"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.message_timestamp_type"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.min_cleanable_dirty_ratio"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.min_compaction_lag_ms"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.min_insync_replicas"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.retention_bytes"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.retention_ms"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.segment_bytes"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.segment_index_bytes"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.segment_jitter_ms"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.segment_ms"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.unclean_leader_election_enable"),
				),
			},
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseKafkaTopicWithConfig, dbConfig, "topic-foobar", 5, 3, "compact", "snappy", 80000),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"digitalocean_database_kafka_topic.foobar", "name", "topic-foobar"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_kafka_topic.foobar", "state", "active"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_kafka_topic.foobar", "replication_factor", "3"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_kafka_topic.foobar", "partition_count", "5"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_kafka_topic.foobar", "config.0.cleanup_policy", "compact"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_kafka_topic.foobar", "config.0.compression_type", "snappy"),
					resource.TestCheckResourceAttr(
						"digitalocean_database_kafka_topic.foobar", "config.0.delete_retention_ms", "80000"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.cleanup_policy"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.compression_type"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.delete_retention_ms"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.file_delete_delay_ms"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.flush_messages"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.flush_ms"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.index_interval_bytes"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.max_compaction_lag_ms"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.message_down_conversion_enable"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.message_format_version"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.message_timestamp_difference_max_ms"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.message_timestamp_type"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.min_cleanable_dirty_ratio"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.min_compaction_lag_ms"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.min_insync_replicas"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.retention_bytes"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.retention_ms"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.segment_bytes"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.segment_index_bytes"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.segment_jitter_ms"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.segment_ms"),
					resource.TestCheckResourceAttrSet(
						"digitalocean_database_kafka_topic.foobar", "config.0.unclean_leader_election_enable"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanDatabaseKafkaTopicDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_database_kafka_topic" {
			continue
		}
		clusterId := rs.Primary.Attributes["cluster_id"]
		name := rs.Primary.Attributes["name"]
		// Try to find the kafka topic
		_, _, err := client.Databases.GetTopic(context.Background(), clusterId, name)

		if err == nil {
			return fmt.Errorf("kafka topic still exists")
		}
	}

	return nil
}

const testAccCheckDigitalOceanDatabaseKafkaTopicBasic = `
%s

resource "digitalocean_database_kafka_topic" "foobar" {
  cluster_id = digitalocean_database_cluster.foobar.id
  name       = "%s"
}`

const testAccCheckDigitalOceanDatabaseKafkaTopicWithConfig = `
%s

resource "digitalocean_database_kafka_topic" "foobar" {
  cluster_id         = digitalocean_database_cluster.foobar.id
  name               = "%s"
  partition_count    = %d
  replication_factor = %d
  config {
    cleanup_policy      = "%s"
    compression_type    = "%s"
    delete_retention_ms = %d
  }
}`
