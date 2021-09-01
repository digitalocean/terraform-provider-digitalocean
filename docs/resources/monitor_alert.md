---
page_title: "DigitalOcean: digital_monitor_alert"
---

# digitalocean_monitor

Provides a [DigitalOcean Monitoring](https://docs.digitalocean.com/reference/api/api-reference/#tag/Monitoring) resource.
Monitor alerts can be configured to alert about, e.g., disk or memory usage exceeding certain threshold, or traffic at certain
limits. Notifications can be sent to either an email address or a Slack channel.

### Basic Example

```hcl
# Create a new Web Droplet in the nyc2 region
resource "digitalocean_droplet" "web" {
  image  = "ubuntu-20-04-x64"
  name   = "web-1"
  region = "nyc2"
  size   = "s-1vcpu-1gb"
}

resource "digitalocean_monitoring" "cpu_alert" {
  alerts      = {
    email = ["benny@digitalocean.com"]
    slack = {
      "channel"   = "Production Alerts",
      "url"       = "https://hooks.slack.com/services/T1234567/AAAAAAAA/ZZZZZZ"
  }
  window      = "5m"
  type        = "v1/insights/droplet/cpu"
  compare     = "GreaterThan"
  value       = 95
  enabled     = true
  entities    = [digitalocean_droplet.web.id]
  description = "Alert about CPU usage"
}
```

## Argument Reference

The following arguments are supported:

* `alerts` - (Required) How to send notifications about the alerts.
* `description` - (Required) The description of the alert.
* `compare` - (Required) The comparison for `value`. 
  This may be either `GreaterThan` or `LessThan`.
* `type` - (Required) The type of the alert.
  This may be either `v1/insights/droplet/load_1`, `v1/insights/droplet/load_5`, `v1/insights/droplet/load_15`,
  `v1/insights/droplet/memory_utilization_percent`, `v1/insights/droplet/disk_utilization_percent`,
  `v1/insights/droplet/cpu`, `v1/insights/droplet/disk_read`, `v1/insights/droplet/disk_write`,
  `v1/insights/droplet/public_outbound_bandwidth`, `v1/insights/droplet/public_inbound_bandwidth`,
  `v1/insights/droplet/private_outbound_bandwidth`, `v1/insights/droplet/private_inbound_bandwidth`.
* `enabled` - (Required) The status of the alert.
* `entities` - (Required) The resources to which the alert policy applies.
* `value` - (Required) The percentage to start alerting at, e.g., 90.
* `tags` - (Required) Tags for the alert.
* `window` - (Required) The time frame of the alert. Either 1m, 5m, 15m or 60m. 

## Attributes Reference

The following attributes are exported.

* `window` - The time frame of the alert.
* `enabled` - The status of the alert.
* `entities` - The resources for which the alert policy applies
* `type` - The type of the alert.
* `value` - The percentage to start alerting at.
* `tags` - Tags for the alert.
* `value` - The percentage to start alerting at, e.g., 90
* `alerts` - The notification policies of the alert policy.
* `description` - The description of the alert.
