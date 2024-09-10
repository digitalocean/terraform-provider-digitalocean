---
page_title: "DigitalOcean: digitalocean_database_logsink"
---

# digitalocean\_database\_logsink

Provides a DigitalOcean database logsink capabilities. Can be configured with rsyslog, elasticsearch, and opensearch.

## Example Usage

```hcl
resource "digitalocean_database_logsink" "logsink-01" {
  cluster_id         = digitalocean_database_cluster.doby.id
  sink_name = "sinkexample"
  sink_type = "opensearch"


  config {
    url= "https://user:passwd@192.168.0.1:25060"
    index_prefix= "opensearch-logs"
    index_days_max= 5
  }
}

resource "digitalocean_database_cluster" "doby" {
  name       = "dobydb"
  engine     = "pg"
  version    = "15"
  size       = "db-s-1vcpu-2gb"
  region     = "nyc1"
  node_count = 1
  tags       = ["production"]
}
```


## Argument Reference

The following arguments are supported. See the [DigitalOcean API documentation](https://docs.digitalocean.com/reference/api/api-reference/#operation/databases_create_logsink)
for additional details on each option.

* `cluster_id` - (Required)  The ID of the target MySQL cluster.
* `sink_name` - (Required) The name of the Logsink.
* `sink_type` - (Required) Sink type. Enum: `rsyslog` `elasticsearch` `opensearch`
* `config` - (Required) Logsink configurations.
    - `rsyslog` configuration options:
        - `server` - (Required) The DNS name or IPv4 address of the rsyslog server.
        - `port` - (Required) An integer of the internal port on which the rsyslog server is listening.
        - `tls` - (Required) A boolean to use TLS (as the messages are not filtered and may contain sensitive information, it is highly recommended to set this to true if the remote server supports it).
        - `format` - (Required) A message format used by the server, this can be either rfc3164 (the old BSD style message format), rfc5424 (current syslog message format) or custom. Enum: `rfc5424`, `rfc3164`, or `custom`.
        - `logline` - (Optional) Only required if format == custom. Syslog log line template for a custom format, supporting limited rsyslog style templating (using %tag%). Supported tags are: HOSTNAME, app-name, msg, msgid, pri, procid, structured-data, timestamp and timestamp:::date-rfc3339.
        - `sd` - (Optional) content of the structured data block of rfc5424 message.
        - `ca` - (Optional) PEM encoded CA certificate.
        - `key` - (Optional) (PEM format) client key if the server requires client authentication
        - `cert` - (Optional) (PEM format) client cert to use
    - `elasticsearch` configuration options:
        - `url` - (Required) Elasticsearch connection URL.
        - `index_prefix` - (Required) Elasticsearch index prefix.
        - `index_days_max` - (Optional) Maximum number of days of logs to keep.
        - `timeout` - (Optional) Elasticsearch request timeout limit.
        - `ca` - (Optional) PEM encoded CA certificate.
    - `opensearch` configuration options:
        - `url` - (Required) Opensearch connection URL.
        - `index_prefix` - (Required) Opensearch index prefix.
        - `index_days_max` - (Optional) Maximum number of days of logs to keep.
        - `timeout` - (Optional) Opensearch request timeout limit.
        - `ca` - (Optional) PEM encoded CA certificate.




## Attributes Reference

All above attributes are exported. If an attribute was set outside of Terraform, it will be computed.

## Import

A MySQL database cluster's configuration can be imported using the `id` the parent cluster, e.g.

```
terraform import digitalocean_database_mysql_config.example 4b62829a-9c42-465b-aaa3-84051048e712
```
