---
page_title: "DigitalOcean: digitalocean_uptime_check"
---

# digitalocean_uptime_check

Provides a [DigitalOcean Uptime Checks](https://docs.digitalocean.com/reference/api/api-reference/#tag/Uptime)
resource. Uptime Checks provide the ability to monitor your endpoints from around the world, and alert you when they're slow, unavailable, or SSL certificates are expiring.


### Basic Example

```hcl
# Create a new check for the target endpoint in a specifc region
resource "digitalocean_uptime_check" "foobar" {
  name    = "example-europe-check"
  target  = "https://www.example.com"
  regions = ["eu_west"]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A human-friendly display name for the check.
* `target` - (Required) The endpoint to perform healthchecks on.
* `type` - The type of health check to perform: 'ping' 'http' 'https'.
* `regions` - An array containing the selected regions to perform healthchecks from: "us_east", "us_west", "eu_west", "se_asia"
* `enabled` - A boolean value indicating whether the check is enabled/disabled.

## Attributes Reference

The following attributes are exported.

* `id` - The id of the check.

## Import

Uptime checks can be imported using the uptime check's `id`, e.g.

```shell
terraform import digitalocean_uptime_check.target 5a4981aa-9653-4bd1-bef5-d6bff52042e4
```
