---
page_title: "DigitalOcean: digitalocean_monitor_alert"
subcategory: "Monitoring"
---

# digitalocean_monitor_alert

Provides a [DigitalOcean Monitoring](https://docs.digitalocean.com/reference/api/api-reference/#tag/Monitoring)
resource. Monitor alerts can be configured to alert about, e.g., disk or memory
usage exceeding a certain threshold or traffic at a certain limit. Notifications
can be sent to either an email address or a Slack channel.

-> **Note** Currently, the [DigitalOcean API](https://docs.digitalocean.com/reference/api/api-reference/#operation/create_alert_policy) only supports creating alerts for Droplets.

### Basic Example

```hcl
# Create a new Web Droplet in the nyc2 region
resource "digitalocean_droplet" "web" {
  image  = "ubuntu-20-04-x64"
  name   = "web-1"
  region = "nyc2"
  size   = "s-1vcpu-1gb"
}

resource "digitalocean_monitor_alert" "cpu_alert" {
  alerts {
    email = ["sammy@digitalocean.com"]
    slack {
      channel = "Production Alerts"
      url     = "https://hooks.slack.com/services/T1234567/AAAAAAAA/ZZZZZZ"
    }
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

* `alerts` - (Required) How to send notifications about the alerts. This is a list with one element, .
  Note that for Slack, the DigitalOcean app needs to have permissions for your workspace. You can
  read more in [Slack's documentation](https://slack.com/intl/en-dk/help/articles/222386767-Manage-app-installation-settings-for-your-workspace)
* `description` - (Required) The description of the alert.
* `compare` - (Required) The comparison for `value`.
  This may be either `GreaterThan` or `LessThan`.
* `type` - (Required) The type of the alert.
  This may be one of `v1/insights/droplet/load_1`, `v1/insights/droplet/load_5`, `v1/insights/droplet/load_15`,
  `v1/insights/droplet/memory_utilization_percent`, `v1/insights/droplet/disk_utilization_percent`,
  `v1/insights/droplet/cpu`, `v1/insights/droplet/disk_read`, `v1/insights/droplet/disk_write`,
  `v1/insights/droplet/public_outbound_bandwidth`, `v1/insights/droplet/public_inbound_bandwidth`,
  `v1/insights/droplet/private_outbound_bandwidth`, `v1/insights/droplet/private_inbound_bandwidth`,
  `v1/insights/lbaas/avg_cpu_utilization_percent`, `v1/insights/lbaas/connection_utilization_percent`,
  `v1/insights/lbaas/droplet_health`, `v1/insights/lbaas/tls_connections_per_second_utilization_percent`,
  `v1/insights/lbaas/increase_in_http_error_rate_percentage_5xx`, `v1/insights/lbaas/increase_in_http_error_rate_percentage_4xx`,
  `v1/insights/lbaas/increase_in_http_error_rate_count_5xx`, `v1/insights/lbaas/increase_in_http_error_rate_count_4xx`,
  `v1/insights/lbaas/high_http_request_response_time`, `v1/insights/lbaas/high_http_request_response_time_50p`,
  `v1/insights/lbaas/high_http_request_response_time_95p`, `v1/insights/lbaas/high_http_request_response_time_99p`,
  `v1/dbaas/alerts/load_15_alerts`, `v1/dbaas/alerts/cpu_alerts`, `v1/dbaas/alerts/memory_utilization_alerts`, or
  `v1/dbaas/alerts/disk_utilization_alerts`.
* `enabled` - (Required) The status of the alert.
* `entities` - A list of IDs for the resources to which the alert policy applies.
* `tags` - A list of tags. When an included tag is added to a resource, the alert policy will apply to it.
* `value` - (Required) The value to start alerting at, e.g., 90% or 85Mbps. This is a floating-point number.
  DigitalOcean will show the correct unit in the web panel.
* `window` - (Required) The time frame of the alert. Either `5m`, `10m`, `30m`, or `1h`.

## Attributes Reference

The following attributes are exported.

* `uuid` - The uuid of the alert.
* `window` - The time frame of the alert.
* `enabled` - The status of the alert.
* `entities` - The resources for which the alert policy applies
* `type` - The type of the alert.
* `value` - The percentage to start alerting at.
* `tags` - Tags for the alert.
* `value` - The percentage to start alerting at, e.g., 90.
* `alerts` - The notification policies of the alert policy.
* `description` - The description of the alert.

## Import

Monitor alerts can be imported using the monitor alert `uuid`, e.g.

```shell
terraform import digitalocean_monitor_alert.cpu_alert b8ecd2ab-2267-4a5e-8692-cbf1d32583e3
```
