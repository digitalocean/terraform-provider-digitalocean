---
page_title: "DigitalOcean: digitalocean_functions_trigger"
subcategory: "Functions"
---

# digitalocean_functions_trigger

Provides a DigitalOcean Functions trigger resource. A trigger invokes a function
in a namespace on a schedule. The function must already be deployed to the
namespace, for example with `doctl serverless deploy`.

## Example Usage

```hcl
resource "digitalocean_functions_namespace" "example" {
  label  = "my-functions"
  region = "nyc1"
}

resource "digitalocean_functions_trigger" "nightly_report" {
  namespace_id = digitalocean_functions_namespace.example.id
  name         = "nightly-report"
  function     = "reports/generate"

  scheduled_details {
    cron = "0 3 * * *"
    body = jsonencode({
      format = "pdf"
    })
  }
}
```

## Argument Reference

The following arguments are supported:

* `namespace_id` - (Required) The ID of the namespace the trigger belongs to. Changing this creates a new trigger.
* `name` - (Required) The trigger's name, unique within the namespace. Changing this creates a new trigger.
* `function` - (Required) The function the trigger invokes, including its package if it has one, e.g. `sample/hello`. Changing this creates a new trigger.
* `type` - (Optional) The trigger type. Only `SCHEDULED` is currently supported, which is the default.
* `enabled` - (Optional) Whether the trigger is enabled. Defaults to `true`.
* `scheduled_details` - (Required) The schedule for the trigger. The block supports:
  - `cron` - (Required) A standard cron expression for when the trigger fires, evaluated in UTC.
  - `body` - (Optional) A JSON object, as a string, passed to the function as its parameters each time the trigger fires. Use `jsonencode` to build it. Differences in whitespace or key order are ignored. Removing the body creates a new trigger.

## Attributes Reference

In addition to the above arguments, the following attributes are exported:

* `id` - The ID of the trigger, in the form `<namespace_id>/<name>`.
* `created_at` - The time the trigger was created.
* `updated_at` - The time the trigger was last updated.
* `last_run_at` - The last time the trigger fired, if it has fired.
* `next_run_at` - The next time the trigger is scheduled to fire.

## Import

Functions triggers can be imported using the namespace ID and the trigger name
joined with a comma, e.g.

```
terraform import digitalocean_functions_trigger.nightly_report fn-b90faf52-2b42-49c2-9792-75edfbb6f397,nightly-report
```
