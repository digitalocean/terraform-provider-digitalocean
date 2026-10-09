---
page_title: "DigitalOcean: digitalocean_microvm_checkpoint"
subcategory: "MicroVMs"
---

# digitalocean_microvm_checkpoint

Provides a DigitalOcean MicroVM checkpoint resource. A checkpoint captures the
memory and disk state of a running MicroVM, and can be used as the source of a
new [`digitalocean_microvm`](microvm).

-> **Note:** MicroVMs are in public preview.

## Example Usage

```hcl
resource "digitalocean_microvm" "web" {
  name   = "example-microvm"
  region = "nyc3"

  size {
    cpu    = 2
    memory = 4096
  }

  source {
    oci_ref = "docker.io/library/nginx:latest"
  }
}

resource "digitalocean_microvm_checkpoint" "web" {
  microvm_id = digitalocean_microvm.web.id
  name       = "web-checkpoint"
}
```

## Argument Reference

The following arguments are supported:

* `microvm_id` - (Required) The ID of the running MicroVM to checkpoint. Changing this forces a new checkpoint.
* `name` - (Optional) A name for the checkpoint. Changing this forces a new checkpoint.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The ID of the checkpoint.
* `microvm_name` - The name of the MicroVM the checkpoint was captured from.
* `region` - The slug of the region the checkpoint is stored in.
* `status` - The status of the checkpoint, such as `CHECKPOINT_AVAILABLE`.
* `size` - The size a MicroVM restored from this checkpoint inherits. Empty when the checkpoint does not record one, in which case the restored MicroVM must set `size`.
  - `cpu` - The number of vCPUs.
  - `memory` - The amount of memory in MiB.
  - `disk` - The size of the disk in GB.
* `memory_bytes` - The size of the persisted memory image in bytes.
* `disk_bytes` - The size of the persisted disk image in bytes.
* `created_at` - The date and time the checkpoint was created.

## Import

A MicroVM checkpoint can be imported using its `id`, e.g.

```
terraform import digitalocean_microvm_checkpoint.web 506f78a4-e098-11e5-ad9f-000f53306ae1
```
