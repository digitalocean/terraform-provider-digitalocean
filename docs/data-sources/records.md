---
page_title: "DigitalOcean: digitalocean_records"
---

# digitalocean_records

Retrieve information about all DNS records within a domain, with the ability to filter and sort the results.
If no filters are specified, all records will be returned.

## Example Usage

Get data for all MX records in a domain:

```hcl
data "digitalocean_records" "example" {
  domain = "example.com"
  filter {
    key    = "type"
    values = ["MX"]
  }
}

output "mail_servers" {
  value = join(",", data.digitalocean_records.example.records[*].value)
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Required) The domain name to search for DNS records

* `filter` - (Optional) Filter the results.
  The `filter` block is documented below.

* `sort` - (Optional) Sort the results.
  The `sort` block is documented below.

`filter` supports the following arguments:

* `key` - (Required) Filter the DNS records by this key. This may be one of `domain`, `flags`, `name`, `port`,
  `priority`, `tag`, `ttl`, `type`, `value`, or `weight`.
  
* `values` - (Required) A list of values to match against the `key` field. Only retrieves DNS records
  where the `key` field takes on one or more of the values provided here.

* `match_by` - (Optional) One of `exact` (default), `re`, or `substring`. For string-typed fields, specify `re` to
  match by using the `values` as regular expressions, or specify `substring` to match by treating the `values` as
  substrings to find within the string field.
  
* `all` - (Optional) Set to `true` to require that a field match all of the `values` instead of just one or more of
  them. This is useful when matching against multi-valued fields such as lists or sets where you want to ensure
  that all of the `values` are present in the list or set.

`sort` supports the following arguments:

* `key` - (Required) Sort the DNS records by this key. This may be one of `domain`, `flags`, `name`, `port`,
  `priority`, `tag`, `ttl`, `type`, `value`, or `weight`.
* `direction` - (Required) The sort direction. This may be either `asc` or `desc`.

## Attributes Reference

The following attributes are exported:

* `id`: The ID of the record.
* `domain`: Domain of the DNS record.
* `name`: The name of the DNS record.
* `type`:	The type of the DNS record.
* `value`:	Variable data depending on record type. For example, the "data" value for an A record would be the IPv4 address to which the domain will be mapped. For a CAA record, it would contain the domain name of the CA being granted permission to issue certificates.
* `priority`:	The priority for SRV and MX records.
* `port`:	The port for SRV records.
* `ttl`: This value is the time to live for the record, in seconds. This defines the time frame that clients can cache queried information before a refresh should be requested.
* `weight`:	The weight for SRV records.
* `flags`: An unsigned integer between 0-255 used for CAA records.
* `tag`: The parameter tag for CAA records.
