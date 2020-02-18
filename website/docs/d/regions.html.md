---
layout: "digitalocean"
page_title: "DigitalOcean: digitalocean_regions"
sidebar_current: "docs-do-datasource-regions"
description: |-
  Get identifiers for all or a filtered subset of DigitalOcean regions.
---

# digitalocean_regions

Get identifiers for all or a filtered subset of DigitalOcean regions. The regions
can be filtered using the `available` and/or `features` attributes. If no filters
are specified, all regions will be returned.

Note: Use the `digitalocean_region` data source to obtain metadata about a single
region.

## Example Usage

```hcl
data "digitalocean_regions" "available" {
  available = true
} 

output "available_regions" {
  value = data.digitalocean_regions.slugs
}
```

## Argument Reference

* `available` - (Optional) A boolean that filters the list of regions by whether or not new Droplets can be created
* `features` - (Optional) A list of features required to be supported by matching regions  

## Attributes Reference

* `slugs` - A list of human-readable strings that are the unique identifier for each matching region. 
