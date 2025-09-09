package database_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/digitalocean/godo"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/acceptance"
	"github.com/digitalocean/terraform-provider-digitalocean/digitalocean/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDigitalOceanDatabaseKafkaSchemaRegistry_Basic(t *testing.T) {
	var databaseKafkaSchemaRegistry godo.DatabaseKafkaSchemaRegistryRequest
	databaseName := acceptance.RandomTestName()
	databaseClusterName := acceptance.RandomTestName()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { acceptance.TestAccPreCheck(t) },
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      testAccCheckDigitalOceanDatabaseKafkaSchemaRegistryDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckDigitalOceanDatabaseKafkaSchemaRegistry, databaseClusterName, databaseName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDigitalOceanDatabaseKafkaSchemaRegistryExists("digitalocean_database_kafka_schema_registry.foobar", &databaseKafkaSchemaRegistry),
					resource.TestCheckResourceAttr(
						"digitalocean_database_kafka_schema_registry.foobar", "name", databaseName),
					resource.TestCheckResourceAttrPair(
						"digitalocean_database_kafka_schema_registry.foobar", "cluster_id",
						"digitalocean_database_cluster.kafka", "id"),
				),
			},
		},
	})
}

func testAccCheckDigitalOceanDatabaseKafkaSchemaRegistryDestroy(s *terraform.State) error {
	client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "digitalocean_database_kafka_schema_registry" {
			continue
		}
		clusterID := rs.Primary.Attributes["cluster_id"]
		name := rs.Primary.Attributes["name"]

		// Try to find the registry
		_, _, err := client.Databases.GetKafkaSchemaRegistry(context.Background(), clusterID, name)

		if err == nil {
			return fmt.Errorf("DatabaseKafkaSchemaRegistry still exists")
		}
	}

	return nil
}

func testAccCheckDigitalOceanDatabaseKafkaSchemaRegistryExists(n string, registry *godo.DatabaseKafkaSchemaRegistryRequest) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No DatabaseKafkaSchemaRegistry ID is set")
		}

		client := acceptance.TestAccProvider.Meta().(*config.CombinedConfig).GodoClient()
		clusterID := rs.Primary.Attributes["cluster_id"]
		name := rs.Primary.Attributes["name"]

		foundRegistry, _, err := client.Databases.GetKafkaSchemaRegistry(context.Background(), clusterID, name)
		if err != nil {
			return err
		}

		*registry = godo.DatabaseKafkaSchemaRegistryRequest{
			SubjectName: foundRegistry.SubjectName,
		}

		return nil
	}
}

const testAccCheckDigitalOceanDatabaseKafkaSchemaRegistry = `
resource "digitalocean_database_cluster" "kafka" {
  name       = "%s"
  engine     = "kafka"
  version    = "3.5"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 3
}

resource "digitalocean_database_kafka_schema_registry" "foobar" {
  cluster_id   = digitalocean_database_cluster.kafka.id
  subject_name = "%s"
  schema_type  = "avro"
  schema       = <<EOF
{
  "type": "record",
  "name": "example",
  "fields": [
    {
      "name": "id",
      "type": "int"
    },
    {
      "name": "name",
      "type": "string"
    }
  ]
}
EOF
}
`
