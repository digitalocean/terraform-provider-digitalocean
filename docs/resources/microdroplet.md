---
page_title: "DigitalOcean: digitalocean_microdroplet"
subcategory: "MicroDroplets"
---

# digitalocean_microdroplet

Provides a [DigitalOcean MicroDroplet](https://docs.digitalocean.com/products/microdroplets/) resource. MicroDroplets run OCI images and can auto-pause when idle to reduce cost.

## Example Usage

```hcl
resource "digitalocean_microdroplet" "example" {
  name   = "example-microdroplet"
  region = "nyc3"
  size   = "microdroplet-1"
  image  = "docker.io/library/nginx:latest"

  http_port     = 8080
  http_protocol = "http"

  auto_pause {
    enabled      = true
    idle_timeout = "5m"
  }

  auto_resume = true

  environment = {
    LOG_LEVEL = "info"
  }

  tags = ["web", "prod"]
}
```

### Pausing and resuming

The `state` attribute is settable and defaults to `running`. Changing it to `paused` calls the MicroDroplet pause action endpoint; changing it back to `running` calls the resume action endpoint. Neither triggers recreation:

```hcl
resource "digitalocean_microdroplet" "example" {
  # ...
  state = "paused"
}
```

Setting `state = "running"` on a paused MicroDroplet resumes it. When `auto_pause` is enabled, the platform may transition the MicroDroplet to `paused` on its own — Terraform suppresses the diff in that case so it does not fight the platform.

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the MicroDroplet. Forces recreation on change.
* `region` - (Required) DigitalOcean region slug for the MicroDroplet's location. Forces recreation on change.
* `size` - (Required) MicroDroplet size slug. Forces recreation on change.
* `image` - (Required) MicroDroplet image reference (UUID or URN). Forces recreation on change.
* `networking` - (Optional) Networking mode: `public` (default) or `vpc`. Forces recreation on change.
* `vpc_uuid` - (Optional) UUID of a VPC to attach to. Only valid when `networking = "vpc"`. Forces recreation on change.
* `http_port` - (Optional) HTTP port to expose. Forces recreation on change.
* `http_protocol` - (Optional) HTTP protocol: `http` or `http2`. Forces recreation on change.
* `environment` - (Optional) Map of environment variables passed to the MicroDroplet at boot. Forces recreation on change.
* `auto_pause` - (Optional) Auto-pause configuration block. Forces recreation on change: the MicroDroplets API has no in-place update path for auto_pause.
  * `enabled` - (Required) Whether auto-pause is enabled.
  * `idle_timeout` - (Optional) Idle timeout as a Go duration string (e.g. `5m`, `30s`).
* `auto_resume` - (Optional) Whether the MicroDroplet auto-resumes when it receives a request. Forces recreation on change: the MicroDroplets API has no in-place update path for auto_resume.
* `tags` - (Optional) Set of tags applied to the MicroDroplet at create time. Forces recreation on change: MicroDroplets are not exposed via the shared `/v2/tags` API, so tags cannot be added, removed, or verified after creation. They also do not populate on `terraform import` or on `digitalocean_microdroplet` / `digitalocean_microdroplets` data sources.
* `state` - (Optional) Desired lifecycle state. One of `running` (default) or `paused`. Changes are applied by invoking the pause / resume action endpoints and do not recreate the resource. This is the only attribute the provider can mutate in place.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The MicroDroplet UUID.
* `urn` - The uniform resource name (URN) for the MicroDroplet.
* `endpoint` - Public endpoint URL for the MicroDroplet.
* `current_state` - Observed lifecycle state of the MicroDroplet (may differ transiently from `state` while a transition is in progress).
* `created_at` - RFC3339 timestamp of when the MicroDroplet was created.

## Import

A MicroDroplet can be imported using its `id`, e.g.

```
terraform import digitalocean_microdroplet.example 506f78a4-e098-11e5-ad9f-000f53306ae1
```
