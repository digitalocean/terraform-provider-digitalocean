---
page_title: "DigitalOcean: digitalocean_database_logsink_rsyslog"
subcategory: "Databases"
---

# digitalocean\_database\_logsink\_rsyslog

Provides a DigitalOcean database logsink resource allowing you to forward logs from a managed database cluster to an external rsyslog server.

This resource supports the following DigitalOcean managed database engines:

* PostgreSQL
* MySQL
* Kafka
* Valkey

**Note**: MongoDB databases use a different log forwarding mechanism and require Datadog logsinks (not currently available in this provider).

## Example Usage

### Basic rsyslog configuration

```hcl
resource "digitalocean_database_logsink_rsyslog" "example" {
  cluster_id = digitalocean_database_cluster.postgres-example.id
  name       = "rsyslog-prod"
  server     = "192.0.2.10"
  port       = 514
  format     = "rfc5424"
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

### TLS-enabled rsyslog configuration

```hcl
resource "digitalocean_database_logsink_rsyslog" "example-tls" {
  cluster_id = digitalocean_database_cluster.postgres-example.id
  name       = "rsyslog-secure"
  server     = "logs.example.com"
  port       = 6514
  tls        = true
  format     = "rfc5424"
  ca_cert    = file("/path/to/ca.pem")
}
```

### mTLS (mutual TLS) configuration

```hcl
resource "digitalocean_database_logsink_rsyslog" "example-mtls" {
  cluster_id  = digitalocean_database_cluster.postgres-example.id
  name        = "rsyslog-mtls"
  server      = "secure-logs.example.com"
  port        = 6514
  tls         = true
  format      = "rfc5424"
  ca_cert     = file("/path/to/ca.pem")
  client_cert = file("/path/to/client.crt")
  client_key  = file("/path/to/client.key")
}
```

### Custom format configuration

```hcl
resource "digitalocean_database_logsink_rsyslog" "example-custom" {
  cluster_id      = digitalocean_database_cluster.postgres-example.id
  name            = "rsyslog-custom"
  server          = "192.0.2.10"
  port            = 514
  format          = "custom"
  logline         = "<%pri%>%timestamp:::date-rfc3339% %HOSTNAME% %app-name% %msg%"
  structured_data = "[example@41058 iut=\"3\"]"
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) UUID of the source database cluster that will forward logs.
* `name` - (Required) Display name for the logsink. **Note**: This is immutable; changing it will force recreation of the resource.
* `server` - (Required) Hostname or IP address of the rsyslog server.
* `port` - (Required) Port number for the rsyslog server. Must be between 1 and 65535.
* `tls` - (Optional) Enable TLS encryption for the rsyslog connection. Defaults to `false`. **Note**: It is highly recommended to enable TLS as log messages may contain sensitive information.
* `format` - (Optional) Log format to use. Must be one of `rfc5424` (default), `rfc3164`, or `custom`.
* `logline` - (Optional) Custom logline template. **Required** when `format` is set to `custom`. Supports rsyslog-style templating with the following tokens: `%HOSTNAME%`, `%app-name%`, `%msg%`, `%msgid%`, `%pri%`, `%procid%`, `%structured-data%`, `%timestamp%`, and `%timestamp:::date-rfc3339%`.
* `structured_data` - (Optional) Content of the structured data block for RFC5424 messages.
* `ca_cert` - (Optional) CA certificate for TLS verification in PEM format. Can be specified using `file()` function.
* `client_cert` - (Optional) Client certificate for mutual TLS authentication in PEM format. **Note**: Requires `tls` to be `true`.
* `client_key` - (Optional) Client private key for mutual TLS authentication in PEM format. **Note**: Requires `tls` to be `true`. This field is marked as sensitive.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The composite ID of the logsink in the format `cluster_id,logsink_id`.
* `logsink_id` - The unique identifier for the logsink as returned by the DigitalOcean API.

## Import

Database logsink rsyslog resources can be imported using the composite ID format `cluster_id,logsink_id`. For example:

```
terraform import digitalocean_database_logsink_rsyslog.example 245bcfd0-7f31-4ce6-a2bc-475a116cca97,f38db7c8-1f31-4ce6-a2bc-475a116cca97
```

**Note**: The cluster ID and logsink ID must be separated by a comma.
