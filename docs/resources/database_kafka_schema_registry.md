---
page_title: "DigitalOcean: digitalocean_database_kafka_schema_registry"
subcategory: "Databases"
---

# digitalocean\_database\_kafka\_schema\_registry

Provides a DigitalOcean Kafka schema registry for Kafka clusters.

## Example Usage

### Create a new Kafka Schema Registry
```hcl
resource "digitalocean_database_kafka_schema_registry" "schema-01" {
  cluster_id   = digitalocean_database_cluster.kafka-example.id
  subject_name = "test-schema"
  schema_type  = "avro"
  schema       = <<EOF
{
  "type": "record",
  "namespace": "example",
  "name": "TestRecord",
  "fields": [
    {"name": "id", "type": "string"},
    {"name": "name", "type": "string"},
    {"name": "value", "type": "int"}
  ]
}
EOF
}

resource "digitalocean_database_cluster" "kafka-example" {
  name       = "example-kafka-cluster"
  engine     = "kafka"
  version    = "3.5"
  size       = "gd-2vcpu-8gb"
  region     = "blr1"
  node_count = 3
  tags       = ["production"]
}
```

## Argument Reference

The following arguments are supported:
* `cluster_id` - (Required) The ID of the target Kafka cluster.
* `subject_name` - (Required) The name of the schema subject.
* `schema_type` - (Required) The schema type. Available values are: avro, json, or protobuf.
* `schema` - (Required) The schema definition as a string.