---
page_title: "DigitalOcean: digitalocean_domain"
subcategory: "Networking"
---

# digitalocean_domain

Get information on a domain. This data source provides the name, TTL, and zone
file as configured on your DigitalOcean account. This is useful if the domain
name in question is not managed by Terraform or you need to utilize TTL or zone
file data.

An error is triggered if the provided domain name is not managed with your
DigitalOcean account.

## Example Usage

Get the zone file for a domain:

```hcl
data "digitalocean_domain" "example" {
  name = "example.com"
}

output "domain_output" {
  value = data.digitalocean_domain.example.zone_file
}
```

```
  $ terraform apply

data.digitalocean_domain.example: Refreshing state...

Apply complete! Resources: 0 added, 0 changed, 0 destroyed.

Outputs:

domain_output = $ORIGIN example.com.
$TTL 1800
example.com. IN SOA ns1.digitalocean.com. hostmaster.example.com. 1516944700 10800 3600 604800 1800
example.com. 1800 IN NS ns1.digitalocean.com.
example.com. 1800 IN NS ns2.digitalocean.com.
example.com. 1800 IN NS ns3.digitalocean.com.
www.example.com. 3600 IN A 176.107.155.137
db.example.com. 3600 IN A 179.189.166.115
jira.example.com. 3600 IN A 207.189.228.15
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the domain.

## Attributes Reference

The following attributes are exported:

* `ttl`: The TTL of the domain.
* `urn` - The uniform resource name of the domain
* `zone_file`: The zone file of the domain.
