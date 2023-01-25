---
page_title: "DigitalOcean: digitalocean_uptime_alert"
---

# digitalocean_uptime_alert

Provides a [DigitalOcean Uptime Alerts](https://docs.digitalocean.com/reference/api/kafka-beta-api-reference/#operation/uptime_alert_create)
resource. Uptime Alerts provide the ability to add alerts to your [DigitalOcean Uptime Checks](https://docs.digitalocean.com/reference/api/kafka-beta-api-reference/#tag/Uptime) when your endpoints are slow, unavailable, or SSL certificates are expiring. 


### Basic Example

```hcl
# Create a new check for the target endpoint in a specifc region
resource "digitalocean_uptime_check" "foobar" {
  name  = "example-europe-check"
  target = "https://www.example.com"
  regions = ["eu_west"]
}


resource "digitalocean_uptime_alert" "alert-example" {
  name  = "latency-alert"
  check_id = "${digitalocean_uptime_check.foobar.id}"
  type = "latency"
	threshold = 300
	comparison = "greater_than"
  period = "2m"
  notifications {
    email = ["sammy@digitalocean.com"]
    slack {
      channel   = "Production Alerts"
      url       = "https://hooks.slack.com/services/T1234567/AAAAAAAA/ZZZZZZ"
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `check_id` - (Required) A unique identifier for a check
* `name` - (Required) A human-friendly display name.
* `notifications` (Required) - The notification settings for a trigger alert.
* `type` - The type of health check to perform: 'ping' 'http' 'https'.
* `threshold` - The comparison operator used against the alert's threshold: "greater_than", "less_than"
* `comparison` - A boolean value indicating whether the check is enabled/disabled.
* `period` - Period of time the threshold must be exceeded to trigger the alert: "2m" "3m" "5m" "10m" "15m" "30m" "1h"

## Attributes Reference

The following attributes are exported.

* `id` - The id of the alert.

## Import

Uptime checks can be imported using the uptime alert's `id`, e.g.

```shell
terraform import digitalocean_uptime_alert.target 5a4981aa-9653-4bd1-bef5-d6bff52042e4
```
