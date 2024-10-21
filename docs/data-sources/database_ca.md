---
page_title: "DigitalOcean: digitalocean_database_ca"
subcategory: "Databases"
---

# digitalocean\_database\_ca

Provides the CA certificate for a DigitalOcean database.

## Example Usage

```hcl
data "digitalocean_database_ca" "ca" {
  cluster_id = "aaa-bbb-ccc-ddd"
}

output "ca_output" {
  value = data.digitalocean_database_ca.ca.certificate
}
```

## Argument Reference

The following arguments are supported:

* `cluster_id` - (Required) The ID of the source database cluster.


## Attributes Reference

The following attributes are exported:

* `certificate` - The CA certificate used to secure database connections decoded to a string.
