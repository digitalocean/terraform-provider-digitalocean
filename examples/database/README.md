# DigitalOcean Database Example

This example demonstrates how to create a PostgreSQL database cluster on DigitalOcean and access its metrics credentials. It showcases:

1. Creating a PostgreSQL database cluster with multiple nodes on the regions default VPC
2. Accessing the account-wide database metrics credentials
3. Outputting the metrics endpoints and credentials

## Prerequisites

You need to export your DigitalOcean API Token as an environment variable:

```
export DIGITALOCEAN_TOKEN="Your API TOKEN"
```

## Configuration

The example uses variables that can be customized:

- `region`: DigitalOcean region (default: "sfo3")
- `db_name`: Name for the database cluster (default: "test-pg")
- `db_engine`: Database engine (default: "pg" for PostgreSQL)
- `db_version`: Database version (default: "17")
- `db_size`: Database size slug (default: "db-s-2vcpu-4gb")
- `db_node_count`: Number of nodes (default: 2)

## Run this example using:

```
terraform init
terraform plan
terraform apply
```

## Notes

The database metrics credentials are account-wide, not cluster-specific.
