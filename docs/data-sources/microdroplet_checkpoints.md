---
page_title: "DigitalOcean: digitalocean_microdroplet_checkpoints"
subcategory: "MicroDroplets"
---

# digitalocean_microdroplet_checkpoints

Fetches the list of checkpoints belonging to a [DigitalOcean MicroDroplet](https://docs.digitalocean.com/products/microdroplets/).

Checkpoints are read-only artifacts: DigitalOcean captures one automatically each time a MicroDroplet is paused, preserving the memory and disk state needed to resume. They cannot be created or deleted directly through the customer API, so this provider only exposes them as a data source.

## Example Usage

```hcl
resource "digitalocean_microdroplet" "example" {
  name   = "example-microdroplet"
  region = "nyc3"
  size   = "microdroplet-1"
  image  = "docker.io/library/nginx:latest"

  auto_pause {
    enabled      = true
    idle_timeout = "5m"
  }
}

data "digitalocean_microdroplet_checkpoints" "example" {
  microdroplet_id = digitalocean_microdroplet.example.id
}

output "checkpoint_ids" {
  value = [for c in data.digitalocean_microdroplet_checkpoints.example.checkpoints : c.id]
}
```

## Argument Reference

* `microdroplet_id` - (Required) ID of the MicroDroplet whose checkpoints should be listed.
* `filter` - (Optional) Repeatable client-side filter block:
  * `key` - (Required) Field to match. Valid keys include `id`, `name`, `status`.
  * `values` - (Required) List of values to match on `key`.
  * `all` - (Optional) Require every value to match. Defaults to `false`.
  * `match_by` - (Optional) `exact`, `re`, or `substring`. Defaults to `exact`.
* `sort` - (Optional) Repeatable client-side sort block:
  * `key` - (Required) Field to sort by (e.g. `created_at`).
  * `direction` - (Optional) `asc` (default) or `desc`.

## Attributes Reference

* `checkpoints` - A list of checkpoint records. Each entry contains:
  * `id` - Checkpoint ID.
  * `microdroplet_id` - ID of the parent MicroDroplet.
  * `name` - Checkpoint name.
  * `status` - Lifecycle status of the checkpoint (e.g. `CHECKPOINT_AVAILABLE`).
  * `memory_bytes` - Size of the persisted memory image, in bytes.
  * `disk_bytes` - Size of the persisted disk image, in bytes.
  * `created_at` - RFC3339 timestamp of when the checkpoint was created.
