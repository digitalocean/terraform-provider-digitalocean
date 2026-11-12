---
page_title: "DigitalOcean: digitalocean_functions_trigger"
subcategory: "Functions"
---

# digitalocean_functions_trigger

Get information about a DigitalOcean Functions trigger.

## Example Usage

```hcl
data "digitalocean_functions_trigger" "example" {
  namespace_id = "fn-b90faf52-2b42-49c2-9792-75edfbb6f397"
  name         = "nightly-report"
}

output "next_run" {
  value = data.digitalocean_functions_trigger.example.next_run_at
}
```

## Argument Reference

The following arguments are supported:

* `namespace_id` - (Required) The ID of the namespace the trigger belongs to.
* `name` - (Required) The trigger's name.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the trigger, in the form `<namespace_id>/<name>`.
* `function` - The function the trigger invokes.
* `type` - The trigger type.
* `enabled` - Whether the trigger is enabled.
* `scheduled_details` - The trigger's schedule:
  - `cron` - The cron expression for when the trigger fires.
  - `body` - The JSON object passed to the function, as a string.
* `created_at` - The time the trigger was created.
* `updated_at` - The time the trigger was last updated.
* `last_run_at` - The last time the trigger fired, if it has fired.
* `next_run_at` - The next time the trigger is scheduled to fire.
