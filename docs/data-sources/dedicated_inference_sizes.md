---
page_title: "DigitalOcean: digitalocean_dedicated_inference_sizes"
subcategory: "Dedicated Inference"
---

# digitalocean\_dedicated\_inference\_sizes

Returns the available GPU sizes and their configurations for dedicated inference
endpoints, including pricing, hardware specifications, and region availability.

## Example Usage

```hcl
data "digitalocean_dedicated_inference_sizes" "available" {}

output "enabled_regions" {
  value = data.digitalocean_dedicated_inference_sizes.available.enabled_regions
}

output "sizes" {
  value = data.digitalocean_dedicated_inference_sizes.available.sizes
}
```

## Argument Reference

There are no arguments for this data source.

## Attributes Reference

The following attributes are exported:

* `enabled_regions` - The list of region slugs where dedicated inference endpoints can be deployed.
* `sizes` - The list of available GPU sizes. Each element contains:
  - `gpu_slug` - The slug identifier for this GPU size.
  - `price_per_hour` - The hourly price for this GPU size.
  - `currency` - The currency for the price.
  - `regions` - The regions where this GPU size is available.
  - `cpu` - The number of vCPUs.
  - `memory` - The amount of memory in MiB.
  - `gpu` - GPU hardware details. Each element contains:
    - `count` - The number of GPUs.
    - `vram_gb` - The VRAM per GPU in GiB.
    - `slug` - The GPU model slug.
  - `size_category` - The category this size belongs to. Each element contains:
    - `name` - The display name of the size category.
    - `fleet_name` - The fleet name associated with the size category.
  - `disks` - The disks attached to this size. Each element contains:
    - `type` - The disk type.
    - `size_gb` - The disk size in GiB.
