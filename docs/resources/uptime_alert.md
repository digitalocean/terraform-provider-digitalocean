---
page_title: "DigitalOcean: digitalocean_uptime_alert"
subcategory: "Monitoring"
---

# digitalocean_uptime_alert

Provides a [DigitalOcean Uptime Alerts](https://docs.digitalocean.com/reference/api/api-reference/#operation/uptime_alert_create)
resource. Uptime Alerts provide the ability to add alerts to your [DigitalOcean Uptime Checks](https://docs.digitalocean.com/reference/api/api-reference/#tag/Uptime) when your endpoints are slow, unavailable, or SSL certificates are expiring.


### Basic Example

```hcl
# Create a new check for the target endpoint in a specific region
resource "digitalocean_uptime_check" "foobar" {
  name    = "example-europe-check"
  target  = "https://www.example.com"
  regions = ["eu_west"]
}

# Create a latency alert for the uptime check
resource "digitalocean_uptime_alert" "alert-example" {
  name       = "latency-alert"
  check_id   = digitalocean_uptime_check.foobar.id
  type       = "latency"
  threshold  = 300
  comparison = "greater_than"
  period     = "2m"
  notifications {
    email = ["sammy@digitalocean.com"]
    slack {
      channel = "Production Alerts"
      url     = "https://hooks.slack.com/services/T1234567/AAAAAAAA/ZZZZZZ"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `check_id` - (Required) A unique identifier for a check
* `name` - (Required) A human-friendly display name.
* `notifications` (Required) - The notification settings for a trigger alert.
* `type` (Required) - The type of health check to perform. Must be one of `latency`, `down`, `down_global` or `ssl_expiry`.
* `threshold` - The threshold at which the alert will enter a trigger state. The specific threshold is dependent on the alert type.
* `comparison` - The comparison operator used against the alert's threshold. Must be one of `greater_than` or `less_than`.
* `period` - Period of time the threshold must be exceeded to trigger the alert. Must be one of `2m`, `3m`, `5m`, `10m`, `15m`, `30m` or `1h`.

`notifications` supports the following:

* `email` - List of email addresses to sent notifications to.
* `slack`
  * `channel` (Required) - The Slack channel to send alerts to.
  * `url` (Required) - The webhook URL for Slack.

## Attributes Reference

The following attributes are exported.

* `id` - The id of the alert.

## Import

Uptime alerts can be imported using both the ID of the alert's parent check and
its own separated by a comma in the format: `check_id,alert_id`. For example:

```shell
terraform import digitalocean_uptime_alert.target 94a7d216-d821-11ee-a327-33d3239ffc4b,5a4981aa-9653-4bd1-bef5-d6bff52042e4
```
