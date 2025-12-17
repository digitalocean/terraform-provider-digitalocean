# DigitalOcean Database Example

This example demonstrates how to create a PostgreSQL database cluster on DigitalOcean with log forwarding configured. It showcases:

1. Creating a PostgreSQL database cluster with multiple nodes on the regions default VPC
2. Configuring an rsyslog logsink to forward database logs to an external rsyslog server
3. Accessing the account-wide database metrics credentials
4. Outputting the metrics endpoints, credentials, and logsink information

## Prerequisites

You need to export your DigitalOcean API Token as an environment variable:

```
export DIGITALOCEAN_TOKEN="Your API TOKEN"
```

## Configuration

The example uses variables that can be customized:

### Database Configuration
- `region`: DigitalOcean region (default: "sfo3")
- `db_name`: Name for the database cluster (default: "test-pg")
- `db_engine`: Database engine (default: "pg" for PostgreSQL)
- `db_version`: Database version (default: "17")
- `db_size`: Database size slug (default: "db-s-2vcpu-4gb")
- `db_node_count`: Number of nodes (default: 2)

### Rsyslog Logsink Configuration
- `rsyslog_server`: Hostname or IP address of your rsyslog server (default: "logs.example.com")
- `rsyslog_port`: Port number for the rsyslog server (default: 514)
- `rsyslog_format`: Log format - rfc5424, rfc3164, or custom (default: "rfc5424")

## Run this example using:

```
terraform init
terraform plan
terraform apply
```

## Notes

- The database metrics credentials are account-wide, not cluster-specific.
- The rsyslog logsink requires a reachable rsyslog server. Update the `rsyslog_server` variable to point to your actual rsyslog server hostname or IP address.
- For production use, consider enabling TLS for the rsyslog connection (see the [rsyslog logsink documentation](../../docs/resources/database_logsink_rsyslog.md) for TLS configuration examples).
