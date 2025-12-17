---
page_title: "DigitalOcean: digitalocean_database_logsink_opensearch"
subcategory: "Databases"
---

# digitalocean\_database\_logsink\_opensearch

Provides a DigitalOcean database logsink resource allowing you to forward logs from a managed database cluster to an external OpenSearch cluster or Elasticsearch endpoint.

This resource is compatible with both OpenSearch and Elasticsearch endpoints due to API compatibility. You can use this resource to connect to either service.

This resource supports the following DigitalOcean managed database engines:

* PostgreSQL
* MySQL
* Kafka
* Valkey

**Note**: MongoDB databases use a different log forwarding mechanism and require Datadog logsinks (not currently available in this provider).

## Example Usage

### Basic OpenSearch configuration

```hcl
resource "digitalocean_database_logsink_opensearch" "example" {
  cluster_id     = digitalocean_database_cluster.postgres-example.id
  name           = "opensearch-logs"
  endpoint       = "https://opensearch.example.com:9200"
  index_prefix   = "db-logs"
  index_days_max = 7
}

resource "digitalocean_database_cluster" "postgres-example" {
  name       = "example-postgres-cluster"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}
```

### OpenSearch with authentication and CA certificate

```hcl
resource "digitalocean_database_logsink_opensearch" "example-secure" {
  cluster_id      = digitalocean_database_cluster.postgres-example.id
  name            = "opensearch-secure"
  endpoint        = "https://user:password@opensearch.example.com:9200"
  index_prefix    = "secure-logs"
  index_days_max  = 14
  ca_cert         = file("/path/to/ca.pem")
  timeout_seconds = 30
}
```

### Elasticsearch endpoint configuration

```hcl
resource "digitalocean_database_logsink_opensearch" "elasticsearch" {
  cluster_id     = digitalocean_database_cluster.postgres-example.id
  name           = "elasticsearch-logs"
  endpoint       = "https://elasticsearch.example.com:9243"
  index_prefix   = "es-logs"
  index_days_max = 30
}
```

### MySQL to OpenSearch configuration

```hcl
resource "digitalocean_database_logsink_opensearch" "mysql" {
  cluster_id     = digitalocean_database_cluster.mysql-example.id
  name           = "mysql-logs"
  endpoint       = "https://opensearch.example.com:9200"
  index_prefix   = "mysql-logs"
  index_days_max = 7
}

resource "digitalocean_database_cluster" "mysql-example" {
  name       = "example-mysql-cluster"
  engine     = "mysql"
  version    = "8"
  size       = "db-s-1vcpu-1gb"
  region     = "nyc1"
  node_count = 1
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) UUID of the source database cluster that will forward logs.
* `name` - (Required) Display name for the logsink. **Note**: This is immutable; changing it will force recreation of the resource.
* `endpoint` - (Required) HTTPS URL to the OpenSearch or Elasticsearch cluster (e.g., `https://host:port`). **Note**: Only HTTPS URLs are supported.
* `index_prefix` - (Required) Prefix for the indices where logs will be stored.
* `index_days_max` - (Optional) Maximum number of days to retain indices. Must be 1 or greater.
* `ca_cert` - (Optional) CA certificate for TLS verification in PEM format. Can be specified using `file()` function. This field is marked as sensitive.
* `timeout_seconds` - (Optional) Request timeout for log deliveries in seconds. Must be 1 or greater.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The composite ID of the logsink in the format `cluster_id,logsink_id`.
* `logsink_id` - The unique identifier for the logsink as returned by the DigitalOcean API.

## Import

Database logsink OpenSearch resources can be imported using the composite ID format `cluster_id,logsink_id`. For example:

```
terraform import digitalocean_database_logsink_opensearch.example 245bcfd0-7f31-4ce6-a2bc-475a116cca97,f38db7c8-1f31-4ce6-a2bc-475a116cca97
```

**Note**: The cluster ID and logsink ID must be separated by a comma.

## Important Notes

### Elasticsearch Compatibility
This resource works with both OpenSearch and Elasticsearch endpoints due to their API compatibility. Use the same resource type regardless of whether you're connecting to OpenSearch or Elasticsearch.

### Managed OpenSearch with Trusted Sources
When forwarding logs to a DigitalOcean Managed OpenSearch cluster with trusted sources enabled, you must manually allow-list the IP addresses of your database cluster nodes.

### Authentication
Include authentication credentials directly in the endpoint URL using the format `https://username:password@host:port`. Alternatively, configure authentication on your OpenSearch/Elasticsearch cluster to accept connections from your database cluster's IP addresses.
