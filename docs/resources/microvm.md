---
page_title: "DigitalOcean: digitalocean_microvm"
subcategory: "MicroVMs"
---

# digitalocean_microvm

Provides a DigitalOcean MicroVM resource. MicroVMs are lightweight sandboxes
that run an OCI container image, pause when idle, and resume on demand.

-> **Note:** MicroVMs are in public preview.

## Example Usage

Create a MicroVM from an OCI image:

```hcl
resource "digitalocean_microvm" "web" {
  name      = "example-microvm"
  region    = "nyc3"
  http_port = 8080

  size {
    cpu    = 2
    memory = 4096
  }

  source {
    oci_ref = "docker.io/library/nginx:latest"
  }

  auto_pause {
    enabled      = true
    idle_timeout = "5m"
  }

  environment = {
    LOG_LEVEL = "info"
  }

  tags = ["web"]
}

output "url" {
  value = digitalocean_microvm.web.urls[0].hostname
}
```

### Pause and Resume

Set `state` to `paused` to pause a running MicroVM, and back to `running` to
resume it. This is the only change that is applied in place; every other
argument forces a new MicroVM.

```hcl
resource "digitalocean_microvm" "worker" {
  name   = "example-worker"
  region = "nyc3"
  state  = "paused"

  size {
    cpu    = 1
    memory = 2048
  }

  source {
    oci_ref = "docker.io/library/alpine:latest"
  }
}
```

### Restore from a Checkpoint

When `source.checkpoint_id` is set, `region` and `size` may be omitted and are
inherited from the checkpoint.

```hcl
resource "digitalocean_microvm_checkpoint" "snapshot" {
  microvm_id = digitalocean_microvm.web.id
  name       = "web-checkpoint"
}

resource "digitalocean_microvm" "restored" {
  name = "example-restored"

  source {
    checkpoint_id = digitalocean_microvm_checkpoint.snapshot.id
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the MicroVM.
* `source` - (Required) What the MicroVM runs. Exactly one of the following must be set:
  - `oci_ref` - (Optional) The OCI reference of the workload container image, such as `docker.io/library/nginx:latest`.
  - `checkpoint_id` - (Optional) The ID of a checkpoint to restore.
* `region` - (Optional) The slug of the region to create the MicroVM in. Required when creating from `source.oci_ref`; inherited from the checkpoint when omitted on a restore.
* `size` - (Optional) The compute size of the MicroVM. Required when creating from `source.oci_ref`; inherited from the checkpoint when omitted on a restore.
  - `cpu` - (Required) The number of vCPUs.
  - `memory` - (Required) The amount of memory in MiB.
* `networking` - (Optional) The networking mode: `public` or `vpc`.
* `vpc_uuid` - (Optional) The ID of the VPC to place the MicroVM in. Only valid when `networking` is `vpc`.
* `http_port` - (Optional) The port the workload serves HTTP on.
* `http_protocol` - (Optional) The HTTP protocol the workload serves: `http` or `http2`.
* `ports` - (Optional) A list of guest ports to open for ingress. Defaults to `http_port` only.
* `environment` - (Optional) A map of environment variables to pass to the workload.
* `auto_pause` - (Optional) Auto-pause configuration. When omitted, the product default (enabled) applies.
  - `enabled` - (Required) Whether the MicroVM pauses after it has been idle for `idle_timeout`.
  - `idle_timeout` - (Optional) How long the MicroVM must be idle before it pauses, as a duration such as `5m` or `30s`.
* `auto_resume` - (Optional) Whether the MicroVM resumes automatically on incoming HTTP traffic. When omitted, the product default (enabled) applies.
* `tags` - (Optional) A list of tag names to apply to the MicroVM.
* `state` - (Optional) The desired lifecycle state: `running` or `paused`. Defaults to `running`. When `auto_pause` is enabled, a MicroVM that the platform has paused is not resumed to match `running`.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The ID of the MicroVM.
* `urn` - The uniform resource name (URN) for the MicroVM.
* `current_state` - The observed lifecycle state of the MicroVM, such as `running` or `paused`.
* `failure_reason` - A description of the failure when `current_state` is `failed`.
* `size.0.disk` - The size of the attached disk in GB.
* `urls` - A list of ingress URLs for the MicroVM:
  - `hostname` - The hostname, without a scheme.
  - `port` - The guest port the URL forwards to.
  - `default` - Whether this is the default URL.
  - `status` - The URL status: `PENDING` or `ACTIVE`.
* `created_at` - The date and time the MicroVM was created.

## Import

A MicroVM can be imported using its `id`, e.g.

```
terraform import digitalocean_microvm.web 506f78a4-e098-11e5-ad9f-000f53306ae1
```
