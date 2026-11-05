---
page_title: "DigitalOcean: digitalocean_dedicated_inference"
subcategory: "Dedicated Inference"
---

# digitalocean\_dedicated\_inference

Get information on a dedicated inference endpoint for use in other resources. This
data source provides all of the endpoint's properties as configured on your
DigitalOcean account.

## Example Usage

```hcl
data "digitalocean_dedicated_inference" "example" {
  id = "endpoint-id"
}

output "endpoint_status" {
  value = data.digitalocean_dedicated_inference.example.status
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the dedicated inference endpoint.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the dedicated inference endpoint.
* `region` - The region where the dedicated inference endpoint is deployed.
* `status` - The current status of the dedicated inference endpoint.
* `vpc_uuid` - The UUID of the VPC the dedicated inference endpoint is deployed in.
* `enable_public_endpoint` - Whether the public HTTPS endpoint is enabled.
* `public_endpoint_fqdn` - The fully-qualified domain name of the public endpoint, if enabled.
* `private_endpoint_fqdn` - The fully-qualified domain name of the private endpoint.
* `model_deployments` - The list of model deployments running on the endpoint. Each element contains:
  - `model_id` - The unique ID of the model.
  - `model_slug` - The slug identifier for the model.
  - `model_provider` - The provider of the model.
  - `accelerators` - The GPU accelerators allocated for this model deployment. Each element contains:
    - `accelerator_slug` - The slug identifier for the GPU accelerator type.
    - `scale` - The number of accelerator units allocated.
    - `type` - The accelerator type.
* `created_at` - The date and time when the dedicated inference endpoint was created.
* `updated_at` - The date and time when the dedicated inference endpoint was last updated.
