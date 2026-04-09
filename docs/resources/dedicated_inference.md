---
page_title: "DigitalOcean: digitalocean_dedicated_inference"
subcategory: "Dedicated Inference"
---

# digitalocean\_dedicated\_inference

Provides a DigitalOcean Dedicated Inference resource. This can be used to create,
modify, and delete dedicated inference endpoints for running GPU-accelerated
model inference.

## Example Usage

```hcl
resource "digitalocean_dedicated_inference" "example" {
  name   = "my-inference-endpoint"
  region = "tor1"

  model_deployments {
    model_slug     = "deepseek-r1-distill-qwen-14b"
    model_provider = "digitalocean"

    accelerators {
      accelerator_slug = "gpu-h100x1-80gb"
      scale            = 1
      type             = "nvidia_h100"
    }
  }
}
```

### With Public Endpoint

```hcl
resource "digitalocean_dedicated_inference" "public" {
  name                   = "my-public-inference"
  region                 = "tor1"
  enable_public_endpoint = true

  model_deployments {
    model_slug     = "deepseek-r1-distill-qwen-14b"
    model_provider = "digitalocean"

    accelerators {
      accelerator_slug = "gpu-h100x1-80gb"
      scale            = 1
      type             = "nvidia_h100"
    }
  }
}
```

### With VPC

```hcl
resource "digitalocean_dedicated_inference" "private" {
  name     = "my-private-inference"
  region   = "tor1"
  vpc_uuid = digitalocean_vpc.example.id

  model_deployments {
    model_slug     = "deepseek-r1-distill-qwen-14b"
    model_provider = "digitalocean"

    accelerators {
      accelerator_slug = "gpu-h100x1-80gb"
      scale            = 1
      type             = "nvidia_h100"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A human-readable name for the dedicated inference endpoint.
* `region` - (Required) The region slug where the dedicated inference endpoint will be deployed. Changing this forces a new resource.
* `enable_public_endpoint` - (Optional) Whether to enable a public HTTPS endpoint for the dedicated inference endpoint. Defaults to `false`. This field is immutable after creation and changing it forces a new resource.
* `vpc_uuid` - (Optional) The UUID of the VPC to deploy the dedicated inference endpoint into. Changing this forces a new resource.
* `model_deployments` - (Required) The list of model deployments to run on the dedicated inference endpoint. Each `model_deployments` block supports:
  - `model_slug` - (Required) The slug identifier for the model to deploy.
  - `model_provider` - (Required) The provider of the model (e.g. `digitalocean`, `huggingface`).
  - `model_id` - (Optional) The unique ID of the model.
  - `accelerators` - (Required) The GPU accelerators to allocate for this model deployment. Each `accelerators` block supports:
    - `accelerator_slug` - (Required) The slug identifier for the GPU accelerator type.
    - `scale` - (Required) The number of accelerator units to allocate. Must be at least 1.
    - `type` - (Required) The accelerator type.
* `hugging_face_token` - (Optional, Sensitive) A HuggingFace token for accessing gated models.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - The unique ID of the dedicated inference endpoint.
* `status` - The current status of the dedicated inference endpoint.
* `public_endpoint_fqdn` - The fully-qualified domain name of the public endpoint, if enabled.
* `private_endpoint_fqdn` - The fully-qualified domain name of the private endpoint.
* `created_at` - The date and time when the dedicated inference endpoint was created.
* `updated_at` - The date and time when the dedicated inference endpoint was last updated.

## Import

Dedicated inference endpoints can be imported using their `id`, e.g.

```
terraform import digitalocean_dedicated_inference.example endpoint-id
```

## Timeouts

`digitalocean_dedicated_inference` provides the following
[Timeouts](https://www.terraform.io/docs/language/resources/syntax.html#operation-timeouts)
configuration options:

- `create` - (Default `60 minutes`) Used for creating the dedicated inference endpoint.
- `update` - (Default `60 minutes`) Used for updating the dedicated inference endpoint.
- `delete` - (Default `60 minutes`) Used for deleting the dedicated inference endpoint.
