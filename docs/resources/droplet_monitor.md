---
page_title: "DigitalOcean: digital_monitor_alert"
---

# digitalocean_monitor

Provides a [DigitalOcean Monitoring](https://docs.digitalocean.com/reference/api/api-reference/#tag/Monitoring) resource.
Monitor alerts can be configured to alert about, e.g., disk or memory usage exceeding certain threshold. Notifications
can be sent to either an email address or a Slack channel.

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
* `description` - (Required) .
* `region` - (Required) The region to start in.
* `compare` - (Required)
* `enabled` - (Required) The status of the alert.
* `entities` - (Required) The resources to which the alert policy applies.
* `value` - (Required) The value of 

## Attributes Reference

