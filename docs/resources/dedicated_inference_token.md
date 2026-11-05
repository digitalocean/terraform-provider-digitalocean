---
page_title: "DigitalOcean: digitalocean_dedicated_inference_token"
subcategory: "Dedicated Inference"
---

# digitalocean\_dedicated\_inference\_token

Provides a DigitalOcean Dedicated Inference Token resource. This can be used to
create and revoke API tokens for dedicated inference endpoints.

~> **Note:** The `token` attribute is only available immediately after creation
and cannot be retrieved afterwards. Make sure to store it securely.

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

resource "digitalocean_dedicated_inference_token" "example" {
  dedicated_inference_id = digitalocean_dedicated_inference.example.id
  name                   = "my-api-token"
}
```

## Argument Reference

The following arguments are supported:

* `dedicated_inference_id` - (Required) The ID of the dedicated inference endpoint this token belongs to. Changing this forces a new resource.
* `name` - (Required) A human-readable name for the token. Changing this forces a new resource.

## Attributes Reference

In addition to the arguments listed above, the following attributes are exported:

* `id` - The composite ID of the token in the format `{dedicated_inference_id}:{token_id}`.
* `token` - (Sensitive) The token value. Only available immediately after creation and not retrievable afterwards.
* `created_at` - The date and time when the token was created.

## Import

Dedicated inference tokens can be imported using the composite ID
`{dedicated_inference_id}:{token_id}`, e.g.

```
terraform import digitalocean_dedicated_inference_token.example endpoint-id:token-id
```
