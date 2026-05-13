---
page_title: "DigitalOcean: digitalocean_dedicated_inference_gpu_model_config"
subcategory: "Dedicated Inference"
---

# digitalocean\_dedicated\_inference\_gpu\_model\_config

Returns the supported GPU and model compatibility matrix for dedicated inference
endpoints. Use this data source to discover which models can be deployed on which
GPU types.

## Example Usage

```hcl
data "digitalocean_dedicated_inference_gpu_model_config" "available" {}

output "gpu_model_configs" {
  value = data.digitalocean_dedicated_inference_gpu_model_config.available.gpu_model_configs
}
```

## Argument Reference

There are no arguments for this data source.

## Attributes Reference

The following attributes are exported:

* `gpu_model_configs` - The list of supported GPU and model combinations. Each element contains:
  - `gpu_slugs` - The GPU slugs that support this model.
  - `model_slug` - The slug identifier for the model.
  - `model_name` - The human-readable name of the model.
  - `is_model_gated` - Whether the model requires gated access (e.g. a HuggingFace token).
