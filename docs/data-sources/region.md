---
page_title: "DigitalOcean: digitalocean_region"
---

# digitalocean_region

Get information on a single DigitalOcean region. This is useful to find out 
what Droplet sizes and features are supported within a region.

## Example Usage

```hcl
data "digitalocean_region" "sfo2" {
  slug = "sfo2"
} 

output "region_name" {
  value = data.digitalocean_region.sfo2.name
}
```

## Argument Reference

* `slug` - (Required) A human-readable string that is used as a unique identifier for each region.

## Attributes Reference

* `slug` - A human-readable string that is used as a unique identifier for each region.
* `name` - The display name of the region.
* `available` -	A boolean value that represents whether new Droplets can be created in this region.
* `sizes` - A set of identifying slugs for the Droplet sizes available in this region.
* `features` - A set of features available in this region.
