# Terraform configuration block specifying required versions and providers
terraform {
  required_version = "~> 1"
  required_providers {
    digitalocean = {
      source  = "digitalocean/digitalocean"
      version = "~> 2"
    }
  }
}

# Configure the DigitalOcean provider
# Authentication is handled via the DIGITALOCEAN_TOKEN environment variable
provider "digitalocean" {
  # You need to set this in your .bashrc
  # export DIGITALOCEAN_TOKEN="Your API TOKEN"
}

# Look up the default VPC for the specified region
# Database clusters are automatically placed in a VPC for network isolation
data "digitalocean_vpc" "default" {
  region = var.region
}

# Retrieve the account-wide database metrics credentials
# These credentials allow you to access Prometheus metrics endpoints for monitoring
# database cluster performance and health
data "digitalocean_database_metrics_credentials" "test" {}

# Create a managed database cluster
# This example creates a PostgreSQL cluster with multiple nodes for high availability
# The cluster is connected to the region's default VPC for secure networking
resource "digitalocean_database_cluster" "test" {
  name                 = var.db_name
  engine               = var.db_engine
  version              = var.db_version
  size                 = var.db_size
  region               = var.region
  node_count           = var.db_node_count
  private_network_uuid = data.digitalocean_vpc.default.id
}

# Configure log forwarding from the database cluster to an external rsyslog server.
# This allows you to aggregate and analyze database logs in your logging infrastructure.
#
# IMPORTANT: This example uses a basic unencrypted rsyslog configuration for simplicity.
# For production environments, you should ALWAYS enable TLS to protect sensitive log data
# in transit. Database logs may contain query patterns, connection details, and other
# sensitive operational information.
#
# For TLS configuration examples, see:
# https://registry.terraform.io/providers/digitalocean/digitalocean/latest/docs/resources/database_logsink_rsyslog
#
# TLS example configuration:
#   tls         = true
#   ca_cert     = file("/path/to/ca.pem")
#   client_cert = file("/path/to/client.crt")  # Optional, for mutual TLS
#   client_key  = file("/path/to/client.key")  # Optional, for mutual TLS
resource "digitalocean_database_logsink_rsyslog" "example" {
  cluster_id = digitalocean_database_cluster.test.id
  name       = "${var.db_name}-logs"
  server     = var.rsyslog_server
  port       = var.rsyslog_port
  format     = var.rsyslog_format
}
