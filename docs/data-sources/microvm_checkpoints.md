---
page_title: "DigitalOcean: digitalocean_microvm_checkpoints"
subcategory: "MicroVMs"
---

# digitalocean_microvm_checkpoints

Get information on MicroVM checkpoints for use in other resources, with the
ability to filter and sort the results. If no arguments are specified, all
checkpoints on the account are returned.

-> **Note:** MicroVMs are in public preview.

## Example Usage

Get the checkpoints captured from a MicroVM and restore the newest one:

```hcl
data "digitalocean_microvm_checkpoints" "web" {
  microvm_id = digitalocean_microvm.web.id

  filter {
    key    = "status"
    values = ["CHECKPOINT_AVAILABLE"]
  }

  sort {
    key       = "created_at"
    direction = "desc"
  }
}

resource "digitalocean_microvm" "restored" {
  name = "example-restored"

  source {
    checkpoint_id = data.digitalocean_microvm_checkpoints.web.checkpoints[0].id
  }
}
```

## Argument Reference

* `microvm_id` - (Optional) Only return checkpoints captured from this MicroVM.
* `filter` - (Optional) Filter the results.
  The `filter` block is documented below.
* `sort` - (Optional) Sort the results.
  The `sort` block is documented below.

`filter` supports the following arguments:

* `key` - (Required) Filter the checkpoints by this key. This may be one of
  `id`, `name`, `microvm_id`, `microvm_name`, `region`, `status`,
  `memory_bytes`, `disk_bytes`, or `created_at`.
* `values` - (Required) A list of values to match against the `key` field.
* `match_by` - (Optional) One of `exact` (default), `re`, or `substring`.
* `all` - (Optional) Set to `true` to require that a field match all of the
  `values` instead of just one or more of them.

`sort` supports the following arguments:

* `key` - (Required) Sort the checkpoints by this key. This may be one of `id`,
  `name`, `microvm_id`, `microvm_name`, `region`, `status`, `memory_bytes`,
  `disk_bytes`, or `created_at`.
* `direction` - (Required) The sort direction. This may be either `asc` or `desc`.

## Attributes Reference

* `checkpoints` - A list of checkpoints satisfying any `filter` and `sort` criteria. Each checkpoint has the following attributes:
  - `id` - The ID of the checkpoint.
  - `name` - The name of the checkpoint.
  - `microvm_id` - The ID of the MicroVM the checkpoint was captured from.
  - `microvm_name` - The name of the MicroVM the checkpoint was captured from.
  - `region` - The slug of the region the checkpoint is stored in.
  - `status` - The status of the checkpoint, such as `CHECKPOINT_AVAILABLE`.
  - `size` - The size a MicroVM restored from this checkpoint inherits (`cpu`, `memory`, `disk`). Empty when the checkpoint does not record one.
  - `memory_bytes` - The size of the persisted memory image in bytes.
  - `disk_bytes` - The size of the persisted disk image in bytes.
  - `created_at` - The date and time the checkpoint was created.
