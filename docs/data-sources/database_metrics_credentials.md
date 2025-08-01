---
page_title: "DigitalOcean: digitalocean_database_metrics_credentials"
subcategory: "Databases"
---

# digitalocean_database_metrics_credentials

Provides access to the metrics credentials for DigitalOcean database clusters. These credentials are account-wide and can be used to access metrics for any database cluster in the account.

## Example Usage

```hcl
data "digitalocean_database_metrics_credentials" "example" {}

output "metrics_username" {
  value = data.digitalocean_database_metrics_credentials.example.username
}

output "metrics_password" {
  sensitive = true
  value     = data.digitalocean_database_metrics_credentials.example.password
}
```

## Argument Reference

This datasource doesn't require any arguments.

## Attributes Reference

The following attributes are exported:

* `username` - The username for accessing database metrics.
* `password` - The password for accessing database metrics. This is marked as sensitive.
